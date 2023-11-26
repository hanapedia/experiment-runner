package file

import (
	"os"

	"sigs.k8s.io/yaml"
)

func WriteKubernetesManifest(obj interface{}, path string) {
	file, err := CreateFile(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	manifest, err := yaml.Marshal(obj)
	_, err = file.WriteString(string(manifest))
	if err != nil {
		panic(err.Error())
	}
}

// CreateFile create and open file in append
func CreateFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
}
