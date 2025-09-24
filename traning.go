package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UnsplashResponse struct {
	Results []struct {
		Urls struct {
			Regular string `json:"regular"`
		} `json:"urls"`
	} `json:"results"`
}

func main() {

	botToken := "###"
	unsplashToken := "###"

	if botToken == "" || unsplashToken == "" {
		fmt.Println("Токены не заданы")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic("Ошибка инициализации бота: ", err)
	}

	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		query := update.Message.Text
		chatID := update.Message.Chat.ID

		if query == "" {
			msg := tgbotapi.NewMessage(chatID, "Пожалуйста, введите слово для поиска картинки.")
			bot.Send(msg)
			continue
		}

		imageURL, err := getImageFromUnsplash(query, unsplashToken)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка при получении картинки: "+err.Error())
			bot.Send(msg)
			continue
		}

		if imageURL == "" {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Картинки по запросу '%s' не найдены.", query))
			bot.Send(msg)
			continue
		}

		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(imageURL))
		photo.Caption = fmt.Sprintf("Картинка по запросу: %s", query)
		_, err = bot.Send(photo)
		if err != nil {
			msg := tgbotapi.NewMessage(chatID, "Ошибка при отправке картинки: "+err.Error())
			bot.Send(msg)
		}
	}
}

func getImageFromUnsplash(query, accessKey string) (string, error) {
	url := fmt.Sprintf("https://api.unsplash.com/search/photos?query=%s&per_page=1", query)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", accessKey))

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unsplash API вернул статус: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var unsplashResp UnsplashResponse
	err = json.Unmarshal(body, &unsplashResp)
	if err != nil {
		return "", err
	}

	if len(unsplashResp.Results) == 0 {
		return "", nil
	}

	return unsplashResp.Results[0].Urls.Regular, nil
}
