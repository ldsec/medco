package util

import (
	"encoding/json"
	"errors"
	"github.com/ldsec/medco-connector/restapi/models"
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
	authorizedQueryTypes, err := extractAuthorizationsFromToken(&token)
	if err != nil {
		return
	}
	logrus.Info("User ", user.ID, " has authorizations for ", len(authorizedQueryTypes), " query types")

	user.Authorizations = &models.UserAuthorizations{
		QueryType: authorizedQueryTypes,
	}
	return
}

// extractAuthorizationsFromToken parsed the token to extract the user's authorizations
func extractAuthorizationsFromToken(token *jwt.Token) (queryTypes []models.QueryType, err error) {

	// retrieve roles, within the keycloak pre-determined structure (this is ugly)
	var extractedRoles []string
	if tokenResourceAccess, ok := token.Get("resource_access"); ok {
		logrus.Trace("1 OK")
		if tokenResourceAccessTyped, ok := tokenResourceAccess.(map[string]interface{}); ok {
			logrus.Trace("2 OK")
			if clientId, ok := tokenResourceAccessTyped[OidcClientID]; ok {
				logrus.Trace("3 OK")
				if clientIdTyped, ok := clientId.(map[string]interface{}); ok {
					logrus.Trace("4 OK")
					if roles, ok := clientIdTyped["roles"]; ok {
						logrus.Trace("5 OK")
						if extractedRolesUntyped, ok := roles.([]interface{}); ok {
							logrus.Trace("6 OK")
							for _, extractedRoleUntyped := range extractedRolesUntyped {
								if extractedRole, ok := extractedRoleUntyped.(string); ok {
									extractedRoles = append(extractedRoles, extractedRole)
								} else {
									logrus.Warn("could not parse authorization", extractedRole)
								}
							}
						}
					}
				}
			}
		}
	}

	if len(extractedRoles) == 0 {
		err = errors.New("error retrieving roles from token, or user has no authorizations")
		logrus.Error(err)
		return
	}

	// match roles to query types
	for _, extractedRole := range extractedRoles {
		switch string(extractedRole) {
		case string(models.QueryTypePatientList):
			queryTypes = append(queryTypes, models.QueryTypePatientList)
		case string(models.QueryTypeCountPerSite):
			queryTypes = append(queryTypes, models.QueryTypeCountPerSite)
		case string(models.QueryTypeCountPerSiteObfuscated):
			queryTypes = append(queryTypes, models.QueryTypeCountPerSiteObfuscated)
		case string(models.QueryTypeCountPerSiteShuffled):
			queryTypes = append(queryTypes, models.QueryTypeCountPerSiteShuffled)
		case string(models.QueryTypeCountPerSiteShuffledObfuscated):
			queryTypes = append(queryTypes, models.QueryTypeCountPerSiteShuffledObfuscated)
		case string(models.QueryTypeCountGlobal):
			queryTypes = append(queryTypes, models.QueryTypeCountGlobal)
		case string(models.QueryTypeCountGlobalObfuscated):
			queryTypes = append(queryTypes, models.QueryTypeCountGlobalObfuscated)

		default:
			logrus.Debug("ignored role ", extractedRole)
		}
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