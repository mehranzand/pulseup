package docker

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	Stdin StdType = iota
	Stdout
	Stderr

	writerPrefixLen = 8
	writerFlagIndex = 0
	writerSizeIndex = 4

	startingBufLen = 1*1024*writerPrefixLen + 1
)

type Log struct {
	Timestamp  int64    `json:"ts"`
	Message    any      `json:"m,omitempty"`
	Data       string   `json:"data,omitempty"`
	ActionTags []string `json:"action_tags,omitempty"`
}

type LogReader struct {
	Events chan *LogEvent
	buffer chan *LogEvent
	reader *bufio.Reader
	wg     sync.WaitGroup
}

func NewLogReader(reader io.Reader, tty bool) *LogReader {
	logReader := &LogReader{
		Events: make(chan *LogEvent),
		reader: bufio.NewReader(reader),
		buffer: make(chan *LogEvent, 100),
	}

	logReader.wg.Add(2)

	go logReader.readerProcess(tty)
	go logReader.bufferEmitter()

	return logReader
}

func (e *LogReader) readerProcess(tty bool) {
	_, err := e.readLog(e.reader, tty)

	if err != nil {
		logrus.Debug(err)
	}

	e.wg.Done()
}

func (e *LogReader) bufferEmitter() {
	for {
		event, ok := <-e.buffer
		if !ok {
			close(e.Events)
			break
		}

		e.Events <- event
	}

	e.wg.Done()
}

func (e *LogReader) readLog(src *bufio.Reader, tty bool) (written int64, err error) {
	var (
		buf       = make([]byte, startingBufLen)
		bufLen    = len(buf)
		nr        int
		er        error
		frameSize int
	)

	for {
		if tty {
			message, err := src.ReadString('\n')
			if err != nil {
				close(e.buffer)

				return 0, err
			}

			e.parseAndPushEvent(message)

		} else {
			for nr < writerPrefixLen {
				var nr2 int
				nr2, er = src.Read(buf[nr:])
				nr += nr2
				logrus.Debugf("nt -> %d", nr)
				if er == io.EOF {
					if nr < writerPrefixLen {
						logrus.Debugf("Corrupted prefix: %v", buf[:nr])
						return written, er
					}
					break
				}
				if er != nil {
					logrus.Debugf("Error reading header: %s", er)
					return 0, er
				}
			}

			switch StdType(buf[writerFlagIndex]) {
			case Stdin:
				fallthrough
			case Stdout:
				// Write on stdout
				//out = dstout
			case Stderr:
				// Write on stderr
				//out = dsterr
			default:
				logrus.Debugf("Error selecting output flag index: (%d)", buf[writerFlagIndex])
				return 0, fmt.Errorf("Unrecognized input header: %d", buf[writerFlagIndex])
			}

			frameSize = int(binary.BigEndian.Uint32(buf[writerSizeIndex : writerSizeIndex+4]))
			logrus.Debugf("framesize: %d", frameSize)

			if frameSize+writerPrefixLen > bufLen {
				logrus.Debugf("Extending buffer cap by %d (was %d)", frameSize+writerPrefixLen-bufLen+1, len(buf))
				buf = append(buf, make([]byte, frameSize+writerPrefixLen-bufLen+1)...)
				bufLen = len(buf)
			}

			for nr < frameSize+writerPrefixLen {
				var nr2 int
				nr2, er = src.Read(buf[nr:])
				nr += nr2
				if er == io.EOF {
					if nr < frameSize+writerPrefixLen {
						logrus.Debugf("Corrupted frame: %v", buf[writerPrefixLen:nr])
						return written, nil
					}
					break
				}
				if er != nil {
					logrus.Debugf("Error reading frame: %s", er)
					return 0, er
				}
			}

			e.parseAndPushEvent(string(buf[writerPrefixLen : frameSize+writerPrefixLen]))

			logrus.Debugf("Channel buffer size: %d", len(e.buffer))

			copy(buf, buf[frameSize+writerPrefixLen:])
			nr -= frameSize + writerPrefixLen
		}

		printMemUsage()
	}
}

func (e *LogReader) parseAndPushEvent(message string) {
	logEvent := &LogEvent{Message: message}

	if index := strings.IndexAny(message, " "); index != -1 {
		stamp := message[:index]
		logEvent.Message = strings.TrimSuffix(message[index+1:], "")

		if timestamp, err := time.Parse(time.RFC3339Nano, stamp); err == nil {
			logEvent.Timestamp = timestamp.UnixMilli()
		}
	}

	e.buffer <- logEvent
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
