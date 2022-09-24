package internalhttp_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/config"
	internalhttp "github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/Vadim-Govorukhin/otus-hw/hw12_13_14_15_calendar/internal/storage"
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
		method      []string
		url         []string
		requestBody []io.Reader
		wantBody    []string
		statusCode  []int
	}{
		{
			name:        "get empty events",
			method:      []string{http.MethodGet},
			url:         []string{"/event/"},
			requestBody: []io.Reader{nil},
			wantBody:    []string{"[]"},
			statusCode:  []int{http.StatusOK},
		},
		{
			name:   "create events",
			method: []string{http.MethodPost, http.MethodPost, http.MethodPost},
			url:    []string{"/event/", "/event/", "/event/"},
			requestBody: []io.Reader{bytes.NewReader(storage.TestEventJson),
				bytes.NewReader(storage.TestEvent2Json),
				bytes.NewReader(storage.TestEvent3Json)},
			wantBody: []string{string(storage.TestEventIDJson),
				string(storage.TestEvent2IDJson),
				string(storage.TestEvent3IDJson)},
			statusCode: []int{http.StatusOK, http.StatusOK, http.StatusOK},
		},
		{
			name:        "get list all events",
			method:      []string{http.MethodGet},
			url:         []string{"/event/"},
			requestBody: []io.Reader{nil},
			wantBody: []string{strings.Join([]string{string(storage.TestEventJson),
				string(storage.TestEvent2Json),
				string(storage.TestEvent3Json)}, "},{")},
			statusCode: []int{http.StatusOK},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logg.Infof("============== start test %s ==========", tc.name)
			for i := range tc.method {
				req := httptest.NewRequest(tc.method[i], tc.url[i], tc.requestBody[i])
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				logg.Infof("status: %d", w.Code)
				logg.Infof("response: %s", w.Body.String())
				require.Equal(t, tc.statusCode[i], w.Code)

				/*
					var listEvents []model.Event
					err := json.NewDecoder(w.Body).Decode(&listEvents)
					require.NoError(t, err)
					require.Equal(t, tc.wantBody[i], listEvents)
				*/
				l, o := time.Now().Zone()
				logg.Info(l, o)
				exp := strings.Split(tc.wantBody[i], "},{")
				act := strings.Split(w.Body.String(), "},{")

				exp = responseBodyReplace(exp)
				act = responseBodyReplace(act)
				require.ElementsMatch(t, exp, act)
			}
		})
	}
}

func responseBodyReplace(str []string) []string {
	res := make([]string, 0)
	replacer := strings.NewReplacer("[", "", "{", "", "}", "", fmt.Sprintf("%v", time.Local), "")

	for _, s := range str {
		res = append(res, replacer.Replace(s))
	}
	return res
}
