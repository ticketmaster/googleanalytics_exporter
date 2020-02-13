[![Build Status](https://travis-ci.org/paha/googleanalytics_exporter.svg?branch=master)](https://travis-ci.org/paha/googleanalytics_exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/paha/googleanalytics_exporter)](https://goreportcard.com/report/github.com/paha/googleanalytics_exporter)
[![Docker Repository on Quay](https://quay.io/repository/paha/ga-prom/status "Docker Repository on Quay")](https://quay.io/repository/paha/ga-prom)

# Google Real Time Analytics to Prometheus

Obtains Google Analytics RealTime metrics, and presents them to prometheus for scraping.

---

## Quick start

1. Copy your [Google creds][2] json file to ./config/ga_creds.json. The email from the json must be added to the GA project permissions, more on that bellow. We recommend you use port 9674 to avoid conflicts as per Prometheus' [default port allocations](https://github.com/prometheus/prometheus/wiki/Default-port-allocations)
1. Create yaml configuration file (`./config/conf.yaml`):.
    ```yaml
    port: 9674
    interval: 60
    viewid: ga:123456789
    metrics:
    - rt:pageviews
    - rt:activeUsers
    ```
1. Install dependencies, compile and run.
    ```bash
    GO111MODULE=on go build ganalytics.go
    ./ganalytics
    ```

### ViewID for the Google Analytics

From your Google Analytics Web UI: *Admin (Low left) ==> View Settings (far right tab, named VIEW)'*

*View ID* should be among *Basic Settings*. Prefix `ga:` must be added to the ID, e.g. `ga:1234556` while adding it to the config.

### Google creds

[Google API manager][2] allows to create OAuth 2.0 credentials for Google APIs. Use *Service account key* credentials type, upon creation a json creds file will be provided. Project RO permissions should be sufficient.

>*The email from GA API creds must be added to analytics project metrics will be obtained from.*>


### Cross compile on a MAC

* [Alpine docker image][3] is used for delivery.
* go should be installed with common compilers - `brew install go --with-cc-common`
* `creds.json` and `config.yaml` expected to be in `./config/`

```bash
CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo ganalytics.go
docker build -t ganalytics .
docker run -it -p 9674:9674 -v $(pwd)/config:/ga/config ganalytics
```

## Author

Pavel Snagovsky, pavel@snagovsky.com

## License

Licensed under the terms of [MIT license][4], see [LICENSE][5] file

[1]: https://github.com/Masterminds/glide
[2]: https://console.developers.google.com/apis/credentials
[3]: https://hub.docker.com/_/alpine/
[4]: https://choosealicense.com/licenses/mit/
[5]: ./LICENSE
