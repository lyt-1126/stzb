package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestStaticRoute_DistAndTeamFallback(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	RegisterRoute(r)

	t.Run("team data.html", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/data.html", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected %d got %d body=%s", http.StatusOK, rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), "率土Data") {
			t.Fatalf("expected html to contain title, got body=%s", rec.Body.String())
		}
	})

	t.Run("team assets fallback", func(t *testing.T) {
		for _, p := range []string{"/assets/index-CTHhHu0d.js", "/assets/index-fk8CytI0.css"} {
			req := httptest.NewRequest(http.MethodGet, p, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("expected %d got %d for %s body=%s", http.StatusOK, rec.Code, p, rec.Body.String())
			}
			if rec.Body.Len() == 0 {
				t.Fatalf("expected non-empty body for %s", p)
			}
		}
	})

	t.Run("dist index", func(t *testing.T) {
		assertHTML200 := func(path string) {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Fatalf("expected %d got %d for %s body=%s", http.StatusOK, rec.Code, path, rec.Body.String())
			}
			bodyLower := strings.ToLower(rec.Body.String())
			if !strings.Contains(bodyLower, "<html") {
				t.Fatalf("expected html for %s, got body=%s", path, rec.Body.String())
			}
		}

		assertHTML200("/")

		path := "/index.html"
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			if rec.Code == http.StatusOK {
				bodyLower := strings.ToLower(rec.Body.String())
				if !strings.Contains(bodyLower, "<html") {
					t.Fatalf("expected html for %s, got body=%s", path, rec.Body.String())
				}
				return
			}
			if rec.Code != http.StatusMovedPermanently && rec.Code != http.StatusFound {
				t.Fatalf("expected 200/30x got %d for %s body=%s", rec.Code, path, rec.Body.String())
			}
			loc := rec.Header().Get("Location")
			if loc == "" {
				t.Fatalf("redirect without Location for %s", path)
			}
			base := &url.URL{Scheme: "http", Host: "example.com", Path: path}
			ref, err := url.Parse(loc)
			if err != nil {
				t.Fatalf("invalid redirect Location %q for %s: %v", loc, path, err)
			}
			path = base.ResolveReference(ref).Path
			if path == "" {
				path = "/"
			}
		}
		t.Fatalf("too many redirects for /index.html")
	})

	t.Run("404 json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/__not_exist__", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected %d got %d body=%s", http.StatusNotFound, rec.Code, rec.Body.String())
		}
		if !strings.Contains(rec.Body.String(), "404 - Page Not Found") {
			t.Fatalf("expected json message, got body=%s", rec.Body.String())
		}
	})
}

