package main

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

var commandRe = regexp.MustCompile(`^\/[a-zA-Z]*`)
var handlers = map[string]HandlerFunc{
	"/help":     helpHandler,
	"/new":      newHandler,
	"/assign":   assignHandler,
	"/unassign": unassignHandler,
	"/resolve":  resolveHandler,
	"/my":       myHandler,
	"/owner":    ownerHandler,
	"/tasks":    tasksHandler,
}

type HandlerFunc func(msg *tgbotapi.Message)

func helpHandler(msg *tgbotapi.Message) {
	replyText := "" +
		"/tasks\n" +
		"/new XXX YYY ZZZ - создаёт новую задачу\n" +
		"/assign_$ID - делает пользователя исполнителем задачи\n" +
		"/unassign_$ID - снимает задачу с текущего исполнителя\n" +
		"/resolve_$ID - выполняет задачу, удаляет её из списка\n" +
		"/my - показывает задачи, которые назначены на меня\n" +
		"/owner - показывает задачи, которые были созданы мной Подробности форматирования смотрите в тестах."
	sendMsg(replyMsg(msg, replyText))
}

func newHandler(msg *tgbotapi.Message) {
	load := msg.Text[4:]
	taskName := strings.TrimSpace(load)
	if taskName == "" {
		sendMsg(replyMsg(msg, "Название задачи не может быть пустым"))
		return
	}
	task, err := createTask(taskName, msg.From)
	if err != nil {
		sendMsg(replyMsg(msg, err.Error()))
		return
	}
	_, err = bot.Send(replyMsg(msg, fmt.Sprintf("Задача \"%s\" создана, id=%d", task.name, task.id)))
	panicIf(err)
}

func assignHandler(msg *tgbotapi.Message) {
	load := msg.Text[8:]
	taskId, err := strconv.ParseInt(strings.TrimSpace(load), 10, 64)
	if err != nil {
		sendMsg(replyMsg(msg, "Неверный ID задачи"))
		return
	}
	SendToOwner := true
	// Плохой код из-за плохой архитектуры из-за плохого понимания нюансов уведомлений, а много кода уже написано
	index := slices.IndexFunc(taskStorage, func(t Task) bool { return t.id == taskId })
	if taskStorage[index].assign.id != msg.From.ID && taskStorage[index].assign.id != taskStorage[index].owner.id && taskStorage[index].assign.id != 0 {
		asgnReply := tgbotapi.NewMessage(taskStorage[index].assign.id, fmt.Sprintf("Задача \"%s\" назначена на @%s", taskStorage[index].name, msg.From.UserName))
		asgnReply.ReplyToMessageID = int(taskStorage[index].assign.id)
		sendMsg(asgnReply)
		SendToOwner = false
	}

	task, err := assignTask(taskId, msg.From)
	log.Println("Tasks:", taskStorage)
	if err != nil {
		sendMsg(replyMsg(msg, err.Error()))
		return
	}
	sendMsg(replyMsg(msg, fmt.Sprintf("Задача \"%s\" назначена на вас", task.name)))
	if msg.From.ID != task.owner.id && task.owner.id != 0 && SendToOwner {
		ownerReply := tgbotapi.NewMessage(task.owner.id, fmt.Sprintf("Задача \"%s\" назначена на @%s", task.name, task.assign.name))
		ownerReply.ReplyToMessageID = int(task.owner.id)
		sendMsg(ownerReply)
	}
}

func unassignHandler(msg *tgbotapi.Message) {
	load := msg.Text[10:]
	taskId, err := strconv.ParseInt(strings.TrimSpace(load), 10, 64)
	if err != nil {
		sendMsg(replyMsg(msg, "Неверный ID задачи"))
		return
	}
	task, err := unassignTask(taskId, msg.From)
	if err != nil {
		sendMsg(replyMsg(msg, err.Error()))
		return
	}
	sendMsg(replyMsg(msg, "Принято"))
	if msg.From.ID != task.owner.id && task.owner.id != 0 {
		ownerReply := tgbotapi.NewMessage(task.owner.id, fmt.Sprintf("Задача \"%s\" осталась без исполнителя", task.name))
		ownerReply.ReplyToMessageID = int(task.owner.id)
		sendMsg(ownerReply)
	}
}

func resolveHandler(msg *tgbotapi.Message) {
	load := msg.Text[9:]
	taskId, err := strconv.ParseInt(strings.TrimSpace(load), 10, 64)
	if err != nil {
		sendMsg(replyMsg(msg, "Неверный ID задачи"))
		return
	}
	task, err := resolveTask(taskId)
	if err != nil {
		sendMsg(replyMsg(msg, err.Error()))
		return
	}
	sendMsg(replyMsg(msg, fmt.Sprintf("Задача \"%s\" выполнена", task.name)))
	if msg.From.ID != task.owner.id && task.owner.id != 0 {
		ownerReply := tgbotapi.NewMessage(task.owner.id, fmt.Sprintf("Задача \"%s\" выполнена @%s", task.name, msg.From.UserName))
		ownerReply.ReplyToMessageID = int(task.owner.id)
		sendMsg(ownerReply)
	}
}

func myHandler(msg *tgbotapi.Message) {
	parts := make([]string, 0)
	for _, task := range taskStorage {
		if task.assign.id == msg.From.ID && !task.done {
			parts = append(parts, fmt.Sprintf("%d. %s by @%s", task.id, task.name, task.owner.name))
			parts = append(parts, fmt.Sprintf("/unassign_%d /resolve_%d", task.id, task.id))
		}
	}
	var reply string
	if len(parts) == 0 {
		reply = ""
	} else {
		reply = strings.Join(parts, "\n")
	}
	sendMsg(replyMsg(msg, reply))
}

func ownerHandler(msg *tgbotapi.Message) {
	parts := make([]string, 0)
	for _, task := range taskStorage {
		if task.owner.id == msg.From.ID && !task.done {
			parts = append(parts, fmt.Sprintf("%d. %s by @%s", task.id, task.name, task.owner.name))
			parts = append(parts, fmt.Sprintf("/assign_%d", task.id))
		}
	}
	var reply string
	if len(parts) == 0 {
		reply = ""
	} else {
		reply = strings.Join(parts, "\n")
	}
	sendMsg(replyMsg(msg, reply))
}

func tasksHandler(msg *tgbotapi.Message) {
	parts := make([]string, 0, len(taskStorage))
	for _, task := range taskStorage {
		if !task.done {
			if len(parts) > 0 {
				parts = append(parts, "")
			}
			parts = append(parts, fmt.Sprintf("%d. %s by @%s", task.id, task.name, task.owner.name))
			if task.assign.id != 0 {
				if task.assign.id == msg.From.ID {
					parts = append(parts, "assignee: я")
					parts = append(parts, fmt.Sprintf("/unassign_%d /resolve_%d", task.id, task.id))
				} else {
					parts = append(parts, fmt.Sprintf("assignee: @%s", task.assign.name))
				}
			} else {
				parts = append(parts, fmt.Sprintf("/assign_%d", task.id))
			}
		}
	}
	var reply string
	if len(parts) == 0 {
		reply = "Нет задач"
	} else {
		reply = strings.Join(parts, "\n")
	}
	sendMsg(replyMsg(msg, reply))
}
