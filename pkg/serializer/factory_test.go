package serializer_test

import (
	"strings"
	"testing"

	"github.com/mikaelchan/hamster/pkg/domain"
	"github.com/mikaelchan/hamster/pkg/serializer"
)

// MockData is a test struct that implements HasType
type MockData struct {
	ID      string
	Content string
}

func (m *MockData) Type() domain.Type {
	return domain.Type("Mockdata")
}

func TestFactory_RegisterAndSerialize(t *testing.T) {
	factory := serializer.NewFactory()
	// Create a registry for MockData
	factory.Register(&MockData{},
		// Serializer
		func(data domain.HasType) ([]byte, error) {
			mockData := data.(*MockData)
			return []byte(mockData.ID + ";" + mockData.Content), nil
		},
		// Deserializer
		func(data []byte, instance domain.HasType) error {
			mockData := instance.(*MockData)
			parts := strings.Split(string(data), ";")
			mockData.ID = parts[0]
			mockData.Content = parts[1]
			return nil
		},
	)

	// Test data
	testData := MockData{
		ID:      "123",
		Content: "test content",
	}

	// Test Serialize
	t.Run("Serialize", func(t *testing.T) {
		serialized, err := factory.Serialize(&testData)
		if err != nil {
			t.Fatalf("Failed to serialize: %v", err)
		}

		if string(serialized) != testData.ID+";"+testData.Content {
			t.Errorf("Serialized data doesn't match original: got %+v, want %+v", serialized, testData.ID)
		}
	})

	// Test Deserialize
	t.Run("Deserialize", func(t *testing.T) {
		serialized := []byte(testData.ID + ";" + testData.Content)
		var deserialized MockData
		err := factory.Deserialize(serialized, &deserialized)
		if err != nil {
			t.Fatalf("Failed to deserialize: %v", err)
		}

		if deserialized.ID != testData.ID ||
			deserialized.Content != testData.Content {
			t.Errorf("Deserialized data doesn't match original: got %+v, want %+v", deserialized, testData)
		}
	})

	// Test DeserializeNew
	t.Run("DeserializeNew", func(t *testing.T) {
		serialized := []byte(testData.ID + ";" + testData.Content)
		instance, err := factory.DeserializeNew("Mockdata", serialized)
		if err != nil {
			t.Fatalf("Failed to get instance: %v", err)
		}

		mockData := instance.(*MockData)
		if mockData.ID != testData.ID ||
			mockData.Content != testData.Content {
			t.Errorf("Retrieved data doesn't match original: got %+v, want %+v", mockData, testData)
		}
	})
}

func TestGlobalFactory_RegisterAndSerialize(t *testing.T) {
	// Test data
	testData := MockData{
		ID:      "123",
		Content: "test content",
	}
	// Create a registry for MockData
	serializer.Register(&MockData{},
		// Serializer
		func(data domain.HasType) ([]byte, error) {
			mockData := data.(*MockData)
			return []byte(mockData.ID + ";" + mockData.Content), nil
		},
		// Deserializer
		func(data []byte, instance domain.HasType) error {
			mockData := instance.(*MockData)
			parts := strings.Split(string(data), ";")
			mockData.ID = parts[0]
			mockData.Content = parts[1]
			return nil
		},
	)

	// Test Serialize
	t.Run("Serialize", func(t *testing.T) {
		serialized, err := serializer.Serialize(&testData)
		if err != nil {
			t.Fatalf("Failed to serialize: %v", err)
		}

		if string(serialized) != testData.ID+";"+testData.Content {
			t.Errorf("Serialized data doesn't match original: got %+v, want %+v", serialized, testData.ID)
		}
	})

	// Test Deserialize
	t.Run("Deserialize", func(t *testing.T) {
		serialized := []byte(testData.ID + ";" + testData.Content)
		var deserialized MockData
		err := serializer.Deserialize(serialized, &deserialized)
		if err != nil {
			t.Fatalf("Failed to deserialize: %v", err)
		}

		if deserialized.ID != testData.ID ||
			deserialized.Content != testData.Content {
			t.Errorf("Deserialized data doesn't match original: got %+v, want %+v", deserialized, testData)
		}
	})

	// Test Deserialize
	t.Run("DeserializeNew", func(t *testing.T) {
		serialized := []byte(testData.ID + ";" + testData.Content)
		instance, err := serializer.DeserializeNew("Mockdata", serialized)
		if err != nil {
			t.Fatalf("Failed to get instance: %v", err)
		}

		mockData := instance.(*MockData)
		if mockData.ID != testData.ID ||
			mockData.Content != testData.Content {
			t.Errorf("Retrieved data doesn't match original: got %+v, want %+v", mockData, testData)
		}
	})
}

func TestFactory_UnregisteredType(t *testing.T) {
	factory := serializer.NewFactory()
	unregisteredType := domain.Type("unregistered.type")

	// Test with unregistered type
	mockData := MockData{
		ID:      "123",
		Content: "test content",
	}

	_, err := factory.Serialize(&mockData)
	if err == nil {
		t.Error("Expected error for unregistered type, got nil")
	}

	_, err = factory.DeserializeNew(unregisteredType, []byte("{}"))
	if err == nil {
		t.Error("Expected error for unregistered type, got nil")
	}
}
