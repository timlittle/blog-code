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

func (g *Game) Update() {

	// New game state
	newState := CreateGameState(len(g.State[0]), len(g.State))

	// Loop through each row
	for indexY, cellY := range g.State {

		// Loop through each column
		for indexX, cellX := range cellY {

			// Count how many neighbours the current cell has
			neighbours := CountNeighbours(indexX, indexY, g.State)

			// Update the new state using the rule based on neighbours
			newState[indexY][indexX] = IsCellAlive(cellX, neighbours)
		}
	}

	// Set the new state to the updated state
	g.State = newState
}

func CountNeighbours(x, y int, gameState [][]int) int {
	// Counter for the neighbours
	count := 0

	// Loop through all the rows
	for cellX := x - 1; cellX <= x+1; cellX++ {

		// Loop through all the columes
		for cellY := y - 1; cellY <= y+1; cellY++ {

			// We want to make sure we do not count past the boundry of the board
			if cellY < 0 || cellX < 0 || cellY >= len(gameState) || cellX >= len(gameState[0]) {
				continue
			}
			// If current cell, we can skip it
			if cellY == y && cellX == x {
				continue
			}
			//  Check if cell is alive
			if gameState[cellY][cellX] == 1 {
				count++
			}

		}
	}
	return count
}

func IsCellAlive(current, neighbours int) int {
	switch {
	// Any live cell with fewer than two live neighbours dies
	// as if by underpopulation.
	case neighbours < 2:
		return 0
	// Any live cell with two or three live neighbours lives
	// on to the next generation.
	// Any dead cell with two neighbours, remains dead
	case neighbours == 2:
		return current
	// Any dead cell with exactly three live neighbours becomes a
	// live cell, as if by reproduction.
	case neighbours == 3:
		return 1
	// Any live cell with more than three live neighbours dies
	// as if by overpopulation.
	case neighbours > 3:
		return 0
	}
	return 0
}

func CreateGameState(newWidth, newHeight int) [][]int {
	// Create a new game state with the right height
	newState := make([][]int, newHeight)

	// Create the rows with the right length
	for i := range newHeight {
		newState[i] = make([]int, newWidth)
	}

	// Return the new state map
	return newState
}

func CreateGliders(x, y int, gameState *[][]int) {
	// Draw the glider patter in the game state
	(*gameState)[y][x+1] = 1
	(*gameState)[y+1][x+2] = 1
	(*gameState)[y+2][x] = 1
	(*gameState)[y+2][x+1] = 1
	(*gameState)[y+2][x+2] = 1
}

func addPattern(x, y int, pattern [][]int, gameState *[][]int) {
	// Loop through the row
	for row := range pattern {

		// Loop through the
		for col := 0; col < len(pattern[row]); col++ {

			// Update the game state if cell alive
			if pattern[row][col] == 1 {
				(*gameState)[y+row][x+col] = 1
			}
		}
	}
}

func CreateGliderGun(x, y int, gameState *[][]int) {
	// Create a slice of the pattern
	pattern := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1},
		{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	addPattern(x, y, pattern, gameState)

}
func CreatePulsar(x, y int, gameState *[][]int) {
	// Create a slice of the pattern
	pattern := [][]int{
		{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0},
	}

	addPattern(x, y, pattern, gameState)
}

func CreatePentadecathlon(x, y int, gameState *[][]int) {
	// Create a slice of the pattern
	pattern := [][]int{
		{0, 0, 1, 0, 0, 0, 0, 1, 0, 0},
		{1, 1, 0, 1, 1, 1, 1, 0, 1, 1},
		{0, 0, 1, 0, 0, 0, 0, 1, 0, 0},
	}

	addPattern(x, y, pattern, gameState)
}

func NewGame(width, height, tileSize int) *Game {
	g := &Game{Width: width, Height: height, tileSize: tileSize}
	g.State = CreateGameState(g.Width/g.tileSize, g.Height/g.tileSize)
	return g
}

func main() {
	// Create the game metadata and state holding object
	var game = NewGame(800, 400, 10)
	CreateGliderGun(0, 0, &game.State)
	CreatePentadecathlon(40, 10, &game.State)
	CreatePulsar(60, 20, &game.State)

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

		// Update the game state before drawing
		game.Update()

		// Draw the game state
		game.Draw()

		// End the drawing
		rl.EndDrawing()
	}
}
