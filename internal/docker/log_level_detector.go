package docker

import (
	"regexp"
	"strings"

	"github.com/mehranzand/pulseup/internal/utils"
)

type LogLevel int

var MapEnumStringToLogType = func() map[string]LogLevel {
	m := make(map[string]LogLevel)
	for i := INFO; i <= WARN; i++ {
		m[i.String()] = i
	}
	return m
}()

const (
	INFO LogLevel = 1 << iota
	DEBUG
	TRACE
	ERROR
	WARN
)

func (l LogLevel) String() string {
	switch l {
	case INFO:
		return "info"
	case DEBUG:
		return "debug"
	case TRACE:
		return "trace"
	case ERROR:
		return "error"
	case WARN:
		return "warn"
	default:
		return "debug"
	}
}

func (l LogLevel) EnumIndex() int {
	return int(l)
}

var keyValueRegex = regexp.MustCompile(`level=(\w+)`)
var plainLevels = map[string]*regexp.Regexp{}
var bracketLevels = map[string]*regexp.Regexp{}
var possibleLogLevels = []string{"debug", "trace", "error", "warn", "info"}

func init() {
	for _, level := range possibleLogLevels {
		plainLevels[level] = regexp.MustCompile("(?i)^" + level + "[^a-z]")
	}

	for _, level := range possibleLogLevels {
		bracketLevels[level] = regexp.MustCompile("(?i)\\[ ?" + level + " ?\\]")
	}
}

func detectLogLevel(message any) LogLevel {
	switch value := message.(type) {
	case string:

		value = utils.StripANSI(value)
		for _, level := range possibleLogLevels {
			if plainLevels[level].MatchString(value) {
				return MapEnumStringToLogType[level]
			}

			if bracketLevels[level].MatchString(value) {
				return MapEnumStringToLogType[level]
			}

			if strings.Contains(value, " "+strings.ToUpper(level)+" ") {
				return MapEnumStringToLogType[level]
			}
		}

		if matches := keyValueRegex.FindStringSubmatch(value); matches != nil {
			return MapEnumStringToLogType[matches[1]]
		}

	default:
		return DEBUG
	}

	return DEBUG
}
