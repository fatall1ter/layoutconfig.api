package permission

import (
	"regexp"
	"strings"

	"git.countmax.ru/countmax/layoutconfig.api/internal/acl"
)

var (
	layoutRE    = regexp.MustCompile(`(?m)[\*\w.]+:data.counting:layouts:(\S+)`)
	stCitiesRE  = regexp.MustCompile(`(?m)[\*\w.]+:data.counting:cities:(\S+)`)
	stRegionRE  = regexp.MustCompile(`(?m)[\*\w.]+:data.counting:regions:(\S+)`)
	stCountryRE = regexp.MustCompile(`(?m)[\*\w.]+:data.counting:countries:(\S+)`)
	stStoreRE   = regexp.MustCompile(`(?m)[\*\w.]+:data.counting:stores:(\S+)`)
)

type Permissions []Permission

type Permission struct {
	Resources  Resources   `json:"resources"`
	Actions    acl.Actions `json:"actions"`
	Effect     acl.Effect  `json:"effect"`
	Conditions interface{} `json:"conditions"`
}

type Resources []string

// [{"resources":["watcom.ru:data.counting:layouts:141466237","watcom.ru:data.counting:cities:961"],"actions":["actions:read"],"effect":"allow","conditions":null}]
// ,W3sicmVzb3VyY2VzIjpbIndhdGNvbS5ydTpkYXRhLmNvdW50aW5nOmxheW91dHM6MTQxNDY2MjM3Iiwid2F0Y29tLnJ1OmRhdGEuY291bnRpbmc6Y2l0aWVzOjk2MSJdLCJhY3Rpb25zIjpbImFjdGlvbnM6cmVhZCJdLCJlZmZlY3QiOiJhbGxvdyIsImNvbmRpdGlvbnMiOm51bGx9XQ==
// {"resources":["watcom.ru:data.counting:layouts:118416189"],
// "actions":["delete"],"effect":"deny","conditions":null}]

func (p Permission) checkLayout(layoutID string) bool {
	for _, r := range p.Resources {
		lids := layoutRE.FindStringSubmatch(r)
		for _, lid := range lids {
			ids := strings.Split(lid, ",")
			for _, id := range ids {
				if id == layoutID || id == "*" { // TODO: reduce nesting level
					return true
				}
			}
		}
	}
	return false
}

func (p Permission) getLayouts() []string {
	res := make([]string, 0, 1)
	for _, r := range p.Resources {
		lids := layoutRE.FindStringSubmatch(r)
		if len(lids) == 0 {
			continue
		}
		ids := strings.Split(lids[len(lids)-1], ",")
		res = append(res, ids...)
	}
	return res
}

func (p Permission) checkStore(layoutID, storeID, cityID, regionID, countryID string) bool {
	if !p.checkLayout(layoutID) {
		return false // by layout
	}
	for _, r := range p.Resources {
		lids := layoutRE.FindStringSubmatch(r)
		for _, lid := range lids {
			ids := strings.Split(lid, ",")
			for _, id := range ids {
				if id == layoutID || id == "*" { // TODO: reduce nesting level
					return true
				}
			}
		}
	}
	return false
}

func (p Permission) getStores() (byStores, byCities, byRegion, byCountry string) {
	for _, r := range p.Resources {
		stids := stStoreRE.FindStringSubmatch(r)
		if len(stids) > 0 {
			byStores = stids[len(stids)-1]
		}
		cids := stCitiesRE.FindStringSubmatch(r)
		if len(cids) > 0 {
			byCities = cids[len(cids)-1]
		}
		rids := stRegionRE.FindStringSubmatch(r)
		if len(rids) > 0 {
			byRegion = rids[len(rids)-1]
		}
		coids := stCountryRE.FindStringSubmatch(r)
		if len(coids) > 0 {
			byCountry = coids[len(coids)-1]
		}
	}
	return
}

func (ps *Permissions) CheckLayout(layoutID string, action acl.Action) bool {
	allow := false
	for _, p := range *ps {
		if p.Actions.HasAction(action) {
			if p.checkLayout(layoutID) {
				allow = p.Effect.IsAllow()
			}
		}
	}
	return allow
}

// StoresFilters returns list of the stores, cities, regions, countries
// if permissions does not contains specific rules returns empty string
func (ps *Permissions) StoresFilters(action string) (byStores, byCities, byRegion, byCountry string) {
	for _, p := range *ps {
		if p.Actions.HasAction(acl.Action(action)) {
			if !p.Effect.IsAllow() {
				continue
			}
			byStores, byCities, byRegion, byCountry = p.getStores()
			if byStores == "" {
				byStores = "*"
			}
			if byCities == "" {
				byCities = "*"
			}
			if byRegion == "" {
				byRegion = "*"
			}
			if byCountry == "" {
				byCountry = "*"
			}
		}
	}
	return
}
