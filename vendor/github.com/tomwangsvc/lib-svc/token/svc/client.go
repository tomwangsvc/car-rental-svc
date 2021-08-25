package svc

import (
	"context"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	lib_certificates "github.com/tomwangsvc/lib-svc/certificates"
	lib_env "github.com/tomwangsvc/lib-svc/env"
	lib_errors "github.com/tomwangsvc/lib-svc/errors"
	lib_log "github.com/tomwangsvc/lib-svc/log"
	lib_secrets "github.com/tomwangsvc/lib-svc/secrets"
	lib_svc "github.com/tomwangsvc/lib-svc/svc"
	lib_token "github.com/tomwangsvc/lib-svc/token"
)

type Client interface {
	NewTokenWithCustomClaims(ctx context.Context, customClaims jwt.Claims) (string, error)
	NewClaims(ctx context.Context, audience string, expiry time.Duration) Claims

	NewToken(ctx context.Context, audience string) (string, error)
	NewTokenWithExpiry(ctx context.Context, audience string, expiry time.Duration) (string, error)
	VerifyJwtToken(jwtToken string) (*jwt.Token, error)
	VerifyTokenAndExtractClaims(ctx context.Context, jwtToken string) (*Claims, error)
	ExtractAndVerifyStandardClaims(tokenClaims jwt.MapClaims) (*jwt.StandardClaims, error)
}

type Config struct {
	Env lib_env.Env
}

func NewClient(ctx context.Context, config Config, certificatesClient lib_certificates.Client, secretsClient lib_secrets.Client) (Client, error) {
	lib_log.Info(ctx, "Initializing", lib_log.FmtAny("config", config))

	apiKey, err := secretsClient.ValueFromBase64AsBytes(config.Env.SvcId, secretTypeKey)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed getting value from base 64 as bytes")
	}

	lib_log.Info(ctx, "Initialized")
	return client{
		apiKey:             apiKey,
		config:             config,
		certificatesClient: certificatesClient,
	}, nil
}

func RequiredCertificates(config Config) lib_certificates.Required {
	if config.Env.SvcId == lib_svc.IamId { // IAM could potentially be called by any service using that service's token, so needs all certs to verify those tokens
		return lib_svc.Services
	}
	return []string{config.Env.SvcId} // Any service may be called by itself from an internal queue which uses an api token for authorization
}

const (
	secretTypeKey = "key"
)

func RequiredSecrets(config Config) lib_secrets.Required {
	return lib_secrets.Required{
		config.Env.SvcId: []string{
			secretTypeKey,
		},
	}
}

type client struct {
	config             Config
	certificatesClient lib_certificates.Client
	apiKey             []byte
}

const (
	TokenExpiry           = 10 * time.Minute
	tokenTimeSyncVariance = -15 * time.Second
)

func (c client) NewToken(ctx context.Context, audience string) (string, error) {
	return c.NewTokenWithExpiry(ctx, audience, TokenExpiry)
}

func (c client) NewTokenWithExpiry(ctx context.Context, audience string, expiry time.Duration) (string, error) {
	return c.NewTokenWithCustomClaims(ctx, newClaims(audience, c.config.Env.Id, c.config.Env.SvcId, c.config.Env.RuntimeId, time.Now(), expiry))
}

type Claims struct {
	ClientId   string `json:"client_id,omitempty"`   // The ID of the client for which the token is generated
	ClientName string `json:"client_name,omitempty"` // The name of the client for which the token is generated
	Env        string `json:"env,omitempty"`         // The environment (dev|prd|stg|uat)
	jwt.StandardClaims
}

func (c Claims) Valid() error {
	return c.StandardClaims.Valid()
}

func (c client) NewClaims(ctx context.Context, audience string, expiry time.Duration) Claims {
	lib_log.Info(ctx, "Generating", lib_log.FmtString("audience", audience), lib_log.FmtDuration("expiry", expiry))

	claims := newClaims(audience, c.config.Env.Id, c.config.Env.SvcId, c.config.Env.RuntimeId, time.Now(), expiry)

	lib_log.Info(ctx, "Generated", lib_log.FmtAny("claims", claims))
	return claims
}

func newClaims(audience, env, SvcId, runtimeId string, now time.Time, expiry time.Duration) Claims {
	validFrom := now.Add(tokenTimeSyncVariance).Unix() // Has to be set back in time to avoid time sync issues; IssuedAt and NotBefore are both set to the same value so that the jwt-go validity checker does not throw an error
	return Claims{
		ClientId:   runtimeId,
		ClientName: SvcId,
		Env:        env,
		StandardClaims: jwt.StandardClaims{
			Audience:  audience,
			ExpiresAt: now.Add(expiry).Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  validFrom,
			Issuer:    SvcId,
			NotBefore: validFrom,
		},
	}
}

