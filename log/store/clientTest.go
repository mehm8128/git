package store

import (
	"encoding/hex"
	"testing"
)

func TestClient_GetObject(t *testing.T) {
	client, err := NewClient("../testrepo")
	if err != nil {
		t.Fatal(err)
	}
	hashString := "436b835a37fec9c43adfa03ba199484893ef6afd"
	hash, err := hex.DecodeString(hashString)
	if err != nil {
		t.Fatal(err)
	}
	obj, err := client.GetObject(hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(obj.Data))
}
