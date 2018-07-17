package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(APIKey)
	if err != nil {
		log.Panic(err)
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if m := update.Message; m.From.ID == 461755327 || (m.Chat.IsPrivate() && m.From.ID == 106568920) {
			go func() {
				if results := CheckText(update.Message.Text); results != nil {
					output := "God damnit, Felix, lern Rechtschreibung!!\n"
					for fail, correction := range results {
						output += "\n- " + fail + "\t→ " + correction
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, output)
					msg.ParseMode = tgbotapi.ModeMarkdown
					msg.ReplyToMessageID = update.Message.MessageID

					bot.Send(msg)
				}
			}()
		}
	}
}
