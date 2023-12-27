package commands

import (
	"time"

	dbot "github.com/CyanFox-Projects/DiscordBot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"
)

var unmute = discord.SlashCommandCreate{
	Name:        "mute",
	Description: "Mute a user",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "User to mute",
			Required:    true,
		},
	},
}

func UnMuteHandler(b *dbot.Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {

		if !e.Client().Caches().MemberPermissions(e.Member().Member).Has(discord.PermissionKickMembers) {
			return e.CreateMessage(dbot.SendNoPermissionMessage(*b))
		}

		user := e.SlashCommandInteractionData().User("user")

		update := json.NewNullablePtr(time.Now().Add(1 * time.Second))

		_, err := e.Client().Rest().UpdateMember(*e.GuildID(), user.ID, discord.MemberUpdate{CommunicationDisabledUntil: update})
		if err != nil {
			b.Logger.Error("failed to mute user: %v", err)
			return err
		}

		embed := discord.NewEmbedBuilder()
		embed.SetTitle("Unmute")
		embed.SetColor(int(dbot.RGBToDecimal(51, 204, 51)))
		embed.AddField("Issuer", e.Member().User.Mention(), false)
		embed.AddField("Target", user.Mention(), false)
		embed.SetFooter(b.Config.EmbedFooterText, b.Config.EmbedFooterURL)
		embed.SetTimestamp(time.Now())
		embed.Build()

		message := discord.MessageCreateBuilder{}
		message.SetEmbeds(embed.Build())

		return e.CreateMessage(message.Build())
	}
}
