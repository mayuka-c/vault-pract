package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	vaultclientgo "github.com/mayuka-c/vault-pract/vault-client-go"
	"github.com/mayuka-c/vault-pract/vaultLib"
)

var (
	ctx = context.Background()
)

const secretFile = "/vault/secrets/cassandra"

type DB struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func readDBSecrets() {
	secBytes, err := os.ReadFile(secretFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	db := DB{}
	json.Unmarshal(secBytes, &db)

	fmt.Printf("DB struct: %+v\n", db)
}

func main() {
	readDBSecrets()

	vaultLib := os.Getenv("VAULT_CLIENT")
	var err error

	if vaultLib == "vaultLib" {
		err = useVault()
	} else {
		err = useVaultClientGo()
	}
	if err != nil {
		panic(err)
	}
}

func useVault() error {
	fmt.Println("Use Vault Lib")
	vault, err := vaultLib.NewClient()
	if err != nil {
		return err
	}

	err = vault.Write(ctx)
	if err != nil {
		return err
	}

	err = vault.Read(ctx)
	if err != nil {
		return err
	}

	return nil
}

func useVaultClientGo() error {
	fmt.Println("Use Vault Client Go")
	vault, err := vaultclientgo.NewClient(ctx)
	if err != nil {
		return err
	}

	err = vault.Write(ctx)
	if err != nil {
		return err
	}

	err = vault.Read(ctx)
	if err != nil {
		return err
	}

	return nil
}
