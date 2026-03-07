package httpstub

import (
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestRoundTripFuncDelegatesToWrappedFunc(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest returned error: %v", err)
	}

	wantErr := errors.New("boom")
	transport := RoundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r != req {
			t.Fatalf("unexpected request pointer: %p", r)
		}
		return nil, wantErr
	})

	_, gotErr := transport.RoundTrip(req)
	if !errors.Is(gotErr, wantErr) {
		t.Fatalf("expected wrapped error %v, got %v", wantErr, gotErr)
	}
}

func TestJSONResponseBuildsJSONHTTPResponse(t *testing.T) {
	t.Parallel()

	resp := JSONResponse(http.StatusCreated, `{"ok":true}`)

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("unexpected content type: %q", got)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("ReadAll returned error: %v", err)
	}
	if string(body) != `{"ok":true}` {
		t.Fatalf("unexpected body: %q", string(body))
	}
}
