package parser

import "fmt"

type PbMessage struct {
	Name      string
	Package   string
	PbPackage string
	GoPackage string
}

// Proto describes a proto file,
type Proto struct {
	Src              string
	Name             string
	Package          Package
	PbPackage        string
	GoPackage        string
	Import           []Import
	Message          []Message
	Service          Services
	ImportMessageMap map[string]PbMessage
}

func (p *Proto) generateImportMessageMap() {
	if p.ImportMessageMap == nil {
		p.ImportMessageMap = make(map[string]PbMessage)
	}
	for _, message := range p.Message {
		p.ImportMessageMap[message.Name] = PbMessage{
			Name:      CamelCase(message.Name),
			Package:   p.Package.Package.Name,
			PbPackage: p.PbPackage,
			GoPackage: p.GoPackage,
		}
	}
	for _, importFile := range p.Import {
		for _, message := range importFile.Proto.Message {
			key := fmt.Sprintf("%s.%s", importFile.Proto.Package.Package.Name, CamelCase(message.Name))
			p.ImportMessageMap[key] = PbMessage{
				Name:      CamelCase(message.Name),
				Package:   importFile.Proto.Package.Package.Name,
				PbPackage: importFile.Proto.PbPackage,
				GoPackage: importFile.Proto.GoPackage,
			}
		}
	}
}

func (p *Proto) GetImportMessage(key string) (msg PbMessage, existed bool) {
	msg, existed = p.ImportMessageMap[key]
	return msg, existed
}

func (p *Proto) HasGrpcService() (hasGrpcService bool) {
	if p == nil {
		return hasGrpcService
	}
	for _, service := range p.Service {
		if len(service.RPC) > 0 {
			hasGrpcService = true
			break
		}
	}

	return hasGrpcService
}
