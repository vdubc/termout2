package room

import (
	"github.com/gdamore/tcell/v2"
	"github.com/vdubc/termout2/common"
)

type Room struct {
	Runes [][]rune
	style tcell.Style
}

func New() *Room {
	room := &Room{Runes: common.OpenRuneFile("data/room1")}
	room.style = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	return room
}

func (r *Room) Show(screen tcell.Screen) {
	// room
	for x, xs := range r.Runes {
		for y, ys := range xs {
			screen.SetContent(y, x, ys, nil, r.style)
		}
	}
}
