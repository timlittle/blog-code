# Building Conway’s Game of Life in Go with raylib-go (Step by Step)
## A step-by-step tutorial on building Conway’s Game of Life in Go, using raylib-go for graphics. Learn how to draw grids, apply the rules of life, and simulate evolving patterns.

I recently started to play around with graphics programming and game engine creation. For that, I was using SDL2, OpenGL and C.

However, I wanted to do something in Go, so I started with [Conway’s Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life). There are various options from C bindings to [OpenGL](https://github.com/go-gl/gl) and [SDL2](https://github.com/veandco/go-sdl2) to full game engines such as [Ebitengine](https://ebitengine.org/). I went with [raylib-go](https://github.com/gen2brain/raylib-go) as it was a good mixture of low-level and abstracted programming.

In this tutorial, we’ll use **Go** and **raylib-go** to create Conway’s Game of Life, a simulation where simple rules create endless patterns. By the end, you’ll have gliders drifting, pulsars pulsing, and even a glider gun firing across your screen.

### The Game of Life

Conway’s Game of Life is a cellular automaton created by John Conway. It is a zero-player game which evolves. Each generation is built on the previous generation using a set of simple rules.

The game requires no players as you start it off, and it starts to evolve without intervention. It is [Turing Complete](https://en.wikipedia.org/wiki/Turing_completeness) and can simulate a wide variety of different patterns and Turing Machines.

It starts with an initial game start, then it uses the following rules to determine which cell lives or dies:

1. Any live cell with **fewer than 2** live neighbours dies, as if by underpopulation.
2. Any **live cell with 2 or 3** live neighbours lives on to the next generation.
3. Any live cell with **more than 3** live neighbours dies, as if by overpopulation.
4. Any dead cell with **exactly 3 live** neighbours becomes a live cell, as if by reproduction.

Each generation or frame in the context of our game will use these rules to determine the next state.

### Demo

In this tutorial, we will be creating a simple cellular automation in Go using RayLib-go. It will start with a state we define, and allow us to modify the state to see different evolutions.

![](static/game-of-life/conways.gif)


### Setting up the project

First thing we need to do is create a go package for our project and pull down raylib-go so we can start to build with it.

```go
mkdir go-gameoflife  
go mod init gameoflife  
go get -v -u github.com/gen2brain/raylib-go/raylib
```

When learning anything new in programming, we always start with the classic [Hello World](https://en.wikipedia.org/wiki/%22Hello,_World!%22_program). To do that, we need to import raylib-go, create a window and display some text.

```go
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
  
 // Close the window at the end of the program  
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
```

We have also created a struct to store the game metadata, and we can later use it to store the state of the game. We use that to initialise the window.

To save a little bit of time, we can create a [Makefile](https://en.wikipedia.org/wiki/Make_%28software%29), which will allow us to simplify the building and running of our code. In this case, we only need one target **run,** however, this can be expanded to run tests, build more files, clean up build states and more.

```make
.DEFAULT_GOAL := run  
  
run:  
 go run .  
  
help:  
 @echo "run - run the game"  
  
.PHONY: run
```

Now that we have the make file, we can run make to start the program. The same could be done by running. `go run .`

```
make
```

![](static/game-of-life/hello.webp)

Hello world running in raylib-go

Win! We now have a window with our Hello World.

In the background, raylib is doing all the hard work for us, working with OpenGL to create the window in the system-specific libraries for your operating system. However, we do not need to worry about that.

For each step of this tutorial, I will provide the code so you can compare. The full code for this step is in the [GitHub Repo](https://github.com/timlittle/blog-code/blob/87e511b14d464d023c870ad3dd77f83b6083350c/go-gameoflife/main.go).

#### Creating a 2D map

The next thing we need to do is to create an initial state for our game, which will consist of a 2D matrix which we use to draw the cells to the screen.

We will need to update our Game struct to store our state and to specify how big we want the cells.

```diff
 type Game struct {  
-       Height int  
-       Width  int  
+       Height   int  
+       Width    int  
+       tileSize int  
+       State    [][]int  
+}
```

We can then create a function to create the new start for us

```go
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
```

Next, we need a way to display this on our window. We can simply loop through the array and check if the cell is alive. If it is, we can draw the cell to the screen. We do this with a `Draw()` method on the **Game** struct.

```go
func (g *Game) Draw() {  
 // Loop through all of the rows  
 for y := range g.State {  
  
  // Loop through all of the columns  
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
```

Now we need to update our `main()`function to use the new game state, then draw it to the window.


```diff
-       game := Game{Width: 800, Height: 400}  
+       var game = NewGame(800, 400, 80)  
  
...  
  
-               // Draw Hello world  
-               rl.DrawText("Hello world!", 350, 200, 20, rl.RayWhite)  
+               // Draw the game state  
+               game.Draw()

```

Now we can run the game again.

```
make
```

![](static/game-of-life/cells.webp)

We now have the starting state for our game. I have used a cell size of 80x80px here to allow us to render it in an 800x400px window, so we can play with the logic and see the results.

Full code for this step is in the [GitHub Repo](https://github.com/timlittle/blog-code/blob/379b93addb38c711377209fef1c7ee7a0569901e/go-gameoflife/main.go)

#### Getting neighbours

Now we get to the part where we bring our game to life. We can start to create the rule to evolve the state.

To do this, we need a way of counting how many live neighbours a cell has. We can do this by looping through all of the neighbours and checking the surrounding cells.

![](static/game-of-life/count_neighbours_diagram.svg)

Logic for checking neighbours

We can start with the top left cell, the move to the next until we have covered all the cells. We will need to exclude the current cell and account for whether the cell we are calculating is at the edge of the window.

```go
func CountNeighbours(x, y int, gameState [][]int) int {  
 // Counter for the neighbours  
 count := 0  
  
 // Loop through all the rows  
 for cellX := x - 1; cellX <= x+1; cellX++ {  
  
  // Loop through all the columns  
  for cellY := y - 1; cellY <= y+1; cellY++ {  
  
   // We want to make sure we do not count past the boundary of the board  
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
```

Next, we need to codify the rules of the game using the count of the neighbours. We can do that with a `switch` statement for each of the rules.

```go
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
```

Now that we have the calculations for our next state implemented, we need to take our existing state and update it based on the rules. We can create an `Update()` method on our Game struct to do this.

```go
func (g *Game) Update() {  
  
 // We can stat with hardcoded state for now  
 // However, we would want to update thisbased on the current state  
 var newState = [][]int{  
  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  
  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  
  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  
  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  
  {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  
 }  
  
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
```

We have used another slicee to build the new state, then we replace the old one. This implementation increases memory usage, however, we could offset and update it in-place with a truth table and two passes of the slice. However, let's keep it simple.

```diff
rl.ClearBackground(rl.Black)  
  
+               // Update the game state before drawing  
+               game.Update()  
+  
                // Draw the game state

```
Now we can add the Update function to our main game loop and start the game again to see if it works.

```
make

```

Run it and you will see the cells start to evolve!

We can see the patterns start to evolve with each generation, and they will converge into a steady state eventually.

![](static/game-of-life/updates.gif)

Conways game of life with updates

Full code for this step is in the [GitHub Repo](https://github.com/timlittle/blog-code/blob/1d8b8944f90046cc7e61d7c13068f493b25f2d21/go-gameoflife/main.go)

#### Update the map

Now that the core functionality of our game is working, we can start to scale the map size to make more complex patterns.

Let's replace the hard-coded map with a function we can use to generate a game state as big as we want.

```go
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
```

We can use this blank game state in our code where we are manually defining the state. We are defining the state in both the `Update()` and `NewGame()` functions.

```diff
        // However, we would want to update this based on the current state  
-       var newState = [][]int{  
-               {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  
-               {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  
-               {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  
-               {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  
-               {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  
-       }  
+       newState := CreateGameState(len(g.State[0]), len(g.State))  
  
        // Loop through each row
```

```diff
 func NewGame(width, height, tileSize int) *Game {  
        g := &Game{Width: width, Height: height, tileSize: tileSize}  
-       g.State = [][]int{  
-               {0, 1, 0, 0, 0, 0, 0, 0, 0, 0},  
-               {0, 0, 1, 0, 0, 0, 0, 1, 0, 0},  
-               {1, 1, 1, 0, 0, 0, 0, 1, 0, 0},  
-               {0, 0, 0, 0, 0, 0, 0, 1, 0, 0},  
-               {0, 0, 0, 0, 0, 0, 0, 0, 0, 0},  
-       }  
+       g.State = CreateGameState(g.Width/g.tileSize, g.Height/g.tileSize)  
        return g  
 }
```

That looks much cleaner. However, we have now lost our initial state, which triggers the game. We can fix this by creating a function that will add a new pattern. We will start with the [Glider](https://en.wikipedia.org/wiki/Glider_%28Conway%27s_Game_of_Life%29) pattern.

```go
func CreateGliders(x, y int, gameState *[][]int) {  
 // Draw the glider patter in the game state  
 (*gameState)[y][x+1] = 1  
 (*gameState)[y+1][x+2] = 1  
 (*gameState)[y+2][x] = 1  
 (*gameState)[y+2][x+1] = 1  
 (*gameState)[y+2][x+2] = 1  
}
```

Now that we have the function to create a Glider, we can add them to our game state in the `main()` function.

```diff
-       var game = NewGame(800, 400, 80)  
+       var game = NewGame(800, 400, 10)  
+       CreateGliders(0, 0, &game.State)  
+       CreateGliders(10, 0, &game.State)  
+       CreateGliders(20, 0, &game.State)  
+       CreateGliders(30, 0, &game.State)
```

Let's run the game again.

```
make
```

Great!

We can see the addition of 4 gliders to a much larger window, and we can see them move across the screen as the generations increase.

![](static/game-of-life/gliders.gif)

Scaled game state with gliders

Full code for this step is in the [GitHub Repo](https://github.com/timlittle/blog-code/blob/4fa23147105c81cdf8a3848db0b90311f110468e/go-gameoflife/main.go)

#### Setting up different simulations

Now that we have the game working, we can play around with different initial states to see how they evolve.

The other patterns are more complex compared to the glider, so I will create a helper function for adding the pattern, then define the patterns separately. More patterns can be found [here](https://conwaylife.appspot.com/library/).

```go
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
```

```go
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
```

```go
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
```

```go
func CreatePentadecathlon(x, y int, gameState *[][]int) {  
 // Create a slice of the pattern  
 pattern := [][]int{  
  {0, 0, 1, 0, 0, 0, 0, 1, 0, 0},  
  {1, 1, 0, 1, 1, 1, 1, 0, 1, 1},  
  {0, 0, 1, 0, 0, 0, 0, 1, 0, 0},  
 }  
  
 addPattern(x, y, pattern, gameState)  
}
```

Now we have our patterns, let's add them to the state.

```diff
var game = NewGame(800, 400, 10)  
-       CreateGliders(0, 0, &game.State)  
-       CreateGliders(10, 0, &game.State)  
-       CreateGliders(20, 0, &game.State)  
-       CreateGliders(30, 0, &game.State)  
+       CreateGliderGun(0, 0, &game.State)  
+       CreatePentadecathlon(40, 10, &game.State)  
+       CreatePulsar(60, 20, &game.State)  
  
        // Create the Raylib window using the state
```

Run our game again.

```
make
```

Now our game is busy, with different interactions and patterns.

![](static/game-of-life/conways.gif)

Full game of life

The full code can be found on the [GitHub repo](https://github.com/timlittle/blog-code/tree/main/go-gameoflife)

#### Testing

Something I have omitted from this tutorial is the testing. When developing this post, I created tests alongside the code to validate that the evolution rules and update function worked as expected. You can find my testing in the [main_test.go](https://github.com/timlittle/blog-code/blob/main/go-gameoflife/main_test.go).

The example, the code to test the IsCellAlive function, looks like:

```go
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
```

In future posts, I can cover testing more and potentially walk through how we could have done this with Test Driven Development (TDD)

#### Conclusion

We have walked through:

- Creating a window with raylib-go
- Adding shapes to the window
- Creating Conway’s Game of Life state
- Updating the state with new generations
- Adding new patterns to the state

We just built Conway’s Game of Life in Go with raylib-go, starting from a blank window all the way to gliders, pulsars, and even a glider gun. Along the way, we covered rendering, updating state, and adding reusable patterns.

You can find the full source code in the [GitHub repo](#). Try running it, tweak the patterns, or create your own.

If you build something cool with it, I’d love to see it! Share your experiments in the comments, or tag me on LinkedIn/GitHub.

