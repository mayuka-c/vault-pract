package vaultclientgo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

var (
	vaultAddress = os.Getenv("VAULT_ADDR")
	vaultRole    = os.Getenv("VAULT_ROLE")
)

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

type VaultClient struct {
	client *vault.Client
}

func NewClient(ctx context.Context) (*VaultClient, error) {
	client, err := vault.New(
		vault.WithAddress(vaultAddress),
		vault.WithRequestTimeout(60*time.Second),
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("client successfully created")

	jwt := fetchServiceAccountToken()

	resp, err := client.Auth.KubernetesLogin(context.Background(), schema.KubernetesLoginRequest{
		Jwt:  jwt,
		Role: vaultRole,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Kubernetes login succeeded")

	if err := client.SetToken(resp.Auth.ClientToken); err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("Successfully set the token")
	return &VaultClient{client: client}, nil
}

func (v *VaultClient) Write(ctx context.Context) error {
	_, err := v.client.Write(ctx, "dpce/app/agena/s1", map[string]any{
		"data": map[string]any{
			"master_tenant_password": "abc123",
			"custom_sso_encKey":      "pass-1",
		},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Successfully wrote the secret into the path")
	return nil
}

func (v *VaultClient) Read(ctx context.Context) error {
	readResp, err := v.client.Read(ctx, "dpce/app/agena/s1")
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("read response: %+v\n", readResp.Data)
	return nil
}
