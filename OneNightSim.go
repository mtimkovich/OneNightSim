package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Player struct {
	Num   int
	Start string
	Card  string
	game  *Game
	// maps player/position number to what they think card is
	Knowledge map[int]info
}

func (p *Player) Knows(position int, card string) {
	p.Knowledge[position] = p.game.info(card)
}

type info struct {
	card string
	// The higher the level, the more likely it is correct
	level int
}

func NewPlayer(n int, card string, game *Game) *Player {
	p := &Player{}
	p.Num = n
	p.Start = card
	p.Card = card
	p.game = game
	p.Knowledge = make(map[int]info)
	p.Knows(p.Num, p.Start)

	return p
}

func (p *Player) String() string {
	return fmt.Sprintf("Player %d: %s", p.Num, p.Card)
}

func shuffle(slice []string) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

type Game struct {
	Players []*Player
	Deck    []string
	Middle  []string
	level   int
}

func NewGame(numPlayers int) *Game {
	g := &Game{}

	// TODO: Handle more than 3 players and more roles
	g.Deck = []string{
		"Werewolf",
		"Werewolf",
		"Seer",
		"Robber",
		"Troublemaker",
		"Villager",
	}

	shuffle(g.Deck)

	for i := 0; i < numPlayers; i++ {
		g.Players = append(g.Players, NewPlayer(i, g.Deck[i], g))
	}

	g.Middle = g.Deck[3:]
	return g
}

func (g *Game) info(card string) info {
	return info{card, g.level}
}

// Pick a player other than given players
func (g *Game) pickOtherPlayer(us ...int) *Player {
	r := rand.Intn(len(g.Players) - len(us))

	for _, n := range us {
		if r >= n {
			r++
		}
	}

	return g.Players[r]
}

func (g *Game) troublemaker() {
	for _, p := range g.Players {
		if p.Start == "Troublemaker" {
			r := g.pickOtherPlayer(p.Num)
			s := g.pickOtherPlayer(p.Num, r.Num)

			r.Card, s.Card = s.Card, r.Card
			p.Knows(r.Num, "*Moved")
			p.Knows(s.Num, "*Moved")

			break
		}
	}
}

func (g *Game) robber() {
	for _, p := range g.Players {
		if p.Start == "Robber" {
			r := g.pickOtherPlayer(p.Num)
			p.Card, r.Card = r.Card, p.Card
			p.Knows(r.Num, p.Start)
			p.Knows(p.Num, p.Card)

			break
		}
	}
}

func (g *Game) seer() {
	for _, p := range g.Players {
		if p.Start == "Seer" {
			if rand.Intn(2) == 0 {
				// Pick a random player
				r := g.pickOtherPlayer(p.Num)

				p.Knows(r.Num, r.Card)
			} else {
				// Pick two random middle cards
				r := rand.Intn(3)
				for i := 0; i < 3; i++ {
					if i == r {
						continue
					}

					p.Knows(i+len(g.Players), g.Middle[i])
				}
			}

			break
		}
	}
}

func (g *Game) werewolves() {
	var wolves []int

	for _, p := range g.Players {
		if p.Start == "Werewolf" {
			wolves = append(wolves, p.Num)
		}
	}

	if len(wolves) == 2 {
		g.Players[wolves[0]].Knows(wolves[1], "Werewolf")
		g.Players[wolves[1]].Knows(wolves[0], "Werewolf")
	}

	for _, w := range wolves {
		for _, p := range g.Players {
			if p.Start != "Werewolf" {
				g.Players[w].Knows(p.Num, "!Werewolf")
			}
		}
	}
}

func (g *Game) Status() {
	for i, p := range g.Players {
		fmt.Println("num:  ", p.Num)
		fmt.Println("start:", p.Start)
		fmt.Println("end:  ", p.Card)
		fmt.Println("know: ", p.Knowledge)

		if i != len(g.Players)-1 {
			fmt.Println()
		}
	}
}

func (g *Game) Play() {
	// 1. Werewolves find each other
	g.werewolves()
	// 2. Seer looks at other player or 2 middle cards
	g.seer()
	g.level++
	// 3. Robber switches with player and looks at card
	g.robber()
	g.level++
	// 4. Troublemaker switches two other players
	g.troublemaker()
	// --- OPTIONAL ---
	// 5. Drunk switches with middle
	// 6. Insomniac looks at their own card
}

func main() {
	game := NewGame(3)
	game.Play()
	game.Status()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
