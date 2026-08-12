package server

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/lissy93/who-dat/internal/domain"
	"github.com/lissy93/who-dat/internal/srcerr"
)

// Error codes returned in the error envelope.
const (
	codeInvalidDomain    = "INVALID_DOMAIN"
	codeUnsupportedTLD   = "UNSUPPORTED_TLD"
	codeRateLimited      = "RATE_LIMITED"
	codeMethodNotAllowed = "METHOD_NOT_ALLOWED"
	codeUpstreamError    = "UPSTREAM_ERROR"
	codeUpstreamTimout   = "UPSTREAM_TIMEOUT"
	codeInternal         = "INTERNAL_ERROR"
)

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Query   string `json:"query"`
	Source  string `json:"source,omitempty"` // lookup backend that failed: "rdap" or "whois"
	Server  string `json:"server,omitempty"` // upstream server queried, when known
	Detail  string `json:"detail,omitempty"` // underlying error chain, for diagnostics
}

// writeError sends the standard error envelope with the given status.
func writeError(w http.ResponseWriter, status int, code, message, query string) {
	writeJSON(w, status, errorBody{Error: errorDetail{Code: code, Message: message, Query: query}})
}

// writeDomainError maps a domain.Parse error to a 400 envelope.
func writeDomainError(w http.ResponseWriter, err error, query string) {
	switch {
	case errors.Is(err, domain.ErrNoPublicSuffix):
		writeError(w, http.StatusBadRequest, codeInvalidDomain, "no valid public suffix", query)
	default:
		writeError(w, http.StatusBadRequest, codeInvalidDomain, "could not parse a registrable domain", query)
	}
}

// lookupError sorts a lookup failure into status + envelope, naming the failing
// source, server, and underlying detail when known.
func lookupError(err error, query string) (status int, retryAfter time.Duration, d errorDetail) {
	d = errorDetail{Query: query, Detail: err.Error()}
	var se *srcerr.SourceError
	if errors.As(err, &se) {
		d.Source, d.Server = se.Source, se.Server
	}

	var rl *srcerr.RateLimited
	switch {
	case errors.As(err, &rl):
		d.Code, d.Message = codeRateLimited, "rate limited; retry later"
		return http.StatusTooManyRequests, rl.RetryAfter, d
	case errors.Is(err, srcerr.ErrTimeout), errors.Is(err, context.DeadlineExceeded):
		d.Code, d.Message = codeUpstreamTimout, "registry timed out"
		return http.StatusGatewayTimeout, 0, d
	case errors.Is(err, srcerr.ErrNoSource):
		d.Code, d.Message = codeUnsupportedTLD, "no rdap or whois source for this tld"
		return http.StatusNotImplemented, 0, d
	case errors.Is(err, srcerr.ErrUpstream):
		d.Code, d.Message = codeUpstreamError, "registry unreachable or returned an invalid response"
		return http.StatusBadGateway, 0, d
	default:
		d.Code, d.Message, d.Detail = codeInternal, "internal error", "" // internals stay internal
		return http.StatusInternalServerError, 0, d
	}
}

// writeLookupError maps an orchestrator/source error to the right status and envelope.
func writeLookupError(w http.ResponseWriter, err error, query string) {
	status, retryAfter, d := lookupError(err, query)
	if retryAfter > 0 {
		w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
	}
	writeJSON(w, status, errorBody{Error: d})
}
