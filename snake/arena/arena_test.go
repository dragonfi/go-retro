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
	type snakeParams struct {
		x, y, size int
		heading    Direction
	}
	addSnake(t, a, width/2, height/2, 5, EAST)
	return a
}

func addSnake(t *testing.T, a Arena, x, y, size int, heading Direction) {
	old := a.State()
	index, err := a.AddSnake(x, y, size, heading)
	if err != nil {
		panic(err)
	}
	state := a.State()
	for i := range old.Snakes {
		if !state.Snakes[i].Equal(old.Snakes[i]) {
			t.Error("Old snakes are not preserved.")
		}
	}
	if len(state.Snakes) != len(old.Snakes)+1 {
		t.Error("New Snake is not appended correctly to Snakes.")
	}
	if index != len(state.Snakes)-1 {
		t.Error("AddSnake returns wrong index.")
	}
	s := state.Snakes[index]
	if s.Heading != heading {
		t.Error("Wrong direction!")
	}
	if s.Length() != size || len(s.Segments) != size {
		t.Error("Wrong snake size: Expected:", size, "Got:", s.Length())
	}
	h := s.Head()
	if h.X != x || h.Y != y {
		t.Error("Wrong position for snake head: Expected:", x, y, "Got:", h.X, h.Y)
	}
}

func checkSnakeMovementHead(t *testing.T, initial Snake, direction Direction, s Snake) {
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
	expected := Position{initial_head.X + dx, initial_head.Y + dy}
	if h.X != expected.X || h.Y != expected.Y {
		t.Error("Wrong position for snake head:", h.X, h.Y, "Expected:", expected.X, expected.Y)
		t.Error("Note: Direction:", direction)
	}
}

func checkSnakeMovementBody(t *testing.T, initial, s Snake) {
	for i := 1; i < len(s.Segments); i++ {
		if s.Segments[i] != initial.Segments[i-1] {
			t.Error("Wrong segment at position:", i, "segment:", s.Segments[i], "expected:", initial.Segments[i-1])
		}
	}
}

func checkDeadSnakeStaysPut(t *testing.T, initial, s Snake) {
	s.Equal(initial)
}

func moveSnakes(t *testing.T, a Arena, directions ...Direction) {
	if len(directions) > len(a.State().Snakes) {
		panic("Error in test: more direction arguments than snakes.")
	}
	initial_snakes := a.State().Snakes
	for i, dir := range directions {
		a.SetSnakeHeading(i, dir)
	}
	a.Tick()
	for i, s := range a.State().Snakes {
		initial := initial_snakes[i]
		direction := initial.Heading
		if i < len(directions) {
			direction = directions[i]
		}
		if s.Heading != direction {
			t.Error("Wrong direction!")
		}
		if s.Length() != initial.Length() {
			t.Error("Wrong snake size: Expected:", initial.Length(), "Got:", s.Length())
		}
		if initial.IsAlive {
			checkSnakeMovementHead(t, initial, direction, s)
			checkSnakeMovementBody(t, initial, s)
		} else {
			checkDeadSnakeStaysPut(t, initial, s)
		}
	}
}

func testSnakeMovement(t *testing.T, a Arena, directions ...Direction) {
	moveSnakes(t, a, directions...)
	s := a.State()
	if s.GameIsOver {
		t.Error("Game should not have ended yet. Head position:", s.Snakes[0].Head())
	}
}

func testSnakeMovementCausesGameOver(t *testing.T, a Arena, directions ...Direction) {
	moveSnakes(t, a, directions...)
	s := a.State()
	if !s.GameIsOver {
		t.Error("Game should have ended. Head position:", s.Snakes[0].Head())
	}
	old_state := s
	a.Tick()
	s = a.State()
	if !s.Equal(old_state) {
		t.Error("Game should not proceed further after game over.")
	}
}

func checkSnakeLength(t *testing.T, size int) {
	s := newSnake(0, 0, size, EAST)
	if s.Length() != len(s.Segments) || false {
		t.Error("Snake.Length returns wrong size: Expected:", len(s.Segments), "Got:", s.Length())
	}
}

func TestSnakeLength(t *testing.T) {
	checkSnakeLength(t, 3)
	checkSnakeLength(t, 4)
	checkSnakeLength(t, 5)
	checkSnakeLength(t, 10)
	checkSnakeLength(t, 100)
}

func TestSnakeMovement(t *testing.T) {
	a := makeArena(t, 40, 20)
	testSnakeMovement(t, a, EAST)
	testSnakeMovement(t, a, NORTH)
	testSnakeMovement(t, a, NORTH)
	testSnakeMovement(t, a, WEST)
	testSnakeMovement(t, a, SOUTH)
}

func TestSnakeMovementHitSelfAndGameOver(t *testing.T) {
	a := makeArena(t, 40, 20)
	testSnakeMovement(t, a, EAST)
	testSnakeMovement(t, a, SOUTH)
	testSnakeMovement(t, a, WEST)
	testSnakeMovementCausesGameOver(t, a, NORTH)
}

