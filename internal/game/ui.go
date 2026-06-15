package game

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/louis-bourgault/chesstd/internal/cfg"
	text1 "github.com/louis-bourgault/chesstd/internal/textRendering"
)

const startButtonStartX = cfg.LeftMargin
const startButtonStartY = cfg.TopMargin + 8*cfg.TileSize + 20
const startButtonWidth = 100
const startButtonHeight = 50

const HealthXPos = cfg.LeftMargin + 8*cfg.TileSize + 20
const HealthYPos = cfg.TopMargin + 40

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
	healthText := "Health: " + fmt.Sprintf("%d", g.Health)
	text1.DrawText(screen, healthText, HealthXPos, HealthYPos, 32, color.RGBA{255, 0, 0, 255})
	text1.DrawText(screen, fmt.Sprintf("Round: %d", g.Round), HealthXPos, HealthYPos+40, 32, color.RGBA{255, 255, 255, 255})

	if g.Shop {
		shopText := "Shop"
		text1.DrawText(screen, shopText, HealthXPos, HealthYPos+40, 24, color.RGBA{255, 255, 0, 255})

		vector.FillRect(screen, HealthXPos, HealthYPos+80, 200, 50, color.RGBA{0, 255, 0, 255}, false)
		text1.DrawText(screen, "Next Round", HealthXPos, HealthYPos+80, 24, color.RGBA{0, 0, 0, 255})
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
		} else if g.Shop && posX >= HealthXPos && posX <= HealthXPos+200 && posY >= HealthYPos+80 && posY <= HealthYPos+80+50 {
			log.Println("starting next round")
			g.Shop = false
			g.StartNextRound()
		}

	}
}
