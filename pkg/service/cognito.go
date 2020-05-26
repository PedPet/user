package service

import (
	"context"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"

	"github.com/PedPet/user/config"
	"github.com/PedPet/user/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
	"github.com/lestrrat/go-jwx/jwk"
	"github.com/pkg/errors"
)

// CognitoClient derscribes the cognito client methods / functionality
type CognitoClient interface {
	Register(ctx context.Context, user *model.User) error
	OTP(ctx context.Context, user *model.User, otp string) error
	ResendConfirmation(ctx context.Context, username string) error
	CheckUsernameTaken(ctx context.Context, username string) (bool, error)
	Login(ctx context.Context, username, password string) (*cognito.AuthenticationResultType, error)
	getWellKnownJWTKs() error
	ParseAndVerifyJWT(ctx context.Context, token string) (*jwt.Token, error)
	GetUserDetails(ctx context.Context, accessToken string) (*model.User, error)
}

const flowUsernamePassword = "USER_PASSWORD_AUTH"
const flowRefreshToken = "REFRESH_TOKEN_AUTH"

type cognitoClient struct {
	cognitoClient *cognito.CognitoIdentityProvider
	userPoolID    string
	appClientID   string
	clientSecret  string
	logger        log.Logger
	wellKnownJWKs *jwk.Set
	region        string
}

// NewCognitoClient creates a cognito type with the required dependencies
func NewCognitoClient(
	identity *cognito.CognitoIdentityProvider,
	cfg config.AWSSettings,
	logger log.Logger,
) (CognitoClient, error) {
	c := &cognitoClient{
		cognitoClient: identity,
		userPoolID:    cfg.CognitoUserPoolID,
		appClientID:   cfg.CognitoAppClientID,
		clientSecret:  cfg.CognitoClientSecret,
		logger:        logger,
		region:        cfg.Region,
	}

	err := c.getWellKnownJWTKs()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create cognito client")
	}

	return c, nil
}

func (c *cognitoClient) getWellKnownJWTKs() error {
	// https://cognito-idp.<region>.amazonaws.com/<pool_id>/.well-known/jwks.json
	var url strings.Builder
	url.WriteString("https://cognito-idp.")
	url.WriteString(c.region)
	url.WriteString(".amazonaws.com/")
	url.WriteString(c.userPoolID)
	url.WriteString("/.well-known/jwks.json")

	set, err := jwk.Fetch(url.String())
	if err != nil {
		return errors.Wrap(err, "Failed to get well known JSON web token key set")
	}

	c.wellKnownJWKs = set
	return nil
}

// Register is self explanatory
func (c cognitoClient) Register(ctx context.Context, user *model.User) error {
	logger := log.With(c.logger, "method", "Register")

	s := calculateSecretHash(user.Username, c.appClientID, c.clientSecret)
	cognitoUser := &cognito.SignUpInput{
		Username:   aws.String(user.Username),
		Password:   aws.String(user.Password),
		ClientId:   aws.String(c.appClientID),
		SecretHash: aws.String(s),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(user.Email),
			},
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(user.PhoneNumber),
			},
		},
	}

	output, err := c.cognitoClient.SignUp(cognitoUser)
	if err != nil {
		return errors.Wrap(err, "")
	}

	logger.Log("Register User", output, user)
	return nil
}

