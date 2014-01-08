package main

import (
	"testing"
)

func TestShortSnakeCreation(t *testing.T) {
	head := Position{0,0}
	length := 1
	snake := NewSnake(head, length)
	if snake[0] != head {
		t.Error(snake[0], "!=", head)
	}
	if len(snake) != 1 {
		t.Error("Snake is of wrong length:", snake, "Expected:", 1)
	}
}

func TestLongSnakeCreation(t *testing.T) {
	head := Position{3, 3}
	length := 6
	snake := NewSnake(head, length)

	expected := make([]Position, length)
	for i:= 0; i<length; i++ {
		expected[i] = Position{head.X, head.Y-i}
	}

	for i := range snake {
		if (snake[i] != expected[i]) {
			t.Error(snake[i], "!=", expected[i])
		}
	}
}
