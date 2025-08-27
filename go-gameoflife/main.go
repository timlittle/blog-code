package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Game struct {
	Height int
	Width  int
}

func main() {
	// Create the game metadata and state holding object
	game := Game{Width: 800, Height: 400}

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

		// Draw Hello world
		rl.DrawText("Hello world!", 350, 200, 20, rl.RayWhite)

		// End the drawing
		rl.EndDrawing()
	}
}
