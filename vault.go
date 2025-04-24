package goapp

import (
	"encoding/json"
	"fmt"
	"os"
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
			app.Logger.Error(app.ctx, fmt.Sprintf("%s_error", base), errors.New(err.Error()))
			duration := time.Minute * 1
			app.Logger.DebugWithData(
				app.ctx,
				fmt.Sprintf("%s_sleeping", base),
				map[string]any{
					"secs": duration.Seconds(),
				},
			)
			time.Sleep(duration)
			continue
		}

		app.Logger.Info(app.ctx, fmt.Sprintf("%s_success", base))
		duration := time.Duration(float64(res.Auth.LeaseDuration)*0.6) * time.Second
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

func (app *App) GetVaultSecret(secretPath string, result any) error {
	res, err := app.vault.Secrets.KvV2Read(
		app.ctx,
		secretPath,
		vault.WithMountPath("services"),
	)
	if err != nil {
		return errors.NewInternalServerError(
			"failed to get vault secret",
			err.Error(),
		)
	}
	b, err := json.Marshal(res.Data.Data)
	if err != nil {
		return errors.NewInternalServerError(
			"failed to parse vault secret",
			err.Error(),
		)
	}

	json.Unmarshal(b, result)
	return nil
}
