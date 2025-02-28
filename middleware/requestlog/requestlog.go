//go:build !solution

package requestlog

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"runtime"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

type RequestInfo struct {
	RequestID string
	Request   *http.Request
	StartTime time.Time
	Wrapper   *ResponseWriter
}

func GenerateRequestID() string {
	return createRandomID()
}

func createRandomID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return getCurrentTimeAsID()
	}
	return hex.EncodeToString(bytes)
}

func getCurrentTimeAsID() string {
	return time.Now().Format(time.RFC3339Nano)
}

func Log(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestInfo := newRequestInfo(r)

			logRequestStart(logger, requestInfo)

			startTime := getCurrentTime()

			wrapper := &ResponseWriter{ResponseWriter: w, status: http.StatusOK}
			requestInfo.Wrapper = wrapper
			requestInfo.StartTime = startTime

			defer logRequestCompletion(logger, requestInfo)

			serveRequest(next, wrapper, r)
		})
	}
}

func newRequestInfo(request *http.Request) *RequestInfo {
	requestID := GenerateRequestID()
	return &RequestInfo{
		RequestID: requestID,
		Request:   request,
	}
}

func logRequestStart(logger *zap.Logger, requestInfo *RequestInfo) {
	message := getRequestStartMessage()
	fields := getRequestLogFields(requestInfo)
	logInfo(logger, message, fields)
}

func getRequestStartMessage() string {
	return "request started"
}

func logRequestCompletion(logger *zap.Logger, requestInfo *RequestInfo) {
	duration := calculateRequestDuration(requestInfo)

	if err := recover(); err != nil {
		logRequestError(logger, requestInfo, duration, err)
		panic(err)
	}

	logRequestSuccess(logger, requestInfo, duration)
}

func calculateRequestDuration(requestInfo *RequestInfo) time.Duration {
	return time.Since(requestInfo.StartTime)
}

func logRequestSuccess(logger *zap.Logger, requestInfo *RequestInfo, duration time.Duration) {
	message := getRequestSuccessMessage()
	fields := append(getRequestLogFields(requestInfo),
		zap.Duration("duration", duration),
		zap.Int("status_code", requestInfo.Wrapper.status),
	)
	logInfo(logger, message, fields)
}

func getRequestSuccessMessage() string {
	return "request finished"
}

func logRequestError(logger *zap.Logger, requestInfo *RequestInfo, duration time.Duration, err interface{}) {
	stackTrace := captureStackTrace()

	message := getRequestErrorMessage()
	fields := append(getRequestLogFields(requestInfo),
		zap.Duration("duration", duration),
		zap.Int("status_code", requestInfo.Wrapper.status),
		zap.String("error", fmt.Sprintf("%v", err)),
		zap.String("stack_trace", stackTrace),
	)
	logError(logger, message, fields)
}

func getRequestErrorMessage() string {
	return "request panicked"
}

func captureStackTrace() string {
	stack := make([]byte, 4096)
	stack = stack[:runtime.Stack(stack, false)]
	return string(stack)
}

func getRequestLogFields(requestInfo *RequestInfo) []zap.Field {
	return []zap.Field{
		zap.String("request_id", requestInfo.RequestID),
		zap.String("path", requestInfo.Request.URL.Path),
		zap.String("method", requestInfo.Request.Method),
	}
}

func logInfo(logger *zap.Logger, message string, fields []zap.Field) {
	logger.Info(message, fields...)
}

func logError(logger *zap.Logger, message string, fields []zap.Field) {
	logger.Error(message, fields...)
}

func getCurrentTime() time.Time {
	return time.Now()
}

func serveRequest(next http.Handler, w http.ResponseWriter, request *http.Request) {
	next.ServeHTTP(w, request)
}
