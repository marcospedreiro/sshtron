package player

import (
	"math/rand"
	"time"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/marcospedreiro/sshtron/config"
	"github.com/marcospedreiro/sshtron/position"
)

// PlayerDirection defines the possible directions for the player as an int
type PlayerDirection int

const (
	PlayerRed     = color.FgRed
	PlayerGreen   = color.FgGreen
	PlayerYellow  = color.FgYellow
	PlayerBlue    = color.FgBlue
	PlayerMagenta = color.FgMagenta
	PlayerCyan    = color.FgCyan
	PlayerWhite   = color.FgWhite

	PlayerUp    PlayerDirection = 0
	PlayerLeft  PlayerDirection = 1
	PlayerDown  PlayerDirection = 2
	PlayerRight PlayerDirection = 3
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

	Trail    []PlayerTrailSegment
	TrailLen int

	BaseSpeed     float64
	SpeedModifier float64
	score         float64
}

// NewPlayer creates a new player. If color is below 1, a random color is chosen
func NewPlayer(worldWidth int, worldHeight int, color color.Attribute) *Player {

	rand.Seed(time.Now().UnixNano())

	startX := rand.Float64() * float64(worldWidth)
	startY := rand.Float64() * float64(worldHeight)

	if color < 0 {
		color = PlayerColors[rand.Intn(len(PlayerColors))]
	}

	spawnDirection, spawnMarker, spawnSpeed := GetSpawnDirection(worldWidth, worldHeight, startX, startY)
	return &Player{
		CreatedAt:     time.Now(),
		Marker:        spawnMarker,
		Direction:     spawnDirection,
		Color:         color,
		Pos:           &position.Position{X: startX, Y: startY},
		TrailLen:      0,
		BaseSpeed:     spawnSpeed,
		SpeedModifier: 0,
	}
}

// GetSpawnDirection determines which quadrant a players start position is in and returns a direction
// that is suitable for the position
func GetSpawnDirection(worldWidth int, worldHeight int, StartX float64, StartY float64) (PlayerDirection, rune, float64) {
	switch {
	case StartX < float64(worldWidth/2) && StartY < float64(worldHeight/2):
		if rand.Float64() > 0.5 {
			return PlayerDown, config.PlayerDownRune, config.VerticalPlayerSpeed
		}
		return PlayerRight, config.PlayerRightRune, config.HorizontalPlayerSpeed
	case StartX > float64(worldWidth/2) && StartY < float64(worldHeight/2):
		if rand.Float64() > 0.5 {
			return PlayerDown, config.PlayerDownRune, config.VerticalPlayerSpeed
		}
		return PlayerLeft, config.PlayerLeftRune, config.HorizontalPlayerSpeed
	case StartX < float64(worldWidth/2) && StartY > float64(worldHeight/2):
		if rand.Float64() > 0.5 {
			return PlayerUp, config.PlayerUpRune, config.VerticalPlayerSpeed
		}
		return PlayerRight, config.PlayerRightRune, config.HorizontalPlayerSpeed
	case StartX > float64(worldWidth/2) && StartY > float64(worldHeight/2):
		if rand.Float64() > 0.5 {
			return PlayerUp, config.PlayerUpRune, config.VerticalPlayerSpeed
		}
		return PlayerLeft, config.PlayerLeftRune, config.HorizontalPlayerSpeed
	default:
		return PlayerDown, config.PlayerDownRune, config.VerticalPlayerSpeed
	}
}

func (p *Player) addTrailSegment(pos position.Position, marker rune) {
	segment := PlayerTrailSegment{marker, pos}
	p.Trail = append([]PlayerTrailSegment{segment}, p.Trail...)
	p.TrailLen++
	return
}

func (p *Player) trimTrailSegments(delta float64) {
	trailSegmentsToKeep := config.PlayerTrailMaxLength
	if config.LimitPlayerTrailByTime {
		trailSegmentsToKeep = int(p.CurrentSpeed()*delta*float64(config.GameLoopsPerSecond)) * config.PlayerTrailMaxTime
	}
	if p.TrailLen > trailSegmentsToKeep {
		p.Trail = p.Trail[:trailSegmentsToKeep]
		p.TrailLen = trailSegmentsToKeep
	}
	return
}

// CurrentSpeed for the player: BaseSpeed + SpeedModifier
func (p *Player) CurrentSpeed() float64 {
	return p.BaseSpeed + p.SpeedModifier
}

