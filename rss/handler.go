package rss

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/url"
	"strconv"
	"tgbot/app"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/mmcdole/gofeed"
	"github.com/tucnak/telebot"
)

func List(m *telebot.Message) {
	msg := ""
	defer func() {
		if len(msg) == 0 {
			return
		}
		if m.FromGroup() {
			msg = "@" + m.Sender.Username + "\n" + msg
		}
		app.Bot.Send(m.Chat, msg, &telebot.SendOptions{
			ParseMode:             telebot.ModeHTML,
			DisableWebPagePreview: true,
		})
	}()

	res, err := GetRssResourceByUserID(m.Sender.ID)
	if err != nil {
		msg = "internal error"
		return
	}

	if res == nil || len(*res) == 0 {
		msg = "rss list is empty"
		return
	}

	for _, rss := range *res {
		msg += fmt.Sprintf("%d: <a href=\"%s\">%s</a> &lt;%s&gt;\n", rss.ID, rss.Link, html.EscapeString(rss.Title), rss.URL)
	}
	return
}

func Add(m *telebot.Message) {
	msg := ""
	defer func() {
		if len(msg) == 0 {
			return
		}
		if m.FromGroup() {
			msg = "@" + m.Sender.Username + "\n" + msg
		}
		app.Bot.Send(m.Chat, msg, &telebot.SendOptions{
			ParseMode:             telebot.ModeHTML,
			DisableWebPagePreview: true,
		})
	}()

	_, err := url.Parse(m.Payload)
	if err != nil {
		log.Println("Error: url is invalid,", m.Payload, m.Sender, err)
		msg = "rss url is invalid."
		return
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(m.Payload)
	if err != nil {
		log.Println("Error: ParseURL failed,", m.Payload, m.Sender, err)
		msg = "rss url parse failed."
		return
	}

	r := RssResource{
		URL:   m.Payload,
		Title: feed.Title,
		Link:  feed.Link,
	}

	if len(feed.Items) > 0 {
		r.LastItemTitle = feed.Items[0].Title
		r.LastItemLink = feed.Items[0].Link
		r.LastItemPublishTime = feed.Items[0].PublishedParsed.Unix()
	}

	rssID, err := AddRssResource(&r)
	if err != nil {
		msg = "internal error"
		return
	}

	err = AddRssUser(rssID, m.Sender.ID)
	if err != nil {
		msg = "internal error"
		return
	}

	if _, ok := rssMgr.Load(r.URL); ok {
		msg = fmt.Sprintf("%s: <a href=\"%s\">%s</a>\n\n---\n%s: <a href=\"%s\">@%s</a>", time.Unix(r.LastItemPublishTime, 0).Format("2006/01/02"), r.LastItemLink, html.EscapeString(r.LastItemTitle), "rss", r.Link, html.EscapeString(r.Title))
		return
	}

	res, err := GetRssResourceByID(rssID)
	if err != nil {
		msg = "internal error"
		return
	}

	mgr := &RssMgr{
		Rss: *res,
	}

	s := gocron.NewScheduler()
	s.Every(app.Cfg.FetchInterval).Seconds().Do(func() {
		fetchAndSend(mgr)
	})
	mgr.Stop = s.Start()

	rssMgr.Store(r.URL, mgr)

	msg = fmt.Sprintf("%s: <a href=\"%s\">%s</a>\n\n---\n%s: <a href=\"%s\">@%s</a>", time.Unix(r.LastItemPublishTime, 0).Format("2006/01/02"), r.LastItemLink, html.EscapeString(r.LastItemTitle), "rss", r.Link, html.EscapeString(r.Title))
	return
}

func Del(m *telebot.Message) {
	msg := ""
	defer func() {
		if len(msg) == 0 {
			return
		}
		if m.FromGroup() {
			msg = "@" + m.Sender.Username + "\n" + msg
		}
		app.Bot.Send(m.Chat, msg, &telebot.SendOptions{
			ParseMode:             telebot.ModeHTML,
			DisableWebPagePreview: true,
		})
	}()

	rssID, err := strconv.ParseInt(m.Payload, 10, 64)
	if err != nil {
		_, err = url.Parse(m.Payload)
		if err != nil {
			msg = "rss id/url is invalid"
			return
		}
		res, err := GetRssResourceByURL(m.Payload)
		if err != nil {
			if err == sql.ErrNoRows {
				msg = "rss url is not in your list"
				return
			} else {
				msg = "internal error"
				return
			}
		}
		rssID = res.ID
	}

	rowsAffected, err := UpdateRssUser(rssID, m.Sender.ID, map[string]interface{}{
		"is_deleted": 1,
	})
	if err != nil {
		msg = "internal error"
		return
	}

	if rowsAffected == 0 {
		msg = "succ"
		return
	}

	msg = "success"
	return
}
