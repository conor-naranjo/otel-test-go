package main

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
)

func reverseStrSvc(ctx context.Context, str string) (string, error) {
	ctx, span := tracer.Start(ctx, "reverse_string_service")
	defer span.End()
	span.SetAttributes(attribute.String("str", str))

	// Apparently reversing a string is a non-trivial problem
	// https://groups.google.com/g/golang-nuts/c/oPuBaYJ17t4

	// get unicode code points
	n := 0
	rune := make([]rune, len(str))
	for _, r := range str {
		rune[n] = r
		n++
	}
	rune = rune[0:n]

	// Reverse
	for i := 0; i < n/2; i++ {
		rune[i], rune[n-1-i] = rune[n-1-i], rune[i]
	}

	// convert back to utf-8
	return string(rune), nil
}
