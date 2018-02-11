package config

// GameServerConf is the config struct for the game properties
type GameServerConf struct {
	VerticalWall   *string
	HorizontalWall *string
	TopLeft        *string
	TopRight       *string
	BottomRight    *string
	BottomLeft     *string
	Grass          *string
	Blocker        *string
}

// Characters for rendering
var (
	VerticalWall   = '║'
	HorizontalWall = '═'
	TopLeft        = '╔'
	TopRight       = '╗'
	BottomRight    = '╝'
	BottomLeft     = '╚'

	Grass   = ' '
	Blocker = '■'
)
