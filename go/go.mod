module github.com/artofthepossible/ghas-me-up/go

// Intentionally stale: Go 1.16 is long out of support.
go 1.16

// Every module is pinned to an old release with published advisories.
require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/text v0.3.0
	gopkg.in/yaml.v2 v2.2.2
)
