// Package acl contains methods for actual access rules for users, based on the high-level permissions
package acl

import (
	"reflect"
	"testing"
)

func TestItems_filterBySlice(t *testing.T) {
	its1 := &Items{
		Item{ID: "11"}, Item{ID: "22"}, Item{ID: "33", Read: true, Update: true}, Item{ID: "44"}, Item{ID: "54"}, Item{ID: "55"}, Item{ID: "66", Read: true}, Item{ID: "77"},
	}
	tests := []struct {
		name    string
		its     *Items
		actions Actions
		input   []string
		want    []string
	}{
		{"simple_single_exists", its1, Actions{ActionRead}, []string{"33"}, []string{"33"}},
		{"simple_single_not_exists", its1, Actions{ActionRead}, []string{"37"}, []string{}},
		{"simple_many_not_exists", its1, Actions{ActionRead}, []string{"1", "37", "47", "99", "122"}, []string{}},
		{"two_exists_not-exists", its1, Actions{ActionRead}, []string{"33", "54"}, []string{"33"}},
		{"two_exists_update_exists", its1, Actions{ActionUpdate}, []string{"33", "54"}, []string{"33"}},
		{"two_exists", its1, Actions{ActionRead}, []string{"33", "66"}, []string{"33", "66"}},
		{"empty_input_all_output", its1, Actions{ActionRead}, []string{}, []string{"11", "22", "33", "44", "54", "55", "66", "77"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.its.filterBySlice(tt.actions, tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Items.filterBySlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItems_Check(t *testing.T) {
	its1 := &Items{
		Item{ID: "11", Read: true},
		Item{ID: "22", Update: true},
		Item{ID: "33", Delete: true},
		Item{ID: "44", Read: true},
		Item{ID: "55", Update: true, Read: true},
		Item{ID: "66"},
		Item{ID: "77"},
	}
	tests := []struct {
		name   string
		its    *Items
		itemID string
		action Action
		want   bool
	}{
		{"simple_delete_allow", its1, "33", "delete", true},
		{"simple_read_allow", its1, "11", "read", true},
		{"simple_notexists_read_deny", its1, "111", "read", false},
		{"simple_notexists_update_deny", its1, "111", "update", false},
		{"simple_notexists_delete_deny", its1, "111", "delete", false},
		{"simple_notexists_create_deny", its1, "111", "create", false},
		{"simple_55update_allow", its1, "55", "update", true},
		{"simple_55read_allow", its1, "55", "read", true},
		{"simple_55delete_deny", its1, "55", "delete", false},
		{"simple_66delete_deny", its1, "66", "delete", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := tt.its.Check(tt.itemID, tt.action); got != tt.want {
				t.Errorf("Items.check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItems_softMerge(t *testing.T) {
	its1 := &Items{
		Item{ID: "1"}, Item{ID: "10"}, Item{ID: "5"},
	}
	its21 := &Items{}
	wnt1 := &Items{
		Item{ID: "1"}, Item{ID: "10"}, Item{ID: "15", Read: true}, Item{ID: "4", Read: true}, Item{ID: "5"}, Item{ID: "9", Read: true},
	}
	wnt2 := &Items{
		Item{ID: "1"}, Item{ID: "10"}, Item{ID: "15", Read: true}, Item{ID: "19", Read: true}, Item{ID: "4", Read: true}, Item{ID: "5"}, Item{ID: "9", Read: true},
	}
	wnt21 := &Items{
		Item{ID: "15", Read: true}, Item{ID: "4", Read: true}, Item{ID: "9", Read: true},
	}
	input1 := []string{"4", "9", "15"}
	input2 := []string{"4", "19", "15"}
	input21 := []string{"9", "4", "15"}
	tests := []struct {
		name    string
		its     *Items
		actions Actions
		allow   bool
		input   []string
		want    *Items
	}{
		{"simple_4-9-15", its1, Actions{"read"}, true, input1, wnt1},
		{"simple_4-19-15", its1, Actions{"read"}, true, input2, wnt2},
		{"empty_9-4-15", its21, Actions{"read"}, true, input21, wnt21},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.its.softMerge(tt.actions, tt.allow, tt.input)
			if tt.want != nil && tt.its != nil {
				if !reflect.DeepEqual(*tt.its, *tt.want) {
					t.Errorf("softMerge = %v, want to %v", *tt.its, *tt.want)
				}
			}
		})
	}
}

func TestItems_upsert(t *testing.T) {
	its1 := &Items{
		Item{ID: "1"}, Item{ID: "10"}, Item{ID: "5"},
	}
	its21 := &Items{}
	wnt1 := &Items{
		Item{ID: "1"}, Item{ID: "10"}, Item{ID: "15", Read: true}, Item{ID: "4", Read: true}, Item{ID: "5"}, Item{ID: "9", Read: true},
	}
	wnt2 := &Items{
		Item{ID: "1", Read: false}, Item{ID: "10"}, Item{ID: "15", Read: false}, Item{ID: "19", Read: false}, Item{ID: "4", Read: true}, Item{ID: "5"}, Item{ID: "9", Read: true},
	}
	wnt21 := &Items{
		Item{ID: "15", Read: true}, Item{ID: "4", Read: true}, Item{ID: "9", Read: true},
	}
	input1 := []string{"4", "9", "15"}
	input2 := []string{"1", "19", "15"}
	input21 := []string{"9", "4", "15"}
	tests := []struct {
		name    string
		its     *Items
		actions Actions
		allow   bool
		input   []string
		want    *Items
	}{
		{"simple_4-9-15", its1, Actions{"read"}, true, input1, wnt1},
		{"simple_4-19-15", its1, Actions{"read"}, false, input2, wnt2},
		{"empty_9-4-15", its21, Actions{"read"}, true, input21, wnt21},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.its.upsert(tt.actions, tt.allow, tt.input)
			if tt.want != nil && tt.its != nil {
				if !reflect.DeepEqual(*tt.its, *tt.want) {
					t.Errorf("upsert = %v, want to %v", *tt.its, *tt.want)
				}
			}
		})
	}
}
