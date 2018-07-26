package console


import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/rivo/tview"
)

func (c *consoleAppImpl) MakeList(title string, items []services.ListItem) services.ConsoleItem {
	list := tview.NewList()
	for _, item := range items {
		list.AddItem(item.Title, item.Description, item.Shortcut, item.Callback)
	}
	list.AddItem("Back", "Return to previous page", 'b', func() {
		c.Back()
	})
	list.AddItem("Main Menu", "Return to main menu", 'r', func() {
		c.BackToMainMenu()
	})
	list.SetBorder(true)
	list.SetBorderPadding(1, 1, 1, 1)
	list.SetTitle(title)
	return list
}
func (c *consoleAppImpl) MakeModal(window services.ModalWindow) services.ConsoleItem {
	modal := tview.NewModal()
	modal.SetBorder(true)
	modal.SetText(window.Content)
	modal.SetBorderPadding(5, 5, 5, 5)
	modal.AddButtons(window.Buttons)
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		c.Back()
		window.Callback(buttonIndex)
	})
	return modal
}
