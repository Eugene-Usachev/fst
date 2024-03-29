package fst

import (
	"crypto/sha256"
	"errors"
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

	converter := NewConverter(&ConverterConfig{
		SecretKey: []byte(`secret`),
	})

	token := converter.NewToken([]byte(`token`))
	if len(token) == 0 {
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
	if len(token) == 0 {
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
	if len(token) == 0 {
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
	if len(token) == 0 {
		t.Error("Token with expire time is nil")
	}

	value, err := converter.ParseToken(token)
	if err != nil {
		t.Error("Token with expire time parse err: ", err)
		return
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
		SecretKey:      []byte(`secret`),
		Postfix:        nil,
		ExpirationTime: time.Second * 1,
		HashType:       sha256.New,
	})

	token := converter.NewToken([]byte(`token`))
	if len(token) == 0 {
		t.Error("Token with expire time is nil")
	}

	time.Sleep(time.Second * 3)

	_, err := converter.ParseToken(token)
	if err == nil {
		t.Error("Token with expire time parse err: ", "token is not expired!")
	} else {
		if !errors.Is(err, TokenExpired) {
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
		SecretKey: []byte(`secret`),
		Postfix:   []byte(`postfix`),
		HashType:  sha256.New,
	})

	token := converter.NewToken([]byte(`token`))
	if len(token) == 0 {
		t.Error("Token is nil")
	}

	value, err := converter.ParseToken(token)
	if err != nil {
		t.Error("Token with postfix parse err: ", err)
		return
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
		SecretKey:      []byte(`secret`),
		Postfix:        []byte(`postfix`),
		ExpirationTime: time.Second * 1,
		HashType:       sha256.New,
	})

	token := converter.NewToken([]byte(`token`))
	if len(token) == 0 {
		t.Error("Token with expire time and postfix is nil")
	}

	time.Sleep(time.Second * 3)

	_, err := converter.ParseToken(token)
	if err == nil {
		t.Error("Token with expire time and postfix parse err: ", "token is not expired!")
	} else {
		if !errors.Is(err, TokenExpired) {
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
	if len(token) == 0 {
		t.Error("Token is nil")
	}

	token = append(token, []byte("admin")...)

	_, err := converter.ParseToken(token)
	if err == nil {
		t.Error("Token parse err: ", "token is not secured!")
	} else {
		if !errors.Is(err, InvalidSignature) {
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
		SecretKey:      []byte(`secret`),
		Postfix:        []byte(`postfix`),
		HashType:       sha256.New,
		ExpirationTime: time.Minute * 5,
	})

	token := converter.NewToken([]byte(`token`))
	if len(token) == 0 {
		t.Error("Token is nil")
	}

	_, err := converter.ParseToken(token)
	if err != nil {
		t.Error("Token parse err: ", err)
	}

	token2 := converter.NewToken([]byte(`token`))
	if len(token2) == 0 {
		t.Error("Token2 is nil")
	}

	token2[0] = 255

	_, err = converter.ParseToken(token2)
	if err == nil {
		t.Error("Token parse err: ", "token is not secured!")
	} else {
		if !errors.Is(err, InvalidSignature) {
			t.Error("Token parse err: ", err)
		}
	}
}
