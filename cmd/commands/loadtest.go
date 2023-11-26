package commands

import (
	"github.com/hanapedia/experiment-runner/internal/application/service"
	"github.com/hanapedia/experiment-runner/internal/infrastructure/config"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"

	k8sInfra "github.com/hanapedia/experiment-runner/internal/infrastructure/kubernetes"
)

var (
	isDryRun bool
)

// loadtestCmd represents the loadtest command
var loadtestCmd = &cobra.Command{
	Use:   "loadtest",
	Short: "Run Loadtest and process metrics",
	Run: func(cmd *cobra.Command, args []string) {
		experimentConfig := config.NewExperimentConfig(isDryRun)
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

		kubernetesAdapter := k8sInfra.NewKubernetesAdapter(kubeConfig)

		runner := service.NewLoadTestRunner(kubernetesAdapter, experimentConfig)
		err = runner.Run()
		if err != nil {
			panic(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(loadtestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadtestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	loadtestCmd.Flags().BoolVarP(&isDryRun, "dry-run", "d", false, "run command dry")
}
