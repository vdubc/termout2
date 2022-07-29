package player

import (
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/vdubc/termout2/common"
)

type Player struct {
	mu    sync.Mutex
	runes [][]rune
	style tcell.Style

	pos   Pos
	speed int64

	moveCh chan struct{}
}

type Pos struct {
	x int
	y int
}

func New() *Player {
	player := &Player{
		runes:  common.OpenRuneFile("data/player1"),
		pos:    Pos{x: 25, y: 54},
		speed:  1200,
		style:  tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack),
		moveCh: make(chan struct{}),
	}
	// wrap spaces left/right // TODO

	player.inactivity()

	return player
}

func (p *Player) Show(screen tcell.Screen) {
	// player
	for x, xs := range p.runes {
		for y, ys := range xs {
			yn := p.pos.y + y /*- len(p.Player[0])/2*/
			xn := p.pos.x + x - len(p.runes)
			screen.SetContent(yn, xn, ys, nil, p.style)
		}
	}

}

func (p *Player) Move(y, x int) {

	// p.moveCh <- struct{}{}

	// xCh := make(chan struct{})
	// go func(ch chan struct{}) {
	// 	defer close(ch)
	// 	for {
	// 		select {
	// 		case <-ch:
	// 			break
	// 		default:
	// 			if x == p.Pos.X {
	// 				break
	// 			}
	// 			if x != p.Pos.X {
	// 				if x > p.Pos.X {
	// 					p.Pos.X += 1
	// 				} else {
	// 					p.Pos.X -= 1
	// 				}
	// 			}
	// 			time.Sleep(time.Duration(p.speed) * time.Millisecond) // speed
	// 		}
	// 	}
	// }(xCh)

	// yCh := make(chan struct{})
	// go func(ch chan struct{}) {
	// 	defer close(ch)
	// 	for {
	// 		select {
	// 		case <-ch:
	// 			break
	// 		default:
	// 			if y == p.Pos.Y {
	// 				break
	// 			}
	// 			if y != p.Pos.Y {
	// 				if y > p.Pos.Y {
	// 					p.Pos.Y += 1
	// 				} else {
	// 					p.Pos.Y -= 1
	// 				}
	// 				time.Sleep(time.Duration(p.speed/6) * time.Millisecond) // speed
	// 			}
	// 		}

	// 	}
	// }(yCh)

	// go func() {
	// 	select {
	// 	case <-p.moveCh:
	// 		xCh <- struct{}{}
	// 		yCh <- struct{}{}
	// 	}
	// }()

}

func replace(runes []rune, old, new string) []rune {
	return []rune(strings.ReplaceAll(string(runes), old, new))
}

func (p *Player) winking() {
	go func() {
		for {
			rand.Seed(time.Now().Unix())
			s := rand.Intn(2) + 1
			time.Sleep(time.Duration(s) * time.Second)
			// p.mu.Lock()
			p.runes[1] = replace(p.runes[1], "Oo", "--")
			// p.mu.Unlock()
			time.Sleep(time.Duration(500) * time.Millisecond)
			// p.mu.Lock()
			p.runes[1] = replace(p.runes[1], "--", "Oo")
			// p.mu.Unlock()
		}
	}()
}

func (p *Player) mouth() {
	go func() {
		texts := []string{"=", "e", "a", "~"}
		for {
			rand.Seed(time.Now().Unix())
			s := rand.Intn(6) + 1
			time.Sleep(time.Duration(s) * time.Second)
			text := texts[rand.Intn(len(texts))]
			// p.mu.Lock()
			p.runes[2] = replace(p.runes[2], "-", text)
			// p.mu.Unlock()
			time.Sleep(time.Duration(500) * time.Millisecond)
			// p.mu.Lock()
			p.runes[2] = replace(p.runes[2], text, "-")
			// p.mu.Unlock()
		}
	}()
}

func (p *Player) text() {
	go func() {
		texts := [][]rune{[]rune("    - ??? "), []rune("    - Huh? "), []rune("    - Who am I? Where am I? ")}
		for {

			rand.Seed(time.Now().Unix())
			s := rand.Intn(5) + 1
			time.Sleep(time.Duration(s) * time.Second)
			text := texts[rand.Intn(len(texts))]
			// p.mu.Lock()
			p.runes[0] = append(p.runes[0], text...)
			// p.mu.Unlock()
			time.Sleep(time.Duration(3) * time.Second)
			// p.mu.Lock()
			p.runes[0] = replace(p.runes[0], string(text), "")
			// p.mu.Unlock()
		}
	}()
}

func (p *Player) inactivity() {
	// for {
	// 	// p.winking()
	// 	// p.saying()
	// }

}

func (p *Player) saying() {
	p.mouth()
	p.text()
}

func (p *Player) walking() {

	// p.winking()

	var a [][][]rune
	for _, n := range []string{"player.walk.rl.1", "player.walk.rl.2", "player.walk.rl.3", "player.walk.rl.4"} {

		data := common.OpenRuneFile("data/" + n)
		for i, rs := range data {
			data[i] = []rune(strings.ReplaceAll(string(rs), "░", " "))
		}

		// a = append(a, data)
		a = append(a, common.ReflectHorizontal(data))
	}
	go func() {
		var i int
		for {
			// p.mu.Lock()
			p.runes = a[i]
			// p.mu.Unlock()
			time.Sleep(time.Duration(350) * time.Millisecond)
			if i == len(a)-1 {
				i = 0
			} else {
				i++
			}
		}
	}()
}
