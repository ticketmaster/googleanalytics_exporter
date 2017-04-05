# Google Real Time Analytics to Prometheus

Obtains Google Analytics RealTime metrics, and presents them to prometheus for scraping.

---

## Quick start

1. Install [Glide][1], if it isn't installed.
1. Copy your [Google creds][2] json file to ./config/ga_creds.json. Make sure the email from the json is added to the GA project permissions. More on this bellow.
1. Create yaml configuration file, see example bellow.
    ```yaml
    promport: 9100
    interval: 60
    viewid: ga:123456789
    metrics:
    - rt:pageviews
    - rt:activeUsers
    ```
1. Install dependencies, compile and run.
    ```bash
    glide install
    go build ganalytics.go
    ./ganalytics
    ```

### ViewID for the Google Analytics

From your Google Analytics Web UI: *Admin (Low left) ==> View Settings (far right tab, named VIEW)'*

*View ID* should be among *Basic Settings*, add prefix `ga:` to the ID listed in the Web UI, e.g. `ga:1234556`.

### Google creds

[Google API manager][2] allows to create OAuth 2.0 credentials for Google APIs. Use *Service account key* credentials type, upon creation a json creds file will be provided. Project RO permissions should be sufficient.

>*The email from GA API creds must be added to analytics project metrics will be obtained from.*>


### Cross compile on a MAC

* go should be installed with common compilers - `brew install go --with-cc-common`
* ensure creds.json and config.yaml in ./config/

```bash
CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo ganalytics.go
docker build -t ganalytics .
docker run -it -p 9100:9100 -v $(pwd)/config:/ga/config ganalytics
```


[1]: https://github.com/Masterminds/glide
[2]: https://console.developers.google.com/apis/credentials
