package fav

import (
	"tgbot/app"

	"github.com/tucnak/telebot"
)

func List(m *telebot.Message) {
	msg := "list:" + m.Payload
	if m.FromGroup() {
		msg = "@" + m.Sender.Username + " " + msg
	}

	app.Bot.Send(m.Chat, msg)
}

func Add(m *telebot.Message) {
	msg := "add:" + m.Payload
	if m.FromGroup() {
		msg = "@" + m.Sender.Username + " " + msg
	}

	app.Bot.Send(m.Chat, msg)
}
