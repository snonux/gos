package main

import "strings"

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

type healthStatus struct {
	alerts map[string]alert
}

func newHealthStatus() healthStatus {
	return healthStatus{
		alerts: make(map[string]alert),
	}
}

func (hs healthStatus) String() string {
	var (
		alertsBySeverity [4][]string
		sb               strings.Builder
	)

	for _, alert := range hs.alerts {
		alertsBySeverity[alert.severity] = append(alertsBySeverity[alert.severity], alert.text)
	}

	possible := [4]alertSeverity{ok, warning, critical, unknown}
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
		return "OK: all is fine"
	}
	return result
}
