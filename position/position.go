package position

// Position contains x and y coordinates
type Position struct {
	X float64
	Y float64
}

// PositionFromInt creates a position struct by converting x and y (int) to float64
func PositionFromInt(x int, y int) Position {
	return Position{float64(x), float64(y)}
}

// RoundX adds .5 to the x coordinate and converts float64 to int
func (p Position) RoundX() int {
	return int(p.X + 0.5)
}

// RoundX adds .5 to the y coordinate and converts float64 to int
func (p Position) RoundY() int {
	return int(p.Y + 0.5)
}
