package cmd

import (
	"fmt"

	"github.com/plunder-app/kube-vip/pkg/kubevip"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// manifests will eventually deprecate the kubeadm set of subcommands
// manifests will be used to generate:
// - Pod spec manifest, mainly used for a static pod (kubeadm)
// - Daemonset manifest, mainly used to run kube-vip as a deamonset within Kubernetes (k3s/rke)

//var inCluster bool
var taint bool

func init() {
	kubeManifest.PersistentFlags().BoolVar(&inCluster, "inCluster", false, "Use the incluster token to authenticate to Kubernetes")
	kubeManifestDaemon.PersistentFlags().BoolVar(&taint, "taint", false, "Taint the manifest for only running on control planes")

	kubeManifest.AddCommand(kubeManifestPod)
	kubeManifest.AddCommand(kubeManifestDaemon)
}

var kubeManifest = &cobra.Command{
	Use:   "manifest",
	Short: "Manifest functions",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		// TODO - A load of text detailing what's actually happening
	},
}

var kubeManifestPod = &cobra.Command{
	Use:   "pod",
	Short: "Generate a Pod Manifest",
	Run: func(cmd *cobra.Command, args []string) {
		// Set the logging level for all subsequent functions
		log.SetLevel(log.Level(logLevel))
		initConfig.LoadBalancers = append(initConfig.LoadBalancers, initLoadBalancer)
		// TODO - A load of text detailing what's actually happening
		kubevip.ParseEnvironment(&initConfig)
		// TODO - check for certain things VIP/interfaces
		if initConfig.Interface == "" {
			cmd.Help()
			log.Fatalln("No interface is specified for kube-vip to bind to")
		}

		// The control plane has a requirement for a VIP being specified
		if initConfig.EnableControlPane && (initConfig.VIP == "" && initConfig.Address == "" && initConfig.DDNS == false) {
			cmd.Help()
			log.Fatalln("No address is specified for kube-vip to expose services on")
		}
		cfg := kubevip.GeneratePodManifestFromConfig(&initConfig, Release.Version, inCluster)

		fmt.Println(cfg)
	},
}

var kubeManifestDaemon = &cobra.Command{
	Use:   "daemonset",
	Short: "Generate a Daemonset Manifest",
	Run: func(cmd *cobra.Command, args []string) {
		// Set the logging level for all subsequent functions
		log.SetLevel(log.Level(logLevel))
		initConfig.LoadBalancers = append(initConfig.LoadBalancers, initLoadBalancer)
		// TODO - A load of text detailing what's actually happening
		kubevip.ParseEnvironment(&initConfig)
		// TODO - check for certain things VIP/interfaces
		if initConfig.Interface == "" {
			cmd.Help()
			log.Fatalln("No interface is specified for kube-vip to bind to")
		}

		// The control plane has a requirement for a VIP being specified
		if initConfig.EnableControlPane && (initConfig.VIP == "" && initConfig.Address == "" && initConfig.DDNS == false) {
			cmd.Help()
			log.Fatalln("No address is specified for kube-vip to expose services on")
		}
		cfg := kubevip.GenerateDeamonsetManifestFromConfig(&initConfig, Release.Version, inCluster, taint)

		fmt.Println(cfg)
	},
}
