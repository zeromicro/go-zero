package kube

import "github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"

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
	Cmd       = cobrax.NewCommand("kube")
	deployCmd = cobrax.NewCommand("deploy", cobrax.WithRunE(deploymentCommand))
)

func init() {
	deployCmdFlags := deployCmd.Flags()
	deployCmdFlags.StringVarP(&varStringName, "name", "n")
	deployCmdFlags.StringVarP(&varStringNamespace, "namespace", "a")
	deployCmdFlags.StringVarP(&varStringImage, "image", "i")
	deployCmdFlags.StringVarP(&varStringSecret, "secret", "s")
	deployCmdFlags.IntVarWithDefaultValue(&varIntRequestCpu, "requestCpu", 500)
	deployCmdFlags.IntVarWithDefaultValue(&varIntRequestMem, "requestMem", 512)
	deployCmdFlags.IntVarWithDefaultValue(&varIntLimitCpu, "limitCpu", 1000)
	deployCmdFlags.IntVarWithDefaultValue(&varIntLimitMem, "limitMem", 1024)
	deployCmdFlags.StringVar(&varStringO, "o")
	deployCmdFlags.IntVarWithDefaultValue(&varIntReplicas, "replicas", 3)
	deployCmdFlags.IntVarWithDefaultValue(&varIntRevisions, "revisions", 5)
	deployCmdFlags.IntVarPWithDefaultValue(&varIntPort, "port", "p", 9100)
	deployCmdFlags.IntVar(&varIntNodePort, "nodePort")
	deployCmdFlags.IntVar(&varIntTargetPort, "targetPort")
	deployCmdFlags.IntVarWithDefaultValue(&varIntMinReplicas, "minReplicas", 3)
	deployCmdFlags.IntVarWithDefaultValue(&varIntMaxReplicas, "maxReplicas", 10)
	deployCmdFlags.StringVar(&varStringImagePullPolicy, "imagePullPolicy")
	deployCmdFlags.StringVar(&varStringHome, "home")
	deployCmdFlags.StringVar(&varStringRemote, "remote")
	deployCmdFlags.StringVar(&varStringBranch, "branch")
	deployCmdFlags.StringVar(&varStringServiceAccount, "serviceAccount")

	_ = deployCmd.MarkFlagRequired("name")
	_ = deployCmd.MarkFlagRequired("namespace")
	_ = deployCmd.MarkFlagRequired("o")
	_ = deployCmd.MarkFlagRequired("port")

	Cmd.AddCommand(deployCmd)
}
