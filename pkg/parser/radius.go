package parser

import (
	"strconv"
	"strings"
)

// RadiusStats contains RADIUS server metrics
type RadiusStats struct {
	ID               string
	IP               string
	State            string
	FailCount        float64
	RequestCount     float64
	QueueLength      float64
	AuthSent         float64
	AuthLostTotal    float64
	AuthLost5m       float64
	AuthLost1m       float64
	AuthAvgTime5m    float64
	AuthAvgTime1m    float64
	AcctSent         float64
	AcctLostTotal    float64
	AcctLost5m       float64
	AcctLost1m       float64
	AcctAvgTime5m    float64
	AcctAvgTime1m    float64
	InterimSent      float64
	InterimLostTotal float64
	InterimLost5m    float64
	InterimLost1m    float64
	InterimAvgTime5m float64
	InterimAvgTime1m float64
}

func parseRadiusSection(radius *RadiusStats, key, value string) {
	f, _ := strconv.ParseFloat(value, 64)
	switch key {
	case "state":
		radius.State = value // State is a string
	case "fail count":
		radius.FailCount = f
	case "request count":
		radius.RequestCount = f
	case "queue length":
		radius.QueueLength = f
	case "auth sent":
		radius.AuthSent = f
	case "auth lost(total/5m/1m)":
		parts := strings.Split(value, "/")
		if len(parts) == 3 {
			total, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			m5, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			m1, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
			radius.AuthLostTotal = total
			radius.AuthLost5m = m5
			radius.AuthLost1m = m1
		}
	case "auth avg time(5m/1m)":
		parts := strings.Split(value, "/")
		if len(parts) == 2 {
			m5, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			m1, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			radius.AuthAvgTime5m = m5
			radius.AuthAvgTime1m = m1
		}
	case "acct sent":
		radius.AcctSent = f
	case "acct lost(total/5m/1m)":
		parts := strings.Split(value, "/")
		if len(parts) == 3 {
			total, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			m5, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			m1, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
			radius.AcctLostTotal = total
			radius.AcctLost5m = m5
			radius.AcctLost1m = m1
		}
	case "acct avg time(5m/1m)":
		parts := strings.Split(value, "/")
		if len(parts) == 2 {
			m5, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			m1, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			radius.AcctAvgTime5m = m5
			radius.AcctAvgTime1m = m1
		}
	case "interim sent":
		radius.InterimSent = f
	case "interim lost(total/5m/1m)":
		parts := strings.Split(value, "/")
		if len(parts) == 3 {
			total, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			m5, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			m1, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
			radius.InterimLostTotal = total
			radius.InterimLost5m = m5
			radius.InterimLost1m = m1
		}
	case "interim avg time(5m/1m)":
		parts := strings.Split(value, "/")
		if len(parts) == 2 {
			m5, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			m1, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			radius.InterimAvgTime5m = m5
			radius.InterimAvgTime1m = m1
		}
	}
}
