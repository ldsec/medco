package utilserver

import (
	"errors"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

// oidcProvider is the definition of an OIDC identity provider
type oidcProvider struct {

	// JwksURL is the URL from which the JWT signing keys are retrieved
	JwksURL string

	// JwtIssuer is the token issuer (for JWT validation)
	JwtIssuer string

	// ClientID is the OIDC client ID
	ClientID string

	// JwtUserIDClaim is the JWT claim containing the user ID
	JwtUserIDClaim string

	// JwksTTL is the TTL of JWKS requests
	JwksTTL time.Duration

	// JwtAcceptableSkew is the acceptable shift in time for checking the JWT expiration
	JwtAcceptableSkew time.Duration

	// cachedJWKSet is the cached set of keys used to establish the trust with the identity provider,
	// valid until cachedJWKSetExpiration
	cachedJWKSet *jwk.Set

	// cachedJWKSetExpiration is the expiration time of cachedJWKSet
	cachedJWKSetExpiration time.Time
}

// retrieveJWKSets retrieves the JWK set (live or from cache if TTL not expired) and cache it
func (oidcProvider *oidcProvider) retrieveJWKSet() (keySet *jwk.Set, err error) {

	if oidcProvider.cachedJWKSet == nil || oidcProvider.cachedJWKSetExpiration.Before(time.Now()) {

		// fetch jwks with custom client to enforce timeout
		oidcProvider.cachedJWKSet, err = jwk.Fetch(
			oidcProvider.JwksURL,
			jwk.WithHTTPClient(&http.Client{
				Timeout: JwksTimeout,
			}),
		)

		if err != nil {
			logrus.Error("Error retrieving JWKS from ", oidcProvider.JwksURL, ": ", err)
			return
		}

		oidcProvider.cachedJWKSetExpiration = time.Now().Add(oidcProvider.JwksTTL)
	}

	return oidcProvider.cachedJWKSet, nil
}

// verifyTokenWithJWKSets verifies a token against several OIDC providers. Verification passes if err is not nil.
func verifyTokenWithJWKSets(token string) (tokenPayload []byte, matchingProvider *oidcProvider, err error) {

	wg := sync.WaitGroup{}
	wg.Add(len(OidcProviders))

	// check concurrently the different OIDC providers
	for _, provider := range OidcProviders {
		go func(provider *oidcProvider) {
			defer wg.Done()

			// get JWT signing keys
			keySet, err := provider.retrieveJWKSet()
			if err != nil {
				logrus.Warn("Failed to retrieve key set for provider ", provider.JwksURL, ": ", err)
				return
			}

			// signature verification attempt
			if attemptedTokenPayload, err := jws.VerifyWithJWKSet([]byte(token), keySet, nil); err == nil {
				logrus.Info("Token validation successful with provider: ", provider.JwksURL)
				if tokenPayload != nil || matchingProvider != nil {
					logrus.Warn("More than one OIDC provider matches")
				}
				tokenPayload = attemptedTokenPayload
				matchingProvider = provider
			} else {
				logrus.Debug("Token validation with provider ", provider.JwksURL, " failed: ", err)
			}

		}(provider)
	}

	wg.Wait()
	if tokenPayload == nil || matchingProvider == nil {
		err = errors.New("authentication failed (signature validation) with all providers")
		logrus.Warn(err)
	}

	return
}
