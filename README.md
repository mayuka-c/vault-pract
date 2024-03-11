# Vault Pract

## PreRequisites
### Vault Helm Installation
```bash
git clone https://github.com/hashicorp/vault-helm.git
helm install --name-template=vault  --set='server.dev.enabled=true' ./vault-helm --namespace="vault-demo"
```

## Build docker image
```bash
    docker build -t vault-pract .
```

Import the image to k3d
```bash
    k3d image import vault-pract -c mycluster
```