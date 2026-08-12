package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lissy93/who-dat/internal/config"
	"github.com/lissy93/who-dat/internal/domain"
	"github.com/lissy93/who-dat/internal/lookup"
	"github.com/lissy93/who-dat/internal/model"
	"github.com/lissy93/who-dat/internal/srcerr"
)

type fakeSource struct {
	result *model.Result
	err    error
}

func (f *fakeSource) Lookup(context.Context, domain.Name) (*model.Result, error) {
	return f.result, f.err
}

// newTestServer builds a handler whose RDAP source returns the given result/err. The WHOIS
// source always reports no source, so it never masks the RDAP outcome under test.
func newTestServer(result *model.Result, err error) http.Handler {
	cfg := &config.Config{LookupTimeout: 2_000_000_000, MaxDomains: 5} // 2s, no rate limit/auth
	svc := lookup.NewService(&fakeSource{result: result, err: err}, &fakeSource{err: srcerr.ErrNoSource}, nil)
	return New(cfg, svc)
}

func registered() *model.Result {
	r := model.New("example.com", "example.com", "com")
	r.IsRegistered = true
	r.Meta = model.Meta{Source: model.SourceRDAP}
	r.Raw = []byte("RAW-RDAP-BODY")
	r.RawContentType = "application/rdap+json"
	return r
}

func TestLookupStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		result     *model.Result
		err        error
		wantStatus int
		wantCode   string
	}{
		{"ok", "/v1/whois/example.com", registered(), nil, http.StatusOK, ""},
		{"invalid domain", "/v1/whois/bad_label.com", nil, nil, http.StatusBadRequest, codeInvalidDomain},
		{"upstream", "/v1/whois/example.com", nil, srcerr.ErrUpstream, http.StatusBadGateway, codeUpstreamError},
		{"timeout", "/v1/whois/example.com", nil, srcerr.ErrTimeout, http.StatusGatewayTimeout, codeUpstreamTimout},
		{"unsupported", "/v1/whois/example.com", nil, srcerr.ErrNoSource, http.StatusNotImplemented, codeUnsupportedTLD},
		{"rate limited", "/v1/whois/example.com", nil, &srcerr.RateLimited{}, http.StatusTooManyRequests, codeRateLimited},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			newTestServer(tt.result, tt.err).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tt.path, nil))

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body %s)", rec.Code, tt.wantStatus, rec.Body)
			}
			if tt.wantCode != "" {
				var body errorBody
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
					t.Fatalf("decode error body: %v", err)
				}
				if body.Error.Code != tt.wantCode {
					t.Errorf("error code = %q, want %q", body.Error.Code, tt.wantCode)
				}
			}
		})
	}
}

func TestLookupErrorEnvelope(t *testing.T) {
	tests := []struct {
		name                             string
		rdapErr, whoisErr                error
		wantStatus                       int
		wantCode, wantSource, wantServer string
		wantDetail                       string // substring; empty means detail must be empty
	}{
		{
			name:     "source and server surface",
			rdapErr:  &srcerr.SourceError{Source: model.SourceRDAP, Server: "https://rdap.nic.li", Err: fmt.Errorf("%w: no response", srcerr.ErrTimeout)},
			whoisErr: srcerr.ErrNoSource, wantStatus: http.StatusGatewayTimeout, wantCode: codeUpstreamTimout,
			wantSource: model.SourceRDAP, wantServer: "https://rdap.nic.li", wantDetail: "no response",
		},
		{
			name:       "fallback explains missing rdap",
			rdapErr:    srcerr.ErrNoSource,
			whoisErr:   &srcerr.SourceError{Source: model.SourceWhois, Err: fmt.Errorf("%w: i/o timeout", srcerr.ErrTimeout)},
			wantStatus: http.StatusGatewayTimeout, wantCode: codeUpstreamTimout,
			wantSource: model.SourceWhois, wantDetail: "no rdap server registered for .com",
		},
		{
			name: "internal errors keep quiet", rdapErr: errors.New("secret internals"),
			whoisErr: srcerr.ErrNoSource, wantStatus: http.StatusInternalServerError, wantCode: codeInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := lookup.NewService(&fakeSource{err: tt.rdapErr}, &fakeSource{err: tt.whoisErr}, nil)
			rec := httptest.NewRecorder()
			New(&config.Config{LookupTimeout: 2_000_000_000, MaxDomains: 5}, svc).
				ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/whois/example.com", nil))

			var body errorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode: %v", err)
			}
			e := body.Error
			if rec.Code != tt.wantStatus || e.Code != tt.wantCode {
				t.Fatalf("status/code = %d/%q, want %d/%q", rec.Code, e.Code, tt.wantStatus, tt.wantCode)
			}
			if e.Source != tt.wantSource || e.Server != tt.wantServer {
				t.Errorf("source/server = %q/%q, want %q/%q", e.Source, e.Server, tt.wantSource, tt.wantServer)
			}
			if !strings.Contains(e.Detail, tt.wantDetail) || (tt.wantDetail == "" && e.Detail != "") {
				t.Errorf("detail = %q, want %q", e.Detail, tt.wantDetail)
			}
		})
	}
}

