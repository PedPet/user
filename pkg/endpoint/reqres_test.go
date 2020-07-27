package endpoint

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"github.com/PedPet/user/config"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
)

var rules config.Validation

func init() {
	validation, err := ioutil.ReadFile("../../config/validation.json")
	if err != nil {
		log.Panicf("Failed to open validation.json: %s", err)
	}
	err = json.Unmarshal(validation, &rules)
	if err != nil {
		log.Panicf("Failed to decode json: %s", err)
	}
}

func TestCreateUserRequestValidate(t *testing.T) {
	testCases := []struct {
		name     string
		payload  CreateUserRequest
		expected string
	}{
		{
			name: "Valid",
			payload: CreateUserRequest{
				Username:    faker.Username(),
				Password:    faker.Password() + "1!",
				Email:       "scrott@gmail.com",
				PhoneNumber: faker.Phonenumber(),
			},
			expected: "",
		},
		{
			name: "Missing username",
			payload: CreateUserRequest{
				Password:    faker.Password() + "!1",
				Email:       "scrott@gmail.com",
				PhoneNumber: faker.Phonenumber(),
			},
			expected: "username: cannot be blank.",
		},
		{
			name: "Missing password",
			payload: CreateUserRequest{
				Username: faker.Username(),
				Email:    "scrott@gmail.com",
			},
			expected: "password: cannot be blank.",
		},
		{
			name: "Missing email",
			payload: CreateUserRequest{
				Username: faker.Username(),
				Password: faker.Password() + "!1",
			},
			expected: "email: cannot be blank.",
		},
		{
			name: "Incorrect username",
			payload: CreateUserRequest{
				Username: "t",
				Password: faker.Password() + "2@",
				Email:    "scrott@gmail.com",
			},
			expected: "username: the length must be between 2 and 100.",
		},
		{
			name: "Incorrect password format",
			payload: CreateUserRequest{
				Username: faker.Username(),
				Password: faker.Password(),
				Email:    "scrott@gmail.com",
			},
			expected: "password: must be in a valid format.",
		},
		{
			name: "Incorrect email format",
			payload: CreateUserRequest{
				Username: faker.Username(),
				Password: faker.Password() + "!4",
				Email:    faker.Email(),
			},
			expected: "email: must be a valid email address.",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.payload.Validate(rules.Password)
			if err != nil {
				assert.EqualError(t, err, tc.expected)
				return
			}

			assert.True(t, err == nil && tc.expected == "", tc.expected)
		})
	}
}

func TestConfirmUserRequestValidate(t *testing.T) {
	testCases := []struct {
		name     string
		payload  ConfirmUserRequest
		expected string
	}{
		{
			name: "Valid",
			payload: ConfirmUserRequest{
				Username: faker.Username(),
				Code:     "123456",
			},
			expected: "",
		},
		{
			name: "Missing username",
			payload: ConfirmUserRequest{
				Code: "123456",
			},
			expected: "username: cannot be blank.",
		},
		{
			name: "Missing code",
			payload: ConfirmUserRequest{
				Username: faker.Username(),
			},
			expected: "code: cannot be blank.",
		},
		{
			name: "Incorrect username",
			payload: ConfirmUserRequest{
				Username: "t",
				Code:     "123456",
			},
			expected: "username: the length must be between 2 and 100.",
		},
		{
			name: "Incorrect code length",
			payload: ConfirmUserRequest{
				Username: faker.Username(),
				Code:     "11",
			},
			expected: "code: the length must be exactly 6.",
		},
		{
			name: "Incorrect code string",
			payload: ConfirmUserRequest{
				Username: faker.Username(),
				Code:     "ss",
			},
			expected: "code: must contain digits only.",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.payload.Validate()
			if err != nil {
				assert.EqualError(t, err, tc.expected)
				return
			}

			assert.True(t, err == nil && tc.expected == "", tc.expected)
		})
	}
}

