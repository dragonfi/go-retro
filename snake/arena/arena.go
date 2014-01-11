package arena

import "math/rand"

type Arena interface {
	State() State
	Tick()
	SetSnakeHeading(h Direction)
}

type Direction int

const (
	EAST = Direction(iota)
	NORTH
	WEST
	SOUTH
)

type Position struct {
	X, Y int
}

type State struct {
	Size      Position
	Snake     Snake
	PointItem Position
	GameIsOver  bool
}

type Snake struct {
	Segments []Position
	Heading  Direction
}

func (s Snake) Head() Position {
	return s.Segments[0]
}

func (s Snake) Length() int {
	return len(s.Segments)
}

func (s *Snake) moveHead() {
	switch s.Heading {
	case EAST:
		s.Segments[0].X += 1
	case NORTH:
		s.Segments[0].Y -= 1
	case WEST:
		s.Segments[0].X -= 1
	case SOUTH:
		s.Segments[0].Y += 1
	}
}

func (s *Snake) extrudeBody() {
	s.Segments = s.Segments[:len(s.Segments)+1]
	for i := len(s.Segments) - 1; i > 0; i-- {
		s.Segments[i] = s.Segments[i-1]
	}
}

func (s *Snake) contractBody() {
	s.Segments = s.Segments[:len(s.Segments)-1]
}

func (s *Snake) extrude() {
	s.extrudeBody()
	s.moveHead()
}

type arena struct {
	size      Position
	snake     Snake
	pointItem Position
	gameIsOver bool
}

func (a arena) State() State {
	segments := make([]Position, len(a.snake.Segments))
	copy(segments, a.snake.Segments)
	snake := Snake{Segments: segments, Heading: a.snake.Heading}
	return State{Size: a.size, Snake: snake, PointItem: a.pointItem, GameIsOver: a.gameIsOver}
}

func inSequence(p Position, sequence []Position) bool {
	for _, item := range sequence {
		if p == item {
			return true
		}
	}
	return false
}

func (a arena) insideArena(p Position) bool {
	if p.X < 0 || p.X >= a.size.X || p.Y < 0 || p.Y >= a.size.Y {
		return false
	}
	return true
}

func (a *arena) endGame() {
	a.gameIsOver = true
}

func (a *arena) Tick() {
	a.snake.extrude()
	if a.snake.Head() == a.pointItem {
		a.setRandomPositionForPointItem()
	} else {
		a.snake.contractBody()
	}

	if inSequence(a.snake.Head(), a.snake.Segments[1:]) {
		a.endGame()
	}

	if !a.insideArena(a.snake.Head()) {
		a.endGame()
	}
}

func (a *arena) SetSnakeHeading(h Direction) {
	a.snake.Heading = h
}

func (a arena) isValidPointItemPosition(p Position) bool {
	if p.X < 0 || p.X >= a.size.X {
		return false
	}
	if p.Y < 0 || p.Y >= a.size.Y {
		return false
	}
	if inSequence(p, a.snake.Segments) {
		return false
	}
	return true
}

func (a *arena) setRandomPositionForPointItem() {
	for counter := 0; counter < 1000; counter++ {
		newPointItem := Position{rand.Intn(a.size.X), rand.Intn(a.size.Y)}
		if a.isValidPointItemPosition(newPointItem) {
			a.pointItem = newPointItem
			return
		}
	}
	panic("Cannot find a place to put point item. Maybe I should see if there are places available at all...")
}

func newSnake(x, y int, size int) Snake {
	segments := make([]Position, size, size*10)
	s := Snake{Segments: segments}
	for i := 0; i < size; i++ {
		s.Segments[i] = Position{x - i, y}
	}
	return s
}

func New(width, height int) Arena {
	s := newSnake(width/2, height/2, 5)
	a := arena{size: Position{width, height}, snake: s}
	return &a
}
