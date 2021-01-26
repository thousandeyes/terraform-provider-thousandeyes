Terraform `thousandeyes` Provider
=========================

- Website: [https://registry.terraform.io/providers/william20111/thousandeyes/latest](https://registry.terraform.io/providers/william20111/thousandeyes/latest)
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is community maintained

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/william20111/terraform-provider-thousandeyes`

```sh
$ git clone git@github.com:william20111/terraform-provider-thousandeyes $GOPATH/src/github.com/william20111/terraform-provider-thousandeyes
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/william20111/terraform-provider-thousandeyes
$ make build
```

Using the provider
----------------------
The provider is now on the terraform registry to pull it as a dependency do the following and run `terraform init`

```hcl
terraform {
  required_providers {
    thousandeyes = {
      source = "william20111/thousandeyes"
      version = "0.3.3"
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

