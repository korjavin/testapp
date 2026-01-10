package main

import (
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// Flags
	token := flag.String("token", "", "Telegram bot token")
	flag.Parse()

	// Setup logging
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	// Get token from env or flag
	botToken := os.Getenv("TG_BOT_TOKEN")
	if *token != "" {
		botToken = *token
	}

	if botToken == "" {
		slog.Error("Bot token is required. Set TG_BOT_TOKEN env or use -token flag")
		os.Exit(1)
	}

	// Create bot
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		slog.Error("Failed to create bot", "error", err)
		os.Exit(1)
	}

	slog.Info("Bot started", "username", bot.Self.UserName)

	// Create update config
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Get updates channel
	updates := bot.GetUpdatesChan(u)

	// Handle updates in goroutine
	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			// Handle /start command
			if update.Message.IsCommand() && update.Message.Command() == "start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to testapp! Click the button below to open the WebApp:")
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("Open WebApp", "open_webapp"),
					),
				)

				if _, err := bot.Send(msg); err != nil {
					slog.Error("Failed to send message", "error", err)
				}
			}

			// Handle callback query
			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "open_webapp" {
					// Answer callback with notification
					callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Opening WebApp...")
					if _, err := bot.AnswerCallbackQuery(callback); err != nil {
						slog.Error("Failed to answer callback", "error", err)
					}

					webAppURL := "https://your-domain.com"
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Click to open: "+webAppURL)
					if _, err := bot.Send(msg); err != nil {
						slog.Error("Failed to send webapp URL", "error", err)
					}
				}
			}

			// Echo other messages
			if !update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You said: "+update.Message.Text)
				if _, err := bot.Send(msg); err != nil {
					slog.Error("Failed to send message", "error", err)
				}
			}
		}
	}()

	slog.Info("Bot is running. Press Ctrl+C to stop.")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Stopping bot...")
	bot.StopReceivingUpdates()
	slog.Info("Bot stopped")
}
