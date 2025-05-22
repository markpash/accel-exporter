package parser

import "strconv"

// L2TPStats contains L2TP protocol metrics
type L2TPStats struct {
	TunnelsStarting                  float64
	TunnelsActive                    float64
	TunnelsFinishing                 float64
	SessionsControlChannelsStarting  float64
	SessionsControlChannelsActive    float64
	SessionsControlChannelsFinishing float64
	SessionsDataChannelsStarting     float64
	SessionsDataChannelsActive       float64
	SessionsDataChannelsFinishing    float64
}

func parseL2TPSection(l2tp *L2TPStats, subsection, key, value string) {
	f, _ := strconv.ParseFloat(value, 64)
	switch subsection {
	case "tunnels":
		switch key {
		case "starting":
			l2tp.TunnelsStarting = f
		case "active":
			l2tp.TunnelsActive = f
		case "finishing":
			l2tp.TunnelsFinishing = f
		}
	case "sessions (control channels)":
		switch key {
		case "starting":
			l2tp.SessionsControlChannelsStarting = f
		case "active":
			l2tp.SessionsControlChannelsActive = f
		case "finishing":
			l2tp.SessionsControlChannelsFinishing = f
		}
	case "sessions (data channels)":
		switch key {
		case "starting":
			l2tp.SessionsDataChannelsStarting = f
		case "active":
			l2tp.SessionsDataChannelsActive = f
		case "finishing":
			l2tp.SessionsDataChannelsFinishing = f
		}
	}
}
