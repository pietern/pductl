package snmp

type Config struct {
	Address  string
	Username string

	// Use SHA authentication with this password.
	Password string

	// Use AES encryption with this key.
	Key string
}
