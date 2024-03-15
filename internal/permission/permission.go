package permission

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"git.countmax.ru/pkg/logging"
	"github.com/pkg/errors"

	"git.countmax.ru/countmax/layoutconfig.api/internal/acl"
	"git.countmax.ru/countmax/layoutconfig.api/internal/connmanager"
	"git.countmax.ru/countmax/layoutconfig.api/internal/permission/cache"
)

const (
	XUserID          string        = "X-User-ID"
	XUserIDlow       string        = "X-User-Id"
	XUserEmail       string        = "X-User-Email"
	XUserEmaillow    string        = "X-User-EMAIL"
	XUserPermission  string        = "X-User-Permissions"
	storesByStores   string        = "byStores"
	storesByCities   string        = "byCities"
	storesByRegion   string        = "byRegion"
	storesByCountry  string        = "byCountry"
	fillStoreTimeout time.Duration = 30 * time.Second
)

var (
	DefaultAllow Permissions = Permissions{
		Permission{
			Resources: []string{"*:data.counting:layouts:*"},
			Actions:   acl.Actions{"*"},
			Effect:    "allow",
		},
	}
	DefaultDeny Permissions = Permissions{
		Permission{
			Resources: []string{"*:data.counting:layouts:*"},
			Actions:   acl.Actions{"*"},
			Effect:    "deny",
		},
	}
)

// Manager is main permission controller.
type Manager struct {
	c      cache.RepoInterface
	repoM  *connmanager.Manager
	policy Permissions
	expire time.Duration
}

// NewManager builder for.
func NewManager(c cache.RepoInterface,
	repoM *connmanager.Manager, policy Permissions, expireCache time.Duration) *Manager {
	return &Manager{
		c: c, repoM: repoM, policy: policy, expire: expireCache,
	}
}

// FromRequest extracts permission value from http request,
// if permissions not passed use policy.
func (m *Manager) FromRequest(r *http.Request) Permissions {
	log := logging.FromContext(r.Context())
	bts, err := b64.StdEncoding.DecodeString(r.Header.Get(XUserPermission))
	if err != nil {
		log.Warnf("extract permissions from request failed %s, set default policy", err)
		return m.policy
	}
	perms := make(Permissions, 0, 2)
	err = json.NewDecoder(bytes.NewReader(bts)).Decode(&perms)
	if err != nil {
		log.Warnf("decode permissions failed %s, set default policy", err)
		return m.policy
	}
	return perms
}

// CheckLayout checks permission to Layout,
// true - if allowed, false - if not.
func (m *Manager) CheckLayout(r *http.Request, layoutID string, action acl.Action) bool {
	ps := m.FromRequest(r)
	return ps.CheckLayout(layoutID, action)
}

// SetPolicy sets policy,
// checks input permissions, if empty then no change.
func (m *Manager) SetPolicy(p Permissions) {
	if len(p) == 0 {
		return
	}
	m.policy = p
}

// CheckStore checks access right to action
// by exists or not in the store list.
func (m *Manager) CheckStore(r *http.Request, layoutID, storeID string, action acl.Action) bool {
	if m.c == nil {
		logging.FromContext(r.Context()).Error("nil cache")
		return false
	}
	u, err := m.getUser(r)
	if err != nil {
		logging.FromContext(r.Context()).Error("getUser error, %s", err)
		return false
	}
	// check by store layer
	if itemsStores, ok := u.ACLs[cache.LayoutID(layoutID)][acl.EntityKindStores]; ok {
		allow, hasItems := itemsStores.Check(storeID, action)
		if hasItems {
			return allow
		}
	}
	// no store permissions, use layout permissions
	ps := m.FromRequest(r)
	return ps.CheckLayout(layoutID, action)
}

// FilteredStores processes request, extracts permissions, user, check permissions for action,
// makes comma separated list of the allowed stores; if allowed all stores returns asterisk (*);
//
// Attention! If not allowed stores returns empty string "".
func (m *Manager) FilteredStores(r *http.Request,
	layoutID, inputList string, action acl.Action) (filteredList string, err error) {
	//
	input := strings.Split(inputList, ",")
	output, err := m.filteredStoresSlice(r, layoutID, action, input)
	if err != nil {
		return strings.Join(output, ","), err
	}
	result := strings.Join(output, ",")
	if result == "*" && len(inputList) > 0 {
		result = inputList
	}
	return result, nil
}

