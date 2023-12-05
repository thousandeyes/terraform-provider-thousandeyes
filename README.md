# Terraform Provider for ThousandEyes [![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/thousandeyes/terraform-provider-thousandeyes?label=release)](https://github.com/thousandeyes/terraform-provider-thousandeyes/releases) [![license](https://img.shields.io/github/license/thousandeyes/terraform-provider-thousandeyes.svg)]()

The Terraform provider for ThousandEyes allows you to manage resources in [ThousandEyes](https://www.thousandeyes.com/).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.17 (to build the provider plugin)

## Usage
The provider is on the Terraform registry. To use it, add the following code and run `terraform init`:

```hcl
terraform {
  required_providers {
    thousandeyes = {
      source = "thousandeyes/thousandeyes"
      version = ">= 1.3.1"
    }
  }
}
```

### Setting up provider
```hcl
provider "thousandeyes" {
  token = "xxxxxx-xxxx-xxxx-xxxx-xxxxxxxxx"
}

```

The provider requires a token. The token can be set on the `token` variable, as shown in the example, or it may instead be passed via the `TE_TOKEN` environment variable.

The provider also supports the following optional settings:

- `account_group_id` may be set to distinguish between affected account groups, if your ThousandEyes account supports more than one.  This may instead be set by the environment variable `TE_AID`.
- `timeout` may be set to specify the number of seconds to wait for responses from the ThousandEyes endpoints.  This may instead be set by the environment variable `TE_TIMEOUT`.  If this is unset or set to `0`, then the thousandeyes-sdk-go library will use its default settings.

### Examples
Example of an HTTP test:

```hcl
data "thousandeyes_agent" "arg_cordoba" {
  agent_name = "Cordoba, Argentina"
}

resource "thousandeyes_http_server" "www_thousandeyes_http_test" {
  test_name      = "Example HTTP test set from Terraform provider"
  interval       = 120
  alerts_enabled = false

  url = "https://www.thousandeyes.com"

  agents {
    agent_id = data.thousandeyes_agent.arg_cordoba.agent_id
  }
}
```

### Supported tests
- [X] agent-to-agent
- [X] agent-to-server
- [X] bgp
- [X] dnssec
- [X] dns-server
- [X] dns-trace
- [X] ftp-server
- [X] http-server
- [X] page-load
- [X] sip-server
- [X] voice (RTP stream)
- [X] web-transactions

## Building The Provider
Clone repository to: `$GOPATH/src/github.com/thousandeyes/terraform-provider-thousandeyes`

```sh
$ git clone git@github.com:thousandeyes/terraform-provider-thousandeyes $GOPATH/src/github.com/thousandeyes/terraform-provider-thousandeyes
```

Enter the provider directory and build the provider:

```sh
$ cd $GOPATH/src/github.com/thousandeyes/terraform-provider-thousandeyes
$ make build
```

Follow the instructions to [install it as a plugin](https://developer.hashicorp.com/terraform/plugin#installing-a-plugin). After placing it into your plugins directory,  run `terraform init` to initialize it.

## Maintainers
This provider plugin is maintained by the ThousandEyes engineering team and accepts community contributions.

> [!NOTE]  
> The `docs/` folder should not be changed manually. Instead, if there are changes to the examples, inputs/outputs of any `data_source` or `resource`, you need to run `go generate` locally and commit the resulting changes to the `.md` files. 


## Acknowledgements
ThousandEyes would like to thank William Fleming, John Dyer, and Joshua Blanchard for their contribution and community maintenance of this project.
