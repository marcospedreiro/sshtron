package config

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
}
