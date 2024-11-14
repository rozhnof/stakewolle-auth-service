package secrets

import "os"

const secretKeyEnv = "SECRET_KEY"

type SecretManager struct{}

func (m SecretManager) SecretKey() []byte {
	return []byte(os.Getenv(secretKeyEnv))
}