// enters

// CheckEnter checks access right to action
// by exists or not in the enter list.
func (m *Manager) CheckEnter(r *http.Request, layoutID, enterID string, action acl.Action) bool {
	if m.c == nil {
		logging.FromContext(r.Context()).Error("nil cache")
		return false
	}
	u, err := m.getUser(r)
	if err != nil {
		logging.FromContext(r.Context()).Error("getUser error, %s", err)
		return false
	}
	// check by store layer
	if itemsEnters, ok := u.ACLs[cache.LayoutID(layoutID)][acl.EntityKindEnters]; ok {
		allow, hasItems := itemsEnters.Check(enterID, action)
		if hasItems {
			return allow
		}
	}
	// no enter permissions, use layout permissions
	ps := m.FromRequest(r)
	return ps.CheckLayout(layoutID, action)
}

// FilteredEnters processes request, extracts permissions, user, check permissions for action,
// makes comma separated list of the allowed enter ids; if allowed all enters returns asterisk (*);
//
// Attention! If not allowed enters returns empty string "".
func (m *Manager) FilteredEnters(r *http.Request,
	layoutID, inputList string, action acl.Action) (filteredList string, err error) {
	//
	input := strings.Split(inputList, ",")
	output, err := m.filteredEnterIDSlice(r, layoutID, action, input)
	if err != nil {
		return strings.Join(output, ","), err
	}
	result := strings.Join(output, ",")
	if result == "*" && len(inputList) > 0 {
		result = inputList
	}
	return result, nil
}

//

//
func (m *Manager) filteredStoresSlice(r *http.Request,
	layoutID string, action acl.Action, input []string) ([]string, error) {
	u, err := m.getUser(r)
	if err != nil {
		return nil, errors.Errorf("getUser error, %s", err)
	}
	if result, ok := u.FilteredStores(cache.LayoutID(layoutID), acl.Actions{action}, input); ok {
		return result, nil
	}
	return []string{}, nil
}

//
func (m *Manager) filteredEnterIDSlice(r *http.Request,
	layoutID string, action acl.Action, input []string) ([]string, error) {
	u, err := m.getUser(r)
	if err != nil {
		return nil, errors.Errorf("getUser error, %s", err)
	}
	if result, ok := u.FilteredEnters(cache.LayoutID(layoutID), acl.Actions{action}, input); ok {
		return result, nil
	}
	return []string{}, nil
}

func (m *Manager) getUser(r *http.Request) (*cache.User, error) {
	uid := m.getUserID(r)
	u, err := m.c.Get(r.Context(), uid)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		return nil, errors.Errorf("get from cache error %s", err)
	}
	if u != nil {
		return u, nil
	}

	ctx, cancel := context.WithTimeout(r.Context(), fillStoreTimeout)
	defer cancel()

	err = m.addUserToCache(ctx, uid, m.FromRequest(r), m.expire)
	if err != nil {
		return nil, errors.Errorf("fillStores error %s", err)
	}

	u, err = m.c.Get(ctx, uid)
	if err != nil && !errors.Is(err, cache.ErrNotFound) {
		return nil, errors.Errorf("get from cache error %s", err)
	}
	if u == nil {
		return nil, errors.Errorf("user %s not found in the cache", uid)
	}
	return u, nil
}

func (m *Manager) getUserID(r *http.Request) string {
	uid := r.Header.Get(XUserID)
	if uid == "" {
		uid = r.Header.Get(XUserIDlow)
	}
	if uid == "" {
		uid = randStringRunes(24)
	}
	return uid
}

// addUserToCache makes user, fills ACLs, puts to cache.
func (m *Manager) addUserToCache(ctx context.Context,
	userID string, permissions Permissions, cacheExpire time.Duration) error {
	//
	// make user
	user := &cache.User{ID: userID, ACLs: make(map[cache.LayoutID]acl.EntityItems, 10)}
	// fill ACLs
	err := m.userFillACL(ctx, user, permissions)
	if err != nil {
		return errors.WithMessage(err, "user fill ACL error")
	}
	// put to cache
	return m.c.Add(ctx, userID, user, cacheExpire)
}

