// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"html/template"
	"math/big"
	"net/http"
	"time"

	"github.com/kasworld/log/genlog/basiclog"
	"github.com/kasworld/log/logflags"
	"github.com/kasworld/weblib"
)

// var Ver = ""

// func init() {
// 	version.Set(Ver)
// }

func main() {
	basiclog.GlobalLogger.SetFlags(basiclog.GlobalLogger.GetFlags().BitClear(logflags.LF_functionname))
	app := NewApp()
	app.ServiceMain(context.Background())
}

func (app WebApp) String() string {
	return fmt.Sprintf("WebApp[%v]", app.addr)
}

type WebApp struct {
	addr    string
	DoClose func() `webformhide:"" stringformhide:""`
}

func NewApp() *WebApp {
	app := &WebApp{
		addr: ":9002",
	}

	app.DoClose = func() {
		basiclog.Fatal("Too early DoClose call %v", app)
	}

	return app
}

func (app *WebApp) ServiceMain(mainctx context.Context) {
	basiclog.Debug("Start ServiceMain %v", app)
	defer func() { basiclog.Debug("End ServiceMain %v", app) }()

	ctx, closeCtx := context.WithCancel(mainctx)
	app.DoClose = closeCtx
	defer app.DoClose()

	go app.serveWebInfo()

	timerInfoTk := time.NewTicker(time.Second)
	defer timerInfoTk.Stop()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop

		case <-timerInfoTk.C:
		}
	}
}

//////////////////////////////////////////////

func (app *WebApp) webFaviconIco(w http.ResponseWriter, r *http.Request) {
}

func (app *WebApp) serveWebInfo() {
	basiclog.Debug("Start serveWebInfo %v", app)
	defer func() { basiclog.Debug("End serveWebInfo %v", app) }()

	authdata := weblib.NewAuthData("WebApp")
	authdata.ReLoadUserData([][2]string{
		{"root", "password"},
	})
	webMux := weblib.NewAuthMux(authdata, basiclog.GlobalLogger)

	// if !version.IsRelease() {
	// 	webprofile.AddWebProfile(webMux)
	// }

	webMux.HandleFuncAuth("/", app.webMCInfo)
	webMux.HandleFunc("/favicon.ico", app.webFaviconIco)
	authdata.AddAllActionName("root")
	basiclog.Debug("%v", webMux)

	var srv http.Server
	srv.Addr = ":9002"
	srv.Handler = webMux
	// srv.TLSConfig = generateTLSConfig()

	srv.ListenAndServe()
	// srv.ListenAndServeTLS("", "")
	// srv.ListenAndServeTLS("localhost.cert", "localhost.key")
}

func (app *WebApp) webMCInfo(w http.ResponseWriter, r *http.Request) {
	weblib.SetFresh(w, r)
	tplIndex, err := template.New("index").Parse(`
        <html>
        <head>
        <title>Info</title>
        </head>
        <body>
        hello world
        <br/>
        </body>
        </html>
    `)
	if err != nil {
		basiclog.Error("%v", err)
	}
	weblib.SetFresh(w, r)
	if err := tplIndex.Execute(w, app); err != nil {
		basiclog.Error("%v", err)
	}
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
