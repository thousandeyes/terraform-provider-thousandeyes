# Endpoint Agent manual acceptance checks

This workspace exercises the import-only `thousandeyes_endpoint_agent` resource against a real Endpoint Agent. It expects credentials and an existing agent to be supplied through environment variables.

## Configure the local provider

From the repository root, build the provider:

```sh
make build
```

Create a Terraform CLI configuration file for the local build, replacing the path with the absolute path to this repository:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/thousandeyes/thousandeyes" = "/absolute/path/to/terraform-provider-thousandeyes"
  }
  direct {}
}
```

Set the CLI config and ThousandEyes credentials. `TE_AID` is optional when the token's current account group is the intended one.

```sh
export TF_CLI_CONFIG_FILE="/absolute/path/to/terraform.rc"
export TE_TOKEN="<your-token>"
export TE_AID="<your-account-group-id>"
export TE_ENDPOINT_AGENT_ID="<existing-endpoint-agent-id>"
```

Terraform development overrides use the local provider binary and do not require `terraform init`. See [Terraform development overrides](https://developer.hashicorp.com/terraform/plugin/debugging).

## Import and read

From this directory, import the existing agent and inspect the refreshed state:

```sh
./import.sh
terraform state show thousandeyes_endpoint_agent.manual
terraform plan
```

The resource configuration intentionally leaves mutable fields unset so the first plan can confirm the imported remote values without requesting updates. The state should include the remote ID, name, license type, and status.

## Check updates and status

Edit `resource.tf` to configure a desired `name`, supported `license_type` (`essentials`, `advantage`, or `embedded`), or `status` (`enabled` or `disabled`), then run:

```sh
terraform plan
terraform apply
terraform state show thousandeyes_endpoint_agent.manual
```

Change `status` to the opposite of its current value to exercise the enable or disable operation. Change it back and apply again if you want to restore the original status. Only use license changes that are appropriate for the selected test agent and account.

For a wireless license, leave `license_type` unset and confirm `terraform state show` reports the remote `embed-wireless-advantage` or `embed-wireless-essentials` value. To check assignment validation, temporarily add one of those values to `resource.tf` and run `terraform plan`; Terraform should report that the value can be read but cannot be assigned. Remove it from configuration afterward.

## Check create rejection

In a separate fresh directory, use `create_rejection/main.tf` with the same local provider override and credentials. Run `terraform plan` and `terraform apply`. The apply should fail with the import-only diagnostic and make no API create request.

## Cleanup

Do not run `terraform destroy` against an agent you want to keep. This resource's delete operation deletes the Endpoint Agent from ThousandEyes. To finish a non-destructive manual check, remove only the local Terraform state with:

```sh
terraform state rm thousandeyes_endpoint_agent.manual
```
