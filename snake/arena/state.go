package arena

type Direction int

const (
	EAST = Direction(iota)
	NORTH
	WEST
	SOUTH
)

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
	IsAlive  bool
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
	return Snake{Segments: segments, Heading: s.Heading, IsAlive: s.IsAlive}
}

func newSnake(x, y, size int, heading Direction) Snake {
	if heading != EAST {
		panic("TODO: Other headings are not implemented.")
	}
	if size < 0 {
		panic("Size should be positive.")
	}
	segments := make([]Position, size, size*10)
	s := Snake{Segments: segments, IsAlive: true}
	for i := 0; i < size; i++ {
		s.Segments[i] = Position{x - i, y}
	}
	return s
}

func inSequence(p Position, sequence []Position) bool {
	for _, item := range sequence {
		if p == item {
			return true
		}
	}
	return false
}
