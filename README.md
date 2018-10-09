# certmonitor

A simple website certificate monitor tool

## How to use?

### For general user

You can download the pre-compiled binaries for the corresponding platform from [release page](https://github.com/Gozap/certmonitor/releases).
Next create a configuration file named `certmonitor.yaml`, like this:

``` yaml
alarm:
- type: smtp
  targets:
  - your_email_address
monitor:
  websites:
  - https://google.com
  cron: '@every 10s'
  beforetime: 168h0m0s
smtp:
  username: your_email_address
  password: "password"
  from: your_email_address
  server: "smtp_server:465"
```

Finally run it(Suppose the file you downloaded is named `certmonitor_linux_amd64`):

``` sh
chmod +x certmonitor_linux_amd64
./certmonitor_linux_amd64
```

### For docker user(Advanced)

build docker image

``` sh
export version=v1.0.1
make docker
```

create a config named `certmonitor.yaml`

``` yaml
alarm:
- type: smtp
  targets:
  - your_email_address
monitor:
  websites:
  - https://google.com
  cron: '@every 10s'
  beforetime: 168h0m0s
smtp:
  username: your_email_address
  password: "password"
  from: your_email_address
  server: "smtp_server:465"
```

run a container

``` sh
docker run -dt --name cermonitor -v ./certmonitor.yaml:/certmonitor.yaml gozap/certmonitor:v1.0.1
```
