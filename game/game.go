package game

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"time"
	"unicode/utf8"

	"github.com/dustinkirkland/golang-petname"
	"github.com/fatih/color"
	"github.com/marcospedreiro/sshtron/config"
	"github.com/marcospedreiro/sshtron/player"
	"github.com/marcospedreiro/sshtron/position"
	"github.com/marcospedreiro/sshtron/session"
)

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

// TileType an int used to describe what type of surface we are
type TileType int

const (
	// TileGrass is the default tiletype that players will move on
	TileGrass TileType = iota
	// TileBlocker is the default tiletype for the barriers to the level
	TileBlocker
)

// Tile is a struct to hold the tile type so we can create a matrix of tiles
type Tile struct {
	Type TileType
}

// Game struct for maintaining state for a game instance
type Game struct {
	Name      string
	Redraw    chan struct{}
	HighScore int

	// Top left is 0,0
	level [][]Tile
	hub   Hub
}

// NewGame generates a new game instance
func NewGame(worldWidth int, worldHeight int) *Game {
	g := &Game{
		Name:   petname.Generate(1, ""),
		Redraw: make(chan struct{}),
		hub:    NewHub(),
	}
	g.initalizeLevel(worldWidth, worldHeight)

	return g
}

func (g *Game) initalizeLevel(width int, height int) {
	g.level = make([][]Tile, width)
	for x := range g.level {
		g.level[x] = make([]Tile, height)
	}
	// default to Grass
	for x := range g.level {
		for y := range g.level[x] {
			g.setTileType(position.PositionFromInt(x, y), TileGrass)
		}
	}
	return
}

func (g *Game) setTileType(pos position.Position, tileType TileType) error {
	outOfBoundsErr := "Given %s value (%s) out of bounds"
	if pos.RoundX() > len(g.level) || pos.RoundX() < 0 {
		return fmt.Errorf(outOfBoundsErr, "X", pos.X)
	}
	// this used to be an else if, but that would mean that if the X was valid, it would never check the Y
	if pos.RoundY() > len(g.level[pos.RoundX()]) || pos.RoundY() < 0 {
		return fmt.Errorf(outOfBoundsErr, "Y", pos.Y)
	}

	g.level[pos.RoundX()][pos.RoundY()].Type = tileType
	return nil

}

func (g *Game) players() map[*player.Player]*session.Session {
	players := make(map[*player.Player]*session.Session)
	for session := range g.hub.Sessions {
		players[session.Player] = session
	}
	return players
}

// WorldWidth the width of the world (len of the x array in [x][y]Tiles)
func (g *Game) WorldWidth() int {
	return len(g.level)
}

// WorldHeight the height of the world (len of the y array in [x][y]Tiles)
func (g *Game) WorldHeight() int {
	return len(g.level[0])
}

// SessionsCount the number of active sessions in a game
func (g *Game) SessionsCount() int {
	return len(g.hub.Sessions)
}

// AddSession calls to the game hubs register channel to add a session
func (g *Game) AddSession(s *session.Session) {
	g.hub.Register <- s
	return
}

// RemoveSession calls to the game hubs unregister channel to remove a session
func (g *Game) RemoveSession(s *session.Session) {
	g.hub.Unregister <- s
	return
}

// Render renders a game world and sends it to the player over the session
func (g *Game) Render(s *session.Session) {
	worldstr := g.worldString(s)

	var b bytes.Buffer
	b.WriteString("\033[H\033[2J")
	b.WriteString(worldstr)

	// send rendered world to player
	io.Copy(s, &b)
}

// AvailableColors returns the colors that are available for players to use in the current game session
func (g *Game) AvailableColors() []color.Attribute {
	usedColors := map[color.Attribute]bool{}
	for _, color := range player.PlayerColors {
		usedColors[color] = false
	}

	for player := range g.players() {
		usedColors[player.Color] = true
	}

	availableColors := []color.Attribute{}
	for color, used := range usedColors {
		if !used {
			availableColors = append(availableColors, color)
		}
	}

	return availableColors
}

