package kube

import (
	"errors"
	"fmt"
	"text/template"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const (
	category           = "kube"
	deployTemplateFile = "deployment.tpl"
	jobTemplateFile    = "job.tpl"
	basePort           = 30000
	portLimit          = 32767
)

// Deployment describes the k8s deployment yaml
type Deployment struct {
	Name        string
	Namespace   string
	Image       string
	Secret      string
	Replicas    int
	Revisions   int
	Port        int
	NodePort    int
	UseNodePort bool
	RequestCpu  int
	RequestMem  int
	LimitCpu    int
	LimitMem    int
	MinReplicas int
	MaxReplicas int
}

// DeploymentCommand is used to generate the kubernetes deployment yaml files.
func DeploymentCommand(c *cli.Context) error {
	nodePort := c.Int("nodePort")
	home := c.String("home")
	remote := c.String("remote")
	if len(remote) > 0 {
		repo, _ := util.CloneIntoGitHome(remote)
		if len(repo) > 0 {
			home = repo
		}
	}

	if len(home) > 0 {
		util.RegisterGoctlHome(home)
	}

	// 0 to disable the nodePort type
	if nodePort != 0 && (nodePort < basePort || nodePort > portLimit) {
		return errors.New("nodePort should be between 30000 and 32767")
	}

	text, err := util.LoadTemplate(category, deployTemplateFile, deploymentTemplate)
	if err != nil {
		return err
	}

	out, err := util.CreateIfNotExist(c.String("o"))
	if err != nil {
		return err
	}
	defer out.Close()

	t := template.Must(template.New("deploymentTemplate").Parse(text))
	err = t.Execute(out, Deployment{
		Name:        c.String("name"),
		Namespace:   c.String("namespace"),
		Image:       c.String("image"),
		Secret:      c.String("secret"),
		Replicas:    c.Int("replicas"),
		Revisions:   c.Int("revisions"),
		Port:        c.Int("port"),
		NodePort:    nodePort,
		UseNodePort: nodePort > 0,
		RequestCpu:  c.Int("requestCpu"),
		RequestMem:  c.Int("requestMem"),
		LimitCpu:    c.Int("limitCpu"),
		LimitMem:    c.Int("limitMem"),
		MinReplicas: c.Int("minReplicas"),
		MaxReplicas: c.Int("maxReplicas"),
	})
	if err != nil {
		return err
	}

	fmt.Println(aurora.Green("Done."))
	return nil
}

// Category returns the category of the deployments.
func Category() string {
	return category
}

// Clean cleans the generated deployment files.
func Clean() error {
	return util.Clean(category)
}

// GenTemplates generates the deployment template files.
func GenTemplates(_ *cli.Context) error {
	return util.InitTemplates(category, map[string]string{
		deployTemplateFile: deploymentTemplate,
		jobTemplateFile:    jobTmeplate,
	})
}

// RevertTemplate reverts the given template file to the default value.
func RevertTemplate(name string) error {
	return util.CreateTemplate(category, name, deploymentTemplate)
}

// Update updates the template files to the templates built in current goctl.
func Update() error {
	err := Clean()
	if err != nil {
		return err
	}

	return util.InitTemplates(category, map[string]string{
		deployTemplateFile: deploymentTemplate,
		jobTemplateFile:    jobTmeplate,
	})
}
