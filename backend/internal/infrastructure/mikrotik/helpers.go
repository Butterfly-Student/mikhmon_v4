package mikrotik

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-routeros/routeros/v3/proto"
)

// parseInt parses an integer from a RouterOS string value.
func parseInt(s string) int64 {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// parseFloat parses a float from a RouterOS string value.
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// ParseRate parses a rate string with unit (bps, kbps, Mbps, Gbps) to bits per second.
// Examples: "0bps" -> 0, "74.3kbps" -> 74300, "2.2Mbps" -> 2200000.
func ParseRate(s string) int64 {
	if s == "" || s == "0" {
		return 0
	}

	var value float64
	var unit string

	for i := len(s) - 1; i >= 0; i-- {
		if (s[i] >= '0' && s[i] <= '9') || s[i] == '.' {
			value = parseFloat(s[:i+1])
			unit = s[i+1:]
			break
		}
	}

	if unit == "" {
		return parseInt(s)
	}

	switch strings.ToLower(unit) {
	case "bps":
		return int64(value)
	case "kbps":
		return int64(value * 1_000)
	case "mbps":
		return int64(value * 1_000_000)
	case "gbps":
		return int64(value * 1_000_000_000)
	default:
		return parseInt(s)
	}
}

// parseBool parses a bool from a RouterOS "yes"/"true" string.
func parseBool(s string) bool {
	return s == "true" || s == "yes"
}

// formatInt formats an int64 to string.
func formatInt(n int64) string {
	return strconv.FormatInt(n, 10)
}

// batchDebounce is the silence window used by listenBatches to detect the gap
// between RouterOS interval ticks. RouterOS sends a burst of !re sentences then
// goes silent until the next tick. 200ms safely fits inside any practical interval.
const batchDebounce = 200 * time.Millisecond

// listenBatches reads sentences from a RouterOS follow stream and emits complete
// batches. A batch is considered complete when no new sentence arrives within
// debounce duration — this detects the gap between RouterOS interval ticks.
//
// RouterOS with =follow= =interval=X sends a burst of !re sentences (one per item)
// then is silent until the next tick. The debounce detects that silence.
//
// The returned channel is closed when sentences is closed or ctx is cancelled.
func listenBatches(
	ctx context.Context,
	sentences <-chan *proto.Sentence,
	debounce time.Duration,
) <-chan []*proto.Sentence {
	out := make(chan []*proto.Sentence, 4)

	go func() {
		defer close(out)

		for {
			// Wait for the first sentence of a new batch.
			var first *proto.Sentence
			select {
			case s, ok := <-sentences:
				if !ok {
					return
				}
				first = s
			case <-ctx.Done():
				return
			}

			// Accumulate the rest of the batch until debounce fires.
			batch := []*proto.Sentence{first}
			timer := time.NewTimer(debounce)

		collect:
			for {
				select {
				case s, ok := <-sentences:
					if !ok {
						timer.Stop()
						break collect
					}
					batch = append(batch, s)
					// Reset debounce on each new sentence.
					if !timer.Stop() {
						select {
						case <-timer.C:
						default:
						}
					}
					timer.Reset(debounce)
				case <-timer.C:
					break collect // silence detected — batch is complete
				case <-ctx.Done():
					timer.Stop()
					return
				}
			}

			select {
			case out <- batch:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
