package stderr

import (
	"log"

	"github.com/utsav2/devexp-monitor/lib/sink"
)

type stderrSink struct{}

func (s *stderrSink) Metric(name string, val float64) {
	log.Printf("%s %f", name, val)
}

func (s *stderrSink) Log(fmt string, val ...interface{}) {
	log.Printf(fmt, val...)
}

func (s *stderrSink) Err(err error) {
	log.Printf("error: %s", err)
}

func New(string) (sink.Sink, error) {
	return &stderrSink{}, nil
}
