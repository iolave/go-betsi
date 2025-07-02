package goapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	vault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/pingolabscl/go-app/pkg/errors"
)

type VaultConfig struct {
	Addr  string
	Token string
}

// determineVaultConfig determines the vault config from the
// environment variables. In no environment variables are
// found, the default config will be returned.
//   - VAULT_ADDR: vault address
//   - VAULT_TOKEN: vault token
func determineVaultConfig(config VaultConfig) VaultConfig {
	envAddr := os.Getenv("VAULT_ADDR")
	envToken := os.Getenv("VAULT_TOKEN")
	if envAddr == "" {
		return config
	}
	if envToken == "" {
		return config
	}
	return VaultConfig{
		Addr:  envAddr,
		Token: envToken,
	}
}

func (app *App) renewVaultToken() {
	base := "vault_token_renewal"
	if app.vault == nil {
		app.Logger.Fatal(
			app.ctx,
			fmt.Sprintf("%s_error", base),
			errors.NewWithName("vault_error", "vault is not configured"),
		)
	}

	for {
		res, err := app.vault.Auth.TokenRenewSelf(app.ctx, schema.TokenRenewSelfRequest{
			Increment: "6h",
		})

		if err != nil {
			duration := time.Minute * 15
			app.Logger.ErrorWithData(app.ctx, fmt.Sprintf("%s_error", base), errors.Wrap(err), map[string]any{
				"sleepingFor": fmt.Sprintf("%ds", int(duration.Seconds())),
			})
			time.Sleep(duration)
			continue
		}

		// Lease duration is in seconds
		lease := res.Auth.LeaseDuration
		sleepTime := float64(lease) * 0.6
		//app.Logger.DebugWithData(
		//	app.ctx,
		//	fmt.Sprintf("%s_success", base),
		//	map[string]any{
		//		"auth":        res.Auth,
		//		"sleepingFor": fmt.Sprintf("%ds", int(sleepTime)),
		//	},
		//)
		time.Sleep(time.Duration(sleepTime * float64(time.Second)))
	}
}

// GetSecret takes an environmnet variable key and checks if it's
// value is a valid vault secret path (format `vault:path/to/secret`)
// and returns its value as map.
//
// If the environment variable value is not a valid vault secret path
// it will assume it's value is a valid json string that's going to
// parsed and returned as a map (THIS IS ONLY MENT TO BE USED FOR
// DEVELOPMENT PURPOSES ONLY).
//
// If it fails to retrieve the secret from vault or parse the json
// string it will return an error.
func (app *App) GetSecret(envKey string) (map[string]any, *errors.Error) {
	env := os.Getenv(envKey)
	if env == "" {
		return nil, errors.NewWithName(
			"vault_error",
			fmt.Sprintf("failed to get secret, environment variable %s is not set", envKey),
		)
	}

	data := map[string]any{}
	if !strings.HasPrefix(env, "vault:") {
		err := json.Unmarshal(bytes.NewBufferString(env).Bytes(), &data)
		if err != nil {
			return nil, errors.NewWithNameAndErr(
				"vault_error",
				"failed to unmarshal env secret",
				err,
			)
		}
		return data, nil
	}

	if app.vault == nil {
		return nil, errors.NewWithName(
			"vault_error",
			"failed to get secret, vault is not configured",
		)
	}

	env = strings.TrimPrefix(env, "vault:")

	res, err := app.vault.Secrets.KvV2Read(
		app.ctx,
		env,
		vault.WithMountPath("services"),
	)
	if err != nil {
		return nil, errors.NewWithNameAndErr(
			"vault_error",
			"failed to get vault secret",
			err,
		)
	}

	return res.Data.Data, nil
}
