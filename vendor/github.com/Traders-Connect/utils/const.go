package utils

const (
	// UTCLayout is the exact layout we need to parse timestamp string to go timestamp
	UTCLayout = "2006-01-02 15:04:05 -0700 MST"
)

// Auth0 roles const
const (
	// Auth0FreeMemberRole is the free member role id in auth0
	Auth0FreeMemberRole  = "rol_IcIRN1EdsYVOddsB"
	Auth0PaidUserRole    = "rol_zTLc1Veg2MttJ5Hi"
	Auth0ExpiredUserRole = "rol_jqJoKSbMAArDKI2e"
)

const (
	ConnectionStatusConnected              = "connected"
	ConnectionStatusReconnecting           = "reconnecting"
	ConnectionStatusConnectionError        = "connection_error"
	ConnectionStatusConnectionDisconnected = "disconnected"
)
