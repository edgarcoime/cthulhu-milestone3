# Terraform directory to setup infrastucture

## Dependencies

Need to install Terraform first to work:
[Terraform Install](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)

Localstack can be installed for local development:
[Localstack Install](https://docs.localstack.cloud/aws/getting-started/installation/)

For Terraform to work with local stack need a wrapper

```bash
git clone https://github.com/localstack/terraform-local
cd terraform-local
make install  # Or go install
```

TFlocal is installed in the .venv as a child of the repo folder. So it must be sourced first before use.

```bash
source .venv/bin/activate
tflocal
```

## How to run

```bash
# 1. Start local stack
localstack start -d

# 2. Activate terraform local venv
source .venv/bin/activate

# 3. Apply from repo
cd terraform
tflocal init
tflocal apply
```

## How to tear down resources

From the repo, with LocalStack running and the terraform-local venv activated:

```bash
cd terraform
tflocal destroy
```

Confirm with `yes` when prompted. This removes the S3 bucket and its CORS configuration from LocalStack. The bucket and any objects in it will no longer exist.