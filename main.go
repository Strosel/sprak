package main

import (
	"log"
	"math/rand"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
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

	bttns = map[string]*widget.Clickable{
		"card":      new(widget.Clickable),
		"correct":   new(widget.Clickable),
		"incorrect": new(widget.Clickable),
	}
	icns = map[string]*widget.Icon{
		"correct":   new(widget.Icon),
		"incorrect": new(widget.Icon),
	}

	deck cards.Deck
	err  error
)

func main() {
	rand.Seed(time.Now().UnixNano())
	trainS.Swap()

	icns["correct"], err = widget.NewIcon(icons.NavigationCheck)
	if err != nil {
		log.Fatal(err)
	}
	icns["incorrect"], err = widget.NewIcon(icons.NavigationClose)
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
	th := material.NewTheme(gofont.Collection())
	var ops op.Ops

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case *system.CommandEvent:
			e.Cancel = true
			w.Invalidate()
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

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
