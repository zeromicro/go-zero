package quickstart

import (
	_ "embed"
	"os"
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
	zrpcWorkDir   string
)

type serviceImpl struct {
	starter func()
}

func (s serviceImpl) Start() {
	s.starter()
}

func (s serviceImpl) Stop() {}

func initRPCProto() error {
	zrpcWorkDir = filepath.Join(projectDir, "rpc")
	if err := pathx.MkdirIfNotExist(zrpcWorkDir); err != nil {
		return err
	}

	protoFilename := filepath.Join(zrpcWorkDir, protoName)
	rpcBytes := []byte(protocContent)
	return os.WriteFile(protoFilename, rpcBytes, 0o666)
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
	execCommand(zrpcWorkDir, arg)
	etcFile := filepath.Join(zrpcWorkDir, "etc", "greet.yaml")
	logx.Must(os.WriteFile(etcFile, []byte(rpcEtcContent), 0o666))
}

func (m micro) start() {
	mono := newMonoService(true)
	goModTidy(projectDir)
	sg := service.NewServiceGroup()
	sg.Add(serviceImpl{func() {
		log.Debug(">> Ready to start a zRPC server...")
		goStart(zrpcWorkDir)
	}})
	sg.Add(serviceImpl{func() {
		mono.start()
	}})
	sg.Start()
}
