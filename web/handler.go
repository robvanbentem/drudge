package web

import (
	"net/http"
	"encoding/json"
	"strconv"
	"drudge/data"
	"time"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	typ := r.URL.Query().Get("type")
	device := r.URL.Query().Get("device")
	interval := r.URL.Query().Get("interval")
	start := r.URL.Query().Get("start")

	iStart, _ := strconv.ParseUint(start, 10, 64)
	iInterval, _ := strconv.ParseUint(interval, 10, 64)

	date := time.Unix(int64(iStart), 0)
	format := date.Format("2006-01-02 15:04:05")

	values, err := data.Fetch(typ, device, iInterval, format)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := data.Group(values, iStart, iInterval)

	json, err := json.Marshal(*result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(json)
}
