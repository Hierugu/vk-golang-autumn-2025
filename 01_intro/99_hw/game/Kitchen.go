package main

import (
	"strings"
)

type Kitchen struct {
	Room
}

func (k Kitchen) lookAroundText() string {
	var result []string = []string{"ты находишься на кухне"}

	if itemsStr := k.getItemsString(); itemsStr != "" {
		result = append(result, itemsStr)
	}

	hasBackpack, hasNotes := false, false
	for _, item := range world.player.inventory {
		if item.name == "рюкзак" {
			hasBackpack = true
		}
		if item.name == "конспекты" {
			hasNotes = true
		}
	}

	if hasBackpack && hasNotes {
		result = append(result, "надо идти в универ")
	} else {
		result = append(result, "надо собрать рюкзак и идти в универ")
	}

	return strings.Join(result, ", ") + ". " + k.getConnectionsString()
}

func (k Kitchen) moveText() string {
	return "кухня, ничего интересного" + ". " + k.getConnectionsString()
}
