package main

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

func main() {
	LoadDictionary()
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
		if m := update.Message; m.From.ID == 461755327 || (m.From.ID == 106568920 && m.Chat.IsPrivate() && !m.IsCommand()) {
			go func() {
				if results := CheckText(m.Text); results != nil {
					output := "God damnit, Felix, lern Rechtschreibung!!\n"
					for fail, correction := range results {
						output += "\n- " + fail + "\t→ " + correction
					}
					msg := tgbotapi.NewMessage(m.Chat.ID, output)
					msg.ParseMode = tgbotapi.ModeMarkdown
					msg.ReplyToMessageID = m.MessageID

					bot.Send(msg)
				}
			}()
		} else if m.From.ID == 106568920 && m.IsCommand() {
			switch m.Command() {
			case "add":
				var output string
				if added := AddToDictionary(m.CommandArguments()); added {
					output = "Erfolgreich hinzgefügt."
				} else {
					output = "Fehler beim Hinzufügen, ist das Wort bereits in der Whitelist?"
				}
				msg := tgbotapi.NewMessage(m.Chat.ID, output)
				bot.Send(msg)
				break
			case "remove":
				var output string
				if RemoveFromDictionary(m.CommandArguments()) {
					output = "Erfolgreich entfernt."
				} else {
					output = "Fehler beim Entfernen, ist das Wort überhaupt in der Whitelist?"
				}
				msg := tgbotapi.NewMessage(m.Chat.ID, output)
				bot.Send(msg)
				break
			case "whitelist":
				msg := tgbotapi.NewMessage(m.Chat.ID, QueryWhitelist())
				bot.Send(msg)
				break
			}
		}
	}
}
