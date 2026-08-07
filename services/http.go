package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"gym/global"
	"gym/global/config"
	"gym/global/log"
	"gym/global/setting"
	"gym/router"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/YiYuhki/ikisocket"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
)

func Http() {
	log.Info().Str("service", "HTTP").Msg("Start Service")

	app := fiber.New(fiber.Config{
		BodyLimit:             500 * 1024 * 1024,
		Prefork:               false,
		CaseSensitive:         true,
		StrictRouting:         true,
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
	})

	// 핸들러 패닉이 프로세스를 죽이지 않도록 가장 먼저 건다.
	// (인증 미들웨어의 슬라이스 범위 초과 패닉으로 서버가 내려갈 수 있었다)
	app.Use(recover.New())

	sites := strings.Join(config.Cors, ", ")
	if sites != "" {
		app.Use(cors.New(cors.Config{
			AllowOrigins: sites,
		}))
	}

	if config.Log.Web {
		app.Use(fiberzerolog.New(fiberzerolog.Config{
			Logger: log.Get(),
		}))
	}

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	if chat.Use {
		app.Use("/ws", func(c *fiber.Ctx) error {
			if websocket.IsWebSocketUpgrade(c) {
				return c.Next()
			}

			return fiber.ErrUpgradeRequired
		})

		app.Get("/ws/:id", ikisocket.New(func(kws *ikisocket.Websocket) {
			id := global.Atol(kws.Params("id"))
			kws.SetAttribute("id", id)

			session := ChatSession{
				Id:    id,
				UUID:  kws.UUID,
				Valid: true,
				Date:  time.Now(),
			}

			for k, v := range chat.Clients {
				if v.Id == id {
					v.Valid = false
					chat.Clients[k] = v
					// 기존 사용자 튕겨줘야 함

					instance := setting.GetInstance()
					if instance.Setting("user.login.single") == "Y" {
						global.SendNotify(id, v.UUID, global.SessionTimeout)
					}
				}
			}

			chat.Clients[kws.UUID] = session
		}))
	}

	app.Static("/", config.DocumentRoot)
	//app.Static("/webdata", config.UploadPath)

	// PATCH/PUT 부분 본문을 기존 값과 병합 (구 Spring 계약 복원 — patchmerge.go).
	// 라우트보다 먼저 걸어야 모든 생성 Update 핸들러에 적용된다.
	app.Use(router.PatchMerge)

	router.SetRouter(app)
	// 손코드 라우트(P1~P5). SetRouter 가 건 /api 미들웨어(JWT 등)를 그대로
	// 통과하도록 반드시 SetRouter 다음, SPA 폴백(/*) 이전에 등록한다.
	router.SetCustomRouter(app)

	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile(fmt.Sprintf("./%v/index.html", config.DocumentRoot), true)
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-c
		log.Info().Msg("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	if config.Mode == "develop" || !config.Tls.Use {
		if err := app.Listen(":" + config.Port); err != nil {
			log.Error().Msg(err.Error())
		}
	} else {
		cer, err := tls.LoadX509KeyPair(config.Tls.Cert, config.Tls.Key)
		if err != nil {
			log.Error().Msg("TLS error")
			log.Error().Msg(err.Error())
			return
		}

		cert := &tls.Config{Certificates: []tls.Certificate{cer}}

		ln, err := tls.Listen("tcp", ":"+config.Port, cert)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}

		if err := app.Listener(ln); err != nil {
			log.Error().Msg(err.Error())
		}
	}
}
