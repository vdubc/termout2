package scene

import (
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/vdubc/termout2/player"
	"github.com/vdubc/termout2/room"
)

type Scene struct {
	screen tcell.Screen
	room   *room.Room
	player *player.Player
}

func New(room *room.Room, player *player.Player) *Scene {

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

	return &Scene{screen: screen, room: room, player: player}
}

func (s *Scene) Run() {
	quitFn := func() {
		s.screen.Fini()
		os.Exit(0)
	}

	evntCh := make(chan tcell.Event)
	envtChQuit := make(chan struct{})
	go s.screen.ChannelEvents(evntCh, envtChQuit)

	mainChQuit := make(chan struct{})
	go func(evntCh chan tcell.Event, envtChQuit chan struct{}, mainChQuit chan struct{}) {
		for {
			select {
			case ev := <-evntCh:

				switch ev := ev.(type) {
				case *tcell.EventResize:
					s.screen.Sync()
				case *tcell.EventKey:
					if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
						envtChQuit <- struct{}{}
						mainChQuit <- struct{}{}
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
			}
		}
	}(evntCh, envtChQuit, mainChQuit)

	for {
		select {
		case <-mainChQuit:
			quitFn()

		case <-time.After(time.Millisecond * 50):
			s.room.Show(s.screen)
			s.player.Show(s.screen)
			s.screen.Show()
		}
	}

}
