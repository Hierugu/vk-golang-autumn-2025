package main

type Hall struct {
	Room
}

func (h Hall) lookAroundText() string {
	return "ничего интересного" + ". " + h.getConnectionsString()
}

func (h Hall) moveText() string {
	return "ничего интересного" + ". " + h.getConnectionsString()
}
