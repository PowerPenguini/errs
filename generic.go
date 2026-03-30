package errs

var GenericInternal = NewError("internal_server_error", "Internal server error", InternalType, nil)
var GenericInvalidJson = NewError("invalid_payload", "Invalid JSON payload", BadRequestType, nil)
var GenericInvalidQuery = NewError("invalid_query", "At least one query argument is invalid", BadRequestType, nil)
var GenericUnauthorized = NewError("unauthorized", "Authentication required", UnauthorizedType, nil)
var GenericForbidden = NewError("forbidden", "Insufficient permissions", ForbiddenType, nil)
