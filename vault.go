package goapp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	vault "github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/pingolabscl/go-app/errors"
)

type VaultConfig struct {
	Addr     string
	Username string
	Password string
}

// determineVaultConfig determines the vault config from the
// environment variables. In no environment variables are
// found, the default config will be returned.
func determineVaultConfig(config VaultConfig) VaultConfig {
	envAddr := os.Getenv("VAULT_ADDR")
	envUsername := os.Getenv("VAULT_USERNAME")
	envPassword := os.Getenv("VAULT_PASSWORD")
	if envAddr == "" {
		return config
	}
	if envUsername == "" {
		return config
	}
	if envPassword == "" {
		return config
	}

	return VaultConfig{
		Addr:     envAddr,
		Username: envUsername,
		Password: envPassword,
	}
}

func (app *App) renewVaultToken(auth *vault.ResponseAuth) {
	if app.vault == nil {
		return
	}

	for {
		if !auth.Renewable {
			app.Logger.Error(
				app.ctx,
				"vault_token_renewal_error",
				errors.New("vault token is not renewable"),
			)
			continue
		}

		secs := time.Duration(float64(auth.LeaseDuration)*0.66) * time.Second
		time.Sleep(secs)
		if _, err := app.vault.Auth.TokenRenewSelf(context.Background(), schema.TokenRenewSelfRequest{
			Increment: fmt.Sprintf("%d", auth.LeaseDuration),
		}); err != nil {
			app.Logger.Error(
				app.ctx,
				"vault_token_renewal_error",
				errors.New(err.Error()),
			)
			continue
		}

		app.Logger.Info(
			app.ctx,
			"vault_token_renewed",
		)
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
