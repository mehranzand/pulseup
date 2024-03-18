package handler

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mehranzand/pulseup/api/middleware"
	"github.com/mehranzand/pulseup/internal/docker"
	log "github.com/sirupsen/logrus"
)

// GetContainers
// @Summary Get list of containers
func (h *Handler) GetContainers(c echo.Context) error {
	cc := c.(*middleware.DockerContext)

	if cc.Client == nil {
		http.Error(c.Response().Writer, "Docker host not found!", http.StatusInternalServerError)
	}

	containers, err := cc.Client.ListContainers()
	if err != nil {
		http.Error(c.Response().Writer, "Continer not found!", http.StatusInternalServerError)
	}

	enc := json.NewEncoder(c.Response())
	for _, l := range containers {
		if err := enc.Encode(l.Name); err != nil {
			return err
		}
		fmt.Fprintf(c.Response().Writer, "| ")
		c.Response().Flush()
	}

	return nil
}

func (h *Handler) StreamLogs(c echo.Context) error {
	var stdTypes docker.StdType
	stdTypes |= docker.STDOUT

	cc := c.(*middleware.DockerContext)
	id := c.Param("id")

	_, ok := cc.Context.Response().Writer.(http.Flusher)
	if !ok {
		http.Error(c.Response().Writer, "Streaming unsupported!", http.StatusInternalServerError)
	}

	container, err := cc.Client.FindContainer(id)
	if err != nil {
		http.Error(c.Response().Writer, "Continer not found!", http.StatusInternalServerError)
	}

	since := time.Now().AddDate(0, 0, 0)
	reader, err := cc.Client.ContainerLogs(cc.Context.Request().Context(), container.ID, strconv.FormatInt(since.Unix(), 10), stdTypes)
	if err != nil {
		http.Error(cc.Context.Response().Writer, "Continer not found!", http.StatusInternalServerError)
	}

	msg, _ := readLog(reader, container.Tty)

	log.Infof("%q\n", msg)

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-transform")
	c.Response().Header().Add(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")

	// f, err := os.Open("log")
	// if err != nil {
	// 	log.Fatalf("unable to read file: %v", err)
	// }
	// defer f.Close()
	// buf := make([]byte, 1024)
	// for {
	// 	n, err := f.Read(buf)
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}

	// 	if n > 0 {
	// 		fmt.Println(string(buf[:n]))
	// 	}

	// }

	return nil
}

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return bytes.NewBuffer(nil)
	},
}

func readLog(r io.Reader, tty bool) (string, error) {
	header := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	reader := bufio.NewReader(r)
	buffer := bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()

	if tty {
		message, err := reader.ReadString('\n')
		if err != nil {
			return message, err
		}
		return message, nil
	} else {
		n, err := io.ReadFull(reader, header)
		if err != nil {
			return "", err
		}
		if n != 8 {
			log.Warnf("unable to read header: %v", header)
			message, _ := reader.ReadString('\n')
			return message, nil
		}

		fmt.Println(header[0])

		switch header[0] {
		case 1:

		case 2:

		default:
			log.Warnf("unknown stream type: %v", header[0])
		}

		count := binary.BigEndian.Uint32(header[4:])
		if count == 0 {
			return "", nil
		}
		_, err = io.CopyN(buffer, reader, int64(count))
		if err != nil {
			return "", err
		}
		return buffer.String(), nil
	}

}