func TestLookupSuccessShape(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestServer(registered(), nil).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/whois/example.com", nil))

	var res model.Result
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if res.Query != "example.com" || !res.IsRegistered || res.Meta.Source != model.SourceRDAP {
		t.Errorf("unexpected result: %+v", res)
	}
}

func TestRawPassthrough(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestServer(registered(), nil).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/whois/example.com?raw=true", nil))

	if got := rec.Header().Get("Content-Type"); got != "application/rdap+json" {
		t.Errorf("content type = %q", got)
	}
	if rec.Body.String() != "RAW-RDAP-BODY" {
		t.Errorf("raw body = %q", rec.Body.String())
	}
}

// rlServer builds a handler with rate limiting on and one valid API key.
func rlServer() http.Handler {
	cfg := &config.Config{
		LookupTimeout: 2_000_000_000,
		MaxDomains:    5,
		RatePerMinute: 60,
		RateBurst:     1,
		APIKeys:       []string{"devkey"},
	}
	svc := lookup.NewService(&fakeSource{result: registered()}, &fakeSource{err: srcerr.ErrNoSource}, nil)
	return New(cfg, svc)
}

func TestRateLimitThrottlesAnonymous(t *testing.T) {
	h := rlServer()
	get := func(key string) int {
		r := httptest.NewRequest(http.MethodGet, "/v1/whois/example.com", nil)
		r.RemoteAddr = "192.0.2.1:1111"
		if key != "" {
			r.Header.Set("X-API-Key", key)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec.Code
	}

	// burst=1: first anonymous request passes, the next is throttled.
	if got := get(""); got != http.StatusOK {
		t.Fatalf("first anon request = %d, want 200", got)
	}
	if got := get(""); got != http.StatusTooManyRequests {
		t.Fatalf("second anon request = %d, want 429", got)
	}
}

func TestValidKeyBypassesRateLimit(t *testing.T) {
	h := rlServer()
	for i := 0; i < 5; i++ {
		r := httptest.NewRequest(http.MethodGet, "/v1/whois/example.com", nil)
		r.RemoteAddr = "192.0.2.2:2222"
		r.Header.Set("X-API-Key", "devkey")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		if rec.Code != http.StatusOK {
			t.Fatalf("keyed request %d = %d, want 200 (should bypass limit)", i, rec.Code)
		}
	}
	// A wrong key gets no bypass and is throttled after the burst.
	throttled := false
	for i := 0; i < 3; i++ {
		r := httptest.NewRequest(http.MethodGet, "/v1/whois/example.com", nil)
		r.RemoteAddr = "192.0.2.3:3333"
		r.Header.Set("Authorization", "Bearer wrong")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		if rec.Code == http.StatusTooManyRequests {
			throttled = true
		}
	}
	if !throttled {
		t.Fatal("wrong key was never throttled; bypass is too permissive")
	}
}

// cacheServer builds a handler with a CDN cache header configured, whose RDAP source
// returns the given result/err.
func cacheServer(result *model.Result, err error) http.Handler {
	cfg := &config.Config{LookupTimeout: 2_000_000_000, MaxDomains: 5, CDNCacheControl: "public, s-maxage=3600"}
	svc := lookup.NewService(&fakeSource{result: result, err: err}, &fakeSource{err: srcerr.ErrNoSource}, nil)
	return New(cfg, svc)
}

func TestCacheHeaderOnlyOnSuccess(t *testing.T) {
	get := func(h http.Handler, path string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		return rec
	}

	if rec := get(cacheServer(registered(), nil), "/v1/whois/example.com"); rec.Header().Get("Cache-Control") == "" {
		t.Error("successful lookup should carry a Cache-Control header")
	}
	if rec := get(cacheServer(nil, srcerr.ErrUpstream), "/v1/whois/example.com"); rec.Header().Get("Cache-Control") != "" {
		t.Errorf("error response must not be cached, got %q", rec.Header().Get("Cache-Control"))
	}
}

func TestMultiCacheHeader(t *testing.T) {
	get := func(h http.Handler) string {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/multi?domains=a.com,b.com", nil))
		return rec.Header().Get("Cache-Control")
	}

	if get(cacheServer(registered(), nil)) == "" {
		t.Error("multi with all lookups succeeding should be cacheable")
	}
	if got := get(cacheServer(nil, srcerr.ErrUpstream)); got != "" {
		t.Errorf("multi with a failed lookup must not be cached, got %q", got)
	}
}

func TestBareLookupRejectsPost(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestServer(registered(), nil).ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/example.com", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestServer(nil, nil).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}
