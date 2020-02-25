package main

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"golang.org/x/image/colornames"
)

func train(gtx *layout.Context, th *material.Theme) {
	if cardCount == 0 {
		deck.Shuffle()
		card = deck.Pick()
		cardCount++
	}

	layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.End,
		}.Layout(gtx,
			layout.Flexed(1, func() {
				s := card.Q
				if flip {
					s = card.A
				}

				bttn := th.Button(s)
				bttn.TextSize = fontSize
				bttn.Background = colornames.Gainsboro
				bttn.Color = colornames.Black
				for bttns["card"].Clicked(gtx) {
					flip = !flip
				}
				bttn.Layout(gtx, bttns["card"])
			}),
			layout.Rigid(func() {
				layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceAround,
				}.Layout(gtx,
					layout.Rigid(func() {
						bttn := th.IconButton(icns["incorrect"])
						bttn.Background = colornames.Indianred
						bttn.Size = unit.Dp(100)
						bttn.Color = colornames.Black
						for bttns["incorrect"].Clicked(gtx) {
							card.Update(false)
						}
						bttn.Layout(gtx, bttns["incorrect"])
					}),
					layout.Rigid(func() {
						bttn := th.IconButton(icns["correct"])
						bttn.Background = colornames.Mediumseagreen
						bttn.Size = unit.Dp(100)
						bttn.Color = colornames.Black
						for bttns["correct"].Clicked(gtx) {
							card.Update(true)
						}
						bttn.Layout(gtx, bttns["correct"])
					}))
			}),
		)
	})
}
