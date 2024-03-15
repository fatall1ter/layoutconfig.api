package infra

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// apiReferences docs
// @Summary Get all references
// @Description get slice of references
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Accept  json
// @Produce  json
// @Tags reference
// @Param layout_id path string true "uuid format, default=*"
// @Success 200 {array} string
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/references [get]
func (s *Server) apiReferences(c echo.Context) error {
	repo, _, err := s.getRepo(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	refs, err := repo.GetReferences(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrServerInternal(fmt.Errorf("database error, get references")))
	}
	return c.JSON(http.StatusOK, refs)
}

// apiRefCategories docs
// @Summary Get renter category reference
// @Description get renter category reference
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Accept  json
// @Produce  json
// @Tags reference
// @Param layout_id path string true "uuid format, default=*"
// @Success 200 {object} reference.RefCategories
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/references/categories [get]
func (s *Server) apiRefCategories(c echo.Context) error {
	repo, _, err := s.getRepo(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	refs, err := repo.GetRefCategories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrServerInternal(fmt.Errorf("database error, get category reference")))
	}
	return c.JSON(http.StatusOK, refs)
}

// apiRefPriceSegments docs
// @Summary Get renter price segments reference
// @Description get renter price segments reference
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Accept  json
// @Produce  json
// @Tags reference
// @Param layout_id path string true "uuid format, default=*"
// @Success 200 {object} reference.RefPrices
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/references/prices [get]
func (s *Server) apiRefPriceSegments(c echo.Context) error {
	repo, _, err := s.getRepo(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	refs, err := repo.GetRefPriceSegments(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrServerInternal(fmt.Errorf("database error, get prices reference")))
	}
	return c.JSON(http.StatusOK, refs)
}

// apiRefKindZones docs
// @Summary Get kind zone reference
// @Description get kind zone reference
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Accept  json
// @Produce  json
// @Tags reference
// @Param layout_id path string true "uuid format, default=*"
// @Success 200 {object} reference.RefPrices
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/references/kindzone [get]
func (s *Server) apiRefKindZones(c echo.Context) error {
	repo, _, err := s.getRepo(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	refs, err := repo.GetRefKindZones(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrServerInternal(fmt.Errorf("database error, get kind zones reference")))
	}
	return c.JSON(http.StatusOK, refs)
}

// apiRefKindEnters docs
// @Summary Get kind enter reference
// @Description get kind enter reference
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Accept  json
// @Produce  json
// @Tags reference
// @Param layout_id path string true "uuid format, default=*"
// @Success 200 {object} reference.RefPrices
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/references/kindenter [get]
func (s *Server) apiRefKindEnters(c echo.Context) error {
	repo, _, err := s.getRepo(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	refs, err := repo.GetRefKindEnters(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			ErrServerInternal(fmt.Errorf("database error, get kind enters reference")))
	}
	return c.JSON(http.StatusOK, refs)
}
