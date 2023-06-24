package fst

import (
	"crypto/sha256"
	"strings"
	"testing"
	"time"
)

func TestNewConverter(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewConverter(&ConverterConfig{
		SecretKey: []byte(`secret`),
	})
	if converter == nil {
		t.Error("Converter is nil")
	}

	converter = NewConverter(&ConverterConfig{
		SecretKey: []byte(`secret`),
		Postfix:   []byte(`postfix`),
		HashType:  sha256.New,
	})
	if converter == nil {
		t.Error("Converter with postfix is nil")
	}

	converter = NewConverter(&ConverterConfig{
		SecretKey:      []byte(`secret`),
		Postfix:        nil,
		ExpirationTime: time.Minute * 5,
		HashType:       sha256.New,
	})
	if converter == nil {
		t.Error("Converter with expiration time is nil")
	}

	converter = NewConverter(&ConverterConfig{
		SecretKey:      []byte(`secret`),
		Postfix:        []byte(`postfix`),
		ExpirationTime: time.Minute * 5,
		HashType:       sha256.New,
	})
	if converter == nil {
		t.Error("Converter with postfix and expiration time is nil")
	}
}

func TestConverter_NewToken(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewConverter(&ConverterConfig{
		SecretKey: []byte(`secret`),
	})

	token := converter.NewToken([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token is nil")
	}
}

func TestConverter_ParseToken(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewConverter(&ConverterConfig{
		SecretKey: []byte(`secret`),
	})

	token := converter.NewToken([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token is nil")
	}

	value, err := converter.ParseToken(token)
	if err != nil {
		t.Error("Token parse err: ", err)
	}
	if string(value) != `token` {
		t.Error("Token parse err: ", string(value), " != ", `token`)
	}
}

func TestConverter_NewTokenWithExpire(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewConverter(&ConverterConfig{
		SecretKey:      []byte(`secret`),
		Postfix:        nil,
		ExpirationTime: time.Minute * 5,
		HashType:       sha256.New,
	})

	token := converter.NewToken([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token with expire time is nil")
	}
}

func TestConverter_ParseTokenWithExpire(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewConverter(&ConverterConfig{
		SecretKey:      []byte(`secret`),
		Postfix:        nil,
		ExpirationTime: time.Minute * 5,
		HashType:       sha256.New,
	})

	token := converter.NewToken([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token with expire time is nil")
	}

	value, err := converter.ParseToken(token)
	if err != nil {
		t.Error("Token with expire time parse err: ", err)
	}

	if string(value) != `token` {
		t.Error("Token with expire time parse err: ", string(value), " != ", `token`)
	}
}

func TestConverter_ExpiredToken(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewConverter(&ConverterConfig{
		SecretKey:          []byte(`secret`),
		Postfix:            nil,
		ExpirationTime:     time.Second * 1,
		HashType:           sha256.New,
		WithExpirationTime: true,
	})

	token := converter.NewToken([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token with expire time is nil")
	}

	time.Sleep(time.Second * 3)

	_, err := converter.ParseToken(token)
	if err == nil {
		t.Error("Token with expire time parse err: ", "token is not expired!")
	} else {
		if err != TokenExpired {
			t.Error("Token with expire time parse err: ", err)
		}
	}
}

func TestConverter_Postfix(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewConverter(&ConverterConfig{
		SecretKey:      []byte(`secret`),
		Postfix:        []byte(`postfix`),
		ExpirationTime: time.Minute * 5,
		HashType:       sha256.New,
	})

	token := converter.NewToken([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token is nil")
	}

	value, err := converter.ParseToken(token)
	if err != nil {
		t.Error("Token with postfix parse err: ", err)
	}

	if string(value) != `token` {
		t.Error("Token with postfix parse err: ", string(value), " != ", `token`)
	}
}

func TestConverter_ExpiredTokenWithPostfix(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewConverter(&ConverterConfig{
		SecretKey:          []byte(`secret`),
		Postfix:            []byte(`postfix`),
		ExpirationTime:     time.Second * 1,
		HashType:           sha256.New,
		WithExpirationTime: true,
	})

	token := converter.NewToken([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token with expire time and postfix is nil")
	}

	time.Sleep(time.Second * 3)

	_, err := converter.ParseToken(token)
	if err == nil {
		t.Error("Token with expire time and postfix parse err: ", "token is not expired!")
	} else {
		if err != TokenExpired {
			t.Error("Token with expire time and postfix parse err: ", err)
		}
	}
}

func TestConverter_SecuredTest(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewConverter(&ConverterConfig{
		SecretKey: []byte(`secret`),
		Postfix:   []byte(`postfix`),
		HashType:  sha256.New,
	})

	token := converter.NewToken([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token is nil")
	}

	components := strings.Split(token, ".")
	components[0] = components[0] + "mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm"

	token = strings.Join(components, ".")

	_, err := converter.ParseToken(token)
	if err == nil {
		t.Error("Token parse err: ", "token is not secured!")
	} else {
		if err != InvalidSignature {
			t.Error("Token parse err: ", err)
		}
	}
}

func TestConverter_SecuredTestWithExpire(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewConverter(&ConverterConfig{
		SecretKey:          []byte(`secret`),
		Postfix:            []byte(`postfix`),
		HashType:           sha256.New,
		ExpirationTime:     time.Minute * 5,
		WithExpirationTime: true,
	})

	token := converter.NewToken([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token is nil")
	}

	components := strings.Split(token, ".")
	components[0] = components[0] + "mmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm"

	token = strings.Join(components, ".")

	_, err := converter.ParseToken(token)
	if err == nil {
		t.Error("Token parse err: ", "token is not secured!")
	} else {
		if err != InvalidSignature {
			t.Error("Token parse err: ", err)
		}
	}

	token2 := converter.NewToken([]byte(`token`))
	if token2 == "" || len(token2) == 0 {
		t.Error("Token2 is nil")
	}

	components = strings.Split(token2, ".")
	components[2] = components[2] + "1"

	token2 = strings.Join(components, ".")

	_, err = converter.ParseToken(token2)
	if err == nil {
		t.Error("Token parse err: ", "token is not secured!")
	} else {
		if err != InvalidSignature {
			t.Error("Token parse err: ", err)
		}
	}
}
