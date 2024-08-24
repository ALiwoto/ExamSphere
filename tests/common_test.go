package tests

import (
	"encoding/json"
	"testing"
	"time"
)

type MyType1 struct {
	DateField time.Time `json:"date_field"`
}

func TestDateJson(t *testing.T) {
	// TestDateJson tests the json serialization of the Date struct
	myValue := &MyType1{
		DateField: time.Now(),
	}
	jsonDate, err := json.Marshal(myValue)
	if err != nil {
		t.Error(err)
	}

	outputStr := string(jsonDate)
	t.Logf("Output: %s", outputStr)
	// now we will try to do it reverse
	myValue2 := &MyType1{}
	err = json.Unmarshal(jsonDate, myValue2)
	if err != nil {
		t.Error(err)
	}

	if myValue.DateField != myValue2.DateField {
		t.Errorf("Expected %s, got %s", myValue.DateField, myValue2.DateField)
	}
}
