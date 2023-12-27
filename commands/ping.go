package commands

import (
	dbot "github.com/CyanFox-Projects/DiscordBot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var ping = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Pong!",
}

func PingHandler(b *dbot.Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {

		message := discord.MessageCreateBuilder{}
		message.SetContent("Pong!")

		return e.CreateMessage(message.Build())
	}
}