func (p *Player) calculateScore(delta float64, playerCount int) float64 {
	rawIncrement := (delta * (float64(playerCount-1) * config.PlayerCountScoreMultiplier))

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
		p.Pos.Y -= p.CurrentSpeed() * delta
	case PlayerLeft:
		p.Pos.X -= p.CurrentSpeed() * delta
	case PlayerDown:
		p.Pos.Y += p.CurrentSpeed() * delta
	case PlayerRight:
		p.Pos.X += p.CurrentSpeed() * delta
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
			p.addTrailSegment(pos, config.PlayerTrailLeftCornerUp)
		case lastSeg != nil &&
			(p.Direction == PlayerUp && endX > lastSegX && endY < lastSegY) ||
			(p.Direction == PlayerLeft && endX < lastSegX && endY > lastSegY):
			p.addTrailSegment(pos, config.PlayerTrailRightCornerDown)
		case lastSeg != nil &&
			(p.Direction == PlayerDown && endX > lastSegX && endY > lastSegY) ||
			(p.Direction == PlayerLeft && endX < lastSegX && endY < lastSegY):
			p.addTrailSegment(pos, config.PlayerTrailRightCornerUp)
		case lastSeg != nil &&
			(p.Direction == PlayerRight && endX > lastSegX && endY > lastSegY) ||
			(p.Direction == PlayerUp && endX < lastSegX && endY < lastSegY):
			p.addTrailSegment(pos, config.PlayerTrailLeftCornerDown)

		// Vertical and horizontal trails
		case endX == startX && endY < startY:
			p.addTrailSegment(pos, config.PlayerTrailVertical)
		case endX < startX && endY == startY:
			p.addTrailSegment(pos, config.PlayerTrailHorizontal)
		case endX == startX && endY > startY:
			p.addTrailSegment(pos, config.PlayerTrailVertical)
		case endX > startX && endY == startY:
			p.addTrailSegment(pos, config.PlayerTrailHorizontal)
		}

		if config.PlayerTrailLengthLimit {
			p.trimTrailSegments(delta)
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

// SetPlayerProperties reads cfg.Game.Player.* and overrides the default player
// properties with values in the configuration json if set
// TODO: There must be a better way to do this?
func SetPlayerProperties(cfg *config.Config) {
	if cfg.Game.Player.VerticalSpeed != nil {
		config.VerticalPlayerSpeed = *cfg.Game.Player.VerticalSpeed
	}
	if cfg.Game.Player.HorizontalSpeed != nil {
		config.HorizontalPlayerSpeed = *cfg.Game.Player.HorizontalSpeed
	}
	if cfg.Game.Player.MaxSpeedIncrease != nil {
		config.MaxSpeedIncrease = *cfg.Game.Player.MaxSpeedIncrease
	}
	if cfg.Game.Player.MaxSpeedDecrease != nil {
		config.MaxSpeedDecrease = *cfg.Game.Player.MaxSpeedDecrease
	}
	if cfg.Game.Player.AccelerationFactor != nil {
		config.AccelerationFactor = *cfg.Game.Player.AccelerationFactor
	}
	if cfg.Game.Player.DecelerationFactor != nil {
		config.DecelerationFactor = *cfg.Game.Player.DecelerationFactor
	}
	if cfg.Game.Player.CountScoreMultiplier != nil {
		config.PlayerCountScoreMultiplier = *cfg.Game.Player.CountScoreMultiplier
	}
	if cfg.Game.Player.TimeoutSeconds != nil {
		config.PlayerTimeout = time.Duration(*cfg.Game.Player.TimeoutSeconds) * time.Second
	}
	if cfg.Game.Player.UpRune != nil {
		config.PlayerUpRune, _ = utf8.DecodeRuneInString(*cfg.Game.Player.UpRune)
	}
	if cfg.Game.Player.DownRune != nil {
		config.PlayerDownRune, _ = utf8.DecodeRuneInString(*cfg.Game.Player.DownRune)
	}
	if cfg.Game.Player.LeftRune != nil {
		config.PlayerLeftRune, _ = utf8.DecodeRuneInString(*cfg.Game.Player.LeftRune)
	}
	if cfg.Game.Player.RightRune != nil {
		config.PlayerRightRune, _ = utf8.DecodeRuneInString(*cfg.Game.Player.RightRune)
	}
	if cfg.Game.Player.TrailHorizontal != nil {
		config.PlayerTrailHorizontal, _ = utf8.DecodeRuneInString(*cfg.Game.Player.TrailHorizontal)
	}
	if cfg.Game.Player.TrailVertical != nil {
		config.PlayerTrailVertical, _ = utf8.DecodeRuneInString(*cfg.Game.Player.TrailVertical)
	}
	if cfg.Game.Player.TrailLeftCornerUp != nil {
		config.PlayerTrailLeftCornerUp, _ = utf8.DecodeRuneInString(*cfg.Game.Player.TrailLeftCornerUp)
	}
	if cfg.Game.Player.TrailLeftCornerDown != nil {
		config.PlayerTrailLeftCornerDown, _ = utf8.DecodeRuneInString(*cfg.Game.Player.TrailLeftCornerDown)
	}
	if cfg.Game.Player.TrailRightCornerDown != nil {
		config.PlayerTrailRightCornerDown, _ = utf8.DecodeRuneInString(*cfg.Game.Player.TrailRightCornerDown)
	}
	if cfg.Game.Player.TrailRightCornerUp != nil {
		config.PlayerTrailRightCornerUp, _ = utf8.DecodeRuneInString(*cfg.Game.Player.TrailRightCornerUp)
	}
	if cfg.Game.Player.PlayerTrailLengthLimit != nil {
		config.PlayerTrailLengthLimit = *cfg.Game.Player.PlayerTrailLengthLimit
	}
	if cfg.Game.Player.LimitPlayerTrailByTime != nil {
		config.LimitPlayerTrailByTime = *cfg.Game.Player.LimitPlayerTrailByTime
	}
	if cfg.Game.Player.PlayerTrailMaxLength != nil {
		config.PlayerTrailMaxLength = *cfg.Game.Player.PlayerTrailMaxLength
	}
	if cfg.Game.Player.PlayerTrailMaxTime != nil {
		config.PlayerTrailMaxTime = *cfg.Game.Player.PlayerTrailMaxTime
	}
	return
}
