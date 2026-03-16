package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
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
		log.Printf("Authorized on account @%s (expected username: @%s)", bot.Self.UserName, username)
	} else {
		log.Printf("Authorized on account @%s", bot.Self.UserName)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Graceful shutdown: stop on SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Println("Bot started with long polling. Press Ctrl+C to stop.")

	for {
		select {
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			handleMessage(bot, update.Message)

		case sig := <-stop:
			log.Printf("Received signal %s, shutting down...", sig)
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
		log.Printf("failed to send message: %v", err)
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
			log.Printf("failed to send /start reply: %v", err)
		}

	case "help":
		text := "Я бот уровня 1 (простые функции).\n\n" +
			"- Отвечаю на /start и /help\n" +
			"- Повторяю твои сообщения (echo)\n\n" +
			"Это шаблон, на основе которого можно делать коммерческие боты."
		reply := tgbotapi.NewMessage(chatID, text)
		reply.ParseMode = tgbotapi.ModeMarkdown
		if _, err := bot.Send(reply); err != nil {
			log.Printf("failed to send /help reply: %v", err)
		}

	default:
		reply := tgbotapi.NewMessage(chatID, "Неизвестная команда. Напиши /help, чтобы узнать, что я умею.")
		if _, err := bot.Send(reply); err != nil {
			log.Printf("failed to send unknown command reply: %v", err)
		}
	}
}
