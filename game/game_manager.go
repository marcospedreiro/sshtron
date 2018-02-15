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
		g = NewGame(config.GameWidth, config.GameHeight)
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
			case config.KeyW, config.KeyZ, config.KeyK, config.KeyComma:
				session.HandleUp()
			case config.KeyA, config.KeyQ, config.KeyH:
				session.HandleLeft()
			case config.KeyS, config.KeyJ, config.KeyO:
				session.HandleDown()
			case config.KeyD, config.KeyL, config.KeyE:
				session.HandleRight()
			case rune(config.KeyAccelerate):
				session.HandleSpeedUp()
			case rune(config.KeyDecelerate):
				session.HandleSlowDown()
			case rune(config.KeyCtrlC), rune(config.KeyEscape):
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
		config.GameWidth = *cfg.Game.Manager.GameWidth
	}
	if cfg.Game.Manager.GameHeight != nil {
		config.GameHeight = *cfg.Game.Manager.GameHeight
	}
	if cfg.Game.Manager.KeyW != nil {
		config.KeyW, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyW)
	}
	if cfg.Game.Manager.KeyA != nil {
		config.KeyA, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyA)
	}
	if cfg.Game.Manager.KeyS != nil {
		config.KeyS, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyS)
	}
	if cfg.Game.Manager.KeyD != nil {
		config.KeyD, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyD)
	}
	if cfg.Game.Manager.KeyZ != nil {
		config.KeyZ, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyZ)
	}
	if cfg.Game.Manager.KeyQ != nil {
		config.KeyQ, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyQ)
	}
	if cfg.Game.Manager.KeyJ != nil {
		config.KeyJ, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyJ)
	}
	if cfg.Game.Manager.KeyK != nil {
		config.KeyK, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyK)
	}
	if cfg.Game.Manager.KeyL != nil {
		config.KeyL, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyL)
	}
	if cfg.Game.Manager.KeyComma != nil {
		config.KeyComma, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyComma)
	}
	if cfg.Game.Manager.KeyO != nil {
		config.KeyO, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyO)
	}
	if cfg.Game.Manager.KeyE != nil {
		config.KeyE, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyE)
	}
	if cfg.Game.Manager.KeyE != nil {
		config.KeyE, _ = utf8.DecodeRuneInString(*cfg.Game.Manager.KeyE)
	}
	if cfg.Game.Manager.KeyCtrlC != nil {
		config.KeyCtrlC = *cfg.Game.Manager.KeyCtrlC
	}
	if cfg.Game.Manager.KeyEscape != nil {
		config.KeyEscape = *cfg.Game.Manager.KeyEscape
	}
	return
}
