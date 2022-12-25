package main

import (
	"fmt"
	"go-telebot/models"
	"image"
)

func imageHandler(message models.Message) {
	var err error
	defer func() {
		if err != nil {
			switch err.(type) {
			default:
				_ = sendMessage("Unexpected error", message.Chat.Id, false)
			case NoImage:
				_ = sendMessage("No image provided!", message.Chat.Id, false)
			}
		}
	}()

	var img image.Image
	img, err = getImage(message)
	fmt.Println(img)
}
