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

	pos         Pos
	speed       int64
	lastEventTs int64

	moveXCh   chan struct{}
	moveYCh   chan struct{}
	moving    *moving
	walkingCh chan struct{}

	animationStandLR   [][]rune
	animationStandRL   [][]rune
	animationWalkingLR [][][]rune
	animationWalkingRL [][][]rune
}

const (
	left int = iota
	right
)

type moving struct {
	mu sync.Mutex
	x  bool
	y  bool
}

func (m *moving) setx(b bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.x = b
}

func (m *moving) sety(b bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.y = b
}

func (m *moving) get() (bool, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.x, m.y
}

type Pos struct {
	x int
	y int
}

func New() *Player {
	player := &Player{
		pos:         Pos{x: 25, y: 54},
		speed:       1200,
		style:       tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack),
		lastEventTs: time.Now().Unix(),
		moveXCh:     make(chan struct{}),
		moveYCh:     make(chan struct{}),
		moving:      new(moving),
	}

	player.animationStandLR = common.OpenRuneFile("data/helley.stay.lr.1")
	player.animationStandRL = common.ReflectHorizontal(common.OpenRuneFile("data/helley.stay.lr.1"))
	player.animationWalkingLR, player.animationWalkingRL = common.OpenRuneFiles(
		"data/helley.walk.lr.1", "data/helley.walk.lr.2", "data/helley.walk.lr.3", "data/helley.walk.lr.4")

	player.runes = player.animationStandLR

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

func (p *Player) stop(direction int) {

	p.lastEventTs = time.Now().Unix()

	xb, yb := p.moving.get()
	if xb {
		p.moveXCh <- struct{}{}
	}
	if yb {
		p.moveYCh <- struct{}{}
	}

	if p.walkingCh != nil {
		p.walkingAnimationStop(direction)
	}
}

func (p *Player) Move(y, x int) {

	p.lastEventTs = time.Now().Unix()

	xb, yb := p.moving.get()

	var direction int
	if y >= p.pos.y {
		direction = right
	} else {
		direction = left
	}

	if xb || yb {
		p.stop(direction)
		return
	}
	p.moving.setx(true)
	p.moving.sety(true)

	move := func() {
		p.walkingCh = p.walkingAnimation(direction)
		go func() {
			for {
				select {
				case <-p.moveXCh:
					p.moving.setx(false)
					if _, yb := p.moving.get(); !yb {
						p.walkingAnimationStop(direction)
					}
					return
				case <-time.After(time.Duration(p.speed) * time.Millisecond):
					if x == p.pos.x {
						p.moving.setx(false)
						if _, yb := p.moving.get(); !yb {
							p.walkingAnimationStop(direction)
						}
						return
					}
					if x != p.pos.x {
						if x > p.pos.x {
							p.pos.x += 1
						} else {
							p.pos.x -= 1
						}
					}
				}
			}
		}()

		go func() {
			for {
				select {
				case <-p.moveYCh:
					p.moving.sety(false)
					if xb, _ := p.moving.get(); !xb {
						p.walkingAnimationStop(direction)
					}
					return
				case <-time.After(time.Duration(p.speed/6) * time.Millisecond):
					if y == p.pos.y {
						p.moving.sety(false)
						if xb, _ := p.moving.get(); !xb {
							p.walkingAnimationStop(direction)
						}
						return
					}
					if y != p.pos.y {
						if y > p.pos.y {
							p.pos.y += 1
						} else {
							p.pos.y -= 1
						}
					}
				}
			}
		}()
	}

	move()
}

func replace(runes []rune, old, new string) []rune {
	return []rune(strings.ReplaceAll(string(runes), old, new))
}

func (p *Player) winking() chan struct{} {
	ch := make(chan struct{})
	go func() {
		for {
			select {
			case <-ch:
				return
			default:
				rand.Seed(time.Now().Unix())
				s := rand.Intn(2) + 1
				time.Sleep(time.Duration(s) * time.Second)
				p.runes[1] = replace(p.runes[1], "Oo", "--")
				time.Sleep(time.Duration(500) * time.Millisecond)
				p.runes[1] = replace(p.runes[1], "--", "Oo")
			}
		}
	}()
	return ch
}

func (p *Player) mouth() chan struct{} {
	ch := make(chan struct{})
	go func() {
		texts := []string{"=", "e", "a", "~"}
		for {
			select {
			case <-ch:
				return
			default:
				rand.Seed(time.Now().Unix())
				s := rand.Intn(6) + 1
				time.Sleep(time.Duration(s) * time.Second)
				text := texts[rand.Intn(len(texts))]
				p.runes[2] = replace(p.runes[2], "-", text)
				time.Sleep(time.Duration(500) * time.Millisecond)
				p.runes[2] = replace(p.runes[2], text, "-")
			}
		}
	}()
	return ch
}

func (p *Player) text() chan struct{} {
	ch := make(chan struct{})
	go func() {
		texts := [][]rune{[]rune("    - ??? "), []rune("    - Huh? "), []rune("    - Who am I? Where am I? ")}
		for {
			select {
			case <-ch:
				return
			default:
				rand.Seed(time.Now().Unix())
				s := rand.Intn(5) + 1
				time.Sleep(time.Duration(s) * time.Second)
				text := texts[rand.Intn(len(texts))]
				p.runes[0] = append(p.runes[0], text...)
				time.Sleep(time.Duration(3) * time.Second)
				p.runes[0] = replace(p.runes[0], string(text), "")
			}
		}
	}()
	return ch
}

func (p *Player) inactivity() {

	go func() {
		var inactive bool
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		var winkingCh chan struct{}
		var sayingCh chan struct{}

		for {
			select {
			case <-ticker.C:
				if time.Unix(p.lastEventTs, 0).Add(5 * time.Second).Before(time.Now()) {
					if !inactive {
						inactive = true
						winkingCh = p.winking()
						sayingCh = p.saying()
					}
				} else if inactive {
					winkingCh <- struct{}{}
					sayingCh <- struct{}{}
					inactive = false
				}
			}
		}
	}()
}

func (p *Player) saying() chan struct{} {
	ch := make(chan struct{})
	mch := p.mouth()
	tch := p.text()
	go func() {
		for {
			select {
			case <-ch:
				mch <- struct{}{}
				tch <- struct{}{}
			}
		}
	}()
	return ch
}

func (p *Player) walkingAnimation(direction int) chan struct{} {

	// p.winking()
	ch := make(chan struct{})

	go func() {

		var runes [][][]rune
		switch direction {
		case left:
			runes = p.animationWalkingRL
		case right:
			runes = p.animationWalkingLR
		}

		var i int
		for {
			select {
			case <-ch:
				return
			default:
				p.runes = runes[i]
				time.Sleep(time.Duration(350) * time.Millisecond)
				if i == len(runes)-1 {
					i = 0
				} else {
					i++
				}
			}
		}
	}()
	return ch
}

func (p *Player) walkingAnimationStop(direction int) {
	p.walkingCh <- struct{}{}
	switch direction {
	case left:
		p.runes = p.animationStandRL
	case right:
		p.runes = p.animationStandLR
	}
}
