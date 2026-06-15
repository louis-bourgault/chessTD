package game

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	text1 "github.com/louis-bourgault/chesstd/internal/textRendering"
)

const startButtonStartX = 0
const startButtonStartY = 0
const startButtonWidth = 100
const startButtonHeight = 50

func (g *Game) RenderUI(screen *ebiten.Image) {
	vector.FillRect(screen, startButtonStartX, startButtonStartY, startButtonWidth, startButtonHeight, color.White, false)
	face := &text.GoTextFace{
		Source: text1.FontSource,
		Size:   32,
	}

	var textColor color.Color = color.RGBA{50, 50, 200, 255}
	op := &text.DrawOptions{}
	op.GeoM.Translate(startButtonStartX, startButtonStartY)
	op.ColorScale.ScaleWithColor(textColor) // Target color modification
	if !g.Running {

		text.Draw(screen, "Start", face, op)
	} else {
		text.Draw(screen, "Pause", face, op)
	}
}

func (g *Game) UpdateUI() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		log.Println("mouse click detected")
		posX, posY := ebiten.CursorPosition()
		if posX >= startButtonStartX && posX <= startButtonStartX+startButtonWidth && posY >= startButtonStartY && posY <= startButtonStartY+startButtonHeight {
			if !g.Running {
				log.Println("starting game")
				g.Running = true
			} else {
				log.Println("pausing game")
				g.Running = false
			}
		}

	}
}
