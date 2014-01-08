package arena

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
	Size  Position
	Snake Snake
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

func (s *Snake) moveBody() {
	for i := len(s.Segments) - 1; i > 0; i-- {
		s.Segments[i] = s.Segments[i-1]
	}
}

func (s *Snake) move() {
	s.moveBody()
	s.moveHead()
}

type arena struct {
	size  Position
	snake Snake
}

func (a arena) State() State {
	segments := make([]Position, len(a.snake.Segments))
	copy(segments, a.snake.Segments)
	return State{Size: a.size, Snake: Snake{Segments: segments, Heading: a.snake.Heading}}
}

func (a *arena) Tick() {
	a.snake.move()
}

func (a *arena) SetSnakeHeading(h Direction) {
	a.snake.Heading = h
}

func newSnake(x, y int, heading Direction, size int) Snake {
	segments := make([]Position, size, size*10)
	s := Snake{Segments: segments}
	for i := 0; i < size; i++ {
		s.Segments[i] = Position{x - i, y}
	}
	return s
}

func New(width, height int) Arena {
	s := newSnake(width/2, height/2, EAST, 5)
	a := arena{size: Position{width, height}, snake: s}
	return &a
}
