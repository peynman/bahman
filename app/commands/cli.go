package main

import (
	"github.com/golang-collections/collections/stack"
	"github.com/peyman-abdi/avalanche/app/interfaces/core"
	"github.com/peyman-abdi/avalanche/app/modules/kernel"
	"github.com/rivo/tview"
	"sort"
)

func main() {
	services := kernel.SetupCLIKernel()

	console = new(ConsoleAppImpl)
	console.services = services

	console.SetupWithConfPath("console")

	if err := console.Run(); err != nil {
		panic(err)
	}
}

type StackPagePair struct {
	Primitive  tview.Primitive
	FullScreen bool
	Name       string
}
type ConsoleAppImpl struct {
	services   core.Services
	app        tview.Application
	pages      []core.ConsolePage
	mainMenu   *tview.List
	priorities map[int]int
	pageStack  *stack.Stack
	pageView   *tview.Pages
}

var _ core.ConsoleApp = (*ConsoleAppImpl)(nil)
var console *ConsoleAppImpl

func (c *ConsoleAppImpl) MakeList(title string, items []core.ListItem) core.ConsoleItem {
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
func (c *ConsoleAppImpl) MakeModal(window core.ModalWindow) core.ConsoleItem {
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
func (c *ConsoleAppImpl) Ask(question string, callback func(bool)) {
	modal := tview.NewModal()
	modal.SetBorder(true)
	modal.SetText(question)
	modal.SetBorderPadding(5, 5, 5, 5)
	modal.AddButtons([]string{"Yes", "No"})
	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		callback(buttonIndex == 0)
	})
	c.SetPage("ask", modal, false)
}

func (c *ConsoleAppImpl) Back() {
	if c.pageStack.Len() >= 2 {
		last := c.pageStack.Pop().(StackPagePair) // current page
		c.pageView.RemovePage(last.Name)
		pair := c.pageStack.Pop().(StackPagePair) // previous page
		c.pageView.SwitchToPage(pair.Name)
		c.pageStack.Push(pair) // push previous as current page again
	} else {
		c.pageView.SwitchToPage("main_menu")
	}

	c.app.SetFocus(c.pageView)
}
func (c *ConsoleAppImpl) BackToMainMenu() {
	for c.pageStack.Len() > 0 {
		page := c.pageStack.Pop().(StackPagePair)
		c.pageView.RemovePage(page.Name)
	}
	c.pageView.SwitchToPage("main_menu")
	c.app.SetFocus(c.pageView)
}
func (c *ConsoleAppImpl) Quit() {
	c.app.Stop()
}
func (c *ConsoleAppImpl) SetPage(name string, item core.ConsoleItem, fullScreen bool) {
	primitive, ok := item.(tview.Primitive)
	if ok {
		c.pageStack.Push(StackPagePair{
			Primitive:  primitive,
			FullScreen: fullScreen,
			Name:       name,
		})

		if fullScreen {
			c.pageView.AddAndSwitchToPage(name, primitive, true)
		} else {
			c.pageView.AddPage(name, primitive, true, true)
		}

		c.app.SetFocus(c.pageView)
	}
}

func (c *ConsoleAppImpl) SetupWithConfPath(confPath string) {
	c.pageStack = stack.New()

	c.mainMenu = tview.NewList()
	c.mainMenu.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("Avalanche Web Serve Cli")

	c.priorities = make(map[int]int)

	modules := c.services.App().InitAvalanchePlugins(c.services.App().ModulesPath("console"), c.services)
	for index, module := range modules {
		page := module.Interface().(core.ConsolePage)
		c.pages = append(c.pages, page)
		c.priorities[page.Priority()] = index
	}

	var priorities []int
	for priority := range c.priorities {
		priorities = append(priorities, priority)
	}
	sort.Ints(priorities)

	shortcuts := []rune{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	}
	indexer := 0
	for index := range priorities {
		page := c.pages[index]
		c.mainMenu.AddItem(page.Title(), page.Description(), shortcuts[indexer], func() {
			page.OnSelected(c)
		})
		indexer++
	}

	c.mainMenu.AddItem("Quite", "Exist Avalanche cli", 'q', func() {
		c.app.Stop()
	})

	c.pageView = tview.NewPages()
	c.pageView.AddPage("main_menu", c.mainMenu, true, true)

	c.app.
		SetRoot(c.pageView, true).
		SetFocus(c.pageView)
}
func (c *ConsoleAppImpl) Run() error {
	return c.app.Run()
}
