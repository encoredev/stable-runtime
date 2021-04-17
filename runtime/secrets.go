package runtime

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

func LoadSecret(key string) string {
	secretsOnce.Do(loadSecrets)
	val, ok := secrets[key]
	if !ok {
		fmt.Fprintln(os.Stderr, "encore: could not find secret", key)
		os.Exit(2)
	}
	return val
}

var (
	secretsOnce sync.Once
	secrets     map[string]string
)

func loadSecrets() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s, err := fetchSecrets(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal: could not load secrets:", err)
		os.Exit(1)
	}
	secrets = s
}
