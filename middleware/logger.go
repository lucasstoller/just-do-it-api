package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Log request
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create custom response writer to capture response
		buf := &bytes.Buffer{}
		rw := &responseWriter{
			ResponseWriter: w,
			body:           buf,
			statusCode:     http.StatusOK, // Default status code
		}

		// Call the next handler
		next.ServeHTTP(rw, r)

		// Format request body for logging
		var prettyRequest string
		if len(requestBody) > 0 {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, requestBody, "", "  "); err == nil {
				prettyRequest = prettyJSON.String()
			} else {
				prettyRequest = string(requestBody)
			}
		}

		// Format response body for logging
		var prettyResponse string
		if rw.body.Len() > 0 {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, rw.body.Bytes(), "", "  "); err == nil {
				prettyResponse = prettyJSON.String()
			} else {
				prettyResponse = rw.body.String()
			}
		}

		// Log the request and response details
		log.Printf(`
Request:
  Method: %s
  Path: %s
  Headers: %v
  Body: %s

Response:
  Status: %d
  Duration: %v
  Body: %s
`,
			r.Method,
			r.URL.Path,
			r.Header,
			prettyRequest,
			rw.statusCode,
			time.Since(startTime),
			prettyResponse,
		)
	}
}
