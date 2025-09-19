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

			// Check if there is water above this cell
			hasWaterAbove := y > 0 && g.State[y-1][x].volume > 0

			g.State[y][x].Draw(x, y, g.tileSize, hasWaterAbove)
		}
	}
}

func (g *Game) Update() {
	// Create a new state to avoid modifying the current state while reading it
	newState := CreateGameState(len(g.State[0]), len(g.State), g.tileSize)

	// Copy current state to new state
	for y := range g.State {
		copy(newState[y], g.State[y])
	}

	// Process the simulation from the bottom upwards
	for y := len(g.State) - 1; y >= 0; y-- {
		for x := range g.State[y] {

			// Only process cells that contain water
			if g.State[y][x].volume > 0 {

				// Check if we are at the bottom boundary
				if y+1 < len(g.State) {
					processWaterCell(x, y, &newState)
				}
			}
		}
	}

	// Replace old state with new calculated state
	g.State = newState

}

func canFlowDown(x, y int, state *[][]Droplet) bool {
	return y+1 < len(*state) && (*state)[y+1][x].volume < 1.0
}

func tryHorizontalFlow(x, y int, state *[][]Droplet) {
	current := &(*state)[y][x]

	// Only cascade if there's water below
	hasWaterBelow := y+1 < len(*state) && (*state)[y+1][x].volume > 0.5
	if !hasWaterBelow {
		return
	}

	// Cascade right - distribute to multiple cells
	for offset := 1; offset <= 3 && x+offset < len((*state)[y]); offset++ {
		target := &(*state)[y][x+offset]
		if target.volume < current.volume {
			flowRate := (current.volume - target.volume) * 0.1 / float64(offset)
			fill(current, target, 1.0, flowRate)
		}
	}

	// Cascade left - distribute to multiple cells
	for offset := 1; offset <= 3 && x-offset >= 0; offset++ {
		target := &(*state)[y][x-offset]
		if target.volume < current.volume {
			flowRate := (current.volume - target.volume) * 0.1 / float64(offset)
			fill(current, target, 1.0, flowRate)
		}
	}
}

func tryDiagonalFlow(x, y int, state *[][]Droplet) {
	current := &(*state)[y][x]

	// Flow diagonally down-right if space is available
	if x+1 < len((*state)[y]) && y+1 < len(*state) && (*state)[y+1][x+1].volume < 1.0 {
		fill(current, &(*state)[y+1][x+1], 1.0, 0.25)
	}

	// Flow diagonally down-left if space is available
	if x > 0 && y+1 < len(*state) && (*state)[y+1][x-1].volume < 1.0 {
		fill(current, &(*state)[y+1][x-1], 1.0, 0.25)
	}
}

func processWaterCell(x, y int, newState *[][]Droplet) {
	// Try to flow downwards, as if by gravity
	fill(&(*newState)[y][x], &(*newState)[y+1][x], 1.0, 0.5)

	// If all water flowed down, no need to continue
	if (*newState)[y][x].volume == 0 {
		return
	}

	// If water can still flow down, don't try other directions yet
	if canFlowDown(x, y, newState) {
		return
	}

	// Water spreads sideways when blocked below
	tryHorizontalFlow(x, y, newState)

	if (*newState)[y][x].volume > 0 {
		tryDiagonalFlow(x, y, newState)
	}

}

// Calculate how much more water a droplet can hold
func remainder(droplet Droplet, maxVolume float64) float64 {
	return maxVolume - droplet.volume
}

// Fill transfers water between two droplets at a controlled rate
func fill(current, target *Droplet, maxVolume, flowRate float64) {

	// Calculate how much water can be transferred
	transfer := remainder(*target, maxVolume)

	// Limit transfer to the flow rate (prevents instant teleportation)
	if transfer > flowRate {
		transfer = flowRate
	}

	// Move water from source to target
	current.volume -= transfer
	target.volume += transfer
}

type Droplet struct {
	volume float64 // How much water this cell contains (0.0 to 1.0)
	size   int
}

func (d *Droplet) Draw(x, y, tileSize int, hasWaterAbove bool) {
	// Convert grid coordinates to pixel coordinates
	pixelX := x * tileSize
	pixelY := y * tileSize

	if d.volume > 0 {
		// Calculate visual height based on water volume
		// Full volume (1.0) = full tile height, half volume (0.5) = half height
		height := int(float64(tileSize) * d.volume)

		// Fill up from the bottom
		offsetY := tileSize - height

		// If water above, fill from the top instead
		if hasWaterAbove {
			offsetY = 0
		}

		// Draw the blue water rectangle
		rl.DrawRectangle(int32(pixelX), int32(pixelY+offsetY), int32(tileSize), int32(height), rl.Blue)
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

func CreateWaterGenerator(x, y, tileSize int, state *[][]Droplet) {
	droplet := Droplet{size: tileSize, volume: 1.0}
	(*state)[y][x] = droplet
}

func main() {
	// Setup the new game
	var game = NewGame(800, 400, 10)

	// Initialize Raylib graphics window
	rl.InitWindow(int32(game.Width), int32(game.Height), "Water simulation")
	defer rl.CloseWindow()

	// Set up a counter, so we can spwan new water at a rate
	frameCount := 0
	flowStartX := 100 / game.tileSize
	flowStartY := 100 / game.tileSize
	CreateWaterGenerator(flowStartX, flowStartY, game.tileSize, &game.State)

	// Setup the frame per second rate
	rl.SetTargetFPS(20)

	// Main loop
	for !rl.WindowShouldClose() {
		frameCount++

		// Begind to draw and set the background to black
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// Add new water every 5 frames (creates continuous water stream)
		if frameCount%5 == 0 {
			CreateWaterGenerator(flowStartX, flowStartY, game.tileSize, &game.State)
			CreateWaterGenerator(flowStartX+1, flowStartY, game.tileSize, &game.State)
			CreateWaterGenerator(flowStartX-1, flowStartY, game.tileSize, &game.State)
		}

		// Draw the game state
		game.Draw()

		// Update the game state based on the rules
		game.Update()

		rl.EndDrawing()
	}

}
