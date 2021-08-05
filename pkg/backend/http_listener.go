package backend

import (
	"crypto/tls"
	"net"

	"github.com/valyala/fasthttp"

	"dynexo.de/pkg/log"
)

const (
	certificateFile = `-----BEGIN CERTIFICATE-----
MIICkjCCAhegAwIBAgIEXLo8PzAKBggqhkjOPQQDAjBAMRswGQYDVQQDDBJ3aXhD
b25uZWN0IFJvb3QgQ0ExFDASBgNVBAoMC2R5bmV4byBHbWJIMQswCQYDVQQGEwJE
RTAeFw0xOTA0MTkyMTIzMTFaFw0yOTA0MTYyMTIzMTFaMF8xFzAVBgNVBAMMDndp
eENPTk5FQ1QgQVBJMRQwEgYDVQQKDAtkeW5leG8gR21iSDETMBEGA1UECwwKd2l4
Q29ubmVjdDEMMAoGA1UECwwDQVBJMQswCQYDVQQGEwJERTB2MBAGByqGSM49AgEG
BSuBBAAiA2IABJF+xLIP1IJI77Vas5FjkuZkanq5Ciz7+wgh0pgEtuQmwtu0Glfp
QSgTo8oTooKgAh8PfF6QUUdtPjoXfJUWrwifpm1KHtxy7R4y/QAKBcqVaCE/BWLP
kT9AeUzijhT8/KOBwjCBvzAJBgNVHRMEAjAAMAsGA1UdDwQEAwIDqDAdBgNVHQ4E
FgQU/lY5GWg5jclDZO6QaunBCNjR9x4wewYDVR0jBHQwcoAUAJGIOErPEtyYCqD0
JnFrBaifFjChRKRCMEAxGzAZBgNVBAMMEndpeENvbm5lY3QgUm9vdCBDQTEUMBIG
A1UECgwLZHluZXhvIEdtYkgxCzAJBgNVBAYTAkRFghQvNSX4/2AI6XTQF2RvSvZs
Z6u1ITAJBgNVHRIEAjAAMAoGCCqGSM49BAMCA2kAMGYCMQCWnhX5ClGiZWVn/mCc
kt/espV0wETIj8dLARH0XLSXpQ9lGTuC2hfOb49Nq1fGNHICMQDE+DI+W4eBSp7g
QCGQFqEVZpsGvymDvWefnH9QQWwU6iC83eO0DQBwF7enIxxNFDA=
-----END CERTIFICATE-----`

	certificateKey = `-----BEGIN EC PRIVATE KEY-----
MIGkAgEBBDC5hUQrSkpMvZf9YNiaxEIOwv53zyp3DxC6ALi8P4aK5wvSrfFQDHxI
Wc0+27PBWYSgBwYFK4EEACKhZANiAASRfsSyD9SCSO+1WrORY5LmZGp6uQos+/sI
IdKYBLbkJsLbtBpX6UEoE6PKE6KCoAIfD3xekFFHbT46F3yVFq8In6ZtSh7ccu0e
Mv0ACgXKlWghPwViz5E/QHlM4o4U/Pw=
-----END EC PRIVATE KEY-----`
)

type (
	WrappedHttpLogger struct {
		logger log.ILogger
	}
)

func (l WrappedHttpLogger) Printf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func ListenAndServeFast(address string, requestHandler fasthttp.RequestHandler, logger log.ILogger) error {

	// load internal api cert and key
	cert, err := tls.X509KeyPair([]byte(certificateFile), []byte(certificateKey))
	if err != nil {
		return err
	}

	// setup TLS connection
	tlsConfig := &tls.Config{
		Certificates:             []tls.Certificate{cert},
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
	}

	tlsConfig.BuildNameToCertificate()

	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	lnTls := tls.NewListener(ln, tlsConfig)

	logger.Infof("Starting Service Listener on %s", address)

	server := &fasthttp.Server{
		Handler: requestHandler,
		Name:    "ufyler",
		Logger:  WrappedHttpLogger{logger},
		//Concurrency: 8 * 1024,

	}

	return server.Serve(lnTls)
}
