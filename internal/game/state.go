package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/louis-bourgault/chesstd/internal/board"
	"github.com/louis-bourgault/chesstd/internal/enemy"
)

type Game struct {
	Board       *board.Board
	MoveBudget  int
	CurrentWave enemy.EnemyWave
}

func NewGame() *Game {
	g := &Game{
		Board:      board.NewBoard(),
		MoveBudget: 10,
	}
	g.CurrentWave = enemy.NewWave(g.Board)
	g.CurrentWave.Begin()
	return g
}

func (g *Game) Update() error {
	// Game logic goes here (e.g., handling input, updating game state)
	g.CurrentWave.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Board.Draw(screen)
	g.CurrentWave.Draw(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	//480x480 pixel arena for the 8x8 chessboard
	return 680, 480
}
