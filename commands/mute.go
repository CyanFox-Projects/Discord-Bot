package commands

import (
	"strconv"
	"strings"
	"time"

	dbot "github.com/CyanFox-Projects/DiscordBot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/json"
)

var mute = discord.SlashCommandCreate{
	Name:        "mute",
	Description: "Mute a user",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "User to mute",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "time",
			Description: "Time to mute the user for (Format: 1d 2h 3m 4s)",
			Required:    true,
		},
	},
}

func MuteHandler(b *dbot.Bot) handler.CommandHandler {
	return func(e *handler.CommandEvent) error {

		if !e.Client().Caches().MemberPermissions(e.Member().Member).Has(discord.PermissionKickMembers) {
			return e.CreateMessage(dbot.SendNoPermissionMessage(*b))
		}

		user := e.SlashCommandInteractionData().User("user")
		muteTime := e.SlashCommandInteractionData().String("time")

		duration, err := stringToDuration(b, muteTime)
		if err != nil {
			b.Logger.Error("failed to parse duration: %v", err)
			return err
		}

		update := json.NewNullablePtr(time.Now().Add(duration))

		_, err = e.Client().Rest().UpdateMember(*e.GuildID(), user.ID, discord.MemberUpdate{CommunicationDisabledUntil: update})
		if err != nil {
			b.Logger.Error("failed to mute user: %v", err)
			return err
		}

		embed := discord.NewEmbedBuilder()
		embed.SetTitle("Mute")
		embed.SetColor(int(dbot.RGBToDecimal(255, 153, 51)))
		embed.AddField("Issuer", e.Member().User.Mention(), false)
		embed.AddField("Target", user.Mention(), false)
		embed.AddField("Until", "<t:"+strconv.FormatInt(time.Now().Add(duration).Unix(), 10)+">", false)
		embed.SetFooter(b.Config.EmbedFooterText, b.Config.EmbedFooterURL)
		embed.SetTimestamp(time.Now())
		embed.Build()

		message := discord.MessageCreateBuilder{}
		message.SetEmbeds(embed.Build())

		return e.CreateMessage(message.Build())
	}
}

func stringToDuration(b *dbot.Bot, input string) (time.Duration, error) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		b.Logger.Error("empty duration")
		return 0, nil
	}

	multipliers := map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": 24 * time.Hour,
	}

	totalDuration := time.Duration(0)
	currentNumber := ""
	for _, char := range input {
		if char >= '0' && char <= '9' {
			currentNumber += string(char)
		} else if multiplier, exists := multipliers[string(char)]; exists {
			if currentNumber == "" {
				b.Logger.Error("no number before %v", char)
				return 0, nil
			}
			number, _ := strconv.Atoi(currentNumber)
			totalDuration += time.Duration(number) * multiplier
			currentNumber = ""
		} else if string(char) != " " {
			b.Logger.Error("unknown character %v", char)
			return 0, nil
		}
	}

	return totalDuration, nil
}
