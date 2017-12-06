package rss

import (
	"fmt"
	"html"
	"log"
	"sync"
	"tgbot/app"

	"github.com/jasonlvhit/gocron"
	"github.com/mmcdole/gofeed"
	"github.com/tucnak/telebot"
)

type RssMgr struct {
	Stop chan<- bool
	Rss  RssResource
}

var rssMgr = &sync.Map{}

func init() {
	err := createTables()
	if err != nil {
		log.Fatalln("Error: createTables failed,", err)
		return
	}

	err = doCron()
	if err != nil {
		log.Fatalln("Error: doCron failed,", err)
		return
	}
}

func doCron() error {
	res, err := GetRssResources()
	if err != nil {
		return err
	}

	if res == nil || len(*res) == 0 {
		return nil
	}

	for _, r := range *res {
		userIDs, err := GetUserIDsByRssID(r.ID)
		if err != nil {
			return err
		}
		if len(userIDs) == 0 {
			continue
		}

		mgr := &RssMgr{
			Rss: r,
		}

		s := gocron.NewScheduler()
		s.Every(app.Cfg.FetchInterval).Seconds().Do(func() {
			fetchAndSend(mgr)
		})
		mgr.Stop = s.Start()

		rssMgr.Store(r.URL, mgr)
	}

	return nil
}

func fetchAndSend(m *RssMgr) error {
	log.Println("Info: fetch start RSS", m.Rss.URL)

	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(m.Rss.URL)
	if err != nil {
		log.Println("Error: fetch failed RSS", m.Rss.URL)
		return nil
	}

	if feed == nil || len(feed.Items) == 0 {
		log.Println("Error: fetch items empty RSS", m.Rss.URL)
		return nil
	}

	userIDs, err := GetUserIDsByRssID(m.Rss.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	if len(userIDs) == 0 {
		m.Stop <- true
		rssMgr.Delete(m.Rss.URL)
		log.Println("Info: fetch cancel RSS", m.Rss.URL)
		return nil
	}

	if (feed.Items[0].PublishedParsed.Unix() == m.Rss.LastItemPublishTime && feed.Items[0].Link == m.Rss.LastItemLink && feed.Items[0].Title == m.Rss.LastItemTitle) || feed.Items[0].PublishedParsed.Unix() < m.Rss.LastItemPublishTime {
		log.Println("Info: fetch old RSS", m.Rss.URL, feed.Items[0].Title, feed.Items[0].Link, feed.Items[0].PublishedParsed.Format("2006/01/02 15:04:05"))
		return nil
	}

	msg := ""
	for _, item := range feed.Items {
		if (item.PublishedParsed.Unix() == m.Rss.LastItemPublishTime && item.Link == m.Rss.LastItemLink && item.Title == m.Rss.LastItemTitle) || item.PublishedParsed.Unix() < m.Rss.LastItemPublishTime {
			break
		}
		msg += fmt.Sprintf("%s: <a href=\"%s\">%s</a>\n", item.PublishedParsed.Format("2006/01/02"), item.Link, html.EscapeString(item.Title))
	}
	msg = fmt.Sprintf("\n---\n%s: <a href=\"%s\">@%s</a>", "rss", m.Rss.Link, html.EscapeString(m.Rss.Title))

	for _, userID := range userIDs {
		_, err = app.Bot.Send(&telebot.User{
			ID: userID,
		}, msg, &telebot.SendOptions{
			ParseMode:             telebot.ModeHTML,
			DisableWebPagePreview: true,
		})
		if err != nil {
			log.Println("Error: send failed RSS", m.Rss.URL, userID)
		}
	}

	m.Rss.LastItemTitle = feed.Items[0].Title
	m.Rss.LastItemLink = feed.Items[0].Link
	m.Rss.LastItemPublishTime = feed.Items[0].PublishedParsed.Unix()

	err = UpdateRssResource(m.Rss.ID, map[string]interface{}{
		"last_item_title":        m.Rss.LastItemTitle,
		"last_item_link":         m.Rss.LastItemLink,
		"last_item_publish_time": m.Rss.LastItemPublishTime,
	})
	if err != nil {
		log.Println(err)
	}

	log.Println("Info: fetch update RSS", m.Rss.URL)
	return nil
}