func calculateSecretHash(username, clientID, clientSecret string) string {
	h := hmac.New(sha256.New, []byte(clientSecret))
	h.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// OTP handles registrations confirmation via verification code
func (c cognitoClient) OTP(ctx context.Context, user *model.User, otp string) error {
	logger := log.With(c.logger, "method", "OTP")

	s := calculateSecretHash(user.Username, c.appClientID, c.clientSecret)

	cu := &cognito.ConfirmSignUpInput{
		ConfirmationCode: aws.String(otp),
		Username:         aws.String(user.Username),
		ClientId:         aws.String(c.appClientID),
		SecretHash:       aws.String(s),
	}
	output, err := c.cognitoClient.ConfirmSignUp(cu)
	if err != nil {
		return err
	}

	logger.Log("Confirm sign up", output)
	return nil
}

// ResendConfirmation handles resending OTP (confirmation codes) to the user's email address
func (c cognitoClient) ResendConfirmation(ctx context.Context, username string) error {
	logger := log.With(c.logger, "method", "ResendConfirmation")

	s := calculateSecretHash(username, c.appClientID, c.clientSecret)
	rc := &cognito.ResendConfirmationCodeInput{
		Username:   aws.String(username),
		ClientId:   aws.String(c.appClientID),
		SecretHash: aws.String(s),
	}
	output, err := c.cognitoClient.ResendConfirmationCode(rc)
	if err != nil {
		return err
	}

	logger.Log("Resend confirmation output:", output)
	return nil
}

// CheckUserNameTaken checks to see if a username has already in use, in the userpool
func (c cognitoClient) CheckUsernameTaken(ctx context.Context, username string) (bool, error) {
	logger := log.With(c.logger, "method", "CheckUsernameTaken")

	cu := &cognito.AdminGetUserInput{
		UserPoolId: aws.String(c.userPoolID),
		Username:   aws.String(username),
	}
	_, err := c.cognitoClient.AdminGetUser(cu)
	if err != nil {
		err2, ok := err.(awserr.Error)
		if ok {
			if err2.Code() == cognito.ErrCodeUserNotFoundException {
				return false, nil
			}
		}

		return false, err
	}

	logger.Log("Username taken")
	return true, nil
}

// Login uses a username and password to log a user in and return their authentication
func (c cognitoClient) Login(ctx context.Context, username, password string) (*cognito.AuthenticationResultType, error) {
	logger := log.With(c.logger, "method", "Login")

	s := calculateSecretHash(username, c.appClientID, c.clientSecret)
	flow := aws.String(flowUsernamePassword)
	params := map[string]*string{
		"USERNAME":    aws.String(username),
		"PASSWORD":    aws.String(password),
		"SECRET_HASH": aws.String(s),
	}

	auth := &cognito.InitiateAuthInput{
		AuthFlow:       flow,
		AuthParameters: params,
		ClientId:       aws.String(c.appClientID),
	}
	output, err := c.cognitoClient.InitiateAuth(auth)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to authenticate user")
	}

	if output.AuthenticationResult == nil {
		return nil, errors.New("Failed to get authenticate user")
	}

	logger.Log("Login output:", output)
	return output.AuthenticationResult, nil
}

// ParseAnVerifyJWT is self explanatory
func (c cognitoClient) ParseAndVerifyJWT(ctx context.Context, token string) (*jwt.Token, error) {
	logger := log.With(c.logger, "method", "ParseAndVerifyJWT")

	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// fmt.Printf("Token: %v\n", t)

		// Looking up the key id will return an array of just one key
		keys := c.wellKnownJWKs.LookupKeyID(token.Header["kid"].(string))
		if len(keys) == 0 {
			return nil, errors.New("Could not find matching `kid` in well known tokens")
		}

		// Build the public RSA key
		key, err := keys[0].Materialize()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create public key")
		}

		rsaPublicKey := key.(*rsa.PublicKey)
		return rsaPublicKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Invalid token")
	}

	if t.Valid != true {
		return nil, errors.New("Token is not valid")
	}

	logger.Log("Token is valid")
	return t, nil
}

// GetUserDetails gets the user's user attributes from aws
func (c cognitoClient) GetUserDetails(ctx context.Context, accessToken string) (*model.User, error) {
	logger := log.With(c.logger, "method", "GetUserDetails")

	ui := &cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	}
	output, err := c.cognitoClient.GetUser(ui)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get user")
	}

	user := &model.User{
		Username: aws.StringValue(output.Username),
	}

	// Put user attributes into user model
	for _, attr := range output.UserAttributes {
		switch aws.StringValue(attr.Name) {
		case "phone_number":
			user.PhoneNumber = aws.StringValue(attr.Value)
		case "email":
			user.Email = aws.StringValue(attr.Value)
		case "email_verified":
			confirmed, err := strconv.ParseBool(aws.StringValue(attr.Value))
			if err != nil {
				user.Confirmed = false
			}
			user.Confirmed = confirmed
		}
	}

	logger.Log("GetUserDetails output:", output)
	return user, nil
}
