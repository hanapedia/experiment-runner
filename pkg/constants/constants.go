package constants

const (
	// ImageName defines the image used by the application.
	ImageName = "docker.io/hiroki11hanada/rca-batch"

	// DefaultImageTag is the defaut value for image tags.
	DefaultImageTag = "latest"

	// ConfigMapName defines the name of the ConfigMap used by the application.
	DefaultConfigMapName = "rca-batch-env"

	// RcaBatchServiceAccountName defines the name of the ServiceAccount used by the application.
	RcaBatchServiceAccountName = "rca-batch-serviceaccount"

	// RcaNamespace defines the name of namespace where this job and the job created will run in.
	RcaNamespace = "rca"

	// ChaosExperimentNamespace defines the name of the namespace that chaos experiment resources are created in.
	ChaosExperimentNamespace = "chaos-mesh"

	// DeploymentIgnoreAnnotaionKey defines the key for annotation that indicates deployments to ignore.
	DeploymentIgnoreAnnotaionKey = "rca"

	// DeploymentIgnoreAnnotaionValue defines the value for annotation that indicates deployments to ignore.
	DeploymentIgnoreAnnotaionValue = "ignore"

	// DefaultLatency is the default latency to be used in experiment.
	DefaultLatency = "15ms"

	// DefaultJitter is the default jitter to be used in experiment.
	DefaultJitter = "5ms"
)
