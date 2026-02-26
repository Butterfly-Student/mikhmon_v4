package mikrotik

import (
	"testing"

	"github.com/go-routeros/routeros/v3/proto"
)

func TestParseRate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"0bps", "0bps", 0},
		{"100bps", "100bps", 100},
		{"1kbps", "1kbps", 1000},
		{"74.3kbps", "74.3kbps", 74300},
		{"1Mbps", "1Mbps", 1000000},
		{"2.2Mbps", "2.2Mbps", 2200000},
		{"1Gbps", "1Gbps", 1000000000},
		{"0.5Gbps", "0.5Gbps", 500000000},
		{"empty", "", 0},
		{"zero", "0", 0},
		{"plain number", "12345", 12345},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseRate(tt.input)
			if result != tt.expected {
				t.Errorf("parseRate(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSplitRateValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantIn   int64
		wantOut  int64
	}{
		{
			name:    "rate dengan unit",
			input:   "74.3kbps/2.2Mbps",
			wantIn:  74300,
			wantOut: 2200000,
		},
		{
			name:    "mixed rate dan plain",
			input:   "1000/2.5Mbps",
			wantIn:  1000,
			wantOut: 2500000,
		},
		{
			name:    "zero rates",
			input:   "0bps/0bps",
			wantIn:  0,
			wantOut: 0,
		},
		{
			name:    "empty",
			input:   "",
			wantIn:  0,
			wantOut: 0,
		},
		{
			name:    "single value (no slash)",
			input:   "1.5Mbps",
			wantIn:  1500000,
			wantOut: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIn, gotOut := splitRateValue(tt.input)
			if gotIn != tt.wantIn {
				t.Errorf("splitRateValue() gotIn = %v, want %v", gotIn, tt.wantIn)
			}
			if gotOut != tt.wantOut {
				t.Errorf("splitRateValue() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestSplitSlashValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantIn   int64
		wantOut  int64
	}{
		{
			name:    "bytes dengan nilai besar",
			input:   "11519824317/96401664078",
			wantIn:  11519824317,
			wantOut: 96401664078,
		},
		{
			name:    "packets dengan nilai besar",
			input:   "51272188/83514119",
			wantIn:  51272188,
			wantOut: 83514119,
		},
		{
			name:    "zero values",
			input:   "0/0",
			wantIn:  0,
			wantOut: 0,
		},
		{
			name:    "empty string",
			input:   "",
			wantIn:  0,
			wantOut: 0,
		},
		{
			name:    "single value (no slash)",
			input:   "12345",
			wantIn:  12345,
			wantOut: 0,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIn, gotOut := splitSlashValue(tt.input)
			if gotIn != tt.wantIn {
				t.Errorf("splitSlashValue() gotIn = %v, want %v", gotIn, tt.wantIn)
			}
			if gotOut != tt.wantOut {
				t.Errorf("splitSlashValue() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestParseQueueStatsSentence(t *testing.T) {
	// Simulasi proto.Sentence dengan Map seperti dari MikroTik
	// proto.Sentence memiliki field Map bertipe map[string]string
	sentence := &proto.Sentence{
		Map: map[string]string{
			"bytes":                 "11519824317/96401664078",
			"packets":               "51272188/83514119",
			"queued-bytes":          "0/0",
			"queued-packets":        "0/0",
			"dropped":               "0/0",
			"rate":                  "74/2200", // simplified, tanpa unit
			"packet-rate":           "97/217",
			"total-bytes":           "0",
			"total-packets":         "0",
			"total-queued-bytes":    "0",
			"total-queued-packets":  "0",
			"total-dropped":         "0",
			"total-rate":            "0",
			"total-packet-rate":     "0",
		},
	}

	result := parseQueueStatsSentence(sentence, "TestQueue")

	if result.Name != "TestQueue" {
		t.Errorf("Expected Name = TestQueue, got %s", result.Name)
	}
	if result.BytesIn != 11519824317 {
		t.Errorf("Expected BytesIn = 11519824317, got %d", result.BytesIn)
	}
	if result.BytesOut != 96401664078 {
		t.Errorf("Expected BytesOut = 96401664078, got %d", result.BytesOut)
	}
	if result.PacketsIn != 51272188 {
		t.Errorf("Expected PacketsIn = 51272188, got %d", result.PacketsIn)
	}
	if result.PacketsOut != 83514119 {
		t.Errorf("Expected PacketsOut = 83514119, got %d", result.PacketsOut)
	}
	if result.TotalBytes != 0 {
		t.Errorf("Expected TotalBytes = 0, got %d", result.TotalBytes)
	}
	if result.TotalPackets != 0 {
		t.Errorf("Expected TotalPackets = 0, got %d", result.TotalPackets)
	}
	// Test Queued fields
	if result.QueuedBytesIn != 0 {
		t.Errorf("Expected QueuedBytesIn = 0, got %d", result.QueuedBytesIn)
	}
	if result.QueuedBytesOut != 0 {
		t.Errorf("Expected QueuedBytesOut = 0, got %d", result.QueuedBytesOut)
	}
	// Test Rate fields
	if result.RateIn != 74 {
		t.Errorf("Expected RateIn = 74, got %d", result.RateIn)
	}
	if result.RateOut != 2200 {
		t.Errorf("Expected RateOut = 2200, got %d", result.RateOut)
	}
}
