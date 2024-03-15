package permission_test

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"math/rand"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/internal/acl"
	cm "git.countmax.ru/countmax/layoutconfig.api/internal/connmanager"
	"git.countmax.ru/countmax/layoutconfig.api/internal/permission"
	"git.countmax.ru/countmax/layoutconfig.api/internal/permission/cache/mem"
	"git.countmax.ru/countmax/layoutconfig.api/repos"
)

/*
[x] CheckLayout(r *http.Request, layoutID, action string) bool

[x] CheckStore(r *http.Request, layoutID, storeID, action string) bool
	s.apiChainRecommendationsQueue
	s.apiChainStoreByID
[ ] CheckDevice(r *http.Request, layoutID, deviceID, action string) bool
	s.apiChainDeviceByID

[ ] CheckEntrance(r *http.Request, layoutID, enterID, action string) bool
	s.apiChainEntranceByID
	s.apiDeleteChainEntrance
[ ] CheckMallEntrance
	s.apiMallEntranceByID read
[ ] CheckZone(r *http.Request, layoutID, zoneID, action string) bool
	s.apiChainZoneByID
[ ] CheckMallZone
	s.apiMallZoneByID
[ ] CheckRenter
	s.apiRenterByID
[x] FilteredStores(r *http.Request, layoutID, action, inputList string) (filteredList string, err error)
	s.apiChainStoresDataAttendance
	s.apiChainStoresDataQueue
	s.apiChainStoresDataQueueNow
	s.apiChainEntrances
	s.apiChainStores
	s.apiChainZones

[ ] FilteredChainZones(r *http.Request, layoutID, action, inputList string) (filteredList string, err error)
	s.apiChainZonesDataAttendance
	s.apiChainZonesDataQueue
	s.apiChainZonesDataQueueNow
	s.apiZoneDataInsideDay
	s.apiZoneDataInsideRange
[ ] FilteredZones mall
	s.apiMallZonesDataAttendance
[ ] FilteredChainEnters(r *http.Request, layoutID, action, inputList string) (filteredList string, err error)
	s.apiChainEntrancesDataAttendance
[ ] FilteredEntersMall
	s.apiMallEntrancesDataAttendance
	s.apiMallEntrances
[ ] FilteretRenters
	s.apiRenterDataAttendance
	s.apiRenters

Test - юзер, пермишн - кладем в реквест и дальше тестируем функцию на ожидаемый результат
добавить базы
кэш всегда в памяти

LayoutID	StoreID		ID_Enter	NameEnter							ID_TypeEnter	NameProject		NameFloor			Name
37664168	124804116	118042379	Вход в "Бренд-2, Магазин-1"			3				Демо-2, сеть	Бренд-2, Магазин-1	Москва
37664168	124804116	99797408	Мимоходящие "Бренд-2, Магазин-1"	11				Демо-2, сеть	Бренд-2, Магазин-1	Москва
37664168	39700940	18130652	Вход в "Бренд-2, Магазин-2"			3				Демо-2, сеть	Бренд-2, Магазин-2	Москва
37664168	39700940	163312801	Мимоходящие "Бренд-2, Магазин-2"	11				Демо-2, сеть	Бренд-2, Магазин-2	Москва
37664168	52944592	133778676	Вход в "Бренд-2, Магазин-3"			3				Демо-2, сеть	Бренд-2, Магазин-3	Санкт-Петербург
118416189	80817080	76469786	Вход в "Бренд-1, Магазин-1"			3				Демо-1, сеть	Бренд-1, Магазин-1	Москва
118416189	80817080	73778046	Мимоходящие "Бренд-1, Магазин-1"	11				Демо-1, сеть	Бренд-1, Магазин-1	Москва
118416189	82949216	90168411	Вход-1 в "Бренд-1, Магазин-2"		3				Демо-1, сеть	Бренд-1, Магазин-2	Санкт-Петербург
118416189	82949216	103594266	Вход-2 в "Бренд-1, Магазин-2"		3				Демо-1, сеть	Бренд-1, Магазин-2	Санкт-Петербург
118416189	109332900	80728066	Вход в "Бренд-1, Магазин-3"			3				Демо-1, сеть	Бренд-1, Магазин-3	Москва
118416189	109332900	160074276	Мимоходящие "Бренд-1, Магазин-3"	11				Демо-1, сеть	Бренд-1, Магазин-3	Москва
118416189	80079091	125067677	Вход в "Бренд-1, Магазин-4"			3				Демо-1, сеть	Бренд-1, Магазин-4	Москва
118416189	147298805	48777088	Вход в "Бренд-1, Магазин-5"			3				Демо-1, сеть	Бренд-1, Магазин-5	Санкт-Петербург
118416189	147298805	142758817	Мимоходящие "Бренд-1, Магазин-5"	11				Демо-1, сеть	Бренд-1, Магазин-5	Санкт-Петербург
*/

