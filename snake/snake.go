package main

import (
	"./arena"
	"math/rand"
	"time"
	"github.com/nsf/termbox-go"
)

func draw(s arena.State) {
	for _, p := range s.Snake.Segments {
		termbox.SetCell(p.X, p.Y, '#', 0, 0)
	}
	p := s.PointItem
	termbox.SetCell(p.X, p.Y, '*', 0, 0)
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
	a := arena.New(x, y)
	tick := time.Tick(50*time.Millisecond)
	event := eventChannel()
	running := true
	for running {
		termbox.Clear(0, 0)
		draw(a.State())
		termbox.Flush()
		select {
		case ev:=<-event:
			if ev.Type == termbox.EventKey {
				switch ev.Key {
				case termbox.KeyEsc:
					running = false
				case termbox.KeyArrowRight:
					a.SetSnakeHeading(arena.EAST)
				case termbox.KeyArrowUp:
					a.SetSnakeHeading(arena.NORTH)
				case termbox.KeyArrowLeft:
					a.SetSnakeHeading(arena.WEST)
				case termbox.KeyArrowDown:
					a.SetSnakeHeading(arena.SOUTH)
				}
			}
		case <-tick:
			a.Tick()
		}
	}
}
