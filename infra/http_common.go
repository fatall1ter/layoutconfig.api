package infra

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"git.countmax.ru/countmax/layoutconfig.api/internal/permission"
	"github.com/labstack/echo/v4"
)

const (
	// PAGESIZE number of document per page by default
	PAGESIZE int = 20
	// swaggerURL string = "/swagger/index.html"
	vueURL string = "/web/index.html"
)

type Version struct {
	Build      string `json:"build,omitempty"`
	GitHash    string `json:"git_hash,omitempty"`
	Version    string `json:"version,omitempty"`
	APIVersion string `json:"api_version,omitempty"`
	DBVersion  string `json:"db_version,omitempty"`
}

func (s *Server) getRepo(c echo.Context) (domain.LayoutRepo, string, error) {
	layoutID := c.QueryParam("layout_id")
	originID := layoutID
	s.log.Debugf("got layout_id %s", layoutID)
	if layoutID == "" {
		layoutID = c.Param("layout_id")
		originID = layoutID
		if layoutID == "" {
			layoutID = "*"
			s.log.Warn("layout_id not passed, set to *")
		}
	}
	repo, ok := s.repoM.RepoByID(layoutID)
	if !ok {
		s.log.Errorf("not found repo for layout_id=%s", layoutID)
		return nil, originID, errLayoutRepoNotFound
	}
	s.log.Debugf("for layout_id %s got repo with Dest %s", layoutID, repo.Dest())
	return repo, originID, nil
}

// getPageParams parse http request parameters and returns offset and limit
func (s *Server) getPageParams(c echo.Context) (offset, limit int64) {

	slimit := c.QueryParam("limit")
	if slimit == "" {
		slimit = "0"
	}
	ilimit, err := strconv.Atoi(slimit)
	if ilimit <= 0 || err != nil {
		ilimit = PAGESIZE
	}

	soffset := c.QueryParam("offset")
	if soffset == "" {
		soffset = "0"
	}
	ioffset, err := strconv.Atoi(soffset)
	if ioffset <= 0 || err != nil {
		ioffset = 0
	}
	limit, offset = int64(ilimit), int64(ioffset)
	return
}

// getFromToParams parse http request parameters and returns from, to time,
// if error happeend return default time and error text
func (s *Server) getFromToParams(c echo.Context) (time.Time, time.Time, error) {
	sFrom := c.QueryParam("from")
	sTo := c.QueryParam("to")
	_, toffset := time.Now().Zone()
	if sFrom == "" {
		sFrom = time.Now().Format("2006-01-02T00:00:00Z07:00")
	}
	from, err := time.Parse(time.RFC3339, sFrom)
	if err != nil {
		preErr := err
		fromNaive, er := time.Parse(naiveTimeFormat, sFrom)
		if er != nil {
			return from, time.Time{}, fmt.Errorf("wrong from [%v] parameter %v/%v", sFrom, preErr, er)
		}
		from = fromNaive.Local().Add(-time.Duration(toffset) * time.Second)
	}
	if sTo == "" {
		sTo = time.Now().Format(time.RFC3339)
	}
	to, err := time.Parse(time.RFC3339, sTo)
	if err != nil {
		preErr := err
		toNaive, err := time.Parse(naiveTimeFormat, sTo)
		if err != nil {
			return from, to, fmt.Errorf("wrong to [%v] parameter %v/%v", sTo, preErr, err)
		}
		to = toNaive.Local().Add(-time.Duration(toffset) * time.Second)
	}
	return from, to, nil
}

// helpers

func (s *Server) joinParam(c echo.Context, single, many string) ([]string, error) {
	list := make([]string, 0, 1)
	one := c.QueryParam(single)
	if len(one) > 0 {
		if !s.listReg.MatchString(one) {
			return nil, fmt.Errorf("wrong %s parameter", single)
		}
		list = append(list, one)
	}
	strList := c.QueryParam(many)
	if len(strList) > 0 {
		if !s.listReg.MatchString(strList) {
			return nil, fmt.Errorf("wrong %s parameter", many)
		}
		strList = strings.Trim(strList, ",")
		parts := strings.Split(strList, ",")
		list = append(list, parts...)
	}
	if len(one) == 0 && len(strList) == 0 {
		list = append(list, "*")
	}
	return list, nil
}

func (s *Server) customHTTPErrorHandler(err error, c echo.Context) {
	if err := c.File(vueURL); err != nil {
		c.Logger().Error(err)
	}
}

// customHTTPLogger - middleware of logger and metric duration
func (s *Server) customHTTPLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		rID := c.Response().Header().Get("x-request-id")
		code := c.Response().Status
		uri := c.Path() // c.Request().URL.EscapedPath()
		query := c.Request().URL.Query().Encode()
		uid, email := s.getUserInfo(c)
		httplog := s.log.With(
			"method", c.Request().Method,
			"proto", c.Request().Proto,
			"remote", c.Request().RemoteAddr,
			"url", uri,
			"query", query,
			"code", code,
			"size", c.Response().Size,
			"duration", time.Since(start).String(),
			"uid", uid,
			"email", email,
			requestIDName, rID)
		host, err := os.Hostname()
		if err != nil {
			s.log.Warnf("os.Hostname error, %s", err)
		}
		switch {
		case code < 400:
			httplog.Infof("%s", host)
		case code >= 400 && code < 500:
			httplog.Warnf("%s", host)
		case code >= 500:
			httplog.Errorf("%s", host)
		}
		s.mAPI.WithLabelValues(uri, strconv.Itoa(code), c.Request().Method).
			Observe(time.Since(start).Seconds())
		return nil
	}
}