const (
	repoTimeout time.Duration = 30 * time.Second
	// layouts
	layoutDemoNet1 string = "118416189"
	layoutDemoNet2 string = "37664168"
	layoutDemoMall string = "73685311"
	// stores
	storeListNet1Spb string = "147298805,82949216"
	storeListNet1    string = "109332900,147298805,80079091,80817080,82949216"
	// 147298805,80079091,82949216
	storeNet1Spb1            string = "82949216"  // Бренд-1, Магазин-2
	storeNet1Spb2            string = "147298805" // Бренд-1, Магазин-5
	storeListNet1SpbAndMsk1  string = "147298805,80079091,82949216"
	storeNet1Msk1            string = "80079091"  // Бренд-1, Магазин-4
	storeNet1Msk2            string = "80817080"  // Бренд-1, Магазин-1
	enterStoreNet1Msk1       string = "125067677" // Вход в "Бренд-1, Магазин-4"
	enterStoreNet1Msk2       string = "76469786"  // Вход в "Бренд-1, Магазин-1"
	enterPassByStoreNet1Msk2 string = "76469786"  // Вход в "Бренд-1, Магазин-1"
	enter1StoreNet1Spb1      string = "90168411"  //  Вход-1 в "Бренд-1, Магазин-2"
	enter2StoreNet1Spb1      string = "103594266" //  Вход-2 в "Бренд-1, Магазин-2"
	enter1and2StoreNet1Spb1  string = "90168411,103594266"
	entersNet1               string = "76469786,73778046,90168411,103594266,80728066,160074276,125067677,48777088,142758817"
	entersNet1SPB            string = "90168411,103594266,48777088,142758817"
	entersNet1SPBAndMSK1     string = "103594266,125067677,142758817,48777088,90168411"
	entersNet1SPB1SPB2       string = "103594266,142758817,48777088,90168411"
	entersNet1MSK            string = "76469786,73778046,80728066,160074276,125067677"

	//
	storeNet2Msk2    string = "39700940"  // Бренд-2, Магазин-2
	storeNet2Msk1    string = "124804116" // Бренд-2, Магазин-1
	storeNet2Spb1    string = "52944592"  // Бренд-2, Магазин-1
	storeListNet2MSK string = "124804116,39700940"
	storeListNet2    string = "124804116,39700940,52944592"
	//
	storeListNet12 string = "124804116,39700940,52944592"
	// 118042379,133778676,163312801,18130652,99797408
	entersNet2            string = "118042379,133778676,163312801,18130652,99797408"
	entersNet2MSK         string = "118042379,163312801,18130652,99797408"
	entersNet2SPB         string = "133778676"
	enter1Net2MSK1        string = "118042379"          //	Вход в "Бренд-2, Магазин-1"
	enterPassByNet2MSK1   string = "99797408"           //	Мимоходящие "Бренд-2, Магазин-1"
	enter1AndPassNet2MKS1 string = "118042379,99797408" // Вход в "Бренд-2, Магазин-1" + Мимоходящие "Бренд-2, Магазин-1"
	enterNet2Spb1         string = "133778676"
	// user
	defaultUID string = "7"

	// permissions
	permAll                   string = `[{"resources":["watcom.ru:data.counting:layouts:*"],"actions":["*"],"effect":"allow"}]`
	permAllReader             string = `[{"resources":["watcom.ru:data.counting:layouts:*"],"actions":["read"],"effect":"allow"}]`
	permChainReaderLDemoNet1  string = `[{"resources":["watcom.ru:data.counting:layouts:118416189"],"actions":["read"],"effect":"allow"}]`
	permChainReaderLDemoNet2  string = `[{"resources":["watcom.ru:data.counting:layouts:37664168"],"actions":["read"],"effect":"allow"}]`
	permChainReaderLDemoNet12 string = `[{"resources":["watcom.ru:data.counting:layouts:37664168,118416189"],"actions":["read"],"effect":"allow"}]`
	// 7
	permChainAndCityMSKReaderLDemoNet1 string = `[{"resources":["watcom.ru:data.counting:layouts:118416189","watcom.ru:data.counting:cities:743"],"actions":["read"],"effect":"allow"}]`

	permChainAndCityMSKReaderLDemoNet2                string = `[{"resources":["watcom.ru:data.counting:layouts:37664168","watcom.ru:data.counting:cities:743"],"actions":["read"],"effect":"allow"}]`
	permChainAndCitySPBReaderLDemoNet1                string = `[{"resources":["watcom.ru:data.counting:layouts:118416189","watcom.ru:data.counting:cities:961"],"actions":["read"],"effect":"allow"}]`
	permChainAndCitySPBPlusStoreReaderLDemoNet1       string = `[{"resources":["watcom.ru:data.counting:layouts:118416189","watcom.ru:data.counting:cities:961","watcom.ru:data.counting:stores:80079091"],"actions":["read"],"effect":"allow"}]`
	permChainAndRegion157MSKReaderLDemoNet1           string = `[{"resources":["watcom.ru:data.counting:layouts:118416189","watcom.ru:data.counting:regions:157"],"actions":["read"],"effect":"allow"}]`
	permChainAndRegion157MSKEx80079091ReaderLDemoNet1 string = `[{"resources":["watcom.ru:data.counting:layouts:118416189","watcom.ru:data.counting:regions:157"],"actions":["read"],"effect":"allow"},{"resources":["watcom.ru:data.counting:layouts:118416189","watcom.ru:data.counting:stores:80079091"],"actions":["*"],"effect":"deny"}]`
	permChainAndCountry1000RUSReaderLDemoNet1         string = `[{"resources":["watcom.ru:data.counting:layouts:118416189","watcom.ru:data.counting:countries:1000"],"actions":["read"],"effect":"allow"}]`
	permChainAndStore82949216ReaderLDemoNet1          string = `[{"resources":["watcom.ru:data.counting:layouts:118416189","watcom.ru:data.counting:stores:82949216"],"actions":["read"],"effect":"allow"}]`
	permChainAndStoreSPB2SPB1LDemoNet1                string = `[{"resources":["watcom.ru:data.counting:layouts:118416189","watcom.ru:data.counting:stores:147298805,82949216"],"actions":["read"],"effect":"allow"}]`

	// conn string
	testCSNets string = "sqlserver://layoutconfig.api:layoutconfig.api@sql-05.watcom.local:1433?database=CM_DemoNetWDA523&connection_timeout=0&encrypt=disable"
	testCSMall string = "sqlserver://layoutconfig.api:layoutconfig.api@sql-05.watcom.local:1433?database=CM_DemoMallWDA523&connection_timeout=0&encrypt=disable"
)

func preparePM(ctx context.Context, css []string, policy permission.Permissions) (*permission.Manager, error) {
	m := &cm.Manager{RWMutex: &sync.RWMutex{}}
	m.InitRepos()
	for _, cs := range css {
		cmr, err := repos.NewLayoutRepo(ctx, repos.LRcountMax523, cs, repoTimeout)
		if err != nil {
			return nil, err
		}
		err = m.RegisterRepo(cmr)
		if err != nil {
			return nil, err
		}
	}
	// cache in mem
	c := mem.New(time.Hour)
	pm := permission.NewManager(c, m, policy, time.Hour)
	return pm, nil
}

