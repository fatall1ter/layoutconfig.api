package infra

import (
	"net/http"

	"git.countmax.ru/countmax/layoutconfig.api/internal/acl"
	"github.com/labstack/echo/v4"
)

// middlewareCheckLayout - sec middleware for check access to layout
func (s *Server) middlewareCheckLayout(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, layoutID, err := s.getRepo(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
		}
		action := extractAction(c)
		if !s.perm.CheckLayout(c.Request(), layoutID, action) {
			s.log.Warnf("not permitted action %s for layout %s", action, layoutID)
			return c.JSON(http.StatusForbidden, ErrForbidden(nil))
		}
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

// middlewareCheckStore - sec middleware for check access to store in specified layout
func (s *Server) middlewareCheckStore(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, layoutID, err := s.getRepo(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
		}
		action := extractAction(c)
		storeID := c.Param("store_id")
		if storeID == "" {
			s.log.Errorf("bad request %v", errEmptyID)
			return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
		}
		if !s.perm.CheckStore(c.Request(), layoutID, storeID, action) {
			s.log.Warnf("not permitted action %s for layout %s and store %s",
				action, layoutID, storeID)
			return c.JSON(http.StatusForbidden, ErrForbidden(nil))
		}
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

// middlewareCheckEnter - sec middleware for check access to entrance in specified layout
func (s *Server) middlewareCheckEnter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, layoutID, err := s.getRepo(c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
		}
		action := extractAction(c)
		entranceID := c.Param("entrance_id")
		if entranceID == "" {
			s.log.Errorf("bad request %v", errEmptyID)
			return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
		}
		if !s.perm.CheckEnter(c.Request(), layoutID, entranceID, action) {
			s.log.Warnf("not permitted action %s for layout %s and entrance %s",
				action, layoutID, entranceID)
			return c.JSON(http.StatusForbidden, ErrForbidden(nil))
		}
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

// helpers

func extractAction(c echo.Context) acl.Action {
	action := acl.ActionRead
	switch c.Request().Method {
	case http.MethodPost:
		action = acl.ActionCreate
	case http.MethodPut, http.MethodPatch:
		action = acl.ActionUpdate
	case http.MethodDelete:
		action = acl.ActionDelete
	}
	return action
}
