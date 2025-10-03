package main

type Bedroom struct {
	Room
}

func (b Bedroom) lookAroundText() string {
	if is := b.getItemsString(); is == "" {
		return "пустая комната" + ". " + b.getConnectionsString()
	}
	return b.getItemsString() + ". " + b.getConnectionsString()
}

func (b Bedroom) moveText() string {
	return "ты в своей комнате" + ". " + b.getConnectionsString()
}
