package iam

import (
	"context"
	"encoding/json"
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
	NewToken(ctx context.Context, tokenData TokenData) (string, error)

	ExtractUnverifiedClaims(ctx context.Context, jwtToken string) (unverifiedClaims *Claims, err error)
	VerifyTokenAndExtractClaims(ctx context.Context, jwtToken string) (*Claims, error)
}

type Config struct {
	Env lib_env.Env
}

func NewClient(ctx context.Context, config Config, certificatesClient lib_certificates.Client, secretsClient lib_secrets.Client) (Client, error) {
	lib_log.Info(ctx, "Initializing", lib_log.FmtAny("config", config))

	var apiKey []byte
	if requiresTokenCreation(config) {
		var err error
		apiKey, err = secretsClient.ValueFromBase64AsBytes(config.Env.SvcId, secretTypeKey)
		if err != nil {
			return nil, lib_errors.Wrap(err, "Failed getting value from base 64 as bytes")
		}
	}

	lib_log.Info(ctx, "Initialized")
	return client{
		apiKey:             apiKey,
		config:             config,
		certificatesClient: certificatesClient,
	}, nil
}

func RequiredCertificates() lib_certificates.Required {
	return lib_certificates.Required{lib_svc.IamId}
}

const (
	secretTypeKey = "key"
)

func RequiredSecrets(config Config) lib_secrets.Required {
	if requiresTokenCreation(config) {
		return lib_secrets.Required{
			lib_svc.IamId: []string{
				secretTypeKey,
			},
		}
	}
	return nil
}

func requiresTokenCreation(config Config) bool {
	return config.Env.SvcId == lib_svc.IamId
}

type client struct {
	apiKey             []byte
	config             Config
	certificatesClient lib_certificates.Client
}

const (
	tokenExpiry           = 10 * time.Minute
	tokenTimeSyncVariance = -15 * time.Second
)

type TokenData struct {
	ClientId           string
	ClientName         string
	Email              string
	IdentityExternalId string
	IdentityId         string
	IdentityProvider   string
	Name               string
	Roles              []string
	Test               bool
}

func (c client) NewToken(ctx context.Context, tokenData TokenData) (string, error) {
	lib_log.Info(ctx, "Creating", lib_log.FmtAny("tokenData", tokenData))

	if !requiresTokenCreation(c.config) {
		return "", lib_errors.Errorf("Unable to generate new token, %q must be configured as requiring token creation", c.config.Env.SvcId)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(c.apiKey)
	if err != nil {
		return "", lib_errors.Wrap(err, "Failed parsing rsa private key from pem")
	}
	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, newClaims(tokenData, c.config.Env.Id, c.config.Env.SvcId, time.Now())).SignedString(key)
	if err != nil {
		return "", lib_errors.Wrap(err, "Failed creating and signing token")
	}

	lib_log.Info(ctx, "Created", lib_log.FmtAny("lib_token.Redact(signedToken)", lib_token.Redact(signedToken)))
	return signedToken, nil
}

type Claims struct {
	ClientId           string   `json:"client_id,omitempty"`            // The ID of the client for which the token is generated
	ClientName         string   `json:"client_name,omitempty"`          // The name of the client for which the token is generated
	Email              string   `json:"email,omitempty"`                // The email of the user
	Env                string   `json:"env,omitempty"`                  // The environment (dev|prd|stg|uat)
	IdentityExternalId string   `json:"identity_external_id,omitempty"` // The IAM identity external ID
	IdentityId         string   `json:"identity_id,omitempty"`          // The IAM identity ID
	IdentityProvider   string   `json:"identity_provider,omitempty"`    // The IAM identity provider
	Name               string   `json:"name,omitempty"`                 // The name of the user
	Roles              []string `json:"roles,omitempty"`                // The roles of the user
	Test               bool     `json:"test"`                           // Whether or not the token is in a test context
	jwt.StandardClaims
}

func (c Claims) Valid() error {
	return c.StandardClaims.Valid()
}

func newClaims(tokenData TokenData, env, SvcId string, now time.Time) Claims {
	validFrom := now.Add(tokenTimeSyncVariance).Unix() // Has to be set back in time to avoid time sync issues; IssuedAt and NotBefore are both set to the same value so that the jwt-go validity checker does not throw an error
	return Claims{
		ClientId:           tokenData.ClientId,
		ClientName:         tokenData.ClientName,
		Email:              tokenData.Email,
		Env:                env,
		IdentityExternalId: tokenData.IdentityExternalId,
		IdentityId:         tokenData.IdentityId,
		IdentityProvider:   tokenData.IdentityProvider,
		Name:               tokenData.Name,
		Roles:              tokenData.Roles,
		Test:               tokenData.Test,
		StandardClaims: jwt.StandardClaims{
			Audience:  SvcId,
			ExpiresAt: now.Add(tokenExpiry).Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  validFrom,
			Issuer:    SvcId,
			NotBefore: validFrom,
		},
	}
}

