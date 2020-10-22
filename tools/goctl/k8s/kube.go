package k8s

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

const (
	ServiceTypeApi ServiceType = "api"
	ServiceTypeRpc ServiceType = "rpc"
	ServiceTypeJob ServiceType = "job"
	envDev                     = "dev"
)

var errUnknownServiceType = errors.New("unknown service type")

type (
	ServiceType string

	KubeRequest struct {
		Env                        string
		ServiceName                string
		ServiceType                ServiceType
		Namespace                  string
		Schedule                   string
		Replicas                   int
		RevisionHistoryLimit       int
		Port                       int
		LimitCpu                   int
		LimitMem                   int
		RequestCpu                 int
		RequestMem                 int
		SuccessfulJobsHistoryLimit int
		HpaMinReplicas             int
		HpaMaxReplicas             int
	}
)

func Gen(req KubeRequest) (string, error) {
	switch req.ServiceType {
	case ServiceTypeApi, ServiceTypeRpc:
		return genApiRpc(req)
	case ServiceTypeJob:
		return genJob(req)
	default:
		return "", errUnknownServiceType
	}
}

func genApiRpc(req KubeRequest) (string, error) {
	t, err := template.New("api_rpc").Parse(apiRpcTmeplate)
	if err != nil {
		return "", err
	}
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, map[string]interface{}{
		"name":                 fmt.Sprintf("%s-%s", req.ServiceName, req.ServiceType),
		"namespace":            req.Namespace,
		"replicas":             req.Replicas,
		"revisionHistoryLimit": req.RevisionHistoryLimit,
		"port":                 req.Port,
		"limitCpu":             req.LimitCpu,
		"limitMem":             req.LimitMem,
		"requestCpu":           req.RequestCpu,
		"requestMem":           req.RequestMem,
		"serviceName":          req.ServiceName,
		"env":                  req.Env,
		"envIsPreOrPro":        req.Env != envDev,
		"envIsDev":             req.Env == envDev,
		"minReplicas":          req.HpaMinReplicas,
		"maxReplicas":          req.HpaMaxReplicas,
	})
	if err != nil {
		return "", nil
	}
	return buffer.String(), nil
}

func genJob(req KubeRequest) (string, error) {
	t, err := template.New("job").Parse(jobTmeplate)
	if err != nil {
		return "", err
	}
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, map[string]interface{}{
		"name":                       fmt.Sprintf("%s-%s", req.ServiceName, req.ServiceType),
		"namespace":                  req.Namespace,
		"schedule":                   req.Schedule,
		"successfulJobsHistoryLimit": req.SuccessfulJobsHistoryLimit,
		"limitCpu":                   req.LimitCpu,
		"limitMem":                   req.LimitMem,
		"requestCpu":                 req.RequestCpu,
		"requestMem":                 req.RequestMem,
		"serviceName":                req.ServiceName,
		"env":                        req.Env,
	})
	if err != nil {
		return "", nil
	}
	return buffer.String(), nil
}
