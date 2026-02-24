package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/crypto/argon2"
)

var (
	magic   = []byte("KUMVAULT")
	version = byte(1)

	ErrBadPasswordOrCorrupt = errors.New("bad password or corrupt vault")
	ErrBadFormat            = errors.New("bad vault format")
)

type ArgonParams struct {
	Time    uint32
	Memory  uint32 // KiB
	Threads uint8
	SaltLen uint8
}

func DefaultParams() ArgonParams {
	return ArgonParams{
		Time:    3,
		Memory:  64 * 1024, // 64 MiB
		Threads: 2,
		SaltLen: 16,
	}
}

func deriveKey(password string, salt []byte, p ArgonParams) []byte {
	return argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Threads, 32)
}

func buildHeader(p ArgonParams, salt, nonce []byte) []byte {
	var b bytes.Buffer
	b.Write(magic)
	b.WriteByte(version)

	_ = binary.Write(&b, binary.BigEndian, p.Time)
	_ = binary.Write(&b, binary.BigEndian, p.Memory)
	b.WriteByte(p.Threads)

	b.WriteByte(byte(len(salt)))
	b.Write(salt)

	b.WriteByte(byte(len(nonce)))
	b.Write(nonce)

	return b.Bytes()
}

func parseHeader(blob []byte) (ArgonParams, []byte, []byte, int, error) {
	if len(blob) < 8+1+4+4+1+1+1 {
		return ArgonParams{}, nil, nil, 0, ErrBadFormat
	}
	if !bytes.Equal(blob[:8], magic) {
		return ArgonParams{}, nil, nil, 0, ErrBadFormat
	}
	if blob[8] != version {
		return ArgonParams{}, nil, nil, 0, ErrBadFormat
	}

	i := 9
	readU32 := func() (uint32, bool) {
		if i+4 > len(blob) {
			return 0, false
		}
		v := binary.BigEndian.Uint32(blob[i : i+4])
		i += 4
		return v, true
	}

	timeCost, ok := readU32()
	if !ok {
		return ArgonParams{}, nil, nil, 0, ErrBadFormat
	}
	memKiB, ok := readU32()
	if !ok {
		return ArgonParams{}, nil, nil, 0, ErrBadFormat
	}
	if i >= len(blob) {
		return ArgonParams{}, nil, nil, 0, ErrBadFormat
	}
	threads := blob[i]
	i++

	if i >= len(blob) {
		return ArgonParams{}, nil, nil, 0, ErrBadFormat
	}
	saltLen := int(blob[i])
	i++
	if i+saltLen > len(blob) {
		return ArgonParams{}, nil, nil, 0, ErrBadFormat
	}
	salt := blob[i : i+saltLen]
	i += saltLen

	if i >= len(blob) {
		return ArgonParams{}, nil, nil, 0, ErrBadFormat
	}
	nonceLen := int(blob[i])
	i++
	if i+nonceLen > len(blob) {
		return ArgonParams{}, nil, nil, 0, ErrBadFormat
	}
	nonce := blob[i : i+nonceLen]
	i += nonceLen

	p := ArgonParams{
		Time:    timeCost,
		Memory:  memKiB,
		Threads: threads,
		SaltLen: uint8(saltLen),
	}
	return p, append([]byte(nil), salt...), append([]byte(nil), nonce...), i, nil
}

func Encrypt(password string, plaintext []byte, p ArgonParams) ([]byte, error) {
	salt := make([]byte, p.SaltLen)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	key := deriveKey(password, salt, p)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	header := buildHeader(p, salt, nonce)
	ciphertext := gcm.Seal(nil, nonce, plaintext, header)
	return append(header, ciphertext...), nil
}

func Decrypt(password string, blob []byte) ([]byte, error) {
	p, salt, nonce, headerLen, err := parseHeader(blob)
	if err != nil {
		return nil, err
	}
	header := blob[:headerLen]
	ciphertext := blob[headerLen:]

	key := deriveKey(password, salt, p)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, header)
	if err != nil {
		return nil, ErrBadPasswordOrCorrupt
	}
	return plaintext, nil
}

type R2 struct {
	Client *s3.Client
	Bucket string
	Key    string
}

func NewR2(ctx context.Context) (*R2, error) {
	accountID := os.Getenv("R2_ACCOUNT_ID")
	ak := os.Getenv("R2_ACCESS_KEY_ID")
	sk := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucket := os.Getenv("R2_BUCKET")
	key := os.Getenv("VAULT_KEY")
	if key == "" {
		key = "vault.bin"
	}

	if accountID == "" || ak == "" || sk == "" || bucket == "" {
		return nil, fmt.Errorf("missing env vars: R2_ACCOUNT_ID, R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_BUCKET")
	}

	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(ak, sk, "")),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &R2{Client: client, Bucket: bucket, Key: key}, nil
}

func Ctx30s() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}
