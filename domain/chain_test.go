package domain

import (
	"encoding/json"
	"testing"
)

var (
	chain1 Chain = Chain{
		LayoutID:  "123",
		Kind:      "chain",
		Title:     "TestChain1",
		Languages: "[\"ru\"]",
		CRMKey:    "001",
		Brands:    "[TestA]",
		Currency:  "rubs",
		Options:   "{\"someopt\":\"data\"}",
		Notes:     "somenotes",
	}

	chain2 Chain = Chain{
		LayoutID:  "456",
		Kind:      "chain",
		Title:     "TestChain2",
		Languages: "[\"ru\"]",
		CRMKey:    "002",
		Brands:    "[TestB]",
		Currency:  "rubs",
		Options:   "{\"someopt\":\"data\"}",
		Notes:     "somenotes",
	}
	testChains Chains = Chains{chain1, chain2}
)

func TestSetZeroValueChains(t *testing.T) {
	fields := []string{"layout_id", "kind", "title"}
	expJson := `[{"layout_id":"123","kind":"chain","title":"TestChain1","is_active":false,"read_only":false},{"layout_id":"456","kind":"chain","title":"TestChain2","is_active":false,"read_only":false}]`
	jsonChBefore, err := json.Marshal(testChains)
	if err != nil {
		t.Errorf("expected success json marshalling, but error %v", err)
	}
	if string(jsonChBefore) == expJson {
		t.Error("oops, start data same as expected, not need continue, wrong input data")
	}
	testChains.SetZeroValue(fields)
	jsonCh, err := json.Marshal(testChains)
	if err != nil {
		t.Errorf("expected success json marshalling, but error %v", err)
	}
	if string(jsonCh) != expJson {
		t.Errorf("expected: %s\nbut got: %s\n", expJson, string(jsonCh))
	}
}

func TestSetZeroValueChain(t *testing.T) {
	fields := []string{"layout_id", "kind", "title"}
	expJson := `{"layout_id":"123","kind":"chain","title":"TestChain1","is_active":false,"read_only":false}`
	fields2 := []string{}
	expJson2 := `{"layout_id":"123","kind":"chain","title":"TestChain1","is_active":false,"read_only":false}`
	jsonChBefore, err := chain1.MarshalJSON()
	if err != nil {
		t.Errorf("expected success json marshalling, but error %v", err)
	}
	if string(jsonChBefore) == expJson {
		t.Error("oops, start data same as expected, not need continue, wrong input data")
	}
	chain1.SetZeroValue(fields)
	jsonCh, err := chain1.MarshalJSON()
	if err != nil {
		t.Errorf("expected success json marshalling, but error %v", err)
	}
	if string(jsonCh) != expJson {
		t.Errorf("expected: %s\nbut got: %s\n", expJson, string(jsonCh))
	}
	chain1.SetZeroValue(fields2)
	jsonCh2, err := chain1.MarshalJSON()
	if err != nil {
		t.Errorf("expected success json marshalling, but error %v", err)
	}
	if string(jsonCh2) != expJson2 {
		t.Errorf("expected: %s\nbut got: %s\n", expJson2, string(jsonCh2))
	}
}
