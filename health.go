package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

type alertSeverity int

const (
	ok alertSeverity = iota
	warning
	critical
	unknown
)

func (s alertSeverity) String() string {
	switch s {
	case ok:
		return "OK"
	case warning:
		return "WARNING"
	case critical:
		return "CRITICAL"
	case unknown:
		return "UNKNOWN"
	default:
		panic("encountered an unknown alertSeverity")
	}
}

type alert struct {
	text     string
	severity alertSeverity
}

func (a alert) String() string {
	return fmt.Sprintf("%s: %s", a.severity, a.text)
}

type healthStatus struct {
	alerts map[string]alert
	mu     sync.Mutex
}

func newHealthStatus() healthStatus {
	return healthStatus{
		alerts: make(map[string]alert),
	}
}

func (hs healthStatus) set(s alertSeverity, what, text string) {
	log.Println("alerting", what, "to", text, "with severity", s)

	hs.mu.Lock()
	defer hs.mu.Unlock()

	hs.alerts[what] = alert{
		text:     text,
		severity: s,
	}
}

func (hs healthStatus) clear(what string) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if _, ok := hs.alerts[what]; ok {
		log.Println("clearing alert for", what)
		delete(hs.alerts, what)
	}
}

func (hs healthStatus) String() string {
	var (
		alerts [4][]string // Alerts by severity
		sb     strings.Builder
	)

	hs.mu.Lock()
	defer hs.mu.Unlock()

	for _, alert := range hs.alerts {
		alerts[alert.severity] = append(alerts[alert.severity], alert.String())
	}

	possible := [4]alertSeverity{unknown, critical, warning, ok}
	for _, severity := range possible {
		if len(alerts[severity]) == 0 {
			continue
		}
		for _, alert := range alerts[severity] {
			sb.WriteString(alert)
			sb.WriteString("\n")
		}
	}

	if result := sb.String(); result != "" {
		return result
	}
	return "OK: all is fine :-)\n"
}
