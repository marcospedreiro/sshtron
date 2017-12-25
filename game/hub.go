package game

import (
	"fmt"

	"github.com/marcospedreiro/sshtron/session"
)

// Hub struct for holding game state
type Hub struct {
	Sessions   map[*session.Session]struct{}
	Redraw     chan struct{}
	Register   chan *session.Session
	Unregister chan *session.Session
}

// NewHub Initializes an empty Hub struct
func NewHub() Hub {
	return Hub{
		Sessions:   make(map[*session.Session]struct{}),
		Redraw:     make(chan struct{}),
		Register:   make(chan *session.Session),
		Unregister: make(chan *session.Session),
	}
}

// HubRun runs the game hub for more meta controls
func (h *Hub) HubRun(g *Game) {
	for {
		select {
		case <-h.Redraw:
			for s := range h.Sessions {
				go g.Render(s)
			}
		case s := <-h.Register:
			// Hide the cursor
			fmt.Fprint(s, "\033[?25l")

			h.Sessions[s] = struct{}{}
		case s := <-h.Unregister:
			if _, ok := h.Sessions[s]; ok {
				fmt.Fprint(s, "\r\n\r\n~ End of Line ~ \r\n\r\nRemember to use WASD to move!\r\n\r\n")

				// Unhide the cursor
				fmt.Fprint(s, "\033[?25h")

				delete(h.Sessions, s)
				s.C.Close()
			}
		}
	}
}
