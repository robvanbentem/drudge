package data

import (
	"github.com/robvanbentem/gocmn"
	"math"
	"time"
)

type Value struct {
	Value     float64
	Timestamp uint64
}

type ValueGroup struct {
	Timestamp uint64
	Min       float64
	Max       float64
	Avg       float64
	Last      float64
	Count     uint64
}

const QRY = "SELECT avg(value) as value, unix_timestamp(date) as timestamp FROM data WHERE type = ? AND device = ? AND `date` > ? group by `date` div ? order by timestamp ASC"

func Fetch(typ string, device string, interval uint64, start string) (*[]Value, error) {
	values := make([]Value, 0, 1024)

	if interval > 1800 {
		interval /= 100
	}

	err := gocmn.GetDB().Select(&values, QRY, typ, device, start, interval)
	if err != nil {
		return &values, err
	}

	return &values, nil
}

func Group(values *[]Value, start uint64, interval uint64) (*[]*ValueGroup, error) {
	ts := start - (start % interval) + interval

	groups := make([]*ValueGroup, 0, 128)
	group := make([]*Value, 0, 32)

	idx := 0
	loc, _ := time.LoadLocation("Europe/Amsterdam")
	now := uint64(time.Now().In(loc).Unix()) - (2 * 60 * 60)

	for n := ts; n < now; n += interval {

		for i := idx; i < len(*values); i++ {
			if (*values)[i].Timestamp < n {
				group = append(group, &(*values)[i])
				idx++
			} else {
				break
			}
		}
		groups = append(groups, calculateGroup(&group, n))
		group = make([]*Value, 0, 32)
	}

	if idx < len(*values) {
		for i := idx; i < len(*values); i++ {
			group = append(group, &(*values)[i])
			idx++
		}
		groups = append(groups, calculateGroup(&group, now))
	}

	return &groups, nil
}

func calculateGroup(group *[]*Value, timestamp uint64) *ValueGroup {
	size := len(*group)

	if size == 0 {
		return &ValueGroup{timestamp, 0, 0, 0, 0, 0}
	}

	min := math.MaxFloat64
	max := 0.0
	total := 0.0

	for _, value := range *group {
		if value.Value > max {
			max = value.Value
		}

		if value.Value < min {
			min = value.Value
		}

		total = total + value.Value
	}

	avg := total / float64(size)

	return &ValueGroup{timestamp, min, max, avg, (*group)[size-1].Value, uint64(size)}
}
