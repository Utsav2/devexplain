package sink

type Sink interface {
	Metric(string, float64)
	Log(string, ...interface{})
	Err(error)
}
