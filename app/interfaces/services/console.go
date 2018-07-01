package services

type ConsoleItem interface{}

type ListItem struct {
	Title       string
	Description string
	Shortcut    rune
	Callback    func()
}
type ModalWindow struct {
	Content  string
	Buttons  []string
	Callback func(int)
}

type ConsoleApp interface {
	Back()
	BackToMainMenu()
	Quit()
	SetPage(name string, item ConsoleItem, fullScreen bool)

	Ask(question string, callback func(yes bool))
	Info(content string, callback func())
	Warn(content string, callback func())
	Error(content string, callback func())
	Success(content string, callback func())

	MakeList(title string, items []ListItem) ConsoleItem
	MakeModal(window ModalWindow) ConsoleItem
}

type ConsolePage interface {
	Title() string
	Description() string
	Priority() int
	OnSelected(console ConsoleApp)
}
