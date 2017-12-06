package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/tucnak/telebot"
)

const (
	RouterJoinSep = "_"
)

var (
	ErrRouterIsNil          = fmt.Errorf("<router> is nil")
	ErrRouterBotIsNil       = fmt.Errorf("<router> bot is nil")
	ErrRouterAppExists      = fmt.Errorf("<router> app exists")
	ErrRouterAppIsNil       = fmt.Errorf("<router> app is nil")
	ErrRouterAppNameInvalid = fmt.Errorf("<router> app name invalid")
	ErrRouterAppItemsEmpty  = fmt.Errorf("<router> app items empty")
)

type RouterItem struct {
	Command     string
	Desc        string
	Handler     interface{}
	FullCommand string
}

type RouterItems []*RouterItem

type RouterApp struct {
	Name  string
	Desc  string
	Items RouterItems
}

func (ra *RouterApp) Usage() string {
	if ra == nil || len(ra.Name) == 0 || len(ra.Items) == 0 {
		return ""
	}

	usage := fmt.Sprintf("%s: %s\n", ra.Name, ra.Desc)

	for _, item := range ra.Items {
		usage += item.FullCommand + " " + item.Desc + "\n"
	}

	usage += "\n"

	return usage
}

type Router struct {
	Bot  *telebot.Bot
	Apps map[string]*RouterApp
}

func (rt *Router) joinPath(path ...string) string {
	return strings.Join(path, RouterJoinSep)
}

func (rt *Router) genFullCommand(path ...string) string {
	return "/" + rt.joinPath(path...)
}

func (rt *Router) RegisterApp(app *RouterApp) error {
	if app == nil {
		return ErrRouterAppIsNil
	}

	if len(app.Name) == 0 {
		return ErrRouterAppNameInvalid
	}

	if rt == nil {
		return ErrRouterIsNil
	}

	if rt.Bot == nil {
		return ErrRouterBotIsNil
	}

	if rt.Apps == nil {
		rt.Apps = make(map[string]*RouterApp)
	}

	if _, ok := rt.Apps[app.Name]; ok {
		return ErrRouterAppExists
	}

	if len(app.Items) == 0 {
		return ErrRouterAppItemsEmpty
	}

	log.Println("App:", app.Name, app.Desc)
	log.Println("---")

	for _, item := range app.Items {
		item.FullCommand = rt.genFullCommand(app.Name, item.Command)
		rt.Bot.Handle(item.FullCommand, item.Handler)
		log.Println("item:", item.FullCommand, item.Desc)
	}

	rt.Apps[app.Name] = app

	return nil
}

func (rt *Router) Usage() string {
	var usage string

	if rt == nil {
		return usage
	}

	for _, app := range rt.Apps {
		usage += app.Usage()
	}

	return usage
}
