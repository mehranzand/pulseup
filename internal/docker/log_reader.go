package docker

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	Stdin StdType = iota
	Stdout
	Stderr

	writerPrefixLen = 8
	writerFlagIndex = 0
	writerSizeIndex = 4

	startingBufLen = 32*1024 + writerPrefixLen + 1
)

type Log struct {
	Message   any   `json:"m,omitempty"`
	Timestamp int64 `json:"ts"`
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
	go logReader.bufferProcess()

	return logReader
}

func (e *LogReader) readerProcess(tty bool) {
	_, err := e.readLog(e.reader, tty)

	if err != nil {
		fmt.Print(err)
	}

	e.wg.Done()
}

func (e *LogReader) bufferProcess() {
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

func (e *LogReader) readLog(src io.Reader, tty bool) (written int64, err error) {
	var (
		buf       = make([]byte, startingBufLen)
		bufLen    = len(buf)
		nr        int
		er        error
		frameSize int
	)

	for {
		// Make sure we have at least a full header
		for nr < writerPrefixLen {
			var nr2 int
			nr2, er = src.Read(buf[nr:])
			nr += nr2
			if er == io.EOF {
				if nr < writerPrefixLen {
					logrus.Debugf("Corrupted prefix: %v", buf[:nr])
					return written, nil
				}
				break
			}
			if er != nil {
				logrus.Debugf("Error reading header: %s", er)
				return 0, er
			}
		}

		// Check the first byte to know where to write
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

		logEvent := LogEvent{Message: string(buf[writerPrefixLen : frameSize+writerPrefixLen])}
		e.buffer <- &logEvent

		logrus.Debugf("Channel buffer size: %d", len(e.buffer))

		copy(buf, buf[frameSize+writerPrefixLen:])
		nr -= frameSize + writerPrefixLen
	}
}
