package parser

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

// At least 5 time series can be extracted from sessions:
// - uptime-raw
// - rx-bytes-raw
// - tx-bytes-raw
// - rx-pkts
// - tx-pkts

// CollectSessions executes accel-cmd and parses its output
func CollectSessions(accelCmdPath string) ([]*Session, error) {
	cmd := exec.Command(accelCmdPath, "show", "sessions", "type,ifname,username,ip,calling-sid,called-sid,sid,uptime-raw,rx-bytes-raw,tx-bytes-raw,rx-pkts,tx-pkts")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return parseSessions(out.String())
}

type Session struct {
	Type       string
	Ifname     string
	Username   string
	IP         string
	CallingSID string
	CalledSID  string
	SID        string
	Uptime     float64
	RxBytes    float64
	TxBytes    float64
	RxPkts     float64
	TxPkts     float64
}

func parseSessions(output string) ([]*Session, error) {
	sessions := []*Session{}
	scanner := bufio.NewScanner(strings.NewReader(output))

	// Skip the header lines (heading and separator)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		lineNum++

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip the header line
		if lineNum == 1 {
			continue
		}

		// Skip the separator line (------+-----+...)
		if strings.Contains(line, "-+") {
			continue
		}

		// Process data rows
		// Split the line using vertical bars and trim whitespace
		fields := strings.Split(line, "|")

		// Trim whitespace from each field
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}

		// Parse numeric fields
		uptime, _ := strconv.ParseFloat(fields[7], 64)
		rxBytes, _ := strconv.ParseFloat(fields[8], 64)
		txBytes, _ := strconv.ParseFloat(fields[9], 64)
		rxPkts, _ := strconv.ParseFloat(fields[10], 64)
		txPkts, _ := strconv.ParseFloat(fields[11], 64)

		// Create a new Session object
		session := &Session{
			Type:       fields[0],
			Ifname:     fields[1],
			Username:   fields[2],
			IP:         fields[3],
			CallingSID: fields[4],
			CalledSID:  fields[5],
			SID:        fields[6],
			Uptime:     uptime,
			RxBytes:    rxBytes,
			TxBytes:    txBytes,
			RxPkts:     rxPkts,
			TxPkts:     txPkts,
		}

		sessions = append(sessions, session)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}
