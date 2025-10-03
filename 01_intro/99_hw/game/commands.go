package main

func lookAround() string {
	return world.player.currentRoom.lookAroundText()
}

func goTo(roomName string) string {
	var currentRoomName string
	switch v := world.player.currentRoom.(type) {
	case *Room:
		currentRoomName = v.name
	case *Hall:
		currentRoomName = v.name
	case *Kitchen:
		currentRoomName = v.name
	case *Bedroom:
		currentRoomName = v.name
	case *Outside:
		currentRoomName = v.name
	}
	if currentRoomName == "коридор" && roomName == "улица" {
		if world.objects[3].status == "open" {
			world.player.currentRoom = world.rooms["улица"]
			return world.player.currentRoom.moveText()
		} else {
			return "дверь закрыта"
		}
	}

	if currentRoomName == "улица" && roomName == "домой" {
		world.player.currentRoom = world.rooms["коридор"]
		return world.player.currentRoom.moveText()
	}

	var connections []RoomInterface
	switch v := world.player.currentRoom.(type) {
	case *Room:
		connections = v.connections
	case *Hall:
		connections = v.connections
	case *Kitchen:
		connections = v.connections
	case *Bedroom:
		connections = v.connections
	case *Outside:
		connections = v.connections
	}
	for _, room := range connections {
		var name string
		switch r := room.(type) {
		case *Room:
			name = r.name
		case *Hall:
			name = r.name
		case *Kitchen:
			name = r.name
		case *Bedroom:
			name = r.name
		case *Outside:
			name = r.name
		}
		if name == roomName {
			world.player.currentRoom = room
			return room.moveText()
		}
	}
	return "нет пути в " + roomName
}

func takeItem(itemName string) string {
	hasBackpack := false
	for _, item := range world.player.inventory {
		if item.name == "рюкзак" {
			hasBackpack = true
			break
		}
	}
	if !hasBackpack {
		return "некуда класть"
	}

	var objects []*Object
	switch v := world.player.currentRoom.(type) {
	case *Room:
		objects = v.objects
	case *Hall:
		objects = v.objects
	case *Kitchen:
		objects = v.objects
	case *Bedroom:
		objects = v.objects
	case *Outside:
		objects = v.objects
	}
	for _, obj := range objects {
		for i, item := range obj.items {
			if item.name == itemName {
				world.player.inventory = append(world.player.inventory, item)
				obj.items = append(obj.items[:i], obj.items[i+1:]...)
				return "предмет добавлен в инвентарь: " + item.name
			}
		}
	}
	return "нет такого"
}

func wearItem(itemName string) string {
	if itemName != "рюкзак" {
		return "нельзя надеть"
	}

	var objects []*Object
	switch v := world.player.currentRoom.(type) {
	case *Room:
		objects = v.objects
	case *Hall:
		objects = v.objects
	case *Kitchen:
		objects = v.objects
	case *Bedroom:
		objects = v.objects
	case *Outside:
		objects = v.objects
	}
	for _, obj := range objects {
		for i, item := range obj.items {
			if item.name == itemName {
				world.player.inventory = append(world.player.inventory, item)
				obj.items = append(obj.items[:i], obj.items[i+1:]...)
				return "вы надели: " + item.name
			}
		}
	}
	return "нельзя надеть"
}

func useItem(itemName, objectName string) string {
	var objects []*Object
	switch v := world.player.currentRoom.(type) {
	case *Room:
		objects = v.objects
	case *Hall:
		objects = v.objects
	case *Kitchen:
		objects = v.objects
	case *Bedroom:
		objects = v.objects
	case *Outside:
		objects = v.objects
	}
	for _, item := range world.player.inventory {
		if item.name == itemName {
			for _, obj := range objects {
				if obj.name == objectName {
					if obj.name == "дверь" && item.name == "ключи" {
						if obj.status == "closed" {
							obj.status = "open"
							return "дверь открыта"
						} else {
							obj.status = "closed"
							return "дверь закрыта"
						}
					}
				}
			}
			return "не к чему применить"
		}
	}
	return "нет предмета в инвентаре - " + itemName
}
