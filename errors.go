package certvault

import "errors"

var (
	ErrClosed       = errors.New("certvault: vault closed")
	ErrNotFound     = errors.New("certvault: entry not found")
	ErrInvalidPEM   = errors.New("certvault: invalid pem")
	ErrNoSigner     = errors.New("certvault: signer not configured")
	ErrRevoked      = errors.New("certvault: certificate revoked")
	ErrPersist      = errors.New("certvault: persist failed")
	ErrCanceled     = errors.New("certvault: canceled")
	ErrBadRequest   = errors.New("certvault: bad request")
	ErrDuplicate    = errors.New("certvault: duplicate name")
	ErrRotateFailed = errors.New("certvault: rotate failed")
)
