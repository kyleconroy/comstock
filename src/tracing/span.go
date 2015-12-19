package tracing

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Span struct {
	TraceID  int64
	ParentID int64
	ID       int64

	annotations map[string]time.Time
	tags        map[string]string
}

func New() Span {
	id := rand.Int63()
	return Span{TraceID: id, ID: id}
}

type encodedSpan struct {
	TraceID  string `json:"root,omitempty"`
	ParentID string `json:"parent,omitempty"`
	ID       string `json:"id,omitempty"`

	Annotations map[string]time.Time `json:"a,omitempty"`
	Tags        map[string]string    `json:"ba,omitempty"`
}

func (s Span) Child() Span {
	return Span{TraceID: s.TraceID, ParentID: s.ID, ID: rand.Int63()}
}

func (s *Span) UnmarshalJSON(in []byte) error {
	var es encodedSpan
	err := json.Unmarshal(in, &es)
	if err != nil {
		return err
	}

	t, terr := strconv.ParseInt(es.TraceID, 16, 64)
	if terr != nil {
		return terr
	}

	p, terr := strconv.ParseInt(es.ParentID, 16, 64)
	if terr != nil {
		return terr
	}

	i, terr := strconv.ParseInt(es.ID, 16, 64)
	if terr != nil {
		return terr
	}

	s.TraceID = t
	s.ParentID = p
	s.ID = i
	s.annotations = es.Annotations
	s.tags = es.Tags
	return nil
}

func (s Span) MarshalJSON() ([]byte, error) {
	es := encodedSpan{
		Annotations: s.annotations,
		Tags:        s.tags,
		TraceID:     fmt.Sprintf("%x", s.TraceID),
		ParentID:    fmt.Sprintf("%x", s.ParentID),
		ID:          fmt.Sprintf("%x", s.ID),
	}
	return json.Marshal(es)
}

// Includes timing
func (s *Span) Annotate(key string) {
	s.annotations[key] = time.Now().UTC()
}

// No timing
func (s *Span) Tag(key, value string) {
	s.tags[key] = value
}
