package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	ColorPrimary     = color.NRGBA{R: 0, G: 170, B: 152, A: 255}
	ColorPrimaryDark = color.NRGBA{R: 0, G: 130, B: 116, A: 255}
	ColorPrimaryLight = color.NRGBA{R: 200, G: 245, B: 238, A: 255}
)

func MakeCard(content fyne.CanvasObject) fyne.CanvasObject {
	return widget.NewCard("", "", content)
}

func MakeGradientHeader(title, subtitle string, height float32) fyne.CanvasObject {
	grad := canvas.NewLinearGradient(ColorPrimary, ColorPrimaryDark, 270)

	titleText := canvas.NewText(title, color.White)
	titleText.TextSize = 22
	titleText.TextStyle = fyne.TextStyle{Bold: true}
	titleText.Alignment = fyne.TextAlignCenter

	var content fyne.CanvasObject
	if subtitle != "" {
		subText := canvas.NewText(subtitle, ColorPrimaryLight)
		subText.TextSize = 13
		subText.Alignment = fyne.TextAlignCenter
		content = container.NewVBox(
			container.NewCenter(titleText),
			container.NewCenter(subText),
		)
	} else {
		content = container.NewCenter(titleText)
	}

	bg := canvas.NewRectangle(color.Transparent)
	bg.SetMinSize(fyne.NewSize(0, height))

	return container.NewStack(grad, bg, container.NewCenter(content))
}

func MakeBadge(text string, bgColor, textColor color.Color) fyne.CanvasObject {
	rect := canvas.NewRectangle(bgColor)
	rect.CornerRadius = 8

	lbl := canvas.NewText(text, textColor)
	lbl.TextSize = 12
	lbl.TextStyle = fyne.TextStyle{Bold: true}
	lbl.Alignment = fyne.TextAlignCenter

	return container.NewStack(rect, container.NewPadded(lbl))
}
