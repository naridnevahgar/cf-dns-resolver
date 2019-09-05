# DNS Resolver CF Plugin
## Overview
The plugin exposes the resolve command, which takes a host name and resolves it to an IP address using the standard `dig` command. 
Out of the list of started apps in the current space, the command attempts to open an ssh tunnel in a sequential manner. 
Once there is a successful ssh connection, the host name is resolved. 

The standard `host` command is used to pre-resolve the host name. 
Then, the `dig` command is executed with the argument supplied to `cf resolve` to echo the IP address on console.

### Usage
`cf resolve <bosh based host name>`

### Cautions
For the plugin to resolve successfully, 
1. There should be a valid session connected via `cf login`
2. The currently logged in `Org` should support ssh connectivity in general

### Installation Instructions
Follow the below instructions to install the plugin into your local (for Mac OS)

    git clone https://github.com/naridnevahgar/cf-dns-resolver.git
    cd cf-dns-resolver-plugin
    cf install-plugin dns_resolver

If you have already installed a previous version, you can force re-installation by adding the `force` flag.

    cf install-plugin dns_resolver -f

To confirm the installation, execute the below command:

    cf resolve --help
    

Happy resolving!