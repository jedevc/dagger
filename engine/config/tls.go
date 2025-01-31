package config

import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

type TLSConfig struct {
	RootCAs []string `json:"ca"`

	KeyPairs     []TLSKeyPair `json:"keypair"`
	TLSConfigDir []string     `json:"tlsconfigdir"`

	SkipVerify bool `json:"skipVerify"`
}

type TLSKeyPair struct {
	Key         string `json:"key"`
	Certificate string `json:"cert"`
}

func (cfg TLSConfig) ToConfig() (*tls.Config, error) {
	rootCAs := append([]string{}, cfg.RootCAs...)
	keyPairs := append([]TLSKeyPair{}, cfg.KeyPairs...)
	for _, d := range cfg.TLSConfigDir {
		fs, err := os.ReadDir(d)
		if err != nil && !errors.Is(err, os.ErrNotExist) && !errors.Is(err, os.ErrPermission) {
			return nil, errors.WithStack(err)
		}
		for _, f := range fs {
			if strings.HasSuffix(f.Name(), ".crt") {
				rootCAs = append(rootCAs, filepath.Join(d, f.Name()))
			}
			if strings.HasSuffix(f.Name(), ".cert") {
				keyPairs = append(keyPairs, TLSKeyPair{
					Certificate: filepath.Join(d, f.Name()),
					Key:         filepath.Join(d, strings.TrimSuffix(f.Name(), ".cert")+".key"),
				})
			}
		}
	}

	tc := &tls.Config{}
	if len(rootCAs) > 0 {
		systemPool, err := x509.SystemCertPool()
		if err != nil {
			if runtime.GOOS == "windows" {
				systemPool = x509.NewCertPool()
			} else {
				return nil, errors.Wrapf(err, "unable to get system cert pool")
			}
		}
		tc.RootCAs = systemPool
	}

	for _, p := range rootCAs {
		dt, err := os.ReadFile(p)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read %s", p)
		}
		tc.RootCAs.AppendCertsFromPEM(dt)
	}

	for _, kp := range keyPairs {
		cert, err := tls.LoadX509KeyPair(kp.Certificate, kp.Key)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to load keypair for %s", kp.Certificate)
		}
		tc.Certificates = append(tc.Certificates, cert)
	}

	if cfg.SkipVerify {
		tc.InsecureSkipVerify = true
	}

	return tc, nil
}
