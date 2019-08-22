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

### Installing a release from GitHub

You'll want to download the appropriate version of the provider

#### Linux/macOS

This shell script should download and install it to the appropriate place:

```sh
mkdir -p ~/.terraform.d/plugins;
unamesys=$(uname -s|tr A-Z a-z);
unamearch=$(uname -m|sed -e 's/x86_64/amd64/;s/i\[3-6\]86/386/');
FILE="${unamesys}_${unamearch}.tar.gz";
RELEASE=$(curl -Ls -o /dev/null -w %{url_effective} https://github.com/GarnerCorp/terraform-provider-zabbix/releases/latest/|sed -e 's#/tag/#/download/#');
curl -Ls -O $RELEASE/$FILE;
tar zxf $FILE;
rsync -a $(echo $FILE |sed -e 's/\..*//') ~/.terraform.d/plugins
```

#### Windows

Visit the [releases](https://github.com/GarnerCorp/terraform-provider-zabbix/releases/latest/) page and download the appropriate file:

* windows_386.tar.gz
* windows_amd64.tar.gz

1. Extract it: `tar zxf windows*.tar.gz`
2. Copy the expanded folder to `%HOME%\.terraform.d\plugins\`

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
