package cards

import (
	"math/rand"
	"time"
)

var timeCap = time.Hour * 24 * 7

type Deck struct {
	Name     string `bson:"name" json:"name"`
	Cards    []Card `bson:"cards" json:"cards"`
	Pool     []int  `bson:"-" json:"-"`
	Shuffled bool   `bson:"-" json:"-"`
}

func (d *Deck) Shuffle() {
	//setup the pool aka the list of card indecies with apperance by odds
	//should be called after any card if !shuffled or after n cards if shuffled

	d.Pool = []int{}
	d.Shuffled = true

	for i, c := range d.Cards {
		if time.Now().Sub(c.Last) > timeCap {
			d.Pool = []int{i}
			d.Shuffled = false
			return
		}
		for j := 0; j < 101-int(c.SuccessRate()*100); j++ {
			d.Pool = append(d.Pool, j)
		}
	}

}

func (d Deck) Categories() []string {
	//use a map to prevent doubles witout nested loops
	catm := map[string]bool{}

	for _, card := range d.Cards {
		for _, c := range card.Categories {
			catm[c] = true
		}
	}

	cat := []string{}

	for c := range catm {
		cat = append(cat, c)
	}

	return cat
}

func (d Deck) Pick() Card {
	if len(d.Pool) == 0 {
		return Card{}
	}
	return d.Cards[d.Pool[rand.Intn(len(d.Pool))]]
}
