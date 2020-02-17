package signature

import (
	"crypto/x509"
	"errors"
	"io"
	"io/ioutil"

	"golang.org/x/crypto/pkcs12"
)

var (
	ErrorFailedToReadFile   = errors.New("Failed to read file")
	ErrorFailedToDecryptPEM = errors.New("Failed to decrypt PEM")
)

func DecodeCertificate(r io.Reader, password string) (*x509.Certificate, error) {
	data, err := readFile(r)
	if err != nil {
		return nil, err
	}

	blocks, err := pkcs12.ToPEM(data, password)
	if err != nil {
		return nil, ErrorFailedToDecryptPEM
	}

	var result *x509.Certificate

	for _, key := range blocks {
		if key.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(key.Bytes)
			if err == nil {
				result = cert
				break
			}
		}
	}

	return result, nil
}

func readFile(r io.Reader) ([]byte, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, ErrorFailedToReadFile
	}

	return data, err
}
