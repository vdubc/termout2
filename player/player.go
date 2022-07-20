package player

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type Player struct {
	mu     sync.Mutex
	Room   [][]rune
	Player [][]rune
	Pos    Pos
	Speed  int64

	Moving    bool
	MovingChX chan struct{}
	MovingChY chan struct{}
}

type Pos struct {
	X int
	Y int
}

func New() *Player {
	player := &Player{
		Room:   room(),
		Player: player(),
	}
	// wrap spaces left/right
	for i := range player.Player {
		player.Player[i] = append([]rune{' '}, player.Player[i]...)
		player.Player[i] = append(player.Player[i], ' ')
	}

	player.Pos.X = 25
	player.Pos.Y = 54
	player.Speed = 1200

	player.animations()
	return player
}

func room() (data [][]rune) {
	return open("data/room1")
}

func player() [][]rune {
	return open("data/player1")
}

func open(name string) (runes [][]rune) {
	file, err := os.Open(name)
	if err != nil {
		panic(err) // TODO
	}
	defer func() {
		if err = file.Close(); err != nil {
			panic(err) // TODO
		}
	}()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err) // TODO
	}

	lines := strings.Split(string(b), "\n")

	for _, line := range lines {
		if line == "\n" || line == "" {
			break
		}
		var xs []rune
		for _, c := range line {
			xs = append(xs, c)
		}
		runes = append(runes, xs)
	}

	return
}

func (p *Player) animations() {

	replace := func(runes []rune, old, new string) []rune {
		return []rune(strings.ReplaceAll(string(runes), old, new))
	}

	// winking
	go func() {
		for {
			rand.Seed(time.Now().Unix())
			s := rand.Intn(2) + 1
			time.Sleep(time.Duration(s) * time.Second)
			p.mu.Lock()
			p.Player[1] = replace(p.Player[1], "Oo", "--")
			p.mu.Unlock()
			time.Sleep(time.Duration(500) * time.Millisecond)
			p.mu.Lock()
			p.Player[1] = replace(p.Player[1], "--", "Oo")
			p.mu.Unlock()
		}
	}()

	// text
	go func() {
		texts := [][]rune{[]rune("    - ??? "), []rune("    - Huh? "), []rune("    - Who am I? Where am I? ")}
		for {
			if p.Moving {
				continue
			}
			rand.Seed(time.Now().Unix())
			s := rand.Intn(5) + 1
			time.Sleep(time.Duration(s) * time.Second)
			text := texts[rand.Intn(len(texts))]
			p.mu.Lock()
			p.Player[0] = append(p.Player[0], text...)
			p.mu.Unlock()
			time.Sleep(time.Duration(3) * time.Second)
			p.mu.Lock()
			p.Player[0] = replace(p.Player[0], string(text), "")
			p.mu.Unlock()
		}
	}()

	// mouth
	go func() {
		texts := []string{"=", "e", "a", "~"}
		for {
			if p.Moving {
				continue
			}
			rand.Seed(time.Now().Unix())
			s := rand.Intn(6) + 1
			time.Sleep(time.Duration(s) * time.Second)
			text := texts[rand.Intn(len(texts))]
			p.mu.Lock()
			p.Player[2] = replace(p.Player[2], "-", text)
			p.mu.Unlock()
			time.Sleep(time.Duration(500) * time.Millisecond)
			p.mu.Lock()
			p.Player[2] = replace(p.Player[2], text, "-")
			p.mu.Unlock()
		}
	}()
}

func (p *Player) Move(y, x int) {

	go func() {
		for {
			if x == p.Pos.X {
				break
			}
			if x != p.Pos.X {
				if x > p.Pos.X {
					p.Pos.X += 1
				} else {
					p.Pos.X -= 1
				}
			}
			time.Sleep(time.Duration(p.Speed) * time.Millisecond) // speed

		}
	}()
	go func() {
		for {
			if y == p.Pos.Y {
				break
			}
			if y != p.Pos.Y {
				if y > p.Pos.Y {
					p.Pos.Y += 1
				} else {
					p.Pos.Y -= 1
				}
				time.Sleep(time.Duration(p.Speed/6) * time.Millisecond) // speed
			}
		}
	}()
}
