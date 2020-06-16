package cards

import (
	"time"
)

type Card struct {
	Q          string    `bson:"question" json:"question"`
	A          string    `bson:"answer" json:"answer"`
	Categories []string  `bson:"categories,omitempty" json:"categories,omitempty"`
	Trained    int       `bson:"trained" json:"trained"`
	Success    int       `bson:"success" json:"success"`
	Last       time.Time `bson:"last" json:"last"`
}

func (c Card) SuccessRate() float64 {
	if c.Trained == 0 {
		return 0.
	}
	return float64(c.Success) / float64(c.Trained)
}

func (c *Card) Update(correct bool) {
	c.Last = time.Now()
	c.Trained++
	if correct {
		c.Success++
	}
}
