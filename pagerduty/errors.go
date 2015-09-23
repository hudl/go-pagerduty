package pagerduty

const (
	ErrInternalError uint = 2000 + iota
	ErrInvalidInputProvided
	ErrArgumentsCausedError
	ErrMissingArguments
	ErrInvalidSinceOrUntilParameterValues
	ErrInvalidQueryDateRange
	ErrAuthenticationFailed
	ErrAccountNotFound
	ErrAccountLocked
	ErrOnlyHTTPSAllowedForThisCall
	ErrAccessDenied
	ErrRequiresRequesterID
	ErrAccountExpired
)
