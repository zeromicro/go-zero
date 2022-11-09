package kube

import "github.com/spf13/cobra"

var (
	varStringName            string
	varStringNamespace       string
	varStringImage           string
	varStringSecret          string
	varIntRequestCpu         int
	varIntRequestMem         int
	varIntLimitCpu           int
	varIntLimitMem           int
	varStringO               string
	varIntReplicas           int
	varIntRevisions          int
	varIntPort               int
	varIntNodePort           int
	varIntTargetPort         int
	varIntMinReplicas        int
	varIntMaxReplicas        int
	varStringHome            string
	varStringRemote          string
	varStringBranch          string
	varStringServiceAccount  string
	varStringImagePullPolicy string

	// Cmd describes a kube command.
	Cmd = &cobra.Command{
		Use:   "kube",
		Short: "Generate kubernetes files",
	}

	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Generate deployment yaml file",
		RunE:  deploymentCommand,
	}
)

func init() {
	deployCmd.Flags().StringVar(&varStringName, "name", "", "The name of deployment (required)")
	deployCmd.Flags().StringVar(&varStringNamespace, "namespace", "", "The namespace of deployment (required)")
	deployCmd.Flags().StringVar(&varStringImage, "image", "", "The docker image of deployment (required)")
	deployCmd.Flags().StringVar(&varStringSecret, "secret", "", "The secret to image pull from registry")
	deployCmd.Flags().IntVar(&varIntRequestCpu, "requestCpu", 500, "The request cpu to deploy")
	deployCmd.Flags().IntVar(&varIntRequestMem, "requestMem", 512, "The request memory to deploy")
	deployCmd.Flags().IntVar(&varIntLimitCpu, "limitCpu", 1000, "The limit cpu to deploy")
	deployCmd.Flags().IntVar(&varIntLimitMem, "limitMem", 1024, "The limit memory to deploy")
	deployCmd.Flags().StringVar(&varStringO, "o", "", "The output yaml file (required)")
	deployCmd.Flags().IntVar(&varIntReplicas, "replicas", 3, "The number of replicas to deploy")
	deployCmd.Flags().IntVar(&varIntRevisions, "revisions", 5, "The number of revision history to limit")
	deployCmd.Flags().IntVar(&varIntPort, "port", 0, "The port of the deployment to listen on pod (required)")
	deployCmd.Flags().IntVar(&varIntNodePort, "nodePort", 0, "The nodePort of the deployment to expose")
	deployCmd.Flags().IntVar(&varIntTargetPort, "targetPort", 0, "The targetPort of the deployment, default to port")
	deployCmd.Flags().IntVar(&varIntMinReplicas, "minReplicas", 3, "The min replicas to deploy")
	deployCmd.Flags().IntVar(&varIntMaxReplicas, "maxReplicas", 10, "The max replicas to deploy")
	deployCmd.Flags().StringVar(&varStringImagePullPolicy, "imagePullPolicy", "", "Image pull policy. One of Always, Never, IfNotPresent")

	deployCmd.Flags().StringVar(&varStringHome, "home", "", "The goctl home path of the template, "+
		"--home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	deployCmd.Flags().StringVar(&varStringRemote, "remote", "", "The remote git repo of the template, "+
		"--home and --remote cannot be set at the same time, if they are, --remote has higher priority\nThe git repo "+
		"directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	deployCmd.Flags().StringVar(&varStringBranch, "branch", "", "The branch of the remote repo, it "+
		"does work with --remote")
	deployCmd.Flags().StringVar(&varStringServiceAccount, "serviceAccount", "", "The ServiceAccount "+
		"for the deployment")
	deployCmd.MarkFlagRequired("name")
	deployCmd.MarkFlagRequired("namespace")
	deployCmd.MarkFlagRequired("o")
	deployCmd.MarkFlagRequired("port")

	Cmd.AddCommand(deployCmd)
}
