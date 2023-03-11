package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
)

type Method struct {
	HttpMethod string
	HttpPath   string
	RpcPath    string
}

// GetMethods returns all methods of the given grpcurl.DescriptorSource.
func GetMethods(source grpcurl.DescriptorSource) ([]Method, error) {
	svcs, err := source.ListServices()
	if err != nil {
		return nil, err
	}

	var methods []Method
	for _, svc := range svcs {
		d, err := source.FindSymbol(svc)
		if err != nil {
			return nil, err
		}

		switch val := d.(type) {
		case *desc.ServiceDescriptor:
			svcMethods := val.GetMethods()
			for _, method := range svcMethods {
				rpcPath := fmt.Sprintf("%s/%s", svc, method.GetName())
				ext := proto.GetExtension(method.GetMethodOptions(), annotations.E_Http)
				switch rule := ext.(type) {
				case *annotations.HttpRule:
					if rule == nil {
						methods = append(methods, Method{
							RpcPath: rpcPath,
						})
						continue
					}

					switch httpRule := rule.GetPattern().(type) {
					case *annotations.HttpRule_Get:
						methods = append(methods, Method{
							HttpMethod: http.MethodGet,
							HttpPath:   adjustHttpPath(httpRule.Get),
							RpcPath:    rpcPath,
						})
					case *annotations.HttpRule_Post:
						methods = append(methods, Method{
							HttpMethod: http.MethodPost,
							HttpPath:   adjustHttpPath(httpRule.Post),
							RpcPath:    rpcPath,
						})
					case *annotations.HttpRule_Put:
						methods = append(methods, Method{
							HttpMethod: http.MethodPut,
							HttpPath:   adjustHttpPath(httpRule.Put),
							RpcPath:    rpcPath,
						})
					case *annotations.HttpRule_Delete:
						methods = append(methods, Method{
							HttpMethod: http.MethodDelete,
							HttpPath:   adjustHttpPath(httpRule.Delete),
							RpcPath:    rpcPath,
						})
					case *annotations.HttpRule_Patch:
						methods = append(methods, Method{
							HttpMethod: http.MethodPatch,
							HttpPath:   adjustHttpPath(httpRule.Patch),
							RpcPath:    rpcPath,
						})
					default:
						methods = append(methods, Method{
							RpcPath: rpcPath,
						})
					}
				default:
					methods = append(methods, Method{
						RpcPath: rpcPath,
					})
				}
			}
		}
	}

	return methods, nil
}

func adjustHttpPath(path string) string {
	path = strings.ReplaceAll(path, "{", ":")
	path = strings.ReplaceAll(path, "}", "")
	return path
}
