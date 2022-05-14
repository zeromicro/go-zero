package quickstart

import (
	_ "embed"
	"io/ioutil"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	//go:embed idl/greet.api
	apiContent string
	//go:embed idl/svc.tpl
	svcContent string
	//go:embed idl/apilogic.tpl
	apiLogicContent string
	//go:embed idl/api.yaml
	apiEtcContent string

	apiWorkDir string
)

func initAPIFlags() error {
	apiWorkDir = filepath.Join(projectDir, "api")
	if err := pathx.MkdirIfNotExist(apiWorkDir); err != nil {
		return err
	}

	apiFilename := filepath.Join(apiWorkDir, "greet.api")
	apiBytes := []byte(apiContent)
	if err := ioutil.WriteFile(apiFilename, apiBytes, 0666); err != nil {
		return err
	}

	gogen.VarStringDir = apiWorkDir
	gogen.VarStringAPI = apiFilename
	return nil
}

type mono struct {
	callRPC bool
}

func newMonoService(callRPC bool) mono {
	m := mono{callRPC}
	m.createAPIProject()
	return m
}

func (m mono) createAPIProject() {
	logx.Must(initAPIFlags())
	log.Debug(">> Generating quickstart api project...")
	logx.Must(gogen.GoCommand(nil, nil))
	etcFile := filepath.Join(apiWorkDir, "etc", "greet.yaml")
	logx.Must(ioutil.WriteFile(etcFile, []byte(apiEtcContent), 0666))
	logicFile := filepath.Join(apiWorkDir, "internal", "logic", "pinglogic.go")

	logx.Must(util.With("logic").Parse(apiLogicContent).SaveTo(map[string]bool{
		"callRPC": m.callRPC,
	}, logicFile, true))

	if m.callRPC {
		svcFile := filepath.Join(apiWorkDir, "internal", "svc", "servicecontext.go")
		logx.Must(ioutil.WriteFile(svcFile, []byte(svcContent), 0666))
	}
}

func (m mono) start() {
	if !m.callRPC {
		goModTidy(projectDir)
	}
	log.Debug(">> Ready to start an API server...")
	log.Debug(">> Run 'curl http://127.0.0.1:8888/ping' after service startup...")
	goStart(apiWorkDir)
}
