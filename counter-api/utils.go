package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func intAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func newJsonApiResponse(code int, message string, data any) map[string]interface{} {
	var responseMessage any = nil
	responseData := data
	success := 200 <= code && code < 300
	if !success {
		responseMessage = message
		responseData = nil
	}
	return map[string]interface{}{
		"success":    success,
		"code":       code,
		"statusText": http.StatusText(code),
		"time":       time.Now().Format("2006/01/02 15:04:05"),
		"message":    responseMessage,
		"data":       responseData,
	}
}

func throwHttpError(writer http.ResponseWriter, code int, message string) {
	writer.WriteHeader(code)
	json.NewEncoder(writer).Encode(newJsonApiResponse(
		code,
		message,
		nil,
	))
}

func truncateToDuration(t time.Time, duration Duration) time.Time {
	switch duration {
	case Second:
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	case Minute:
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
	case Hour:
		return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
	case Day:
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	case Month:
		return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	default:
		return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	}
}

func substractDuration(t time.Time, duration Duration, count int) time.Time {
	countDuration := time.Duration(-count)
	switch duration {
	case Second:
		return t.Add(countDuration * time.Second)
	case Minute:
		return t.Add(countDuration * time.Minute)
	case Hour:
		return t.Add(countDuration * time.Hour)
	case Day:
		return t.AddDate(0, 0, -count)
	case Month:
		return t.AddDate(0, -count, 0)
	case Year:
		return t.AddDate(-count, 0, 0)
	case Decade:
		return t.AddDate(-count*10, 0, 0)
	case Century:
		return t.AddDate(-count*100, 0, 0)
	default:
		return t
	}
}

func setupLogFile() error {

	logfile, err := os.OpenFile(
		filepath.Join(
			*logsPath,
			time.Now().Format("2006-01-02_15-04-05")+".log",
		),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		return err
	}
	defer logfile.Close()

	log.SetOutput(logfile)

	return nil
}

func envWithDefaultString(key string, defaultValue string) string {
	var value, isPresent = os.LookupEnv(key)
	if isPresent {
		return value
	}
	return defaultValue
}

func envWithDefaultDuration(key string, defaultValue time.Duration) time.Duration {
	var value, isPresent = os.LookupEnv(key)
	if isPresent {
		var parsed, err = time.ParseDuration(value)
		if err == nil {
			return parsed
		}
	}
	return defaultValue
}
