package main

type World struct {
	player  Player
	rooms   map[string]RoomInterface
	objects map[int]*Object
	items   map[string]*Item
}

type Player struct {
	currentRoom RoomInterface
	inventory   []*Item
}

type Room struct {
	id          int
	name        string
	objects     []*Object
	connections []RoomInterface
}

type Object struct {
	id                int
	name              string
	prepositionalName string
	items             []*Item
	status            string
}

type Item struct {
	id   int
	name string
}
