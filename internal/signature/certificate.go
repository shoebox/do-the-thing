package signature

import (
	"context"
	"crypto/x509"
	"dothething/internal/api"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/pkcs12"
)

var (
	ErrorFailedToReadFile   = errors.New("Failed to read file")
	ErrorFailedToDecryptPEM = errors.New("Failed to decrypt PEM")
)

// CertificateService service interface definition
type certService struct {
	api.API
	*api.Config
}

// NewCertificateService create a new instance of the certificate service
func NewCertificateService(api api.API, cfg *api.Config) api.CertificateService {
	return certService{api, cfg}
}

// DecDecodeCertificate allow to validate/decode the content of the io.Reader into a P12Certificate
func (xs certService) DecodeCertificate(r io.Reader, password string) (api.P12Certificate, error) {
	var result api.P12Certificate

	// We read the content of the file
	data, err := readFile(r)
	if err != nil {
		return result, err
	}

	// And use pkcs12 API to convert all PEM safe bags to PEM blocks
	blocks, err := pkcs12.ToPEM(data, password)
	if err != nil {
		return result, ErrorFailedToDecryptPEM
	}

	// We iterate on all blocks
	for _, key := range blocks {
		// We check for certificate type
		if key.Type == "CERTIFICATE" {
			// And we parse the single certificate
			if cert, err := x509.ParseCertificate(key.Bytes); err == nil {
				result = api.P12Certificate{Certificate: cert}
				break
			}
		}
	}

	return result, nil
}

// readFile is a convenient helper to read the content of a reader and returns it's content or an error
func readFile(r io.Reader) ([]byte, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, ErrorFailedToReadFile
	}

	return data, err
}

// isCertificateFile do the file is a validate p12 certificate file
func isCertificateFile(info os.FileInfo) bool {
	return info.Mode().IsRegular() && // Is it a regular file
		!info.IsDir() && // And not a directory
		strings.HasSuffix(info.Name(), ".p12") // And has the right extension
}

func (xs certService) readCertificateFile(path string) (*api.P12Certificate, error) {
	// Open the file content to a reader
	f, err := xs.API.FileService().Open(path)
	if err != nil {
		return nil, err
	}

	// Defer close the file
	defer f.Close()

	// Decodes it (temp test password)
	cert, err := xs.DecodeCertificate(f, xs.Config.CodeSignOption.CertificatePassword)
	if err != nil {
		return nil, err
	}

	cert.FilePath = path

	return &cert, nil
}

func (xs certService) worker(wg *sync.WaitGroup, paths <-chan string, out *[]*api.P12Certificate) {
	wg.Add(1)
	for value := range paths {
		c, err := xs.readCertificateFile(value)
		if err != nil {
			fmt.Println("error", err)
		}

		*out = append(*out, c)
	}
	wg.Done()
}

// ResolveInFolder is used to resolve all certificate files into the provided path
func (xs certService) ResolveInFolder(ctx context.Context, root string) []*api.P12Certificate {
	var res []*api.P12Certificate
	var wg sync.WaitGroup
	// Increment waitgroup counter and create go routines
	paths := make(chan string)
	for i := 0; i < 8; i++ {
		go xs.worker(&wg, paths, &res)
	}

	err := xs.API.FileService().Walk(ctx, root, isCertificateFile, paths, &wg)
	if err != nil {
		fmt.Println("err", err)
	}

	wg.Wait()

	return res
}