func (c client) ExtractUnverifiedClaims(ctx context.Context, jwtToken string) (unverifiedClaims *Claims, err error) {
	lib_log.Info(ctx, "Extracting", lib_log.FmtString("lib_token.Redact(jwtToken)", lib_token.Redact(jwtToken)))

	var token *jwt.Token
	if _, err = jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		token = t
		return nil, nil
	}); err == nil {
		err = lib_errors.Errorf("Token successfully parsed but we DID NOT provide a key. An error was expected because we were decoding the token, not verifying it: %v", token)
		return
	} else if token == nil {
		err = lib_errors.Errorf("Token not recognized: %q", jwtToken)
		return
	}

	tokenClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = lib_errors.Errorf("Claims not recognized: %v", token.Claims)
		return
	}

	tokenClaimsJson, err := json.Marshal(tokenClaims)
	if err != nil {
		err = lib_errors.Wrap(err, "Failed marshalling token claims")
		return
	}
	var claims Claims
	if err = json.Unmarshal(tokenClaimsJson, &claims); err != nil {
		err = lib_errors.Wrap(err, "Failed unmarshalling token claims json into Claims")
		return
	}
	unverifiedClaims = &claims

	lib_log.Info(ctx, "Extracted", lib_log.FmtAny("unverifiedClaims", unverifiedClaims))
	return unverifiedClaims, nil
}

func (c client) VerifyTokenAndExtractClaims(ctx context.Context, jwtToken string) (*Claims, error) {
	lib_log.Info(ctx, "Generating", lib_log.FmtString("lib_token.Redact(jwtToken)", lib_token.Redact(jwtToken)))

	token, err := c.verifyJWTToken(jwtToken)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed verifying jwt token")
	}
	verifiedClaims, err := extractAndVerifyClaims(c.config.Env.Id, *token)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed extracting claims from token")
	}

	lib_log.Info(ctx, "Generated", lib_log.FmtAny("verifiedClaims", verifiedClaims))
	return verifiedClaims, nil
}

func (c client) verifyJWTToken(jwtToken string) (*jwt.Token, error) {
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
func extractAndVerifyClaims(env string, token jwt.Token) (*Claims, error) {
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
	email, _ := tokenClaims["email"].(string) // Optional
	tokenEnv, ok := tokenClaims["env"].(string)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'env' not provided in token: %v", tokenClaims)
	} else if strings.ToLower(tokenEnv) != strings.ToLower(env) {
		return nil, lib_errors.Errorf("Claim 'env' %q does not match expected environment id %q", tokenEnv, env)
	}
	identityExternalId, ok := tokenClaims["identity_external_id"].(string)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'identity_external_id' not provided in token: %v", tokenClaims)
	}
	identityId, ok := tokenClaims["identity_id"].(string)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'identity_id' not provided in token: %v", tokenClaims)
	}
	identityProvider, ok := tokenClaims["identity_provider"].(string)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'identity_provider' not provided in token: %v", tokenClaims)
	}
	name, ok := tokenClaims["name"].(string)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'name' not provided in token: %v", tokenClaims)
	}
	var roles []string
	if r, ok := tokenClaims["roles"].([]interface{}); ok {
		for i, v := range r {
			role, ok := v.(string)
			if !ok {
				return nil, lib_errors.Errorf("Claim 'role' %v not recognized at [%d]=%v", r, i, v)
			}
			roles = append(roles, role)
		}
	}
	test, ok := tokenClaims["test"].(bool)
	if !ok {
		return nil, lib_errors.Errorf("Claim 'test' not provided in token: %v", tokenClaims)
	}
	standardClaims, err := extractAndVerifyStandardClaims(tokenClaims)
	if err != nil {
		return nil, lib_errors.Wrap(err, "Failed extracting standard claims from token")
	} else if standardClaims.Audience != lib_svc.IamId {
		return nil, lib_errors.Errorf("Claim 'aud' %q does not match expected aud %q", standardClaims.Audience, lib_svc.IamId)
	} else if standardClaims.Issuer != lib_svc.IamId {
		return nil, lib_errors.Errorf("Claim 'iss' %q does not match expected iss %q", standardClaims.Issuer, lib_svc.IamId)
	}
	return &Claims{
		ClientId:           clientId,
		ClientName:         clientName,
		Email:              email,
		Env:                tokenEnv,
		IdentityExternalId: identityExternalId,
		IdentityId:         identityId,
		IdentityProvider:   identityProvider,
		Name:               name,
		Roles:              roles,
		StandardClaims:     *standardClaims,
		Test:               test,
	}, nil
	//revive:enable:cyclomatic
}

func extractAndVerifyStandardClaims(tokenClaims jwt.MapClaims) (*jwt.StandardClaims, error) {
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
