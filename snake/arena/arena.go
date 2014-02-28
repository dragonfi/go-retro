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
	Size       Position
	Snakes     []Snake
	PointItem  Position
	GameIsOver bool
}

func (s State) Equal(other State) bool {
	if s.Size != other.Size {
		return false
	}
	if !s.Snakes[0].Equal(other.Snakes[0]) {
		return false
	}
	if s.GameIsOver != other.GameIsOver {
		return false
	}
	if s.PointItem != other.PointItem {
		return false
	}
	return true
}

type Snake struct {
	Segments []Position
	Heading  Direction
}

func (s Snake) Equal(other Snake) bool {
	if s.Heading != other.Heading {
		return false
	}
	if s.Length() != other.Length() {
		return false
	}
	for i := range s.Segments {
		if s.Segments[i] != other.Segments[i] {
			return false
		}
	}
	return true
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

func (s Snake) Copy() Snake {
	segments := make([]Position, len(s.Segments))
	copy(segments, s.Segments)
	return Snake{Segments: segments, Heading: s.Heading}
}

type arena struct {
	size       Position
	snake      Snake
	pointItem  Position
	gameIsOver bool
}

func (a arena) State() State {
	snakes := make([]Snake, 1, 1)
	snakes[0] = a.snake.Copy()
	return State{
		Size:       a.size,
		Snakes:     snakes,
		PointItem:  a.pointItem,
		GameIsOver: a.gameIsOver,
	}
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
	if a.gameIsOver {
		return
	}
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

func (a arena) getValidPositions() []Position {
	valid_positions := make([]Position, 0, a.size.X*a.size.Y)
	for i := 0; i < a.size.X; i++ {
		for j := 0; j < a.size.Y; j++ {
			p := Position{i, j}
			if a.isValidPointItemPosition(p) {
				valid_positions = append(valid_positions, p)
			}
		}
	}
	return valid_positions

}

func (a *arena) setRandomPositionForPointItem() {
	valid_positions := a.getValidPositions()
	if len(valid_positions) == 0 {
		a.endGame()
	} else {
		a.pointItem = valid_positions[rand.Intn(len(valid_positions))]
	}
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
	a.setRandomPositionForPointItem()
	return &a
}
