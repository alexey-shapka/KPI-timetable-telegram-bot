package  main

import (
	"strconv"
	"strings"
	"regexp"
	"fmt"
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

var storage map[int64]string
var memory map[int64]string
var logCount map[int64][]string

func main() {

	storage = make(map[int64]string)
	memory = make(map[int64]string)
	logCount = make(map[int64][]string)

	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)


		_, id := logCount[msg.ChatID]
		if id {
		}	else{
			logCount[msg.ChatID] = []string{"0", update.Message.From.FirstName}
		}

		if msg.Text == "/start" {
			msgText := "*Сomands:*\n/table\n/today\n/tomorrow\n/setgroup"
			message := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
			message.ParseMode = "markdown"
			bot.Send(message)

		}	else if msg.Text == "/table" {

			action := countActions(logCount[msg.ChatID])
			logCount[msg.ChatID] = []string{action, update.Message.From.FirstName}
			
			text := ""
			_, id := storage[msg.ChatID]
					if id {
						if storage[msg.ChatID] == "inputerror" {
							text = "Невірно вказана група."
						}	else{
							text = table(storage[msg.ChatID], checkWeek())
						}
					}	else {
						text = "Введіть назву групи\n/setgroup"
					}
			
			message := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			message.ParseMode = "markdown"
			bot.Send(message)

		}	else if msg.Text == "/today" {

			action := countActions(logCount[msg.ChatID])
			logCount[msg.ChatID] = []string{action, update.Message.From.FirstName}

			message := tgbotapi.NewMessage(update.Message.Chat.ID, checkData(msg.ChatID, checkDay()))
			message.ParseMode = "markdown"
			bot.Send(message)
			
		}	else if msg.Text == "/tomorrow" {

			action := countActions(logCount[msg.ChatID])
			logCount[msg.ChatID] = []string{action, update.Message.From.FirstName}

			day := checkDay()+1
			if day == 7{
				day = 0
			}
			message := tgbotapi.NewMessage(update.Message.Chat.ID,checkData(msg.ChatID, day))
			message.ParseMode = "markdown"
			bot.Send(message)

		}	else if msg.Text == "/setgroup"{

			action := countActions(logCount[msg.ChatID])
			logCount[msg.ChatID] = []string{action, update.Message.From.FirstName}

			storage[msg.ChatID]="input"
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Введіть назву групи"))
		
		}	else if msg.Text == "/rating"{

			action := countActions(logCount[msg.ChatID])
			logCount[msg.ChatID] = []string{action, update.Message.From.FirstName}

			text := ""
			_, id := storage[msg.ChatID]
					if id {
						if storage[msg.ChatID] == "inputerror" {
							text = "Невірно вказана група."
						}	else{
							text = getRating(storage[msg.ChatID])
						}
					}	else {
						text = "Введіть назву групи\n/setgroup"
					}
			message := tgbotapi.NewMessage(update.Message.Chat.ID,text)
			message.ParseMode = "markdown"
			bot.Send(message)

		}	else if msg.Text == "/teacher"{

			action := countActions(logCount[msg.ChatID])
			logCount[msg.ChatID] = []string{action, update.Message.From.FirstName}

			_, id := storage[msg.ChatID]
			if id {
			memory[msg.ChatID] = storage[msg.ChatID]
			}

			storage[msg.ChatID] = "teacherInput"
			message := tgbotapi.NewMessage(update.Message.Chat.ID, "Введіть викладача")
			message.ParseMode = "markdown"
			bot.Send(message)

		}	else if msg.Text == "/log"{

			viewLog := ""
			for k, v := range logCount{
				viewLog += fmt.Sprintf("%s : %s : %s\n",v[0],strconv.FormatInt(k, 10),v[1])
			}

			message := tgbotapi.NewMessage(update.Message.Chat.ID, viewLog)
			bot.Send(message)
		}

		if storage[msg.ChatID] == "input" {
			if msg.Text != "/setgroup" {
				text := ""
				re := regexp.MustCompile(`[^a-zA-zа-яА-я0-9 і -]`)
				group := re.ReplaceAllString(strings.Replace(strings.ToLower(msg.Text), "и", "i", -1), "")

				if len(table(group,1)) == 217 {
					storage[msg.ChatID] = "inputerror"
					text = "Групу не знайдено.\nСпробуйте ще раз.\n/setgroup"

				}	else{
					text = "Група успішно змінена"
					storage[msg.ChatID] = group
				}
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
			}
		}	else if storage[msg.ChatID] == "teacherInput"{
				if msg.Text != "/teacher" {
					teacherInfo := getTeacherschedule(msg.Text)

					if len(teacherInfo)<10{
						teacherInfo = "Викладача не знайдено."
					}

					message := tgbotapi.NewMessage(update.Message.Chat.ID, teacherInfo)
					message.ParseMode = "markdown"

					_, id := memory[msg.ChatID]
					if id {
						storage[msg.ChatID] = memory[msg.ChatID]
						bot.Send(message)
					}	else{
						delete(storage, msg.ChatID)
						bot.Send(message)
					}
			}
		}
	}
}

func table(group string, week int) string{
	result := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",getDay(getinfo(group,week,1)),getDay(getinfo(group,week,2)),getDay(getinfo(group,week,3)),
	getDay(getinfo(group,week,4)),getDay(getinfo(group,week,5)),getDay(getinfo(group,week,6)))

	return result
}

func checkData(chatId int64, day int) string{
	text := ""
	_, id := storage[chatId]
			if id {
				if storage[chatId] == "inputerror" {
					text = "Невірно вказана група."
				}	else {
					text =  getDay(getinfo(storage[chatId],checkWeek(),day))
				}
			}	else{
					text = "Введіть назву групи\n/setgroup"
				}
	return text
}

func countActions(info[] string) string{
	actions, _ := strconv.Atoi(info[0])
	actions++
	return strconv.Itoa(actions)
}
