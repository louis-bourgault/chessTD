package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/louis-bourgault/chesstd/internal/game"
)

func main() {
	//without this, looks weird on a high DPI retina display
	m := ebiten.Monitor()
	scale := m.DeviceScaleFactor() * 0.75 //this is good for my mbp, but you may want to adjust this for your display, espeically without a hidpi display
	ebiten.SetWindowSize(int(scale*680), int(480*scale))
	ebiten.SetWindowTitle("Chess Tower Defense")

	g := game.NewGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
