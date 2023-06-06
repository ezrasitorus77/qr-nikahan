package helper

import (
	"qr-nikahan/internal/consts"
)

func CreateMessage(name string) (message string) {
	var prefix string = "Halo " + name

	return prefix + consts.MessageBody
}
