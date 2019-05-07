// Copyright 2015,2016,2017,2018 SeukWon Kang (kasworld@gmail.com)

package weblib

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func Use(
	h http.HandlerFunc,
	middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {

	for _, m := range middleware {
		h = m(h)
	}
	return h
}

type NoAuth struct {
}

func NewNoAuth() *NoAuth {
	return &NoAuth{}
}

func (na *NoAuth) Auth(h http.HandlerFunc) http.HandlerFunc {
	return h
	// return h.ServeHTTP
}

type BasicAuth struct {
	realm    string
	name     string
	password string
	realmstr string
}

func NewBasicAuth(realm string, name string, password string) *BasicAuth {
	return &BasicAuth{
		realm:    realm,
		name:     name,
		password: password,
		realmstr: `Basic realm="` + realm + `"`,
	}
}

// Leverages nemo's answer in http://stackoverflow.com/a/21937924/556573
func (ba *BasicAuth) Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", ba.realmstr)
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}
		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}
		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}
		if pair[0] != ba.name || pair[1] != ba.password {
			http.Error(w, "Not authorized", 401)
			return
		}
		h(w, r)
		// h.ServeHTTP(w, r)
	}
}
