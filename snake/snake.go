package main

import (
	"fmt"
	"github.com/dragonfi/go-retro/snake/arena"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

func putString(x, y int, s string) {
	for i, r := range s {
		termbox.SetCell(x+i, y, r, 0, 0)
	}
}

type Position struct {
	X, Y int
}

type ArenaWidget struct {
	arena  arena.Arena
	offset Position
	size   Position
	state  arena.State
}

func (w *ArenaWidget) Tick() {
	w.arena.Tick()
	w.state = w.arena.State()
}

func (w *ArenaWidget) SetSnakeHeading(direction arena.Direction) {
	w.arena.SetSnakeHeading(direction)
}

func (w ArenaWidget) setCell(x, y int, r rune, fg, bg termbox.Attribute) {
	termbox.SetCell(w.offset.X+x, w.offset.Y+y, r, fg, bg)
}

func (w ArenaWidget) putString(x, y int, str string) {
	putString(w.offset.X+x, w.offset.Y+y, str)
}

func (w ArenaWidget) drawBorder() {
	s := w.state
	for i := -1; i <= s.Size.X; i++ {
		for j := -1; j <= s.Size.Y; j++ {
			if i == -1 || i == s.Size.X || j == -1 || j == s.Size.Y {
				w.setCell(i, j, '#', 0, 0)
			}
		}
	}
}

func (w ArenaWidget) drawSnakes() {
	for _, snake := range w.state.Snakes {
		w.drawSnake(snake)
	}
}
func (w ArenaWidget) drawSnake(snake arena.Snake) {
	for _, p := range snake.Segments {
		w.setCell(p.X, p.Y, '#', 0, 0)
	}
}

func (w ArenaWidget) drawPointItem() {
	p := w.state.PointItem
	w.setCell(p.X, p.Y, '*', 0, 0)
}

func (w ArenaWidget) putGameOverText() {
	s := w.state
	w.putString(s.Size.X/2-9, s.Size.Y/2-3, "##################")
	w.putString(s.Size.X/2-9, s.Size.Y/2-2, "#    Game Over   #")
	w.putString(s.Size.X/2-9, s.Size.Y/2-1, "#                #")
	w.putString(s.Size.X/2-9, s.Size.Y/2+0, "# Enter: Restart #")
	w.putString(s.Size.X/2-9, s.Size.Y/2+1, "# ESC: Exit      #")
	w.putString(s.Size.X/2-9, s.Size.Y/2+2, "##################")
}

func (w ArenaWidget) putScore() {
	s := w.state
	for i, snake := range s.Snakes {
		w.putString(1, 1+i, fmt.Sprintf("Score: %d", len(snake.Segments)))
	}
}

func (w ArenaWidget) Draw() {
	w.drawBorder()
	w.putScore()
	w.drawSnakes()
	w.drawPointItem()
	if w.state.GameIsOver {
		w.putGameOverText()
	}
}

func (w *ArenaWidget) ResetArena() {
	w.arena = arena.New(w.size.X, w.size.Y)
	w.arena.AddSnake(w.size.X/2, w.size.Y/2, 5, arena.EAST)
}

func NewArenaWidget(ox, oy, x, y int) ArenaWidget {
	w := ArenaWidget{offset: Position{ox, oy}, size: Position{x, y}}
	w.ResetArena()
	return w
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
	aw := NewArenaWidget(offsetx, offsety, x-2*offsetx, y-2*offsety)
	tick := time.Tick(100 * time.Millisecond)
	event := eventChannel()
	running := true

	handleKey := map[termbox.Key]func(){
		termbox.KeyEsc:        func() { running = false },
		termbox.KeyEnter:      func() { aw.ResetArena() },
		termbox.KeyArrowRight: func() { aw.SetSnakeHeading(arena.EAST) },
		termbox.KeyArrowUp:    func() { aw.SetSnakeHeading(arena.NORTH) },
		termbox.KeyArrowLeft:  func() { aw.SetSnakeHeading(arena.WEST) },
		termbox.KeyArrowDown:  func() { aw.SetSnakeHeading(arena.SOUTH) },
	}

	for running {
		termbox.Clear(0, 0)
		aw.Draw()
		termbox.Flush()
		select {
		case ev := <-event:
			if ev.Type == termbox.EventKey {
				f := handleKey[ev.Key]
				if f != nil {
					f()
				}
			}
		case <-tick:
			aw.Tick()
		}
	}
}
