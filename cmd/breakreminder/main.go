package main

import (
	fyneApp "fyne.io/fyne/v2/app"

	"github.com/zhangxinyu/breakreminder/assets"
	"github.com/zhangxinyu/breakreminder/internal/app"
)

func main() {
	a := fyneApp.NewWithID("com.zhangxinyu.breakreminder")
	a.SetIcon(assets.AppIcon)

	application := app.New(a)
	application.Run()
}
