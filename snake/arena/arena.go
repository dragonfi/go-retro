package arena

import (
	"errors"
	"math/rand"
)

type Arena interface {
	State() State
	Tick()
	SetSnakeHeading(snake int, h Direction)
	AddSnake(x, y, size int, h Direction) (snake int, err error)
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

// TODO: consider providing deep Copy for state.

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

func (s State) Copy() State {
	return State{
		Size:       s.Size,
		Snakes:     s.copySnakes(),
		PointItem:  s.PointItem,
		GameIsOver: s.GameIsOver,
	}
}

func (s State) copySnakes() []Snake {
	snakes := make([]Snake, len(s.Snakes))
	for i, snake := range s.Snakes {
		snakes[i] = snake.Copy()
	}
	return snakes
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
	s State
}

func (a arena) State() State {
	return a.s.Copy()
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
	if p.X < 0 || p.X >= a.s.Size.X || p.Y < 0 || p.Y >= a.s.Size.Y {
		return false
	}
	return true
}

func (a *arena) endGame() {
	a.s.GameIsOver = true
}

func (a *arena) Tick() {
	if a.s.GameIsOver {
		return
	}
	for id := range a.s.Snakes {
		snake := &a.s.Snakes[id]
		snake.extrude()
		if snake.Head() == a.s.PointItem {
			a.setRandomPositionForPointItem()
		} else {
			snake.contractBody()
		}

		for other_id, other_snake := range a.s.Snakes {
			if inSequence(snake.Head(), other_snake.Segments) {
				if id != other_id {
					a.endGame()
				} else if inSequence(snake.Head(), other_snake.Segments[1:]) {
					a.endGame()
				}
			}
		}

		if !a.insideArena(snake.Head()) {
			a.endGame()
		}
	}
}

func isOpposingDirections(h1, h2 Direction) bool {
	if h1 == EAST && h2 == WEST {
		return true
	}
	if h1 == WEST && h2 == EAST {
		return true
	}
	if h1 == NORTH && h2 == SOUTH {
		return true
	}
	if h1 == SOUTH && h2 == NORTH {
		return true
	}
	return false
}

func (a *arena) SetSnakeHeading(snake int, h Direction) {
	if isOpposingDirections(a.s.Snakes[snake].Heading, h) {
		return
	}
	a.s.Snakes[snake].Heading = h
}

func (a arena) isValidPlacementPosition(p Position) bool {
	if p.X < 0 || p.X >= a.s.Size.X {
		return false
	}
	if p.Y < 0 || p.Y >= a.s.Size.Y {
		return false
	}
	for _, snake := range a.s.Snakes {
		if inSequence(p, snake.Segments) {
			return false
		}
	}
	return true
}

func (a arena) getValidPositions() []Position {
	valid_positions := make([]Position, 0, a.s.Size.X*a.s.Size.Y)
	for i := 0; i < a.s.Size.X; i++ {
		for j := 0; j < a.s.Size.Y; j++ {
			p := Position{i, j}
			if a.isValidPlacementPosition(p) {
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
		a.s.PointItem = valid_positions[rand.Intn(len(valid_positions))]
	}
}

func (a *arena) AddSnake(x, y, size int, heading Direction) (int, error) {
	if !a.isValidPlacementPosition(Position{x, y}) {
		return -1, errors.New("Invalid position for snake head.")
	}
	new_snake := newSnake(x, y, size, heading)
	for _, snake := range a.s.Snakes {
		if inSequence(snake.Head(), new_snake.Segments) {
			return -1, errors.New("Snake segment makes another snake head position invalid.")
		}
	}
	a.s.Snakes = append(a.s.Snakes, new_snake)
	return len(a.s.Snakes) - 1, nil
}

func newSnake(x, y, size int, heading Direction) Snake {
	if heading != EAST {
		panic("TODO: Other headings are not implemented.")
	}
	if size < 0 {
		panic("Size should be positive.")
	}
	segments := make([]Position, size, size*10)
	s := Snake{Segments: segments}
	for i := 0; i < size; i++ {
		s.Segments[i] = Position{x - i, y}
	}
	return s
}

func New(width, height int) Arena {
	a := arena{s: State{Size: Position{width, height}}}
	a.setRandomPositionForPointItem()
	return &a
}
