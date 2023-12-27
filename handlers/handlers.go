package handlers

import (
	"strconv"
	"time"

	dbot "github.com/CyanFox-Projects/DiscordBot"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func WelcomeHandler(b *dbot.Bot) bot.EventListener {
	return bot.NewListenerFunc(func(e *events.GuildMemberJoin) {

		b.Logger.Info("New member joined")
		err := e.Client().Rest().AddMemberRole(e.GuildID, e.Member.User.ID, b.Config.MemberRole)
		if err != nil {
			b.Logger.Error(err)
		}

		embed := discord.NewEmbedBuilder()
		embed.SetTitle("Welcome")
		embed.SetDescription("Welcome to the CyanFox-Projects Discord Server!")
		embed.SetColor(int(dbot.RGBToDecimal(51, 204, 255)))
		embed.SetThumbnail(*e.Member.User.AvatarURL())
		embed.AddField("User", e.Member.User.Mention(), false)
		embed.AddField("ID", strconv.FormatUint(uint64(e.Member.User.ID), 10), false)
		embed.SetFooter(b.Config.EmbedFooterText, b.Config.EmbedFooterURL)
		embed.SetTimestamp(time.Now())
		embed.Build()

		message := discord.MessageCreateBuilder{}
		message.SetEmbeds(embed.Build())

		_, err = b.Client.Rest().CreateMessage(b.Config.WelcomeChannel, message.Build())
		if err != nil {
			b.Logger.Error(err)
		}
	})
}
