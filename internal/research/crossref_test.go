package research

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCrossrefClientFetchesDOIMetadata(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/works/10.1103/PhysRevC.106.L031301" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("User-Agent"); got == "" {
			t.Fatal("expected User-Agent header")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"status": "ok",
			"message": {
				"DOI": "10.1103/PhysRevC.106.L031301",
				"title": ["First experiment at the Super Heavy Element Factory"],
				"container-title": ["Physical Review C"],
				"issued": {"date-parts": [[2022, 9, 12]]}
			}
		}`))
	}))
	defer server.Close()

	client := NewCrossrefClient(WithBaseURL(server.URL), WithUserAgent("moscovium-statera-go-test"))
	work, err := client.FetchWork(context.Background(), "10.1103/PhysRevC.106.L031301")
	if err != nil {
		t.Fatalf("FetchWork returned error: %v", err)
	}

	if work.DOI != "10.1103/PhysRevC.106.L031301" {
		t.Fatalf("DOI = %q", work.DOI)
	}
	if work.Title != "First experiment at the Super Heavy Element Factory" {
		t.Fatalf("Title = %q", work.Title)
	}
	if work.Journal != "Physical Review C" {
		t.Fatalf("Journal = %q", work.Journal)
	}
	if work.Year != 2022 {
		t.Fatalf("Year = %d", work.Year)
	}
}
