package vaultLib

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/kubernetes"
)

var (
	vaultAddress = os.Getenv("VAULT_ADDR")
	vaultRole    = os.Getenv("VAULT_ROLE")
)

type VaultClient struct {
	client *api.Client
}

func NewClient() (*VaultClient, error) {
	client, err := api.NewClient(&api.Config{Address: vaultAddress})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	k8sAuth, err := auth.NewKubernetesAuth(
		vaultRole,
		auth.WithServiceAccountTokenPath("/var/run/secrets/kubernetes.io/serviceaccount/token"),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Kubernetes auth method: %w", err)
	}

	authInfo, err := client.Auth().Login(context.Background(), k8sAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to log in with Kubernetes auth: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	fmt.Println("Succesfully created Vault session")
	return &VaultClient{client: client}, nil
}

func (v *VaultClient) Write(ctx context.Context) error {
	_, err := v.client.Logical().Write("dpce/app/agena/s1", map[string]interface{}{
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
	secret, err := v.client.Logical().Read("dpce/app/agena/s1")
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("read response: %+v\n", secret.Data)
	return nil
}
