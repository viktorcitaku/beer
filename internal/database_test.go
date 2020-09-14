package internal

import (
	"hash/fnv"
	"io"
	"math"
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

func TestExecuteDatabaseCleanup(t *testing.T) {

	t.SkipNow() // This test is ignored

	InitializeDatabase()

	u1 := &User{
		Email:       "drake@owl",
		LastUpdated: time.Now().Add(-10 * time.Minute),
	}

	u2 := &User{
		Email:       "fero@gang",
		LastUpdated: time.Now().Add(-10 * time.Minute),
	}

	SaveUser(u1)
	SaveUser(u2)

	t.Log(math.Abs(u1.LastUpdated.Sub(time.Now()).Minutes()) > 5)

	ExecuteDatabaseCleanup()
}

