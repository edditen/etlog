package utils

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

var (
	// MinuteSeconds seconds of one minute
	MinuteSeconds = 60
	// HourSeconds seconds of one hour
	HourSeconds = 60 * MinuteSeconds
	// DaySeconds seconds of one day
	DaySeconds = 24 * HourSeconds
	// KBytes size of k
	KBytes = 1024
	// MBytes size of m
	MBytes = 1024 * KBytes
	// GBytes size of g
	GBytes = 1024 * MBytes

	// ErrTimeParse time parse error
	ErrTimeParse = errors.New("time parse error")
	// ErrSizeParse size parse error
	ErrSizeParse = errors.New("file size parse error")
)

var timeUnits = map[byte]int{
	'h': HourSeconds,
	'd': DaySeconds,
	'm': MinuteSeconds,
}

var sizeUnits = map[byte]int{
	'k': KBytes,
	'm': MBytes,
	'g': GBytes,
}

func LastSubstring(s string, sep string) string {
	idx := strings.LastIndex(s, sep)
	return s[idx+1:]
}

// ParseSize parse the file size
func ParseSize(fileSize string) (int, error) {
	if len(fileSize) == 0 {
		return -1, errors.New("fileSize is blank")
	}

	sizeStr := strings.ToLower(fileSize)
	buf := []byte(sizeStr)

	i, err := strconv.Atoi(string(buf[0 : len(buf)-1]))
	if err != nil || i < 0 {
		return -1, ErrSizeParse
	}
	unitByte := buf[len(buf)-1]
	units, ok := sizeUnits[unitByte]
	if !ok {
		return -1, ErrSizeParse
	}

	sd := units * i
	if sd <= 0 {
		return -1, ErrSizeParse
	}

	return sd, nil
}

// ParseSeconds from 1h 10m 10d
func ParseSeconds(interval string) (int, error) {
	if len(interval) == 0 {
		return -1, errors.New("interval is blank")
	}

	timeStr := strings.ToLower(interval)
	buf := []byte(timeStr)
	i, err := strconv.Atoi(string(buf[0 : len(buf)-1]))
	if err != nil || i < 0 {
		return -1, ErrTimeParse
	}
	unitByte := buf[len(buf)-1]
	units, ok := timeUnits[unitByte]
	if !ok {
		return -1, ErrTimeParse
	}
	sd := units * i
	if sd <= 0 {
		return -1, ErrTimeParse
	}
	return sd, nil
}
