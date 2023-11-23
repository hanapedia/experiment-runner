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

	// metrics query
	METRICS_QUERY_CONFIG_MAP_NAME     string
	METRICS_QUERY_IMAGE               string
	METRICS_QUERY_IMAGE_TAG           string
	METRICS_QUERY_ENDPOINT            string
	METRICS_QUERY_END_TIME            string
	METRICS_QUERY_DURATION            string
	METRICS_QUERY_STEP                string
	METRICS_QUERY_AWS_REGION          string
	METRICS_QUERY_S3_BUCKET           string
	METRICS_QUERY_WORKLOAD_CONTAINERS string

	// rca
	RCA_NORMAL_DURATION        string
	RCA_INJECTION_DURATION     string
	RCA_LATENCY                string
	RCA_JITTER                 string
	RCA_INJECTION_IGNORE_KEY   string
	RCA_INJECTION_IGNORE_VALUE string
}

var defaults = EnvVars{
	EXPERIMENT_NAME:      "test",
	TARGET_NAMESPACE:     "emulation",
	EXPERIMENT_NAMESPACE: "experiment",

	METRICS_QUERY_CONFIG_MAP_NAME:     "metrics-query-env",
	METRICS_QUERY_IMAGE:               "docker.io/hiroki11hanada/metrics-query",
	METRICS_QUERY_IMAGE_TAG:           "latest",
	METRICS_QUERY_ENDPOINT:            "http://localhost:9090",
	METRICS_QUERY_END_TIME:            "",
	METRICS_QUERY_DURATION:            "30m",
	METRICS_QUERY_STEP:                "15s",
	METRICS_QUERY_AWS_REGION:          "ap-northeast-1",
	METRICS_QUERY_S3_BUCKET:           "test",
	METRICS_QUERY_WORKLOAD_CONTAINERS: "server|redis",

	RCA_NORMAL_DURATION:        "",
	RCA_INJECTION_DURATION:     "",
	RCA_LATENCY:                "",
	RCA_JITTER:                 "",
	RCA_INJECTION_IGNORE_KEY:   "injection",
	RCA_INJECTION_IGNORE_VALUE: "1",
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

		METRICS_QUERY_ENDPOINT:            readEnv("METRICS_QUERY_ENDPOINT", defaults.METRICS_QUERY_ENDPOINT),
		METRICS_QUERY_END_TIME:            readEnv("END_TIME", defaults.METRICS_QUERY_END_TIME),
		METRICS_QUERY_DURATION:            readEnv("DURATION", defaults.METRICS_QUERY_DURATION),
		METRICS_QUERY_STEP:                readEnv("STEP", defaults.METRICS_QUERY_STEP),
		METRICS_QUERY_AWS_REGION:          readEnv("AWS_REGION", defaults.METRICS_QUERY_AWS_REGION),
		METRICS_QUERY_S3_BUCKET:           readEnv("S3_BUCKET", defaults.METRICS_QUERY_S3_BUCKET),
		METRICS_QUERY_WORKLOAD_CONTAINERS: readEnv("WORKLOAD_CONTAINERS", defaults.METRICS_QUERY_WORKLOAD_CONTAINERS),

		RCA_NORMAL_DURATION:    readEnv("RCA_NORMAL_DURATION", defaults.RCA_NORMAL_DURATION),
		RCA_INJECTION_DURATION: readEnv("RCA_INJECTION_DURATION", defaults.RCA_INJECTION_DURATION),
		RCA_LATENCY:            readEnv("RCA_LATENCY", defaults.RCA_LATENCY),
		RCA_JITTER:             readEnv("RCA_JITTER", defaults.RCA_JITTER),
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
