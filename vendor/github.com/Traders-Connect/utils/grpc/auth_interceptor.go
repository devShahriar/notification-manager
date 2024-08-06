package grpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// UserMetadataClaims contains traders connect specific metadata fields
type UserMetadataClaims struct {
	//IsSignUpFinished string `json:"isSignUpFinished"`
	OnBoardingComplete string `json:"onboardingComplete"`
}

// AppMetadataClaims contains internal info of users
type AppMetadataClaims struct {
	InternalUserID string     `json:"internal_user_id"`
	AccountLimit   UsageLimit `json:"account_limit"`
}

// UsageLimit is the trades connect usage limit of a user
type UsageLimit struct {
	Monthly int64 `json:"monthly"`
	Yearly  int64 `json:"yearly"`
}

// CustomClaims contains a custom version of the claims we handle
type CustomClaims struct {
	Scope        string             `json:"scope"`
	Permissions  []string           `json:"permissions"`
	Audience     []string           `json:"aud,omitempty"`
	Roles        []string           `json:"https://tradersconnect.com/roles"`
	UserMetadata UserMetadataClaims `json:"https://tradersconnect.com/user_metadata"`
	AppMetadata  AppMetadataClaims  `json:"https://tradersconnect.com/app_metadata"`
	UserID       string             `json:"https://tradersconnect.com/user_id"`
	jwt.StandardClaims
}

// JSONWebKeys represents web keys for the token
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// jwks is the JWKS structure
type jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type RPCRules struct {
	Rules map[string]RPCRule
}

type RPCRule struct {
	AllowedPermissions []string
	NoAuthRequired     bool
	ValidateTCUser     bool
}

// GRPCAuthInterceptor is used to intercept requests and authorize against Auth0 with JWT tokens
type GRPCAuthInterceptor struct {
	domain      string
	rpcRules    RPCRules
	jwks        jwks
	authEnabled bool
	log         *zap.SugaredLogger
}

// ContextKey is a data type used to pass the context around
type ContextKey string

var (
	ClaimsContextKey = ContextKey("claims")
)

// RequestCompany is used to parse a generic request to extract the company ID. This will be used to validate the user owns the company
type RequestCompany struct {
	CompanyId string
}

// RequestStore is used to parse a generic request to extract the store ID. This will be used to validate the user owns the company
type RequestStore struct {
	StoreID string
}

// NewGRPCAuthInterceptor returns a new GRPCAuthInterceptor
func NewGRPCAuthInterceptor(domain string, rpcRoles RPCRules, authEnabled bool, logger *zap.SugaredLogger) (*GRPCAuthInterceptor, error) {
	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	jwks, err := fetchJWKS(domain)
	if err != nil {
		return nil, err
	}

	return &GRPCAuthInterceptor{
		rpcRules:    rpcRoles,
		domain:      domain,
		jwks:        *jwks,
		authEnabled: authEnabled,
		log:         logger,
	}, nil
}

// GetCustomClaims tries to extract custom claims added to a context, if present
func GetCustomClaims(ctx context.Context) (*CustomClaims, error) {
	token := ctx.Value(ClaimsContextKey)
	claims, ok := token.(*CustomClaims)
	if !ok {
		return nil, errors.New("unable to read token from request context")
	}
	return claims, nil
}

// GetUserID returns - if present - the user ID from the token
func (cc *CustomClaims) GetUserID() (string, error) {
	if cc.UserID == "" {
		return "", errors.New("empty userID")
	}
	return cc.UserID, nil
}

// GetUserID returns the UserID in the token. These calls are unsafe and should only be called for non important things, like logs
func GetUserID(ctx context.Context) string {
	cc, err := GetCustomClaims(ctx)
	if err != nil {
		return ""
	}
	uid, _ := cc.GetUserID()
	return uid
}

