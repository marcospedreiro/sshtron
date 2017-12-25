package session

import (
	"time"

	"github.com/fatih/color"
	"github.com/marcospedreiro/sshtron/player"
	"golang.org/x/crypto/ssh"
)

// Session holds all data for a connected player
type Session struct {
	C ssh.Channel

	LastAction time.Time
	HighScore  int
	Player     *player.Player
}

// NewSession creates a new Session
func NewSession(c ssh.Channel, worldWidth, worldHeight int,
	color color.Attribute) *Session {

	s := Session{C: c, LastAction: time.Now()}
	s.newGame(worldWidth, worldHeight, color)

	return &s
}

func (s *Session) newGame(worldWidth, worldHeight int, color color.Attribute) {
	s.Player = player.NewPlayer(worldWidth, worldHeight, color)
}

func (s *Session) didAction() {
	s.LastAction = time.Now()
}

// StartOver resets the current player in the game
func (s *Session) StartOver(worldWidth int, worldHeight int) {
	s.newGame(worldWidth, worldHeight, s.Player.Color)
}

// Read input over the connection channel
func (s *Session) Read(p []byte) (int, error) {
	return s.C.Read(p)
}

// Write output over the connection channel
func (s *Session) Write(p []byte) (int, error) {
	return s.C.Write(p)
}

// HandleUp responds to the player pressing the up direction key
func (s *Session) HandleUp() {
	if s.Player.Direction == player.PlayerDown {
		return
	}
	s.Player.Direction = player.PlayerUp
	s.Player.Marker = player.PlayerUpRune
	s.didAction()
}

// HandleDown responds to the player pressing the down direction key
func (s *Session) HandleDown() {
	if s.Player.Direction == player.PlayerUp {
		return
	}
	s.Player.Direction = player.PlayerDown
	s.Player.Marker = player.PlayerDownRune
	s.didAction()
}

// HandleLeft responds to the player pressing the left direction key
func (s *Session) HandleLeft() {
	if s.Player.Direction == player.PlayerRight {
		return
	}
	s.Player.Direction = player.PlayerLeft
	s.Player.Marker = player.PlayerLeftRune
	s.didAction()
}

// HandleRight responds to the player pressing the right direction key
func (s *Session) HandleRight() {
	if s.Player.Direction == player.PlayerLeft {
		return
	}
	s.Player.Direction = player.PlayerRight
	s.Player.Marker = player.PlayerRightRune
	s.didAction()
}