// Only works with square worlds
func (g *Game) worldString(s *session.Session) string {
	worldWidth := g.WorldWidth()
	worldHeight := g.WorldHeight()

	/* create 2d slice of strings to represent the world
	two chars longer in each direction to accomodate for walls
	*/
	strWorld := make([][]string, worldWidth+2)
	for x := range strWorld {
		strWorld[x] = make([]string, worldHeight+2)
	}

	// load walls into rune slice
	borderColorizer := color.New(player.PlayerBorderColors[s.Player.Color]).SprintFunc()
	for x := range strWorld {
		strWorld[x][0] = borderColorizer(string(HorizontalWall))
		strWorld[x][worldHeight+1] = borderColorizer(string(HorizontalWall))
	}
	for y := range strWorld[0] {
		strWorld[0][y] = borderColorizer(string(VerticalWall))
		strWorld[worldWidth+1][y] = borderColorizer(string(VerticalWall))
	}

	// colorize corners
	strWorld[0][0] = borderColorizer(string(TopLeft))
	strWorld[worldWidth+1][0] = borderColorizer(string(TopRight))
	strWorld[0][worldHeight+1] = borderColorizer(string(BottomLeft))
	strWorld[worldWidth+1][worldHeight+1] = borderColorizer(string(BottomRight))

	// Draw the player's score
	scoreStr := fmt.Sprintf(
		" Score: %d : Your High Score: %d : Game High Score: %d ",
		s.Player.Score(),
		s.HighScore,
		g.HighScore,
	)
	for i, r := range scoreStr {
		strWorld[3+i][0] = borderColorizer(string(r))
	}

	// Draw the player's color
	colorStr := fmt.Sprintf(" %s ", player.PlayerColorNames[s.Player.Color])
	colorStrColorizer := color.New(s.Player.Color).SprintFunc()
	for i, r := range colorStr {
		charsRemaining := len(colorStr) - i
		strWorld[len(strWorld)-3-charsRemaining][0] = colorStrColorizer(string(r))
	}

	// Draw everyone's scores
	if len(g.players()) > 1 {
		// Sort the players by color name
		players := []*player.Player{}

		for player := range g.players() {
			if player == s.Player {
				continue
			}

			players = append(players, player)
		}

		sort.Sort(player.ByColor(players))
		startX := 3

		// Actually draw their scores
		for _, p := range players {
			colorizer := color.New(p.Color).SprintFunc()
			scoreStr := fmt.Sprintf(" %s: %d",
				player.PlayerColorNames[p.Color],
				p.Score(),
			)
			for _, r := range scoreStr {
				strWorld[startX][len(strWorld[0])-1] = colorizer(string(r))
				startX++
			}
		}

		// Add final spacing next to wall
		strWorld[startX][len(strWorld[0])-1] = " "
	} else {
		warning :=
			" Warning: Other Players Must be in This Game for You to Score! "
		for i, r := range warning {
			strWorld[3+i][len(strWorld[0])-1] = borderColorizer(string(r))
		}
	}

	// Draw the game's name
	nameStr := fmt.Sprintf(" %s ", g.Name)
	for i, r := range nameStr {
		charsRemaining := len(nameStr) - i
		strWorld[len(strWorld)-3-charsRemaining][len(strWorld[0])-1] =
			borderColorizer(string(r))
	}

	// Load the level into the string slice
	for x := 0; x < worldWidth; x++ {
		for y := 0; y < worldHeight; y++ {
			tile := g.level[x][y]

			switch tile.Type {
			case TileGrass:
				strWorld[x+1][y+1] = string(Grass)
			case TileBlocker:
				strWorld[x+1][y+1] = string(Blocker)
			}
		}
	}

	// Load the players into the rune slice
	for player := range g.players() {
		colorizer := color.New(player.Color).SprintFunc()

		pos := player.Pos
		strWorld[pos.RoundX()+1][pos.RoundY()+1] = colorizer(string(player.Marker))

		// Load the player's trail into the rune slice
		for _, segment := range player.Trail {
			x, y := segment.Pos.RoundX()+1, segment.Pos.RoundY()+1
			strWorld[x][y] = colorizer(string(segment.Marker))
		}
	}

	// Convert the rune slice to a string
	buffer := bytes.NewBuffer(make([]byte, 0, worldWidth*worldHeight*2))
	for y := 0; y < len(strWorld[0]); y++ {
		for x := 0; x < len(strWorld); x++ {
			buffer.WriteString(strWorld[x][y])
		}

		// Don't add an extra newline if we're on the last iteration
		if y != len(strWorld[0])-1 {
			buffer.WriteString("\r\n")
		}
	}

	return buffer.String()
}

