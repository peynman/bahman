package main

import (
	"github.com/rivo/tview"
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"github.com/peyman-abdi/avalanche/app/core/app"
	"github.com/peyman-abdi/avalanche/app/core/database"
	"github.com/peyman-abdi/avalanche/app/core/logger"
	"sort"
	"github.com/golang-collections/collections/stack"
	"github.com/peyman-abdi/avalanche/app/core/config"
	"github.com/peyman-abdi/avalanche/app/core/trans"
	"github.com/peyman-abdi/avalanche/app/core/modules"
	"github.com/peyman-abdi/avalanche/app/core"
)

func main() {
	application := app.Initialize()
	appConfig := config.Initialize(application)
	appLogger := logger.Initialize()
	repo, migrator := database.Initialize(appConfig)
	defer database.Close()
	localizations := trans.Initialize(appConfig, application, appLogger)

	mm := modules.Initialize(appConfig)

	s := core.Initialize(application, appConfig, appLogger, repo, migrator, localizations, mm)

	appLogger.LoadChannels(s)
	mm.LoadModules(s)

	consoleApp = new(ConsoleAppImpl)
	consoleApp.services = s

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
	services interfaces.Services
	app tview.Application
	pages []interfaces.ConsolePage
	mainMenu *tview.List
	priorities map[int]int
	pageStack *stack.Stack
}
var _ interfaces.ConsoleApp = (*ConsoleAppImpl)(nil)
var consoleApp *ConsoleAppImpl

func (c *ConsoleAppImpl) MakeList(title string, items []interfaces.ListItem) interfaces.ConsoleItem {
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
func (c *ConsoleAppImpl) Back()  {
	if c.pageStack.Len() >= 1 {
		pair := c.pageStack.Pop().(StackPagePair)
		c.app.SetRoot(pair.Primitive, pair.FullScreen)
		c.app.SetFocus(pair.Primitive)
	} else {
		c.app.SetRoot(c.mainMenu, true)
		c.app.SetFocus(c.mainMenu)
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
func (c *ConsoleAppImpl) SetPage(item interfaces.ConsoleItem, fullScreen bool) {
	primitive, ok := item.(tview.Primitive)
	if ok {
		c.pageStack.Push(StackPagePair{
			Primitive: primitive,
			FullScreen: fullScreen,
		})
		c.app.SetRoot(primitive, fullScreen)
		c.app.SetFocus(primitive)
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