func TestManager_CheckLayout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	css := []string{testCSNets}
	m, err := preparePM(ctx, css, permission.DefaultAllow)
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		return
	}
	//
	tests := []struct {
		name   string
		uid    string
		perm   string
		layout string
		action acl.Action
		policy permission.Permissions
		allow  bool
	}{
		// random
		{"1.demoNet1_allow_read_chain_randomID", randStringRunes(7), permChainReaderLDemoNet1, "00000000", acl.ActionRead, permission.DefaultAllow, false},
		{"2.demoNet1_allow_read_chain_check_create_randomID", randStringRunes(7), permChainReaderLDemoNet1, "33333333", acl.ActionCreate, permission.DefaultAllow, false},
		// net1
		{"3.demoNet1_allow_read_chain", randStringRunes(7), permChainReaderLDemoNet1, layoutDemoNet1, acl.ActionRead, permission.DefaultAllow, true},
		{"4.demoNet1_allow_read_chain_check_create", randStringRunes(7), permChainReaderLDemoNet1, layoutDemoNet1, acl.ActionCreate, permission.DefaultAllow, false},
		{"5.demoNet1_allow_read_chain_check_update", randStringRunes(7), permChainReaderLDemoNet1, layoutDemoNet1, acl.ActionUpdate, permission.DefaultAllow, false},
		{"6.demoNet1_allow_read_chain_check_delete", randStringRunes(7), permChainReaderLDemoNet1, layoutDemoNet1, acl.ActionDelete, permission.DefaultAllow, false},
		{"7.demoNet1_checkDemoNet2_read_chain", randStringRunes(7), permChainReaderLDemoNet1, layoutDemoNet2, acl.ActionRead, permission.DefaultAllow, false},
		{"8.demoNet1_checkDemoNet2_delete_chain", randStringRunes(7), permChainReaderLDemoNet1, layoutDemoNet2, acl.ActionDelete, permission.DefaultAllow, false},
		// net2
		{"9.demoNet2_allow_read_chain", randStringRunes(7), permChainReaderLDemoNet2, layoutDemoNet2, acl.ActionRead, permission.DefaultAllow, true},
		{"10.demoNet2_allow_read_chain_check_create", randStringRunes(7), permChainReaderLDemoNet2, layoutDemoNet2, acl.ActionCreate, permission.DefaultAllow, false},
		{"11.demoNet2_allow_read_chain_check_update", randStringRunes(7), permChainReaderLDemoNet2, layoutDemoNet2, acl.ActionUpdate, permission.DefaultAllow, false},
		{"12.demoNet2_allow_read_chain_check_delete", randStringRunes(7), permChainReaderLDemoNet2, layoutDemoNet2, acl.ActionDelete, permission.DefaultAllow, false},
		{"13.demoNet2_checkDemoNet1_read_chain", randStringRunes(7), permChainReaderLDemoNet2, layoutDemoNet1, acl.ActionRead, permission.DefaultAllow, false},
		{"14.demoNet2_checkDemoNet1_delete_chain", randStringRunes(7), permChainReaderLDemoNet2, layoutDemoNet1, acl.ActionDelete, permission.DefaultAllow, false},
		// default policy
		{"15.defPolicyAllow_read_chain", randStringRunes(7), "", layoutDemoNet1, acl.ActionRead, permission.DefaultAllow, true},
		{"16.defPolicyAllow_read_EmptyChain", "", "", "", acl.ActionRead, permission.DefaultAllow, true},
		{"17.defPolicyAllow_create_chain", randStringRunes(7), "", layoutDemoNet1, acl.ActionCreate, permission.DefaultAllow, true},
		{"18.defPolicyAllow_update_chain", randStringRunes(7), "", layoutDemoNet1, acl.ActionUpdate, permission.DefaultAllow, true},
		{"19.defPolicyAllow_delete_chain", randStringRunes(7), "", layoutDemoNet1, acl.ActionDelete, permission.DefaultAllow, true},
		{"20.defPolicyDeny_create_chain", randStringRunes(7), "", layoutDemoNet1, acl.ActionCreate, permission.DefaultDeny, false},
		{"21.defPolicyDeny_read_chain", randStringRunes(7), "", layoutDemoNet1, acl.ActionRead, permission.DefaultDeny, false},
		{"22.defPolicyDeny_update_chain", randStringRunes(7), "", layoutDemoNet1, acl.ActionUpdate, permission.DefaultDeny, false},
		{"23.defPolicyDeny_delete_chain", randStringRunes(7), "", layoutDemoNet1, acl.ActionDelete, permission.DefaultDeny, false},
		//
		{"24.all_reader_allow_read_chain", randStringRunes(7), permAllReader, layoutDemoNet1, acl.ActionRead, permission.DefaultAllow, true},
		{"25.all_reader_deny_create_chain", randStringRunes(7), permAllReader, layoutDemoNet1, acl.ActionCreate, permission.DefaultAllow, false},
		{"26.all_reader_deny_update_chain", randStringRunes(7), permAllReader, layoutDemoNet1, acl.ActionUpdate, permission.DefaultAllow, false},
		{"27.all_reader_deny_delete_chain", randStringRunes(7), permAllReader, layoutDemoNet1, acl.ActionDelete, permission.DefaultAllow, false},
		{"28.all_reader_allow_read_chain2", randStringRunes(7), permAllReader, layoutDemoNet2, acl.ActionRead, permission.DefaultAllow, true},
		{"29.all_reader_deny_create_chain2", randStringRunes(7), permAllReader, layoutDemoNet2, acl.ActionCreate, permission.DefaultAllow, false},
		{"30.all_reader_deny_update_chain2", randStringRunes(7), permAllReader, layoutDemoNet2, acl.ActionUpdate, permission.DefaultAllow, false},
		{"31.all_reader_deny_delete_chain2", randStringRunes(7), permAllReader, layoutDemoNet2, acl.ActionDelete, permission.DefaultAllow, false},
		//
		{"32.all_allow_read_chain", randStringRunes(7), permAll, layoutDemoNet1, acl.ActionRead, permission.DefaultAllow, true},
		{"33.all_allow_create_chain", randStringRunes(7), permAll, layoutDemoNet1, acl.ActionCreate, permission.DefaultAllow, true},
		{"34.all_allow_update_chain", randStringRunes(7), permAll, layoutDemoNet1, acl.ActionUpdate, permission.DefaultAllow, true},
		{"35.all_allow_delete_chain", randStringRunes(7), permAll, layoutDemoNet1, acl.ActionDelete, permission.DefaultAllow, true},
		{"36.all_allow_read_chain2", randStringRunes(7), permAll, layoutDemoNet1, acl.ActionRead, permission.DefaultAllow, true},
		{"37.all_allow_create_chain2", randStringRunes(7), permAll, layoutDemoNet1, acl.ActionCreate, permission.DefaultAllow, true},
		{"38.all_allow_update_chain2", randStringRunes(7), permAll, layoutDemoNet1, acl.ActionUpdate, permission.DefaultAllow, true},
		{"39.all_allow_delete_chain2", randStringRunes(7), permAll, layoutDemoNet1, acl.ActionDelete, permission.DefaultAllow, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "localhost", nil)
			if err != nil {
				t.Errorf("make request error, %s", err)
				return
			}
			m.SetPolicy(tt.policy)
			request.Header.Add(permission.XUserID, tt.uid)
			request.Header.Add(permission.XUserPermission, b64.StdEncoding.EncodeToString([]byte(tt.perm)))
			got := m.CheckLayout(request, tt.layout, tt.action)
			if got != tt.allow {
				t.Errorf("m.CheckLayout=%t, want %t", got, tt.allow)
			}
		})
	}
}

