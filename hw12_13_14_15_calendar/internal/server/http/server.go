package internalhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	jsontime "github.com/liamylian/jsontime/v2/v2"

	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	jsontime.AddTimeFormatAlias("sql_datetime", "2006-01-02 15:04:05")
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
	var json = jsontime.ConfigWithCustomTimeFormat

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	var e model.Event
	err = json.Unmarshal(jsonData, &e)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id, err := calendarApp.Create(e)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"Id": id})
}

func updateEventHandler(c *gin.Context) {
	fmt.Println("[server] updateEventHandler")
	var json = jsontime.ConfigWithCustomTimeFormat

	id, err := uuid.Parse(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	var e model.Event
	err = json.Unmarshal(jsonData, &e)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = calendarApp.Update(id, e)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"Id": id})
}

func deleteEventHandler(c *gin.Context) {
	fmt.Println("[server] deleteEventHandler")

	id, err := uuid.Parse(c.Params.ByName("id"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	calendarApp.Delete(id)

	c.JSON(http.StatusOK, "deleted")
}

func dueDayHandler(c *gin.Context) {
	fmt.Println("[server] dueDayHandler")
	params := []string{"day", "month", "year"}

	m := make(map[string]int)
	for _, d := range params {
		val, err := strconv.Atoi(c.Params.ByName(d))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		m[d] = val
	}
	date := time.Date(m["year"], time.Month(m["month"]), m["day"], 1, 2, 3, 0, time.Local)

	events, err := calendarApp.ListEventsByDay(date)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, events)
}

func dueMonthHandler(c *gin.Context) {
	fmt.Println("[server] dueMonthHandler")
	params := []string{"month", "year"}

	m := make(map[string]int)
	for _, d := range params {
		val, err := strconv.Atoi(c.Params.ByName(d))
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		m[d] = val
	}
	date := time.Date(m["year"], time.Month(m["month"]), m["day"], 1, 2, 3, 0, time.Local)

	events, err := calendarApp.ListEventsByDay(date)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, events)
}

func getAllEventsHandler(c *gin.Context) {
	fmt.Println("[server] getAllEventsHandler")

	allEvents, err := calendarApp.ListAllEvents()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		//c.Writer.WriteHeader(statusServerError)
		return
	}
	c.JSON(http.StatusOK, allEvents)
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
