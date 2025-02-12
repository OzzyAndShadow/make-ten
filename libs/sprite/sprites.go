package sprites

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var targetFPS int

func SetTargetFPS(newValue int) {
	targetFPS = newValue
}

type Sprite struct {
	Frames      []int
	FPS         int
	frame       int
	spriteSheet SpriteSheet
}

type SpriteSheet struct {
	texture    rl.Texture2D
	spriteSize int
}

// frameTime = targetFPS / sprite.FPS
// index = floor(sprite.frame / frameTime) % sprite.Order.length
// reset when frame > frameTime * sprite.Order.length

func (sprite *Sprite) AppendFrame(frame int) {
	sprite.Frames = append(sprite.Frames, frame)
}

func (sprite *Sprite) SetSpriteSheet(texture rl.Texture2D, spriteSize int) {
	sprite.spriteSheet = SpriteSheet{texture: texture, spriteSize: spriteSize}
}

func (sprite *Sprite) Update() {
	sprite.frame++
	frameTime := targetFPS / sprite.FPS
	if sprite.frame > frameTime*len(sprite.Frames) {
		sprite.frame = 0
	}
}

func (sprite *Sprite) Display(x float32, y float32, width float32, height float32) {
	frameTime := targetFPS / sprite.FPS
	index := int(sprite.frame/frameTime) % len(sprite.Frames)
	currentFrame := sprite.Frames[index]
	rl.DrawTexturePro(sprite.spriteSheet.texture, textureIndexToRectangle(sprite, currentFrame), rl.NewRectangle(x, y, width, height), rl.NewVector2(0, 0), 0, rl.White)
}

func textureIndexToRectangle(sprite *Sprite, index int) rl.Rectangle {
	x := (index * sprite.spriteSheet.spriteSize) % int(sprite.spriteSheet.texture.Width)
	y := (index / (int(sprite.spriteSheet.texture.Width) / sprite.spriteSheet.spriteSize)) * sprite.spriteSheet.spriteSize
	return rl.NewRectangle(float32(x), float32(y), float32(sprite.spriteSheet.spriteSize), float32(sprite.spriteSheet.spriteSize))
}
