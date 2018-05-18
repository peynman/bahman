package interfaces

type ConsoleItem interface {}

type ListItem struct {
	Title string
	Description string
	Shortcut rune
	Callback func()
}

type ConsoleApp interface {
	Back()
	BackToMainMenu()
	Quit()
	SetPage(item ConsoleItem, fullScreen bool)

	MakeList(title string, items []ListItem) ConsoleItem
}

type ConsolePage interface {
	Title() string
	Description() string
	Priority() int
	OnSelected(console ConsoleApp)
}