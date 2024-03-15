// Package connmanager contains methods for management of the domain.LayoutRepos.
package connmanager

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"git.countmax.ru/countmax/layoutconfig.api/internal/commonapiclient"
	"git.countmax.ru/countmax/layoutconfig.api/repos"
	"git.countmax.ru/countmax/pkg/logging"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	dtsConfig            string        = "config"
	dtsAPI               string        = "api"
	periodLostRetry      time.Duration = 120 * time.Second
	timeoutRepoOperation time.Duration = 30 * time.Second
)

var (
	allowedDST []string = []string{dtsConfig, dtsAPI}
)

// ExtServiceInterface behavior of the some external service for healthcheck
type ExtServiceInterface interface {
	Dest() string
	Scope() string
	Health(context.Context) error
}

// Manager wraps map of the LayoutRepos by RWmutex.
type Manager struct {
	*sync.RWMutex
	repos        map[string]domain.LayoutRepo
	lostAndFound *[]lostDB
	extSvc       []ExtServiceInterface
}

// lostDB params for make LayoutRepo.
type lostDB struct {
	cs      string
	timeout time.Duration
	scope   string
}

// New builder for the Manager.
func New(ctx context.Context, cfg *viper.Viper) (*Manager, error) {
	lf := make([]lostDB, 0)
	m := &Manager{
		RWMutex:      &sync.RWMutex{},
		repos:        make(map[string]domain.LayoutRepo),
		lostAndFound: &lf,
	}
	err := m.fillRepos(ctx, cfg)
	return m, err
}

func (m *Manager) InitRepos() {
	m.repos = make(map[string]domain.LayoutRepo)
}

// fillRepos make LayoutRepo(s) and fill map in the Manager by source,
// config - single database; api - by list ids through commonapi
func (m *Manager) fillRepos(ctx context.Context, cfg *viper.Viper) error {
	const cfgKey string = "countmax.source"

	path := cfg.GetString(cfgKey)
	if path == "" {
		return errors.Errorf("config doesn't contains [%s] key", cfgKey)
	}

	timeout := cfg.GetDuration("countmax.timeout")
	url := cfg.GetString("countmax.url")
	token := cfg.GetString("countmax.token")
	scope := cfg.GetString("countmax.version")
	listIDs := cfg.GetString("countmax.ids")

	switch path {
	case dtsConfig:
		return m.regRepoByConfig(ctx, scope, url, timeout)
	case dtsAPI:
		if listIDs == "" {
			return errors.Errorf("config countmax.ids key is empty")
		}
		return m.regRepoByAPI(ctx, scope, url, token, listIDs, timeout)
	default:
		return errors.Errorf("config %s is %s, unsupported, allow only %v", cfgKey, path, allowedDST)
	}
}

// RepoByID returns repo for specified layoutID,
// if layoutID == *, then will be returned last added repo
func (m *Manager) RepoByID(layoutID string) (domain.LayoutRepo, bool) {
	m.RLock()
	defer m.RUnlock()
	repo, ok := m.repos[layoutID]
	return repo, ok
}

// Repos returns slice of the all registered domain.LayoutRepo(s).
func (m *Manager) Repos() []domain.LayoutRepo {
	m.RLock()
	defer m.RUnlock()
	repos := make([]domain.LayoutRepo, 0, len(m.repos))
	_tmp := make(map[string]struct{}, len(m.repos))
	for _, repo := range m.repos {
		if _, ok := _tmp[repo.Dest()]; !ok {
			repos = append(repos, repo)
			_tmp[repo.Dest()] = struct{}{}
		}
	}
	return repos
}

// GetExtServices returns all external services for health checking.
func (m *Manager) GetExtServices() []ExtServiceInterface {
	m.RLock()
	defer m.RUnlock()
	return m.extSvc
}

