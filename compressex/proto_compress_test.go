package compressex

import "testing"

type testData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestProtoMarshalUnmarshal(t *testing.T) {
	original := testData{Name: "test", Value: 42}

	compressed, err := ProtoMarshal(original)
	if err != nil {
		t.Fatalf("ProtoMarshal error: %v", err)
	}
	if len(compressed) == 0 {
		t.Fatal("ProtoMarshal returned empty bytes")
	}

	var result testData
	err = ProtoUnmarshal(compressed, &result)
	if err != nil {
		t.Fatalf("ProtoUnmarshal error: %v", err)
	}

	if result.Name != original.Name || result.Value != original.Value {
		t.Fatalf("mismatch: got %+v, want %+v", result, original)
	}
}

func TestProtoUnmarshalEmpty(t *testing.T) {
	var result testData
	err := ProtoUnmarshal([]byte{}, &result)
	if err != nil {
		t.Fatalf("ProtoUnmarshal empty: %v", err)
	}
}

func TestProtoUnmarshalNil(t *testing.T) {
	var result testData
	err := ProtoUnmarshal(nil, &result)
	if err != nil {
		t.Fatalf("ProtoUnmarshal nil: %v", err)
	}
}

func TestProtoMarshalCompression(t *testing.T) {
	original := testData{Name: "a very long string that should compress well aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Value: 12345}

	compressed, err := ProtoMarshal(original)
	if err != nil {
		t.Fatalf("ProtoMarshal error: %v", err)
	}

	var result testData
	err = ProtoUnmarshal(compressed, &result)
	if err != nil {
		t.Fatalf("ProtoUnmarshal error: %v", err)
	}

	if result.Name != original.Name {
		t.Fatalf("name mismatch: got %q, want %q", result.Name, original.Name)
	}
}
