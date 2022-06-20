package datadog

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	datadog "github.com/DataDog/datadog-api-client-go/api/v2/datadog"
	"github.com/utsav2/devexplain/lib/sink"
)

type datadogSink struct {
	progname string
	apiKey    string
	apiClient *datadog.APIClient
}

func (s *datadogSink) context() context.Context {
	return context.WithValue(context.Background(), datadog.ContextAPIKeys, map[string]datadog.APIKey{
		"apiKeyAuth": {
			Key: s.apiKey,
		},
	})
}

func ptrInt(s int64) *int64 {
	return &s
}

func ptrFloat(f float64) *float64 {
	return &f
}

func (s *datadogSink) Metric(name string, val float64) {
	accepted, _, err := s.apiClient.MetricsApi.SubmitMetrics(
		s.context(),
		*datadog.NewMetricPayload([]datadog.MetricSeries{
			{
				Metric: fmt.Sprintf("%s.%s", s.progname, name),
				Points: []datadog.MetricPoint{
					{
						Timestamp: ptrInt(time.Now().Unix()),
						Value:     ptrFloat(val),
					},
				},
			},
		}))
	if err != nil {
		log.Printf("error sending datadog metrics: %s\n", err)
	} else {
		status, ok := accepted.GetStatusOk()
		log.Printf("datadog metrics submission status: %v %v\n", status, ok)
	}
}

func (s *datadogSink) Log(fmtstr string, val ...interface{}) {
	_, resp, err := s.apiClient.LogsApi.SubmitLog(s.context(), []datadog.HTTPLogItem{
		*datadog.NewHTTPLogItem(fmt.Sprintf(fmtstr, val...)),
	})
	if err != nil {
		log.Printf("error sending datadog logs : %s\n", err)
	} else {
		log.Printf("datadog logs submission HTTP status: %v\n", resp.StatusCode)
	}
}

func (s *datadogSink) Err(err error) {
	s.Log("error: %s", err)
}

var envVar = "DATADOG_API_KEY"

func New(progname string) (sink.Sink, error) {
	key, ok := os.LookupEnv(envVar)
	if !ok {
		return nil, fmt.Errorf("datadog sink: API key not provided. Use %s environment variable to provide", envVar)
	}
	return &datadogSink{
		progname: progname,
		apiClient: datadog.NewAPIClient(datadog.NewConfiguration()),
		apiKey:    key,
	}, nil
}
