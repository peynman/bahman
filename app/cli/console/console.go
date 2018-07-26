package console

import (
	"github.com/peyman-abdi/bahman/app/interfaces/services"
	"github.com/rivo/tview"
	"sort"
	"github.com/golang-collections/collections/stack"
)

type StackPagePair struct {
	Primitive  tview.Primitive
	FullScreen bool
	Name       string
}
type consoleAppImpl struct {
	services   services.Services
	app        tview.Application
	pages      []services.ConsolePage
	mainMenu   *tview.List
	priorities map[int]int
	pageStack  *stack.Stack
	pageView   *tview.Pages
}

var _ services.ConsoleApp = (*consoleAppImpl)(nil)

func New(confPath string, services services.Services) services.ConsoleApp {
	console := new(consoleAppImpl)
	console.services = services
	console.SetupWithConfPath(confPath)
	return console
}

func (c *consoleAppImpl) SetupWithConfPath(confPath string) {
	c.pageStack = stack.New()

	c.mainMenu = tview.NewList()
	c.mainMenu.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("bahman Web Serve Cli")

	c.priorities = make(map[int]int)

	modules := c.services.App().InitbahmanPlugins(c.services.App().ModulesPath("console"), c.services)
	for index, module := range modules {
		page := module.Interface().(services.ConsolePage)
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

	c.mainMenu.AddItem("Quite", "Exist bahman cli", 'q', func() {
		c.app.Stop()
	})

	c.pageView = tview.NewPages()
	c.pageView.AddPage("main_menu", c.mainMenu, true, true)

	c.app.
		SetRoot(c.pageView, true).
		SetFocus(c.pageView)
}
func (c *consoleAppImpl) Run() error {
	return c.app.Run()
}


func (c *consoleAppImpl) Back() {
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
func (c *consoleAppImpl) BackToMainMenu() {
	for c.pageStack.Len() > 0 {
		page := c.pageStack.Pop().(StackPagePair)
		c.pageView.RemovePage(page.Name)
	}
	c.pageView.SwitchToPage("main_menu")
	c.app.SetFocus(c.pageView)
}
func (c *consoleAppImpl) Quit() {
	c.app.Stop()
}
func (c *consoleAppImpl) SetPage(name string, item services.ConsoleItem, fullScreen bool) {
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
