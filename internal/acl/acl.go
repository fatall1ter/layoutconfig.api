// Package acl contains common methods for actual access rules for users, based on the high-level permissions
package acl

import (
	"sort"
)

const (
	ActionCreate   Action = "create"
	ActionRead     Action = "read"
	ActionUpdate   Action = "update"
	ActionDelete   Action = "delete"
	ActionAsterisk Action = "*"
	EffectAllow    Effect = "allow"
	EffectDeny     Effect = "deny"
	//
	EntityKindLayouts EntityKind = "layouts"
	EntityKindStores  EntityKind = "stores"
	EntityKindEnters  EntityKind = "enters"
	EntityKindZones   EntityKind = "zones"
	EntityKindDevices EntityKind = "devices"
	EntityKindRenters EntityKind = "renters"
)

type EntityKind string

type EntityItems map[EntityKind]Items

//

type Action string

type Actions []Action

func (as Actions) HasAction(act Action) bool {
	for _, a := range as {
		if a == act || a == ActionAsterisk {
			return true
		}
	}
	return false
}

type Effect string

func (e Effect) IsAllow() bool {
	return e == EffectAllow
}

type Item struct {
	ID string
	// Kind   string
	Create bool
	Read   bool
	Update bool
	Delete bool
}

func (it *Item) setActions(actions Actions, allow bool) {
	for _, action := range actions {
		switch action {
		case ActionCreate:
			it.Create = allow
		case ActionRead:
			it.Read = allow
		case ActionUpdate:
			it.Update = allow
		case ActionDelete:
			it.Delete = allow
		case ActionAsterisk:
			it.Create = allow
			it.Read = allow
			it.Update = allow
			it.Delete = allow
		}
	}
}

type Items []Item

func NewItems(actions Actions, allow bool, list []string, kind EntityKind) Items {
	res := make(Items, len(list))
	for i, el := range list {
		_item := Item{ID: el}
		_item.setActions(actions, allow)
		res[i] = _item
	}
	sort.Sort(&res)
	return res
}

func NewEntityItems(actions Actions, allow bool, list []string, kind EntityKind) EntityItems {
	resmap := make(map[EntityKind]Items, 1)
	resmap[kind] = NewItems(actions, allow, list, kind)
	return resmap
}

// Merge for "allow" effect only add new items,
// for "deny" makes replace.
// TODO: make more clear merge.
func (its *Items) Merge(actions Actions, allow bool, input []string) {
	if allow {
		its.softMerge(actions, allow, input)
		return
	}
	its.upsert(actions, allow, input)
}

func (its *Items) softMerge(actions Actions, allow bool, input []string) {
	for _, in := range input {
		ids := its.ListIDs()
		ix := sort.SearchStrings(ids, in)
		if ix < its.Len() && ids[ix] == in {
			continue
		}
		_item := Item{ID: in}
		_item.setActions(actions, allow)
		(*its) = append((*its), _item)
	}
	sort.Sort(its)
}

func (its *Items) upsert(actions Actions, allow bool, input []string) {
	for _, in := range input {
		ids := its.ListIDs()
		_item := Item{ID: in}
		_item.setActions(actions, allow)
		ix := sort.SearchStrings(ids, in)
		if ix < its.Len() && ids[ix] == in {
			(*its)[ix] = _item
			continue
		}
		(*its) = append((*its), _item)
	}
	sort.Sort(its)
}

// implement sort interface

func (its *Items) Len() int {
	return len(*its)
}

func (its *Items) Swap(i, j int) {
	(*its)[i], (*its)[j] = (*its)[j], (*its)[i]
}

func (its *Items) Less(i, j int) bool {
	return (*its)[i].ID < (*its)[j].ID
}

// custom methods

func (its *Items) ListIDs() []string {
	res := make([]string, its.Len())
	for i, item := range *its {
		res[i] = item.ID
	}
	return res
}

func (its *Items) FilteredList(actions Actions, filter []string) []string {
	return its.filterBySlice(actions, filter)
}

func (its *Items) filterBySlice(actions Actions, input []string) []string {
	filtered := make([]string, 0, len(input))
	ids := its.ListIDs()
	if len(input) == 0 {
		return ids
	}
	if result, ok := its.inputAsterisk(actions, input); ok {
		return result
	}
	for _, in := range input {
		ix := sort.SearchStrings(ids, in)
		if ix < its.Len() && ids[ix] == in {
			filtered = append(filtered, its.filterByActions(ix, actions)...)
		}
	}
	return filtered
}

func (its *Items) inputAsterisk(actions Actions, input []string) ([]string, bool) {
	if len(input) == 1 {
		if input[0] == "" || input[0] == "*" {
			return its.filterBySliceAsterisk(actions), true
		}
	}
	return nil, false
}

func (its *Items) filterByActions(ix int, actions Actions) []string {
	filtered := make([]string, 0, 4)
	_tm := make(map[string]struct{}, its.Len())
	for _, action := range actions {
		switch action {
		case ActionCreate:
			if (*its)[ix].Create {
				_tm[(*its)[ix].ID] = struct{}{}
			}
			continue
		case ActionRead:
			if (*its)[ix].Read {
				_tm[(*its)[ix].ID] = struct{}{}
			}
			continue
		case ActionUpdate:
			if (*its)[ix].Update {
				_tm[(*its)[ix].ID] = struct{}{}
			}
			continue
		case ActionDelete:
			if (*its)[ix].Delete {
				_tm[(*its)[ix].ID] = struct{}{}
			}
			continue
		}
	}
	for id := range _tm {
		filtered = append(filtered, id)
	}
	return filtered
}

func (its *Items) filterBySliceAsterisk(actions Actions) []string {
	result := make([]string, 0, its.Len())
	for i := range *its {
		result = append(result, its.filterByActions(i, actions)...)
	}
	return result
}

func (its *Items) Exists(itemID string) bool {
	ids := its.ListIDs()
	ix := sort.SearchStrings(ids, itemID)
	if ix < its.Len() && ids[ix] == itemID {
		return true
	}
	return false
}

func (its *Items) Check(itemID string, action Action) (allow, hasItems bool) {
	if its.Len() == 0 {
		return true, false
	}
	ids := its.ListIDs()

	ix := sort.SearchStrings(ids, itemID)
	if aix, ok := its.checkAsterisk(); ok {
		ix = aix
	}
	if ix < its.Len() && (ids[ix] == itemID || ids[ix] == "*") {
		switch action {
		case ActionCreate:
			return (*its)[ix].Create, true
		case ActionRead:
			return (*its)[ix].Read, true
		case ActionUpdate:
			return (*its)[ix].Update, true
		case ActionDelete:
			return (*its)[ix].Delete, true
		default:
			return false, true
		}
	}
	return false, true
}

func (its *Items) checkAsterisk() (int, bool) {
	if len(*its) == 1 {
		return 0, (*its)[0].ID == "*"
	}
	return its.Len(), false
}
