package main

import (
	"drudge/common"
	"drudge/web"
	"encoding/json"
	"github.com/robvanbentem/gocmn"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Payload struct {
	Type string
	Data json.RawMessage
}

func main() {
	common.LoadConfig()

	gocmn.InitLogger(common.ConfigRoot.LogFile)
	defer gocmn.ShutdownLogger()

	gocmn.InitDB(common.ConfigRoot.Database)
	defer gocmn.CloseDB()

	r := mux.NewRouter()
	registerRoutes(r)

	serve := fmt.Sprintf("%s:%d", common.ConfigRoot.Host, common.ConfigRoot.Port)

	gocmn.Log.Infof("Server started on %s\n", serve)
	gocmn.Log.Fatal(http.ListenAndServe(serve, r))
}

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/data", web.Handle).Methods("GET")
}
