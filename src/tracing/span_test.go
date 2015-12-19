package tracing

import (
	"encoding/json"
	"testing"
)

func TestSpan(t *testing.T) {
	span := New()
	blob, err := json.Marshal(span)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(blob))

	var s Span
	err = json.Unmarshal(blob, &s)
	if err != nil {
		t.Fatal(err)
	}

	if s.TraceID != span.TraceID {
		t.Errorf("incorrect span serialization: %+v", s)
	}
}
