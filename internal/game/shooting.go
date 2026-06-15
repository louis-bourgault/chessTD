package game

import (
	"image/color"
	"log"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/louis-bourgault/chesstd/internal/cfg"
	"github.com/louis-bourgault/chesstd/internal/types"
)

func getPieceDamage(p types.PieceType) int {
	switch p {
	case types.Pawn:
		return 2
	case types.King:
		return 3
	case types.Bishop:
		return 4
	case types.Rook:
		return 5
	case types.Knight:
		return 6
	case types.Queen:
		return 6
	}
	return 1
}

func getPieceCooldown(p types.PieceType) int {
	switch p {
	case types.Pawn:
		return 90
	case types.King:
		return 60
	case types.Bishop:
		return 45
	case types.Rook:
		return 60
	case types.Knight:
		return 90
	case types.Queen:
		return 20
	}
	return 60
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (g *Game) isPathBlocked(px, py, ex, ey int) bool {
	dx := ex - px
	dy := ey - py

	stepX := 0
	if dx > 0 {
		stepX = 1
	} else if dx < 0 {
		stepX = -1
	}

	stepY := 0
	if dy > 0 {
		stepY = 1
	} else if dy < 0 {
		stepY = -1
	}

	x := px + stepX
	y := py + stepY

	for x != ex || y != ey {
		if g.Board.Grid[y][x] != nil && g.Board.Grid[y][x].PieceType != types.None {
			return true
		}
		x += stepX
		y += stepY
	}
	return false
}

func (g *Game) isTargetable(px, py, ex, ey int, pieceType types.PieceType) bool {
	dx := abs(ex - px)
	dy := abs(ey - py)

	switch pieceType {
	case types.Rook:
		if dx != 0 && dy != 0 {
			return false
		}
		return !g.isPathBlocked(px, py, ex, ey)

	case types.Bishop:
		if dx != dy {
			return false
		}
		return !g.isPathBlocked(px, py, ex, ey)

	case types.Queen:
		if dx != 0 && dy != 0 && dx != dy {
			return false
		}
		return !g.isPathBlocked(px, py, ex, ey)

	case types.Knight:
		return (dx == 1 && dy == 2) || (dx == 2 && dy == 1)

	case types.Pawn:
		// White pawns move up (decreasing Y) and attack diagonally
		return (ey == py-1) && (dx == 1)

	case types.King:
		return dx <= 1 && dy <= 1 && (dx > 0 || dy > 0)
	}
	return false
}

func (g *Game) findTargetForPiece(px, py int, pieceType types.PieceType) *types.Enemy {
	var bestTarget *types.Enemy
	for i := range g.CurrentWave.Enemies {
		enemy := &g.CurrentWave.Enemies[i]
		ex := int(enemy.XPos)
		ey := int(enemy.YPos)

		if g.isTargetable(px, py, ex, ey, pieceType) {
			if bestTarget == nil || enemy.CurrentTileIndex > bestTarget.CurrentTileIndex {
				bestTarget = enemy
			}
		}
	}
	return bestTarget
}

func (g *Game) spawnProjectile(sx, sy float64, target *types.Enemy, damage int) {
	g.NextProjectileId++
	p := types.VisualProjectile{
		Id:            g.NextProjectileId,
		X:             sx,
		Y:             sy,
		TargetEnemyId: target.Id,
		TargetLastX:   target.XPos,
		TargetLastY:   target.YPos,
		Speed:         10.0, //tiles per second
		Damage:        damage,
	}
	g.Projectiles = append(g.Projectiles, p)
}

func (g *Game) UpdateVisualProjectiles() {
	for i := len(g.Projectiles) - 1; i >= 0; i-- {
		proj := &g.Projectiles[i]

		// Find target enemy
		var targetEnemy *types.Enemy
		for j := range g.CurrentWave.Enemies {
			if g.CurrentWave.Enemies[j].Id == proj.TargetEnemyId {
				targetEnemy = &g.CurrentWave.Enemies[j]
				break
			}
		}

		var tx, ty float64
		if targetEnemy != nil {
			tx = targetEnemy.XPos
			ty = targetEnemy.YPos
			proj.TargetLastX = tx
			proj.TargetLastY = ty
		} else {
			tx = proj.TargetLastX
			ty = proj.TargetLastY
		}

		dx := tx - proj.X
		dy := ty - proj.Y
		dist := math.Sqrt(dx*dx + dy*dy)

		moveDist := proj.Speed * cfg.Dt
		if dist <= moveDist {
			// Arrived!
			if targetEnemy != nil {
				targetEnemy.Health -= proj.Damage
				log.Printf("en %d hit, now %d\n", targetEnemy.Id, targetEnemy.Health)
				if targetEnemy.Health <= 0 {
					for idx := range g.CurrentWave.Enemies {
						if g.CurrentWave.Enemies[idx].Id == proj.TargetEnemyId {
							g.CurrentWave.Enemies = slices.Delete(g.CurrentWave.Enemies, idx, idx+1)
							log.Printf("enemy %d dead!\n", proj.TargetEnemyId)
							break
						}
					}
				}
			}
			g.Projectiles = slices.Delete(g.Projectiles, i, i+1)
		} else {
			proj.X += (dx / dist) * moveDist
			proj.Y += (dy / dist) * moveDist
		}
	}
}

func (g *Game) UpdateTowerAttacks() {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			piece := g.Board.Grid[y][x]
			if piece == nil || piece.PieceType == types.None {
				continue
			}
			if piece.CooldownFrames > 0 {
				piece.CooldownFrames--
			}
		}
	}

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			piece := g.Board.Grid[y][x]
			if piece == nil || piece.PieceType == types.None {
				continue
			}

			if piece.CooldownFrames == 0 {
				target := g.findTargetForPiece(x, y, piece.PieceType)
				if target != nil {
					damage := getPieceDamage(piece.PieceType)
					g.spawnProjectile(float64(x)+0.5, float64(y)+0.5, target, damage)
					piece.CooldownFrames = getPieceCooldown(piece.PieceType)
				}
			}
		}
	}
}

func (g *Game) DrawProjectiles(screen *ebiten.Image) {
	for _, proj := range g.Projectiles {
		vector.FillCircle(
			screen,
			float32(proj.X*cfg.TileSize+cfg.LeftMargin),
			float32(proj.Y*cfg.TileSize+cfg.TopMargin),
			float32(6),
			color.RGBA{255, 215, 0, 255}, // Gold
			false,
		)
	}
}
