package app

import (
	"github.com/gin-gonic/gin"
	"github.com/pabgn/ElMundo-Backend/handlers"
	"strconv"
)

type App struct {
	engine *gin.Engine
}

// Run starts the app at the given port
func (app *App) Run(port int) {
	app.engine.Run(":" + strconv.Itoa(port))
}

// add the routes to the router engine
func (app *App) addRoutes() {
	app.engine.GET("/channel/:channel", handlers.GetChannel)
	app.engine.GET("/tag/:tag", handlers.GetTag)
	app.engine.GET("/tweets/:keywords", handlers.GetTweets)
	app.engine.GET("/ads", handlers.GetAds)
	app.engine.POST("/ads/seen", handlers.AdSeen)
}

// NewApp creates a new instance of the app
func NewApp() *App {
	app := &App{}
	app.engine = gin.Default()
	app.addRoutes()

	// TODO: Set default channel urls to storage

	return app
}
