package kubernetes

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JobArgs struct {
	Name            string
	TargetNamespace string
	ConfigMapName   string
	S3BucketDir     string
	K6TestName      string
	JobImageName    string
	Duration        string
}

// ConstructEnvFromConfigMap creates an EnvFromSource object from a config map.
func ConstructEnvFromConfigMap(name string) corev1.EnvFromSource {
	return corev1.EnvFromSource{
		ConfigMapRef: &corev1.ConfigMapEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: name,
			},
		},
	}
}

// ConstructEnvFromSecret creates an EnvVar object from a secret.
func ConstructEnvFromSecret(envName, secretName, key string) corev1.EnvVar {
	return corev1.EnvVar{
		Name: envName,
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: secretName,
				},
				Key: key,
			},
		},
	}
}

// ConstructEnvFromString creates an EnvVar object from a string.
func ConstructEnvFromString(envName, value string) corev1.EnvVar {
	return corev1.EnvVar{
		Name:  envName,
		Value: value,
	}
}

// ConstructContainer creates a Container object.
func ConstructContainer(name, imageName string, envFrom []corev1.EnvFromSource, env []corev1.EnvVar) corev1.Container {
	return corev1.Container{
		Name:    name,
		Image:   imageName,
		EnvFrom: envFrom,
		Env:     env,
	}
}

// ConstructJob creates a batchv1.Job object equivalent to the provided yaml manifest.
func ConstructJob(args JobArgs) *batchv1.Job {
	envFrom := []corev1.EnvFromSource{
		ConstructEnvFromConfigMap(args.ConfigMapName),
	}
	env := []corev1.EnvVar{
		ConstructEnvFromSecret("AWS_ACCESS_KEY_ID", "aws-credentials", "aws_access_key_id"),
		ConstructEnvFromSecret("AWS_SECRET_ACCESS_KEY", "aws-credentials", "aws_secret_access_key"),
		ConstructEnvFromString("S3_BUCKET_DIR", args.S3BucketDir),
		ConstructEnvFromString("K6_TEST_NAME", args.K6TestName),
		ConstructEnvFromString("DURATION", args.Duration),
		ConstructEnvFromString("NAMESPACE", args.TargetNamespace),
		ConstructEnvFromString("TZ", "Asia/Tokyo"),
	}
	container := ConstructContainer(args.Name, args.JobImageName, envFrom, env)
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: args.Name,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers:    []corev1.Container{container},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
}
