package web

import (
	"net/http"
	"github.com/robvanbentem/gocmn"
	"encoding/json"
	"time"
	"strconv"
)

type Result struct {
	Type     string `json:"type"`
	Device   string `json:"device"`
	Datetime string `json:"date"`
	Avg      float64 `json:"avg"`
	Min      float64 `json:"min"`
	Max      float64`json:"max"`
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

	timestamp, _ := strconv.ParseInt(start, 10, 64)
	start = time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

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
