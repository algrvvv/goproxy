# GoProxy - Simple proxy server in golang

![GoProxy](https://i.imgur.com/uPDYT2n_d.webp?maxwidth=760&fidelity=grand)

### Installation

The installation has two options:

1. clone the repository and make a build
2. install the binary for your system. Exceptions: macOS

```shell
# download the installer
$ curl -LO https://raw.githubusercontent.com/algrvvv/goproxy/main/install.sh
# make the file executable
$ chmod +x install.sh
# start the installation
$ ./install.sh
```

After installation, edit `config.yaml` to suit your needs and start the server:

```shell
# if you chose the first installation option
$ ./bin/goproxy
# if you chose the first installation option
$ ./goproxy
```

### Additions

in addition, you can view a [browser extension](https://github.com/algrvvv/goproxy-ext) that can be used in conjunction
with this proxy server
