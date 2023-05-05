package kube

import (
	_ "embed"
	"errors"
	"fmt"
	"text/template"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const (
	category           = "kube"
	deployTemplateFile = "deployment.tpl"
	jobTemplateFile    = "job.tpl"
	basePort           = 30000
	portLimit          = 32767
)

var (
	//go:embed deployment.tpl
	deploymentTemplate string
	//go:embed job.tpl
	jobTemplate string
)

// Deployment describes the k8s deployment yaml
type Deployment struct {
	Name            string
	Namespace       string
	Image           string
	Secret          string
	Replicas        int
	Revisions       int
	Port            int
	TargetPort      int
	NodePort        int
	UseNodePort     bool
	RequestCpu      int
	RequestMem      int
	LimitCpu        int
	LimitMem        int
	MinReplicas     int
	MaxReplicas     int
	ServiceAccount  string
	ImagePullPolicy string
}

// deploymentCommand is used to generate the kubernetes deployment yaml files.
func deploymentCommand(_ *cobra.Command, _ []string) error {
	nodePort := varIntNodePort
	home := varStringHome
	remote := varStringRemote
	branch := varStringBranch
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote, branch)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		pathx.RegisterGoctlHome(home)
	}

	// 0 to disable the nodePort type
	if nodePort != 0 && (nodePort < basePort || nodePort > portLimit) {
		return errors.New("nodePort should be between 30000 and 32767")
	}

	text, err := pathx.LoadTemplate(category, deployTemplateFile, deploymentTemplate)
	if err != nil {
		return err
	}

	out, err := pathx.CreateIfNotExist(varStringO)
	if err != nil {
		return err
	}
	defer out.Close()

	if varIntTargetPort == 0 {
		varIntTargetPort = varIntPort
	}

	t := template.Must(template.New("deploymentTemplate").Parse(text))
	err = t.Execute(out, Deployment{
		Name:            varStringName,
		Namespace:       varStringNamespace,
		Image:           varStringImage,
		Secret:          varStringSecret,
		Replicas:        varIntReplicas,
		Revisions:       varIntRevisions,
		Port:            varIntPort,
		TargetPort:      varIntTargetPort,
		NodePort:        nodePort,
		UseNodePort:     nodePort > 0,
		RequestCpu:      varIntRequestCpu,
		RequestMem:      varIntRequestMem,
		LimitCpu:        varIntLimitCpu,
		LimitMem:        varIntLimitMem,
		MinReplicas:     varIntMinReplicas,
		MaxReplicas:     varIntMaxReplicas,
		ServiceAccount:  varStringServiceAccount,
		ImagePullPolicy: varStringImagePullPolicy,
	})
	if err != nil {
		return err
	}

	fmt.Println(color.Green.Render("Done."))
	return nil
}

// Category returns the category of the deployments.
func Category() string {
	return category
}

// Clean cleans the generated deployment files.
func Clean() error {
	return pathx.Clean(category)
}

// GenTemplates generates the deployment template files.
func GenTemplates() error {
	return pathx.InitTemplates(category, map[string]string{
		deployTemplateFile: deploymentTemplate,
		jobTemplateFile:    jobTemplate,
	})
}

// RevertTemplate reverts the given template file to the default value.
func RevertTemplate(name string) error {
	return pathx.CreateTemplate(category, name, deploymentTemplate)
}

// Update updates the template files to the templates built in current goctl.
func Update() error {
	err := Clean()
	if err != nil {
		return err
	}

	return pathx.InitTemplates(category, map[string]string{
		deployTemplateFile: deploymentTemplate,
		jobTemplateFile:    jobTemplate,
	})
}
