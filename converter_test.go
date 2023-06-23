package fst

import (
	"crypto/sha256"
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
		SecretKey:   []byte(`secret`),
		Postfix:     nil,
		ExpireTime:  time.Minute * 5,
		HashType:    sha256.New,
		DisableLogs: true,
	})
	if converter == nil {
		t.Error("Converter is nil")
	}
}

func TestConverter_NewToken(t *testing.T) {
	defer func() {
		if pn := recover(); pn != nil {
			t.Error("panic handled: ", pn)
		}
	}()

	converter := NewConverter(&ConverterConfig{
		SecretKey:   []byte(`secret`),
		Postfix:     nil,
		ExpireTime:  time.Minute * 5,
		HashType:    sha256.New,
		DisableLogs: true,
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
		SecretKey:   []byte(`secret`),
		Postfix:     nil,
		ExpireTime:  time.Minute * 5,
		HashType:    sha256.New,
		DisableLogs: true,
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
		SecretKey:   []byte(`secret`),
		Postfix:     nil,
		ExpireTime:  time.Minute * 5,
		HashType:    sha256.New,
		DisableLogs: true,
	})

	token := converter.NewTokenWithExpire([]byte(`token`))
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
		SecretKey:   []byte(`secret`),
		Postfix:     nil,
		ExpireTime:  time.Minute * 5,
		HashType:    sha256.New,
		DisableLogs: true,
	})

	token := converter.NewTokenWithExpire([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token with expire time is nil")
	}

	value, err := converter.ParseTokenWithExpire(token)
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
		SecretKey:   []byte(`secret`),
		Postfix:     nil,
		ExpireTime:  time.Second * 1,
		HashType:    sha256.New,
		DisableLogs: true,
	})

	token := converter.NewTokenWithExpire([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token with expire time is nil")
	}

	time.Sleep(time.Second * 3)

	_, err := converter.ParseTokenWithExpire(token)
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
		SecretKey:   []byte(`secret`),
		Postfix:     []byte(`postfix`),
		ExpireTime:  time.Minute * 5,
		HashType:    sha256.New,
		DisableLogs: true,
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
		SecretKey:   []byte(`secret`),
		Postfix:     []byte(`postfix`),
		ExpireTime:  time.Second * 1,
		HashType:    sha256.New,
		DisableLogs: true,
	})

	token := converter.NewTokenWithExpire([]byte(`token`))
	if token == "" || len(token) == 0 {
		t.Error("Token with expire time and postfix is nil")
	}

	time.Sleep(time.Second * 3)

	_, err := converter.ParseTokenWithExpire(token)
	if err == nil {
		t.Error("Token with expire time and postfix parse err: ", "token is not expired!")
	} else {
		if err != TokenExpired {
			t.Error("Token with expire time and postfix parse err: ", err)
		}
	}
}
