package board

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/louis-bourgault/chesstd/internal/cfg"
	"github.com/louis-bourgault/chesstd/internal/font"
	"github.com/louis-bourgault/chesstd/internal/types"
)

type Board struct {
	//grid is where they currently are, starting is where they're reset to at the round's end, pieces are the info about them
	Grid                  [8][8]*types.Piece //x then y
	Pieces                []types.Piece      //contains all special things.
	StartingGridLocations [8][8]*types.Piece
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

	//add more later.

	// for i := 0; i < 8; i++ {
	// 	b.Grid[6][i] = types.Piece{PieceType: types.Pawn}
	// }

	return b
}

func (b *Board) Draw(screen *ebiten.Image) {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			var tileColor color.Color = color.RGBA{230, 217, 181, 255} // Light
			if (x+y)%2 != 0 {
				tileColor = color.RGBA{139, 115, 85, 255} // Dark
			}
			vector.FillRect(screen, float32(x*cfg.TileSize+cfg.LeftMargin), float32(y*cfg.TileSize+cfg.TopMargin), cfg.TileSize, cfg.TileSize, tileColor, false)

			face := &text.GoTextFace{
				Source: font.FontSource,
				Size:   32,
			}

			if b.Grid[y][x] != nil && b.Grid[y][x].PieceType != types.None {
				var pieceColor color.Color = color.RGBA{50, 50, 200, 255}
				op := &text.DrawOptions{}
				op.GeoM.Translate(float64(x*cfg.TileSize+cfg.LeftMargin), float64(y*cfg.TileSize+cfg.TopMargin))
				op.ColorScale.ScaleWithColor(pieceColor) // Target color modification

				//string function defined in types
				text.Draw(screen, b.Grid[y][x].String()[:2], face, op)
			}
		}
	}
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
