package arena_test

import (
	"testing"
	"../arena"
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

func testSnakeState(t *testing.T, initial arena.Snake, direction arena.Direction, s arena.Snake) {
	initial_head := initial.Head()
	if s.Heading != direction {
		t.Error("Wrong direction!")
	}
	if s.Length() != 5 || len(s.Segments) != 5 {
		t.Error("Wrong snake size: Expected:", 5, "Got:", s.Length())
	}
	h := s.Head()
	if h.X != initial_head.X + 1 || h.Y != initial_head.Y {
		t.Error("Wrong position for snake head:", h.X, h.Y, "Expected:", initial_head.X+1, initial_head.Y)
		t.Error("Note: Direction:", direction)
	}
}

/*
func TestSnakeMovement(t *testing.T) {
	const width = 40
	const height = 20
	a := makeArena(t, width, height)
	var initial = a
	initial_snake := a.Snake
	t.Log( initial_snake.Segment[0].X)
	t.Log( initial_snake.Head())

	a.Tick()
	t.Log(initial)
	t.Log(a)
	t.Log(initial_snake.Head())

	testSnakeState(t, initial_snake, arena.EAST, a.Snake)


	initial_snake = a.Snake
	a.Snake.Heading = arena.NORTH
	a.Tick()
	testSnakeState(t, initial_snake, arena.NORTH, a.Snake)
}
*/

