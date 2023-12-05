package config

import (
	"os"
	"strconv"
	"sync"
)

type EnvVars struct {
	// common
	EXPERIMENT_NAME      string
	TARGET_NAMESPACE     string
	EXPERIMENT_NAMESPACE string
	K6_TEST_NAME         string
	DURATION             string
	ARRIVAL_RATES        string

	// metrics processor
	METRICS_PROCESSOR_CONFIG_MAP_NAME string
	METRICS_PROCESSOR_IMAGE           string
	METRICS_PROCESSOR_IMAGE_TAG       string
	METRICS_PROCESSOR_S3_BUCKET_DIR   string

	// rca
	RCA_NORMAL_DURATION        string
	RCA_INJECTION_DURATION     string
	RCA_LATENCY                string
	RCA_JITTER                 string
	RCA_INJECTION_IGNORE_KEY   string
	RCA_INJECTION_IGNORE_VALUE string

	// loadgenerator
	LG_CONFIG_MAP_NAME              string
	LG_IMAGE                        string
	LG_IMAGE_TAG                    string
	LG_FRONTEND_ADDR                string
	LG_INDEX_ROUTE                  string
	LG_K6_PROMETHEUS_RW_SERVER_URL  string
	LG_K6_PROMETHEUS_RW_TREND_STATS string
}

var defaults = EnvVars{
	EXPERIMENT_NAME:      "test",
	TARGET_NAMESPACE:     "emulation",
	EXPERIMENT_NAMESPACE: "experiment",
	K6_TEST_NAME:         "test",
	DURATION:             "30s",
	ARRIVAL_RATES:        "10,50,100,500,1000",

	METRICS_PROCESSOR_CONFIG_MAP_NAME: "metrics-processor-env",
	METRICS_PROCESSOR_IMAGE:           "docker.io/hiroki11hanada/metrics-processor",
	METRICS_PROCESSOR_IMAGE_TAG:       "latest",
	METRICS_PROCESSOR_S3_BUCKET_DIR:   "test",

	RCA_NORMAL_DURATION:        "",
	RCA_INJECTION_DURATION:     "",
	RCA_LATENCY:                "",
	RCA_JITTER:                 "",
	RCA_INJECTION_IGNORE_KEY:   "rca",
	RCA_INJECTION_IGNORE_VALUE: "ignore",

	LG_CONFIG_MAP_NAME:              "lg-script",
	LG_IMAGE:                        "grafana/k6",
	LG_IMAGE_TAG:                    "0.47.0",
	LG_FRONTEND_ADDR:                "frontend.emulation.svc.cluster.local:80",
	LG_INDEX_ROUTE:                  "/",
	LG_K6_PROMETHEUS_RW_SERVER_URL:  "http://prometheus-kube-prometheus-prometheus.monitoring.svc.cluster.local:9090/api/v1/write",
	LG_K6_PROMETHEUS_RW_TREND_STATS: "p(95),p(99),avg",
}

var envVars *EnvVars
var once sync.Once

func GetEnvs() *EnvVars {
	once.Do(func() {
		envVars = loadEnvVariables()
	})
	return envVars
}

func loadEnvVariables() *EnvVars {
	return &EnvVars{
		EXPERIMENT_NAME:      readEnv("EXPERIMENT_NAME", defaults.EXPERIMENT_NAME),
		TARGET_NAMESPACE:     readEnv("TARGET_NAMESPACE", defaults.TARGET_NAMESPACE),
		EXPERIMENT_NAMESPACE: readEnv("EXPERIMENT_NAMESPACE", defaults.EXPERIMENT_NAMESPACE),
		K6_TEST_NAME:         readEnv("K6_TEST_NAME", defaults.K6_TEST_NAME),
		DURATION:             readEnv("DURATION", defaults.DURATION),
		ARRIVAL_RATES:        readEnv("ARRIVAL_RATES", defaults.ARRIVAL_RATES),

		METRICS_PROCESSOR_CONFIG_MAP_NAME: readEnv("METRICS_PROCESSOR_CONFIG_MAP_NAME", defaults.METRICS_PROCESSOR_CONFIG_MAP_NAME),
		METRICS_PROCESSOR_IMAGE:           readEnv("METRICS_PROCESSOR_IMAGE", defaults.METRICS_PROCESSOR_IMAGE),
		METRICS_PROCESSOR_IMAGE_TAG:       readEnv("METRICS_PROCESSOR_IMAGE_TAG", defaults.METRICS_PROCESSOR_IMAGE_TAG),
		METRICS_PROCESSOR_S3_BUCKET_DIR:   readEnv("METRICS_PROCESSOR_S3_BUCKET_DIR", defaults.METRICS_PROCESSOR_S3_BUCKET_DIR),

		RCA_NORMAL_DURATION:        readEnv("RCA_NORMAL_DURATION", defaults.RCA_NORMAL_DURATION),
		RCA_INJECTION_DURATION:     readEnv("RCA_INJECTION_DURATION", defaults.RCA_INJECTION_DURATION),
		RCA_LATENCY:                readEnv("RCA_LATENCY", defaults.RCA_LATENCY),
		RCA_JITTER:                 readEnv("RCA_JITTER", defaults.RCA_JITTER),
		RCA_INJECTION_IGNORE_KEY:   readEnv("RCA_INJECTION_IGNORE_KEY", defaults.RCA_INJECTION_IGNORE_KEY),
		RCA_INJECTION_IGNORE_VALUE: readEnv("RCA_INJECTION_IGNORE_VALUE", defaults.RCA_INJECTION_IGNORE_VALUE),

		LG_CONFIG_MAP_NAME:              readEnv("LG_CONFIG_MAP_NAME", defaults.LG_CONFIG_MAP_NAME),
		LG_IMAGE:                        readEnv("LG_IMAGE", defaults.LG_IMAGE),
		LG_IMAGE_TAG:                    readEnv("LG_IMAGE_TAG", defaults.LG_IMAGE_TAG),
		LG_FRONTEND_ADDR:                readEnv("LG_FRONTEND_ADDR", defaults.LG_FRONTEND_ADDR),
		LG_INDEX_ROUTE:                  readEnv("LG_INDEX_ROUTE", defaults.LG_INDEX_ROUTE),
		LG_K6_PROMETHEUS_RW_SERVER_URL:  readEnv("LG_K6_PROMETHEUS_RW_SERVER_URL", defaults.LG_K6_PROMETHEUS_RW_SERVER_URL),
		LG_K6_PROMETHEUS_RW_TREND_STATS: readEnv("LG_K6_PROMETHEUS_RW_TREND_STATS", defaults.LG_K6_PROMETHEUS_RW_TREND_STATS),
	}
}

func readEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func readBoolEnv(key string, defaultValue bool) bool {
	boolValue := defaultValue
	if value, ok := os.LookupEnv(key); ok {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return boolValue
		}
		return parsed
	}
	return boolValue
}
