package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/alekslesik/telegram-bot-simple/internal/bot"
	"github.com/alekslesik/telegram-bot-simple/internal/logging"
	"github.com/alekslesik/telegram-bot-simple/internal/telegram"
)

var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func main() {
	logger := logging.NewFromEnv()

	buildDate := BuildDate
	if t, err := time.Parse(time.RFC3339, BuildDate); err == nil {
		if loc, locErr := time.LoadLocation("Europe/Moscow"); locErr == nil {
			t = t.In(loc)
		}
		buildDate = t.Format("02/01/2006 15:04:05")
	} else if t, err := time.Parse(time.RFC3339Nano, BuildDate); err == nil {
		if loc, locErr := time.LoadLocation("Europe/Moscow"); locErr == nil {
			t = t.In(loc)
		}
		buildDate = t.Format("02/01/2006 15:04:05")
	}

	logger.Info("starting",
		"version", Version,
		"commit", Commit,
		"build_date", buildDate,
	)

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("env var TOKEN is not set (see .env)")
	}

	username := os.Getenv("USERNAME")

	tg, err := telegram.New(token)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	if username != "" {
		logger.Info("authorized",
			"username", tg.Self.UserName,
			"expected_username", username,
		)
	} else {
		logger.Info("authorized",
			"username", tg.Self.UserName,
		)
	}

	commands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{Command: "start", Description: "🚀 Старт"},
		tgbotapi.BotCommand{Command: "menu", Description: "📋 Демо-меню"},
		tgbotapi.BotCommand{Command: "help", Description: "📋 Меню команд"},
		tgbotapi.BotCommand{Command: "about", Description: "ℹ️ О боте"},
		tgbotapi.BotCommand{Command: "usecases", Description: "💼 Примеры задач"},
		tgbotapi.BotCommand{Command: "features", Description: "🧩 Возможности"},
		tgbotapi.BotCommand{Command: "ping", Description: "✅ Проверка статуса"},
		tgbotapi.BotCommand{Command: "echo", Description: "🗣️ Повторить текст"},
	)
	if _, err := tg.Request(commands); err != nil {
		logger.Error("failed to register bot commands", "err", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tg.GetUpdatesChan(u)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	h := bot.Handlers{
		Bot:    tg,
		Logger: logger,
	}

	logger.Info("bot started with long polling, press Ctrl+C to stop")

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			h.HandleMessage(update.Message)

		case sig := <-stop:
			logger.Info("received signal, shutting down", "signal", sig.String())
			return
		}
	}
}
