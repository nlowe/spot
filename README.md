# Spot - The Watchdog for your Build Agents

Spot is a watchdog for build agents in Jenkins and Bamboo

## Building

You need [`dep`](https://github.com/golang/dep). The easiest way to build
is to run `make`, which will generate linux and windows binaries in `dist/`.

If you don't have `make`, you can build manually:

```bash
# linux
dep ensure
go build -o dist/mmcmd -v ./cmd/mmcmd

# windows
dep ensure
go build -o dist/mmcmd.exe -v ./cmd/mmcmd
```

## Usage

```txt
alerts for disconnected build agents
Usage: main.exe [--jenkins JENKINS] [--verbosity VERBOSITY] [--once]

Options:
  --jenkins JENKINS, -j JENKINS
                         Jenkins Url & credentials in the form of https://jenkins/,username,password
  --verbosity VERBOSITY, -v VERBOSITY
                         Verbosity [panic, fatal, error, warn, info, debug] [default: info]
  --once, -o             Run checks once and exit
  --help, -h             display this help and exit
```

### Example

```txt
$ go run .\cmd\spot\main.go --jenkins "https://devops.jenkins.hylandqa.net,username,password" --jenkins "https://csp.jenkins.hylandqa.net/" --once --verbosity debug
INFO[0000] Hello, World!
DEBU[0000] Trying to parse jenkins instance              jenkins="https://devops.jenkins.hylandqa.net,username,password"
DEBU[0000] Trying to parse jenkins instance              jenkins="https://csp.jenkins.hylandqa.net/"
INFO[0000] Running Watchdog Task
DEBU[0000] Checking for offline agents                   detector="[jenkins] https://devops.jenkins.hylandqa.net"
DEBU[0001] Node is online                                agent=master detector="[jenkins] https://devops.jenkins.hylandqa.net"
DEBU[0001] Node is online                                agent=RDV-003960.hylandqa.net detector="[jenkins] https://devops.jenkins.hylandqa.net"
DEBU[0001] Node is online                                agent=RDV-004063.hylandqa.net detector="[jenkins] https://devops.jenkins.hylandqa.net"
DEBU[0001] Check Complete                                detector="[jenkins] https://devops.jenkins.hylandqa.net"
DEBU[0001] Checking for offline agents                   detector="[jenkins] https://csp.jenkins.hylandqa.net"
DEBU[0001] Node is online                                agent=master detector="[jenkins] https://csp.jenkins.hylandqa.net"
DEBU[0001] Node is online                                agent="RDV-004097.hylandqa.net (QAV Performance)" detector="[jenkins] https://csp.jenkins.hylandqa.net"
DEBU[0001] Node is online                                agent="Ubuntu-Docker (10.40.0.120)" detector="[jenkins] https://csp.jenkins.hylandqa.net"
DEBU[0001] Node is online                                agent=Windows-Server-1709-Docker-0 detector="[jenkins] https://csp.jenkins.hylandqa.net"
DEBU[0001] Node is online                                agent=Windows-Server-1709-Docker-1 detector="[jenkins] https://csp.jenkins.hylandqa.net"
DEBU[0001] Node is online                                agent=Windows-Server-1709-Docker-2 detector="[jenkins] https://csp.jenkins.hylandqa.net"
DEBU[0001] Node is online                                agent=Windows-Server-1709-Docker-3 detector="[jenkins] https://csp.jenkins.hylandqa.net"
DEBU[0001] Node is online                                agent=Windows-Server-1709-Docker-4 detector="[jenkins] https://csp.jenkins.hylandqa.net"
DEBU[0001] Check Complete                                detector="[jenkins] https://csp.jenkins.hylandqa.net"
INFO[0001] Goodbye
```