func TestManager_CheckStore(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	css := []string{testCSNets, testCSMall}
	m, err := preparePM(ctx, css, permission.DefaultAllow)
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		return
	}
	//
	tests := []struct {
		name   string
		uid    string
		perm   string
		layout string
		store  string
		action acl.Action
		policy permission.Permissions
		allow  bool
	}{
		// random
		{"1.net1_allow_read_chain_randomID", randStringRunes(7), permChainReaderLDemoNet1, "00000000", storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, false},
		{"2.net1_allow_read_chain_check_create_randomID", randStringRunes(7), permChainReaderLDemoNet1, "33333333", storeNet1Spb1, acl.ActionCreate, permission.DefaultAllow, false},
		// net1
		{"3.net1_allow_read_storeSPb", randStringRunes(7), permChainAndCitySPBReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, true},
		{"4.net1_allow_read_anyStoreInChain", randStringRunes(7), permChainReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, true},
		{"5.net1_deny_create_anyStoreInChain", randStringRunes(7), permChainReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionCreate, permission.DefaultAllow, false},
		{"6.net1_deny_update_anyStoreInChain", randStringRunes(7), permChainReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionUpdate, permission.DefaultAllow, false},
		{"7.net1_deny_delete_anyStoreInChain", randStringRunes(7), permChainReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionDelete, permission.DefaultAllow, false},
		// net1 store in msk perm in spb + 1 store
		{"8.net1_allow_read_MskStoreInChain", randStringRunes(7), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, true},
		{"9.net1_deny_create_anyStoreInChain", randStringRunes(7), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionCreate, permission.DefaultAllow, false},
		{"10.net1_deny_update_anyStoreInChain", randStringRunes(7), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionUpdate, permission.DefaultAllow, false},
		{"11.net1_deny_delete_anyStoreInChain", randStringRunes(7), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionDelete, permission.DefaultAllow, false},
		// net1 store in msk perm in spb + 1 store
		{"12.net1_allow_read_SpbStoreInChain", randStringRunes(7), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, true},
		{"13.net1_deny_create_SpbStoreInChain", randStringRunes(7), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionCreate, permission.DefaultAllow, false},
		{"14.net1_deny_update_SpbStoreInChain", randStringRunes(7), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionUpdate, permission.DefaultAllow, false},
		{"15.net1_deny_delete_SpbStoreInChain", randStringRunes(7), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionDelete, permission.DefaultAllow, false},
		// net1 store in msk region check store in MSK
		{"16.net1_allow_read_MskStoreInMskRegion", randStringRunes(7), permChainAndRegion157MSKReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, true},
		{"17.net1_deny_create_MskStoreInMskRegion", randStringRunes(7), permChainAndRegion157MSKReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionCreate, permission.DefaultAllow, false},
		{"18.net1_deny_update_MskStoreInMskRegion", randStringRunes(7), permChainAndRegion157MSKReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionUpdate, permission.DefaultAllow, false},
		{"19.net1_deny_delete_MskStoreInMskRegion", randStringRunes(7), permChainAndRegion157MSKReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionDelete, permission.DefaultAllow, false},
		// net1 store in msk region check store in SPB
		{"20.net1_deny_read_SpbStoreInMskRegion", randStringRunes(7), permChainAndRegion157MSKReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, false},
		{"21.net1_deny_create_SpbStoreInMskRegion", randStringRunes(7), permChainAndRegion157MSKReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionCreate, permission.DefaultAllow, false},
		{"22.net1_deny_update_SpbStoreInMskRegion", randStringRunes(7), permChainAndRegion157MSKReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionUpdate, permission.DefaultAllow, false},
		{"23.net1_deny_delete_SpbStoreInMskRegion", randStringRunes(7), permChainAndRegion157MSKReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionDelete, permission.DefaultAllow, false},
		// net1 store in Russia country check store in MSK
		{"24.net1_allow_read_MskStoreInRUS", randStringRunes(7), permChainAndCountry1000RUSReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, true},
		{"25net1_deny_create_MskStoreInRUS", randStringRunes(7), permChainAndCountry1000RUSReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionCreate, permission.DefaultAllow, false},
		{"26.net1_deny_update_MskStoreInRUS", randStringRunes(7), permChainAndCountry1000RUSReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionUpdate, permission.DefaultAllow, false},
		{"27.net1_deny_delete_MskStoreInRUS", randStringRunes(7), permChainAndCountry1000RUSReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionDelete, permission.DefaultAllow, false},
		// net1 store in Russia country check store in MSK
		{"28.net1_allow_read_SpbStoreInRUS", randStringRunes(7), permChainAndCountry1000RUSReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, true},
		{"29.net1_deny_create_SpbStoreInRUS", randStringRunes(7), permChainAndCountry1000RUSReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionCreate, permission.DefaultAllow, false},
		{"30.net1_deny_update_SpbStoreInRUS", randStringRunes(7), permChainAndCountry1000RUSReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionUpdate, permission.DefaultAllow, false},
		{"31.net1_deny_delete_SpbStoreInRUS", randStringRunes(7), permChainAndCountry1000RUSReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionDelete, permission.DefaultAllow, false},
		// net1 store in storeList check that store
		{"32.net1_allow_read_StoreInList", randStringRunes(7), permChainAndStore82949216ReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, true},
		{"33.net1_deny_create_StoreInList", randStringRunes(7), permChainAndStore82949216ReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionCreate, permission.DefaultAllow, false},
		{"34.net1_deny_update_StoreInList", randStringRunes(7), permChainAndStore82949216ReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionUpdate, permission.DefaultAllow, false},
		{"35.net1_deny_delete_StoreInList", randStringRunes(7), permChainAndStore82949216ReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionDelete, permission.DefaultAllow, false},
		// net1 store in storeList check another stores
		{"36.net1_deny_read_StoreInListNotInList", randStringRunes(7), permChainAndStore82949216ReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, false},
		{"37.net1_deny_update_StoreInListNotInList", randStringRunes(7), permChainAndStore82949216ReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionUpdate, permission.DefaultAllow, false},
		{"38.net1_deny_delete_StoreInListNotInList", randStringRunes(7), permChainAndStore82949216ReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionDelete, permission.DefaultAllow, false},
		{"39.net1_deny_read_StoreInListNotInList", randStringRunes(7), permChainAndStore82949216ReaderLDemoNet1, layoutDemoNet1, "111", acl.ActionRead, permission.DefaultAllow, false},
		{"40.net1_deny_update_StoreInListNotInList", randStringRunes(7), permChainAndStore82949216ReaderLDemoNet1, layoutDemoNet1, "111", acl.ActionUpdate, permission.DefaultAllow, false},
		{"41.net1_deny_delete_StoreInListNotInList", randStringRunes(7), permChainAndStore82949216ReaderLDemoNet1, layoutDemoNet1, "111", acl.ActionDelete, permission.DefaultAllow, false},
		// net1 stores in msk region with exclude list: 80079091
		{"42.net1_allow_read_StoreInMSKRegionWithEx", randStringRunes(7), permChainAndRegion157MSKEx80079091ReaderLDemoNet1, layoutDemoNet1, storeNet1Msk2, acl.ActionRead, permission.DefaultAllow, true},
		{"43.net1_deny_update_StoreInMSKRegionWithEx", randStringRunes(7), permChainAndRegion157MSKEx80079091ReaderLDemoNet1, layoutDemoNet1, storeNet1Msk2, acl.ActionUpdate, permission.DefaultAllow, false},
		{"44.net1_deny_read_StoreInMSKRegionWithEx_exStore", randStringRunes(7), permChainAndRegion157MSKEx80079091ReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, false},
		{"45.net1_deny_update_StoreInMSKRegionWithEx_exStore", randStringRunes(7), permChainAndRegion157MSKEx80079091ReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionUpdate, permission.DefaultAllow, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "localhost", nil)
			if err != nil {
				t.Errorf("make request error, %s", err)
				return
			}
			m.SetPolicy(tt.policy)
			request.Header.Add(permission.XUserID, tt.uid)
			request.Header.Add(permission.XUserPermission, b64.StdEncoding.EncodeToString([]byte(tt.perm)))
			got := m.CheckStore(request, tt.layout, tt.store, tt.action)
			if got != tt.allow {
				t.Errorf("m.CheckStore=%t, want %t", got, tt.allow)
			}
		})
	}
}

