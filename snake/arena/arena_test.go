package arena_test

import (
	"../arena"
	"testing"
)

func TestArenaCreation(t *testing.T) {
	makeArena(t, 40, 20)
	makeArena(t, 50, 30)
}

func makeArena(t *testing.T, width, height int) arena.Arena {
	a := arena.New(width, height)
	state := a.State()
	if state.Size.X != width || state.Size.Y != height {
		t.Error("Wrong width or height. Expected:", width, height, "Got:", state.Size.X, state.Size.Y)
	}
	s := state.Snake
	if s.Heading != arena.EAST {
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

func testSnakeMovementHead(t *testing.T, initial arena.Snake, direction arena.Direction, s arena.Snake) {
	initial_head := initial.Head()
	h := s.Head()
	var dx, dy int
	switch direction {
	case arena.EAST:
		dx = 1
	case arena.NORTH:
		dy = -1
	case arena.WEST:
		dx = -1
	case arena.SOUTH:
		dy = 1
	}
	if h.X != initial_head.X+dx || h.Y != initial_head.Y+dy {
		t.Error("Wrong position for snake head:", h.X, h.Y, "Expected:", initial_head.X+dx, initial_head.Y+dy)
		t.Error("Note: Direction:", direction)
	}
}

func testSnakeMovementBody(t *testing.T, initial, s arena.Snake) {
	for i := 1; i < len(s.Segments); i++ {
		if s.Segments[i] != initial.Segments[i-1] {
			t.Error("Wrong segment at position:", i, "segment:", s.Segments[i], "expected:", initial.Segments[i-1])
		}
	}
}

func testSnakeMovement(t *testing.T, a arena.Arena, direction arena.Direction) {
	initial := a.State().Snake
	a.SetSnakeHeading(direction)
	a.Tick()
	s := a.State().Snake
	if s.Heading != direction {
		t.Error("Wrong direction!")
	}
	if s.Length() != 5 || len(s.Segments) != 5 {
		t.Error("Wrong snake size: Expected:", 5, "Got:", s.Length())
	}
	testSnakeMovementHead(t, initial, direction, s)
	testSnakeMovementBody(t, initial, s)
}

func TestSnakeMovement(t *testing.T) {
	const width = 40
	const height = 20
	a := makeArena(t, width, height)
	testSnakeMovement(t, a, arena.EAST)
	testSnakeMovement(t, a, arena.NORTH)
	testSnakeMovement(t, a, arena.WEST)
	testSnakeMovement(t, a, arena.SOUTH)
}
