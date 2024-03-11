package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

var (
	ctx = context.Background()
)

const secretFile = "/vault/secrets/database"

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

func fetchServiceAccountToken() string {
	serviceAccountTokenPath := "/var/run/secrets/kubernetes.io/serviceaccount/token"

	fileBytes, err := os.ReadFile(serviceAccountTokenPath)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Println(string(fileBytes))
	return string(fileBytes)
}

func main() {
	readDBSecrets()
	// prepare a client with the given base address
	client, err := vault.New(
		vault.WithAddress("http://vault:8200"),
		vault.WithRequestTimeout(30*time.Second),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("client successfully created")

	jwt := fetchServiceAccountToken()

	resp, err := client.Auth.KubernetesLogin(context.Background(), schema.KubernetesLoginRequest{
		Jwt:  jwt,
		Role: "test-app",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := client.SetToken(resp.Auth.ClientToken); err != nil {
		fmt.Println(err)
		return
	}

	_, err = client.Write(ctx, "app/secret/agena/s1", map[string]any{
		"data": map[string]any{
			"master_tenant_password": "abc123",
			"custom_sso_encKey":      "pass-1",
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	readResp, err := client.Read(ctx, "app/secret/agena/s1")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("read response: %+v\n", readResp.Data)
}
