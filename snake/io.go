package main

import (
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

func putString(x, y int, s string) {
	for i, r := range s {
		termbox.SetCell(x+i, y, r, 0, 0)
	}
}

func Init() (x, y int) {
	rand.Seed(time.Now().UnixNano())
	termbox.Init()
	return termbox.Size()
}

func Close() {
	termbox.Close()
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

type KeyMap map[termbox.Key]func()
type RuneMap map[rune]func()

func handleEvent(ev termbox.Event, keyMap KeyMap, runeMap RuneMap) {
	if ev.Type == termbox.EventKey {
		f := runeMap[ev.Ch]
		if ev.Ch == 0 {
			f = keyMap[ev.Key]
		}
		if f != nil {
			f()
		}
	}
}
