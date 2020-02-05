Terraform `thousandeyes` Provider
=========================

- Website: https://www.terraform.io
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
If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

Supported tests right now:

- [x] http-server
- [x] page-load
- [ ] agent-to-server
- [ ] agent-to-agent
- [ ] bgp
- [ ] http-server
- [ ] page-load
- [ ] transactions
- [ ] web-transactions
- [ ] ftp-server
- [ ] dns-trace
- [ ] dns-server
- [ ] dns-dnssec
- [ ] dnsp-domain
- [ ] dnsp-server
- [ ] sip-server
- [ ] voice (RTP Stream)
- [ ] voice-call

