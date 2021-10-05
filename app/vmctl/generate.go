package main

import (
	"fmt"
	"os"
	"time"
)

type generateCfg struct {
	timeFrom, timeTo string
	step             time.Duration
	metric           string
	value            float64
	funcName         string
	fileName         string
}

func Generate(cfg generateCfg) error {
	start, from, err := parseTime(cfg.timeFrom, cfg.timeTo)
	if err != nil {
		return fmt.Errorf("failed to parse time in filter: %s", err)
	}
	fu, ok := funcs[cfg.funcName]
	if !ok {
		return fmt.Errorf("unknown func name %q", cfg.funcName)
	}

	fi, err := os.OpenFile(cfg.fileName,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}
	v := cfg.value
	for s := start; s < from; s += cfg.step.Milliseconds() {
		r := fmt.Sprintf("%s %.2f %d", cfg.metric, v, s)
		if _, err := fi.WriteString(r); err != nil {
			return fmt.Errorf("failed to write to file: %s", err)
		}
		fi.Write([]byte{'\n'})
		v = fu(v)
	}
	return fi.Close()
}

type function func(v float64) float64

var funcs = map[string]function{
	"const": func(v float64) float64 { return v },
	"inc":   func(v float64) float64 { return v + 1 },
}

func parseTime(start, end string) (int64, int64, error) {
	var s, e int64
	if start == "" && end == "" {
		return 0, 0, nil
	}
	if start != "" {
		v, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse %q: %s", start, err)
		}
		s = v.UnixNano() / int64(time.Millisecond)
	}
	if end != "" {
		v, err := time.Parse(time.RFC3339, end)
		if err != nil {
			return 0, 0, fmt.Errorf("failed to parse %q: %s", end, err)
		}
		e = v.UnixNano() / int64(time.Millisecond)
	}
	return s, e, nil
}
