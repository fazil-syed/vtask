package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

func ParseTimeString(timeStr string) (*time.Time, error) {
	if strings.TrimSpace(timeStr) == "" {
		return nil, nil
	}
	t, err := dateparse.ParseLocal(timeStr)
	if err != nil {
		return nil, err
	}
	fmt.Println(t)
	return &t, nil
}