package commands

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/dynamic"
	k8sClientGo "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/hanapedia/experiment-runner/internal/application/service"
	"github.com/hanapedia/experiment-runner/internal/infrastructure/chaosmesh"
	"github.com/hanapedia/experiment-runner/internal/infrastructure/config"
	k8sInfra "github.com/hanapedia/experiment-runner/internal/infrastructure/kubernetes"
)

// rcaCmd represents the rca command
var rcaCmd = &cobra.Command{
	Use:   "rca",
	Short: "Run a RCA experiment.",
	Run: func(cmd *cobra.Command, args []string) {
		// Prepare experiment configs
		config := config.NewRCAExperimentConfig()

		// Load kubeconfig from KUBECONFIG
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}
		clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			loadingRules,
			configOverrides,
		)
		kubeConfig, err := clientConfig.ClientConfig()
		if err != nil {
			panic(err.Error())
		}

		// prepare kube client for kubernetes API
		clientset, err := k8sClientGo.NewForConfig(kubeConfig)
		if err != nil {
			panic(err.Error())
		}
		kubernetesAdapter := k8sInfra.NewKubernetesAdapter(clientset, config)

		// Prepare kube dynamic config for chaos mesh resource
		dynamicClient, err := dynamic.NewForConfig(kubeConfig)
		if err != nil {
			panic(err.Error())
		}
		chaosAdapter := chaosmesh.NewChaosMeshAdapter(dynamicClient, config)

		experimentRunner := service.NewExperimentRunner(config, kubernetesAdapter, chaosAdapter)
		err = experimentRunner.Run()
		if err != nil {
			panic(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(rcaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rcaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rcaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
