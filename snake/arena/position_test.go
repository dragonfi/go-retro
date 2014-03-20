package arena_test

import (
	"github.com/dragonfi/go-retro/snake/arena"
	"testing"
)

var positions = []arena.Position{
	{10, 20}, {-3, 1}, {3, -1}, {-10, -10},
}

var differing_positions = [][2]arena.Position{
	{{10, 20}, {10, 21}},
	{{-3, 1}, {-4, 1}},
	{{3, -1}, {15, 24}},
	{{-10, -10}, {10, 10}},
	{{8, 2}, {2, 8}},
}

func TestPositionsAreEqual(t *testing.T) {
	for _, p1 := range positions {
		p2 := arena.Position{p1.X, p1.Y}
		assertPositionsAreEqual(t, p1, p2)
	}
}

func TestPositionsDiffer(t *testing.T) {
	for _, p := range differing_positions {
		assertPositionsDiffer(t, p[0], p[1])
	}
}

func assertPositionsAreEqual(t *testing.T, p1, p2 arena.Position) {
	if !p1.Equal(p2) {
		t.Error("Positions should be equal:", p1, "\t=", p2)
	}
}

func assertPositionsDiffer(t *testing.T, p1, p2 arena.Position) {
	if p1.Equal(p2) {
		t.Error("Positions should not be equal:", p1, "\t!=", p2)
	}
}
