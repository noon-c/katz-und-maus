package main

import (
	"bytes"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func runWrite() {
	ctx, cancel := Ctx30s()
	defer cancel()

	master := os.Getenv("MASTER_PASSWORD")
	if master == "" {
		log.Fatal("set MASTER_PASSWORD")
	}

	text := os.Getenv("VAULT_TEXT")
	if text == "" {
		text = "hello from encrypted vault"
	}

	r2, err := NewR2(ctx)
	if err != nil {
		log.Fatal(err)
	}

	blob, err := Encrypt(master, []byte(text), DefaultParams())
	if err != nil {
		log.Fatal(err)
	}

	_, err = r2.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r2.Bucket),
		Key:         aws.String(r2.Key),
		Body:        bytes.NewReader(blob),
		ContentType: aws.String("application/octet-stream"),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("OK: uploaded r2://%s/%s (%d bytes)", r2.Bucket, r2.Key, len(blob))
}
