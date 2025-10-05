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

	// Load textures
	texBackground = rl.LoadTexture("resources/space_background.png")
}

var (
	texBackground rl.Texture2D
)

func draw() {
	rl.BeginDrawing()

	// Set the background to a nebula
	bgSource := rl.Rectangle{X: 0, Y: 0, Width: float32(texBackground.Width), Height: float32(texBackground.Height)}
	bgDest := rl.Rectangle{X: 0, Y: 0, Width: screenWidth, Height: screenHeight}
	rl.DrawTexturePro(texBackground, bgSource, bgDest, rl.Vector2{X: 0, Y: 0}, 0, rl.White)

	// Draw the score to the screen
	rl.DrawText("Score 0", 10, 10, 20, rl.Gray)

	rl.EndDrawing()

}

func update() {
	//TODO:  Update the state

}

func deinit() {
	rl.CloseWindow()

	// Unload textures when the game closes
	rl.UnloadTexture(texBackground)
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
