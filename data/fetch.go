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

const QRY = "SELECT value, unix_timestamp(`date`) as timestamp FROM data WHERE type = ? AND device = ? AND `date` > ? order by `date` ASC"

func Fetch(typ string, device string, interval uint64, start string) (*[]Value, error) {
	values := make([]Value, 0, 1024)
	err := gocmn.GetDB().Select(&values, QRY, typ, device, start)
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
	now := uint64(time.Now().Unix())

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
	/*
		for idx, value := range *values {

			if value.Timestamp > ts {
				groups = append(groups, calculateGroup(&group, ts))

				group = make([]*Value, 0, 32)
				ts += interval
			}

			if value.Timestamp < ts+interval {
				group = append(group, &(*values)[idx])
			}
		}

		groups = append(groups, calculateGroup(&group, ts))*/

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
	last := 0.0

	for _, value := range *group {
		if value.Value > max {
			max = value.Value
		}

		if value.Value < min {
			min = value.Value
		}

		total = total + value.Value
		last = value.Value

	}

	avg := total / float64(size)

	return &ValueGroup{timestamp, min, max, avg, last, uint64(size)}
}
