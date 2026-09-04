package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Config struct {
	URL         string
	Method      string
	Body        string
	Token       string
	Concurrency int
	RPS         float64
	Duration    time.Duration
	Timeout     time.Duration
	Interval    time.Duration

	// Ramp format:
	// 10:30,25:30,50:30,100:60
	// means:
	// 10 concurrency for 30s
	// 25 concurrency for 30s
	// 50 concurrency for 30s
	// 100 concurrency for 60s
	Ramp string

	CSVFile  string
	JSONFile string
}

type Result struct {
	Latency     time.Duration
	StatusCode  int
	Error       error
	Bytes       int64
	StartedAt   time.Time
	CompletedAt time.Time
}

type Metrics struct {
	mu sync.Mutex

	Latencies []time.Duration
	Status    map[int]int64

	Requests int64
	Success  int64
	Failed   int64
	Timeouts int64

	Bytes int64

	InFlight int64

	Start time.Time
}

func NewMetrics() *Metrics {
	return &Metrics{
		Status: make(map[int]int64),
		Start:  time.Now(),
	}
}

func (m *Metrics) Record(r Result) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Requests++

	if r.Error != nil {
		m.Failed++

		if os.IsTimeout(r.Error) {
			m.Timeouts++
		}
	} else if r.StatusCode >= 200 && r.StatusCode < 400 {
		m.Success++
	} else {
		m.Failed++
	}

	if r.StatusCode > 0 {
		m.Status[r.StatusCode]++
	}

	m.Latencies = append(m.Latencies, r.Latency)
	m.Bytes += r.Bytes
}

func (m *Metrics) IncInFlight() {
	m.mu.Lock()
	m.InFlight++
	m.mu.Unlock()
}

func (m *Metrics) DecInFlight() {
	m.mu.Lock()
	m.InFlight--
	m.mu.Unlock()
}

func percentile(values []time.Duration, p float64) time.Duration {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]time.Duration, len(values))
	copy(sorted, values)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	index := int(math.Ceil(p*float64(len(sorted)))) - 1

	if index < 0 {
		index = 0
	}

	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

func (m *Metrics) Snapshot() {
	m.mu.Lock()
	defer m.mu.Unlock()

	elapsed := time.Since(m.Start)

	rps := 0.0
	if elapsed.Seconds() > 0 {
		rps = float64(m.Requests) / elapsed.Seconds()
	}

	var avg time.Duration

	if len(m.Latencies) > 0 {
		var total time.Duration

		for _, l := range m.Latencies {
			total += l
		}

		avg = total / time.Duration(len(m.Latencies))
	}

	var min, max time.Duration

	if len(m.Latencies) > 0 {
		min = m.Latencies[0]
		max = m.Latencies[0]

		for _, l := range m.Latencies {
			if l < min {
				min = l
			}

			if l > max {
				max = l
			}
		}
	}

	fmt.Printf(
		"\r[%s] req=%d success=%d failed=%d in-flight=%d RPS=%.1f "+
			"avg=%s p50=%s p95=%s p99=%s max=%s",
		time.Since(m.Start).Round(time.Second),
		m.Requests,
		m.Success,
		m.Failed,
		m.InFlight,
		rps,
		avg.Round(time.Millisecond),
		percentile(m.Latencies, 0.50).Round(time.Millisecond),
		percentile(m.Latencies, 0.95).Round(time.Millisecond),
		percentile(m.Latencies, 0.99).Round(time.Millisecond),
		max.Round(time.Millisecond),
	)
}

