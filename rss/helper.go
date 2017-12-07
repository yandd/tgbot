package rss

import (
	"time"

	"github.com/mmcdole/gofeed"
)

func getFeedItemUpdateTime(item *gofeed.Item) *time.Time {
	if item.UpdatedParsed != nil {
		return item.UpdatedParsed
	}

	if item.PublishedParsed != nil {
		return item.PublishedParsed
	}

	return nil
}