func TestSnakeMovementHitWallAndGameOver(t *testing.T) {
	width, height := 40, 20
	a := makeArena(t, width, height)
	iterations := width - 1 - a.State().Snakes[0].Head().X
	for i := 0; i < iterations; i++ {
		testSnakeMovement(t, a, EAST)
	}
	testSnakeMovementCausesGameOver(t, a, EAST)
}

func TestSnakeMovementForTwoSnakes(t *testing.T) {
	a := makeArena(t, 40, 20)
	addSnake(t, a, 30, 15, 5, EAST)
	testSnakeMovement(t, a, EAST, EAST)
	testSnakeMovement(t, a, NORTH, SOUTH)
	testSnakeMovement(t, a, NORTH, SOUTH)
	testSnakeMovement(t, a, WEST, WEST)
	testSnakeMovement(t, a, SOUTH, NORTH)
}

func TestStateCopy(t *testing.T) {
	a := makeArena(t, 40, 20)
	s1 := a.State()
	s2 := s1.Copy()
	if s1.GameIsOver != s2.GameIsOver {
		t.Fail()
	}
	if s1.PointItem != s2.PointItem {
		t.Fail()
	}
	for i, s1snake := range s1.Snakes {
		s2snake := s2.Snakes[i]
		if s1snake.Heading != s2snake.Heading {
			t.Fail()
		}
		for j := range s1snake.Segments {
			if s1snake.Segments[j] != s2snake.Segments[j] {
				t.Fail()
			}
		}
	}
	if s1.Size != s2.Size {
		t.Fail()
	}
	s1.Snakes[0].Segments[0] = Position{30, 30}
	if s1.Snakes[0].Segments[0] == s2.Snakes[0].Segments[0] {
		t.Error("Copied state should not be able to modify original data.")
	}
}

func TestSnakeMovementEatPointItemAndGrow(t *testing.T) {
	a := makeArena(t, 40, 20)
	initial := a.State().Snakes[0]

	a.(*arena).s.PointItem = Position{initial.Head().X + 1, initial.Head().Y}
	a.SetSnakeHeading(0, EAST)

	a.Tick()
	s := a.State().Snakes[0]
	if s.Heading != EAST {
		t.Error("Wrong direction!")
	}
	if s.Length() != initial.Length()+1 {
		t.Error("Wrong snake size: Expected:", initial.Length()+1, "Got:", s.Length())
	}
	if s.Head() == a.State().PointItem {
		t.Error("Point item is not eaten correctly:", a.State().PointItem)
	}
	checkSnakeMovementHead(t, initial, EAST, s)
	checkSnakeMovementBody(t, initial, s)
}

func TestDeclareGameOverWhenCannotPlaceMorePointItems(t *testing.T) {
	a := makeArena(t, 2, 1).(*arena)
	h := a.s.Snakes[0].Head()
	a.s.PointItem = Position{h.X + 1, h.Y}
	a.Tick()
	state := a.State()
	if !state.GameIsOver {
		t.Error("Game should have ended.")
	}
}

func TestSnakesAreEqual(t *testing.T) {
	s1, s2 := makeSnakes()
	if !s1.Equal(s2) {
		t.Error("Snakes should not differ.")
	}
}

func TestSnakesDifferInSegment(t *testing.T) {
	s1, s2 := makeSnakes()
	s2.Segments[4] = Position{0, 0}
	assertSnakesDiffer(t, s1, s2)
}

func TestSnakesDifferInSize(t *testing.T) {
	s1, s2 := makeSnakes()
	s2.Segments = append(s2.Segments, Position{0, 0})
	assertSnakesDiffer(t, s1, s2)
}

func TestSnakesDifferInHeading(t *testing.T) {
	s1, s2 := makeSnakes()
	s2.Heading = NORTH
	assertSnakesDiffer(t, s1, s2)
}

func makeSnakes() (Snake, Snake) {
	x, y, size, heading := 10, 15, 5, EAST
	return newSnake(x, y, size, heading), newSnake(x, y, size, heading)
}

func assertSnakesDiffer(t *testing.T, s1, s2 Snake) {
	if s1.Equal(s2) {
		t.Error("Snakes should differ:", s1, s2)
	}
}

func TestStatesAreEqual(t *testing.T) {
	s1, s2 := makeStates(t)
	if !s1.Equal(s2) {
		t.Error("States should not differ:", s1, s2)
	}
}

func TestStatesDifferInSize(t *testing.T) {
	s1, s2 := makeStates(t)

	s2.Size.X += 1
	assertStatesDiffer(t, s1, s2)

	s2.Size.X = s1.Size.X
	s2.Size.Y -= 1
	assertStatesDiffer(t, s1, s2)
}

func TestStatesDifferInGameOverFlag(t *testing.T) {
	s1, s2 := makeStates(t)
	s2.GameIsOver = true
	assertStatesDiffer(t, s1, s2)
}

func TestStatesDifferInPointItemPosition(t *testing.T) {
	s1, s2 := makeStates(t)

	s2.PointItem.X += 1
	assertStatesDiffer(t, s1, s2)

	s2.PointItem.X = s1.PointItem.X
	s2.PointItem.Y += 1
	assertStatesDiffer(t, s1, s2)
}

