package json_test

import (
	"testing"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/serializer"
	"github.com/mikaelchan/hamster/pkg/serializer/json"
)

// MockData is a test struct that implements Serializable
type MockData struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func (m *MockData) Type() domain.Type {
	return "Mockdata"
}

func TestRegister(t *testing.T) {
	factory := serializer.NewFactory()
	json.Register(factory, &MockData{})

	// Test data
	testData := &MockData{
		ID:      "123",
		Content: "test content",
	}

	// Test serialization
	serialized, err := factory.Serialize(testData)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	// Test deserialization
	deserialized := &MockData{}
	err = factory.Deserialize(serialized, deserialized)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	// Verify data
	if deserialized.ID != testData.ID || deserialized.Content != testData.Content {
		t.Errorf("Data mismatch after serialization/deserialization:\ngot: %+v\nwant: %+v",
			deserialized, testData)
	}
}

func TestRegistryJSON(t *testing.T) {
	json.RegisterJSON(&MockData{})

	// Test data
	testData := &MockData{
		ID:      "123",
		Content: "test content",
	}

	// Test serialization
	serialized, err := serializer.Serialize(testData)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	// Test deserialization
	deserialized := &MockData{}
	err = serializer.Deserialize(serialized, deserialized)
	if err != nil {
		t.Fatalf("Failed to deserialize: %v", err)
	}

	// Verify data
	if deserialized.ID != testData.ID || deserialized.Content != testData.Content {
		t.Errorf("Data mismatch after serialization/deserialization:\ngot: %+v\nwant: %+v",
			deserialized, testData)
	}
}
