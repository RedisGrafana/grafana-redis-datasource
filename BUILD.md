# How to install and build Grafana with Redis Datasource on RPM-based Linux

## Clone repository

```bash
git clone https://github.com/RedisTimeSeries/grafana-redis-datasource.git
```

## Install Grafana

- Follow [Install on RPM-based Linux](https://grafana.com/docs/grafana/latest/installation/rpm/) to install and start Grafana

- Open Grafana in web-browser `http://X.X.X.X:3000`

## Build Datasource

Redis datasource consists of both frontend and backend components.

### React Frontend

- Install latest version of Node.js using [Node Version Manager](https://github.com/nvm-sh/nvm)

- Install `yarn` to build Datasource

```bash
npm install yarn -g
```

- Install Datasource dependencies

```bash
yarn install
```

- Build Datasource

```bash
yarn build
```

### Golang Backend

- Install Golang

```bash
yum install go
```

- Install [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency

```bash
go get -u github.com/grafana/grafana-plugin-sdk-go
```

- Install mage (make-like build tool using Go)

```bash
git clone https://github.com/magefile/mage
cd mage
go run bootstrap.go
```

- Build backend plugin binaries for Linux, Windows and MacOS

```bash
mage -v
```

## Update Grafana Configuration

- Move distribution to Grafana's `plugins/` folder

```bash
mv dist/ /var/lib/grafana/plugins/redis-datasource
```

- Add `redis-datasource` to allowed unsigned plugins

```bash
vi /etc/grafana/grafana.ini
```

```
[plugins]
;enable_alpha = false
;app_tls_skip_verify_insecure = false
# Enter a comma-separated list of plugin identifiers to identify plugins that are allowed to be loaded even if they lack a valid signature.
allow_loading_unsigned_plugins = redis-datasource
```

- Verify that plugin registered

```bash
tail -100 /var/log/grafana/grafana.log
```

```
t=2020-07-01T06:03:38+0000 lvl=info msg="Starting plugin search" logger=plugins
t=2020-07-01T06:03:38+0000 lvl=warn msg="Running an unsigned backend plugin" logger=plugins pluginID=redis-datasource pluginDir=/var/lib/grafana/plugins/redis-datasource
t=2020-07-01T06:03:38+0000 lvl=info msg="Registering plugin" logger=plugins name=redis-datasource
t=2020-07-01T06:03:38+0000 lvl=info msg="HTTP Server Listen" logger=http.server address=[::]:3000 protocol=http subUrl= socket=
```

- Add new Datasource to Grafana using `Configuration` -> `Data Sources`

![Datasource](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/images/datasource.png)

- Import Redis monitoring dashboard from `dashbords/` folder

![Dashboard](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/images/redis-dashboard.png)

If you have questions, enhancement ideas or running into issues, please just open an issue on the repository: https://github.com/RedisTimeSeries/grafana-redis-datasource
