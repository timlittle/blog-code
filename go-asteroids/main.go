package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 400
)

func init() {
	//Built go function which runs before main()

	// Setup the raylib window
	rl.InitWindow(screenWidth, screenHeight, "Asteroids")
	rl.SetTargetFPS(60)
}

func draw() {
	rl.BeginDrawing()

	// Set the background to black
	rl.ClearBackground(rl.Black)

	// Draw the score to the screen
	rl.DrawText("Score 0", 10, 10, 20, rl.Gray)

	rl.EndDrawing()

}

func update() {
	//TODO:  Update the state

}

func deinit() {
	rl.CloseWindow()
}

func main() {
	// When the main function ends, call the deinit() function
	defer deinit()

	// Continue the loop until the window is closed or ESC is pressed
	for !rl.WindowShouldClose() {
		draw()
		update()
	}
}
