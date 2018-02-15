package config

import "time"

// PlayerConf is the config struct for player properties
type PlayerConf struct {
	VerticalSpeed        *float64
	HorizontalSpeed      *float64
	CountScoreMultiplier *float64
	TimeoutSeconds       *int
	UpRune               *string
	DownRune             *string
	LeftRune             *string
	RightRune            *string
	TrailHorizontal      *string
	TrailVertical        *string
	TrailLeftCornerUp    *string
	TrailLeftCornerDown  *string
	TrailRightCornerDown *string
	TrailRightCornerUp   *string

	PlayerTrailLengthLimit *bool
	LimitPlayerTrailByTime *bool
	PlayerTrailMaxLength   *int
	PlayerTrailMaxTime     *int
}

// default values if not provided in config file
var (
	VerticalPlayerSpeed        = 0.007
	HorizontalPlayerSpeed      = 0.01
	PlayerCountScoreMultiplier = 1.25
	PlayerTimeout              = 15 * time.Second

	PlayerUpRune    = '⇡'
	PlayerDownRune  = '⇣'
	PlayerLeftRune  = '⇠'
	PlayerRightRune = '⇢'

	PlayerTrailHorizontal      = '┄'
	PlayerTrailVertical        = '┆'
	PlayerTrailLeftCornerUp    = '╭'
	PlayerTrailLeftCornerDown  = '╰'
	PlayerTrailRightCornerDown = '╯'
	PlayerTrailRightCornerUp   = '╮'

	PlayerTrailLengthLimit = true
	LimitPlayerTrailByTime = false
	PlayerTrailMaxLength   = 20
	PlayerTrailMaxTime     = 3
)
