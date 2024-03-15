package permission

import (
	"testing"

	"git.countmax.ru/countmax/layoutconfig.api/internal/acl"
)

func TestPermissions_CheckLayout(t *testing.T) {
	p1 := &Permissions{
		Permission{
			Resources: []string{"watcom.ru:data.counting:layouts:118416189",
				"watcom.ru:data.counting:cities:961",
				"watcom.ru:data.counting:stores:80817080,80079091"},
			Actions:    acl.Actions{"read"},
			Effect:     "allow",
			Conditions: nil,
		},
		Permission{
			Resources:  []string{"watcom.ru:data.counting:layouts:118416189"},
			Actions:    acl.Actions{"delete"},
			Effect:     "deny",
			Conditions: nil,
		},
	}
	//
	p2 := &DefaultAllow
	p3 := &DefaultDeny
	type args struct {
		action   acl.Action
		layoutID string
	}
	tests := []struct {
		name string
		ps   *Permissions
		args args
		want bool
	}{
		{"allow_read_118416189", p1, args{action: "read", layoutID: "118416189"}, true},
		{"deny_delete_118416189", p1, args{action: "delete", layoutID: "118416189"}, false},
		{"deny_delete_77", p1, args{action: "delete", layoutID: "77"}, false},
		{"allow_read_any_118416189", p2, args{action: "read", layoutID: "118416189"}, true},
		{"allow_read_any_22", p2, args{action: "read", layoutID: "22"}, true},
		{"allow_delete_any_22", p2, args{action: "delete", layoutID: "22"}, true},
		{"deny_delete_any_22", p3, args{action: "delete", layoutID: "22"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ps.CheckLayout(tt.args.layoutID, tt.args.action); got != tt.want {
				t.Errorf("Permissions.CheckLayout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissions_StoresFilters(t *testing.T) {
	p1 := &Permissions{
		Permission{
			Resources: []string{"watcom.ru:data.counting:layouts:118416189",
				"watcom.ru:data.counting:cities:961",
				"watcom.ru:data.counting:stores:80817080,80079091"},
			Actions:    acl.Actions{"read"},
			Effect:     "allow",
			Conditions: nil,
		},
		Permission{
			Resources:  []string{"watcom.ru:data.counting:layouts:118416189"},
			Actions:    acl.Actions{"delete"},
			Effect:     "deny",
			Conditions: nil,
		},
	}
	p2 := &Permissions{
		Permission{
			Resources: []string{"watcom.ru:data.counting:layouts:118416189",
				"watcom.ru:data.counting:countries:1",
				"watcom.ru:data.counting:stores:77"},
			Actions:    acl.Actions{"read"},
			Effect:     "allow",
			Conditions: nil,
		},
		Permission{
			Resources:  []string{"watcom.ru:data.counting:layouts:118416189"},
			Actions:    acl.Actions{"delete"},
			Effect:     "deny",
			Conditions: nil,
		},
	}
	tests := []struct {
		name          string
		ps            *Permissions
		action        string
		wantByStores  string
		wantByCities  string
		wantByRegion  string
		wantByCountry string
	}{
		{"stores_cities_1", p1, "read", "80817080,80079091", "961", "*", "*"},
		{"stores_countries_1", p2, "read", "77", "*", "*", "1"},
		{"default_allow", &DefaultAllow, "read", "*", "*", "*", "*"},
		{"default_deny", &DefaultDeny, "read", "", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotByStores, gotByCities, gotByRegion, gotByCountry := tt.ps.StoresFilters(tt.action)
			if gotByStores != tt.wantByStores {
				t.Errorf("Permissions.FilteredStores() gotByStores = %v, want %v", gotByStores, tt.wantByStores)
			}
			if gotByCities != tt.wantByCities {
				t.Errorf("Permissions.FilteredStores() gotByCities = %v, want %v", gotByCities, tt.wantByCities)
			}
			if gotByRegion != tt.wantByRegion {
				t.Errorf("Permissions.FilteredStores() gotByRegion = %v, want %v", gotByRegion, tt.wantByRegion)
			}
			if gotByCountry != tt.wantByCountry {
				t.Errorf("Permissions.FilteredStores() gotByCountry = %v, want %v", gotByCountry, tt.wantByCountry)
			}
		})
	}
}
