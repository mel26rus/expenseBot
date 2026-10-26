package exchangerate

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
func TestFetchRates(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })
	for _, tc := range []struct {
		name, body string
		status     int
		valid      bool
	}{
		{"valid", `{"timestamp":1790899200,"base":"RUB","rates":{"USD":0.01,"EUR":0.009}}`, 200, true},
		{"HTTP error", ``, 503, false},
		{"invalid JSON", `{`, 200, false},
		{"missing USD", `{"rates":{"EUR":0.009}}`, 200, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			http.DefaultTransport = transportFunc(func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body))}, nil
			})
			rates := FetchRates(context.Background())
			if (rates != nil) != tc.valid {
				t.Fatalf("rates=%v", rates)
			}
			if tc.valid && (rates.Values["USDT"] != 1 || rates.Values["USD"] != 1 || rates.Values["RUB"] != 100 || rates.Values["EUR"] < 0.899 || rates.Values["EUR"] > 0.901 || !rates.Date.Equal(time.Unix(1790899200, 0).UTC())) {
				t.Fatalf("incorrect conversion: %+v", rates)
			}
		})
	}
}
