package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func NewServer(logger Logger, app *app.App, addres string) *Server {
	return &Server{server: &http.Server{
		Addr:    addres,
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
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
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
	fmt.Println("[server] createEventHandler")

	type RequestEvent struct {
		ID             model.EventID `json:"event_id,omitempty"`
		Title          string        `json:"title"`
		StartDate      time.Time     `json:"start_date"`
		EndDate        time.Time     `json:"end_date"`
		Description    string        `json:"descr,omitempty"`
		UserID         model.UserID  `json:"user_id"`
		NotifyUserTime float64       `json:"notify_user_time,omitempty"`
	}

	var rev RequestEvent
	if err := c.ShouldBindJSON(&rev); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	rev.ID = uuid.New()
	id := calendarApp.Create(model.Event{})
	c.JSON(http.StatusOK, gin.H{"Id": id})
}

func updateEventHandler(c *gin.Context) {
	fmt.Println("[server] updateEventHandler")
	// TODO
}

func deleteEventHandler(c *gin.Context) {
	fmt.Println("[server] deleteEventHandler")
	// TODO
}

func dueDayHandler(c *gin.Context) {
	fmt.Println("[server] dueDayHandler")
	// TODO
}

func dueMonthHandler(c *gin.Context) {
	fmt.Println("[server] dueMonthHandler")
	// TODO
}

func getAllEventsHandler(c *gin.Context) {
	fmt.Println("[server] getAllEventsHandler")

	allEvents, err := calendarApp.ListAllEvents()
	if err != nil {
		//c.String(http.StatusBadRequest, err.Error())
		c.Writer.WriteHeader(statusServerError)
		return
	}
	c.JSON(statusOK, allEvents)
}

func getAllUserEventsHandler(c *gin.Context) {
	fmt.Println("[server] getAllUserEventsHandler")

	id, err := strconv.Atoi(c.Params.ByName("uid"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	events, err := calendarApp.ListUserEvents(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, events)
}
