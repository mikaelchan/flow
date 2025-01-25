package serializer

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/mikaelchan/hamster/pkg/domain"
)

var factory *Factory

func init() {
	factory = NewFactory()
}

func GetFactory() *Factory {
	return factory
}

type Serializer func(domain.HasType) ([]byte, error)
type Deserializer func([]byte, domain.HasType) error
type Getter func() domain.HasType

func ReflectGetter(val domain.HasType) Getter {
	typ := reflect.TypeOf(val)
	return func() domain.HasType {
		x := reflect.New(typ.Elem()).Interface()
		return x.(domain.HasType)
	}
}

type Registry struct {
	serializer   Serializer
	deserializer Deserializer
	getter       Getter
}

func Register(val domain.HasType, serializer Serializer, deserializer Deserializer) {
	factory.Register(val, serializer, deserializer)
}

func Serialize(val domain.HasType) ([]byte, error) {
	return factory.Serialize(val)
}

func Deserialize(data []byte, val domain.HasType) error {
	return factory.Deserialize(data, val)
}

func DeserializeNew(typ domain.Type, data []byte) (domain.HasType, error) {
	return factory.DeserializeNew(typ, data)
}

type Factory struct {
	registries sync.Map
}

func NewFactory() *Factory {
	return &Factory{
		registries: sync.Map{},
	}
}

func (f *Factory) Register(val domain.HasType, serializer Serializer, deserializer Deserializer) {
	f.registries.Store(val.Type(), &Registry{
		serializer:   serializer,
		deserializer: deserializer,
		getter:       ReflectGetter(val),
	})
}

func (f *Factory) Serialize(val domain.HasType) ([]byte, error) {
	reg, err := f.getRegistry(val.Type())
	if err != nil {
		return nil, err
	}
	return reg.serializer(val)
}

func (f *Factory) Deserialize(data []byte, val domain.HasType) error {
	reg, err := f.getRegistry(val.Type())
	if err != nil {
		return err
	}
	return reg.deserializer(data, val)
}

func (f *Factory) DeserializeNew(typ domain.Type, data []byte) (domain.HasType, error) {
	reg, err := f.getRegistry(typ)
	if err != nil {
		return nil, err
	}

	val := reg.getter()
	if err := reg.deserializer(data, val); err != nil {
		return nil, fmt.Errorf("failed to deserialize data: %w", err)
	}
	return val, nil
}

// getRegistry safely retrieves and type asserts the registry
func (f *Factory) getRegistry(typ domain.Type) (*Registry, error) {
	value, ok := f.registries.Load(typ)
	if !ok {
		return nil, fmt.Errorf("no serializer registered for type: %s", typ.String())
	}

	reg, ok := value.(*Registry)
	if !ok {
		return nil, fmt.Errorf("invalid registry type for: %s", typ.String())
	}
	return reg, nil
}
