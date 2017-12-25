package game

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/marcospedreiro/sshtron/player"
	"github.com/marcospedreiro/sshtron/session"
	"golang.org/x/crypto/ssh"
)

const (
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

// NewGameManager returns a GameManager for when setting up a new server
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
			case KeyCtrlC, KeyEscape:
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