func (m *Metrics) PrintFinal() {
	m.mu.Lock()
	defer m.mu.Unlock()

	elapsed := time.Since(m.Start)

	rps := float64(m.Requests) / elapsed.Seconds()

	var avg time.Duration

	if len(m.Latencies) > 0 {
		var total time.Duration

		for _, l := range m.Latencies {
			total += l
		}

		avg = total / time.Duration(len(m.Latencies))
	}

	var min, max time.Duration

	if len(m.Latencies) > 0 {
		min = m.Latencies[0]
		max = m.Latencies[0]

		for _, l := range m.Latencies {
			if l < min {
				min = l
			}

			if l > max {
				max = l
			}
		}
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("             LOAD TEST RESULT")
	fmt.Println("========================================")

	fmt.Printf("Duration:       %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("Requests:       %d\n", m.Requests)
	fmt.Printf("Successful:     %d (%.2f%%)\n",
		m.Success,
		percentage(m.Success, m.Requests),
	)
	fmt.Printf("Failed:         %d (%.2f%%)\n",
		m.Failed,
		percentage(m.Failed, m.Requests),
	)
	fmt.Printf("Timeouts:       %d\n", m.Timeouts)
	fmt.Printf("Actual RPS:     %.2f\n", rps)
	fmt.Printf("Bytes received: %d\n", m.Bytes)

	fmt.Println()
	fmt.Println("Latency:")
	fmt.Printf("  Min:          %s\n", min.Round(time.Millisecond))
	fmt.Printf("  Avg:          %s\n", avg.Round(time.Millisecond))
	fmt.Printf("  P50:          %s\n", percentile(m.Latencies, 0.50).Round(time.Millisecond))
	fmt.Printf("  P90:          %s\n", percentile(m.Latencies, 0.90).Round(time.Millisecond))
	fmt.Printf("  P95:          %s\n", percentile(m.Latencies, 0.95).Round(time.Millisecond))
	fmt.Printf("  P99:          %s\n", percentile(m.Latencies, 0.99).Round(time.Millisecond))
	fmt.Printf("  Max:          %s\n", max.Round(time.Millisecond))

	fmt.Println()
	fmt.Println("HTTP Status Codes:")

	statusCodes := make([]int, 0, len(m.Status))

	for code := range m.Status {
		statusCodes = append(statusCodes, code)
	}

	sort.Ints(statusCodes)

	for _, code := range statusCodes {
		fmt.Printf("  %d:            %d\n", code, m.Status[code])
	}

	fmt.Println("========================================")
}

func percentage(value, total int64) float64 {
	if total == 0 {
		return 0
	}

	return float64(value) / float64(total) * 100
}

func parseHeaders() http.Header {
	headers := http.Header{}

	raw := os.Getenv("LOADTEST_HEADERS")

	if raw == "" {
		return headers
	}

	// Example:
	// LOADTEST_HEADERS="X-Test-Run=abc,Accept=application/json"
	for _, item := range strings.Split(raw, ",") {
		parts := strings.SplitN(item, "=", 2)

		if len(parts) == 2 {
			headers.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}

	return headers
}

func doRequest(
	client *http.Client,
	cfg Config,
	headers http.Header,
	metrics *Metrics,
) {
	start := time.Now()

	metrics.IncInFlight()
	defer metrics.DecInFlight()

	var body io.Reader

	if cfg.Body != "" {
		body = bytes.NewBufferString(cfg.Body)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.Timeout,
	)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		cfg.Method,
		cfg.URL,
		body,
	)

	if err != nil {
		metrics.Record(Result{
			Latency: time.Since(start),
			Error:   err,
		})
		return
	}

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	if cfg.Token != "" {
		req.Header.Set(
			"Authorization",
			"Bearer "+cfg.Token,
		)
	}

	if cfg.Body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := client.Do(req)

	if err != nil {
		metrics.Record(Result{
			Latency: time.Since(start),
			Error:   err,
		})
		return
	}

	defer resp.Body.Close()

	n, _ := io.Copy(io.Discard, resp.Body)

	metrics.Record(Result{
		Latency:    time.Since(start),
		StatusCode: resp.StatusCode,
		Bytes:      n,
	})
}

type RampStage struct {
	Concurrency int
	Duration    time.Duration
}

func parseRamp(input string) ([]RampStage, error) {
	if input == "" {
		return nil, nil
	}

	var stages []RampStage

	for _, stage := range strings.Split(input, ",") {
		parts := strings.Split(stage, ":")

		if len(parts) != 2 {
			return nil, fmt.Errorf(
				"invalid ramp stage %q; expected concurrency:seconds",
				stage,
			)
		}

		concurrency, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, err
		}

		seconds, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		stages = append(stages, RampStage{
			Concurrency: concurrency,
			Duration:    time.Duration(seconds) * time.Second,
		})
	}

	return stages, nil
}

func runStage(
	client *http.Client,
	cfg Config,
	metrics *Metrics,
	concurrency int,
	duration time.Duration,
	headers http.Header,
) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		duration,
	)
	defer cancel()

	var wg sync.WaitGroup

	sem := make(chan struct{}, concurrency)

	var interval <-chan time.Time

	var ticker *time.Ticker

	if cfg.RPS > 0 {
		intervalDuration := time.Duration(
			float64(time.Second) / cfg.RPS,
		)

		ticker = time.NewTicker(intervalDuration)
		defer ticker.Stop()

		interval = ticker.C
	}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return

		default:
		}

		if interval != nil {
			select {
			case <-ctx.Done():
				wg.Wait()
				return

			case <-interval:
			}
		}

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Wait()
			return
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			defer func() {
				<-sem
			}()

			doRequest(
				client,
				cfg,
				headers,
				metrics,
			)
		}()
	}
}

