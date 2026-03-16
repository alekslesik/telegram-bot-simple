package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
}

var logger = newLogger()

func newLogger() *slog.Logger {
	// Уровень логирования: LOG_LEVEL=debug включает debug, всё остальное — info.
	level := slog.LevelInfo
	if strings.EqualFold(strings.TrimSpace(os.Getenv("LOG_LEVEL")), "debug") {
		level = slog.LevelDebug
	}

	// Формат времени: 16.06.2026 15:31:30
	timeReplacer := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			if t, ok := a.Value.Any().(time.Time); ok {
				a.Value = slog.StringValue(t.Format("02.01.2006 15:04:05"))
			}
		}
		return a
	}

	opts := &slog.HandlerOptions{
		Level:       level,
		ReplaceAttr: timeReplacer,
	}

	// Формат логов: LOG_FORMAT=json → JSON, иначе человекочитаемый текст.
	format := strings.ToLower(strings.TrimSpace(os.Getenv("LOG_FORMAT")))
	if format == "json" {
		return slog.New(slog.NewJSONHandler(os.Stdout, opts))
	}

	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

func main() {
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("env var TOKEN is not set (see .env)")
	}

	username := os.Getenv("USERNAME")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("failed to create bot: %v", err)
	}

	if username != "" {
		logger.Info("authorized",
			"username", bot.Self.UserName,
			"expected_username", username,
		)
	} else {
		logger.Info("authorized",
			"username", bot.Self.UserName,
		)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Graceful shutdown: stop on SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	logger.Info("bot started with long polling, press Ctrl+C to stop")

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			handleMessage(bot, update.Message)

		case sig := <-stop:
			logger.Info("received signal, shutting down", "signal", sig.String())
			return
		}
	}
}

func handleMessage(bot Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	if msg.IsCommand() {
		handleCommand(bot, msg)
		return
	}

	// Simple echo for non-command messages
	reply := tgbotapi.NewMessage(chatID, "Ты написал: "+msg.Text)
	if _, err := bot.Send(reply); err != nil {
		logger.Error("failed to send message", "err", err)
	}
}

func handleCommand(bot Bot, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	switch msg.Command() {
	case "start":
		text := "Привет! Я простой бот на Go.\n\n" +
			"Доступные команды:\n" +
			"/start - приветствие\n" +
			"/help - подсказка\n" +
			"Напиши любое сообщение — я его повторю."
		reply := tgbotapi.NewMessage(chatID, text)
		reply.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.Send(reply); err != nil {
			logger.Error("failed to send /start reply", "err", err)
		}

	case "help":
		text := "Я бот уровня 1 (простые функции).\n\n" +
			"- Отвечаю на /start и /help\n" +
			"- Повторяю твои сообщения (echo)\n\n" +
			"Это шаблон, на основе которого можно делать коммерческие боты."
		reply := tgbotapi.NewMessage(chatID, text)
		reply.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.Send(reply); err != nil {
			logger.Error("failed to send /help reply", "err", err)
		}

	default:
		reply := tgbotapi.NewMessage(chatID, "Неизвестная команда. Напиши /help, чтобы узнать, что я умею.")
		if _, err := bot.Send(reply); err != nil {
			logger.Error("failed to send unknown command reply", "err", err)
		}
	}
}
