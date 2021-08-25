package certificates

import (
	"context"
	"fmt"

	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_storage "github.com/tomwangsvc/lib-svc/storage"
)

type Client interface {
	ForIssuer(issuer string) ([]byte, error)
}

type Config struct {
	BucketName string
	Env        lib_env.Env
	Required   Required
}

func NewClient(ctx context.Context, config Config, storageClient lib_storage.Client) (Client, error) {
	lib_log.Info(ctx, "Initializing", lib_log.FmtAny("config", config))

	if config.BucketName == "" {
		return nil, lib_errors.New("BucketName in config not supplied, please add a bucket name for correction configuration")
	}

	c := client{
		config:              config,
		certificateByIssuer: make(map[string][]byte),
	}
	if err := c.downloadAndCacheCertificates(ctx, storageClient); err != nil {
		return nil, lib_errors.Wrap(err, "Failed downloading and caching certificates")
	}

	lib_log.Info(ctx, "Initialized")
	return &c, nil
}

type Required []string // List of svc ids for issuers for which certificates are required

func ReduceRequired(required ...Required) Required {
	var reduced Required
	for i := 0; i < len(required); i++ {
		for _, newV := range required[i] {
			var exists bool
			for _, reducedV := range reduced {
				if newV == reducedV {
					exists = true
					break
				}
			}
			if !exists {
				reduced = append(reduced, newV)
			}
		}
	}
	return reduced
}

type client struct {
	certificateByIssuer map[string][]byte
	config              Config
}

func (c client) ForIssuer(issuer string) ([]byte, error) {
	if v, ok := c.certificateByIssuer[issuer]; ok {
		return v, nil
	}
	return nil, lib_errors.Errorf("No cert for %q", issuer)
}

func (c *client) downloadAndCacheCertificates(ctx context.Context, storageClient lib_storage.Client) error {
	bucketName := c.config.BucketName

	type certificateAndIssuerAndError struct {
		Err         error
		Issuer      string
		Certificate []byte
	}
	ch := make(chan certificateAndIssuerAndError)
	for _, issuer := range c.config.Required {
		go func(issuer string) {
			if issuer == "" {
				ch <- certificateAndIssuerAndError{
					Err: lib_errors.Errorf("Invalid required certificate issuer %q", issuer),
				}
				return
			}
			bucketKey := c.certificateFileName(issuer)
			certificate, _, err := storageClient.ReadFromBucket(ctx, bucketName, bucketKey)
			if err != nil {
				ch <- certificateAndIssuerAndError{
					Err: lib_errors.Wrap(err, "Failed reading certificate from bucket"),
				}
				return
			}

			ch <- certificateAndIssuerAndError{
				Issuer:      issuer,
				Certificate: certificate,
			}
		}(issuer)
	}

	certificateByIssuer := make(map[string][]byte)
	var err error
	for i := 0; i < len(c.config.Required); i++ {
		certificateAndIssuerAndError := <-ch
		if certificateAndIssuerAndError.Err != nil {
			lib_log.Error(ctx, "Failed downloading certificate in goroutine", lib_log.FmtError(certificateAndIssuerAndError.Err))
			err = lib_errors.Wrap(certificateAndIssuerAndError.Err, "Failed downloading certificate in goroutine")
		}
		certificateByIssuer[certificateAndIssuerAndError.Issuer] = certificateAndIssuerAndError.Certificate
	}
	if err != nil {
		return lib_errors.Wrap(err, "Failed downloading one or more certificates")
	}
	c.certificateByIssuer = certificateByIssuer

	return nil
}

func (c client) certificateFileName(issuer string) string {
	return fmt.Sprintf("%s_cert_%s.pem", issuer, c.config.Env.Id)
}
