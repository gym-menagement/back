package router

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"

	"gym/global/log"
	"gym/models"
	"gym/services/withdraw"

	"github.com/gofiber/fiber/v2"
)

// 회원 탈퇴 (계정 삭제).
//
// 문제 1 — 생성된 `DELETE /api/user/:id`·`/api/user/batch` 핸들러는 요청자
// 신원을 보지 않는다. 로그인한 아무 회원이 남의 계정(관리자 포함)을 지울 수
// 있었다. OwnershipGuard 는 GET 만 본다.
// 문제 2 — 그 핸들러는 user_tb 행만 지워서 신체정보·출석·운동기록 등 개인
// 기록이 그대로 남았다 (개인정보 처리방침의 "탈퇴 시 삭제" 고지와 불일치).
//
// 생성 파일은 make server 로 덮어써지므로 다른 가드와 같은 전략으로
// 미들웨어가 요청을 가로채 직접 처리한다 (생성 핸들러까지 가지 않는다):
//   - DELETE /api/user/me     본인 탈퇴 (SetupWithdrawRoutes)
//   - DELETE /api/user/:id    본인 또는 플랫폼관리자
//   - DELETE /api/user/batch  플랫폼관리자 전용
//
// 어느 경로든 정리 로직은 services/withdraw 하나를 쓴다.

var userDeletePathRe = regexp.MustCompile(`^/api/user/(\d+|batch)$`)

// UserDeleteGuard 는 SetRouter(=JwtAuthRequired) 보다 먼저 등록되므로
// c.Locals("user") 가 아직 없다 — tokenUser 로 요청자를 직접 파싱한다.
var UserDeleteGuard = func(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodDelete {
		return c.Next()
	}
	m := userDeletePathRe.FindStringSubmatch(c.Path())
	if m == nil {
		return c.Next()
	}
	me := tokenUser(c)
	if me == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
	}

	if m[1] == "batch" {
		if !isPlatformAdmin(me) {
			return forbidUserDelete(c, me.Id, 0)
		}
		ids, ok := batchUserIDs(c.Body())
		if !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"code": "bad request"})
		}
		for _, id := range ids {
			if resp := withdrawUser(c, me, id); resp != nil {
				return resp(c)
			}
		}
		return c.JSON(fiber.Map{"code": "ok"})
	}

	id, _ := strconv.ParseInt(m[1], 10, 64)
	if id != me.Id && !isPlatformAdmin(me) {
		return forbidUserDelete(c, me.Id, id)
	}
	if resp := withdrawUser(c, me, id); resp != nil {
		return resp(c)
	}
	return c.JSON(fiber.Map{"code": "ok"})
}

// SetupWithdrawRoutes — 본인 탈퇴. id 를 받지 않고 토큰의 사용자만 지운다.
func SetupWithdrawRoutes(group fiber.Router) {
	group.Delete("/user/me", func(c *fiber.Ctx) error {
		me := currentUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": "unauthorized"})
		}
		if resp := withdrawUser(c, me, me.Id); resp != nil {
			return resp(c)
		}
		return c.JSON(fiber.Map{"code": "ok"})
	})
}

// withdrawUser 는 권한 판정이 끝난 대상을 탈퇴 처리한다. 성공이면 nil,
// 실패면 에러 응답을 쓰는 함수를 돌려준다.
func withdrawUser(c *fiber.Ctx, me *models.User, targetID int64) func(*fiber.Ctx) error {
	conn := models.NewConnection()
	if conn == nil {
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
	}
	defer conn.Close()

	// 헬스장 소유자가 스스로 탈퇴하면 헬스장(회원·결제·구독)이 주인 없이
	// 남는다. 본인 탈퇴는 막고 운영자 처리(플랫폼관리자 삭제)로 넘긴다.
	if targetID == me.Id {
		owner, err := withdraw.OwnsGym(conn, targetID)
		if err != nil {
			log.Error().Int64("target", targetID).Msg("withdraw: " + err.Error())
			return func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
			}
		}
		if owner {
			return func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"code":    "gym owner",
					"message": "헬스장 소유자 계정은 직접 탈퇴할 수 없습니다. 고객센터로 문의해 주세요.",
				})
			}
		}
	}

	err := withdraw.Withdraw(conn, targetID)
	if errors.Is(err, withdraw.ErrNotFound) {
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": "not found"})
		}
	}
	if err != nil {
		log.Error().Int64("actor", me.Id).Int64("target", targetID).Msg("withdraw: " + err.Error())
		return func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": "error"})
		}
	}
	log.Info().Int64("actor", me.Id).Int64("target", targetID).Msg("user withdrawn")
	return nil
}

// batchUserIDs 는 일괄 삭제 본문을 읽는다. 웹 프론트는 id 배열([1,2])을,
// 생성 핸들러 계약은 객체 배열([{"id":1}])을 쓰므로 둘 다 받는다.
func batchUserIDs(body []byte) ([]int64, bool) {
	var raw []json.RawMessage
	if json.Unmarshal(body, &raw) != nil {
		return nil, false
	}
	ids := make([]int64, 0, len(raw))
	for _, r := range raw {
		var id int64
		if json.Unmarshal(r, &id) != nil {
			var item struct {
				Id int64 `json:"id"`
			}
			if json.Unmarshal(r, &item) != nil {
				return nil, false
			}
			id = item.Id
		}
		if id <= 0 {
			return nil, false
		}
		ids = append(ids, id)
	}
	return ids, true
}

func forbidUserDelete(c *fiber.Ctx, actor, target int64) error {
	log.Info().Int64("actor", actor).Int64("target", target).Msg("user delete blocked")
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": "forbidden"})
}
