package util

import (
	"encoding/json"
	"errors"
	"github.com/lca1/medco-connector/models"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/sirupsen/logrus"
	"time"
)

var cachedKeySet struct {
	keySet *jwk.Set
	expirationTime time.Time
}

// AuthorizeUser authorizes user and populate principal with user information, including its authorizations
// returns error if user is not authorized
func AuthorizeUser(credentials models.ResourceCredentials, user *models.User) (err error) {

	// get JWT signing keys
	keySet, err := retrieveJWKSet()
	if err != nil {
		logrus.Warn("failed to retrieve key set: ", err)
		return
	}

	// signature verification
	tokenPayload, err := jws.VerifyWithJWKSet([]byte(credentials.MEDCOTOKEN), keySet, nil)
	if err != nil {
		logrus.Warn("authentication failed (signature validation): ", err)
		return
	}

	// parse and validate claims
	var token jwt.Token
	if err = json.Unmarshal(tokenPayload, &token); err != nil {
		logrus.Warn("authentication failed (token parsing error): ", err)
		return
	}

	err = token.Verify(
		jwt.WithIssuer(OidcJwtIssuer),
		jwt.WithAudience(OidcClientID),
		jwt.WithAcceptableSkew(30 * time.Second),
	)
	if err != nil {
		logrus.Warn("authentication failed (invalid claim): ", err)
		return
	}

	// extract user name
	if userID, ok := token.Get(OidcJwtUserIDClaim); ok {
		user.ID = userID.(string)
		user.Token = credentials.MEDCOTOKEN
		logrus.Info("authenticated user " + user.ID)
	} else {
		err = errors.New("authentication failed (user ID claim not present)")
		logrus.Warn(err)
	}

	// extract user authorizations
	// todo: implement the extraction
	user.Authorizations = &models.UserAuthorizations{
		QueryType: []models.QueryType{
			models.QueryTypePatientList,
			models.QueryTypeCountPerSite,
			models.QueryTypeCountPerSiteObfuscated,
			models.QueryTypeCountPerSiteShuffled,
			models.QueryTypeCountPerSiteShuffledObfuscated,
			models.QueryTypeCountGlobal,
			models.QueryTypeCountGlobalObfuscated,
		},
	}
	return
}

// AuthorizeQueryType authorizes the query type requested by the user
func AuthorizeQueryType(user models.User, requestedQueryType models.QueryType) (err error) {
	for _, userQueryType := range user.Authorizations.QueryType {
		if userQueryType == requestedQueryType {
			logrus.Info("user is authorized to execute the query type " + string(requestedQueryType))
			return nil
		}
	}

	err = errors.New("user is not authorized to execute the query type " + string(requestedQueryType))
	logrus.Warn(err)
	return
}

// retrieveJWKSet retrieves the JWK set (live or from cache if TTL not expired) and cache it
func retrieveJWKSet() (keySet *jwk.Set, err error) {

	if cachedKeySet.keySet == nil || cachedKeySet.expirationTime.Before(time.Now()) {
		cachedKeySet.keySet, err = jwk.Fetch(JwksURL)
		if err != nil {
			return
		}

		cachedKeySet.expirationTime = time.Now().Add(time.Duration(JwksTTLSeconds) * time.Second)
	}
	return cachedKeySet.keySet, nil
}