# Terraform zabbix provider

Allows to manage zabbix hosts

### Sample config

```
provider "zabbix" {
  user = "admin"
  password = "zabbix"
  server_url = "http://localhost/zabbix/api_jsonrpc.php"
}

resource "zabbix_host" "zabbix1" {
  host = "127.0.0.1"
  interfaces = [{
    ip = "127.0.0.1"
    main = true
  }]
  groups = ["Linux servers", "${zabbix_host_group.zabbix.name}"]
  templates = ["Template ICMP Ping"]
}

resource "zabbix_host_group" "zabbix" {
  name = "something"
}
```

### Build and publish to terraform

```
go get github.com/GarnerCorp/terraform-provider-zabbix
# If you're on Windows / Plan9, this won't work
# This assumes that you only have a single item in your gopath
cd $(go env GOPATH)/src/github.com/GarnerCorp/terraform-provider-zabbix
make -s release -j10

mkdir -p ~/.terraform/plugins/
rsync -a release/terraform-provider-zabbix/* ~/.terraform/plugins/
```

### Terraform

```
cd $YOUR_TERRAFORM_PROJECT
terraform init
...
```
