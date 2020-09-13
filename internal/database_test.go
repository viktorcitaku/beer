package internal

import (
	"hash/fnv"
	"io"
	"testing"
	"time"
)

func TestHash(t *testing.T) {
	h1 := fnv.New32a()
	_, _ = io.WriteString(h1, "drake@owl")
	t.Logf("H1: %d", h1.Sum32() + 12)
	h2 := fnv.New32a()
	_, _ = io.WriteString(h2, "fero@gang")
	t.Logf("H2: %d", h2.Sum32() + 12)
}

func TestTicker(t *testing.T) {
	go func() {
		for range time.Tick(time.Second) {
			t.Log(time.Now())
		}
	}()
	time.Sleep(10 * time.Second)
}

