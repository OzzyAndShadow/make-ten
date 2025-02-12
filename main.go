package main

import (
	"math"
	"math/rand"

	sprites "github.com/OzzyAndShadow/make-ten/libs/sprite"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenSize = 600
	boardSize  = 10
	time       = 30
)

func main() {
	rl.InitWindow(screenSize, screenSize+100, "Make Ten")
	defer rl.CloseWindow()

	dragStart := rl.NewVector2(-1, -1)
	palette := getPalette()
	spriteSheet := rl.LoadTexture("assets/spritesheet.png")
	var grid [boardSize][boardSize]int
	refreshGrid(&grid)

	var numberSprites [9]sprites.Sprite

	for i := 0; i < 9; i++ {
		numberSprites[i] = sprites.Sprite{
			Frames: []int{i},
			FPS:    1,
		}
		numberSprites[i].SetSpriteSheet(spriteSheet, 8)
	}

	rl.SetTargetFPS(60)
	sprites.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(palette.Background)

		if dragStart.X == -1 && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			pos := rl.GetMousePosition()
			dragStart.X = float32(math.Floor(float64(pos.X / screenSize * boardSize)))
			dragStart.Y = float32(math.Floor(float64(pos.Y / screenSize * boardSize)))
		} else if dragStart.X != -1 && !rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			dragStart.X = -1
			dragStart.Y = -1
		}

		if dragStart.X != -1 && dragStart.Y != -1 {
			selection := getMouseSelection(dragStart)
			cellSize := screenSize / boardSize
			rl.DrawRectangle(
				int32(selection.position.X*float32(cellSize)),
				int32(selection.position.Y*float32(cellSize)),
				int32(selection.size.X*float32(cellSize)),
				int32(selection.size.Y*float32(cellSize)),
				palette.Selection,
			)
		}

		cellSize := screenSize / boardSize
		for x := 0; x < boardSize; x++ {
			for y := 0; y < boardSize; y++ {
				cell := grid[x][y]
				if cell != 0 {
					numberSprites[cell-1].Display(
						float32(x*cellSize),
						float32(y*cellSize),
						float32(cellSize),
						float32(cellSize),
					)
				}
			}
		}

		rl.EndDrawing()
	}
}

type selection struct {
	position rl.Vector2
	size     rl.Vector2
}

type palette struct {
	Background rl.Color
	Selection  rl.Color
}

func getPalette() palette {
	var p palette

	p.Background = rl.NewColor(17, 20, 40, 255)
	p.Selection = rl.NewColor(29, 34, 68, 255)

	return p
}

func getMouseSelection(dragStart rl.Vector2) selection {
	var sel selection
	sel.position = dragStart
	mousePosition := rl.GetMousePosition()
	cellX := int(math.Round(float64(mousePosition.X) / screenSize * boardSize))
	cellY := int(math.Round(float64(mousePosition.Y) / screenSize * boardSize))
	size := rl.NewVector2(float32(cellX)-dragStart.X, float32(cellY)-dragStart.Y)

	if size.X < 0 {
		sel.position = rl.Vector2Add(sel.position, rl.NewVector2(size.X, 0))
		size.X = -size.X + 1
	}

	if size.Y < 0 {
		sel.position = rl.Vector2Add(sel.position, rl.NewVector2(0, size.Y))
		size.Y = -size.Y + 1
	}

	sel.size = size

	if sel.position.X < 0 {
		sel.size = rl.Vector2Add(sel.size, rl.NewVector2(sel.position.X, 0))
		sel.position.X = 0
	}

	if sel.position.Y < 0 {
		sel.size = rl.Vector2Add(sel.size, rl.NewVector2(0, sel.position.Y))
		sel.position.Y = 0
	}

	if sel.position.X+sel.size.X > boardSize {
		sel.size = rl.Vector2Subtract(sel.size, rl.NewVector2(sel.position.X+sel.size.X-boardSize, 0))
	}

	if sel.position.Y+sel.size.Y > boardSize {
		sel.size = rl.Vector2Subtract(sel.size, rl.NewVector2(0, sel.position.Y+sel.size.Y-boardSize))
	}

	return sel
}

func refreshGrid(grid *[boardSize][boardSize]int) {
	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			grid[i][j] = rand.Intn(9) + 1
		}
	}
}
