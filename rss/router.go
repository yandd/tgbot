package rss

import "tgbot/app"

func init() {
	app.Route.RegisterApp(&app.RouterApp{
		Name: "rss",
		Desc: "rss订阅",
		Items: app.RouterItems{
			&app.RouterItem{
				Command: "list",
				Desc:    "显示订阅列表",
				Handler: List,
			},
			&app.RouterItem{
				Command: "add",
				Desc:    "<url>: 添加订阅url",
				Handler: Add,
			},
			&app.RouterItem{
				Command: "del",
				Desc:    "<id/url>: 删除订阅id/url",
				Handler: Del,
			},
		},
	})
}
