package docker

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"io"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

const (
	Stdin StdType = iota
	Stdout
	Stderr

	headerLen = 8
	flagIndex = 0
	sizeIndex = 4
)

var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

type Log struct {
	StdType    StdType  `json:"type"`
	Timestamp  int64    `json:"ts"`
	Message    any      `json:"m,omitempty"`
	Data       string   `json:"d,omitempty"`
	ActionTags []string `json:"tag,omitempty"`
}

type LogReader struct {
	Events chan *LogEvent
	buffer chan *LogEvent
	reader *bufio.Reader
	wg     sync.WaitGroup
	Errors chan error
}

var ErrBadHeader = fmt.Errorf("unable to read header")

func NewLogReader(reader io.Reader, tty bool) *LogReader {
	logReader := &LogReader{
		Events: make(chan *LogEvent),
		reader: bufio.NewReader(reader),
		buffer: make(chan *LogEvent, 100),
		Errors: make(chan error, 1),
	}

	logReader.wg.Add(2)

	go logReader.readerProcess(tty)
	go logReader.bufferEmitter()

	return logReader
}

func (e *LogReader) readerProcess(tty bool) {
	err := e.readLog(e.reader, tty)

	if err != nil {
		logrus.Tracef("log reader error %v", err)
		e.Errors <- err
		close(e.buffer)
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

func (e *LogReader) readLog(reader *bufio.Reader, tty bool) error {
	for {
		var frameSize int
		var stdType StdType = STDOUT
		buffer := bufPool.Get().(*bytes.Buffer)
		header := make([]byte, headerLen)
		defer bufPool.Put(buffer)
		buffer.Reset()

		if tty {
			message, err := reader.ReadString('\n')
			if err != nil {
				return err
			}

			e.parseAndPushEvent(message, stdType)

		} else {
			n, err := io.ReadFull(reader, header)

			if err != nil {
				return err
			}

			if n != 8 {
				log.Warnf("unable to read header: %v", header)
				message, _ := reader.ReadString('\n')

				e.parseAndPushEvent(message, stdType)

				return ErrBadHeader
			}

			switch header[0] {
			case 1:
				stdType = STDOUT
			case 2:
				stdType = STDERR
			default:
				log.Warnf("undefined std type: %v", header[0])
			}

			frameSize = int(binary.BigEndian.Uint32(header[4:]))
			//logrus.Warnf("framesize: %d", frameSize)

			if frameSize == 0 {
				continue
			}

			_, err = io.CopyN(buffer, reader, int64(frameSize))
			if err != nil {
				return err
			}

			e.parseAndPushEvent(buffer.String(), stdType)
		}

		if false {
			printMemUsage()
		}
	}
}

func (e *LogReader) parseAndPushEvent(message string, stdType StdType) {
	hash := fnv.New32a()
	hash.Write([]byte(message))
	logEvent := &LogEvent{Id: hash.Sum32(), Message: message, StdType: stdType.String()}

	if index := strings.IndexAny(message, " "); index != -1 {
		stamp := message[:index]
		logEvent.Message = strings.TrimSuffix(message[index+1:], "")

		if timestamp, err := time.Parse(time.RFC3339Nano, stamp); err == nil {
			logEvent.Timestamp = timestamp.UnixMilli()
		}
	}

	logEvent.Level = detectLogLevel(message)

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
