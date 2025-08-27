package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Game struct {
	Height   int
	Width    int
	tileSize int
	State    [][]int
}

func (g *Game) Draw() {
	// Loop through all of the rows
	for y := range g.State {

		// Loop through all of the columes
		for x := 0; x < len(g.State[y]); x++ {

			// If we have marked the column as a 1, draw it as white
			if g.State[y][x] == 1 {

				// We will need to scale our blocks to the size of the window
				pixelX := x * g.tileSize
				pixelY := y * g.tileSize

				// Draw the block to the screen
				rl.DrawRectangle(int32(pixelX), int32(pixelY), int32(g.tileSize), int32(g.tileSize), rl.RayWhite)
			}
		}
	}
}

func NewGame(width, height, tileSize int) *Game {
	g := &Game{Width: width, Height: height, tileSize: tileSize}
	g.State = [][]int{
		{0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 1, 0, 0},
		{1, 1, 1, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	return g
}

func main() {
	// Create the game metadata and state holding object
	var game = NewGame(800, 400, 80)

	// Create the Raylib window using the state
	rl.InitWindow(int32(game.Width), int32(game.Height), "Game of life")

	// Close the window at the end of the progrm
	defer rl.CloseWindow()

	// We dont need a high FPS for the game, so 10 should be enough
	rl.SetTargetFPS(10)

	// Loop until the window needs to close
	for !rl.WindowShouldClose() {
		// Starting drawing to the canvas
		rl.BeginDrawing()

		// Create a black background
		rl.ClearBackground(rl.Black)

		// Draw the game state
		game.Draw()

		// End the drawing
		rl.EndDrawing()
	}
}
