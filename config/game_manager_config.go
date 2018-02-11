package config

// GameManagerConf is the config struct for the game server properties
type GameManagerConf struct {
	GameWidth  *int
	GameHeight *int
	KeyW       *string
	KeyA       *string
	KeyS       *string
	KeyD       *string
	KeyZ       *string
	KeyQ       *string
	KeyH       *string
	KeyJ       *string
	KeyK       *string
	KeyL       *string
	KeyComma   *string
	KeyO       *string
	KeyE       *string
	KeyCtrlC   *int
	KeyEscape  *int
}

var (
	GameWidth  = 78
	GameHeight = 22

	KeyW = 'w'
	KeyA = 'a'
	KeyS = 's'
	KeyD = 'd'

	KeyZ = 'z'
	KeyQ = 'q'
	// KeyS and KeyD are already defined

	KeyH = 'h'
	KeyJ = 'j'
	KeyK = 'k'
	KeyL = 'l'

	KeyComma = ','
	KeyO     = 'o'
	KeyE     = 'e'

	KeyCtrlC  = 3
	KeyEscape = 27
)
