package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

func handleGet(writer http.ResponseWriter, request *http.Request) {

	counter := request.Context().Value("counter").(string)
	per := request.Context().Value("per").(Duration)
	count := request.Context().Value("count").(int)
	sum := request.Context().Value("sum").(bool)

	result, err := dbExecuteOperation(
		counter,
		func(db *clover.DB) (any, error) {

			if per == Infinity || count <= 1 {
				if doc, err := db.FindFirst(
					query.NewQuery(counter).Sort(query.SortOption{Field: "time", Direction: -1}),
				); err != nil {
					return nil, err
				} else if doc == nil {
					return 0, nil
				} else {
					return doc.Get("value"), nil
				}
			}

			var thresholds []int64

			now := time.Now()
			thresholds = append(thresholds, now.UnixMilli())

			if count >= 2 {
				truncated := truncateToDuration(now, per)
				thresholds = append(thresholds, truncated.UnixMilli())

				for i := 1; i <= count-2; i++ {
					thresholds = append(thresholds, substractDuration(truncated, per, i).UnixMilli())
				}
			}

			var times []int64

			for i, threshold := range thresholds {

				criteria := query.Field("time").Lt(threshold)
				if i != len(thresholds)-1 {
					criteria = criteria.And(query.Field("time").GtEq(thresholds[i+1]))
				}

				if doc, err := db.FindFirst(
					query.NewQuery(counter).
						Sort(query.SortOption{Field: "time", Direction: -1}).
						Where(criteria),
				); err != nil {
					return nil, err
				} else if doc == nil {
					if sum {
						if doc, err := db.FindFirst(
							query.NewQuery(counter).Sort(query.SortOption{Field: "time", Direction: -1}).Where(query.Field("time").Lt(threshold)),
						); err != nil {
							return nil, err
						} else if doc == nil {
							times = append(times, 0)
						} else {
							times = append(times, doc.Get("value").(int64))
						}
					} else {
						times = append(times, 0)
					}
				} else {
					times = append(times, doc.Get("value").(int64))
				}
			}

			return times, nil
		},
	)
	if err != nil {
		throwHttpError(writer, 500, "Something went wrong while looking up this counter")
		log.Printf("[ERROR] %s\n", err.Error())
		return
	}

	json.NewEncoder(writer).Encode(newJsonApiResponse(
		200,
		"",
		result,
	))
}

func handleOperation(writer http.ResponseWriter, request *http.Request) {
	operation := request.Context().Value("operation").(Operation)
	counter := request.Context().Value("counter").(string)
	amount := int64(request.Context().Value("amount").(int))

	result, err := dbExecuteOperation(
		counter,
		func(db *clover.DB) (any, error) {

			var value int64

			if doc, err := db.FindFirst(
				query.NewQuery(counter).Sort(query.SortOption{Field: "time", Direction: -1}),
			); err != nil {
				return nil, err
			} else if doc == nil {
				value = 0
			} else {
				value = doc.Get("value").(int64)
			}

			switch operation {
			case Add:
				value += amount
			case Sub:
				value -= amount
			case Set:
				value = amount
			}

			doc := document.NewDocument()
			doc.Set("time", time.Now().UnixMilli())
			doc.Set("counter", counter)
			doc.Set("value", value)
			doc.Set("operation", operation.toString())
			doc.Set("amount", amount)

			_, err := db.InsertOne(counter, doc)
			return value, err
		},
	)
	if err != nil {
		throwHttpError(writer, 500, "Something went wrong while updating this counter")
		log.Printf("[ERROR] %s\n", err.Error())
		return
	}

	json.NewEncoder(writer).Encode(newJsonApiResponse(
		200,
		"",
		result,
	))
}
