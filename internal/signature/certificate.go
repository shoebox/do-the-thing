package signature

import (
	"context"
	"crypto/x509"
	"dothething/internal/config"
	"dothething/internal/util"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/crypto/pkcs12"
)

// P12Certificate is a more convenient alias
type P12Certificate struct {
	*x509.Certificate
	FilePath string
}

var (
	ErrorFailedToReadFile   = errors.New("Failed to read file")
	ErrorFailedToDecryptPEM = errors.New("Failed to decrypt PEM")
)

// CertificateService service interface definition
type CertificateService interface {
	DecodeCertificate(r io.Reader, password string) (P12Certificate, error)
	ResolveInFolder(ctx context.Context, root string) []P12Certificate
}

type certService struct {
	config config.Config
	fs     util.FileService
}

// NewCertificateService create a new instance of the certificate service
func NewCertificateService(cfg config.Config, fs util.FileService) CertificateService {
	return certService{config: cfg, fs: fs}
}

// DecDecodeCertificate allow to validate/decode the content of the io.Reader into a P12Certificate
func (xs certService) DecodeCertificate(r io.Reader, password string) (P12Certificate, error) {
	var result P12Certificate

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
				result = P12Certificate{Certificate: cert}
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

func (xs certService) readCertificateFile(
	ctx context.Context,
	path string,
	res chan P12Certificate,
) error {
	// Open the file content to a reader
	f, err := xs.fs.Open(path)
	if err != nil {
		return err
	}

	// Defer close the file
	defer f.Close()

	return xs.decodeRawCertificate(ctx, f, res, path)
}

func (xs certService) decodeRawCertificate(
	ctx context.Context,
	reader io.Reader,
	cres chan P12Certificate,
	filepath string,
) error {
	// Decodes it (temp test password)
	cert, err := xs.DecodeCertificate(reader, xs.config.CodeSignOption.CertificatePassword)
	if err != nil {
		return err
	}

	cert.FilePath = filepath

	// And finally select the action for result
	select {
	case cres <- cert: // Populate the certs channel with the valid certificate
	case <-ctx.Done(): // In case of context cancelation
		return ctx.Err()
	}

	return nil
}

// ResolveInFolder is used to resolve all certificate files into the provided path
func (xs certService) ResolveInFolder(ctx context.Context, root string) []P12Certificate {
	certs := make(chan P12Certificate)

	// Use the file service walk helper to find candidates certificate file
	errgroup := xs.fs.Walk(ctx, root, isCertificateFile,
		// And for each candidate
		func(ctx context.Context, path string) error {
			// We try to read a certificate file
			return xs.readCertificateFile(ctx, path, certs)
		})

	// Waiting for all the sub goroutines to complete
	go func() {
		// For the errgroup
		errgroup.Wait()

		// And we finally close the certs channel
		close(certs)
	}()

	// Finally we convert the channel items to an array
	var res []P12Certificate
	for c := range certs {
		res = append(res, c)
	}

	// Check whether any of the goroutines failed. Since the error group is accumulating the
	// errors, we don't need to send them (or check for them) in the individual
	// results sent on the channel.
	if err := errgroup.Wait(); err != nil {
		return res
	}

	return res
}
