package commands

import (
	"log"
	"strconv"
	"time"

	dbot "github.com/CyanFox-Projects/DiscordBot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var purge = discord.SlashCommandCreate{
	Name:        "embed",
	Description: "Send an pre-made embed",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionInt{
			Name:        "amount",
			Description: "Amount of messages to purge",
			Required:    true,
		},
	},
}

func PurgeHandler(b *dbot.Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {

		roleId := b.Config.StaffRoleID
		roleExists := false

		for _, rID := range e.Member().RoleIDs {
			if rID == roleId {
				roleExists = true
				break
			}
		}

		if !roleExists {
			return e.CreateMessage(dbot.SendNoPermissionMessage(*b))
		}

		err := e.DeferCreateMessage(false)
		if err != nil {
			b.Logger.Error("failed to defer create message: %s", err.Error())
			return err
		}

		amount := e.SlashCommandInteractionData().Int("amount")

		if amount > 100 {
			amount = 99
		}

		if amount < 1 {
			amount = 1
		}

		messages, err := e.Client().Rest().GetMessages(e.ChannelID(), 0, 0, 0, amount)
		if err != nil {
			log.Println("Errors fetching messages:", err)
			return err
		}

		for _, message := range messages {
			err := e.Client().Rest().DeleteMessage(e.ChannelID(), message.ID)
			if err != nil {
				b.Logger.Error("failed to delete message: %s", err.Error())
				return err
			}
		}

		embed := discord.NewEmbedBuilder()
		embed.SetTitle("Purge")
		embed.SetColor(int(dbot.RGBToDecimal(255, 153, 0)))
		embed.AddField("Issuer", e.Member().User.Mention(), false)
		embed.AddField("Amount", strconv.Itoa(amount), false)
		embed.AddField("Channel", "<#"+strconv.FormatUint(uint64(e.ChannelID()), 10)+">", false)
		embed.SetFooter(b.Config.EmbedFooterText, b.Config.EmbedFooterURL)
		embed.SetTimestamp(time.Now())
		embed.Build()

		message := discord.MessageCreateBuilder{}
		message.SetEmbeds(embed.Build())

		_, err = b.Client.Rest().CreateMessage(b.Config.WelcomeChannel, message.Build())
		if err != nil {
			b.Logger.Error(err)
		}

		return nil
	}
}
