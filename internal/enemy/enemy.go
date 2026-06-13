package enemy

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/louis-bourgault/chesstd/internal/board"
	"github.com/louis-bourgault/chesstd/internal/cfg"
	"github.com/louis-bourgault/chesstd/internal/types"
)

func NewWave(board *board.Board) EnemyWave {
	wave := EnemyWave{}
	//generate a path where there aren't pieces.
	//always starts somewhere on the left edge

	//we might want something more complex later, for example trying to avoid the range of pieces, or trying to stay away from the player's pieces, but for now anything goes
	//then we add obstacles that are made to protect the path
	//then enemies after that

	foundPath := false

	for !foundPath {
		foundStart := false
		candidatePath := []types.Position{}

		//first find a start pos, and then go from there. if we hit a dead end, start over with a new start pos. if we get to the right edge, we're done.
		for !foundStart {
			startY := rand.Intn(8)
			if board.GetPieceAt(types.Position{X: 0, Y: startY}).PieceType == types.None {
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
				if step.X >= 0 && step.X < 8 && step.Y >= 0 && step.Y < 8 && board.GetPieceAt(step).PieceType == types.None {
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

	return wave
}

type EnemyWave struct {
	Path      []types.Position
	Enemies   []types.Enemy
	Obstacles []types.Position
}

func (w *EnemyWave) Update() {
	//move enemies along the path, handle interactions, etc.

}

func (w *EnemyWave) Draw(screen *ebiten.Image) {
	//draw enemies and obstacles on the screen
	for _, tile := range w.Path {
		//draw path for testing
		vector.FillRect(screen, float32(tile.X*cfg.TileSize+cfg.LeftMargin), float32(tile.Y*cfg.TileSize+cfg.TopMargin), float32(cfg.TileSize), float32(cfg.TileSize), color.RGBA{255, 0, 0, 128}, false)
	}
}
