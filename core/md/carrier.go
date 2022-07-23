package md

// Carrier represents an ability to convert other data into Metadata.
type Carrier interface {
	Carrier() (Metadata, error)
}
