package types

type PieceType int

type Piece struct {
	PieceType PieceType
}

type Position struct {
	X int
	Y int
}

type Enemy struct {
	Type             string
	Health           int
	XPos             float64
	YPos             float64
	CurrentTileIndex int
}

const (
	None = iota
	King
	Queen
	Bishop
	Knight
	Rook
	Pawn
)

func (p Piece) String() string {
	switch p.PieceType {
	case None:
		return "None"
	case King:
		return "King"
	case Queen:
		return "Queen"
	case Bishop:
		return "Bishop"
	case Knight:
		return "Knight"
	case Rook:
		return "Rook"
	case Pawn:
		return "Pawn"
	default:
		return "Unknown"
	}
}
