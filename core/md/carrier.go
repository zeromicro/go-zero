package md

import "context"

// Extractor represents an ability to extract Metadata.
type Extractor interface {
	Extract(ctx context.Context) (context.Context, error)
}

// Injector represents an ability that can be injected into Metadata.
type Injector interface {
	Inject(ctx context.Context) error
}
