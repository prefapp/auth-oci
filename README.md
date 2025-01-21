[![Pr verify](https://github.com/prefapp/auth-oci/actions/workflows/pr_verify.yaml/badge.svg)](https://github.com/prefapp/auth-oci/actions/workflows/pr_verify.yaml)
# OCI Registry Authentication Tool

## Install

### Linux (amd64)

```shell
wget -O auth-oci.tar.gz https://github.com/prefapp/auth-oci/releases/download/v1.0.0/auth-oci_1.0.0_linux_amd64.tar.gz
tar -xvzf ./auth-oci.tar.gz
sudo mv auth-oci /usr/local/bin/
auth-oci login --help
```

## Commands

### login Command
The login command authenticates you to the Helm registries as defined in the YAML configuration files.

#### Syntax
```bash
auth-oci login [options]
```

#### Available Options:
- `--registries-dir <path>`: Directory where registry YAML files are stored.
- `--releases-registry <url>`: Host for the releases registry.
- `--snapshots-registry <url>`: Host for the snapshots registry.
- `--types`: Technology types. For now the supported technologies are: `helm` (more coming soon).
- `--releases-registry-username <username>`: Username for the releases registry (only for DockerHub, GHCR, or generic).
- `--releases-registry-password <password>`: Password for the releases registry (only for DockerHub, GHCR, or generic).
- `--snapshots-registry-username <username>`: Username for the snapshots registry (only for DockerHub, GHCR, or generic).
- `--snapshots-registry-password <password>`: Password for the snapshots registry (only for DockerHub, GHCR, or generic).

### Examples

#### Authenticating with DockerHub, GHCR, or Generic Registry
If you're using DockerHub, GitHub Container Registry (GHCR), or a generic registry, you can provide the username and password directly in the command line:
```bash
auth-oci login --registries-dir "/path/to/registries" \
  --types helm \
  --releases-registry "my-release-registry.com" \
  --snapshots-registry "my-snapshot-registry.com" \
  --releases-registry-username "myuser" \
  --releases-registry-password "mypassword" \
  --snapshots-registry-username "myuser" \
  --snapshots-registry-password "mypassword"
```

#### Authenticating with AWS ECR (OIDC) or Azure ACR (OIDC)
For AWS or Azure authentication, ensure that your environment is set up with the proper credentials. You do not need to pass username and password.
```bash
auth-oci login --registries-dir "/path/to/registries" \
  --types helm \
  --releases-registry "<account-id>.dkr.ecr.<region>.amazonaws.com" \
  --snapshots-registry "<acr-name>.azurecr.io" 
```

## Registry YAML Configuration
Each registry should be defined in a YAML file with the following structure. The `authStrategy` key specifies which authentication method to use.

### Example: AWS ECR (OIDC)
```yaml
name: myregistry
registry: <account-id>.dkr.ecr.<region>.amazonaws.com
image_type:
  - snapshots
  - releases
default: true
auth_strategy: aws_oidc
base_paths:
  services: "projects"
  charts: "charts"
```

### Example: Azure ACR (OIDC)
```yaml
name: myregistry
registry: prefappacr.azurecr.io
image_type:
  - snapshots
  - releases
default: true
auth_strategy: azure_oidc
base_paths:
  services: "projects"
  charts: "charts"
```

### Example: DockerHub
```yaml
name: myregistry
registry: docker.io
image_type:
  - snapshots
  - releases
default: true
auth_strategy: dockerhub
base_paths:
  services: "projects"
  charts: "charts"
```

### Example: GitHub Container Registry (GHCR)
```yaml
name: myregistry
registry: ghcr.io
image_type:
  - snapshots
  - releases
default: true
auth_strategy: ghcr
base_paths:
  services: "projects"
  charts: "charts"
```

### Example: Generic Registry
```yaml
name: myregistry
registry: myregistry.com
image_type:
  - snapshots
  - releases
default: true
auth_strategy: generic
base_paths:
  services: "projects"
  charts: "charts"
```

## Notes
- **AWS and Azure Authentication**: These strategies use OIDC tokens. Ensure your environment is properly configured for AWS or Azure credentials.
- **DockerHub, GHCR, and Generic Authentication**: For these, you will need to provide a username and password either through the command line or the YAML configuration files.
- **Error Handling**: If a registry is not found or authentication fails, the program will terminate with an error message.



