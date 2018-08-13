package tracer

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
)

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
//func (z *span) Msgsize() (s int) {
//s = 1 + 5 + msgp.StringPrefixSize + len(z.Name) + 8 + msgp.StringPrefixSize + len(z.Service) + 9 + msgp.StringPrefixSize + len(z.Resource) + 5 + msgp.StringPrefixSize + len(z.Type) + 6 + msgp.Int64Size + 9 + msgp.Int64Size + 5 + msgp.MapHeaderSize
//if z.Meta != nil {
//for za0001, za0002 := range z.Meta {
//_ = za0002
//s += msgp.StringPrefixSize + len(za0001) + msgp.StringPrefixSize + len(za0002)
//}
//}
//s += 8 + msgp.MapHeaderSize
//if z.Metrics != nil {
//for za0003, za0004 := range z.Metrics {
//_ = za0004
//s += msgp.StringPrefixSize + len(za0003) + msgp.Float64Size
//}
//}
//s += 8 + msgp.Uint64Size + 9 + msgp.Uint64Size + 10 + msgp.Uint64Size + 6 + msgp.Int32Size
//return
//}

// span represents a computation. Callers must call Finish when a span is
// complete to ensure it's submitted.
type span struct {
	sync.RWMutex `msg:"-"`

	Name     string             `msg:"name"`              // operation name
	Service  string             `msg:"service"`           // service name (i.e. "grpc.server", "http.request")
	Resource string             `msg:"resource"`          // resource name (i.e. "/user?id=123", "SELECT * FROM users")
	Type     string             `msg:"type"`              // protocol associated with the span (i.e. "web", "db", "cache")
	Start    int64              `msg:"start"`             // span start time expressed in nanoseconds since epoch
	Duration int64              `msg:"duration"`          // duration of the span expressed in nanoseconds
	Meta     map[string]string  `msg:"meta,omitempty"`    // arbitrary map of metadata
	Metrics  map[string]float64 `msg:"metrics,omitempty"` // arbitrary map of numeric metrics
	SpanID   uint64             `msg:"span_id"`           // identifier of this span
	TraceID  uint64             `msg:"trace_id"`          // identifier of the root span
	ParentID uint64             `msg:"parent_id"`         // identifier of the span's direct parent
	Error    int32              `msg:"error"`             // error status of the span; 0 means no errors

	context *spanContext `msg:"-"` // span propagation context
}

// Context yields the SpanContext for this Span. Note that the return
// value of Context() is still valid after a call to Finish(). This is
// called the span context and it is different from Go's context.
func (s *span) Context() ddtrace.SpanContext { return s.context }

// SetBaggageItem sets a key/value pair as baggage on the span. Baggage items
// are propagated down to descendant spans and injected cross-process. Use with
// care as it adds extra load onto your tracing layer.
func (s *span) SetBaggageItem(key, val string) {
	s.context.setBaggageItem(key, val)
}

// BaggageItem gets the value for a baggage item given its key. Returns the
// empty string if the value isn't found in this Span.
func (s *span) BaggageItem(key string) string {
	return s.context.baggageItem(key)
}

// SetTag adds a set of key/value metadata to the span.
func (s *span) SetTag(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()
	if key == ext.Error {
		s.setTagError(value)
		return
	}
	if v, ok := value.(string); ok {
		s.setTagString(key, v)
		return
	}
	if v, ok := toFloat64(value); ok {
		s.setTagNumeric(key, v)
		return
	}
	// not numeric, not a string and not an error, the likelihood of this
	// happening is close to zero, but we should nevertheless account for it.
	s.Meta[key] = fmt.Sprint(value)
}

// setTagError sets the error tag. It accounts for various valid scenarios.
// This method is not safe for concurrent use.
func (s *span) setTagError(value interface{}) {
	switch v := value.(type) {
	case bool:
		// bool value as per Opentracing spec.
		if !v {
			s.Error = 0
		} else {
			s.Error = 1
		}
	case error:
		// if anyone sets an error value as the tag, be nice here
		// and provide all the benefits.
		s.Error = 1
		s.Meta[ext.ErrorMsg] = v.Error()
		s.Meta[ext.ErrorType] = reflect.TypeOf(v).String()
		s.Meta[ext.ErrorStack] = string(debug.Stack())
	case nil:
		// no error
		s.Error = 0
	default:
		// in all other cases, let's assume that setting this tag
		// is the result of an error.
		s.Error = 1
	}
}

// setTagString sets a string tag. This method is not safe for concurrent use.
func (s *span) setTagString(key, v string) {
	switch key {
	case ext.ServiceName:
		s.Service = v
	case ext.ResourceName:
		s.Resource = v
	case ext.SpanType:
		s.Type = v
	default:
		s.Meta[key] = v
	}
}

// setTagNumeric sets a numeric tag, in our case called a metric. This method
// is not safe for concurrent use.
func (s *span) setTagNumeric(key string, v float64) {
	switch key {
	case ext.SamplingPriority:
		// setting sampling priority per spec
		s.Metrics[samplingPriorityKey] = v
		s.context.setSamplingPriority(int(v))
	default:
		s.Metrics[key] = v
	}
}

// Finish closes this Span (but not its children) providing the duration
// of its part of the tracing session.
func (s *span) Finish(opts ...ddtrace.FinishOption) {
	var cfg ddtrace.FinishConfig
	for _, fn := range opts {
		fn(&cfg)
	}
	var t int64
	if cfg.FinishTime.IsZero() {
		t = now()
	} else {
		t = cfg.FinishTime.UnixNano()
	}
	if cfg.Error != nil {
		s.SetTag(ext.Error, cfg.Error)
	}
	s.finish(t)
}

// SetOperationName sets or changes the operation name.
func (s *span) SetOperationName(operationName string) {
	s.Lock()
	defer s.Unlock()

	s.Name = operationName
}

func (s *span) finish(finishTime int64) {
	s.Lock()
	defer s.Unlock()
	if s.Duration == 0 {
		s.Duration = finishTime - s.Start
	}
	if !s.context.sampled {
		// not sampled
		return
	}
	// todo
}

// String returns a human readable representation of the span. Not for
// production, just debugging.
func (s *span) String() string {
	lines := []string{
		fmt.Sprintf("Name: %s", s.Name),
		fmt.Sprintf("Service: %s", s.Service),
		fmt.Sprintf("Resource: %s", s.Resource),
		fmt.Sprintf("TraceID: %d", s.TraceID),
		fmt.Sprintf("SpanID: %d", s.SpanID),
		fmt.Sprintf("ParentID: %d", s.ParentID),
		fmt.Sprintf("Start: %s", time.Unix(0, s.Start)),
		fmt.Sprintf("Duration: %s", time.Duration(s.Duration)),
		fmt.Sprintf("Error: %d", s.Error),
		fmt.Sprintf("Type: %s", s.Type),
		"Tags:",
	}
	s.RLock()
	for key, val := range s.Meta {
		lines = append(lines, fmt.Sprintf("\t%s:%s", key, val))
	}
	for key, val := range s.Metrics {
		lines = append(lines, fmt.Sprintf("\t%s:%f", key, val))
	}
	s.RUnlock()
	return strings.Join(lines, "\n")
}

const samplingPriorityKey = "_sampling_priority_v1"
