package provider

import (
	"context"
	"crypto/rsa"
	"os"
	"time"

	"github.com/dylEasydev/go-oauth2-easyclass/db"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/token/jwt"
)

func InitProvider(store *db.Store, key *rsa.PrivateKey) fosite.OAuth2Provider {
	keyGetter := func(context.Context) (interface{}, error) {
		return key, nil
	}

	secret := []byte(os.Getenv("SECRET"))

	conf := &fosite.Config{
		GlobalSecret: secret,

		AccessTokenLifespan:                 1 * time.Hour,
		RefreshTokenLifespan:                24 * time.Hour,
		AuthorizeCodeLifespan:               5 * time.Minute,
		VerifiableCredentialsNonceLifespan:  1 * time.Hour,
		IDTokenLifespan:                     1 * time.Hour,
		EnforcePKCE:                         true,
		GrantTypeJWTBearerCanSkipClientAuth: false,
		EnablePKCEPlainChallengeMethod:      true,
		IDTokenIssuer:                       "esasy-class",
		PushedAuthorizeContextLifespan:      5 * time.Minute,
		SendDebugMessagesToClients:          true,
	}

	return compose.Compose(
		conf,
		store,
		&compose.CommonStrategy{
			CoreStrategy:               compose.NewOAuth2JWTStrategy(keyGetter, compose.NewOAuth2HMACStrategy(conf), conf),
			OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(keyGetter, conf),
			Signer:                     &jwt.DefaultSigner{GetPrivateKey: keyGetter},
		},
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2AuthorizeImplicitFactory,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		compose.OAuth2TokenIntrospectionFactory,
		compose.RFC7523AssertionGrantFactory,

		compose.OpenIDConnectExplicitFactory,
		compose.OpenIDConnectImplicitFactory,
		compose.OpenIDConnectHybridFactory,
		compose.OpenIDConnectRefreshFactory,

		compose.OAuth2TokenIntrospectionFactory,
		compose.OAuth2TokenRevocationFactory,

		compose.OAuth2PKCEFactory,
		compose.PushedAuthorizeHandlerFactory,
		compose.OIDCUserinfoVerifiableCredentialFactory,
	)
}