func (m *Manager) userFillACL(ctx context.Context, user *cache.User, permissions Permissions) error {
	for _, permission := range permissions {
		layoutIDs := permission.getLayouts()
		for _, layoutID := range layoutIDs {
			err := m.userFillACLByLayout(ctx, user, layoutID, layoutIDs, permission)
			if err != nil {
				return errors.WithMessage(err, "userFillACLByLayout error")
			}
		}
	}
	return nil
}

func (m *Manager) userFillACLByLayout(ctx context.Context,
	user *cache.User, layoutID string, layoutIDs []string, permission Permission) error {
	//
	allow := permission.Effect == "allow"
	// asterisk
	if layoutID == "*" {
		user.AddItems(cache.LayoutID(layoutID), allow, permission.Actions, layoutIDs, acl.EntityKindLayouts)
		return nil
	}
	storeIDs, enranceIDs, err := m.getEntitiesFromRepo(ctx, layoutID, permission)
	if err != nil {
		return errors.WithMessage(err, "getEntitiesFromRepo error")
	}

	user.AddItems(cache.LayoutID(layoutID), allow, permission.Actions, layoutIDs, acl.EntityKindLayouts)
	user.AddItems(cache.LayoutID(layoutID), allow, permission.Actions, storeIDs, acl.EntityKindStores)
	user.AddItems(cache.LayoutID(layoutID), allow, permission.Actions, enranceIDs, acl.EntityKindEnters)
	return nil
}

func (m *Manager) getEntitiesFromRepo(ctx context.Context,
	layoutID string, permission Permission) (storeIDs, entranceIDs []string, err error) {
	//
	byStores, byCities, byRegion, byCountry := permission.getStores()
	// byStores
	storeIDs, err = m.getStores(ctx, storesByStores, byStores, layoutID)
	if err != nil {
		return storeIDs, entranceIDs, errors.Errorf("get stores byStores error, %s", err)
	}
	// byCities
	listByCities, err := m.getStores(ctx, storesByCities, byCities, layoutID)
	if err != nil {
		return storeIDs, entranceIDs, errors.Errorf("get stores byCities error, %s", err)
	}
	storeIDs = append(storeIDs, listByCities...)
	// byRegion
	listByRegion, err := m.getStores(ctx, storesByRegion, byRegion, layoutID)
	if err != nil {
		return storeIDs, entranceIDs, errors.Errorf("get stores byRegion error, %s", err)
	}
	storeIDs = append(storeIDs, listByRegion...)
	// byCountry
	listByCountry, err := m.getStores(ctx, storesByCountry, byCountry, layoutID)
	if err != nil {
		return storeIDs, entranceIDs, errors.Errorf("get stores byCountry error, %s", err)
	}
	storeIDs = append(storeIDs, listByCountry...)
	entranceIDs, err = m.getEntrances(ctx, layoutID, storeIDs)
	if err != nil {
		return storeIDs, entranceIDs, errors.Errorf("get entrances error, %s", err)
	}
	return storeIDs, entranceIDs, nil
}

func (m *Manager) getEntrances(ctx context.Context, layoutID string, stores []string) ([]string, error) {
	repo, ok := m.repoM.RepoByID(layoutID)
	if !ok {
		return nil, errors.Errorf("not found repo for layoutID %s", layoutID)
	}
	return repo.FindEntrances(ctx, layoutID, strings.Join(stores, ","))
}

// getStores finds storeID list by kind: byList, byCity, byRegion, byCountry.
func (m *Manager) getStores(ctx context.Context, kind, list, layoutID string) ([]string, error) {
	if list == "" {
		return nil, nil
	}
	repo, ok := m.repoM.RepoByID(layoutID)
	if !ok {
		return nil, errors.Errorf("not found repo for layoutID %s", layoutID)
	}
	switch kind {
	case storesByStores:
		return strings.Split(list, ","), nil
	case storesByCities:
		return repo.FindStoresByCities(ctx, layoutID, list)
	case storesByRegion:
		return repo.FindStoresByRegions(ctx, layoutID, list)
	case storesByCountry:
		return repo.FindStoresByCountries(ctx, layoutID, list)
	}
	return nil, errors.Errorf("undefined kind %s", kind)
}

// helpers

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890-+!#$%^&*()_~")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		// nolint:gosec
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
