package main

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/strosel/sprak/cards"
	"golang.org/x/image/colornames"
)

func drawCard(th *material.Theme, card cards.Card) layout.Widget {
	text := card.Q
	if flip {
		text = card.A
	}

	bttn := material.Button(th, bttns["card"], text)
	bttn.TextSize = fontSize
	bttn.Background = colornames.Gainsboro
	bttn.Color = colornames.Black
	for bttns["card"].Clicked() {
		flip = !flip
	}
	return bttn.Layout
}

func drawIncorrect(th *material.Theme, card cards.Card) layout.Widget {
	bttn := material.IconButton(th, bttns["incorrect"], icns["incorrect"])
	bttn.Background = colornames.Indianred
	bttn.Size = unit.Dp(100)
	bttn.Color = colornames.Black
	for bttns["incorrect"].Clicked() {
		card.Update(false)
	}
	return bttn.Layout
}

func drawCorrect(th *material.Theme, card cards.Card) layout.Widget {
	bttn := material.IconButton(th, bttns["correct"], icns["correct"])
	bttn.Background = colornames.Mediumseagreen
	bttn.Size = unit.Dp(100)
	bttn.Color = colornames.Black
	for bttns["correct"].Clicked() {
		card.Update(true)
	}
	return bttn.Layout
}

func drawAns(th *material.Theme, card cards.Card) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.UniformInset(unit.Dp(16)).Layout(gtx,
			func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceAround,
				}.Layout(gtx,
					layout.Rigid(drawIncorrect(th, card)),
					layout.Rigid(drawCorrect(th, card)),
				)
			})
	}
}

func train(gtx layout.Context, th *material.Theme) {
	card := cards.Card{Q: "Q", A: "A"}
	if cardCount == 0 {
		deck.Shuffle()
		// card = deck.Pick()
		cardCount++
	}

	layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.End,
		}.Layout(gtx,
			layout.Flexed(1, drawCard(th, card)),
			layout.Rigid(drawAns(th, card)),
		)
	})
}
