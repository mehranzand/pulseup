package docker

import (
	"regexp"
)

type LogLevel int

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
		return "unknown"
	}
}

func (l LogLevel) EnumIndex() int {
	return int(l)
}

var plainLevels = map[string]*regexp.Regexp{}

var possibleLogLevels = []string{"debug", "trace", "error", "warn", "info"}

func detectLogLevel(message string) {

	//stripedMessage := utils.StripANSI(message)

	// for _, level := range possiableLogLevels {
	// 	if plainLevels[level].MatchString(stripedMessage) {
	// 		return LogLevel
	// 	}

	// 	if bracketLevels[level].MatchString(value) {
	// 		return level
	// 	}

	// 	if strings.Contains(value, " "+strings.ToUpper(level)+" ") {
	// 		return level
	// 	}
	// }

	// if matches := keyValueRegex.FindStringSubmatch(value); matches != nil {
	// 	return matches[1]
	// }

}
