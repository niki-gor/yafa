package env

import (
	"errors"
	"os"
)

var requiredEnv = [...]string{
	"POSTGRES_USER",
	"POSTGRES_PASSWORD",
	"POSTGRES_DB",
}

func GetRequired() (map[string]string, error) {
	env := make(map[string]string)
	for _, name := range requiredEnv {
		var exists bool
		env[name], exists = os.LookupEnv(name)
		if !exists {
			return nil, errors.New("missing required env " + name)
		}
	}
	return env, nil
}
