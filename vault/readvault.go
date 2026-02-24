package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

func runRead() {
	ctx, cancel := Ctx30s()
	defer cancel()

	master := os.Getenv("MASTER_PASSWORD")
	if master == "" {
		log.Fatal("set MASTER_PASSWORD")
	}

	r2, err := NewR2(ctx)
	if err != nil {
		log.Fatal(err)
	}

	out, err := r2.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r2.Bucket),
		Key:    aws.String(r2.Key),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			log.Fatalf("GetObject error: %s (%s)", apiErr.ErrorMessage(), apiErr.ErrorCode())
		}
		log.Fatal(err)
	}
	defer out.Body.Close()

	blob, err := io.ReadAll(out.Body)
	if err != nil {
		log.Fatal(err)
	}

	plain, err := Decrypt(master, blob)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.TrimRight(string(plain), "\n"))
}
