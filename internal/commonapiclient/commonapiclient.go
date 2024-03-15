// Package commonapiclient simple commonapi rest client
package commonapiclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	maxIdleConns int = 10
)

// API repository for commonapi
type API struct {
	url     string
	token   string
	handler *http.Client
}

// New builder for API
func New(url, token string, timeout time.Duration) (*API, error) {
	api := &API{
		url:   url,
		token: token,
	}
	// accept google invalid *.watcom.ru cert error
	// nolint:gosec
	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	// set timeouts and border for opened connections
	api.handler = &http.Client{Timeout: timeout, Transport: &http.Transport{
		MaxIdleConns:    maxIdleConns,
		IdleConnTimeout: timeout,
		TLSClientConfig: cfg,
	}}
	return api, api.healthz(context.Background())
}

// GetConnections returns list of the connection strings to countmax dbs
func (a *API) GetConnections(list string) ([]string, error) {
	//http://elk-01.watcom.local:7007/v2/projects/10984
	ids := strings.Split(list, ",")
	if len(ids) == 0 {
		return nil, fmt.Errorf("empty list not allowed")
	}
	result := make([]string, 0, len(ids))
	for _, id := range ids {
		dsn, err := a.getConnectionByID(id)
		if err != nil {
			return nil, err
		}
		result = append(result, dsn)
	}
	return result, nil
}

// Project - entity of project, business object
type Project struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	TypeID      int64  `json:"typeId"`
	TypeName    string `json:"typeName"`
	ParentID    int64  `json:"parentId"`
	ManagerID   int64  `json:"managerId"`
	ManagerName string `json:"managerName"`
	IsEnabled   bool   `json:"isEnabled"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	DBName      string `json:"dbName"`
	Login       string `json:"login"`
	Password    string `json:"password"`
}

func (p *Project) makeURL() string {
	//sqlserver://%s:%s@%s:%d?database=%s&connection_timeout=0&encrypt=disable
	query := url.Values{}
	query.Add("database", p.DBName)
	query.Add("encrypt", "disable")
	query.Add("connection_timeout", "0") //
	u := &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(p.Login, p.Password),
		Host:   fmt.Sprintf("%s:%d", p.IP, p.Port),
		// Path:  instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	return u.String()
}

func (a *API) getConnectionByID(id string) (string, error) {
	uri := fmt.Sprintf("%s/v2/projects/%s", a.url, id)
	request, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return "", err
	}
	request.Header.Set("accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+a.token)
	resp, err := a.handler.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("uri: %s, bad response, httpcode=%d", uri, resp.StatusCode)
		return "", err
	}
	project := &Project{}
	err = json.NewDecoder(resp.Body).Decode(project)
	if err != nil {
		err = fmt.Errorf("uri: %s, decode error %v, httpcode=%d", uri, err, resp.StatusCode)
		return "", err
	}
	return project.makeURL(), nil
}

func (a *API) healthz(ctx context.Context) error {
	uri := fmt.Sprintf("%s/health", a.url)
	request, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return err
	}
	resp, err := a.handler.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("uri: %s, bad response, httpcode=%d", uri, resp.StatusCode)
	}
	return nil
}

func (a *API) Health(ctx context.Context) error {
	return a.healthz(ctx)
}

func (a *API) Scope() string {
	return "commonapi"
}
func (a *API) Dest() string {
	return a.url
}
