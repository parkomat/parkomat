package web

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ResponseWriter with added metrics
type HttpLog struct {
	http.ResponseWriter

	IP        string
	Datetime  time.Time
	Method    string
	Host      string
	URI       string
	Status    int
	BytesSent uint64
	Referer   string
	Duration  time.Duration
}

// Implement http.ResponseWriter method and count bytes written
func (httpLog *HttpLog) Write(data []byte) (int, error) {
	n, err := httpLog.ResponseWriter.Write(data)

	httpLog.BytesSent += uint64(n)

	return n, err
}

// Let's hijack response's status code
func (httpLog *HttpLog) WriteHeader(status int) {
	httpLog.Status = status
	httpLog.ResponseWriter.WriteHeader(status)
}

func (httpLog *HttpLog) WriteLog(output io.Writer) {
	datetime := httpLog.Datetime.Format("02/Jan/2006 15:04:05")

	// TODO: this should be configurable
	fmt.Fprintf(output, "[%s] %d %s %s %s %s %s %d %.4f\r\n",
		datetime,
		httpLog.Status,
		httpLog.IP,
		httpLog.Method,
		httpLog.Host,
		httpLog.URI,
		httpLog.Referer,
		httpLog.BytesSent,
		httpLog.Duration.Seconds())
}

type HttpLogHandler struct {
	Handler http.Handler
	Output  io.Writer
}

func NewHttpLogHandler(handler http.Handler, output io.Writer) http.Handler {
	return &HttpLogHandler{
		Handler: handler,
		Output:  output,
	}
}

// Logging request parameters and measuring time spent on dealing with the request
func (httpLogHandler *HttpLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	httpLog := &HttpLog{
		ResponseWriter: w,
		IP:             strings.Split(r.RemoteAddr, ":")[0],
		Datetime:       time.Now().UTC(),
		Method:         r.Method,
		URI:            r.RequestURI,
		Host:           r.Host,
		Referer:        r.Referer(),
		BytesSent:      0,
	}

	httpLogHandler.Handler.ServeHTTP(httpLog, r)
	timeStop := time.Now().UTC()

	httpLog.Duration = timeStop.Sub(httpLog.Datetime)

	httpLog.WriteLog(httpLogHandler.Output)
}