func writeJSON(filename string, metrics *Metrics) error {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	result := map[string]interface{}{
		"requests":   metrics.Requests,
		"successful": metrics.Success,
		"failed":     metrics.Failed,
		"timeouts":   metrics.Timeouts,
		"bytes":      metrics.Bytes,
		"status":     metrics.Status,
		"latency": map[string]string{
			"min": percentile(metrics.Latencies, 0).String(),
			"p50": percentile(metrics.Latencies, 0.50).String(),
			"p90": percentile(metrics.Latencies, 0.90).String(),
			"p95": percentile(metrics.Latencies, 0.95).String(),
			"p99": percentile(metrics.Latencies, 0.99).String(),
			"max": percentile(metrics.Latencies, 1).String(),
		},
	}

	data, err := json.MarshalIndent(result, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func main() {
	var cfg Config

	flag.StringVar(&cfg.URL, "url", "", "Target URL")
	flag.StringVar(&cfg.Method, "method", "POST", "HTTP method")
	flag.StringVar(&cfg.Body, "body", "", "JSON request body")
	flag.StringVar(&cfg.Token, "token", "", "Bearer token")

	flag.IntVar(
		&cfg.Concurrency,
		"concurrency",
		10,
		"Maximum concurrent requests",
	)

	flag.Float64Var(
		&cfg.RPS,
		"rps",
		0,
		"Target requests/sec; 0 = as fast as concurrency allows",
	)

	flag.DurationVar(
		&cfg.Duration,
		"duration",
		60*time.Second,
		"Test duration",
	)

	flag.DurationVar(
		&cfg.Timeout,
		"timeout",
		10*time.Second,
		"Request timeout",
	)

	flag.DurationVar(
		&cfg.Interval,
		"interval",
		5*time.Second,
		"Live metrics interval",
	)

	flag.StringVar(
		&cfg.Ramp,
		"ramp",
		"",
		"Ramp stages, e.g. 10:30,25:30,50:30,100:60",
	)

	flag.StringVar(
		&cfg.CSVFile,
		"csv",
		"",
		"CSV output file",
	)

	flag.StringVar(
		&cfg.JSONFile,
		"json",
		"",
		"JSON summary output file",
	)

	flag.Parse()

	if cfg.URL == "" {
		log.Fatal("missing -url")
	}

	headers := parseHeaders()

	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
		MaxConnsPerHost:     cfg.Concurrency,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
	}

	metrics := NewMetrics()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("              GO LOAD TEST")
	fmt.Println("========================================")
	fmt.Printf("Target:         %s\n", cfg.URL)
	fmt.Printf("Method:         %s\n", cfg.Method)
	fmt.Printf("Concurrency:    %d\n", cfg.Concurrency)
	fmt.Printf("RPS:            %.2f\n", cfg.RPS)
	fmt.Printf("Duration:       %s\n", cfg.Duration)
	fmt.Printf("Timeout:        %s\n", cfg.Timeout)
	fmt.Println("========================================")
	fmt.Println()

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// Live metrics.
	go func() {
		ticker := time.NewTicker(cfg.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				metrics.Snapshot()

			case <-ctx.Done():
				return
			}
		}
	}()

	ramp, err := parseRamp(cfg.Ramp)

	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	if len(ramp) > 0 {
		for _, stage := range ramp {
			fmt.Printf(
				"\nStarting stage: concurrency=%d duration=%s\n",
				stage.Concurrency,
				stage.Duration,
			)

			stageCfg := cfg
			stageCfg.Concurrency = stage.Concurrency

			runStage(
				client,
				stageCfg,
				metrics,
				stage.Concurrency,
				stage.Duration,
				headers,
			)

			select {
			case <-ctx.Done():
				break

			default:
			}
		}
	} else {
		runStage(
			client,
			cfg,
			metrics,
			cfg.Concurrency,
			cfg.Duration,
			headers,
		)
	}

	fmt.Printf(
		"\nTest completed in %s\n",
		time.Since(start).Round(time.Millisecond),
	)

	metrics.PrintFinal()

	if cfg.JSONFile != "" {
		if err := writeJSON(cfg.JSONFile, metrics); err != nil {
			log.Printf("failed to write JSON: %v", err)
		} else {
			fmt.Printf("JSON written to %s\n", cfg.JSONFile)
		}
	}

	if cfg.CSVFile != "" {
		if err := writeCSV(cfg.CSVFile, metrics); err != nil {
			log.Printf("failed to write CSV: %v", err)
		} else {
			fmt.Printf("CSV written to %s\n", cfg.CSVFile)
		}
	}
}

func writeCSV(filename string, metrics *Metrics) error {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{
		"latency_ms",
	}); err != nil {
		return err
	}

	for _, latency := range metrics.Latencies {
		if err := w.Write([]string{
			fmt.Sprintf("%.3f", float64(latency.Microseconds())/1000),
		}); err != nil {
			return err
		}
	}

	return w.Error()
}
