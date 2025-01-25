package json

import (
	"encoding/json"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/serializer"
)

func RegisterJSON(val ...domain.HasType) {
	Register(serializer.GetFactory(), val...)
}

func Register(f *serializer.Factory, val ...domain.HasType) {
	for _, v := range val {
		f.Register(v, func(hasType domain.HasType) ([]byte, error) {
			return json.Marshal(hasType)
		}, func(data []byte, hasType domain.HasType) error {
			return json.Unmarshal(data, hasType)
		})
	}

}
