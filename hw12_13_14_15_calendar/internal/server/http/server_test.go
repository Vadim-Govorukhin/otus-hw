package internalhttp_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
	internalhttp "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/server/http"
	basestorage "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage/base"
	ginzap "github.com/akath19/gin-zap"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var calendarApp *app.App

func testStart() (*zap.SugaredLogger, http.Handler) {
	conf := config.NewConfig()
	toml.DecodeFile("../../../configs/config.toml", &conf)

	logger, _ := zap.NewDevelopment()
	logg := logger.Sugar()

	conf.Storage.Type = "memory" // tests based on memory storage
	storage, _ := basestorage.InitStorage(conf.Storage, logg)
	calendarApp = app.New(logg, storage)

	gin.SetMode(gin.ReleaseMode)
	router := internalhttp.CreateHandler(calendarApp, ginzap.Logger(3*time.Second, logg.Desugar()))
	return logg, router
}

func TestHandler(t *testing.T) {
	logg, router := testStart()

	testCases := []struct {
		name        string
		method      string
		url         string
		requestBody io.Reader
		wantBody    string
		statusCode  int
	}{
		{
			name:        "get empty events",
			method:      http.MethodGet,
			url:         "/event/",
			requestBody: nil,
			wantBody:    "[]",
			statusCode:  http.StatusOK,
		},
		{
			name:        "create events and get all of them",
			method:      http.MethodGet,
			url:         "/event/",
			requestBody: nil,
			wantBody:    "[]",
			statusCode:  http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logg.Infof("============== start test %s ==========", tc.name)
			req := httptest.NewRequest(tc.method, tc.url, tc.requestBody)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			logg.Infof("status: %d", w.Code)
			logg.Infof("response: %s", w.Body.String())
			require.Equal(t, tc.statusCode, w.Code)
			require.Equal(t, tc.wantBody, w.Body.String())
		})
	}
}
