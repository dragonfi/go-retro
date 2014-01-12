package main

import (
	"./arena"
	"fmt"
	"math/rand"
	"time"
	"github.com/nsf/termbox-go"
)

func putString(x, y int, s string) {
	for i, r := range s {
		termbox.SetCell(x+i, y, r, 0, 0)
	}
}

func draw(ox, oy int, s arena.State) {
	str := fmt.Sprintf("Score: %d", len(s.Snake.Segments))
	putString(ox+1, oy+1, str)
	for i := -1; i<=s.Size.X; i++ {
		for j := -1; j<=s.Size.Y; j++ {
			if i == -1 || i == s.Size.X || j == -1 || j == s.Size.Y {
				termbox.SetCell(ox+i, oy+j, '#', 0, 0)
			}
		}
	}
	for _, p := range s.Snake.Segments {
		termbox.SetCell(ox+p.X, oy+p.Y, '#', 0, 0)
	}
	p := s.PointItem
	termbox.SetCell(ox+p.X, oy+p.Y, '*', 0, 0)
	if s.GameIsOver {
		putString(ox+s.Size.X/2 - 8, oy+s.Size.Y/2-2, "   Game Over  ")
		putString(ox+s.Size.X/2 - 8, oy+s.Size.Y/2+0, "Enter: Restart")
		putString(ox+s.Size.X/2 - 8, oy+s.Size.Y/2+1, "ESC: Exit")
	}

}

func eventChannel() <-chan termbox.Event {
	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()
	return events
}

func main() {
	rand.Seed(time.Now().UnixNano())
	termbox.Init()
	defer termbox.Close()
	x, y := termbox.Size()
	offsetx, offsety := 2, 2
	ax, ay := x - 2*offsetx, y - 2*offsety
	a := arena.New(ax, ay)
	tick := time.Tick(100*time.Millisecond)
	event := eventChannel()
	running := true

	handleKey := map[termbox.Key]func() {
		termbox.KeyEsc: func(){running = false},
		termbox.KeyEnter: func(){a = arena.New(ax, ay)},
		termbox.KeyArrowRight: func(){a.SetSnakeHeading(arena.EAST)},
		termbox.KeyArrowUp: func(){a.SetSnakeHeading(arena.NORTH)},
		termbox.KeyArrowLeft: func(){a.SetSnakeHeading(arena.WEST)},
		termbox.KeyArrowDown: func(){a.SetSnakeHeading(arena.SOUTH)},
	}

	for running {
		termbox.Clear(0, 0)
		draw(offsetx, offsety, a.State())
		termbox.Flush()
		select {
		case ev:=<-event:
			if ev.Type == termbox.EventKey {
				f := handleKey[ev.Key]
				if f != nil {
					f()
				}
			}
		case <-tick:
			a.Tick()
		}
	}
}
