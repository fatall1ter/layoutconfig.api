package infra

import (
	"net/http"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/labstack/echo/v4"
)

// EntitiesResponse http wrapper with metadata
type EntitiesResponse struct {
	Data domain.Entities `json:"data"`
	Metadata
}

// apiEntities docs
// @Summary Get all entities
// @Description get slice of entities with loc, entity_key, parent_key, kind, offset, limit parameters
// @Description layout_id (recommended parameter), if not pass datasource may be not correct
// @Accept  json
// @Produce  json
// @Tags common
// @Param layout_id path string true "uuid format, default=*"
// @Param entity_key query string false "default=*"
// @Param parent_key query string false "default=*"
// @Param kind query string false "default=*"
// @Param loc query string false "location, default=ru"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Success 200 {object} infra.EntitiesResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/entities [get]
func (s *Server) apiEntities(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	loc := c.QueryParam("loc")
	if loc == "" {
		loc = "ru"
	}
	entityKey := c.QueryParam("entity_key")
	parentKey := c.QueryParam("parent_key")
	kind := c.QueryParam("kind")
	repo, _, err := s.getRepo(c)
	if err != nil {
		s.log.Errorf("getRepo error, %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errInvalidDataSource))
	}
	// TODO: add perm CheckLayout read
	entities, count, err := repo.GetEntities(c.Request().Context(), loc, entityKey, parentKey, kind, offset, limit)
	if err != nil {
		s.log.Errorf("repo.GetEntities error %s", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	response := EntitiesResponse{
		Data: entities,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(entities)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}
