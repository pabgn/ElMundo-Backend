package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/pabgn/ElMundo-Backend/models"
	"github.com/pabgn/ElMundo-Backend/services"
	"time"
)

var TIMEOUT_MAX_AGE = 3600

func GetChannel(c *gin.Context) {
	var (
		err  error
		news []models.News
	)

	channel := c.Params.ByName("channel")
	storage := c.MustGet("storage").(services.Storage)

	lastRefresh := storage.GetLastChannelRefresh(channel)
	if int64(lastRefresh)+int64(TIMEOUT_MAX_AGE) > time.Now().Unix() {
		var newsContent string

		newsContent, err = storage.GetNewsByChannel(channel)
		err = services.DecodeFromJSON(newsContent, &news)
	} else {
		// Erase the channel info (tweets, news and refreshes)
		// because it has timed out or it does not exist
		storage.EraseChannel(channel)

		var url string
		url, err := storage.GetChannelURL(channel)
		if err == nil {
			news, lastRefresh, err = models.NewsFromURL(url)

			// Store the retrieved result
			go storeCachedResults(storage, lastRefresh, channel, news)
		}
	}

	if err != nil {
		c.JSON(404, map[string]interface{}{
			"message": err.Error(),
		})
	}

	c.JSON(200, map[string]interface{}{
		"news":         news,
		"next_refresh": lastRefresh,
	})
}

func storeCachedResults(storage services.Storage, lastRefresh uint64, channel string, news []models.News) {
	storage.SetLastChannelRefresh(channel, lastRefresh)
	var encodedNews string
	encodedNews, _ = services.EncodeJSON(news)
	storage.SetNewsByChannel(channel, encodedNews)

	// Store tags
}
