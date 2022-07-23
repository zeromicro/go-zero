package md

import (
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/metadata"
)

var _ Carrier = (*GRPCMetadataCarrier)(nil)

// GRPCMetadataCarrier represents that the data in the metadata of grpc is converted into Metadata.
type GRPCMetadataCarrier metadata.MD

func (g GRPCMetadataCarrier) Carrier() (Metadata, error) {
	m := Metadata{}
	md, ok := g["metadata"]
	if !ok {
		return m, nil
	}

	if len(md) == 0 {
		return m, nil
	}

	err := json.Unmarshal([]byte(md[0]), &m)
	if err != nil {
		return nil, fmt.Errorf("metadata deserialization error:%w", err)
	}

	return m, nil
}
