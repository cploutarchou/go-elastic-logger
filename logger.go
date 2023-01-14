package go_elastic_logger

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var index = "logs"
var logLevel LogLevel = ERROR
var esClient *elasticsearch.Client
var mu sync.Mutex
var once sync.Once
var layout = "2023-01-14 13:04:05"

func SetIndex(newIndex string) {
	index = newIndex
}
func SetLogLevel(level LogLevel) {
	logLevel = level
}

func SetElasticsearchClient(client *elasticsearch.Client) {
	once.Do(func() {
		esClient = client
	})
}
func SetTimeFormat(customLayout string) {
	layout = customLayout
}
func log(level LogLevel, message string) {
	timestamp := time.Now().Format(layout)
	levelStr := level.String()
	switch level {
	case DEBUG:
		levelStr = "DEBUG"
	case INFO:
		levelStr = "INFO"
	case WARNING:
		levelStr = "WARNING"
	case ERROR:
		levelStr = "ERROR"
	case FATAL:
		levelStr = "FATAL"
	}
	
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	function := runtime.FuncForPC(pc).Name()
	parts := strings.Split(function, ".")
	function = parts[1]
	logLine := map[string]interface{}{
		"timestamp": timestamp,
		"level":     levelStr,
		"message":   message,
		"file":      fmt.Sprintf("%s:%d", file, line),
		"function":  function + "()",
	}
	jsonLogLine, err := json.Marshal(logLine)
	if err != nil {
		_, _ = io.WriteString(os.Stderr, fmt.Sprintf("Error marshalling logLine to json: %s", err))
	}
	_, _ = io.WriteString(os.Stdout, string(jsonLogLine))
	mu.Lock()
	defer mu.Unlock()
	if level < logLevel {
		return
	}
	if esClient != nil {
		_, err := esClient.Index(
			index,
			strings.NewReader(string(jsonLogLine)),
			esClient.Index.WithContext(context.Background()),
		)
		if err != nil {
			_, _ = io.WriteString(os.Stderr, fmt.Sprintf("Error indexing log to Elasticsearch: %s", err))
		}
	}
}
func (level LogLevel) String() string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "ERROR"
	}
}

func Debug(message string) {
	log(DEBUG, message)
}

func Debugf(format string, v ...interface{}) {
	Debug(fmt.Sprintf(format, v...))
}
func Info(message string) {
	log(INFO, message)
}
func Infof(format string, v ...interface{}) {
	Info(fmt.Sprintf(format, v...))
}
func Warning(message string) {
	log(WARNING, message)
}

func Warningf(format string, v ...interface{}) {
	Warning(fmt.Sprintf(format, v...))
}
func Error(message string) {
	log(ERROR, message)
}

func Errorf(format string, v ...interface{}) {
	Error(fmt.Sprintf(format, v...))
}
func Fatal(message string) {
	log(FATAL, message)
	os.Exit(1)
}
func Fatalf(format string, v ...interface{}) {
	Fatal(fmt.Sprintf(format, v...))
}
