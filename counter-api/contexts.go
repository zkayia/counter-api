package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func counterContext(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if counter := chi.URLParam(request, "counter"); counter != "" {
			handler.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), "counter", counter)))
		} else {
			throwHttpError(writer, 400, "Invalid counter name, should match `[a-zA-Z0-9-_]+`")
		}

	})
}

func getQueryContext(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var per Duration
		var count int
		var sum bool
		var err error

		if valueStr := request.URL.Query().Get("per"); valueStr != "" {
			per, err = durationFromString(valueStr)
		} else {
			per = Infinity
		}
		if err != nil {
			throwHttpError(writer, 422, "Invalid parameter 'per', should be `infinity,second,minute,hour,day,month,year,decade,century`")
			return
		}

		err = nil

		if valueStr := request.URL.Query().Get("count"); valueStr != "" {
			count, err = strconv.Atoi(valueStr)
		} else {
			count = 1
		}
		if err != nil {
			throwHttpError(writer, 422, "Invalid parameter 'count', should match `[0-9]+`")
			return
		}

		if valueStr := request.URL.Query().Get("sum"); valueStr == "" || valueStr == "false" {
			sum = false
		} else if valueStr == "true" {
			sum = true
		} else {
			throwHttpError(writer, 422, "Invalid parameter 'sum', should match `true,false`")
			return
		}

		getQueryContext := context.WithValue(request.Context(), "per", per)
		getQueryContext = context.WithValue(getQueryContext, "count", intAbs(count))
		getQueryContext = context.WithValue(getQueryContext, "sum", sum)
		handler.ServeHTTP(writer, request.WithContext(getQueryContext))
	})
}

func valueQueryContext(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var amount int
		var err error

		if amountStr := request.URL.Query().Get("amount"); amountStr != "" {
			amount, err = strconv.Atoi(amountStr)
		} else {
			amount = 1
		}
		if err != nil {
			throwHttpError(writer, 422, "Invalid amount, should be `[0-9]+`")
			return
		}

		handler.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), "amount", amount)))
	})
}

func operationContext(operation Operation) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			handler.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), "operation", operation)))
		})
	}
}
