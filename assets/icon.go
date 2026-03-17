package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed icon.png
var iconBytes []byte

var AppIcon = &fyne.StaticResource{
	StaticName:    "icon.png",
	StaticContent: iconBytes,
}
