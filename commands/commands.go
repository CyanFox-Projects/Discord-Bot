package commands

import "github.com/disgoorg/disgo/discord"

var Commands = []discord.ApplicationCommandCreate{
	mute,
	unmute,
	purge,
	embed,
	ping,
}
