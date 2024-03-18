package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mehranzand/pulseup/api/middleware"
	"github.com/mehranzand/pulseup/internal/docker"
)

// GetContainers
// @Summary Get list of containers
func (h *Handler) GetContainers(c echo.Context) error {
	cc := c.(*middleware.DockerContext)

	if cc.Client == nil {
		http.Error(c.Response().Writer, "Docker host not found!", http.StatusInternalServerError)

		return nil
	}

	containers, err := cc.Client.ListContainers()
	if err != nil {
		http.Error(c.Response().Writer, "Continer not found!", http.StatusInternalServerError)

		return nil
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

	_, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		http.Error(c.Response().Writer, "Streaming unsupported!", http.StatusInternalServerError)
	}

	since := time.Now().AddDate(0, 0, -2)

	reader, err := cc.Client.ContainerLogs(c.Request().Context(), id, strconv.FormatInt(since.Unix(), 10), stdTypes)
	if err != nil {
		http.Error(c.Response().Writer, "Continer not found!", http.StatusInternalServerError)
	}

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-transform")
	c.Response().Header().Add("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")

	defer reader.Close()

	data := make([]byte, 10000)
	for {
		bufReader := bufio.NewReader(reader)
		_, err = bufReader.Read(data)

		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println(string(data))
	}

	return nil
}
