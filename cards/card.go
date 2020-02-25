package cards

import "time"

type Card struct {
	Q          string
	A          string
	Categories []string
	Trained    int
	Success    int
	Last       time.Time
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
