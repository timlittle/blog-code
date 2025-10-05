package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 400
	tileSize     = 64
)

var (
	texTiles      rl.Texture2D
	texBackground rl.Texture2D
	spriteRec     rl.Rectangle
	player        Player
)

type Player struct {
	position     rl.Vector2
	speed        rl.Vector2
	size         rl.Vector2
	acceleration float32
	rotation     float32
	isBoosting   bool
}

func (p *Player) Draw() {
	destTexture := rl.Rectangle{X: p.position.X, Y: p.position.Y, Width: p.size.X, Height: p.size.Y}
	rl.DrawTexturePro(
		texTiles,
		spriteRec,
		destTexture,
		rl.Vector2{X: p.size.X / 2, Y: p.size.Y / 2},
		p.rotation,
		rl.White,
	)
}

func initGame() {
	player = Player{
		position:     rl.Vector2{X: 400, Y: 200},
		speed:        rl.Vector2{X: 0.0, Y: 0.0},
		size:         rl.Vector2{X: tileSize, Y: tileSize},
		rotation:     0.0,
		acceleration: 0.0,
		isBoosting:   false,
	}
}

func init() {
	//Built go function which runs before main()

	// Setup the raylib window
	rl.InitWindow(screenWidth, screenHeight, "Asteroids")
	rl.SetTargetFPS(60)

	// Load textures
	texTiles = rl.LoadTexture("resources/tilesheet.png")
	texBackground = rl.LoadTexture("resources/space_background.png")

	spriteRec = rl.Rectangle{X: tileSize * 0, Y: tileSize * 2, Width: tileSize, Height: tileSize}

	initGame()
}

func draw() {
	rl.BeginDrawing()

	// Set the background to a nebula
	bgSource := rl.Rectangle{X: 0, Y: 0, Width: float32(texBackground.Width), Height: float32(texBackground.Height)}
	bgDest := rl.Rectangle{X: 0, Y: 0, Width: screenWidth, Height: screenHeight}
	rl.DrawTexturePro(texBackground, bgSource, bgDest, rl.Vector2{X: 0, Y: 0}, 0, rl.White)

	//Draw the player
	player.Draw()

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
	rl.UnloadTexture(texTiles)
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
