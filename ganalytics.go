/*
Obtains Google Analytics RealTime API metrics, and presents them to
prometheus for scraping.
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/analytics/v3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	credsfile = "./config/ga_creds.json"
	conffile  = "./config/conf.yaml"
	promGauge = make(map[string]*prometheus.GaugeVec)
	config    = new(conf)
)

// conf defines configuration parameters
type conf struct {
	Interval int          `yaml:"interval"`
	Metrics  []metricConf `yaml:"metrics"`
	ViewID   string       `yaml:"viewid"`
	PromPort string       `yaml:"port"`
}

type metricConf struct {
	Name       string   `yaml:"name"`
	Dimensions []string `yaml:"dimensions"`
	Filters    []string `yaml:"filters"`
	Sort       []string `yaml:"sort"`
	Limit      string   `yaml:"limit"`
}

func init() {
	config.getConf(conffile)

	// All metrics are registered as Prometheus Gauge with the possibility of dimensions (labels)
	for _, metric := range config.Metrics {
		promMetricName := getMetricName(metric)
		promMetricLabelName := getMetricLabels(metric)
		if promGauge[promMetricName] == nil {
			promGauge[promMetricName] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Name: fmt.Sprintf("ga_%s", promMetricName),
				Help: fmt.Sprintf("Google Analytics %s", metric.Name),
			}, promMetricLabelName)

			prometheus.MustRegister(promGauge[promMetricName])
		}
	}
}

func main() {
	creds := getCreds(credsfile)

	// JSON web token configuration
	jwtc := jwt.Config{
		Email:        creds["client_email"],
		PrivateKey:   []byte(creds["private_key"]),
		PrivateKeyID: creds["private_key_id"],
		Scopes:       []string{analytics.AnalyticsReadonlyScope},
		TokenURL:     creds["token_uri"],
		// Expires:      time.Duration(1) * time.Hour, // Expire in 1 hour
	}

	httpClient := jwtc.Client(oauth2.NoContext)
	as, err := analytics.New(httpClient)
	if err != nil {
		panic(err)
	}

	// Authenticated RealTime Google Analytics API service
	rts := analytics.NewDataRealtimeService(as)

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Google Analytics Exporter</title></head>
			<body>
			<h1>Google Analytics Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
	go http.ListenAndServe(fmt.Sprintf(":%s", config.PromPort), nil)

	for {
		for _, metric := range config.Metrics {
			// Go routine per metric dimension
			go func(metric metricConf) {
				metricName := getMetricName(metric)
				val := getMetric(rts, metric)
				for _, vald := range val {
					valf, _ := strconv.ParseFloat(vald[len(vald)-1], 64)
					promGauge[metricName].WithLabelValues(vald[0:len(metric.Dimensions)]...).Set(valf)
				}
			}(metric)
		}

		time.Sleep(time.Second * time.Duration(config.Interval))
	}
}

func getMetric(rts *analytics.DataRealtimeService, metric metricConf) [][]string {
	getc := rts.Get(config.ViewID, metric.Name)
	if metric.Dimensions != nil || len(metric.Dimensions) > 0 {
		getc = getc.Dimensions(strings.Join(metric.Dimensions, ","))
	}
	if metric.Filters != nil || len(metric.Filters) > 0 {
		getc = getc.Filters(strings.Join(metric.Filters, ","))
	}
	if metric.Sort != nil || len(metric.Sort) > 0 {
		getc = getc.Sort(strings.Join(metric.Sort, ","))
	}
	if len(metric.Limit) > 0 {
		i, err := strconv.ParseInt(metric.Limit, 10, 64)
		if err == nil {
			getc = getc.MaxResults(i)
		}
	}
	m, err := getc.Do()
	if err != nil {
		panic(err)
	}

	return m.Rows
}

// conf.getConf reads yaml configuration file
func (c *conf) getConf(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(data, &c); err != nil {
		panic(err)
	}
}

// https://console.developers.google.com/apis/credentials
// 'Service account keys' creds formated file is expected.
// NOTE: the email from the creds has to be added to the Analytics permissions
func getCreds(filename string) (r map[string]string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(data, &r); err != nil {
		panic(err)
	}

	return r
}

func getMetricName(dimension metricConf) string {
	promMetricName := append([]string{dimension.Name}, dimension.Dimensions...)
	return strings.Replace(strings.Join(promMetricName, "_"), ":", "_", -1)
}
func getMetricLabels(dimension metricConf) []string {
	result := make([]string, len(dimension.Dimensions))
	for i, item := range dimension.Dimensions {
		result[i] = strings.Replace(item, ":", "_", -1)
	}
	return result
}
