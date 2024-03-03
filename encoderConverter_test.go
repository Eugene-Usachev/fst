package fst

import (
	"crypto/sha256"
	"errors"
	"testing"
	"time"
)

// We can not to test all functions, because EncoderConverter is built on top of Converter, that is full covered by tests.
func TestEncoderConverter(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewEncodedConverter(&ConverterConfig{
		SecretKey:      []byte(`secret`),
		Postfix:        []byte(`postfix`),
		HashType:       sha256.New,
		ExpirationTime: time.Minute * 5,
	})

	if converter == nil {
		t.Error("Converter is nil")
		return
	}

	token := converter.NewToken([]byte(`token`))
	if len(token) < 1 {
		t.Error("token is empty")
		return
	}

	value, err := converter.ParseToken(token)
	if err != nil {
		t.Error("error: ", err)
		return
	}

	if string(value) != `token` {
		t.Error("value is not token, but ", string(value))
	}
}

func TestEncoderConverterExpired(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewEncodedConverter(&ConverterConfig{
		SecretKey:      []byte(`secret`),
		Postfix:        []byte(`postfix`),
		HashType:       sha256.New,
		ExpirationTime: time.Second * 5,
	})

	if converter == nil {
		t.Error("Converter is nil")
		return
	}

	token := converter.NewToken([]byte(`token`))
	if len(token) < 1 {
		t.Error("token is empty")
		return
	}

	time.Sleep(time.Second * 4)

	value, err := converter.ParseToken(token)
	if err != nil {
		t.Error("error: ", err)
		return
	}

	if string(value) != `token` {
		t.Error("value is not token, but ", string(value))
	}

	time.Sleep(time.Second * 3)

	_, err = converter.ParseToken(token)
	if err == nil {
		t.Error("token is not expired!")
		return
	}

	if !errors.Is(err, TokenExpired) {
		t.Error("unexpected error: ", err)
		return
	}
}
