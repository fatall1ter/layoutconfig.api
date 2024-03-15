package infra

import (
	"errors"
	"net/http"
	"strconv"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"git.countmax.ru/countmax/layoutconfig.api/internal/acl"
	"github.com/labstack/echo/v4"
)

var (
	errRepo error = errors.New("storage error")
)

// BehaviorsResponse http wrapper with metadata
type BehaviorsResponse struct {
	Data domain.Behaviors `json:"data"`
	Metadata
}

// apiBehaviors docs
// @Summary Get behaviors
// @Description get behaviors for layouts
// @Produce  json
// @Tags common
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Success 200 {object} infra.BehaviorsResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/layouts/behaviors [get]
func (s *Server) apiBehaviors(c echo.Context) error {
	offset, limit := s.getPageParams(c)
	var count int64
	bhvs := make(domain.Behaviors, 0, limit)
LOOP:
	for _, repo := range s.repoM.Repos() {
		layouts, _, err := repo.FindLayouts(c.Request().Context(), "", "", 0, 9999)
		if err != nil {
			s.log.Errorf("repo.FindLayouts error, %s", err)
			return c.JSON(http.StatusInternalServerError, ErrServerInternal(errors.New("find layouts error")))
		}
		for _, layout := range layouts {
			if !s.perm.CheckLayout(c.Request(), layout.ID, acl.ActionRead) {
				s.log.Warnf("not permitted action %s for layout %s", acl.ActionRead, layout.ID)
				continue LOOP
			}
		}
		s.log.Debugf("find behavior from repo %s", repo.GetSrvPortDB())
		_bhvs, _count, err := repo.FindBehaviors(c.Request().Context(), offset, limit)
		if err != nil {
			s.log.Errorf("repo.FindBehaviors erro, %s", err)
			return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
		}
		s.log.Debugf("got behaviors %v", _bhvs)
		currLen := int64(len(bhvs))
		for _, _bhv := range _bhvs {

			if currLen < limit {
				bhvs = append(bhvs, _bhv)
				currLen++
			}
		}
		count += _count
	}
	response := BehaviorsResponse{
		Data: bhvs,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(bhvs)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

// apiBehaviorByLayoutID docs
// @Summary Get behavior for specified layout
// @Description get behaviors for for specified layout_id
// @Produce  json
// @Tags common
// @Param layout_id path string true "digit/uuid format"
// @Success 200 {object} domain.Behavior
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/layouts/{layout_id}/behaviors [get]
func (s *Server) apiBehaviorByLayoutID(c echo.Context) error {
	repo, id, err := s.getRepo(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	if id == "" {
		s.log.Warnf("bad request apiBehaviorByLayoutID, %v", errEmptyID)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
	}
	bhv, err := repo.FindBehaviorByLayoutID(c.Request().Context(), id)
	if err != nil {
		s.log.Errorf("repo.FindBehaviorByLayoutID(%s) error, %s", id, err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusOK, bhv)
}

// apiBehaviorUpd docs
// @Summary Updates behavior for specified layout
// @Description Updates behavior for specified layout
// @Description example upd behavior, layout_id must exists in the database:
// @Description {
// @Description 	"Open": "10:00:00",
// @Description 	"Close": "23:00:00",
// @Description 	"time_zone": "Europe/Moscow",
// @Description 	"behavior_config": {
// @Description 	  "queue": {
// @Description 		"layouts": [
// @Description 		  {
// @Description 			"layout_id": "89884907",
// @Description 			"title": "retail name",
// @Description 			"threshold": 99.99,
// @Description 			"service_channel": {
// @Description 			  "indexes": [
// @Description 				{
// @Description 				  "kind": "work_time",
// @Description 				  "weight": 1,
// @Description 				  "op": "*"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "queue_length",
// @Description 				  "weight": 40,
// @Description 				  "op": "+"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "cashiers_activities",
// @Description 				  "weight": 100,
// @Description 				  "op": "+"
// @Description 				}
// @Description 			  ]
// @Description 			}
// @Description 		  }
// @Description 		],
// @Description 		"stores": [
// @Description 		  {
// @Description 			"store_id": null,
// @Description 			"title": "store name",
// @Description 			"threshold": 99.99,
// @Description 			"service_channel": {
// @Description 			  "indexes": [
// @Description 				{
// @Description 				  "kind": "work_time",
// @Description 				  "weight": 1,
// @Description 				  "op": "*"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "queue_length",
// @Description 				  "weight": 40,
// @Description 				  "op": "+"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "cashiers_activities",
// @Description 				  "weight": 90,
// @Description 				  "op": "+"
// @Description 				}
// @Description 			  ]
// @Description 			}
// @Description 		  }
// @Description 		],
// @Description 		"service_channels": [
// @Description 		  {
// @Description 			"service_channel_id": "131928214",
// @Description 			"title": "Касса №10",
// @Description 			"threshold": 99.99,
// @Description 			"service_channel": {
// @Description 			  "indexes": [
// @Description 				{
// @Description 				  "kind": "work_time",
// @Description 				  "weight": 1,
// @Description 				  "op": "*"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "queue_length",
// @Description 				  "weight": 100,
// @Description 				  "op": "+"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "cashiers_activities",
// @Description 				  "weight": 100,
// @Description 				  "op": "+"
// @Description 				}
// @Description 			  ]
// @Description 			}
// @Description 		  },
// @Description 		  {
// @Description 			"service_channel_id": "66932267",
// @Description 			"title": "Касса №2",
// @Description 			"threshold": 99.99,
// @Description 			"service_channel": {
// @Description 			  "indexes": [
// @Description 				{
// @Description 				  "kind": "work_time",
// @Description 				  "weight": 1,
// @Description 				  "op": "*"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "queue_length",
// @Description 				  "weight": 100,
// @Description 				  "op": "+"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "cashiers_activities",
// @Description 				  "weight": 100,
// @Description 				  "op": "+"
// @Description 				}
// @Description 			  ]
// @Description 			}
// @Description 		  },
// @Description 		  {
// @Description 			"service_channel_id": "73500176",
// @Description 			"title": "Касса №5",
// @Description 			"threshold": 99.99,
// @Description 			"service_channel": {
// @Description 			  "indexes": [
// @Description 				{
// @Description 				  "kind": "work_time",
// @Description 				  "weight": 1,
// @Description 				  "op": "*"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "queue_length",
// @Description 				  "weight": 100,
// @Description 				  "op": "+"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "cashiers_activities",
// @Description 				  "weight": 100,
// @Description 				  "op": "+"
// @Description 				}
// @Description 			  ]
// @Description 			}
// @Description 		  },
// @Description 		  {
// @Description 			"service_channel_id": "80541959",
// @Description 			"title": "Касса №7",
// @Description 			"threshold": 99.99,
// @Description 			"service_channel": {
// @Description 			  "indexes": [
// @Description 				{
// @Description 				  "kind": "work_time",
// @Description 				  "weight": 1,
// @Description 				  "op": "*"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "queue_length",
// @Description 				  "weight": 100,
// @Description 				  "op": "+"
// @Description 				},
// @Description 				{
// @Description 				  "kind": "cashiers_activities",
// @Description 				  "weight": 100,
// @Description 				  "op": "+"
// @Description 				}
// @Description 			  ]
// @Description 			}
// @Description 		  }
// @Description 		]
// @Description 	  },
// @Description 	  "recommendations": {
// @Description 		"layouts": [
// @Description 		  {
// @Description 			"layout_id": "123124324",
// @Description 			"title": "retail name",
// @Description 			"std_coef": 0.3,
// @Description 			"queue_multiplier": 0.6,
// @Description 			"pred_minutes": 10,
// @Description 			"hist_minutes": 30,
// @Description 			"checkout_productivity": 0.8
// @Description 		  }
// @Description 		],
// @Description 		"stores": [
// @Description 		  {
// @Description 			"store_id": "234534545",
// @Description 			"title": "store name",
// @Description 			"std_coef": 0.3,
// @Description 			"queue_multiplier": 0.6,
// @Description 			"pred_minutes": 10,
// @Description 			"hist_minutes": 30,
// @Description 			"checkout_productivity": 0.8
// @Description 		  },
// @Description 		  {
// @Description 			"store_id": "234556487545",
// @Description 			"title": "store name another",
// @Description 			"std_coef": 0.3,
// @Description 			"queue_multiplier": 0.6,
// @Description 			"pred_minutes": 10,
// @Description 			"hist_minutes": 30,
// @Description 			"checkout_productivity": 0.8
// @Description 		  }
// @Description 		]
// @Description 	  },
// @Description 	  "queue_thresholds": {
// @Description 		"layouts": [
// @Description 		  {
// @Description 			"layout_id": "89884907",
// @Description 			"title": "retail name",
// @Description 			"threshold": 3,
// @Description 			"sequence_length": 2
// @Description 		  }
// @Description 		],
// @Description 		"stores": [
// @Description 		  {
// @Description 			"store_id": "242342343",
// @Description 			"title": "store name",
// @Description 			"threshold": 2.99,
// @Description 			"sequence_length": 2
// @Description 		  }
// @Description 		],
// @Description 		"blocks_service_channels": [
// @Description 		  {
// @Description 			"block_service_chanels_id": "2343245",
// @Description 			"title": "cash block name",
// @Description 			"threshold": 2.89,
// @Description 			"sequence_length": 2
// @Description 		  }
// @Description 		]
// @Description 	  }
// @Description 	}
// @Description   }
// @Accept  json
// @Produce  json
// @Tags common
// @Param layout_id path string true "digit/uuid format"
// @Param behavior body domain.Behavior true "behavior properties"
// @Success 200 {object} infra.SuccessResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/layouts/{layout_id}/behaviors [put]
func (s *Server) apiBehaviorUpd(c echo.Context) error {
	bhv := &domain.Behavior{}
	if err := c.Bind(bhv); err != nil {
		s.log.Errorf("apiBehaviorUpd, bad request error %s", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	repo, id, err := s.getRepo(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	if id == "" {
		s.log.Errorf("bad request apiBehaviorUpd, %s", errEmptyID)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(errEmptyID))
	}
	cu, err := repo.UpdBehavior(c.Request().Context(), id, *bhv)
	if err != nil {
		s.log.Errorf("UpdBehavior failed %v", err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(errRepo))
	}
	return c.JSON(http.StatusCreated, OkStatus(`updated `+strconv.FormatInt(cu, 10)))
}
