package internalhttp

import (
	"context"
	"net/http"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/gin-gonic/gin"
)

var (
	statusOK          = http.StatusOK
	statusServerError = http.StatusInternalServerError
)

var calendarApp *app.App

type Server struct { // TODO
	server *http.Server
}

type Logger interface { // TODO
}

func NewServer(logger Logger, app *app.App) *Server {
	return &Server{server: &http.Server{
		Addr:    "",
		Handler: createHandler(app),
	}}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

func createHandler(calendar *app.App) http.Handler {
	router := gin.Default()
	calendarApp = calendar

	router.POST("/event/", createEventHandler) // Create
	//router.GET("/event/:id", getEventHandler)
	router.PUT("/event/:id", updateEventHandler)        // Update
	router.DELETE("/event/:id", deleteEventHandler)     // Delete
	router.GET("/due/:year/:month/:day", dueDayHandler) // ListEventsByDay
	router.GET("/due/:year/:month", dueMonthHandler)    // ListEventsByMonth
	router.GET("/event/", getAllEventsHandler)          // ListAllEvents
	router.GET("/user/:uid", getAllUserEventsHandler)   // ListUserEvents
	return router
}

func createEventHandler(c *gin.Context) {
	// TODO
}

func updateEventHandler(c *gin.Context) {
	// TODO
}

func deleteEventHandler(c *gin.Context) {
	// TODO
}

func dueDayHandler(c *gin.Context) {
	// TODO
}

func dueMonthHandler(c *gin.Context) {
	// TODO
}

func getAllEventsHandler(c *gin.Context) {
	allEvents, err := calendarApp.ListAllEvents()
	if err != nil {
		//c.String(http.StatusBadRequest, err.Error())
		c.Writer.WriteHeader(statusServerError)
		return
	}
	c.JSON(statusOK, allEvents)
}

func getAllUserEventsHandler(c *gin.Context) {
	// TODO
}
