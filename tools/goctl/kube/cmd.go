package kube

import "github.com/spf13/cobra"

var (
	varStringName           string
	varStringNamespace      string
	varStringImage          string
	varStringSecret         string
	varIntRequestCpu        int
	varIntRequestMem        int
	varIntLimitCpu          int
	varIntLimitMem          int
	varStringO              string
	varIntReplicas          int
	varIntRevisions         int
	varIntPort              int
	varIntNodePort          int
	varIntMinReplicas       int
	varIntMaxReplicas       int
	varStringHome           string
	varStringRemote         string
	varStringBranch         string
	varStringServiceAccount string
	Cmd                     = &cobra.Command{
		Use:   "kube",
		Short: "generate kubernetes files",
		RunE:  deploymentCommand,
	}

	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "generate deployment yaml file",
	}
)

func init() {
	deployCmd.Flags().StringVar(&varStringName, "name", "", "the name of deployment (required)")
	deployCmd.Flags().StringVar(&varStringNamespace, "namespace", "", "the namespace of deployment (required)")
	deployCmd.Flags().StringVar(&varStringImage, "image", "", "the docker image of deployment (required)")
	deployCmd.Flags().StringVar(&varStringSecret, "secret", "", "the secret to image pull from registry")
	deployCmd.Flags().IntVar(&varIntRequestCpu, "requestCpu", 500, "the request cpu to deploy")
	deployCmd.Flags().IntVar(&varIntRequestMem, "requestMem", 512, "the request memory to deploy")
	deployCmd.Flags().IntVar(&varIntLimitCpu, "limitCpu", 1000, "the limit cpu to deploy")
	deployCmd.Flags().IntVar(&varIntLimitMem, "limitMem", 1024, "the limit memory to deploy")
	deployCmd.Flags().StringVar(&varStringO, "o", "", "the output yaml file (required)")
	deployCmd.Flags().IntVar(&varIntReplicas, "replicas", 3, "the number of replicas to deploy")
	deployCmd.Flags().IntVar(&varIntRevisions, "revisions", 5, "the number of revision history to limit")
	deployCmd.Flags().IntVar(&varIntPort, "port", 0, "the port of the deployment to listen on pod (required)")
	deployCmd.Flags().IntVar(&varIntNodePort, "nodePort", 0, "the nodePort of the deployment to expose")
	deployCmd.Flags().IntVar(&varIntMinReplicas, "minReplicas", 3, "the min replicas to deploy")
	deployCmd.Flags().IntVar(&varIntMaxReplicas, "maxReplicas", 10, "the max replicas to deploy")

	deployCmd.Flags().StringVar(&varStringHome, "home", "", "the goctl home path of the template, "+
		"--home and --remote cannot be set at the same time, if they are, --remote has higher priority")
	deployCmd.Flags().StringVar(&varStringRemote, "remote", "", "the remote git repo of the template, "+
		"--home and --remote cannot be set at the same time, if they are, --remote has higher priority\n\tThe git repo "+
		"directory must be consistent with the https://github.com/zeromicro/go-zero-template directory structure")
	deployCmd.Flags().StringVar(&varStringBranch, "branch", "", "the branch of the remote repo, it "+
		"does work with --remote")
	deployCmd.Flags().StringVar(&varStringServiceAccount, "serviceAccount", "", "the ServiceAccount "+
		"for the deployment")
	deployCmd.MarkFlagRequired("name")
	deployCmd.MarkFlagRequired("namespace")
	deployCmd.MarkFlagRequired("o")
	deployCmd.MarkFlagRequired("port")

	Cmd.AddCommand(deployCmd)
}
