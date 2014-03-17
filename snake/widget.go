package main

import (
	"fmt"
	"github.com/dragonfi/go-retro/snake/arena"
	"github.com/nsf/termbox-go"
	"time"
)

var colors = map[string]termbox.Attribute{
	"snake1":    termbox.ColorGreen | termbox.AttrBold,
	"snake2":    termbox.ColorYellow | termbox.AttrBold,
	"pointItem": termbox.ColorCyan | termbox.AttrBold,
}

func getSnakeColor(i int) termbox.Attribute {
	key := ""
	switch i {
	case 0:
		key = "snake1"
	case 1:
		key = "snake2"
	}
	return colors[key]
}

type Position struct {
	X, Y int
}

type ArenaWidget struct {
	arena   arena.Arena
	offset  Position
	size    Position
	state   arena.State
	running bool
	KeyMap  KeyMap
	RuneMap RuneMap
}

func (w *ArenaWidget) Tick() {
	w.arena.Tick()
	w.state = w.arena.State()
}

func (w *ArenaWidget) SetSnakeHeading(snake int, direction arena.Direction) {
	w.arena.SetSnakeHeading(snake, direction)
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
	for i, snake := range w.state.Snakes {
		w.drawSnake(getSnakeColor(i), snake)
	}
}
func (w ArenaWidget) drawSnake(color termbox.Attribute, snake arena.Snake) {
	for i, p := range snake.Segments {
		char := '#'
		if i == 0 {
			char = 'O'
		}
		if !snake.IsAlive {
			char = 'X'
		}
		w.setCell(p.X, p.Y, char, color, 0)
	}

}

func (w ArenaWidget) drawPointItem() {
	p := w.state.PointItem
	w.setCell(p.X, p.Y, '*', colors["pointItem"], 0)
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
		w.putString(1, 1+i, fmt.Sprintf("Player %d: %d", i+1, len(snake.Segments)))
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
	w.arena.AddSnake(w.size.X/3, w.size.Y/3, 5, arena.EAST)
}

func (w *ArenaWidget) Run() {
	tick := time.Tick(100 * time.Millisecond)
	event := eventChannel()
	w.running = true

	for w.running {
		termbox.Clear(0, 0)
		w.Draw()
		termbox.Flush()
		select {
		case ev := <-event:
			handleEvent(ev, w.KeyMap, w.RuneMap)
		case <-tick:
			w.Tick()
		}
	}
}

func (w *ArenaWidget) Exit() {
	w.running = false
}

func NewArenaWidget(ox, oy, x, y int) ArenaWidget {
	w := ArenaWidget{offset: Position{ox, oy}, size: Position{x, y}}
	w.ResetArena()
	return w
}
