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
	"github.com/pingolabscl/go-app/errors"
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
			errors.New("vault is not configured"),
		)
	}

	for {
		res, err := app.vault.Auth.TokenRenewSelf(app.ctx, schema.TokenRenewSelfRequest{
			Increment: "60m",
		})
		if err != nil {
			duration := time.Minute * 1
			app.Logger.ErrorWithData(app.ctx, fmt.Sprintf("%s_error", base), errors.New(err.Error()), map[string]any{
				"sleeping_for": fmt.Sprintf("%ds", int(duration.Seconds())),
			})
			time.Sleep(duration)
			continue
		}

		app.Logger.Info(app.ctx, fmt.Sprintf("%s_success", base))
		duration := time.Duration(float64(res.Auth.LeaseDuration) * 0.1)
		app.Logger.DebugWithData(
			app.ctx,
			fmt.Sprintf("%s_sleeping", base),
			map[string]any{
				"secs": duration.Seconds(),
			},
		)

		time.Sleep(duration)
	}
}

// GetSecret checks if the environment variable value is a valid
// vault secret path and returns the value as a map[string]any.
//
// If the environment variable value is not a valid vault secret path
// it will try to parse the value as a JSON string and return the
// value as a map[string]any (THIS IS ONLY MENT TO BE USED
// FOR DEVELOPMENT PURPOSES ONLY).
func (app *App) GetSecret(envKey string) (map[string]any, error) {
	env := os.Getenv(envKey)
	if env == "" {
		return nil, errors.NewInternalServerError(
			"failed to get secret",
			fmt.Sprintf("environment variable %s is not set", envKey),
		)
	}

	data := map[string]any{}
	if !strings.HasPrefix(env, "vault:") {
		err := json.Unmarshal(bytes.NewBufferString(env).Bytes(), &data)
		if err != nil {
			return nil, errors.NewInternalServerError(
				"failed to parse env secret)",
				err.Error(),
			)
		}
		return data, nil
	}

	if app.vault == nil {
		return nil, errors.NewInternalServerError(
			"failed to get secret",
			"vault is not configured",
		)
	}

	env = strings.TrimPrefix(env, "vault:")

	res, err := app.vault.Secrets.KvV2Read(
		app.ctx,
		env,
		vault.WithMountPath("services"),
	)
	if err != nil {
		return nil, errors.NewInternalServerError(
			"failed to get vault secret",
			err.Error(),
		)
	}

	return res.Data.Data, nil
}
