package errs

import "github.com/zeromicro/x/errors"

var UnAuthorizationException = errors.New(401, "Unauthorized")

var NotFoundException = errors.New(404, "Not found")

var DatabaseError = errors.New(500, "Database error")