// Run runs a game instance
func (g *Game) Run() {
	go func() {
		for {
			g.hub.Redraw <- <-g.Redraw
		}
	}()

	// run game loop
	go func() {
		var lastUpdate time.Time
		c := time.Tick(time.Second / 60) //TODO: make this configurable
		for now := range c {
			g.Update(float64(now.Sub(lastUpdate)) / float64(time.Millisecond))
			lastUpdate = now
		}
	}()

	// redraw
	// potential optimization: use diffs to only redraw when needed
	go func() {
		c := time.Tick(time.Second / 10)
		for range c {
			g.Redraw <- struct{}{}
		}
	}()

	g.hub.HubRun(g)
}

// Update is the main game logic loop. delta is the time since the last update in milliseconds
func (g *Game) Update(delta float64) {
	// set of all coordinates occupied by trails
	trailCoordinateMap := make(map[string]bool)

	// update player data
	for p, s := range g.players() {
		p.Update(len(g.players()), delta)

		// update session high score
		if p.Score() > s.HighScore {
			s.HighScore = p.Score()
		}
		// Update global high score
		if p.Score() > g.HighScore {
			g.HighScore = p.Score()
		}

		// handle out of player out of bounds
		if p.IsOutOfBounds(0, len(g.level), 0, len(g.level[0])) {
			s.StartOver(g.WorldWidth(), g.WorldHeight())
		}

		// Kick the player if they've timed out
		if time.Now().Sub(s.LastAction) > player.PlayerTimeout {
			fmt.Fprint(s, "\r\n\r\nYou were terminated due to inactivity\r\n")
			g.RemoveSession(s)
			return
		}

		for _, seg := range p.Trail {
			coordStr := fmt.Sprintf("%d,%d", seg.Pos.RoundX(), seg.Pos.RoundY())
			trailCoordinateMap[coordStr] = true
		}
	}

	// Check if any players collide with a trail and restart them if so
	for p, s := range g.players() {
		playerPos := fmt.Sprintf("%d,%d", p.Pos.RoundX(), p.Pos.RoundY())
		if collided := trailCoordinateMap[playerPos]; collided {
			s.StartOver(g.WorldWidth(), g.WorldHeight())
		}
	}
}

// SetGameServerProperties reads cfg.Game.Server.* and overrides the default player
// properties with values in the configuration json if set
// TODO: There must be a better way to do this?
func SetGameServerProperties(cfg *config.Config) {
	if cfg.Game.Server.VerticalWall != nil {
		VerticalWall, _ = utf8.DecodeRuneInString(*cfg.Game.Server.VerticalWall)
	}
	if cfg.Game.Server.HorizontalWall != nil {
		HorizontalWall, _ = utf8.DecodeRuneInString(*cfg.Game.Server.HorizontalWall)
	}
	if cfg.Game.Server.TopLeft != nil {
		TopLeft, _ = utf8.DecodeRuneInString(*cfg.Game.Server.TopLeft)
	}
	if cfg.Game.Server.TopRight != nil {
		TopRight, _ = utf8.DecodeRuneInString(*cfg.Game.Server.TopRight)
	}
	if cfg.Game.Server.BottomRight != nil {
		BottomRight, _ = utf8.DecodeRuneInString(*cfg.Game.Server.BottomRight)
	}
	if cfg.Game.Server.BottomLeft != nil {
		BottomLeft, _ = utf8.DecodeRuneInString(*cfg.Game.Server.BottomLeft)
	}
	if cfg.Game.Server.Grass != nil {
		Grass, _ = utf8.DecodeRuneInString(*cfg.Game.Server.Grass)
	}
	if cfg.Game.Server.Blocker != nil {
		Blocker, _ = utf8.DecodeRuneInString(*cfg.Game.Server.Blocker)
	}
	return
}
