package otelemetry

import "fmt"

func handleErr(err error, s string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", s, err))
	}
}
