package web

import (
	"net/http"
	"github.com/robvanbentem/gocmn"
	"encoding/json"
)

type Result struct {
	Type     string
	Device   string
	Datetime string
	Avg      float64
	Min      float64
	Max      float64
}

const QRY = "SELECT type, device, from_unixtime(avg(unix_timestamp(`date`))) as `datetime`, " +
	"avg(value) as `avg`, min(value) as `min`, max(value) as `max` " +
	"FROM data WHERE type = ? AND device = ? AND `date` > ? " +
	"GROUP BY device, ROUND(UNIX_TIMESTAMP(`date`) / ?);"

func Handle(w http.ResponseWriter, r *http.Request) {
	typ := r.URL.Query().Get("type")
	device := r.URL.Query().Get("device")
	interval := r.URL.Query().Get("interval")
	start := r.URL.Query().Get("start")

	results := []Result{}

	err := gocmn.GetDB().Select(&results, QRY, typ, device, start, interval)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(json)
}
