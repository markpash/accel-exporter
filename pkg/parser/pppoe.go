package parser

import "strconv"

// PPPoEStats contains PPPoE protocol metrics
type PPPoEStats struct {
	Starting    float64
	Active      float64
	DelayedPADO float64
	RecvPADI    float64
	DropPADI    float64
	SentPADO    float64
	RecvPADR    float64
	RecvPADRDup float64
	SentPADS    float64
	Filtered    float64
}

func parsePPPoESection(pppoe *PPPoEStats, key, value string) {
	f, _ := strconv.ParseFloat(value, 64)
	switch key {
	case "starting":
		pppoe.Starting = f
	case "active":
		pppoe.Active = f
	case "delayed PADO":
		pppoe.DelayedPADO = f
	case "recv PADI":
		pppoe.RecvPADI = f
	case "drop PADI":
		pppoe.DropPADI = f
	case "sent PADO":
		pppoe.SentPADO = f
	case "recv PADR":
		pppoe.RecvPADR = f
	case "recv PADR(dup)":
		pppoe.RecvPADRDup = f
	case "sent PADS":
		pppoe.SentPADS = f
	case "filtered":
		pppoe.Filtered = f
	}
}
