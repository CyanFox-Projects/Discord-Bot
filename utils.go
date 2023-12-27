package DiscordBot

import (
	"time"

	"github.com/disgoorg/disgo/discord"
)

func RGBToDecimal(r, g, b int) int64 {
	return int64(r)<<16 | int64(g)<<8 | int64(b)
}

func SendNoPermissionMessage(b Bot) discord.MessageCreate {
	embed := discord.NewEmbedBuilder()
	embed.SetTitle("No Permission")
	embed.SetDescription("You do not have permission to use this command")
	embed.SetColor(int(RGBToDecimal(255, 51, 0)))
	embed.SetFooter(b.Config.EmbedFooterText, b.Config.EmbedFooterURL)
	embed.SetTimestamp(time.Now())
	embed.Build()

	message := discord.MessageCreateBuilder{}
	message.SetEmbeds(embed.Build())

	return message.Build()
}