func (s *Server) getUserInfo(c echo.Context) (string, string) {
	uid := c.Request().Header.Get(permission.XUserID)
	if uid == "" {
		uid = c.Request().Header.Get(permission.XUserIDlow)
	}
	email := c.Request().Header.Get(permission.XUserEmail)
	if uid == "" {
		email = c.Request().Header.Get(permission.XUserEmaillow)
	}
	return uid, email
}

func (s *Server) responserMIME(c echo.Context, code int, payload interface{}) error {
	accept := c.Request().Header.Get("Accept")
	if strings.Contains(accept, "xml") {
		return c.XML(code, payload)
	}
	return c.JSON(code, payload)

}

// apiHealthCheck returk 200 ok if repository is connected
// @Summary Healthcheck service eq repository connected
// @Tags health
// @Success 200 {object} infra.SuccessResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /health [get]
func (s *Server) apiHealthCheck(c echo.Context) error {
	err := s.healthCheck(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	return c.JSON(http.StatusOK, OkStatus("OK"))
}

func (s *Server) metrics(c echo.Context) error {
	return c.Redirect(http.StatusPermanentRedirect, "/metrics")
}

// apiVersion returs version info
// @Summary Get version info
// @Tags health
// @Success 200 {object} infra.Version
// @Router /v2/version [get]
func (s *Server) apiVersion(c echo.Context) error {
	ver := Version{
		Build:      s.build,
		GitHash:    s.githash,
		Version:    s.version,
		APIVersion: apiVersion,
		DBVersion:  s.scope,
	}
	return c.JSON(http.StatusOK, ver)
}

// Metadata - metadata, page, limit, offset... etc...
type Metadata struct {
	ResultSet ResultSet `json:"result_set"`
}

// ResultSet - values total, limit....
type ResultSet struct {
	Count  int64 `json:"count"`
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
	Total  int64 `json:"total"`
}

// SuccessResponse structure for json response success results
type SuccessResponse struct {
	Message        string `json:"message"`  // text of message
	HTTPStatusCode int    `json:"httpcode"` // http response status code
	StatusText     string `json:"status"`   // user-level status message
}

// OkStatus wrapper HTTP 200 OK response
func OkStatus(message string) SuccessResponse {
	return SuccessResponse{
		Message:        message,
		HTTPStatusCode: http.StatusOK,
		StatusText:     http.StatusText(http.StatusOK),
	}
}

// CreatedStatus wrapper HTTP 201 Create response
func CreatedStatus(message string) SuccessResponse {
	return SuccessResponse{
		Message:        message,
		HTTPStatusCode: http.StatusCreated,
		StatusText:     http.StatusText(http.StatusCreated),
	}
}

// ErrResponse structure for common response with some error
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

// ErrInvalidRequest - wrapper for make err structure
func ErrInvalidRequest(err error) ErrResponse {
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     http.StatusText(http.StatusBadRequest),
		ErrorText:      fmt.Sprintf("%v", err),
	}
}

// ErrServerInternal - wrapper for make err structure
func ErrServerInternal(err error) ErrResponse {
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     http.StatusText(http.StatusInternalServerError),
		ErrorText:      fmt.Sprintf("%v", err),
	}
}

// ErrNotFound - wrapper for make err structure for empty result
func ErrNotFound(err error) ErrResponse {
	Error := ""
	if err != nil {
		Error = err.Error()
	}
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     http.StatusText(http.StatusNotFound),
		ErrorText:      Error,
	}
}

// ErrNotFound - wrapper for make err structure for empty result
func ErrUnAuthorized(err error) ErrResponse {
	Error := ""
	if err != nil {
		Error = err.Error()
	}
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     http.StatusText(http.StatusUnauthorized),
		ErrorText:      Error,
	}
}

// ErrForbidden - wrapper for make err structure for forbidden resource
func ErrForbidden(err error) ErrResponse {
	er := ""
	if err != nil {
		er = err.Error()
	}
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusForbidden,
		StatusText:     http.StatusText(http.StatusForbidden),
		ErrorText:      er,
	}
}

// ErrPayloadTooLarge - wrapper for make err structure
func ErrPayloadTooLarge(err error) ErrResponse {
	return ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusRequestEntityTooLarge,
		StatusText:     http.StatusText(http.StatusRequestEntityTooLarge),
		ErrorText:      fmt.Sprintf("%v", err),
	}
}

// ErrUnsupportedFormat - 415 error implementation
var ErrUnsupportedFormat = &ErrResponse{HTTPStatusCode: http.StatusUnsupportedMediaType, StatusText: "415 - Unsupported Media Type."}

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Code     int         `json:"-"`
	Message  interface{} `json:"message"`
	Internal error       `json:"-"` // Stores the error returned by an external dependency
}
