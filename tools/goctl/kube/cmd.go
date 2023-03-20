package kube

import (
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/internal/flags"
)

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
		Short: flags.Get("kube.short"),
	}

	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: flags.Get("kube.deploy.short"),
		RunE:  deploymentCommand,
	}
)

func init() {
	deployCmdFlags := deployCmd.Flags()
	deployCmdFlags.StringVar(&varStringName, "name", "", flags.Get("kube.deploy.name"))
	deployCmdFlags.StringVar(&varStringNamespace, "namespace", "", flags.Get("kube.deploy.namespace"))
	deployCmdFlags.StringVar(&varStringImage, "image", "", flags.Get("kube.deploy.image"))
	deployCmdFlags.StringVar(&varStringSecret, "secret", "", flags.Get("kube.deploy.secret"))
	deployCmdFlags.IntVar(&varIntRequestCpu, "requestCpu", 500, flags.Get("kube.deploy.requestCpu"))
	deployCmdFlags.IntVar(&varIntRequestMem, "requestMem", 512, flags.Get("kube.deploy.requestMem"))
	deployCmdFlags.IntVar(&varIntLimitCpu, "limitCpu", 1000, flags.Get("kube.deploy.limitCpu"))
	deployCmdFlags.IntVar(&varIntLimitMem, "limitMem", 1024, flags.Get("kube.deploy.limitMem"))
	deployCmdFlags.StringVar(&varStringO, "o", "", flags.Get("kube.deploy.o"))
	deployCmdFlags.IntVar(&varIntReplicas, "replicas", 3, flags.Get("kube.deploy.replicas"))
	deployCmdFlags.IntVar(&varIntRevisions, "revisions", 5, flags.Get("kube.deploy.revisions"))
	deployCmdFlags.IntVar(&varIntPort, "port", 0, flags.Get("kube.deploy.port"))
	deployCmdFlags.IntVar(&varIntNodePort, "nodePort", 0, flags.Get("kube.deploy.nodePort"))
	deployCmdFlags.IntVar(&varIntTargetPort, "targetPort", 0, flags.Get("kube.deploy.targetPort"))
	deployCmdFlags.IntVar(&varIntMinReplicas, "minReplicas", 3, flags.Get("kube.deploy.minReplicas"))
	deployCmdFlags.IntVar(&varIntMaxReplicas, "maxReplicas", 10, flags.Get("kube.deploy.maxReplicas"))
	deployCmdFlags.StringVar(&varStringImagePullPolicy, "imagePullPolicy", "", flags.Get("kube.deploy.imagePullPolicy"))
	deployCmdFlags.StringVar(&varStringHome, "home", "", flags.Get("kube.deploy.home"))
	deployCmdFlags.StringVar(&varStringRemote, "remote", "", flags.Get("kube.deploy.remote"))
	deployCmdFlags.StringVar(&varStringBranch, "branch", "", flags.Get("kube.deploy.branch"))
	deployCmdFlags.StringVar(&varStringServiceAccount, "serviceAccount", "", flags.Get("kube.deploy.serviceAccount"))

	_ = deployCmd.MarkFlagRequired("name")
	_ = deployCmd.MarkFlagRequired("namespace")
	_ = deployCmd.MarkFlagRequired("o")
	_ = deployCmd.MarkFlagRequired("port")

	Cmd.AddCommand(deployCmd)
}
