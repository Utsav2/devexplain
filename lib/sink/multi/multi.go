package multi

import (
	"fmt"
	"strings"

	"github.com/utsav2/devexplain/lib/sink"
	"github.com/utsav2/devexplain/lib/sink/datadog"
	"github.com/utsav2/devexplain/lib/sink/stderr"
)

type multiSink struct {
	sinks []sink.Sink
}

func (m *multiSink) Metric(name string, val float64) {
	for _, s := range m.sinks {
		s.Metric(name, val)
	}
}

func (m *multiSink) Log(fmt string, val ...interface{}) {
	for _, s := range m.sinks {
		s.Log(fmt, val...)
	}
}

func (m *multiSink) Err(err error) {
	for _, s := range m.sinks {
		s.Err(err)
	}
}

var validSinks = map[string]func(string) (sink.Sink, error) {
	"stderr":  stderr.New,
	"datadog": datadog.New,
}

var MultiSink sink.Sink = &multiSink{}

func New(progname string, sinkNames []string) (sink.Sink, error) {
	sinks := []sink.Sink{}
	for _, s := range sinkNames {
		fn, ok := validSinks[s]
		if !ok {
			sinks := []string{}
			for vs := range validSinks {
				sinks = append(sinks, vs)
			}
			return nil, fmt.Errorf("invalid sink %s. Valid sinks: %s", s, strings.Join(sinks, ","))
		}
		sink, err := fn(progname)
		if err != nil {
			return nil, fmt.Errorf("error during sink creation: %s", err)
		}
		sinks = append(sinks, sink)
	}
	return &multiSink{
		sinks: sinks,
	}, nil
}
