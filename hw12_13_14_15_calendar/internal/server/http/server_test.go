package internalhttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
	basestorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/base"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func test_start() {
	conf := config.NewConfig()
	toml.DecodeFile("../../../configs/config.toml", &conf)

	logger, _ := zap.NewDevelopment()
	logg := logger.Sugar()

	storage, _ := basestorage.InitStorage(conf.Storage, logg)
	calendarApp = app.New(logg, storage)
}

func TestHandler(t *testing.T) {
	test_start()
	rPath := "/event/"
	router := gin.Default()
	router.GET(rPath, getAllEventsHandler)
	req, _ := http.NewRequest("GET", rPath, strings.NewReader(`{"id": "1","name": "joe"}`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	fmt.Printf("status: %d", w.Code)
	fmt.Printf("response: %s", w.Body.String())
	require.Equal(t, w.Code, 200)
	require.Equal(t, w.Code, 500)
}

/*

func TestEventHandler(t *testing.T) {

	handler := internalhttp.CreateHandler(app, "logPath", logger)

	r := httptest.NewRequest("GET", "http://127.0.0.1:80/user?id=42", nil)
	w := httptest.NewRecorder()


	/*
		tt := []struct {
			name       string
			method     string
			input      *[]model.Event
			want       string
			statusCode int
		}{
			{
				name:       "empty store",
				method:     http.MethodGet,
				input:      &[]model.Event{},
				want:       "{}",
				statusCode: http.StatusOK,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				request := httptest.NewRequest(tc.method, "/event/", nil)
				responseRecorder := httptest.NewRecorder()

			})
		}

}
*/
