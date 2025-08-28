package main

import (
	"reflect"
	"testing"
)

func TestCountNeighbours(t *testing.T) {
	testCases := []struct {
		name string
		x    int
		y    int
		want int
	}{
		{name: "Four neighbours", x: 1, y: 1, want: 4},
		{name: "Three neighbours - Bottom bound", x: 1, y: 3, want: 3},
		{name: "Three neighbours - Top bound", x: 2, y: 0, want: 3},
		{name: "Four neighbours - Right bound", x: 3, y: 2, want: 4},
		{name: "One neighbour - Left bound", x: 0, y: 2, want: 1},
	}
	var gameState = [][]int{
		{0, 1, 0, 0},
		{0, 0, 1, 1},
		{0, 1, 1, 0},
		{0, 0, 1, 0},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := CountNeighbours(tt.x, tt.y, gameState)

			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})

	}

}

func TestIsCellAlive(t *testing.T) {
	testCases := []struct {
		name       string
		want       int
		current    int
		neighbours int
	}{
		{name: "Live cell should not live with <2", neighbours: 1, current: 1, want: 0},
		{name: "Live cell should live with 2", neighbours: 2, current: 1, want: 1},
		{name: "Live cells should live with 3", neighbours: 3, current: 1, want: 1},
		{name: "Dead cells should live with 3", neighbours: 3, current: 0, want: 1},
		{name: "Live cell should not live with >3", neighbours: 4, current: 1, want: 0},
		{name: "Dead cells should not live if already dead", neighbours: 0, want: 0},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCellAlive(tt.current, tt.neighbours)

			if got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestUpdateState(t *testing.T) {
	expectedGameState := [][]int{
		{0, 0, 1, 0},
		{0, 0, 0, 1},
		{0, 1, 0, 0},
		{0, 1, 1, 0},
	}
	var game = Game{Width: 40, Height: 40, tileSize: 10}
	game.State = [][]int{
		{0, 1, 0, 0},
		{0, 0, 1, 1},
		{0, 1, 1, 0},
		{0, 0, 1, 0},
	}
	game.Update()

	if !reflect.DeepEqual(game.State, expectedGameState) {
		t.Errorf("got %v, want %v", game.State, expectedGameState)
	}
}
