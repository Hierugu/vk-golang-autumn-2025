package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var world World

func main() {
	initGame()
	fmt.Println("Добро пожаловать в игру!")
	fmt.Println("Доступные команды: осмотреться, идти <куда>, взять <что>, надеть <что>, применить <что> <к чему>")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := strings.TrimSpace(scanner.Text())
		if command == "" {
			continue
		}
		result := handleCommand(command)
		fmt.Println(result)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка ввода:", err)
	}
}

func initGame() {
	world = World{items: make(map[string]*Item), objects: make(map[int]*Object), rooms: make(map[string]RoomInterface)}

	world.items["чай"] = &Item{id: 0, name: "чай"}
	world.items["ключи"] = &Item{id: 1, name: "ключи"}
	world.items["конспекты"] = &Item{id: 2, name: "конспекты"}
	world.items["рюкзак"] = &Item{id: 3, name: "рюкзак"}

	world.objects[0] = &Object{id: 0, name: "стол", prepositionalName: "на столе", items: []*Item{world.items["чай"]}, status: ""}
	world.objects[1] = &Object{id: 1, name: "стол", prepositionalName: "на столе", items: []*Item{world.items["ключи"], world.items["конспекты"]}, status: ""}
	world.objects[2] = &Object{id: 2, name: "стул", prepositionalName: "на стуле", items: []*Item{world.items["рюкзак"]}, status: ""}
	world.objects[3] = &Object{id: 3, name: "дверь", prepositionalName: "на двери", items: []*Item{}, status: "closed"}

	world.rooms["кухня"] = &Kitchen{
		Room: Room{
			id:          0,
			name:        "кухня",
			objects:     []*Object{world.objects[0]},
			connections: []RoomInterface{},
		},
	}
	world.rooms["коридор"] = &Hall{
		Room: Room{
			id:          1,
			name:        "коридор",
			objects:     []*Object{world.objects[3]},
			connections: []RoomInterface{},
		},
	}
	world.rooms["комната"] = &Bedroom{
		Room: Room{
			id:          2,
			name:        "комната",
			objects:     []*Object{world.objects[1], world.objects[2]},
			connections: []RoomInterface{},
		},
	}
	world.rooms["улица"] = &Outside{
		Room: Room{
			id:          3,
			name:        "улица",
			objects:     []*Object{},
			connections: []RoomInterface{},
		},
	}

	world.rooms["кухня"].(*Kitchen).Room.connections = append(world.rooms["кухня"].(*Kitchen).Room.connections, world.rooms["коридор"])
	world.rooms["комната"].(*Bedroom).Room.connections = append(world.rooms["комната"].(*Bedroom).Room.connections, world.rooms["коридор"])
	world.rooms["коридор"].(*Hall).Room.connections = append(world.rooms["коридор"].(*Hall).Room.connections, world.rooms["кухня"])
	world.rooms["коридор"].(*Hall).Room.connections = append(world.rooms["коридор"].(*Hall).Room.connections, world.rooms["комната"])
	world.rooms["коридор"].(*Hall).Room.connections = append(world.rooms["коридор"].(*Hall).Room.connections, world.rooms["улица"])
	world.rooms["улица"].(*Outside).Room.connections = append(world.rooms["улица"].(*Outside).Room.connections, world.rooms["коридор"])

	world.player = Player{currentRoom: world.rooms["кухня"], inventory: []*Item{}}
}

func handleCommand(command string) string {
	cmd := strings.Fields(command)
	if len(cmd) == 0 {
		return "неизвестная команда"
	}

	switch {
	case cmd[0] == "осмотреться" && len(cmd) == 1:
		return lookAround()
	case cmd[0] == "идти" && len(cmd) == 2:
		return goTo(cmd[1])
	case cmd[0] == "взять" && len(cmd) == 2:
		return takeItem(cmd[1])
	case cmd[0] == "надеть" && len(cmd) == 2:
		return wearItem(cmd[1])
	case cmd[0] == "применить" && len(cmd) == 3:
		return useItem(cmd[1], cmd[2])
	default:
		return "неизвестная команда"
	}
}
