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

type arena struct {
	s State
}

func (a arena) State() State {
	return a.s.Copy()
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

func (a *arena) endGameIfAllSnakesAreDead() {
	for _, snake := range a.s.Snakes {
		if snake.IsAlive {
			return
		}
	}
	a.endGame()
}

func (a *arena) killSnake(snake int) {
	a.s.Snakes[snake].IsAlive = false
	a.endGameIfAllSnakesAreDead()
}

func (a *arena) Tick() {
	if a.s.GameIsOver {
		return
	}
	for id := range a.s.Snakes {
		snake := &a.s.Snakes[id]
		if !snake.IsAlive {
			continue
		}
		snake.extrude()
		if snake.Head() == a.s.PointItem {
			a.setRandomPositionForPointItem()
		} else {
			snake.contractBody()
		}

		for other_id, other_snake := range a.s.Snakes {
			if inSequence(snake.Head(), other_snake.Segments) {
				if id != other_id {
					a.killSnake(id)
				} else if inSequence(snake.Head(), other_snake.Segments[1:]) {
					a.killSnake(id)
				}
			}
		}

		if !a.insideArena(snake.Head()) {
			a.killSnake(id)
		}
	}
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

func New(width, height int) Arena {
	if width < 0 || height < 0 {
		panic("Arena width and height must be positive.")
	}
	a := arena{s: State{Size: Position{width, height}}}
	a.setRandomPositionForPointItem()
	return &a
}
