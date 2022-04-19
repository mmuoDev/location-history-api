package app

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mmuoDev/location-history-api.git/internal/db"
	conn "github.com/mmuoDev/location-history-api.git/pkg/db"
)

//App contains handlers for this app
type App struct {
	AddLocationHandler     http.HandlerFunc
	RetrieveHistoryHandler http.HandlerFunc
	DeleteHistoryHandler   http.HandlerFunc
}

//Handler returns the main handler for this app
func (a App) Handler() http.HandlerFunc {
	router := httprouter.New()
	router.HandlerFunc(http.MethodPost, fmt.Sprintf("/location/:%s/now", orderID), a.AddLocationHandler)
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("/location/:%s", orderID), a.RetrieveHistoryHandler)
	router.HandlerFunc(http.MethodDelete, fmt.Sprintf("/location/:%s", orderID), a.DeleteHistoryHandler)

	return http.HandlerFunc(router.ServeHTTP)
}

// /OptionalArgs optional arguments for this app
type OptionalArgs struct {
	AddLocation     db.AddLocationFunc
	RetrieveHistory db.RetrieveHistoryFunc
	DeleteHistory   db.DeleteHistoryFunc
}

// Options is a type for application options to modify the app
type Options func(o *OptionalArgs)

//New creates a new instance of the App
func New(storage *conn.Storage, options ...Options) App {
	o := OptionalArgs{
		AddLocation:     db.AddLocation(*storage),
		RetrieveHistory: db.RetrieveHistory(*storage),
		DeleteHistory:   db.DeleteHistory(*storage),
	}
	for _, option := range options {
		option(&o)
	}
	addLocationHandler := AddLocationHandler(o.AddLocation)
	retrieveHistoryHandler := RetrieveHistoryHandler(o.RetrieveHistory)
	deleteHistoryHandler := DeleteHistoryHandler(o.DeleteHistory)
	return App{
		AddLocationHandler:     addLocationHandler,
		RetrieveHistoryHandler: retrieveHistoryHandler,
		DeleteHistoryHandler:   deleteHistoryHandler,
	}
}