func TestManager_FromRequest(t *testing.T) {
	request, err := http.NewRequest("GET", "localhost", nil)
	if err != nil {
		t.Errorf("make request error, %s", err)
		return
	}
	// [{"resources":["watcom.ru:data.counting:layouts:118416189","watcom.ru:data.counting:cities:961","watcom.ru:data.counting:stores:80817080,80079091"],
	// "actions":["read"],"effect":"allow","conditions":null},{"resources":["watcom.ru:data.counting:layouts:118416189"],"actions":["delete"],"effect":"deny","conditions":null}]
	p1 := permission.Permissions{
		permission.Permission{
			Resources:  []string{"watcom.ru:data.counting:layouts:118416189", "watcom.ru:data.counting:cities:961", "watcom.ru:data.counting:stores:80817080,80079091"},
			Actions:    acl.Actions{"read"},
			Effect:     "allow",
			Conditions: nil,
		},
		permission.Permission{
			Resources:  []string{"watcom.ru:data.counting:layouts:118416189"},
			Actions:    acl.Actions{"delete"},
			Effect:     "deny",
			Conditions: nil,
		},
	}
	request.Header.Add(permission.XUserPermission, makePerm(p1))

	tests := []struct {
		name string
		m    *permission.Manager
		r    *http.Request
		want permission.Permissions
	}{
		{"simple_1", &permission.Manager{}, request, p1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.m.FromRequest(tt.r)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.getPermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_FilteredStores(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	css := []string{testCSNets, testCSMall}
	m, err := preparePM(ctx, css, permission.DefaultAllow)
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		return
	}
	//
	u1 := randStringRunes(18)
	u2 := randStringRunes(18)
	//
	tests := []struct {
		name             string
		uid              string
		perm             string
		layout           string
		inputList        string
		action           acl.Action
		policy           permission.Permissions
		wantFilteredList string
		wantErr          bool
	}{
		{"1.allowAll_read_MSKStoreNet1", randStringRunes(18), permAll, layoutDemoNet1, storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, storeNet1Msk1, false},
		{"2.allowAll_create_MSKStoreNet1", randStringRunes(18), permAll, layoutDemoNet1, storeNet1Msk1, acl.ActionCreate, permission.DefaultAllow, storeNet1Msk1, false},
		{"3.allowAll_update_MSKStoreNet1", randStringRunes(18), permAll, layoutDemoNet1, storeNet1Msk1, acl.ActionUpdate, permission.DefaultAllow, storeNet1Msk1, false},
		{"4.allowAll_delete_MSKStoreNet1", randStringRunes(18), permAll, layoutDemoNet1, storeNet1Msk1, acl.ActionDelete, permission.DefaultAllow, storeNet1Msk1, false},
		//
		{"5.allowAll_read_SpbStoreNet1", randStringRunes(18), permAll, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, storeNet1Spb1, false},
		{"6.allowAll_read_Any", randStringRunes(18), permAll, layoutDemoNet1, "", acl.ActionRead, permission.DefaultAllow, "*", false},
		//
		{"7.allowReadAll_read_SpbStoreNet1", randStringRunes(18), permAll, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, storeNet1Spb1, false},
		{"8.allowReadAll_read_Any", randStringRunes(18), permAll, layoutDemoNet1, "", acl.ActionRead, permission.DefaultAllow, "*", false},
		// permChainAndCity961PlusStoreReaderLDemoNet1
		{"9.allowNet1SPBCity+StoreMSK_read_SpbStoreNet1", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, storeNet1Spb1, false},
		{"10.allowNet1SPBCity+StoreMSK_read_MskStore1Net1", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, storeNet1Msk1, false},
		{"11.allowNet1SPBCity+StoreMSK_read_empty", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, "", acl.ActionRead, permission.DefaultAllow, storeListNet1SpbAndMsk1, false},
		{"12.allowNet1SPBCity+StoreMSK_read_asterisk", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, "*", acl.ActionRead, permission.DefaultAllow, storeListNet1SpbAndMsk1, false},
		{"13.allowNet1SPBCity+StoreMSK_read_StoreSpb2", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeNet1Spb2, acl.ActionRead, permission.DefaultAllow, storeNet1Spb2, false},
		{"14.allowNet1SPBCity+StoreMSK_read_StoreSpb1+2", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, storeListNet1Spb, acl.ActionRead, permission.DefaultAllow, storeListNet1Spb, false},
		// permChainAndStore2IDReaderLDemoNet1
		{"15.allowNet1TwoStore_read_1st", randStringRunes(18), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, storeNet1Spb1, false},
		{"16.allowNet1TwoStore_read_2nd", randStringRunes(18), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, storeNet1Spb2, acl.ActionRead, permission.DefaultAllow, storeNet1Spb2, false},
		{"17.allowNet1TwoStore_read_Both", randStringRunes(18), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, storeListNet1Spb, acl.ActionRead, permission.DefaultAllow, storeListNet1Spb, false},
		{"18.allowNet1TwoStore_read_OneToOne", randStringRunes(18), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, storeNet1Spb1 + "," + storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, storeNet1Spb1, false},
		{"19.allowNet1TwoStore_read_NotInList", randStringRunes(18), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, "", false},
		// permChainReaderLDemoNet2
		// TODO: проверить разрешение на один лэйаут,  вывод списка по другому!?
		{"20.allowNet2_read_*", randStringRunes(18), permChainReaderLDemoNet2, layoutDemoNet2, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, storeNet1Spb1, false},
		{"21.allowNet2_read_ByL1", randStringRunes(18), permChainReaderLDemoNet2, layoutDemoNet1, "*", acl.ActionRead, permission.DefaultAllow, "", false},
		// permChainAndCity742ReaderLDemoNet2 msk in net2
		{"22.allowNet2MSK_read_1stMSK", randStringRunes(18), permChainAndCityMSKReaderLDemoNet2, layoutDemoNet2, storeNet2Msk1, acl.ActionRead, permission.DefaultAllow, storeNet2Msk1, false},
		{"23.allowNet2MSK_read_2ndMSK", randStringRunes(18), permChainAndCityMSKReaderLDemoNet2, layoutDemoNet2, storeNet2Msk2, acl.ActionRead, permission.DefaultAllow, storeNet2Msk2, false},
		{"24.allowNet2MSK_read_Both", randStringRunes(18), permChainAndCityMSKReaderLDemoNet2, layoutDemoNet2, storeListNet2MSK, acl.ActionRead, permission.DefaultAllow, storeListNet2MSK, false},
		{"25.allowNet2MSK_read_OneMSKOneSPB", randStringRunes(18), permChainAndCityMSKReaderLDemoNet2, layoutDemoNet2, storeNet2Spb1 + "," + storeNet2Msk1, acl.ActionRead, permission.DefaultAllow, storeNet2Msk1, false},
		{"26.allowNet2MSK_read_NotInList", randStringRunes(18), permChainAndCityMSKReaderLDemoNet2, layoutDemoNet2, storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, "", false},
		// permChainReaderLDemoNet12 two chains in list
		{"27.allowNet12_read_Net2Msk1", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, storeNet2Msk1, acl.ActionRead, permission.DefaultAllow, storeNet2Msk1, false},
		{"28.allowNet12_read_Net2Msk2", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, storeNet2Msk2, acl.ActionRead, permission.DefaultAllow, storeNet2Msk2, false},
		{"29.allowNet12_read_Net2Spb1", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, storeNet2Spb1, acl.ActionRead, permission.DefaultAllow, storeNet2Spb1, false},
		{"30.allowNet12_read_Net2All", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, storeListNet2, acl.ActionRead, permission.DefaultAllow, storeListNet2, false},
		{"31.allowNet12_read_Net1Msk1", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, storeNet1Msk1, false},
		{"32.allowNet12_read_Net1Msk2", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, storeNet1Msk2, acl.ActionRead, permission.DefaultAllow, storeNet1Msk2, false},
		{"33.allowNet12_read_Net1Spb1", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, storeNet1Spb1, false},
		{"34.allowNet12_read_Net1All", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, storeListNet1, acl.ActionRead, permission.DefaultAllow, storeListNet1, false},
		{"35.allowNet12_read_Net2ByNet1", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, storeListNet2, acl.ActionRead, permission.DefaultAllow, storeListNet2, false},
		{"36.allowNet12_read_Net1ByNet2", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, storeListNet1, acl.ActionRead, permission.DefaultAllow, storeListNet1, false},
		{"37.allowNet12_read_ByNet1_empty", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, "", acl.ActionRead, permission.DefaultAllow, "*", false},
		{"38.allowNet12_read_ByNet2_*", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, "*", acl.ActionRead, permission.DefaultAllow, "*", false},
		// default policy with two layout in one datasource
		{"39.allowPolicy_read_Net1", u1, "", layoutDemoNet1, "", acl.ActionRead, permission.DefaultAllow, "*", false},
		{"40.allowPolicy_read_Net2", u1, "", layoutDemoNet2, "", acl.ActionRead, permission.DefaultAllow, "*", false},
		//
		{"41.allowAll_readByNet1_*", u2, permAll, layoutDemoNet1, "", acl.ActionRead, permission.DefaultAllow, "*", false},
		{"42.allowAll_readByNet2_*", u2, permAll, layoutDemoNet2, "", acl.ActionRead, permission.DefaultAllow, "*", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "localhost", nil)
			if err != nil {
				t.Errorf("make request error, %s", err)
				return
			}
			m.SetPolicy(tt.policy)
			request.Header.Add(permission.XUserID, tt.uid)
			request.Header.Add(permission.XUserPermission, b64.StdEncoding.EncodeToString([]byte(tt.perm)))
			gotFilteredList, err := m.FilteredStores(request, tt.layout, tt.inputList, tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.FilteredStores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFilteredList != tt.wantFilteredList {
				t.Errorf("Manager.FilteredStores() = %v, want %v", gotFilteredList, tt.wantFilteredList)
			}
		})
	}
}

// enters test

func TestManager_CheckEnter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	css := []string{testCSNets, testCSMall}
	m, err := preparePM(ctx, css, permission.DefaultAllow)
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		return
	}
	//
	tests := []struct {
		name        string
		uid         string
		permissions string
		layout      string
		enter       string
		action      acl.Action
		policy      permission.Permissions
		allow       bool
	}{
		// net1Msk_enterNet1MSK1_allow_read_check_read_true
		// net1Msk - permissions for retail1 with stores by city MSK
		// enterNet1MSK1 - enterID1 in store in Net1 in MSK city
		// allow_read - action typed in permissions
		// check_read - action for check
		// true/false - result of the check
		{"1.  net1Msk_enterNet1MSK1_allow_read_check_read_true", randStringRunes(7), permChainAndCityMSKReaderLDemoNet1, layoutDemoNet1, enterStoreNet1Msk1, acl.ActionRead, permission.DefaultAllow, true},
		{"2.  net1Msk_enterNet1MSK1_allow_read_check_delete_false", randStringRunes(7), permChainAndCityMSKReaderLDemoNet1, layoutDemoNet1, enterStoreNet1Msk1, acl.ActionDelete, permission.DefaultAllow, false},
		{"3.  net2Msk_enterNet1MSK1_allow_read_check_read_false", randStringRunes(7), permChainAndCityMSKReaderLDemoNet2, layoutDemoNet1, enterStoreNet1Msk1, acl.ActionRead, permission.DefaultAllow, false},
		{"4.  net1Msk_enterNet1SPB1_allow_read_check_read_false", randStringRunes(7), permChainAndCityMSKReaderLDemoNet1, layoutDemoNet1, enter1StoreNet1Spb1, acl.ActionRead, permission.DefaultAllow, false},
		{"5.  net1_enterNet1SPB1_allow_read_check_read_true", randStringRunes(7), permChainReaderLDemoNet1, layoutDemoNet1, enter1StoreNet1Spb1, acl.ActionRead, permission.DefaultAllow, true},
		{"6.  net2_enterNet1SPB1_allow_read_check_read_false", randStringRunes(7), permChainReaderLDemoNet2, layoutDemoNet1, enter1StoreNet1Spb1, acl.ActionRead, permission.DefaultAllow, false},
		{"7.  net1ListStores_enter1InList_allow_read_check_read_true", randStringRunes(7), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, enter1StoreNet1Spb1, acl.ActionRead, permission.DefaultAllow, true},
		{"8.  net1ListStores_enter2InList_allow_read_check_read_true", randStringRunes(7), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, enter2StoreNet1Spb1, acl.ActionRead, permission.DefaultAllow, true},
		{"9.  net1ListStores_enterNotInList_allow_read_check_read_false", randStringRunes(7), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, enterStoreNet1Msk1, acl.ActionRead, permission.DefaultAllow, false},
		{"10. net1ListStores_enter2NotInList_allow_read_check_read_false", randStringRunes(7), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, enterStoreNet1Msk2, acl.ActionRead, permission.DefaultAllow, false},
		// permChainAndRegion157MSKEx80079091ReaderLDemoNet1
		{"11.  net1MskExcludeStore1_enterNet1MSK1_allow_read_check_read_false", randStringRunes(7), permChainAndRegion157MSKEx80079091ReaderLDemoNet1, layoutDemoNet1, enterStoreNet1Msk1, acl.ActionRead, permission.DefaultAllow, false},
		{"12.  net1MskExcludeStore1_enterNet1MSK2_allow_read_check_read_false", randStringRunes(7), permChainAndRegion157MSKEx80079091ReaderLDemoNet1, layoutDemoNet1, enterStoreNet1Msk2, acl.ActionRead, permission.DefaultAllow, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "localhost", nil)
			if err != nil {
				t.Errorf("make request error, %s", err)
				return
			}
			m.SetPolicy(tt.policy)
			request.Header.Add(permission.XUserID, tt.uid)
			request.Header.Add(permission.XUserPermission, b64.StdEncoding.EncodeToString([]byte(tt.permissions)))
			got := m.CheckEnter(request, tt.layout, tt.enter, tt.action)
			if got != tt.allow {
				t.Errorf("m.CheckEnter=%t, want %t", got, tt.allow)
			}
		})
	}
}

