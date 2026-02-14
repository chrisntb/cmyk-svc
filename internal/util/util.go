package util

import (
	"encoding/json"
	"time"

	"github.com/tailscale/hujson"
)

func DaysElapsedAbsolute(a, b time.Time) int {
	if a.Before(b) {
		a, b = b, a
	}
	daysFloat := a.Sub(b).Seconds()
	days := int(daysFloat / 86400)
	if daysFloat/1 >= 0 {
		days++
	}
	return days
}

func StandardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()
	return ast.Pack(), nil
}

func Print(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}
