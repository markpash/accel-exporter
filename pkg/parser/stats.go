package parser

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// Stats represents all statistics gathered from accel-cmd
type Stats struct {
	Uptime        float64
	CPUPercent    float64
	MemRSS        float64
	MemVirt       float64
	Core          CoreStats
	Sessions      SessionStats
	PPPoE         PPPoEStats
	L2TP          L2TPStats
	RadiusServers map[string]RadiusStats
}

// CoreStats contains core metrics
type CoreStats struct {
	MempoolAllocated float64
	MempoolAvailable float64
	ThreadCount      float64
	ThreadActive     float64
	ContextCount     float64
	ContextSleeping  float64
	ContextPending   float64
	MDHandlerCount   float64
	MDHandlerPending float64
	TimerCount       float64
	TimerPending     float64
}

// SessionStats contains session metrics
type SessionStats struct {
	Starting  float64
	Active    float64
	Finishing float64
}

// CollectStats executes accel-cmd and parses its output
func CollectStats(accelCmdPath string) (*Stats, error) {
	cmd := exec.Command(accelCmdPath, "show", "stat")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return parseStats(out.String())
}

var sections = map[string]struct{}{
	"core":     {},
	"sessions": {},
	"pppoe":    {},
	"l2tp":     {},
}

var subsections = map[string]struct{}{
	"tunnels":                     {},
	"sessions (control channels)": {},
	"sessions (data channels)":    {},
}

// parseStats parses the output of accel-cmd show stat
func parseStats(output string) (*Stats, error) {
	stats := &Stats{
		RadiusServers: make(map[string]RadiusStats),
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	var currentSection string
	var currentSubsection string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		// Ignore empty lines
		if line == "" {
			continue
		}

		// Determine if line is section header
		if strings.HasSuffix(line, ":") {
			section := strings.TrimSuffix(line, ":")

			// Determine if it's a main section header
			if _, ok := sections[section]; ok {
				currentSection = section
				continue
			}

			// Determine if it's a subsection header
			if _, ok := subsections[section]; ok {
				currentSubsection = section
				continue
			}
			currentSection = strings.TrimSuffix(line, ":")
			continue
		}

		// Split line into key and value
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Parse key-value pairs based on section
		switch currentSection {
		case "":
			parseMainSection(stats, key, value)
		case "core":
			parseCoreSection(&stats.Core, key, value)
			parseCoreSection_1_13(&stats.Core, key, value)
		case "sessions":
			parseSessionsSection(&stats.Sessions, key, value)
		case "pppoe":
			parsePPPoESection(&stats.PPPoE, key, value)
		case "l2tp":
			parseL2TPSection(&stats.L2TP, currentSubsection, key, value)
		default:
			if strings.HasPrefix(currentSection, "radius") {
				radiusMatch := regexp.MustCompile(`radius\((\d+), ([\d\.]+)\)`).FindStringSubmatch(currentSection)
				if len(radiusMatch) == 3 {
					radiusID := radiusMatch[1]
					radiusIP := radiusMatch[2]
					if _, exists := stats.RadiusServers[radiusID]; !exists {
						stats.RadiusServers[radiusID] = RadiusStats{
							ID: radiusID,
							IP: radiusIP,
						}
					}

					rs := stats.RadiusServers[radiusID]
					parseRadiusSection(&rs, key, value)
					stats.RadiusServers[radiusID] = rs
				}
			}
		}
	}

	return stats, scanner.Err()
}

// Helper functions to parse each section...
func parseMainSection(stats *Stats, key, value string) {
	switch key {
	case "uptime":
		stats.Uptime = parseUptime(value)
	case "cpu":
		stats.CPUPercent = parsePercentage(value)
	case "mem(rss/virt)":
		parseMemory(stats, value)
	}
}

func parseUptime(value string) float64 {
	// Parse uptime in format "138.00:05:20" (days.hours:minutes:seconds)
	parts := strings.Split(value, ".")
	if len(parts) != 2 {
		return 0
	}

	days, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}

	timeParts := strings.Split(parts[1], ":")
	if len(timeParts) != 3 {
		return 0
	}

	hours, _ := strconv.ParseFloat(timeParts[0], 64)
	minutes, _ := strconv.ParseFloat(timeParts[1], 64)
	seconds, _ := strconv.ParseFloat(timeParts[2], 64)

	return days*86400 + hours*3600 + minutes*60 + seconds
}

func parsePercentage(value string) float64 {
	// Example: "1.23%"
	trimmed := strings.TrimSuffix(value, "%")
	f, _ := strconv.ParseFloat(trimmed, 64)
	return f
}

func parseMemory(stats *Stats, value string) {
	// Example: "12345 / 67890 K"
	parts := strings.Split(strings.TrimSuffix(value, " K"), "/")
	if len(parts) == 2 {
		rss, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		virt, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		stats.MemRSS = rss
		stats.MemVirt = virt
	}
}

func parseCoreSection(core *CoreStats, key, value string) {
	// Implement parsing logic based on key
	switch key {
	case "mempool(allocated/available)":
		// Example: "1024 / 2048"
		parts := strings.Split(value, "/")
		if len(parts) == 2 {
			alloc, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			avail, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			core.MempoolAllocated = alloc
			core.MempoolAvailable = avail
		}
	case "threads(count/active)":
		parts := strings.Split(value, "/")
		if len(parts) == 2 {
			count, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			active, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			core.ThreadCount = count
			core.ThreadActive = active
		}
	case "context(count/sleep/pending)":
		parts := strings.Split(value, "/")
		if len(parts) == 3 {
			count, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			sleep, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			pending, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
			core.ContextCount = count
			core.ContextSleeping = sleep
			core.ContextPending = pending
		}
	case "md_handler(count/pending)":
		parts := strings.Split(value, "/")
		if len(parts) == 2 {
			count, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			pending, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			core.MDHandlerCount = count
			core.MDHandlerPending = pending
		}
	case "timer(count/pending)":
		parts := strings.Split(value, "/")
		if len(parts) == 2 {
			count, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			pending, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			core.TimerCount = count
			core.TimerPending = pending
		}
	}
}

func parseCoreSection_1_13(core *CoreStats, key, value string) {
	f, _ := strconv.ParseFloat(strings.TrimSpace(value), 64)
	// Implement parsing logic based on key
	switch key {
	case "mempool_allocated":
		core.MempoolAllocated = f
	case "mempool_available":
		core.MempoolAvailable = f
	case "thread_count":
		core.ThreadCount = f
	case "thread_active":
		core.ThreadActive = f
	case "context_count":
		core.ContextCount = f
	case "context_sleeping":
		core.ContextSleeping = f
	case "context_pending":
		core.ContextPending = f
	case "md_handler_count":
		core.MDHandlerCount = f
	case "md_handler_pending":
		core.MDHandlerPending = f
	case "timer_count":
		core.TimerCount = f
	case "timer_pending":
		core.TimerPending = f
	}
}

func parseSessionsSection(sessions *SessionStats, key, value string) {
	f, _ := strconv.ParseFloat(value, 64)
	switch key {
	case "starting":
		sessions.Starting = f
	case "active":
		sessions.Active = f
	case "finishing":
		sessions.Finishing = f
	}
}
