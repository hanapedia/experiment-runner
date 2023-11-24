package kubernetes

import "k8s.io/apimachinery/pkg/api/resource"

var LG_REPLICAS int32 = 1
var LG_CPU_LIMIT = resource.MustParse("250m")
var LG_CPU_REQUEST = resource.MustParse("125m")
var LG_MEMORY_LIMIT = resource.MustParse("1G")
var LG_MEMORY_REQUEST = resource.MustParse("64Mi")
