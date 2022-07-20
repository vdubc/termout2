package scene

import (
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/vdubc/termout2/player"
)

type Scene struct {
	screen tcell.Screen
	player *player.Player
}

func New() *Scene {

	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err) // TODO
	}
	if err := screen.Init(); err != nil {
		panic(err) // TODO
	}

	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	screen.SetStyle(style)
	screen.EnableMouse()
	screen.Clear()

	return &Scene{screen: screen}
}

func (s *Scene) Add(player *player.Player) {
	s.player = player
}

func (s *Scene) Run() {

	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack) // TODO

	quitFn := func() {
		s.screen.Fini()
		os.Exit(0)
	}

	evch := make(chan tcell.Event)
	quit := make(chan struct{})
	go s.screen.ChannelEvents(evch, quit)

	for {
		select {
		case <-quit:
			quitFn()
			break
		case ev := <-evch:

			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.screen.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					quit <- struct{}{}
					quitFn()
				} else if ev.Key() == tcell.KeyCtrlL {
					s.screen.Sync()
				} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
					s.screen.Clear()
				}
			case *tcell.EventMouse:
				switch ev.Buttons() {
				case tcell.Button1:
					y, x := ev.Position()
					s.player.Move(y, x)
				}
			}

		case <-time.After(time.Millisecond * 50):

			// room
			for x, xs := range s.player.Room {
				for y, ys := range xs {
					s.screen.SetContent(y, x, ys, nil, style)
				}
			}

			// player
			for x, xs := range s.player.Player {
				for y, ys := range xs {
					s.screen.SetContent(s.player.Pos.Y+y, s.player.Pos.X+x-len(s.player.Player), ys, nil, style)
				}
			}

			// Update screen
			s.screen.Show()

		}
	}

	// for {

	// 	// Poll event
	// 	ev := s.screen.PollEvent()

	// 	// Process event
	// 	switch ev := ev.(type) {
	// 	case *tcell.EventResize:
	// 		s.screen.Sync()
	// 	case *tcell.EventKey:
	// 		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
	// 			quit()
	// 		} else if ev.Key() == tcell.KeyCtrlL {
	// 			s.screen.Sync()
	// 		} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
	// 			s.screen.Clear()
	// 		}
	// 	case *tcell.EventMouse:
	// 	}
	// 	// case <-time.After(time.Millisecond * 50):

	// }
}
