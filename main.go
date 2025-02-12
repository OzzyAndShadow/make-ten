package main

import (
	"bytes"
	_ "embed"
	"image/png"
	"math"
	"math/rand"
	"strconv"

	sprites "github.com/OzzyAndShadow/make-ten/libs/sprite"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenSize = 600
	boardSize  = 10
	maxTime    = 90
	targetFPS  = 60
)

//go:embed assets/spritesheet.png
var rawSpriteSheet []byte

func main() {
	rl.InitWindow(screenSize, screenSize+screenSize/10, "Make Ten")
	defer rl.CloseWindow()

	dragStart := rl.NewVector2(-1, -1)
	palette := getPalette()

	spriteSheet := loadEmbeddableTexture(rawSpriteSheet)
	defer rl.UnloadTexture(spriteSheet)

	var grid [boardSize][boardSize]int
	refreshGrid(&grid)
	var score int
	time := maxTime

	var numberSprites [9]sprites.Sprite
	for i := 0; i < 9; i++ {
		numberSprites[i] = sprites.Sprite{
			Frames: []int{i},
			FPS:    1,
		}
		numberSprites[i].SetSpriteSheet(spriteSheet, 8)
	}

	var bottomSprites [10]sprites.Sprite
	barFrames := []int{9, 10, 11, 12, 9, 11, 12, 9, 10, 11}
	for i := 0; i < 10; i++ {
		bottomSprites[i] = sprites.Sprite{
			Frames: []int{barFrames[i]},
			FPS:    1,
		}
		bottomSprites[i].SetSpriteSheet(spriteSheet, 8)
	}

	rl.SetTargetFPS(targetFPS)
	sprites.SetTargetFPS(targetFPS)

	frames := 0
	gameOver := false

	for !rl.WindowShouldClose() {
		frames++
		if frames == targetFPS {
			frames = 0
			time--
			if time < 0 {
				time = 0
				gameOver = true
				dragStart.X = -1
			}
		}

		resetButtonCollision := rl.NewRectangle(0, screenSize, screenSize/10*3, screenSize/10)
		if rl.CheckCollisionPointRec(rl.GetMousePosition(), resetButtonCollision) {
			if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
				refreshGrid(&grid)
				score = 0
				frames = 0
				time = maxTime
				gameOver = false
			}
		}

		if dragStart.X == -1 && rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			pos := rl.GetMousePosition()
			dragStart.X = float32(math.Floor(float64(pos.X / screenSize * boardSize)))
			dragStart.Y = float32(math.Floor(float64(pos.Y / screenSize * boardSize)))
		} else if dragStart.X != -1 && !rl.IsMouseButtonDown(rl.MouseButtonLeft) {
			selection := getMouseSelection(dragStart)

			sum := getSumOfSelection(selection, &grid)
			if sum == 10 {
				score += int(selection.size.X) * int(selection.size.Y)
				clearSelection(selection, &grid)
			}

			dragStart.X = -1
			dragStart.Y = -1
		}

		rl.BeginDrawing()

		rl.ClearBackground(palette.Background)

		if dragStart.X != -1 && dragStart.Y != -1 && !gameOver {
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

		bottomCellSize := screenSize / 10
		for i := 0; i < 10; i++ {
			bottomSprites[i].Display(
				float32(i*bottomCellSize),
				float32(boardSize*cellSize),
				float32(bottomCellSize),
				float32(bottomCellSize),
			)
		}

		resetX := int32(bottomCellSize/2*3) - rl.MeasureText("RESET", 24)/2
		rl.DrawText("RESET", resetX, int32(boardSize*cellSize+bottomCellSize/3), 24, palette.Text)

		scoreX := int32(screenSize/2) - rl.MeasureText(strconv.Itoa(score), 24)/2
		rl.DrawText(strconv.Itoa(score), scoreX, int32(boardSize*cellSize+bottomCellSize/3), 24, palette.Text)

		timeX := int32(screenSize-bottomCellSize*3/2) - rl.MeasureText(strconv.Itoa(time), 24)/2
		rl.DrawText(strconv.Itoa(time), timeX, int32(boardSize*cellSize+bottomCellSize/3), 24, palette.Text)

		if gameOver {
			clone := palette.Background
			clone.A = 200
			rl.DrawRectangle(0, 0, screenSize, screenSize, clone)
		}

		rl.EndDrawing()
	}
}

func loadEmbeddableTexture(raw []byte) rl.Texture2D {
	img, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		panic(err)
	}
	return rl.LoadTextureFromImage(rl.NewImageFromImage(img))
}

type selection struct {
	position rl.Vector2
	size     rl.Vector2
}

func getSumOfSelection(sel selection, grid *[boardSize][boardSize]int) int {
	sum := 0
	for x := sel.position.X; x < sel.position.X+sel.size.X; x++ {
		for y := sel.position.Y; y < sel.position.Y+sel.size.Y; y++ {
			sum += grid[int(x)][int(y)]
		}
	}
	return sum
}

func clearSelection(sel selection, grid *[boardSize][boardSize]int) {
	for x := sel.position.X; x < sel.position.X+sel.size.X; x++ {
		for y := sel.position.Y; y < sel.position.Y+sel.size.Y; y++ {
			grid[int(x)][int(y)] = 0
		}
	}
}

type palette struct {
	Background rl.Color
	Selection  rl.Color
	Text       rl.Color
}

func getPalette() palette {
	var p palette

	p.Background = rl.NewColor(30, 30, 46, 255)
	p.Selection = rl.NewColor(49, 50, 68, 255)
	p.Text = rl.NewColor(205, 214, 244, 255)

	return p
}

func getMouseSelection(dragStart rl.Vector2) selection {
	var sel selection
	sel.position = dragStart
	mousePosition := rl.GetMousePosition()
	cellX := float64(mousePosition.X) / screenSize * boardSize
	cellY := float64(mousePosition.Y) / screenSize * boardSize
	if cellX < float64(dragStart.X) {
		cellX = math.Floor(cellX)
	} else {
		cellX = math.Ceil(cellX)
	}
	if cellY < float64(dragStart.Y) {
		cellY = math.Floor(cellY)
	} else {
		cellY = math.Ceil(cellY)
	}
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
