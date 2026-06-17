// Package errcode provides unified error codes with gRPC status bidirectional mapping.
package errcode

import "google.golang.org/grpc/codes"

// Code represents a unified error code used across all services.
type Code int

const (
	OK               Code = 0
	InvalidParam     Code = 400
	Unauthenticated  Code = 401
	PermissionDenied Code = 403
	NotFound         Code = 404
	AlreadyExists    Code = 409
	Internal         Code = 500
)

// String returns the human-readable name of the error code.
func (c Code) String() string {
	switch c {
	case OK:
		return "OK"
	case InvalidParam:
		return "InvalidParam"
	case Unauthenticated:
		return "Unauthenticated"
	case PermissionDenied:
		return "PermissionDenied"
	case NotFound:
		return "NotFound"
	case AlreadyExists:
		return "AlreadyExists"
	case Internal:
		return "Internal"
	default:
		return "Unknown"
	}
}

// Message returns the default Chinese message for the error code.
func (c Code) Message() string {
	switch c {
	case OK:
		return "OK"
	case InvalidParam:
		return "参数错误"
	case Unauthenticated:
		return "未登录或登录已过期"
	case PermissionDenied:
		return "没有权限"
	case NotFound:
		return "资源不存在"
	case AlreadyExists:
		return "资源已存在"
	case Internal:
		return "服务器内部错误"
	default:
		return "未知错误"
	}
}

// GRPC maps the unified error code to a gRPC status code.
func (c Code) GRPC() codes.Code {
	switch c {
	case OK:
		return codes.OK
	case InvalidParam:
		return codes.InvalidArgument
	case Unauthenticated:
		return codes.Unauthenticated
	case PermissionDenied:
		return codes.PermissionDenied
	case NotFound:
		return codes.NotFound
	case AlreadyExists:
		return codes.AlreadyExists
	case Internal:
		return codes.Internal
	default:
		return codes.Internal
	}
}

// FromGRPC maps a gRPC status code back to a unified error code.
func FromGRPC(c codes.Code) Code {
	switch c {
	case codes.OK:
		return OK
	case codes.InvalidArgument:
		return InvalidParam
	case codes.Unauthenticated:
		return Unauthenticated
	case codes.PermissionDenied:
		return PermissionDenied
	case codes.NotFound:
		return NotFound
	case codes.AlreadyExists:
		return AlreadyExists
	case codes.Internal:
		return Internal
	default:
		return Internal
	}
}
