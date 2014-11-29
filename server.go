package main

import (
	"encoding/json"
	"fmt"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mvader/gorss"
)

type News struct {
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Author    string   `json:"author"`
	URL       string   `json:"url"`
	Tags      []string `json:"categories,omitempty"`
	Media     []Media  `json:"media,omitempty"`
	Thumbnail string   `json:"thumbnail,omitempty"`
}

type Media struct {
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

func DecodeFromJSON(content string, value interface{}) error {
	return json.Unmarshal([]byte(content), value)
}

var basePath = "http://www.elmundo.es/rss/hackathon/"

func rssUrl(category string) string {
	return basePath + category + ".xml"
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

func NewsFromURL(url string) []News {
	var news []News
	feed, err := rss.LoadFeed(url)

	if len(feed.Channels) < 1 {
		err = errors.New("no channels in feed")
	}

	if err == nil {
		news = make([]News, len(feed.Channels[0].Items))
		for i, v := range feed.Channels[0].Items {
			news[i] = newsFromRSSItem(v)
		}
	}

	return news
}

func serveJson(c *gin.Context) {
	var category = c.Params.ByName("category")
	news := NewsFromURL(rssUrl(category))
	if len(news) > 0 {
		c.JSON(200, map[string]interface{}{
			"news": news,
		})
		return
	}

	c.JSON(404, map[string]interface{}{
		"message": "category not found",
	})
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			fmt.Println("options")
			c.Abort(200)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/categories/:category", serveJson)
	r.Static("/file/", basePath)

	r.Run(":3001")
}
