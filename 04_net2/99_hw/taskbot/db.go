package main

import (
	"fmt"
	"log"
	"slices"

	tgbotapi "github.com/skinass/telegram-bot-api/v5"
)

type User struct {
	id   int64
	name string
}

type Task struct {
	id     int64
	name   string
	owner  User
	assign User
	done   bool
}

var taskStorage = []Task{}

func createTask(name string, owner *tgbotapi.User) (*Task, error) {
	if slices.IndexFunc(taskStorage, func(t Task) bool { return t.name == name }) != -1 {
		return nil, fmt.Errorf("Задача с таким названием уже существует")
	}
	newID := int64(len(taskStorage) + 1)
	taskStorage = append(taskStorage, Task{id: newID, name: name, owner: User{name: owner.UserName, id: owner.ID}, assign: User{}, done: false})
	return &taskStorage[len(taskStorage)-1], nil
}

func assignTask(id int64, assignee *tgbotapi.User) (*Task, error) {
	log.Println("Tasks:", taskStorage)
	log.Println("Assigning task:", id, "to user:", assignee.UserName)
	index := slices.IndexFunc(taskStorage, func(t Task) bool { return t.id == id })
	if index == -1 {
		return nil, fmt.Errorf("Задача с таким ID не существует")
	}
	taskStorage[index].assign = User{name: assignee.UserName, id: assignee.ID}
	return &taskStorage[index], nil
}

func unassignTask(id int64, assignee *tgbotapi.User) (*Task, error) {
	index := slices.IndexFunc(taskStorage, func(t Task) bool { return t.id == id })
	if index == -1 {
		return nil, fmt.Errorf("Задача с таким ID не существует")
	}
	if taskStorage[index].assign.id != assignee.ID {
		return nil, fmt.Errorf("Задача не на вас")
	}
	taskStorage[index].assign = User{}
	return &taskStorage[index], nil
}

func resolveTask(id int64) (*Task, error) {
	index := slices.IndexFunc(taskStorage, func(t Task) bool { return t.id == id })
	if index == -1 {
		return nil, fmt.Errorf("Задача с таким ID не существует")
	}
	if taskStorage[index].done {
		return nil, fmt.Errorf("Задача уже выполнена")
	}
	taskStorage[index].done = true
	return &taskStorage[index], nil
}
