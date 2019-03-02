Parkomat.io [![Build Status](https://travis-ci.org/parkomat/parkomat.svg?branch=master)](https://travis-ci.org/parkomat/parkomat)
-----------------------------------------------------------------------------------------------------------------------------------

### What is it?

DNS + Web + WebDav server in one package.

### Features

- DNS server with catch-all function 
- Web server with SSL support (can run many certificates on one IP)
- WebDav for easy upload of files to the web

### Why ?

Parkomat is useful when you have a lot of domains and managing them via typical hosting panel becomes too complex.

### Installation

Parkomat at the moment doesn't provide pre-built binaries, so you need to have Go 1.5+ installed. Latest version of Go is recommended.

To build, issue:

```
go get github.com/parkomat/parkomat
```

### Setting up

As a configuration format Parkomat uses [TOML](https://github.com/toml-lang/toml)

### Try with Docker

```
docker pull parkomat/parkomat
```

```
docker run -d -e PARKOMAT_CONFIG_FILE=/opt/parkomat/config.toml -v /your/parkomat/directory:/opt/parkomat -p 53:53/udp parkomat/parkomat
```

Remember to have `config.toml` file in your `/your/parkomat/directory` path.

### Example Configuration:

Note: instead of `127.0.0.1` use your external IP.

```
# if you set it to true, Parkomat will serve any domain pointing at it
catch_all = true

[[domains]]
name = "example.domain"

[[domains]]
name = "parkomat.io"
	# supports per domain zone settings
	[domains.zone]
	A = "192.168.0.1"
	MX = """
1 better.mail.server
"""
	TXT = """
hello world
"""

# each domain will use following zone settings
[zone]
# for both .domain and www.domain
A = "127.0.0.1"
MX = '''
1 test1.mail.server
10 test2.mail.server
'''

[web]
ip = "0.0.0.0"
port = 80
path = "./www"

# make sure that path exists
# for example issue mkdir -p /var/log/parkomat
access_log = "/var/log/parkomat/access.log"

[webdav]
enabled = true
username = "hello"
password = "world"
# your share will be under http://example.domain/dav/
mount = "/dav/"

[dns]
ip = "127.0.0.1"
port = 53

# details of dns servers for NS record
[[dns.servers]]
name = "ns1.parkomat.co"
ip = "127.0.0.1"

[[dns.servers]]
name = "ns2.parkomat.co"
ip = "127.0.0.1"
```

Make sure to create `GLUE` record for each dns server listed in `[[dns.servers]]`. You need to follow your registrar documentation on how to do it.

You can run multiple parkomat nodes for DNS server. Make sure they use the same configuration file (for example mounted via NFS).

To run parkomat in DNS only mode, use:

```
./parkomat -dns_only=true -config_file=/path/to/config.toml
```

You can also use following environment variables, that will overwrite passed arguments:

`PARKOMAT_CONFIG_FILE` - path to the configuration file, for example `/path/to/config.toml`

`PARKOMAT_DNS_ONLY` - `true` or `false` for DNS only mode

### Web server directory structure

You `./web` path could look like this:

```
.
├── default
│   └── public_html
│       └── index.html
├── parkomat.io
|   ├── parkomat.io.crt
|   ├── parkomat.io.key
|   └── public_html
|       └── index.html
└── config.toml
```

To add new domain, simply create new directory with that domain name.
If you want to use SSL, just copy `domain.crt` and `domain.key` files to that domain directory (be careful - do not upload them to `public_html` directory). You need to restart parkomat afterwards (SSL at the moment is not reloaded at runtime).

All your html and other files go to `public_html` directory.

### WebDav

If you want to use WebDav with windows, the domain you will be using it with should have certificates uploaded. Apparently WebDav doesn't work without SSL on Windows.

### TO DO

- Better documentation
- API
- Stats
- Mail forwarding
- Live reload of configuration
- ???

