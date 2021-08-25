package secrets

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_storage "github.com/tomwangsvc/lib-svc/storage"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

type Client interface {
	ValueFromBase64(secretDomain, secretType string) (string, error)
	ValueFromBase64AsBytes(secretDomain, secretType string) ([]byte, error)
	ValueFromBase64WithNewLinesStripped(secretDomain, secretType string) (string, error)
	Value(secretDomain, secretType string) (string, error)
}

type Config struct {
	BucketName string
	Env        lib_env.Env
	Required
}

func NewClient(ctx context.Context, config Config, storageClient lib_storage.Client) (Client, error) {
	lib_log.Info(ctx, "Initializing", lib_log.FmtAny("config", config))

	if config.BucketName == "" {
		return nil, lib_errors.New("BucketName in config not supplied, please add a bucket name for correction configuration")
	}
	c := client{
		config:            config,
		secretBySecretKey: make(map[string]string),
	}
	if err := c.downloadAndDecryptAndCacheSecrets(ctx, storageClient); err != nil {
		return nil, lib_errors.Wrap(err, "Failed downloading, decrypting, and caching secrets")
	}

	lib_log.Info(ctx, "Initialized")
	return c, nil
}

type Required map[string][]string // Map of secret domain to list of secret types for that domain

func ReduceRequired(required ...Required) Required {
	reduced := make(Required)
	for i := 0; i < len(required); i++ {
		for newK, newV := range required[i] {
			newV = reduceRequiredTypes(nil, newV)
			if reducedV, ok := reduced[newK]; ok {
				reduced[newK] = reduceRequiredTypes(reducedV, newV)
				continue
			}
			reduced[newK] = newV
		}
	}
	return reduced
}

func reduceRequiredTypes(reduced, new []string) []string {
	for _, newv := range new {
		var exists bool
		for _, reducedV := range reduced {
			if newv == reducedV {
				exists = true
				break
			}
		}
		if !exists {
			reduced = append(reduced, newv)
		}
	}
	return reduced
}

type client struct {
	config            Config
	secretBySecretKey map[string]string
}

func (c client) Value(secretDomain, secretType string) (string, error) {
	k := secretKey(secretDomain, secretType)
	v, ok := c.secretBySecretKey[k]
	if !ok {
		return "", lib_errors.Errorf("No secret for domain %q and type %q", secretDomain, secretType)
	}
	return v, nil
}

func secretKey(secretDomain, secretType string) string {
	return secretDomain + secretType
}

func (c client) ValueFromBase64AsBytes(secretDomain, secretType string) ([]byte, error) {
	v, err := c.Value(secretDomain, secretType)
	if err != nil {
		return nil, err
	}
	buf, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed decoding base64 plaintext to byte array")
	}
	return buf, nil
}

func (c client) ValueFromBase64(secretDomain, secretType string) (string, error) {
	v, err := c.ValueFromBase64AsBytes(secretDomain, secretType)
	if err != nil {
		return "", err
	}
	return string(v), nil
}

func (c client) ValueFromBase64WithNewLinesStripped(secretDomain, secretType string) (string, error) {
	v, err := c.ValueFromBase64(secretDomain, secretType)
	if err != nil {
		return "", err
	}
	return strings.Replace(v, "\n", "", -1), nil
}

func (c *client) downloadAndDecryptAndCacheSecrets(ctx context.Context, storageClient lib_storage.Client) error {
	googleClient, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return lib_errors.Wrap(err, "Failed creating Google client")
	}
	kmsClient, err := cloudkms.New(googleClient)
	if err != nil {
		return lib_errors.Wrap(err, "Failed creating KMS client")
	}
	bucketName := c.config.BucketName

	type secretAndSecretKeyAndError struct {
		Err       error
		Secret    string
		SecretKey string
	}
	ch := make(chan secretAndSecretKeyAndError)
	for secretDomain, secretTypes := range c.config.Required {
		for _, secretType := range secretTypes {
			go func(secretDomain, secretType string) {
				if secretDomain == "" {
					ch <- secretAndSecretKeyAndError{
						Err: lib_errors.Errorf("Invalid required secret domain %q", secretDomain),
					}
					return
				}
				if secretType == "" {
					ch <- secretAndSecretKeyAndError{
						Err: lib_errors.Errorf("Invalid required secret type %q", secretType),
					}
					return
				}
				bucketKey := c.secretFileName(secretDomain, secretType)
				object, _, err := storageClient.ReadFromBucket(ctx, bucketName, bucketKey)
				if err != nil {
					ch <- secretAndSecretKeyAndError{
						Err: lib_errors.Wrap(err, "Failed reading object from bucket"),
					}
					return
				}

				var secret struct {
					Ciphertext string `json:"ciphertext"`
				}
				err = json.Unmarshal(object, &secret)
				if err != nil {
					ch <- secretAndSecretKeyAndError{
						Err: lib_errors.Wrap(err, "Failed unmarshalling secret"),
					}
					return
				}
				plaintext, err := c.decrypt(secret.Ciphertext, kmsClient)
				if err != nil {
					ch <- secretAndSecretKeyAndError{
						Err: lib_errors.Wrap(err, fmt.Sprintf("Failed decrypting: secretDomain=%s, secretType=%s, bucketName=%s, bucketKey=%s", secretDomain, secretType, bucketName, bucketKey)),
					}
					return
				}

				ch <- secretAndSecretKeyAndError{
					SecretKey: secretKey(secretDomain, secretType),
					Secret:    plaintext,
				}
			}(secretDomain, secretType)
		}
	}

	secretBySecretKey := make(map[string]string)
	for _, secretTypes := range c.config.Required {
		for range secretTypes {
			secretAndSecretKeyAndError := <-ch
			if secretAndSecretKeyAndError.Err != nil {
				lib_log.Error(ctx, "Failed downloading secret in goroutine", lib_log.FmtError(secretAndSecretKeyAndError.Err))
				err = lib_errors.Wrap(secretAndSecretKeyAndError.Err, "Failed downloading secret in goroutine")
			}
			secretBySecretKey[secretAndSecretKeyAndError.SecretKey] = secretAndSecretKeyAndError.Secret
		}
	}
	if err != nil {
		return lib_errors.Wrap(err, "Failed downloading one or more secrets")
	}
	c.secretBySecretKey = secretBySecretKey

	return nil
}

func (c client) secretFileName(secretDomain, secretType string) string {
	return fmt.Sprintf("%s_%s_cloudkms-%s.json", secretDomain, secretType, c.config.Env.Id)
}

func (c client) decrypt(ciphertext string, kmsClient *cloudkms.Service) (string, error) {
	keyRing := c.config.Env.KeyRing()
	resp, err := kmsClient.Projects.Locations.KeyRings.CryptoKeys.Decrypt(keyRing, &cloudkms.DecryptRequest{Ciphertext: ciphertext}).Do()
	if err != nil {
		return "", lib_errors.Wrapf(err, "Failed to decrypt ciphertext using key ring '%s'", keyRing)
	}
	return resp.Plaintext, nil
}
