package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-telebot/models"
	"image"
	"io"
	"net/http"
	"os"
)

var baseUrl string
var auth string
var offset int
var handlers map[string]func(message models.Message)

func init() {
	token := os.Getenv("TOKEN")
	baseUrl = "https://api.telegram.org"
	auth = fmt.Sprintf("/bot%s", token)
	offset = 0

	handlers = map[string]func(message models.Message){}

}

func main() {
	for {
		updates, err := getUpdates()
		if err != nil {
			fmt.Println("Ошибка получения обновлений ->", err)
			continue
		}

		fmt.Printf("Получены обновления %d\n", len(updates))

		for _, update := range updates {
			go handleUpdate(update)
		}
	}
}

func getUpdates() ([]models.Update, error) {
	fmt.Println("Получение обновлений -> offset =", offset)
	url := fmt.Sprintf("%s%s/getUpdates?offset=%d&timeout=15", baseUrl, auth, offset)
	response, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	var updatesResult models.UpdateResult
	if err := json.Unmarshal(body, &updatesResult); err != nil {
		return nil, err
	}

	return updatesResult.Result, nil
}

func sendMessage(text string, chatId int, keyboard bool) error {

	fmt.Println("Отправка сообщения")

	url := fmt.Sprintf("%s%s/sendMessage", baseUrl, auth)
	message := map[string]any{
		"chat_id":    chatId,
		"text":       text,
		"parse_mode": "markdown",
	}

	if keyboard {
		message["reply_markup"] = makeKeyboard()
	} else {
		message["reply_markup"] = map[string]bool{
			"remove_keyboard": true,
		}
	}

	jsonBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, sendErr := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if sendErr != nil {
		return err
	}
	fmt.Printf("Отправлено сообщение -> %+v\n", message)

	return nil
}

func getImage(message models.Message) (image.Image, error) {
	if message.Photo == nil {
		return nil, NoImage{}
	}

	photo := message.Photo[0]
	file, _ := getFileInfo(photo.FileId)

	url := fmt.Sprintf("%s/file/%s/%s", baseUrl, auth, file.FilePath)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	img, _, imgErr := image.Decode(response.Body)
	fmt.Println(imgErr)
	return img, nil
}

func getFileInfo(fileId string) (*models.File, error) {

	url := fmt.Sprintf("%s%s/getFile?file_id=%s", baseUrl, auth, fileId)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	var fileResult models.FileResult
	if err = json.Unmarshal(body, &fileResult); err != nil {
		return nil, err
	}

	return &fileResult.Result, nil
}

func makeKeyboard() models.Keyboard {
	keyboard := models.Keyboard{
		Keyboard: [][]models.Key{
			{{Text: "1"}, {Text: "2"}, {Text: "3"}},
			{{Text: "4"}, {Text: "5"}, {Text: "6"}},
			{{Text: "7"}, {Text: "8"}, {Text: "9"}},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: false,
	}
	return keyboard
}

func handleUpdate(update models.Update) {
	fmt.Printf("Обработка обновления -> %+v\n", update)

	if update.UpdateId >= offset {
		offset = update.UpdateId + 1
		handler := handlers[update.Message.Text]
		if handler == nil {
			handler = imageHandler
		}
		handler(update.Message)
	}
}
