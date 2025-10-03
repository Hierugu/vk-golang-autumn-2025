package main

type Outside struct {
	Room
}

func (o Outside) lookAroundText() string {
	return "на улице весна. можно пройти - домой"
}

func (o Outside) moveText() string {
	return "на улице весна. можно пройти - домой"
}
