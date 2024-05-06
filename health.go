package main

import (
	"fmt"
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

	delete(hs.alerts, what)
}

func (hs healthStatus) String() string {
	var (
		alertsBySeverity [4][]string
		sb               strings.Builder
	)

	hs.mu.Lock()
	defer hs.mu.Unlock()

	for _, alert := range hs.alerts {
		alertsBySeverity[alert.severity] = append(alertsBySeverity[alert.severity], alert.String())
	}

	possible := [4]alertSeverity{unknown, critical, warning, ok}
	for _, severity := range possible {
		if len(alertsBySeverity[severity]) == 0 {
			continue
		}
		for _, alert := range alertsBySeverity[severity] {
			sb.WriteString(alert)
			sb.WriteString("\n")
		}
	}

	result := sb.String()
	if result == "" {
		return "OK: all is fine\n"
	}
	return result
}
