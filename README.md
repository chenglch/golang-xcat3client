## golang-xcat3client
This is a command line client written in golang to work with [xcat3](https://github.com/chenglch/xcat3)
for the technique discussion purpose. A python version client is located at
[python-xcat3client](https://github.com/chenglch/python-xcat3client).

## Supported operating systems

* Ubuntu 14.04, 16.04
* Redhat >= 7.0


## Setup

Setup golang and the environment of GOPATH and GOBIN at first. Then running the following command:

```
go get -v github.com/chenglch/golang-xcat3client
mv $GOBIN/golang-xcat3client $GOBIN/xcat3

```
If you setup xcat3 with [xcat-play](https://github.com/chenglch/xcat-play), please run the following
command:
```
source /etc/profile.d/xcat3.sh
```
Otherwise:
```
export XCAT3_URL=http://<xcat3 daenib ip>:<xcat3 port>
```

## Usage

```
[root@xcat3 ~]# xcat3 help
xcat3 --help and xcat3 help COMMAND to see the usage for specfied
           command.

Usage:
  xcat3 [command]

Available Commands:
  bootdev     Set/Get next boot device (net or disk or cdrom).
  create      Enroll node(s) into xCAT3 service
  delete      Unregister node(s) from the xCAT3 service.
  deploy      Deploy node(s) into specified state.
  export      Export node(s) information as a specific json data file.
  help        Help about any command
  import      Import node(s) information from json data file.
  list        List node(s) in xCAT3 service
  network     This is network child command for xcat3
  nic         This is nic child command for xcat3
  osimage     This is osimage child command for xcat3
  passwd      This is passwds child command for xcat3
  power       Power operation on/off/reset/status for nodes.
  service     This is service child command for xcat3
  show        Show detailed information about node(s).
  update      Update information about registered node(s).

Flags:
  -h, --help   Help message for xcat3

Use "xcat3 [command] --help" for more information about a command.

[root@xcat3 ~]# xcat3 help network
xcat3 network --help and xcat3 network help COMMAND to see the usage for specfied
           command.

Usage:
  xcat3 network [command]

Available Commands:
  create      Register network into xCAT3 service.
  delete      Unregister network from xCAT3 service.
  list        List network(s) in xCAT3 service
  show        Show detailed infomation about network.
  update      Update information about registered network.

Flags:
  -h, --help   help for network

Use "xcat3 network [command] --help" for more information about a command.
```

## Example

```
xcat3 create node0 ipmi=ipmi netboot=pxe arch=x86_64 \
  --nic mac=43:87:0a:05:00:00,ip=12.0.0.1,name=eth0 \
  --nic mac=43:87:0a:05:00:01,ip=13.0.0.1,name=eth1 \
  --control bmc_address=11.0.0.0,bmc_password=password,bmc_username=admin

xcat3 show node0
xcat3 bootdev node0 net
xcat3 power node0 boot
```