package internal

import (
	"fmt"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/desc"
)

// GetMethods returns all methods of the given grpcurl.DescriptorSource.
func GetMethods(source grpcurl.DescriptorSource) ([]string, error) {
	svcs, err := source.ListServices()
	if err != nil {
		return nil, err
	}

	var methods []string
	for _, svc := range svcs {
		d, err := source.FindSymbol(svc)
		if err != nil {
			return nil, err
		}

		switch val := d.(type) {
		case *desc.ServiceDescriptor:
			svcMethods := val.GetMethods()
			for _, method := range svcMethods {
				methods = append(methods, fmt.Sprintf("%s/%s", svc, method.GetName()))
			}
		}
	}

	return methods, nil
}
