package main

import (
	"6-fyne-chat/view"
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
)

func main() {
	os.Mkdir("log", 0777)
	file, err := os.OpenFile("log/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	log.SetOutput(file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("Chat")
	w.Resize(fyne.NewSize(900, 500))

	container := view.Chat(w)
	w.SetContent(container)
	w.ShowAndRun()

}
