package uuid

import "github.com/google/uuid"

type UUIDGenerator struct {
}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

func (gen *UUIDGenerator) NewString() string {
	return uuid.NewString()
}
