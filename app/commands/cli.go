package main

import (
	"github.com/rivo/tview"
	"avalanche/app/core/interfaces"
	"avalanche/app/core/app"
	"avalanche/app/core/database"
	"avalanche/app/core/logger"
	"sort"
	"github.com/golang-collections/collections/stack"
	"avalanche/app/core/config"
	"avalanche/app/core/trans"
)

func main() {
	config.Initialize()
	logger.Initialize()
	database.Initialize()
	defer database.Close()
	trans.Initialize()

	consoleApp = new(ConsoleAppImpl)
	consoleApp.SetupWithConfPath("console")
	
	if err := consoleApp.Run(); err != nil {
		panic(err)
	}
}

type StackPagePair struct {
	Primitive tview.Primitive
	FullScreen bool
}
type ConsoleAppImpl struct {
	app tview.Application
	pages []interfaces.ConsolePage
	mainMenu *tview.List
	priorities map[int]int
	pageStack *stack.Stack
}
var _ interfaces.ConsoleApp = (*ConsoleAppImpl)(nil)
var consoleApp *ConsoleAppImpl

func (c *ConsoleAppImpl) Back()  {
	if c.pageStack.Len() >= 1 {
		pair := c.pageStack.Pop().(StackPagePair)
		c.app.SetRoot(pair.Primitive, pair.FullScreen)
	} else {
		c.app.SetRoot(c.mainMenu, true)
	}
}
func (c *ConsoleAppImpl) BackToMainMenu()  {
	for c.pageStack.Len() > 0 {
		c.pageStack.Pop()
	}

	c.app.SetRoot(c.mainMenu, true)
	c.app.SetFocus(c.mainMenu)
}
func (c *ConsoleAppImpl) Quit()  {
	c.app.Stop()
}
func (c *ConsoleAppImpl) SetPage(primitive tview.Primitive, fullScreen bool) {
	c.pageStack.Push(StackPagePair{
		Primitive: primitive,
		FullScreen: fullScreen,
	})
	c.app.SetRoot(primitive, fullScreen)
	c.app.SetFocus(primitive)
}
func (c *ConsoleAppImpl) SetupWithConfPath(confPath string) {
	c.pageStack = stack.New()

	c.mainMenu = tview.NewList()
	c.mainMenu.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle("Avalanche Web Serve Cli")

	c.priorities = make(map[int]int)

	modules := app.InitAvalanchePlugins(app.ModulesPath("console"))
	for index, module := range modules {
		page := module.Interface().(interfaces.ConsolePage)
		c.pages = append(c.pages, page)
		c.priorities[page.Priority()] = index
	}

	var priorities []int
	for priority := range c.priorities {
		priorities = append(priorities, priority)
	}
	sort.Ints(priorities)

	shortcuts := []rune {
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

	c.app.
		SetRoot(c.mainMenu, true).
		SetFocus(c.mainMenu)
}
func (c *ConsoleAppImpl) Run() error  {
	return c.app.Run()
}
