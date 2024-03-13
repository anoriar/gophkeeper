package uuid

//go:generate mockgen -source=uuid_generator_interface.go -destination=mock_uuid_generator/mock_uuid_generator.go -package=mock_uuid_generator
type UUIDGeneratorInterface interface {
	NewString() string
}
