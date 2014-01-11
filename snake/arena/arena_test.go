package arena

import (
	"testing"
)

func TestArenaCreation(t *testing.T) {
	makeArena(t, 40, 20)
	makeArena(t, 50, 30)
}

func makeArena(t *testing.T, width, height int) Arena {
	a := New(width, height)
	state := a.State()
	if state.Size.X != width || state.Size.Y != height {
		t.Error("Wrong width or height. Expected:", width, height, "Got:", state.Size.X, state.Size.Y)
	}
	s := state.Snake
	if s.Heading != EAST {
		t.Error("Wrong direction!")
	}
	if s.Length() != 5 || len(s.Segments) != 5 {
		t.Error("Wrong snake size: Expected:", 5, "Got:", s.Length())
	}
	h := s.Head()
	if h.X != width/2 || h.Y != height/2 {
		t.Error("Wrong position for snake head: Expected:", width/2, height/2, "Got:", h.X, h.Y)
	}
	return a
}

func testSnakeMovementHead(t *testing.T, initial Snake, direction Direction, s Snake) {
	initial_head := initial.Head()
	h := s.Head()
	var dx, dy int
	switch direction {
	case EAST:
		dx = 1
	case NORTH:
		dy = -1
	case WEST:
		dx = -1
	case SOUTH:
		dy = 1
	}
	if h.X != initial_head.X+dx || h.Y != initial_head.Y+dy {
		t.Error("Wrong position for snake head:", h.X, h.Y, "Expected:", initial_head.X+dx, initial_head.Y+dy)
		t.Error("Note: Direction:", direction)
	}
}

func testSnakeMovementBody(t *testing.T, initial, s Snake) {
	for i := 1; i < len(s.Segments); i++ {
		if s.Segments[i] != initial.Segments[i-1] {
			t.Error("Wrong segment at position:", i, "segment:", s.Segments[i], "expected:", initial.Segments[i-1])
		}
	}
}

func testSnakeMovement(t *testing.T, a Arena, direction Direction) {
	initial := a.State().Snake
	a.SetSnakeHeading(direction)
	a.Tick()
	s := a.State().Snake
	if s.Heading != direction {
		t.Error("Wrong direction!")
	}
	if s.Length() != initial.Length() {
		t.Error("Wrong snake size: Expected:", initial.Length(), "Got:", s.Length())
	}
	testSnakeMovementHead(t, initial, direction, s)
	testSnakeMovementBody(t, initial, s)
}

func testSnakeLength(t *testing.T, size int) {
	s := newSnake(0, 0, size)
	if s.Length() != len(s.Segments) || false {
		t.Error("Snake.Length returns wrong size: Expected:", len(s.Segments), "Got:", s.Length())
	}
}

func TestSnakeLength(t *testing.T) {
	testSnakeLength(t, 3)
	testSnakeLength(t, 4)
	testSnakeLength(t, 5)
	testSnakeLength(t, 10)
	testSnakeLength(t, 100)
}

func TestSnakeMovement(t *testing.T) {
	a := makeArena(t, 40, 20)
	testSnakeMovement(t, a, EAST)
	testSnakeMovement(t, a, NORTH)
	testSnakeMovement(t, a, NORTH)
	testSnakeMovement(t, a, WEST)
	testSnakeMovement(t, a, SOUTH)
	//testSnakeMovementHitWallAndDie(t, a, EAST)
	//testSnakeMovementHitSelfAndDie(t, a, EAST)
}

func TestState(t *testing.T) {
	a := makeArena(t, 40, 20).(*arena)
	s := a.State()
	if s.PointItem != a.pointItem {
		t.Fail()
	}
	if s.Snake.Heading != a.snake.Heading {
		t.Fail()
	}
	for i := range s.Snake.Segments {
		if s.Snake.Segments[i] != a.snake.Segments[i] {
			t.Fail()
		}
	}
	if s.Size != a.size {
		t.Fail()
	}
}

func TestSnakeMovementEatPointItemAndGrow(t *testing.T) {
	a := makeArena(t, 40, 20)
	initial := a.State().Snake

	a.(*arena).pointItem = Position{initial.Head().X + 1, initial.Head().Y}
	a.SetSnakeHeading(EAST)

	a.Tick()
	s := a.State().Snake
	if s.Heading != EAST {
		t.Error("Wrong direction!")
	}
	if s.Length() != initial.Length()+1 {
		t.Error("Wrong snake size: Expected:", initial.Length()+1, "Got:", s.Length())
	}
	if s.Head() == a.State().PointItem {
		t.Error("Point item is not eaten correctly:", a.State().PointItem)
	}
	testSnakeMovementHead(t, initial, EAST, s)
	testSnakeMovementBody(t, initial, s)
}

func TestValidPointItemPositions(t *testing.T) {
	width := 40
	height := 20
	a := makeArena(t, width, height).(*arena)
	valid_positions := []Position {
		{0, 0}, {0, height-1}, {width-1, 0}, {width-1, height-1},
		{1, 1}, {21, 15}, {17, 18},
	}
	for _, position := range valid_positions {
		if !a.isValidPointItemPosition(position) {
			t.Error("Point item position should be valid:", position)
		}
	}
}

func TestInvalidPointItemPositionsOutOfBounds(t *testing.T) {
	width := 40
	height := 20
	a := makeArena(t, width, height).(*arena)
	invalid_positions := []Position {
		{-1, 0}, {0, -1}, {width, 0}, {0, width}, {width, height},
		{-54, -36}, {-32, 100}, {-32, 11}, {11, -30},
	}
	for _, position := range invalid_positions {
		if a.isValidPointItemPosition(position) {
			t.Error("Point item position should be invalid:", position)
		}
	}
}

func TestInvalidPointItemPositionsOnSnake(t *testing.T) {
	width := 40
	height := 20
	a := makeArena(t, width, height).(*arena)
	for _, position := range a.snake.Segments {
		if a.isValidPointItemPosition(position) {
			t.Error("Point item position should be invalid:", position)
		}
	}
}
