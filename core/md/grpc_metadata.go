package md

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/metadata"
)

// GRPCMetadataCarrier represents that the data in the metadata of grpc is converted into Metadata.
type GRPCMetadataCarrier metadata.MD

func (g GRPCMetadataCarrier) Extract(ctx context.Context) (context.Context, error) {
	m, err := g.extract()
	if err != nil {
		return ctx, err
	}

	md := FromContext(ctx)
	md = md.Clone()
	for k, v := range m {
		md.Append(k, v...)
	}

	return NewContext(ctx, md), nil
}

func (g GRPCMetadataCarrier) extract() (Metadata, error) {
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
		return nil, fmt.Errorf("metadata:%w", err)
	}

	return m, nil
}

func (g GRPCMetadataCarrier) Inject(ctx context.Context) error {
	md := FromContext(ctx)
	if len(md) == 0 {
		return nil
	}

	bytes, err := json.Marshal(md)
	if err != nil {
		return fmt.Errorf("metadata:%w", err)
	}

	g["metadata"] = []string{string(bytes)}
	return nil
}