func (c client) NewTokenWithCustomClaims(ctx context.Context, customClaims jwt.Claims) (string, error) {
	lib_log.Info(ctx, "Creating", lib_log.FmtAny("customClaims", customClaims))

	key, err := jwt.ParseRSAPrivateKeyFromPEM(c.apiKey)
	if err != nil {
		return "", lib_errors.Wrap(err, "Failed parsing rsa private key from pem")
	}
	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, customClaims).SignedString(key)
	if err != nil {
		return "", lib_errors.Wrap(err, "Failed creating and signing token")
	}

	lib_log.Info(ctx, "Created", lib_log.FmtAny("lib_token.Redact(signedToken)", lib_token.Redact(signedToken)))
	return signedToken, nil
}

func (c client) VerifyTokenAndExtractClaims(ctx context.Context, jwtToken string) (*Claims, error) {
	lib_log.Info(ctx, "Generating", lib_log.FmtString("lib_token.Redact(jwtToken)", lib_token.Redact(jwtToken)))

	token, err := c.VerifyJwtToken(jwtToken)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed verifying jwt token")
	}
	verifiedClaims, err := c.extractAndVerifyClaims(c.config.Env.Id, c.config.Env.SvcId, *token)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed extracting claims from token")
	}

	lib_log.Info(ctx, "Generated", lib_log.FmtAny("verifiedClaims", verifiedClaims))
	return verifiedClaims, nil
}

func (c client) VerifyJwtToken(jwtToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, lib_errors.Errorf("Invalid signing method: %v", token.Header["alg"])
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, lib_errors.Errorf("Claims not recognized: %v", token.Claims)
		}
		if err := claims.Valid(); err != nil {
			v, _ := err.(*jwt.ValidationError)
			if v.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, lib_errors.Wrapf(err, "Claims not valid due to expired: %v", claims)
			}
			return nil, lib_errors.Wrapf(err, "Claims not valid: %v", claims)
		}
		v, ok := claims["iss"].(string)
		if !ok {
			return nil, lib_errors.Errorf("Claim 'iss' not provided in token: %v", claims)
		}
		cert, err := c.certificatesClient.ForIssuer(v)
		if err != nil {
			return nil, lib_errors.Wrap(err, "Failed reading cert for claim")
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM(cert)
		if err != nil {
			return nil, lib_errors.Wrap(err, "Failed parsing RSA public key from PEM")
		}
		return key, nil
	})
	if err != nil {
		return nil, lib_errors.Wrapf(err, "Error parsing token: %s", jwtToken)
	}
	return token, nil
}

//revive:disable:cyclomatic
func (c client) extractAndVerifyClaims(env, svcId string, token jwt.Token) (*Claims, error) {
	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, lib_errors.Errorf("Claims not recognized: %v", token.Claims)
	}
	clientId, ok := tokenClaims["client_id"].(string)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'client_id' not provided in token: %v", tokenClaims)
	}
	clientName, ok := tokenClaims["client_name"].(string)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'client_name' not provided in token: %v", tokenClaims)
	}
	tokenEnv, ok := tokenClaims["env"].(string)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'env' not provided in token: %v", tokenClaims)
	} else if strings.ToLower(tokenEnv) != strings.ToLower(env) {
		return nil, lib_errors.Errorf("Claim 'env' %q does not match expected environment id %q", tokenEnv, env)
	}

	standardClaims, err := c.ExtractAndVerifyStandardClaims(tokenClaims)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed extracting standard claims from token")
	} else if standardClaims.Audience != svcId {
		return nil, lib_errors.Errorf("Claim 'aud' %q does not match expected aud %q", standardClaims.Audience, svcId)
	}

	return &Claims{
		ClientId:       clientId,
		ClientName:     clientName,
		Env:            tokenEnv,
		StandardClaims: *standardClaims,
	}, nil
	//revive:enable:cyclomatic
}

func (c client) ExtractAndVerifyStandardClaims(tokenClaims jwt.MapClaims) (*jwt.StandardClaims, error) {
	audience, ok := tokenClaims["aud"].(string)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'aud' not provided in token: %v", tokenClaims)
	}
	exp, ok := tokenClaims["exp"].(float64)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'exp' not provided in token: %v", tokenClaims)
	}
	jti, ok := tokenClaims["jti"].(string)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'jti' not provided in token: %v", tokenClaims)
	}
	iat, ok := tokenClaims["iat"].(float64)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'iat' not provided in token: %v", tokenClaims)
	}
	issuer, ok := tokenClaims["iss"].(string)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'iss' not provided in token: %v", tokenClaims)
	}
	nbf, ok := tokenClaims["nbf"].(float64)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'nbf' not provided in token: %v", tokenClaims)
	}
	sub, ok := tokenClaims["sub"].(string)
	return &jwt.StandardClaims{
		Audience:  audience,
		ExpiresAt: int64(exp),
		Id:        jti,
		IssuedAt:  int64(iat),
		Issuer:    issuer,
		NotBefore: int64(nbf),
		Subject:   sub,
	}, nil
}