func TestResendConfirmationRequest(t *testing.T) {
	testCases := []struct {
		name     string
		payload  ResendConfirmationRequest
		expected string
	}{
		{
			name: "Valid",
			payload: ResendConfirmationRequest{
				Username: faker.Username(),
			},
			expected: "",
		},
		{
			name:     "Missing username",
			payload:  ResendConfirmationRequest{},
			expected: "username: cannot be blank.",
		},
		{
			name: "Invalid username",
			payload: ResendConfirmationRequest{
				Username: "t",
			},
			expected: "username: the length must be between 2 and 100.",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.payload.Validate()
			if err != nil {
				assert.EqualError(t, err, tc.expected)
				return
			}

			assert.True(t, err == nil && tc.expected == "", tc.expected)
		})
	}
}

func TestUsernameTakenRequestValidate(t *testing.T) {
	testCases := []struct {
		name     string
		payload  UsernameTakenRequest
		expected string
	}{
		{
			name: "Valid",
			payload: UsernameTakenRequest{
				Username: faker.Username(),
			},
			expected: "",
		},
		{
			name:     "Missing username",
			payload:  UsernameTakenRequest{},
			expected: "username: cannot be blank.",
		},
		{
			name: "Invalid username",
			payload: UsernameTakenRequest{
				Username: "t",
			},
			expected: "username: the length must be between 2 and 100.",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.payload.Validate()
			if err != nil {
				assert.EqualError(t, err, tc.expected)
				return
			}

			assert.True(t, err == nil && tc.expected == "", tc.expected)
		})

	}
}

func TestLoginRequestValidate(t *testing.T) {
	testCases := []struct {
		name     string
		payload  LoginRequest
		expected string
	}{
		{
			name: "Valid",
			payload: LoginRequest{
				Username: faker.Username(),
				Password: faker.Password() + "1!",
			},
			expected: "",
		},
		{
			name: "Missing username",
			payload: LoginRequest{
				Password: faker.Password() + "1!",
			},
			expected: "username: cannot be blank.",
		},
		{
			name: "Missing password",
			payload: LoginRequest{
				Username: faker.Username(),
			},
			expected: "password: cannot be blank.",
		},
		{
			name: "Invalid username",
			payload: LoginRequest{
				Username: "t",
				Password: faker.Password() + "1!",
			},
			expected: "username: the length must be between 2 and 100.",
		},
		{
			name: "Invalid password",
			payload: LoginRequest{
				Username: faker.Username(),
				Password: faker.Password(),
			},
			expected: "password: must be in a valid format.",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.payload.Validate(rules.Password)
			if err != nil {
				assert.EqualError(t, err, tc.expected)
				return
			}

			assert.True(t, err == nil && tc.expected == "", tc.expected)
		})
	}
}

func TestVerifyJWTRequestValidate(t *testing.T) {
	testCases := []struct {
		name     string
		payload  VerifyJWTRequest
		expected string
	}{
		{
			name: "Valid",
			payload: VerifyJWTRequest{
				Jwt: faker.Word(),
			},
			expected: "",
		},
		{
			name:     "Missing JWT",
			payload:  VerifyJWTRequest{},
			expected: "jwt: cannot be blank.",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.payload.Validate()
			if err != nil {
				assert.EqualError(t, err, tc.expected)
				return
			}

			assert.True(t, err == nil && tc.expected == "", tc.expected)
		})
	}
}

func TestUserDetailsRequest(t *testing.T) {
	testCases := []struct {
		name     string
		payload  UserDetailsRequest
		expected string
	}{
		{
			name: "Valid",
			payload: UserDetailsRequest{
				Jwt: faker.Word(),
			},
			expected: "",
		},
		{
			name:     "Missing JWT",
			payload:  UserDetailsRequest{},
			expected: "jwt: cannot be blank.",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.payload.Validate()
			if err != nil {
				assert.EqualError(t, err, tc.expected)
				return
			}

			assert.True(t, err == nil && tc.expected == "", tc.expected)
		})
	}
}
