package errcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestCodeString(t *testing.T) {
	assert.Equal(t, "OK", OK.String())
	assert.Equal(t, "InvalidParam", InvalidParam.String())
	assert.Equal(t, "Unauthenticated", Unauthenticated.String())
	assert.Equal(t, "NotFound", NotFound.String())
	assert.Equal(t, "Internal", Internal.String())
}

func TestCodeToGRPC(t *testing.T) {
	tests := []struct {
		name     string
		code     Code
		expected codes.Code
	}{
		{"OK maps to OK", OK, codes.OK},
		{"InvalidParam maps to InvalidArgument", InvalidParam, codes.InvalidArgument},
		{"Unauthenticated maps to Unauthenticated", Unauthenticated, codes.Unauthenticated},
		{"NotFound maps to NotFound", NotFound, codes.NotFound},
		{"Internal maps to Internal", Internal, codes.Internal},
		{"AlreadyExists maps to AlreadyExists", AlreadyExists, codes.AlreadyExists},
		{"PermissionDenied maps to PermissionDenied", PermissionDenied, codes.PermissionDenied},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.code.GRPC())
		})
	}
}

func TestGRPCToCode(t *testing.T) {
	tests := []struct {
		name     string
		grpcCode codes.Code
		expected Code
	}{
		{"OK -> OK", codes.OK, OK},
		{"InvalidArgument -> InvalidParam", codes.InvalidArgument, InvalidParam},
		{"Unauthenticated -> Unauthenticated", codes.Unauthenticated, Unauthenticated},
		{"NotFound -> NotFound", codes.NotFound, NotFound},
		{"Internal -> Internal", codes.Internal, Internal},
		{"AlreadyExists -> AlreadyExists", codes.AlreadyExists, AlreadyExists},
		{"PermissionDenied -> PermissionDenied", codes.PermissionDenied, PermissionDenied},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, FromGRPC(tt.grpcCode))
		})
	}
}

func TestCodeMessage(t *testing.T) {
	assert.Equal(t, "OK", OK.Message())
	assert.NotEmpty(t, InvalidParam.Message())
	assert.NotEmpty(t, Unauthenticated.Message())
}
