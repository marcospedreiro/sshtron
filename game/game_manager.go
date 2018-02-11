package game

import (
	"bufio"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/marcospedreiro/sshtron/config"
	"github.com/marcospedreiro/sshtron/player"
	"github.com/marcospedreiro/sshtron/session"
	"golang.org/x/crypto/ssh"
)

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

// GameManager maintains the list of running games
type GameManager struct {
	Games         map[string]*Game
	HandleChannel chan ssh.Channel
}

// NewGameManager returns a Manager for when setting up a new server
func NewGameManager() *GameManager {
	return &GameManager{
		Games:         map[string]*Game{},
		HandleChannel: make(chan ssh.Channel),
	}
}

/*HandleNewChannel handles a new player connection either by joining a game or creating a new one
Then handles a connected players actions up through the player leaving the game
*/
func (gm *GameManager) HandleNewChannel(sc ssh.Channel, color string) {
	g := gm.getAvailableGame()
	if g == nil {
		g = NewGame(GameWidth, GameHeight)
		gm.Games[g.Name] = g

		go g.Run()
	}

	colorOptions := g.AvailableColors()
	finalColor := colorOptions[0]

	// choose the requested color if available
	color = strings.ToLower(color)
	for _, clr := range colorOptions {
		if strings.ToLower(player.PlayerColorNames[clr]) == color {
			finalColor = clr
			break
		}
	}

	session := session.NewSession(sc, g.WorldWidth(), g.WorldHeight(), finalColor)
	g.AddSession(session)

	go func() {
		reader := bufio.NewReader(sc)
		for {
			r, _, err := reader.ReadRune()
			if err != nil {
				fmt.Println(err)
				break
			}

			switch r {
			case KeyW, KeyZ, KeyK, KeyComma:
				session.HandleUp()
			case KeyA, KeyQ, KeyH:
				session.HandleLeft()
			case KeyS, KeyJ, KeyO:
				session.HandleDown()
			case KeyD, KeyL, KeyE:
				session.HandleRight()
			case rune(KeyCtrlC), rune(KeyEscape):
				if g.SessionsCount() == 1 {
					delete(gm.Games, g.Name)
				}
				g.RemoveSession(session)
			}
		}
	}()

}

// SessionCount all the sessions known by the game manager
func (gm *GameManager) SessionCount() int {
	sum := 0
	for _, game := range gm.Games {
		sum += game.SessionsCount()
	}
	return sum
}

// GameCount all the games known by the game manager
func (gm *GameManager) GameCount() int {
	return len(gm.Games)
}

// getAvailableGame returns a reference to a game with available spots for
// players. If one does not exist, nil is returned.
func (gm *GameManager) getAvailableGame() *Game {
	var g *Game

	for _, game := range gm.Games {
		spots := game.AvailableColors()
		if len(spots) > 0 {
			g = game
			break
		}
	}

	return g
}

// SetGameManagerProperties reads cfg.Game.Manager.* and overrides the default
// game manager properties with values in the configuration json if set
// TODO: There must be a better way to do this?
func SetGameManagerProperties(cfg *config.Config) {
	if cfg.Game.Manager.GameWidth != nil {
		GameWidth = *cfg.Game.Manager.GameWidth
	}
	if cfg.Game.Manager.GameHeight != nil {
		GameHeight = *cfg.Game.Manager.GameHeight
	}
	if cfg.Game.Manager.KeyW != nil {
		KeyW, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyW)
	}
	if cfg.Game.Manager.KeyA != nil {
		KeyA, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyA)
	}
	if cfg.Game.Manager.KeyS != nil {
		KeyS, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyS)
	}
	if cfg.Game.Manager.KeyD != nil {
		KeyD, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyD)
	}
	if cfg.Game.Manager.KeyZ != nil {
		KeyZ, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyZ)
	}
	if cfg.Game.Manager.KeyQ != nil {
		KeyQ, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyQ)
	}
	if cfg.Game.Manager.KeyJ != nil {
		KeyJ, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyJ)
	}
	if cfg.Game.Manager.KeyK != nil {
		KeyK, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyK)
	}
	if cfg.Game.Manager.KeyL != nil {
		KeyL, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyL)
	}
	if cfg.Game.Manager.KeyComma != nil {
		KeyComma, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyComma)
	}
	if cfg.Game.Manager.KeyO != nil {
		KeyO, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyO)
	}
	if cfg.Game.Manager.KeyE != nil {
		KeyE, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyE)
	}
	if cfg.Game.Manager.KeyE != nil {
		KeyE, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyE)
	}
	if cfg.Game.Manager.KeyCtrlC != nil {
		KeyCtrlC = *cfg.Game.Manager.KeyCtrlC
	}
	if cfg.Game.Manager.KeyEscape != nil {
		KeyEscape = *cfg.Game.Manager.KeyEscape
	}
	return
}
