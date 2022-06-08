package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
	"tgbotcrypwallet/internal/getprice"
)

// Напишем свою базу данных и определим для нее тип и переменную
type wallet map[string]float64 // [валюта]количество валюты
var db = map[int64]wallet{}    // [Chat.ID]слайс map

func main() {
	// Create object bot
	bot, err := tgbotapi.NewBotAPI("5466103665:AAFqXxh8GwD_-a963U9Zq2CCK8pxj2P2Wj4")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true // Add debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0) // Create chanel
	u.Timeout = 60             // This is necessary so that the connection is permanent and does not fall off

	updates := bot.GetUpdatesChan(u) // Create a channel for updates

	for update := range updates {
		if update.Message == nil { // ignore any non-message Updates
			continue
		}
		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		command := strings.Split(update.Message.Text, " ")
		switch command[0] {
		case "ADD":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат ввода данных"))
			}
			amount, err := strconv.ParseFloat(command[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			}
			if _, ok := db[update.Message.Chat.ID]; !ok { // Проверка существоввания чта с таким ID
				db[update.Message.Chat.ID] = wallet{}
			}
			db[update.Message.Chat.ID][command[1]] += amount
			balanceText := fmt.Sprintf("%f\n", db[update.Message.Chat.ID][command[1]])
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))
			// bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Валюта добавлена"))
		case "SUB":
			if len(command) != 3 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат ввода данных"))
			}
			amount, err := strconv.ParseFloat(command[2], 64)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			}
			if _, ok := db[update.Message.Chat.ID]; !ok { // Проверка существоввания чта с таким ID
				continue // Выходим, т.к. не нужно отнимать валюту из кошелька, который еще не создан
			}
			// Сделать проверку на  < 0
			db[update.Message.Chat.ID][command[1]] -= amount
			balanceText := fmt.Sprintf("%0.f\n", db[update.Message.Chat.ID][command[1]])
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, balanceText))
		case "DEL":
			if len(command) != 2 {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат ввода данных"))
			}
			delete(db[update.Message.Chat.ID], command[1])
		case "SHOW":
			msg := ""
			var sum float64 = 0
			for key, value := range db[update.Message.Chat.ID] {
				price, _ := getprice.GetPrice(key)
				sum += value * price
				msg += fmt.Sprintf("%s: %f [%.2f]\n", key, value, value*price)
			}
			msg += fmt.Sprintf("Total: %.2f\n", sum)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда не найдена"))
		}
	}
}
