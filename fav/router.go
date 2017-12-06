package fav

import "tgbot/app"

func init() {
	app.Route.RegisterApp(&app.RouterApp{
		Name: "fav",
		Desc: "收藏",
		Items: app.RouterItems{
			&app.RouterItem{
				Command: "list",
				Desc:    "显示收藏列表",
				Handler: List,
			},
			&app.RouterItem{
				Command: "add",
				Desc:    "<text>: 添加收藏",
				Handler: Add,
			},
		},
	})
}
