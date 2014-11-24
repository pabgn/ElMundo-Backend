package models

import (
	"errors"
	"github.com/mvader/gorss"
	"time"
)

type News struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Author    string   `json:"author"`
	URL       string   `json:"url"`
	Tags      []string `json:"categories,omitempty"`
	Media     []Media  `json:"media,omitempty"`
	Thumbnail string   `json:"thumbnail,omitempty"`
	Tweets    []Tweet  `json:"tweets,omitempty"`
}

type Media struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

func newsFromRSSItem(item rss.Item) News {
	var (
		thumbnail string
		media     = make([]Media, len(item.MediaContent))
		tags      = make([]string, len(item.Categories))
	)

	if item.Author == "" {
		item.Author = item.Creator
	}

	for i, v := range item.MediaContent {
		media[i] = Media{
			URL:   v.URL,
			Title: v.Title.Value,
		}
	}

	for i, v := range item.Categories {
		tags[i] = v.Value
	}

	return News{
		Title:     item.Title,
		Content:   item.Description,
		Author:    item.Author,
		URL:       item.Guid,
		Tags:      tags,
		Media:     media,
		Thumbnail: thumbnail,
	}
}

func NewsFromURL(url string) ([]News, uint64, error) {
	var (
		news        []News
		lastRefresh uint64
	)
	feed, err := rss.LoadFeed(url)

	if len(feed.Channels) < 1 {
		err = errors.New("no channels in feed")
	}

	if err == nil {
		news = make([]News, len(feed.Channels[0].Items))
		for i, v := range feed.Channels[0].Items {
			news[i] = newsFromRSSItem(v)
		}

		t, _ := time.Parse(time.RFC1123Z, feed.Channels[0].LastBuildDate)
		lastRefresh = uint64(t.Unix())
	}

	return news, lastRefresh, err
}
