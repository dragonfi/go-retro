package main

import (
	"github.com/nsf/termbox-go"
)

type Position struct {
	X, Y int
}

type Snake []Position

// NewSnake returns a new Snake pointing in the +X direction
func NewSnake(head Position, length int) Snake {
	snake := make(Snake, length, 100)
	snake[0] = head
	return snake
}

func main() {
	err := termbox.Init()
	if err != nil {
			panic(err)
	}
	defer termbox.Close()

	_ = NewSnake(Position{0,0}, 1)
	
}
