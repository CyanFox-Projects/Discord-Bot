package commands

import (
	"io/ioutil"
	"net/http"
	"time"

	dbot "github.com/CyanFox-Projects/DiscordBot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var embed = discord.SlashCommandCreate{
	Name:        "purge",
	Description: "Purge messages from a channel",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "type",
			Description: "Type of embed to send (info, rules)",
			Required:    true,
		},
	},
}

func EmbedHandler(b *dbot.Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {

		option := e.SlashCommandInteractionData().String("type")

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

		embed := discord.NewEmbedBuilder()
		resp, err := http.Get(b.Config.InfoURL)
		embed.SetTitle("Information's")

		if option == "rules" {
			embed.SetTitle("Rules")
			resp, err = http.Get(b.Config.RulesURL)
		}

		if err != nil {
			b.Logger.Error("failed to get content from url: %s", err.Error())
			return err
		}

		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			b.Logger.Error("failed to read content: %s", err.Error())
			return err
		}

		embed.SetDescription(string(content))
		embed.SetColor(int(dbot.RGBToDecimal(51, 204, 255)))
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
