package main

import (
	"github.com/nsf/termbox-go"
	"github.com/dragonfi/go-retro/snake/arena"
)

func main() {
	x, y := Init()
	defer Close()
	offsetx, offsety := 2, 2
	aw := NewArenaWidget(offsetx, offsety, x-2*offsetx, y-2*offsety)

	aw.KeyMap = KeyMap{
		termbox.KeyEsc:        func() { aw.Exit() },
		termbox.KeyEnter:      func() { aw.ResetArena() },
		termbox.KeyArrowRight: func() { aw.SetSnakeHeading(0, arena.EAST) },
		termbox.KeyArrowUp:    func() { aw.SetSnakeHeading(0, arena.NORTH) },
		termbox.KeyArrowLeft:  func() { aw.SetSnakeHeading(0, arena.WEST) },
		termbox.KeyArrowDown:  func() { aw.SetSnakeHeading(0, arena.SOUTH) },
	}

	aw.RuneMap = RuneMap{
		'D': func() { aw.SetSnakeHeading(1, arena.EAST) },
		'W': func() { aw.SetSnakeHeading(1, arena.NORTH) },
		'A': func() { aw.SetSnakeHeading(1, arena.WEST) },
		'S': func() { aw.SetSnakeHeading(1, arena.SOUTH) },
		'd': func() { aw.SetSnakeHeading(1, arena.EAST) },
		'w': func() { aw.SetSnakeHeading(1, arena.NORTH) },
		'a': func() { aw.SetSnakeHeading(1, arena.WEST) },
		's': func() { aw.SetSnakeHeading(1, arena.SOUTH) },
	}

	aw.Run()
}
