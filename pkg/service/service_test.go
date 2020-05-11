package service

import (
	"context"
	"testing"

	"github.com/PedPet/user/model"
    "github.com/PedPet/user/pkg/endpoint"
    "github.com/onsi/gomega"
)

var logger log.Logger
{
    logger = log.NewLogfmtLogger(os.Stderr)
    logger = log.NewSyncLogger(logger)
    logger = log.With(logger,
        "service", "login",
        "time:", log.DefaultTimestampUTC,
        "caller", log.DefaultCaller,
    )
}

func TestCreateUser(t *testing.T) {
    testCases := []struct{
        name string
        req *endpoint.CreateUserRequest
        expectedErr bool
    }{
        {
            name: "req ok",
            req: *endpoint.CreateUserRequest{
                Username: "sc7639",
                Email: "scrott@gmail.com",
                Password: "testpass",
                PhoneNumber: "07700000000",
            },
            expectedErr: false,
        },
        {
            name: "req with empty request",
            req: *endpoint.CreateUserRequest{},
            expectedErr: true,
        },
        {
            name: "nil request",
            req: nil,
            expectedErr: true,
        }
    }

    for _, tc := range testCases {
        testCase := tc
        t.Run(testCase.name, func(t *testing.T) {
            t.Parallel()
            g := gomega.
            ctx := context.Background()

            // call
            service := NewUserService()

        });
    }
}
