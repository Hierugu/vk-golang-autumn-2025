package main

import (
	"strings"
)

type RoomInterface interface {
	lookAroundText() string
	moveText() string
	getConnectionsString() string
	getItemsString() string
}

func (r Room) getConnectionsString() string {
	var names []string
	for _, room := range r.connections {
		if room != nil {
			switch v := room.(type) {
			case *Room:
				names = append(names, v.name)
			case *Hall:
				names = append(names, v.name)
			case *Kitchen:
				names = append(names, v.name)
			case *Bedroom:
				names = append(names, v.name)
			case *Outside:
				names = append(names, v.name)
			default:
				names = append(names, "unknown")
			}
		}
	}
	return "можно пройти - " + strings.Join(names, ", ")
}

func (r Room) getItemsString() string {
	var parts []string
	for _, obj := range r.objects {
		if len(obj.items) == 0 {
			continue
		}
		var itemNames []string
		for _, item := range obj.items {
			itemNames = append(itemNames, item.name)
		}
		parts = append(parts, obj.prepositionalName+": "+strings.Join(itemNames, ", "))
	}
	return strings.Join(parts, ", ")
}

func (r Room) lookAroundText() string {
	return "lookAroundText"
}

func (r Room) moveText() string {
	return "moveText"
}
