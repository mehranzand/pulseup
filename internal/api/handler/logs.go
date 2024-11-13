package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mehranzand/pulseup/internal/api/middleware"
	"github.com/mehranzand/pulseup/internal/docker"
	"github.com/sirupsen/logrus"
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
	fmt.Fprintf(c.Response().Writer, "[")

	enc := json.NewEncoder(c.Response())
	for i, l := range containers {

		if err := enc.Encode(l); err != nil {
			return err
		}

		if len(containers) != i+1 {
			fmt.Fprintf(c.Response().Writer, ",")
		}

		c.Response().Flush()
	}
	fmt.Fprintf(c.Response().Writer, "]")

	return nil
}

func (h *Handler) StreamLogs(c echo.Context) error {
	var stdTypes docker.StdType
	stdTypes |= docker.STDOUT
	stdTypes |= docker.STDERR

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-transform")
	c.Response().Header().Add(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")

	cc := c.(*middleware.DockerContext)
	id := c.Param("id")

	f, ok := cc.Context.Response().Writer.(http.Flusher)
	if !ok {
		http.Error(c.Response().Writer, "Streaming unsupported!", http.StatusInternalServerError)
		return nil
	}

	container, err := cc.Client.FindContainer(id)
	if err != nil {
		http.Error(c.Response().Writer, "Continer not found!", http.StatusNotFound)
		return nil
	}

	since := time.Now().AddDate(0, 0, -10)
	reader, err := cc.Client.ContainerLogs(c.Request().Context(), container.ID, strconv.FormatInt(since.Unix(), 10), stdTypes)
	if err != nil {
		if err == io.EOF {
			fmt.Fprintf(c.Response().Writer, "event: container-stopped\ndata: end of stream\n\n")
			f.Flush()
		} else {
			http.Error(c.Response().Writer, err.Error(), http.StatusInternalServerError)
		}
		return nil
	}

	lr := docker.NewLogReader(reader, container.Tty)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
outerloop:
	for {
		select {
		case event, ok := <-lr.Events:
			if !ok {
				break outerloop
			}

			if buffer, err := json.Marshal(event); err != nil {
				logrus.Errorf("json encoding error while streaming %v", err.Error())
			} else {
				fmt.Fprintf(c.Response().Writer, "data: %s\n", buffer)
				f.Flush()
			}

			fmt.Fprintf(c.Response().Writer, "\n")
			f.Flush()

		case <-ticker.C:
			fmt.Fprintln(c.Response().Writer, "ping:")
			f.Flush()
		}
	}

	select {
	case err := <-lr.Errors:
		if err != nil {
			if err == io.EOF {
				logrus.Debugf("container stopped: %v", container.ID)
				fmt.Fprintf(c.Response().Writer, "event: container-stopped\ndata: end of stream\n\n")
				f.Flush()
			} else if err != context.Canceled {
				logrus.Errorf("unknown error while streaming %v", err.Error())
			}
		}
	default:
	}

	return nil
}
