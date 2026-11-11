package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type EventStream struct {
	w       http.ResponseWriter
	flusher http.Flusher
}

func NewEventStream(w http.ResponseWriter) (*EventStream, error) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, errors.New("Streaming unsupported")
	}

	return &EventStream{w: w, flusher: flusher}, nil
}

func (s *EventStream) Start() {
	s.w.Header().Set("Content-Type", "text/event-stream")
	s.w.Header().Set("Cache-Control", "no-cache")
	s.w.Header().Set("X-Accel-Buffering", "no")
}

func (s *EventStream) Send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintf(s.w, "data: %s\n\n", data); err != nil {
		return err
	}
	s.flusher.Flush()

	return nil
}

func (s *EventStream) Ping() error {
	if _, err := fmt.Fprint(s.w, ": ping\n\n"); err != nil {
		return err
	}
	s.flusher.Flush()

	return nil
}
