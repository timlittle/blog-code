package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Game struct {
	Width    int
	Height   int
	State    [][]Droplet // 2D grid of water droplets [y][x]
	tileSize int
}

func (g *Game) Draw() {
	// Loop through each cell and call the cells Draw() method
	for y := range g.State {
		for x := 0; x < len(g.State[y]); x++ {
			g.State[y][x].Draw(x, y, g.tileSize)
		}
	}
}

type Droplet struct {
	volume float64 // How much water this cell contains (0.0 to 1.0)
	size   int
}

func (d *Droplet) Draw(x, y, tileSize int) {
	// Convert grid coordinates to pixel coordinates
	pixelX := x * tileSize
	pixelY := y * tileSize

	if d.volume > 0 {
		// Draw the blue water rectangle
		rl.DrawRectangle(int32(pixelX), int32(pixelY), int32(tileSize), int32(tileSize), rl.Blue)
	}
}

func NewGame(width, height, tileSize int) *Game {
	// Create a new game
	g := &Game{Width: width, Height: height, tileSize: tileSize}

	// Create the new game state
	// divide pixel dimensions by tile size to get grid size
	g.State = CreateGameState(g.Width/g.tileSize, g.Height/g.tileSize, tileSize)

	return g
}

func CreateGameState(newWidth, newHeight, tileSize int) [][]Droplet {
	// Create a new game state
	newState := make([][]Droplet, newHeight)

	// Loop through each row
	for y := range newHeight {

		// Create the columns
		newState[y] = make([]Droplet, newWidth)

		// Loop through each cell and create a Droplet
		for x := range newState[y] {
			newState[y][x] = Droplet{
				size: tileSize,
			}
		}
	}
	return newState
}

func main() {
	// Setup the new game
	var game = NewGame(800, 400, 10)

	// Initialize Raylib graphics window
	rl.InitWindow(int32(game.Width), int32(game.Height), "Water simulation")
	defer rl.CloseWindow()

	// Create a single water droplet and add it to the screen
	droplet := Droplet{size: game.tileSize, volume: 1.0}
	game.State[100/game.tileSize][400/game.tileSize] = droplet

	// Setup the frame per second rate
	rl.SetTargetFPS(20)

	// Main loop
	for !rl.WindowShouldClose() {

		// Begind to draw and set the background to black
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Draw the game state
		game.Draw()

		rl.EndDrawing()
	}

}