func TestStatesDifferInSnakes(t *testing.T) {
	s1, s2 := makeStates(t)
	s2.Snakes[0].Heading = NORTH
	assertStatesDiffer(t, s1, s2)
}

func makeStates(t *testing.T) (State, State) {
	x, y := 30, 15
	a := makeArena(t, x, y)
	return a.State(), a.State()
}

func assertStatesDiffer(t *testing.T, s1, s2 State) {
	if s1.Equal(s2) {
		t.Error("States should differ:", s1, s2)
	}
}

func TestValidPointItemPositions(t *testing.T) {
	width := 40
	height := 20
	a := makeArena(t, width, height).(*arena)
	valid_positions := []Position{
		{0, 0}, {0, height - 1}, {width - 1, 0}, {width - 1, height - 1},
		{1, 1}, {21, 15}, {17, 18},
	}
	for _, position := range valid_positions {
		if !a.isValidPlacementPosition(position) {
			t.Error("Point item position should be valid:", position)
		}
	}
}

func TestInvalidPointItemPositionsOutOfBounds(t *testing.T) {
	width := 40
	height := 20
	a := makeArena(t, width, height).(*arena)
	invalid_positions := []Position{
		{-1, 0}, {0, -1}, {width, 0}, {0, width}, {width, height},
		{-54, -36}, {-32, 100}, {-32, 11}, {11, -30},
	}
	for _, position := range invalid_positions {
		if a.isValidPlacementPosition(position) {
			t.Error("Point item position should be invalid:", position)
		}
	}
}

func TestInvalidPointItemPositionsOnSnakes(t *testing.T) {
	width := 40
	height := 20
	a := makeArena(t, width, height).(*arena)
	for _, snake := range a.State().Snakes {
		for _, position := range snake.Segments {
			if a.isValidPlacementPosition(position) {
				t.Error("Point item position should be invalid:", position)
			}
		}
	}
}

func TestNewSnakeHeadCannotBeAtInvalidPosition(t *testing.T) {
	width, height := 40, 20
	a := makeArena(t, width, height)
	a.AddSnake(20, 10, 5, EAST)
	i, err := a.AddSnake(20, 10, 5, EAST)
	if i != -1 || err == nil {
		t.Error("Same snake head position should be invalid.")
	}
	i, err = a.AddSnake(19, 10, 5, EAST)
	if i != -1 || err == nil {
		t.Error("Snake head position inside other snake should be invalid.")
	}
	i, err = a.AddSnake(21, 10, 5, EAST)
	if i != -1 || err == nil {
		t.Error("Placing new segment on another snake's head should be invalid.")
	}
	i, err = a.AddSnake(-1, 10, 5, EAST)
	if i != -1 || err == nil {
		t.Error("Placing snake head outside arena boundary should be invalid.")
	}
	i, err = a.AddSnake(10, -1, 5, EAST)
	if i != -1 || err == nil {
		t.Error("Placing snake head outside arena boundary should be invalid.")
	}
	i, err = a.AddSnake(width+1, 10, 5, EAST)
	if i != -1 || err == nil {
		t.Error("Placing snake head outside arena boundary should be invalid.")
	}
	i, err = a.AddSnake(3, height+1, 5, EAST)
	if i != -1 || err == nil {
		t.Error("Placing snake head outside arena boundary should be invalid.")
	}
	if len(a.State().Snakes) != 1 {
		t.Error("Bad snakes should not be added at all.")
	}
}

func TestReverseMotionSuicideIsInvalid(t *testing.T) {
	a := makeArena(t, 40, 20)
	testSnakeMovement(t, a, EAST)
	a.SetSnakeHeading(0, WEST)
	a.Tick()
	testSnakeMovement(t, a, EAST)

	testSnakeMovement(t, a, NORTH)
	a.SetSnakeHeading(0, SOUTH)
	a.Tick()
	testSnakeMovement(t, a, NORTH)

	testSnakeMovement(t, a, WEST)
	a.SetSnakeHeading(0, EAST)
	a.Tick()
	testSnakeMovement(t, a, WEST)

	testSnakeMovement(t, a, SOUTH)
	a.SetSnakeHeading(0, NORTH)
	a.Tick()
	testSnakeMovement(t, a, SOUTH)
}

func TestGameIsOnlyOverWhenAllSnakeDies(t *testing.T) {
	a := makeArena(t, 40, 20)
	addSnake(t, a, 30, 15, 5, EAST)
	testSnakeMovement(t, a, EAST, EAST)
	testSnakeMovement(t, a, NORTH, NORTH)
	testSnakeMovement(t, a, NORTH, WEST)
	testSnakeMovement(t, a, WEST, SOUTH)
	testSnakeMovement(t, a, SOUTH)
	testSnakeMovementCausesGameOver(t, a, EAST)
	h := a.State().Snakes[1].Head()
	if h.X != 30 || h.Y != 15 {
		t.Error("Dead snakes should not move.")
	}
}