func TestManager_FilteredEnters(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	css := []string{testCSNets, testCSMall}
	m, err := preparePM(ctx, css, permission.DefaultAllow)
	if err != nil {
		t.Errorf("unexpected error, %s", err)
		return
	}
	//
	u1 := randStringRunes(18)
	u2 := randStringRunes(18)
	u3 := randStringRunes(18)
	//
	tests := []struct {
		name             string
		uid              string
		permissions      string
		layout           string
		inputList        string
		action           acl.Action
		policy           permission.Permissions
		wantFilteredList string
		wantErr          bool
	}{
		{"1.allowAll_read_MSKStoreNet1", randStringRunes(18), permAll, layoutDemoNet1, "", acl.ActionRead, permission.DefaultAllow, "*", false},
		{"2.allowAll_read_MSKStoreNet2", randStringRunes(18), permAll, layoutDemoNet2, "*", acl.ActionRead, permission.DefaultAllow, "*", false},
		{"3.allowAll_update_MSKStoreNet1", randStringRunes(18), permAll, layoutDemoNet1, enter1StoreNet1Spb1, acl.ActionUpdate, permission.DefaultAllow, enter1StoreNet1Spb1, false},
		{"4.allowAll_delete_MSKStoreNet1", randStringRunes(18), permAll, layoutDemoNet1, enterPassByStoreNet1Msk2, acl.ActionDelete, permission.DefaultAllow, enterPassByStoreNet1Msk2, false},
		//
		{"5.allowAll_read_EntersNet1", randStringRunes(18), permAll, layoutDemoNet1, entersNet1, acl.ActionRead, permission.DefaultAllow, entersNet1, false},
		{"6.allowAll_read_EntersNet2", randStringRunes(18), permAll, layoutDemoNet2, entersNet2, acl.ActionRead, permission.DefaultAllow, entersNet2, false},
		{"7.allowReadAll_read_SPBStoreNet1", randStringRunes(18), permAll, layoutDemoNet1, entersNet1SPB, acl.ActionRead, permission.DefaultAllow, entersNet1SPB, false},
		{"8.allowReadAll_read_MSKStoreNet2", randStringRunes(18), permAll, layoutDemoNet2, entersNet2MSK, acl.ActionRead, permission.DefaultAllow, entersNet2MSK, false},
		{"9.allowNet1SPBCity+StoreMSK_read_SpbStoreNet1", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, entersNet1SPB, acl.ActionRead, permission.DefaultAllow, entersNet1SPB, false},
		{"10.allowNet1SPBCity+StoreMSK_read_MskStore1Net1", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, enterStoreNet1Msk1, acl.ActionRead, permission.DefaultAllow, enterStoreNet1Msk1, false},

		{"11.allowNet1SPBCity+StoreMSK_read_empty", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, "", acl.ActionRead, permission.DefaultAllow, entersNet1SPBAndMSK1, false},
		{"12.allowNet1SPBCity+StoreMSK_read_asterisk", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, "*", acl.ActionRead, permission.DefaultAllow, entersNet1SPBAndMSK1, false},
		{"13.allowNet1SPBCity+StoreMSK_read_StoreSpb2", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, enter1StoreNet1Spb1, acl.ActionRead, permission.DefaultAllow, enter1StoreNet1Spb1, false},
		{"14.allowNet1SPBCity+StoreMSK_read_StoreSpb1+2", randStringRunes(18), permChainAndCitySPBPlusStoreReaderLDemoNet1, layoutDemoNet1, entersNet1SPB1SPB2, acl.ActionRead, permission.DefaultAllow, entersNet1SPB1SPB2, false},

		// permChainAndStoreSPB2SPB1LDemoNet1
		{"15.allowNet1TwoStore_read_1st", randStringRunes(18), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, enter1StoreNet1Spb1, acl.ActionRead, permission.DefaultAllow, enter1StoreNet1Spb1, false},
		{"16.allowNet1TwoStore_read_2nd", randStringRunes(18), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, enter2StoreNet1Spb1, acl.ActionRead, permission.DefaultAllow, enter2StoreNet1Spb1, false},
		{"17.allowNet1TwoStore_read_Both", randStringRunes(18), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, enter1and2StoreNet1Spb1, acl.ActionRead, permission.DefaultAllow, enter1and2StoreNet1Spb1, false},
		{"18.allowNet1TwoStore_read_OneToOne", randStringRunes(18), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, enter1StoreNet1Spb1 + "," + enterStoreNet1Msk1, acl.ActionRead, permission.DefaultAllow, enter1StoreNet1Spb1, false},
		{"19.allowNet1TwoStore_read_NotInList", randStringRunes(18), permChainAndStoreSPB2SPB1LDemoNet1, layoutDemoNet1, entersNet1MSK, acl.ActionRead, permission.DefaultAllow, "", false},

		// permChainReaderLDemoNet2
		{"20.allowNet2_read_*", randStringRunes(18), permChainReaderLDemoNet2, layoutDemoNet2, enter1StoreNet1Spb1, acl.ActionRead, permission.DefaultAllow, "", false},
		{"21.allowNet2_read_ByL1", randStringRunes(18), permChainReaderLDemoNet2, layoutDemoNet1, "*", acl.ActionRead, permission.DefaultAllow, "", false},

		// permChainAndCity742ReaderLDemoNet2 msk in net2
		{"22.allowNet2MSK_read_1stMSK", randStringRunes(18), permChainAndCityMSKReaderLDemoNet2, layoutDemoNet2, enter1Net2MSK1, acl.ActionRead, permission.DefaultAllow, enter1Net2MSK1, false},
		{"23.allowNet2MSK_read_2ndMSK", randStringRunes(18), permChainAndCityMSKReaderLDemoNet2, layoutDemoNet2, enterPassByNet2MSK1, acl.ActionRead, permission.DefaultAllow, enterPassByNet2MSK1, false},
		{"24.allowNet2MSK_read_Both", randStringRunes(18), permChainAndCityMSKReaderLDemoNet2, layoutDemoNet2, enter1AndPassNet2MKS1, acl.ActionRead, permission.DefaultAllow, enter1AndPassNet2MKS1, false},
		{"25.allowNet2MSK_read_OneMSKOneSPB", randStringRunes(18), permChainAndCityMSKReaderLDemoNet2, layoutDemoNet2, enterNet2Spb1 + "," + enter1Net2MSK1, acl.ActionRead, permission.DefaultAllow, enter1Net2MSK1, false},
		{"26.allowNet2MSK_read_NotInList", randStringRunes(18), permChainAndCityMSKReaderLDemoNet2, layoutDemoNet2, entersNet1SPB, acl.ActionRead, permission.DefaultAllow, "", false},
		/*
			// permChainReaderLDemoNet12 two chains in list
			{"27.allowNet12_read_Net2Msk1", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, storeNet2Msk1, acl.ActionRead, permission.DefaultAllow, storeNet2Msk1, false},
			{"28.allowNet12_read_Net2Msk2", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, storeNet2Msk2, acl.ActionRead, permission.DefaultAllow, storeNet2Msk2, false},
			{"29.allowNet12_read_Net2Spb1", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, storeNet2Spb1, acl.ActionRead, permission.DefaultAllow, storeNet2Spb1, false},
			{"30.allowNet12_read_Net2All", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, storeListNet2, acl.ActionRead, permission.DefaultAllow, storeListNet2, false},
			{"31.allowNet12_read_Net1Msk1", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, storeNet1Msk1, acl.ActionRead, permission.DefaultAllow, storeNet1Msk1, false},
			{"32.allowNet12_read_Net1Msk2", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, storeNet1Msk2, acl.ActionRead, permission.DefaultAllow, storeNet1Msk2, false},
			{"33.allowNet12_read_Net1Spb1", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, storeNet1Spb1, acl.ActionRead, permission.DefaultAllow, storeNet1Spb1, false},
			{"34.allowNet12_read_Net1All", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, storeListNet1, acl.ActionRead, permission.DefaultAllow, storeListNet1, false},
			{"35.allowNet12_read_Net2ByNet1", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, storeListNet2, acl.ActionRead, permission.DefaultAllow, storeListNet2, false},
			{"36.allowNet12_read_Net1ByNet2", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, storeListNet1, acl.ActionRead, permission.DefaultAllow, storeListNet1, false},
			{"37.allowNet12_read_ByNet1_empty", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet1, "", acl.ActionRead, permission.DefaultAllow, "*", false},
			{"38.allowNet12_read_ByNet2_*", randStringRunes(18), permChainReaderLDemoNet12, layoutDemoNet2, "*", acl.ActionRead, permission.DefaultAllow, "*", false},

		*/
		// default policy with two layout in one datasource
		{"39.allowPolicy_read_Net1", u1, "", layoutDemoNet1, "", acl.ActionRead, permission.DefaultAllow, "*", false},
		{"40.allowPolicy_read_Net2", u1, "", layoutDemoNet2, "", acl.ActionRead, permission.DefaultAllow, "*", false},
		//
		{"41.allowAll_readByNet1_*", u2, permAll, layoutDemoNet1, "", acl.ActionRead, permission.DefaultAllow, "*", false},
		{"42.allowAll_readByNet2_*", u2, permAll, layoutDemoNet2, "", acl.ActionRead, permission.DefaultAllow, "*", false},
		//
		{"43.denyPolicy_readByNet1_*", u3, "", layoutDemoNet1, "", acl.ActionRead, permission.DefaultDeny, "", false},
		{"44.denyPolicy_readByNet2_*", u3, "", layoutDemoNet2, "*", acl.ActionRead, permission.DefaultDeny, "", false},
		/**/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "localhost", nil)
			if err != nil {
				t.Errorf("make request error, %s", err)
				return
			}
			m.SetPolicy(tt.policy)
			request.Header.Add(permission.XUserID, tt.uid)
			request.Header.Add(permission.XUserPermission, b64.StdEncoding.EncodeToString([]byte(tt.permissions)))
			gotFilteredList, err := m.FilteredEnters(request, tt.layout, tt.inputList, tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.FilteredEnters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFilteredList != tt.wantFilteredList {
				t.Errorf("Manager.FilteredEnters() = %v, want %v", gotFilteredList, tt.wantFilteredList)
			}
		})
	}
}

//

func makePerm(ps permission.Permissions) string {
	bts, err := json.Marshal(ps)
	if err != nil {
		return ""
	}
	return b64.StdEncoding.EncodeToString(bts)
}

//

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// permission layout1SPB+MSK1
// [{"resources":["watcom.ru:data.counting:layouts:118416189","watcom.ru:data.counting:cities:961","watcom.ru:data.counting:stores:80079091"],"actions":["read"],"effect":"allow"}]
// b64
// W3sicmVzb3VyY2VzIjpbIndhdGNvbS5ydTpkYXRhLmNvdW50aW5nOmxheW91dHM6MTE4NDE2MTg5Iiwid2F0Y29tLnJ1OmRhdGEuY291bnRpbmc6Y2l0aWVzOjk2MSIsIndhdGNvbS5ydTpkYXRhLmNvdW50aW5nOnN0b3Jlczo4MDA3OTA5MSJdLCJhY3Rpb25zIjpbInJlYWQiXSwiZWZmZWN0IjoiYWxsb3cifV0=
