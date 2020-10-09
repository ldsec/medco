package utilserver

import (
	"encoding/json"
	"errors"
	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/sirupsen/logrus"
)

// AuthenticateUser authenticates user and creates principal with user information, including its authorizations
// returns error if user is not authorized
func AuthenticateUser(token string) (user *models.User, err error) {

	// verify signature
	tokenPayload, matchingProvider, err := verifyTokenWithJWKSets(token)
	if err != nil {
		logrus.Warn("authentication failed (signature validation): ", err)
		return
	}

	// parse and validate claims
	var parsedToken jwt.Token
	if err = json.Unmarshal(tokenPayload, &parsedToken); err != nil {
		logrus.Warn("authentication failed (token parsing error): ", err)
		return
	}

	err = parsedToken.Verify(
		jwt.WithIssuer(matchingProvider.JwtIssuer),
		jwt.WithAudience(matchingProvider.ClientID),
		jwt.WithAcceptableSkew(matchingProvider.JwtAcceptableSkew),
	)
	if err != nil {
		logrus.Warn("authentication failed (invalid claim): ", err)
		return
	}

	// extract user name
	user = &models.User{}
	if userID, ok := parsedToken.Get(matchingProvider.JwtUserIDClaim); ok {
		user.ID = userID.(string)
		user.Token = token
		logrus.Info("authenticated user " + user.ID)
	} else {
		err = errors.New("authentication failed (user ID claim not present)")
		logrus.Warn(err)
	}

	// extract user authorizations
	user.Authorizations, err = extractAuthorizationsFromToken(&parsedToken, matchingProvider)
	return
}

// extractAuthorizationsFromToken parsed the token to extract the user's authorizations
func extractAuthorizationsFromToken(token *jwt.Token, provider *oidcProvider) (ua *models.UserAuthorizations, err error) {

	// retrieve roles, within the keycloak pre-determined structure (this is ugly)
	var extractedRoles []string
	if tokenResourceAccess, ok := token.Get("resource_access"); ok {
		logrus.Trace("1 OK")
		if tokenResourceAccessTyped, ok := tokenResourceAccess.(map[string]interface{}); ok {
			logrus.Trace("2 OK")
			if clientID, ok := tokenResourceAccessTyped[provider.ClientID]; ok {
				logrus.Trace("3 OK")
				if clientIDTyped, ok := clientID.(map[string]interface{}); ok {
					logrus.Trace("4 OK")
					if roles, ok := clientIDTyped["roles"]; ok {
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

	// match roles to authorizations
	ua = &models.UserAuthorizations{}
	authCount := len(extractedRoles)
	for _, extractedRole := range extractedRoles {
		switch extractedRole {

		// rest api authorizations
		case string(models.RestAPIAuthorizationMedcoNetwork):
			ua.RestAPI = append(ua.RestAPI, models.RestAPIAuthorizationMedcoNetwork)
		case string(models.RestAPIAuthorizationMedcoExplore):
			ua.RestAPI = append(ua.RestAPI, models.RestAPIAuthorizationMedcoExplore)
		case string(models.RestAPIAuthorizationMedcoGenomicAnnotations):
			ua.RestAPI = append(ua.RestAPI, models.RestAPIAuthorizationMedcoGenomicAnnotations)
		case string(models.RestAPIAuthorizationMedcoSurvivalAnalysis):
			ua.RestAPI = append(ua.RestAPI, models.RestAPIAuthorizationMedcoSurvivalAnalysis)

		// explore query type authorizations
		case string(models.ExploreQueryTypePatientList):
			ua.ExploreQuery = append(ua.ExploreQuery, models.ExploreQueryTypePatientList)
		case string(models.ExploreQueryTypeCountPerSite):
			ua.ExploreQuery = append(ua.ExploreQuery, models.ExploreQueryTypeCountPerSite)
		case string(models.ExploreQueryTypeCountPerSiteObfuscated):
			ua.ExploreQuery = append(ua.ExploreQuery, models.ExploreQueryTypeCountPerSiteObfuscated)
		case string(models.ExploreQueryTypeCountPerSiteShuffled):
			ua.ExploreQuery = append(ua.ExploreQuery, models.ExploreQueryTypeCountPerSiteShuffled)
		case string(models.ExploreQueryTypeCountPerSiteShuffledObfuscated):
			ua.ExploreQuery = append(ua.ExploreQuery, models.ExploreQueryTypeCountPerSiteShuffledObfuscated)
		case string(models.ExploreQueryTypeCountGlobal):
			ua.ExploreQuery = append(ua.ExploreQuery, models.ExploreQueryTypeCountGlobal)
		case string(models.ExploreQueryTypeCountGlobalObfuscated):
			ua.ExploreQuery = append(ua.ExploreQuery, models.ExploreQueryTypeCountGlobalObfuscated)

		default:
			logrus.Debug("ignored role ", extractedRole)
			authCount--
		}
	}

	logrus.Debug("User has ", authCount, " authorizations")
	return
}

// AuthorizeRestAPIEndpoint authorizes the REST API endpoint requested by the user
func AuthorizeRestAPIEndpoint(user *models.User, requiredAuthorization models.RestAPIAuthorization) (err error) {
	for _, userAuth := range user.Authorizations.RestAPI {
		if userAuth == requiredAuthorization {
			logrus.Info("user is authorized to request the endpoint with authorization " + string(requiredAuthorization))
			return nil
		}
	}

	err = errors.New("user is not authorized to request the endpoint with authorization " + string(requiredAuthorization))
	logrus.Warn(err)
	return
}

// AuthorizeExploreQueryType authorizes the explore query type requested by the user
func AuthorizeExploreQueryType(user *models.User, requestedQueryType models.ExploreQueryType) (err error) {
	for _, userQueryType := range user.Authorizations.ExploreQuery {
		if userQueryType == requestedQueryType {
			logrus.Info("user is authorized to execute the query type " + string(requestedQueryType))
			return nil
		}
	}

	err = errors.New("user is not authorized to execute the query type " + string(requestedQueryType))
	logrus.Warn(err)
	return
}
