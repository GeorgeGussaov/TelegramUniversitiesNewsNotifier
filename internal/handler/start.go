package handler

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStartChatCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI){
	log.Print(update.ChannelPost.Chat.ID)	//просто выводим айдишник чата, чтобы записать его вручную в env и в последующем пушить туда уведомления
}