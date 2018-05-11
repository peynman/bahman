package interfaces

import "github.com/rivo/tview"

type ConsoleApp interface {
	Back()
	BackToMainMenu()
	Quit()
	SetPage(primitive tview.Primitive, fullScreen bool)
}

type ConsolePage interface {
	Title() string
	Description() string
	Priority() int
	OnSelected(console ConsoleApp)
}