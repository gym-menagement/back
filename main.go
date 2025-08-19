package main

import (
	"fmt"
	"gym/global/config"
	"gym/global/log"
	"gym/global/setting"
	"gym/models"
	"gym/services"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.New(rand.NewSource(time.Now().UnixNano()))

	log.Info().Str("Version", config.Version).Str("Mode", config.Mode).Msg("Start")

	models.InitCache()

	tempPath := fmt.Sprintf("%v/temp", config.UploadPath)
	os.MkdirAll(tempPath, 777)
	os.Chmod(tempPath, os.FileMode(0755))

	setting.GetInstance()

	services.Cron()
	// services.Fcm()
	services.Chat()
	services.Notify()
	services.Http()
}
