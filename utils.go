package DiscordBot

import "strconv"

func HexToColorInt(s string) int {
	colorInt, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		return 0
	}
	return int(colorInt)
}
