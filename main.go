package main

import (
	"log"
	"math/rand"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/strosel/sprak/cards"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Screen int

const (
	listS Screen = iota
	trainS
	deckS //new and edit
	cardS
)

var (
	fontSize  = unit.Dp(32)
	cardCount = 0
	flip      = false

	bttns = map[string]*widget.Button{
		"card":      new(widget.Button),
		"correct":   new(widget.Button),
		"incorrect": new(widget.Button),
	}
	icns = map[string]*material.Icon{
		"correct":   new(material.Icon),
		"incorrect": new(material.Icon),
	}

	deck cards.Deck
	card *cards.Card
	err  error
)

func main() {
	rand.Seed(time.Now().UnixNano())
	gofont.Register()
	trainS.Swap()

	icns["correct"], err = material.NewIcon(icons.NavigationCheck)
	if err != nil {
		log.Fatal(err)
	}
	icns["incorrect"], err = material.NewIcon(icons.NavigationClose)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		win := app.NewWindow()
		if err := loop(win); err != nil {
			log.Fatal(err)
		}
	}()

	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	gtx := layout.NewContext(w.Queue())
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case *system.CommandEvent:
			e.Cancel = true
			w.Invalidate()
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)

			train(gtx, th)

			e.Frame(gtx.Ops)
		}
	}
}

func (s Screen) Swap() {
	switch s {
	case trainS:
		cardCount = 0
		//TODO
		deck = cards.Deck{}
	}
}