// HasPermission returns true if the token has the permission
func (cc *CustomClaims) HasPermission(perm string) bool {
	for _, p := range cc.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

// MonthlyAccountLimit returns the monthly account limit of a user
func MonthlyAccountLimit(ctx context.Context) (int64, error) {
	cc, err := GetCustomClaims(ctx)
	if err != nil {
		return 0, err
	}
	return cc.AppMetadata.AccountLimit.Monthly, nil
}

// YearlyAccountLimit returns the yearly account limit of a user
func YearlyAccountLimit(ctx context.Context) (int64, error) {
	cc, err := GetCustomClaims(ctx)
	if err != nil {
		return 0, err
	}
	return cc.AppMetadata.AccountLimit.Yearly, nil
}

// GetInternalUserID returns the tc user id associated with user
func (cc *CustomClaims) GetInternalUserID() (string, error) {
	if cc.AppMetadata.InternalUserID == "" {
		return "", errors.New("no tc user associated with this auth0 user")
	}
	return cc.AppMetadata.InternalUserID, nil
}

// GetTCUserID returns the tc user in the token. These calls are unsafe and should only be called for non important things, like logs
func GetTCUserID(ctx context.Context) string {
	cc, err := GetCustomClaims(ctx)
	if err != nil {
		return ""
	}
	user, _ := cc.GetInternalUserID()
	return user
}

// IsRequestUserIDValid checks whether the company passed in the request matches the one in the token.
// These calls are unsafe and should only be called for non important things, like logs
func IsRequestUserIDValid(ctx context.Context, userID string) bool {
	cc, err := GetCustomClaims(ctx)
	if err != nil {
		return false
	}

	claimsUser, err := cc.GetInternalUserID()
	if err != nil {
		return false
	}

	if claimsUser != userID {
		return false
	}
	return true
}

// Interceptor is the auth interceptor
func (i *GRPCAuthInterceptor) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// we always want to deny access if there are no permissions defined, so we always check
	rules, ok := i.rpcRules.Rules[info.FullMethod]
	if !ok {
		i.log.Errorw("no access rule defined for endpoint", "rpc", info.FullMethod)
		return nil, status.Error(codes.PermissionDenied, "invalid endpoint")
	}

	// by default we will authorize every rpc endpoint request. but we will skip authorization for those which are
	// explicitly marked NoAuthRequired = true
	if rules.NoAuthRequired {
		return handler(ctx, req)
	}

	claims, err := i.authorize(ctx, rules.AllowedPermissions)
	if err != nil {
		i.log.Error(err)
		return nil, err
	}

	if claims != nil {
		ctx = context.WithValue(ctx, ClaimsContextKey, claims)
	}

	if rules.ValidateTCUser {
		// check the tc user id passed in the request is owned by the user
		cPtr := reflect.ValueOf(req)
		userID := reflect.Indirect(cPtr).FieldByName("UserId").String()
		if !IsRequestUserIDValid(ctx, userID) {
			i.log.Errorw("the requested user id is not valid for this user", "rpc", info.FullMethod, "user_id", GetUserID(ctx), "internal_user_id", userID)
			return nil, status.Error(codes.PermissionDenied, "unauthorized")
		}
	}

	return handler(ctx, req)
}

// fetchJWKS retrieves the public JWKS for the tokens
func fetchJWKS(domain string) (*jwks, error) {
	resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", domain))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks = jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	return &jwks, err
}

// getPemCert converts a token in to PEM format if it belongs to the public keys in JWKS
func (i *GRPCAuthInterceptor) getPemCert(token *jwt.Token) (string, error) {
	var cert string
	for k := range i.jwks.Keys {
		if token.Header["kid"] == i.jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + i.jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		return cert, errors.New("unable to find appropriate key")
	}

	return cert, nil
}

// getAccessToken returns the access token from a GRPC context
func (i *GRPCAuthInterceptor) getAccessToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	return values[0], nil
}

// authorize autorizes an rpc endpoint based on the rpcRoles
func (i *GRPCAuthInterceptor) authorize(ctx context.Context, requiredPermissions []string) (*CustomClaims, error) {
	bearerToken, err := i.getAccessToken(ctx)
	if err != nil {
		i.log.Errorf("could not read bearer token: %s", err)
		return nil, err
	}

	accessToken := strings.TrimPrefix(bearerToken, "Bearer ")
	claims, err := i.getClaims(accessToken)
	if err != nil {
		i.log.Errorf("invalid access token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid access token: %v", err)
	}

	for _, permission := range claims.Permissions {
		for _, role := range requiredPermissions {
			if role == permission {
				return claims, nil
			}
		}
	}

	return nil, status.Error(codes.PermissionDenied, "no permission to access this RPC")
}

// WithRoles overwrites the roles for a CustomClaim
func (cc *CustomClaims) WithRoles(roles []string) *CustomClaims {
	cc.Roles = roles
	return cc
}

// WithPermissions overwrites the permissions for a CustomClaim
func (cc *CustomClaims) WithPermissions(perms []string) *CustomClaims {
	cc.Permissions = perms
	return cc
}

// getClaims returns a CustomClaims structure extracted froom the pem version of the token
func (i *GRPCAuthInterceptor) getClaims(accessToken string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := i.getPemCert(token)
		if err != nil {
			return nil, err
		}
		return jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	})

	// if we're in test mode we don't care about the token being valid
	if i.authEnabled {
		if err != nil || !token.Valid {
			return nil, fmt.Errorf("invalid token: %w", err)
		}
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
