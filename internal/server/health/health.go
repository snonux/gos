package health

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

type Severity int

const (
	OK Severity = iota
	Warning
	Critical
	Unknown
)

func (s Severity) String() string {
	switch s {
	case OK:
		return "OK"
	case Warning:
		return "WARNING"
	case Critical:
		return "CRITICAL"
	case Unknown:
		fallthrough
	default:
		return "UNKNOWN"
	}
}

type alert struct {
	text     string
	severity Severity
}

func (a alert) String() string {
	return fmt.Sprintf("%s: %s", a.severity, a.text)
}

type Status struct {
	alerts map[string]alert
	mu     *sync.Mutex
}

func NewStatus() Status {
	return Status{
		alerts: make(map[string]alert),
		mu:     &sync.Mutex{},
	}
}

func (hs Status) Set(s Severity, healthStatusKey string, info any) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	infoStr := fmt.Sprintf("%v", info)
	log.Printf("status: alerting %s as %s: %s", healthStatusKey, s, infoStr)

	hs.alerts[healthStatusKey] = alert{
		text:     infoStr,
		severity: s,
	}
}

func (hs Status) Clear(healthStatusKey string) {
	hs.mu.Lock()
	defer hs.mu.Unlock()

	if _, ok := hs.alerts[healthStatusKey]; ok {
		log.Println("status: clearing ", healthStatusKey)
		delete(hs.alerts, healthStatusKey)
	}
}

func (hs Status) String() string {
	var (
		alerts [4][]string // Alerts by severity
		sb     strings.Builder
	)

	hs.mu.Lock()
	defer hs.mu.Unlock()

	for healthStatusKey, alert := range hs.alerts {
		str := fmt.Sprintf("%s (handler %s)", alert, healthStatusKey)
		alerts[alert.severity] = append(alerts[alert.severity], str)
	}

	possible := [4]Severity{Unknown, Critical, Warning, OK}
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
