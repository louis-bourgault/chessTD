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
	Shop             bool
}

func NewGame() *Game {
	g := &Game{
		Board:  board.NewBoard(),
		Health: 20,
		Round:  1,
	}
	g.CurrentWave = enemy.NewWave(g.Board.Grid, g.OnEnemyEscape)
	g.Board.IsTileOnEnemyPath = func(pos types.Position) bool {
		for _, pathPos := range g.CurrentWave.Path {
			if pathPos == pos {
				return true
			}
		}
		return false
	}
	// g.CurrentWave.Spawn()

	return g
}

func (g *Game) OnEnemyEscape(escaped types.Enemy) {
	g.Health -= 1
	log.Printf("enemy escaped, health is now %d\n", g.Health)

}

func (g *Game) Update() error {
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
	if g.CurrentWave.Finished {
		g.Shop = true
	}
	return nil

}

func (g *Game) StartNextRound() {
	g.Round += 1
	g.CurrentWave = enemy.NewWave(g.Board.Grid, g.OnEnemyEscape)
	g.Board.IsTileOnEnemyPath = func(pos types.Position) bool {
		for _, pathPos := range g.CurrentWave.Path {
			if pathPos == pos {
				return true
			}
		}
		return false
	}
	g.Health += 5
	g.Health = min(g.Health, 20) //cap health at 20
	g.Board.Reset()              //resets the board to its initial state for the new round
	g.Board.MoveBudget += 5
	g.Running = false
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
