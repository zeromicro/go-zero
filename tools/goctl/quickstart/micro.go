package quickstart

import (
	_ "embed"
	"io/ioutil"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const protoName = "greet.proto"

var (
	//go:embed idl/greet.proto
	protocContent string
	//go:embed idl/rpc.yaml
	rpcEtcContent string

	zRPCWorkDir string
)

type serviceImpl struct {
	starter func()
}

func (s serviceImpl) Start() {
	s.starter()
}

func (s serviceImpl) Stop() {}

func initRPCProto() error {
	zRPCWorkDir = filepath.Join(projectDir, "rpc")
	if err := pathx.MkdirIfNotExist(zRPCWorkDir); err != nil {
		return err
	}

	protoFilename := filepath.Join(zRPCWorkDir, protoName)
	rpcBytes := []byte(protocContent)
	return ioutil.WriteFile(protoFilename, rpcBytes, 0666)
}

type micro struct{}

func newMicroService() micro {
	m := micro{}
	m.mustStartRPCProject()
	return m
}

func (m micro) mustStartRPCProject() {
	logx.Must(initRPCProto())
	log.Debug(">> Generating quickstart zRPC project...")
	arg := "goctl rpc protoc " + protoName + " --go_out=. --go-grpc_out=. --zrpc_out=. --verbose"
	execCommand(zRPCWorkDir, arg)
	etcFile := filepath.Join(zRPCWorkDir, "etc", "greet.yaml")
	logx.Must(ioutil.WriteFile(etcFile, []byte(rpcEtcContent), 0666))
}

func (m micro) start() {
	mono := newMonoService(true)
	goModTidy(projectDir)
	sg := service.NewServiceGroup()
	sg.Add(serviceImpl{func() {
		log.Debug(">> Ready to start a zRPC server...")
		goStart(zRPCWorkDir)
	}})
	sg.Add(serviceImpl{func() {
		mono.start()
	}})
	sg.Start()
}
