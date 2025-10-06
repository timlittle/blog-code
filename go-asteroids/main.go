package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth      = 800
	screenHeight     = 400
	tileSize         = 64
	rotationSpeed    = 2.0
	playerSpeed      = 6.0
	initialAsteroids = 5
)

var (
	texTiles      rl.Texture2D
	texBackground rl.Texture2D
	spriteRec     rl.Rectangle
	boostRec      rl.Rectangle
	asteroidRec   rl.Rectangle
	asteroids     []Asteroid
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
	if p.isBoosting {
		rl.DrawTexturePro(
			texTiles,
			boostRec,
			destTexture,
			rl.Vector2{X: p.size.X / 2, Y: p.size.Y/2 - 40},
			p.rotation,
			rl.White,
		)
	}
	rl.DrawTexturePro(
		texTiles,
		spriteRec,
		destTexture,
		rl.Vector2{X: p.size.X / 2, Y: p.size.Y / 2},
		p.rotation,
		rl.White,
	)
}

func (p *Player) Update() {

	// Rotate the player with the arrow keys
	if rl.IsKeyDown(rl.KeyLeft) {
		player.rotation -= rotationSpeed
	}
	if rl.IsKeyDown(rl.KeyRight) {
		player.rotation += rotationSpeed
	}
	// Default to not boosting
	player.isBoosting = false

	// Accelerate the player with up
	if rl.IsKeyDown(rl.KeyUp) {
		if player.acceleration < 0.9 {
			player.acceleration += 0.1
		}
		player.isBoosting = true
	}
	// Decellerate the player with down
	if rl.IsKeyDown(rl.KeyDown) {
		if player.acceleration > 0 {
			player.acceleration -= 0.05
		}
		if player.acceleration < 0 {
			player.acceleration = 0

		}
	}

	// Get the direction the sprite is pointing
	direction := getDirectionVector(player.rotation)

	// Start to move to the direction
	player.speed = rl.Vector2Scale(direction, playerSpeed)

	// Accelerate in that direction
	player.position.X += player.speed.X * player.acceleration
	player.position.Y -= player.speed.Y * player.acceleration

	// To void losing our ship, we wrap around the screen
	wrapPosition(&p.position, tileSize)
}

// Enum for storing the size of the asteroid
type AsteroidSize int

const (
	Large AsteroidSize = iota
	Medium
	Small
)

type Asteroid struct {
	position     rl.Vector2
	speed        rl.Vector2
	size         rl.Vector2
	asteroidSize AsteroidSize
}

func (a *Asteroid) Draw() {
	// Draw the asteroid to the screen
	destTexture := rl.Rectangle{X: a.position.X, Y: a.position.Y, Width: a.size.X, Height: a.size.Y}
	rl.DrawTexturePro(
		texTiles,
		asteroidRec,
		destTexture,
		rl.Vector2{X: a.size.X / 2, Y: a.size.Y / 2},
		0.0,
		rl.White,
	)
}

func (a *Asteroid) Update() {
	// Move the asteroid in its direction
	a.position = rl.Vector2Add(a.position, a.speed)

	// Wrap the position, so they are always on screen
	wrapPosition(&a.position, a.size.X)
}

// Asteroid helper functions
func createLargeAsteroid() Asteroid {

	// Generate a random edge of the screen to spawn
	randomEdge := rl.GetRandomValue(0, 3)
	var position rl.Vector2

	// Generate a random position on screen
	randomX := float32(rl.GetRandomValue(0, screenWidth))
	randomY := float32(rl.GetRandomValue(0, screenHeight))

	switch randomEdge {
	case 0:
		position = rl.Vector2{X: randomX, Y: +tileSize}
	case 1:
		position = rl.Vector2{X: screenWidth + tileSize, Y: randomY}
	case 2:
		position = rl.Vector2{X: randomX, Y: screenHeight + tileSize}
	case 3:
		position = rl.Vector2{X: -tileSize, Y: randomY}
	}

	// Generate a random speed and direction for the asteroid
	speed := rl.Vector2{
		X: float32(rl.GetRandomValue(-10, 10)) / 10,
		Y: float32(rl.GetRandomValue(-10, 10)) / 10,
	}

	// Create the large asteroid
	return createAsteroid(Large, position, speed)
}

func createAsteroid(asteroidSize AsteroidSize, position, speed rl.Vector2) Asteroid {

	// Scale the image of the asteroid based on the asteroidSize
	var size rl.Vector2
	switch asteroidSize {
	case Large:
		size = rl.Vector2{X: tileSize * 1.0, Y: tileSize * 1.0}
	case Medium:
		size = rl.Vector2{X: tileSize * 0.7, Y: tileSize * 0.7}
	case Small:
		size = rl.Vector2{X: tileSize * 0.4, Y: tileSize * 0.4}
	}

	// Create the asteroid
	return Asteroid{
		position:     position,
		speed:        speed,
		size:         size,
		asteroidSize: asteroidSize,
	}
}

// Helper functions
func getDirectionVector(rotation float32) rl.Vector2 {
	// Convert the rotation to radians
	radians := float64(rotation) * rl.Deg2rad

	// Return the vector of the direction we are pointing at
	return rl.Vector2{
		X: float32(math.Sin(radians)),
		Y: float32(math.Cos(radians)),
	}
}

func wrapPosition(pos *rl.Vector2, objectSize float32) {
	// If we go off the left side of the screen
	if pos.X > screenWidth+objectSize {
		pos.X = -objectSize
	}
	// If we go off the right side of the screen
	if pos.X < -objectSize {
		pos.X = screenWidth + objectSize
	}
	// If we go off the bottom of the screen
	if pos.Y > screenHeight+objectSize {
		pos.Y = -objectSize
	}
	// If we go off the top of the screen
	if pos.Y < -objectSize {
		pos.Y = screenHeight + objectSize
	}
}

func initGame() {
	// Create the asteroids field
	asteroids = nil
	for range initialAsteroids {
		asteroids = append(asteroids, createLargeAsteroid())
	}

	// Create the player
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

	// Sprites for the ship and it boost
	spriteRec = rl.Rectangle{X: tileSize * 0, Y: tileSize * 2, Width: tileSize, Height: tileSize}
	boostRec = rl.Rectangle{X: tileSize * 7, Y: tileSize * 5, Width: tileSize, Height: tileSize}

	// Sprite for the asteroid
	asteroidRec = rl.Rectangle{X: tileSize * 1, Y: tileSize * 4, Width: tileSize, Height: tileSize}

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

	// Draw the asteroid field
	for i := range asteroids {
		asteroids[i].Draw()
	}

	// Draw the score to the screen
	rl.DrawText("Score 0", 10, 10, 20, rl.Gray)

	rl.EndDrawing()

}

func update() {
	// Update the player
	player.Update()

	// Update the asteroid field
	for i := range asteroids {
		asteroids[i].Update()
	}

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
