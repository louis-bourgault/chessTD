package board

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/louis-bourgault/chesstd/internal/cfg"
	"github.com/louis-bourgault/chesstd/internal/textRendering"
	text1 "github.com/louis-bourgault/chesstd/internal/textRendering"
	"github.com/louis-bourgault/chesstd/internal/types"
)

const MoveBudgetXPos = cfg.LeftMargin + 8*cfg.TileSize + 20
const MoveBudgetYPos = cfg.TopMargin

type Board struct {
	//grid is where they currently are, starting is where they're reset to at the round's end, pieces are the info about them
	Grid                  [8][8]*types.Piece //x then y
	Pieces                []types.Piece      //contains all special things.
	StartingGridLocations [8][8]*types.Piece
	MoveBudget            int
	SelectedTile          *types.Position
	IsTileOnEnemyPath     func(pos types.Position) bool
}

func NewBoard() *Board {
	b := &Board{}

	p := &types.Piece{PieceType: types.King}
	b.Pieces = append(b.Pieces, *p)
	b.Grid[7][4] = p
	b.StartingGridLocations[7][4] = p

	p = &types.Piece{PieceType: types.Queen}
	b.Pieces = append(b.Pieces, *p)
	b.Grid[7][3] = p
	b.StartingGridLocations[7][3] = p

	p = &types.Piece{PieceType: types.Bishop}
	b.Pieces = append(b.Pieces, *p)
	b.Grid[7][2] = p
	b.StartingGridLocations[7][2] = p

	p = &types.Piece{PieceType: types.Knight}
	b.Pieces = append(b.Pieces, *p)
	b.Grid[7][1] = p
	b.StartingGridLocations[7][1] = p

	p = &types.Piece{PieceType: types.Rook}
	b.Pieces = append(b.Pieces, *p)
	b.Grid[7][0] = p
	b.StartingGridLocations[7][0] = p

	p = &types.Piece{PieceType: types.Bishop}
	b.Pieces = append(b.Pieces, *p)
	b.Grid[7][5] = p
	b.StartingGridLocations[7][5] = p

	p = &types.Piece{PieceType: types.Knight}
	b.Pieces = append(b.Pieces, *p)
	b.Grid[7][6] = p
	b.StartingGridLocations[7][6] = p

	p = &types.Piece{PieceType: types.Rook}
	b.Pieces = append(b.Pieces, *p)
	b.Grid[7][7] = p
	b.StartingGridLocations[7][7] = p

	//add more later.

	for i := 0; i < 8; i++ {
		p = &types.Piece{PieceType: types.Pawn}
		b.Pieces = append(b.Pieces, *p)
		b.Grid[6][i] = p
		b.StartingGridLocations[6][i] = p
	}

	b.MoveBudget = 10

	return b
}

func (b *Board) Reset() {
	b.Grid = b.StartingGridLocations
	b.MoveBudget = 10
}

func (b *Board) Draw(screen *ebiten.Image) {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			var tileColor color.Color = color.RGBA{230, 217, 181, 255}
			if (x+y)%2 != 0 {
				tileColor = color.RGBA{139, 115, 85, 255}
			}
			if b.SelectedTile != nil && b.SelectedTile.X == x && b.SelectedTile.Y == y {
				tileColor = color.RGBA{255, 255, 0, 255}
			}
			vector.FillRect(screen, float32(x*cfg.TileSize+cfg.LeftMargin), float32(y*cfg.TileSize+cfg.TopMargin), cfg.TileSize, cfg.TileSize, tileColor, false)

			face := &text.GoTextFace{
				Source: text1.FontSource,
				Size:   32,
			}

			if b.Grid[y][x] != nil && b.Grid[y][x].PieceType != types.None {
				var pieceColor color.Color = color.RGBA{50, 50, 200, 255}
				op := &text.DrawOptions{}
				op.GeoM.Translate(float64(x*cfg.TileSize+cfg.LeftMargin), float64(y*cfg.TileSize+cfg.TopMargin))
				op.ColorScale.ScaleWithColor(pieceColor)

				//string function defined in types
				text.Draw(screen, b.Grid[y][x].String()[:2], face, op)
			}
		}
	}
	moveText := fmt.Sprintf("Move Budget: %d", b.MoveBudget)
	// log.Printf(moveText)
	textRendering.DrawText(screen, moveText, MoveBudgetXPos, MoveBudgetYPos, 32, color.RGBA{255, 255, 255, 255})

}

func (b *Board) Update() {
	if b.MoveBudget <= 0 {
		return
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		posX, posY := ebiten.CursorPosition()
		gridX := (posX - cfg.LeftMargin) / cfg.TileSize
		gridY := (posY - cfg.TopMargin) / cfg.TileSize
		if b.SelectedTile == nil {
			b.SelectedTile = &types.Position{X: gridX, Y: gridY}
		} else {
			from := *b.SelectedTile
			to := types.Position{X: gridX, Y: gridY}
			if b.IsValidMove(from, to) {
				piece := b.GetPieceAt(from)
				b.Grid[to.Y][to.X] = piece
				b.Grid[from.Y][from.X] = nil
				b.SelectedTile = nil
				b.MoveBudget -= 1
			} else {
				b.SelectedTile = nil
			}
		}
	}
}

func (b *Board) IsValidMove(from, to types.Position) bool {
	if b.IsTileOnEnemyPath(to) {
		return false
	}
	if from.X < 0 || from.X >= 8 || from.Y < 0 || from.Y >= 8 {
		return false
	}
	if to.X < 0 || to.X >= 8 || to.Y < 0 || to.Y >= 8 {
		return false
	}
	piece := b.GetPieceAt(from)
	if piece.PieceType == types.None {
		return false
	}
	//switch case here for piece type
	return true
}

func (b *Board) GetPieceAt(pos types.Position) *types.Piece {
	if pos.X >= 0 && pos.X < 8 && pos.Y >= 0 && pos.Y < 8 {
		//the array is y, x, which is slightly weird.
		p := b.Grid[pos.Y][pos.X]
		if p == nil {
			return &types.Piece{PieceType: types.None}
		}
		return p
	}
	return &types.Piece{PieceType: types.None}
	//we can't return nil, there would be null pointer exceptions all over the place
}
