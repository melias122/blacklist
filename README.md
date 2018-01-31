# **UBNT edgeos-dnsmasq-blacklist dnsmasq DNS Blacklisting and Redirection**

[![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/britannic/blacklist/blob/master/LICENSE.txt)[![Version](https://img.shields.io/badge/version-v0.0.12-green.svg)](https://github.com/britannic/blacklist)[![GoDoc](https://godoc.org/github.com/britannic/blacklist?status.svg)](https://godoc.org/github.com/britannic/blacklist)[![Build Status](https://travis-ci.org/britannic/blacklist.svg?branch=master)](https://travis-ci.org/britannic/blacklist)[![Coverage Status](https://coveralls.io/repos/github/britannic/blacklist/badge.svg?branch=master)](https://coveralls.io/github/britannic/blacklist?branch=master)[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/github.com/britannic/blacklist)

[community.ubnt.com](https://community.ubnt.com/t5/EdgeMAX/Self-Installer-to-configure-Ad-Server-and-Blacklist-Blocking/td-p/1337892)

NOTE: THIS IS NOT OFFICIAL UBIQUITI SOFTWARE AND THEREFORE NOT SUPPORTED OR ENDORSED BY Ubiquiti Networks®

## **Overview**

EdgeMax dnsmasq DNS blacklisting and redirection is inspired by the users at [EdgeMAX Community](https://community.ubnt.com/t5/EdgeMAX/bd-p/EdgeMAX)

## **Copyright**

* Copyright © 2018 Helm Rock Consulting

## **Licenses**

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
1. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
    ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
    WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
    DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
    ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
    (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
    LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
    ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
    (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
    SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

    The views and conclusions contained in the software and documentation are those
    of the authors and should not be interpreted as representing official policies,
    either expressed or implied, of the FreeBSD Project.

## **Latest Release**

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;[![Latest Release](https://img.shields.io/badge/-Latest%20Release-green.svg)](https://github.com/britannic/blacklist/releases/latest)

## **Change Log**

* See [changelog](CHANGELOG.md) for details.

## **Features**

* Adds DNS blacklisting integration to the EdgeRouter configuration
* Generates configuration files used directly by dnsmasq to redirect dns lookups
* Integrated with the EdgeMax OS CLI
* Any FQDN in the blacklist will force dnsmasq to return the configured dns redirect IP address

## **Compatibility**

* edgeos-dnsmasq-blacklist has been tested on the EdgeRouter ERLite-3, ERPoe-5, ER-X and UniFi Security Gateway USG-3 routers
  * versions EdgeMAX: v1.7.0-v1.9.7+hotfix.4, UniFi: v4.4.12-v4.4.18
* integration could be adapted to work on VyOS and Vyatta derived ports, since  EdgeOS is a fork and port of Vyatta 6.3

## **Installation**

### **EdgeRouter ERLite-3, ERPoe-5 & UniFi-Gateway-3**

```bash
curl https://community.ubnt.com/ubnt/attachments/ubnt/EdgeMAX/194030/5/edgeos-dnsmasq-blacklist_0.0.12_mips.deb.tgz | tar -xvz
sudo dpkg -i edgeos-dnsmasq-blacklist_0.0.12_mips.deb
```

### **EdgeRouter ER-X & ER-X-SFP**

```bash
curl
https://community.ubnt.com/ubnt/attachments/ubnt/EdgeMAX/194030/6/edgeos-dnsmasq-blacklist_0.0.12_mipsel.deb.tgz | tar -xvz
sudo dpkg -i edgeos-dnsmasq-blacklist_0.0.12_mipsel.deb
```

## **Upgrade**

* Since dpkg cannot upgrade packages, follow the instructions under [Installation](#installation) and the previous package version will be automatically removed before the new package version is installed

## **Removal**

### **EdgeMAX - All Platforms**

```bash
sudo apt-get remove edgeos-dnsmasq-blacklist
```

## **Frequently Asked Questions**

**How do I disable/enable dnsmasq blacklisting?**

* Use these CLI configure commands:
* Disable:

```bash
configure
set service dns forwarding blacklist disabled true
commit;save;exit
```

* Enable:

```bash
configure
set service dns forwarding blacklist disabled false
commit;save;exit
```

**How do I back up my blacklist configuration before upgrading and restore it afterward?**

* use the following commands (make a note of the file name):

```bash
export DATE=$(date +'%FT%H%M%S'); echo "Backing up blacklist configuration to: /config/user-data/blacklist.${DATE}.cmds"; show configuration commands | grep blacklist > /config/user-data/blacklist.$(date +'%FT%H%M%S').cmds
```

* After installing the latest version, you can merge your backed up configuration:

```bash
configure
.  /config/user-data/blacklist.[date string].cmds
commit;save;exit
```

* If you prefer to delete the default configuration and restore your previous configuration, run these commands:

```bash
configure
delete service dns forwarding blacklist
.  /config/user-data/blacklist.[date string].cmds
commit;save;exit
```

**Which blacklist sources are installed by default?**

* You can use this command in the CLI shell to view the current sources after installation or view the log and see previous downloads:

```bash
show configuration commands | match blacklist | match source
more /var/log/update-dnsmasq.log
```

**How do I add or delete sources?**

* Using the CLI configure command, to delete domains and hosts sources:

```bash
configure
delete service dns forwarding blacklist domains source malc0de
delete service dns forwarding blacklist hosts source yoyo.org
commit;save;exit
```

* To add a source, first check it can serve a text list and also note the prefix (if any) before the hosts or domains, e.g. [http://www.malwaredomainlist.com/](http://www.malwaredomainlist.com/) has this format:

```text
#               MalwareDomainList.com Hosts List           #
#   http://www.malwaredomainlist.com/hostslist/hosts.txt   #
#         Last updated: Mon, 04 Dec 17 19:18:42 +0000      #


127.0.0.1  localhost
127.0.0.1  0koryu0.easter.ne.jp
127.0.0.1  109-204-26-16.netconnexion.managedbroadband.co.uk
127.0.0.1  1866809.securefastserver.com
```

* So the prefix is "127.0.0.1  "
* Here's how to creating the source in the CLI:

```bash
configure
set service dns forwarding blacklist hosts source malwaredomainlist description '127.0.0.1 based host and domain list'
set service dns forwarding blacklist hosts source malwaredomainlist prefix '127.0.0.1  '
set service dns forwarding blacklist hosts source malwaredomainlist url 'http://www.malwaredomainlist.com/hostslist/hosts.txt'
commit;save;exit
```

**How do I globally exclude or include hosts or a domains?**

* Use these commands to globally include or exclude blacklisted entries:

```bash
configure
set service dns forwarding blacklist exclude cdn.visiblemeasures.com
set service dns forwarding blacklist include www.nastywebsites.com
commit;save;exit
```

**How do I exclude or include a host or a domain?**

* Use these commands to include or exclude blacklisted entries:

```bash
configure
set service dns forwarding blacklist domains exclude visiblemeasures.com
set service dns forwarding blacklist domains include domainsnastywebsites.com
set service dns forwarding blacklist hosts exclude cdn.visiblemeasures.com
set service dns forwarding blacklist hosts include www.nastywebsites.com
commit;save;exit
```

**Does update-dnsmasq run automatically?**

* Yes, a scheduled task is created and run daily at midnight
* It can be reconfigured using these CLI configuration commands:

```bash
set system task-scheduler task update_blacklists executable path /config/scripts/update-dnsmasq
set system task-scheduler task update_blacklists interval 1d
```

* For example, to change the execution interval to every 6 hours, use this command:

```bash
set system task-scheduler task update_blacklists interval 6h
```

* In daily use, no additional interaction with update-dnsmasq is required. By default, cron will run update-dnsmasq at midnight each day to download the blacklist sources and update the dnsmasq configuration files in /etc/dnsmasq.d. dnsmasq will automatically be reloaded after the configuration file update is completed.

**How I use the command line switches?**

* update-dnsmasq has the following commandline switches available:

```bash
/config/scripts/update-dnsmasq.amd64 -h
    -dir string
            Override dnsmasq directory (default "/etc/dnsmasq.d")
    -f [full file path]
            [full file path] # Load a config.boot file
    -h   Display help
    -v   Verbose display
    -version
            Show version
```

### **EdgeOS dnsmasq Configuration**

* dnsmasq may need to be configured to ensure blacklisting works correctly
  * Here is an example using the EdgeOS configuration shell

```bash
configure
set service dns forwarding cache-size 2048
set service dns forwarding except-interface [Your WAN i/f]
set service dns forwarding name-server [Your choice of IPv4 Internet Name-Server]
set service dns forwarding name-server [Your choice of IPv4 Internet Name-Server]
set service dns forwarding name-server [Your choice of IPv6 Internet Name-Server]
set service dns forwarding name-server [Your choice of IPv6 Internet Name-Server]
set service dns forwarding options bogus-priv
set service dns forwarding options domain-needed
set service dns forwarding options domain=mydomain.local
set service dns forwarding options enable-ra
set service dns forwarding options expand-hosts
set service dns forwarding options localise-queries
set service dns forwarding options strict-order
set service dns forwarding system
set system name-server 127.0.0.1
set system name-server '::1'
commit; save; exit
```

> blacklist





- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