// regRepoByConfig add single repo by connection string in the config.
func (m *Manager) regRepoByConfig(ctx context.Context, scope, url string, timeout time.Duration) error {
	cmr, err := repos.NewLayoutRepo(ctx, scope, url, timeout)
	if err != nil {
		return errors.Wrap(err, "repos.NewLayoutRepo failed")
	}
	err = m.RegisterRepo(cmr)
	if err != nil {
		return errors.Wrap(err, "registerRepo failed")
	}
	return nil
}

// regRepoByAPI adds many repos by get layouts from commonapi and process them connection strings.
func (m *Manager) regRepoByAPI(ctx context.Context, scope, url, token, ids string, timeout time.Duration) error {
	log := logging.FromContext(ctx)
	api, err := commonapiclient.New(url, token, timeout)
	if err != nil {
		return errors.Wrap(err, "commonapiclient.New failed")
	}
	m.extSvc = append(m.extSvc, api)
	css, err := api.GetConnections(ids)
	if err != nil {
		return errors.WithMessage(err, "config countmax.ids key is empty")
	}
	for _, cs := range css {
		cmr, err := repos.NewLayoutRepo(ctx, scope, cs, timeout)
		if err != nil {
			log.Errorf("repos.NewLayoutRepo for %s failed: %s, add to lostAndFound", getSrvPortDB(cs), err)
			m.addLF(lostDB{cs: cs, timeout: timeout, scope: scope})
		}
		err = m.RegisterRepo(cmr)
		if err != nil {
			log.Errorf("registerRepo failed, %s", err)
		}
	}
	go m.lostRetrier(ctx, periodLostRetry)
	return nil
}

// RegisterRepo requests layouts from repo and fill repos maps in the Manager's hidden field.
func (m *Manager) RegisterRepo(repo domain.LayoutRepo) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutRepoOperation)
	defer cancel()
	layouts, _, err := repo.FindLayouts(ctx, "ru", "*", 0, 199)
	if err != nil {
		return errors.Wrap(err, "repo.FindLayouts failed")
	}
	m.Lock()
	defer m.Unlock()
	for _, l := range layouts {
		m.repos[l.ID] = repo
	}
	m.repos["*"] = repo
	m.extSvc = append(m.extSvc, repo)
	return nil
}

// getSrvPortDB returns uniq string with server.port.dbname parts.
func getSrvPortDB(connStr string) string {
	// "server=some.domain.ip;user id=root;password=master;port=1433;database=CM_Net523"
	// "postgres://commonapi:commonapi@elk-01:15432/evolution?sslmode=disable&pool_max_conns=2"
	// sqlserver://root:master@study-app.watcom.local:1433?database=CM_Karpov523&connection_timeout=0&encrypt=disable
	connStr = strings.ReplaceAll(connStr, "|", "%7C")
	connStr = strings.ReplaceAll(connStr, "[", "%5B")
	connStr = strings.ReplaceAll(connStr, "]", "%5D")
	connStr = strings.ReplaceAll(connStr, "{", "%7B")
	connStr = strings.ReplaceAll(connStr, "}", "%7D")
	connStr = strings.ReplaceAll(connStr, "`", "%60")
	connStr = strings.ReplaceAll(connStr, "#", "%23")
	connStr = strings.ReplaceAll(connStr, "%", "%25")
	connStr = strings.ReplaceAll(connStr, "^", "%5E")

	u, err := url.Parse(connStr)
	if err != nil || u.Scheme == "" {
		return connStr
	}

	var srv, port, db string
	// postgres
	if u.Scheme == "postgres" {
		db = strings.ReplaceAll(u.Path, "/", "")
	}
	if u.Scheme == "sqlserver" {
		params := strings.Split(u.RawQuery, "&")
		for _, par := range params {
			if strings.Contains(par, "database=") {
				dbs := strings.Split(par, "=")
				if len(dbs) == 2 {
					db = dbs[1]
				}
			}
		}
	}
	srv = u.Hostname()
	port = u.Port()
	return fmt.Sprintf("[%s].[%s].[%s]", srv, port, db)
}
