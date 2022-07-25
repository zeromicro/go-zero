package md

import "context"

// Carrier represents an ability to convert other data into Metadata.
type Carrier interface {
	Extractor
	Injector
}

// Extractor represents an ability to extract Metadata.
type Extractor interface {
	Extract(ctx context.Context) (context.Context, error)
}

// Injector represents an ability that can be injected into Metadata.
type Injector interface {
	Inject(ctx context.Context) error
}
