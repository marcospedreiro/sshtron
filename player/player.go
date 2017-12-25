package player

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
	"github.com/marcospedreiro/sshtron/position"
)

// PlayerDirection defines the possible directions for the player as an int
type PlayerDirection int

const (
	VerticalPlayerSpeed        = 0.007
	HorizontalPlayerSpeed      = 0.01
	PlayerCountScoreMultiplier = 1.25
	PlayerTimeout              = 15 * time.Second

	PlayerUpRune    = '⇡'
	PlayerLeftRune  = '⇠'
	PlayerDownRune  = '⇣'
	PlayerRightRune = '⇢'

	PlayerTrailHorizontal      = '┄'
	PlayerTrailVertical        = '┆'
	PlayerTrailLeftCornerUp    = '╭'
	PlayerTrailLeftCornerDown  = '╰'
	PlayerTrailRightCornerDown = '╯'
	PlayerTrailRightCornerUp   = '╮'

	PlayerRed     = color.FgRed
	PlayerGreen   = color.FgGreen
	PlayerYellow  = color.FgYellow
	PlayerBlue    = color.FgBlue
	PlayerMagenta = color.FgMagenta
	PlayerCyan    = color.FgCyan
	PlayerWhite   = color.FgWhite

	PlayerUp PlayerDirection = iota
	PlayerLeft
	PlayerDown
	PlayerRight
)

var PlayerColors = []color.Attribute{
	PlayerRed, PlayerGreen, PlayerYellow, PlayerBlue,
	PlayerMagenta, PlayerCyan, PlayerWhite,
}

var PlayerBorderColors = map[color.Attribute]color.Attribute{
	PlayerRed:     color.FgHiRed,
	PlayerGreen:   color.FgHiGreen,
	PlayerYellow:  color.FgHiYellow,
	PlayerBlue:    color.FgHiBlue,
	PlayerMagenta: color.FgHiMagenta,
	PlayerCyan:    color.FgHiCyan,
	PlayerWhite:   color.FgHiWhite,
}

var PlayerColorNames = map[color.Attribute]string{
	PlayerRed:     "Red",
	PlayerGreen:   "Green",
	PlayerYellow:  "Yellow",
	PlayerBlue:    "Blue",
	PlayerMagenta: "Magenta",
	PlayerCyan:    "Cyan",
	PlayerWhite:   "White",
}

// PlayerTrailSegment is the type for the trail left by the player
type PlayerTrailSegment struct {
	Marker rune
	Pos    position.Position
}

// Player is the struct for player specific data within the session
type Player struct {
	Name      string
	CreatedAt time.Time
	Direction PlayerDirection
	Marker    rune
	Color     color.Attribute
	Pos       *position.Position

	Trail []PlayerTrailSegment

	score float64
}

// NewPlayer creates a new player. If color is below 1, a random color is chosen
func NewPlayer(worldWidth int, worldHeight int, color color.Attribute) *Player {

	rand.Seed(time.Now().UnixNano())

	startX := rand.Float64() * float64(worldWidth)
	startY := rand.Float64() * float64(worldHeight)

	if color < 0 {
		color = PlayerColors[rand.Intn(len(PlayerColors))]
	}

	return &Player{
		CreatedAt: time.Now(),
		Marker:    PlayerDownRune,
		Direction: PlayerDown,
		Color:     color,
		Pos:       &position.Position{X: startX, Y: startY},
	}
}

func (p *Player) addTrailSegment(pos position.Position, marker rune) {
	segment := PlayerTrailSegment{marker, pos}
	p.Trail = append([]PlayerTrailSegment{segment}, p.Trail...)
}

func (p *Player) calculateScore(delta float64, playerCount int) float64 {
	rawIncrement := (delta * (float64(playerCount-1) * PlayerCountScoreMultiplier))

	// Convert millisecond increment to seconds
	actualIncrement := rawIncrement / 1000

	return p.score + actualIncrement
}

// IsOutOfBounds returns true if the current positions.Round[X|Y] are out of bounds
func (p *Player) IsOutOfBounds(minX int, maxX int, minY int, maxY int) bool {
	oob := p.Pos.RoundX() < minX || p.Pos.RoundX() >= maxX || p.Pos.RoundY() < minY || p.Pos.RoundY() >= maxY
	if oob {
		return true
	}
	return false
}

// Score of the player at the current tick of the game clock
func (p *Player) Score() int {
	return int(p.score)
}

// Update a player each tick of the game clock
func (p *Player) Update(numPlayers int, delta float64) {
	startX, startY := p.Pos.RoundX(), p.Pos.RoundY()

	switch p.Direction {
	case PlayerUp:
		p.Pos.Y -= VerticalPlayerSpeed * delta
	case PlayerLeft:
		p.Pos.X -= HorizontalPlayerSpeed * delta
	case PlayerDown:
		p.Pos.Y += VerticalPlayerSpeed * delta
	case PlayerRight:
		p.Pos.X += HorizontalPlayerSpeed * delta
	}

	endX, endY := p.Pos.RoundX(), p.Pos.RoundY()

	// If we moved, add a trail segment.
	if endX != startX || endY != startY {
		var lastSeg *PlayerTrailSegment
		var lastSegX, lastSegY int
		if len(p.Trail) > 0 {
			lastSeg = &p.Trail[0]
			lastSegX = lastSeg.Pos.RoundX()
			lastSegY = lastSeg.Pos.RoundY()
		}

		pos := position.PositionFromInt(startX, startY)

		switch {
		// Handle corners. This took an ungodly amount of time to figure out. Highly
		// recommend you don't touch.
		case lastSeg != nil &&
			(p.Direction == PlayerRight && endX > lastSegX && endY < lastSegY) ||
			(p.Direction == PlayerDown && endX < lastSegX && endY > lastSegY):
			p.addTrailSegment(pos, PlayerTrailLeftCornerUp)
		case lastSeg != nil &&
			(p.Direction == PlayerUp && endX > lastSegX && endY < lastSegY) ||
			(p.Direction == PlayerLeft && endX < lastSegX && endY > lastSegY):
			p.addTrailSegment(pos, PlayerTrailRightCornerDown)
		case lastSeg != nil &&
			(p.Direction == PlayerDown && endX > lastSegX && endY > lastSegY) ||
			(p.Direction == PlayerLeft && endX < lastSegX && endY < lastSegY):
			p.addTrailSegment(pos, PlayerTrailRightCornerUp)
		case lastSeg != nil &&
			(p.Direction == PlayerRight && endX > lastSegX && endY > lastSegY) ||
			(p.Direction == PlayerUp && endX < lastSegX && endY < lastSegY):
			p.addTrailSegment(pos, PlayerTrailLeftCornerDown)

		// Vertical and horizontal trails
		case endX == startX && endY < startY:
			p.addTrailSegment(pos, PlayerTrailVertical)
		case endX < startX && endY == startY:
			p.addTrailSegment(pos, PlayerTrailHorizontal)
		case endX == startX && endY > startY:
			p.addTrailSegment(pos, PlayerTrailVertical)
		case endX > startX && endY == startY:
			p.addTrailSegment(pos, PlayerTrailHorizontal)
		}
	}

	p.score = p.calculateScore(delta, numPlayers)
}

type ByColor []*Player

func (slice ByColor) Len() int {
	return len(slice)
}

func (slice ByColor) Less(i, j int) bool {
	return PlayerColorNames[slice[i].Color] < PlayerColorNames[slice[j].Color]
}

func (slice ByColor) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
