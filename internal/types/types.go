package types

type PieceType int

type Piece struct {
	PieceType      PieceType
	CooldownFrames int
}

type Position struct {
	X int
	Y int
}

type Enemy struct {
	Id               uint64
	Type             string
	Health           int
	XPos             float64
	YPos             float64
	CurrentTileIndex int
}

type VisualProjectile struct {
	Id            uint64
	X             float64
	Y             float64
	TargetEnemyId uint64
	TargetLastX   float64
	TargetLastY   float64
	Speed         float64
	Damage        int
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
