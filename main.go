package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/BenjaminB64/discordgoblindtest/discord"
	"github.com/BenjaminB64/discordgoblindtest/utils"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var loc *time.Location

func main() {
	loc, _ = time.LoadLocation("Europe/Paris")
	// Env + logs
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	if utils.IsDevelopment() {
		logrus.SetLevel(logrus.DebugLevel)
	}
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	ctx, ctxCcl := context.WithCancel(context.Background())

	// Discord
	d, err := discord.InitDiscord(os.Getenv("DISCORD_TOKEN"), os.Getenv("DISCORD_CHANNEL_ID"), utils.IsDevelopment(), ctx, loc)
	if err != nil {
		logrus.Fatal(err)
	}
	interruptChannel := make(chan os.Signal)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	logrus.Info("Bot running...")
	err = d.SendDebugMessage("Hello !")
	if err != nil {
		logrus.WithError(err).Error("Error sending hello")
	}

	wg := sync.WaitGroup{}
	// Loop
forLabel:
	for {
		select {
		case intr := <-interruptChannel:
			logrus.Debug(intr)
			d.SendDebugMessage("Good bye.")
			ctxCcl()
			wg.Wait()
			d.Close()
			logrus.Info("Closed")
			break forLabel
		}
	}

}
