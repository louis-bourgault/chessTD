package enemy

import (
	"image/color"
	"math/rand"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/louis-bourgault/chesstd/internal/cfg"
	"github.com/louis-bourgault/chesstd/internal/types"
)

const BasicEnemyMoveSpeed = 5.0 //tiles per second
const FramesPerSpawn = 60

func GetPieceAt(board [8][8]*types.Piece, pos types.Position) bool {
	if pos.X < 0 || pos.X >= 8 || pos.Y < 0 || pos.Y >= 8 {
		return false
	}
	return board[pos.Y][pos.X] != nil
}

func NewWave(board [8][8]*types.Piece, escFunc func(escaped types.Enemy)) EnemyWave {
	wave := EnemyWave{}
	//generate a path where there aren't pieces.
	//always starts somewhere on the left edge

	//we might want something more complex later, for example trying to avoid the range of pieces, or trying to stay away from the player's pieces, but for now anything goes
	//then we add obstacles that are made to protect the path
	//then enemies after that

	//could also take a param for difficulty which effects numbers of enemies, obstacles, and how much they try to avoid pieces, etc.

	foundPath := false

	for !foundPath {
		foundStart := false
		candidatePath := []types.Position{}

		//first find a start pos, and then go from there. if we hit a dead end, start over with a new start pos. if we get to the right edge, we're done.
		for !foundStart {
			startY := rand.Intn(8)
			if !GetPieceAt(board, types.Position{X: 0, Y: startY}) {
				candidatePath = append(candidatePath, types.Position{X: 0, Y: startY})
				foundStart = true
			}
		}

		for {
			currentPos := candidatePath[len(candidatePath)-1]

			if currentPos.X == 7 {
				foundPath = true
				break
			}

			nextSteps := []types.Position{
				{X: currentPos.X + 1, Y: currentPos.Y},
				{X: currentPos.X, Y: currentPos.Y - 1},
				{X: currentPos.X, Y: currentPos.Y + 1},
			}

			validNextSteps := []types.Position{}
			for _, step := range nextSteps {
				if step.X >= 0 && step.X < 8 && step.Y >= 0 && step.Y < 8 && !GetPieceAt(board, step) {
					validNextSteps = append(validNextSteps, step)
				}
				//if we've already been there, don't go back.
				for _, pos := range candidatePath {
					if pos == step {
						for i, s := range validNextSteps {
							if s == step {
								validNextSteps = append(validNextSteps[:i], validNextSteps[i+1:]...)
								break
							}
						}
					}
				}
			}
			if len(validNextSteps) == 0 {
				break
			}
			//randomyl choose.
			nextStep := validNextSteps[rand.Intn(len(validNextSteps))]
			candidatePath = append(candidatePath, nextStep)
		}

		if foundPath {
			wave.Path = candidatePath
		}

	}
	wave.OnEnemyEscaped = escFunc
	wave.SpawnCounter = 0
	wave.ToSpawn = 10
	return wave

}

func (w *EnemyWave) Spawn() {
	//starts spawning enemies

	w.Enemies = append(w.Enemies, types.Enemy{
		Type:             "Basic", //NOT based on chess pieces, we want more of a bloons td style where they're differentiated by health, speed, resistance, etc.
		Health:           10,
		XPos:             float64(w.Path[0].X) + 0.5, //the float represents how far through the tile they are, not their absolute pos.
		YPos:             float64(w.Path[0].Y) + 0.5,
		CurrentTileIndex: 0,
		Id:               w.SpawnCounter,
	})
	w.SpawnCounter++
}

type EnemyWave struct {
	Path                 []types.Position
	Enemies              []types.Enemy
	Obstacles            []types.Position
	framesSinceLastSpawn int
	OnEnemyEscaped       func(escaped types.Enemy)
	SpawnCounter         uint64
	Finished             bool
	ToSpawn              int
}

func (w *EnemyWave) Update() {
	//if we go backwards, we can delete them better
	if w.framesSinceLastSpawn >= FramesPerSpawn && w.SpawnCounter < uint64(w.ToSpawn) {
		w.framesSinceLastSpawn = 0
		w.Spawn()
	}
	w.framesSinceLastSpawn++
	for i := len(w.Enemies) - 1; i >= 0; i-- {
		enemy := &w.Enemies[i]

		if enemy.CurrentTileIndex < len(w.Path) {
			currentTile := w.Path[enemy.CurrentTileIndex]
			var nextTile types.Position

			if enemy.CurrentTileIndex+1 >= len(w.Path) {
				nextTile = types.Position{
					Y: currentTile.Y,
					X: currentTile.X + 1,
				}
			} else {
				nextTile = w.Path[enemy.CurrentTileIndex+1]
			}

			axisX := nextTile.X - currentTile.X
			axisY := nextTile.Y - currentTile.Y
			enemy.XPos += float64(axisX) * cfg.Dt * BasicEnemyMoveSpeed
			enemy.YPos += float64(axisY) * cfg.Dt * BasicEnemyMoveSpeed
			targetX := float64(nextTile.X) + 0.5
			targetY := float64(nextTile.Y) + 0.5

			if (axisX > 0 && enemy.XPos >= targetX) || (axisX < 0 && enemy.XPos <= targetX) ||
				(axisY > 0 && enemy.YPos >= targetY) || (axisY < 0 && enemy.YPos <= targetY) {

				enemy.XPos = targetX
				enemy.YPos = targetY
				enemy.CurrentTileIndex++
			}
		} else {
			//because we are iterating backwards, qwe can use slices.delete safely
			copyEnemy := *enemy
			w.Enemies = slices.Delete(w.Enemies, i, i+1)
			w.OnEnemyEscaped(copyEnemy)
		}
	}
	if len(w.Enemies) == 0 && w.SpawnCounter >= uint64(w.ToSpawn) {
		w.Finished = true
	}

}

func (w *EnemyWave) Draw(screen *ebiten.Image) {
	//draw enemies and obstacles on the screen
	for _, tile := range w.Path {
		//draw path for testing
		vector.FillRect(screen, float32(tile.X*cfg.TileSize+cfg.LeftMargin), float32(tile.Y*cfg.TileSize+cfg.TopMargin), float32(cfg.TileSize), float32(cfg.TileSize), color.RGBA{255, 0, 0, 128}, false)
	}
	for _, enemy := range w.Enemies {
		vector.FillCircle(screen, float32(enemy.XPos*cfg.TileSize+cfg.LeftMargin), float32(enemy.YPos*cfg.TileSize+cfg.TopMargin), float32(20), color.RGBA{227, 149, 189, 128}, false)
	}
	vector.FillCircle(screen, float32(1*cfg.TileSize+cfg.LeftMargin-0.5*cfg.TileSize), float32(1*cfg.TileSize+cfg.TopMargin-0.5*cfg.TileSize), float32(20), color.RGBA{227, 149, 189, 128}, false)

}
