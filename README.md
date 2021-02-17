# Goserve

A simple webserver for static files served from localhost so that you don't need to use npm packages. TLS encryption included.

## Prerequisites

To use TLS you'll need an installed `openssl` package for certificate generation. However, TLS will only be used on the loopback interface. Certificates will be stored in `$USERCACHE/goserve-certs` (which corresponds to `/home/<user>/.cache` on Linux) and are valid for 365 days.

To suppress the "insecure certificate" warning on Chrom(e/ium), go to [chrome://flags/#allow-insecure-localhost](chrome://flags/#allow-insecure-localhost) and enable the setting.

## Usage

`make` will build the binary and install it in `$GOPATH/bin`. Call `goserve [OPTIONS] <FOLDER>` to start serving. If FOLDER is left out, it will use the current folder.

## Flags

```bash
Usage of goserve:
  -a	Bind to all interfaces (default: Loopback only)
  -c	Use http caching (default: false)
  -p string
    	Port (default "9000")
  -s	Use TLS. Will generate certs if they are not present (default: false)
```
