Terraform `thousandeyes` Provider
=========================

- Website: [https://registry.terraform.io/providers/thousandeyes/thousandeyes/latest](https://registry.terraform.io/providers/thousandeyes/thousandeyes/latest)
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by the ThousandEyes engineering team and accepts community contributions.

Acknowledgements
----------------

ThousandEyes would like to thank William Fleming, John Dyer, and Joshua Blanchard for their contribution and community maintenance of this project.

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.16 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/thousandeyes/terraform-provider-thousandeyes`

```sh
$ git clone git@github.com:thousandeyes/terraform-provider-thousandeyes $GOPATH/src/github.com/thousandeyes/terraform-provider-thousandeyes
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/thousandeyes/terraform-provider-thousandeyes
$ make build
```

Using the provider
----------------------
The provider is now on the terraform registry to pull it as a dependency do the following and run `terraform init`

```hcl
terraform {
  required_providers {
    thousandeyes = {
      source = "thousandeyes/thousandeyes"
    }
  }
}
```

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

### Setting up provider

```hcl

provider "thousandeyes" {
  token = "xxxxxx-xxxx-xxxx-xxxx-xxxxxxxxx"
}

```

The provider requires that the token be set, the token may instead be passed via the environment variable `TE_TOKEN`.

The provider also supports the following optional settings:

- `account_group_id` may be set to distinguish between affected account groups, if your ThousandEyes account supports more than one.  This may instead be set by the environment variable `TE_AID`.
- `timeout` may be set to specify the number of seconds to wait for responses from the ThousandEyes endpoints.  This may instead be set by the environment variable `TE_TIMEOUT`.  If this is unset or set to `0`, then the go-thousandeyes library will use its default settings.

### HTTP Test

```hcl

data "thousandeyes_agent" "test_agent" {
  name  = "na-sjc-2-te [VS01]"
}

resource "thousandeyes_http_server" "google_http_test" {
  name = "google test"
  interval = 120
  url = "https://google.com"
  agents {
      agent_id = data.thousandeyes_agent.test_agent.agent_id
  }
  agents {

      agent_id = 12345
  }
}
```

### Agent to Server

```hcl
data "thousandeyes_agent" "test_agent_example" {
  name  = "na-sjc-2-te [VS01]"
}

resource "thousandeyes_agent_to_server" "agent_test_example" {
  name = "my agent test"
  interval = 120
  server = "8.8.8.8"
  agents {
      agent_id = data.thousandeyes_agent.test_agent_example.agent_id
  }

}
```

Supported tests right now:

- [x] http-server
- [x] page-load
- [x] web-transactions
- [x] agent-to-server
- [x] agent-to-agent
- [x] bgp
- [ ] transactions
- [x] ftp-server
- [x] dns-trace
- [x] dns-server
- [x] dns-dnssec
- [ ] dnsp-domain
- [ ] dnsp-server
- [x] sip-server
- [x] voice (RTP Stream)
- [x] voice-call

