package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/louis-bourgault/chesstd/internal/board"
	"github.com/louis-bourgault/chesstd/internal/enemy"
	"github.com/louis-bourgault/chesstd/internal/types"
)

type Game struct {
	Board            *board.Board
	CurrentWave      enemy.EnemyWave
	Health           int
	Projectiles      []types.VisualProjectile
	NextProjectileId uint64
	Running          bool
	Lost             bool
	Round            int
}

func NewGame() *Game {
	g := &Game{
		Board:  board.NewBoard(),
		Health: 20,
		Round:  1,
	}
	g.CurrentWave = enemy.NewWave(g.Board.Grid, g.OnEnemyEscape)
	// g.CurrentWave.Spawn()

	return g
}

func (g *Game) OnEnemyEscape(escaped types.Enemy) {
	g.Health -= 1
	log.Printf("enemy escaped, health is now %d\n", g.Health)

}

func (g *Game) Update() error {
	// Game logic goes here (e.g., handling input, updating game state)
	if g.Health <= 0 {
		g.Running = false
		log.Println("Game Over!")
		return nil
	}
	g.UpdateUI()
	g.Board.Update()
	if !g.Running { //we're not in the middle of a wave, so don't do enemy stuff.
		return nil
	}
	g.CurrentWave.Update()
	g.UpdateTowerAttacks()
	g.UpdateVisualProjectiles()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Board.Draw(screen)
	g.CurrentWave.Draw(screen)
	g.DrawProjectiles(screen)
	g.RenderUI(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	//480x480 pixel arena for the 8x8 chessboard
	return 680, 480
}
