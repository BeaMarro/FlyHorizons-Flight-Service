package load_test_utils

import (
	"fmt"
	"html/template"
	"os"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type LoadTestUtils struct {
}

func (utils *LoadTestUtils) LogMetrics(t *testing.T, metrics *vegeta.Metrics) {
	t.Logf("Requests: %d", metrics.Requests)
	t.Logf("Success Rate: %.2f%%", metrics.Success*100)
	t.Logf("Latency (mean): %s", metrics.Latencies.Mean)
	t.Logf("Latency (50th): %s", metrics.Latencies.P50)
	t.Logf("Latency (95th): %s", metrics.Latencies.P95)
	t.Logf("Latency (99th): %s", metrics.Latencies.P99)
	t.Logf("Throughput: %.2f req/s", metrics.Throughput)
	t.Logf("Errors: %v", metrics.Errors)
}

func (utils *LoadTestUtils) EvaluateMetricsSuccess(t *testing.T, metrics *vegeta.Metrics) {
	if metrics.Success < 0.95 { // Fails if success rate is below 95%
		t.Errorf("Success rate too low: %.2f%% (expected >= 95%%)", metrics.Success*100)
	}
	if len(metrics.Errors) > 0 {
		t.Errorf("Encountered errors: %v", metrics.Errors)
	}
}

func (utils *LoadTestUtils) GenerateHTMLReport(t *testing.T, metrics *vegeta.Metrics, outputFile, title string) {
	// Create output file
	fullPath := fmt.Sprintf("./report/%s", outputFile)

	file, err := os.Create(fullPath)
	if err != nil {
		t.Fatalf("Failed to create HTML report file: %v", err)
	}
	defer file.Close()

	// Format the metrics for the template
	data := map[string]interface{}{
		"Title":       title,
		"Timestamp":   time.Now().Format(time.RFC3339),
		"Requests":    metrics.Requests,
		"SuccessRate": fmt.Sprintf("%.2f%%", metrics.Success*100),
		"LatencyMean": metrics.Latencies.Mean.String(),
		"Latency50th": metrics.Latencies.P50.String(),
		"Latency95th": metrics.Latencies.P95.String(),
		"Latency99th": metrics.Latencies.P99.String(),
		"LatencyMax":  metrics.Latencies.Max.String(),
		"Throughput":  fmt.Sprintf("%.2f", metrics.Throughput),
		"BytesIn":     metrics.BytesIn.Mean,
		"BytesOut":    metrics.BytesOut.Mean,
		"Errors":      metrics.Errors,
		"ErrorRate":   fmt.Sprintf("%.2f%%", (1-metrics.Success)*100),
		"Duration":    metrics.Duration.String(),
		"StatusCodes": metrics.StatusCodes,
	}

	tmpl, err := template.ParseFiles("templates/report_template.html")
	if err != nil {
		t.Fatalf("Failed to parse HTML template: %v", err)
	}

	if err := tmpl.Execute(file, data); err != nil {
		t.Fatalf("Failed to execute HTML template: %v", err)
	}

	t.Logf("HTML report generated successfully: %s", outputFile)
}
