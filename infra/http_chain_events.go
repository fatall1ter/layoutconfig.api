package infra

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/domain"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// ChainEventsResponse http wrapper with metadata
type ChainEventsResponse struct {
	Data domain.Events `json:"data"`
	Metadata
}

// nolint:lll
// apiChainEvents docs
// @Summary Get all events for the retail schema
// @Description get slice of events with layout_id, store_id, key, kind, severity; from - to datetime range and offset, limit parameters
// @Description key can be: queue.threshold.exceeded, ...
// @Description kind can be: business, system, user
// @Description severity can be: info, warn, alarm
// @Description from/to can be: YYYY-MM-DDTHH:mm:ss+07:00 or naive YYYY-MM-DD HH:mm:ss then the server's local timezone is applied
// @Description fields - comma separated values of field names, can be id,key,event_time,kind,message,severity,layout_id,store_id,source,creator,created_at... all of them described at the model
// @Produce json
// @Tags chains/events
// @Param layout_id query string false "default=*"
// @Param store_id query string false "default=*"
// @Param key query string false "default=*"
// @Param kind query string false "default=*"
// @Param severity query string false "default=*"
// @Param from query string false "ISO8601 datetime, default begin of day"
// @Param to query string false "ISO8601 datetime, dafault current time"
// @Param offset query integer false "default=0"
// @Param limit query integer false "default=20"
// @Param fields query string false "id,key,event_time,kind... default=id,key,event_time,message,layout_id,store_id"
// @Success 200 {object} infra.ChainEventsResponse
// @Failure 400 {object} infra.HTTPError
// @Failure 401 {object} infra.HTTPError
// @Failure 405 {object} infra.HTTPError
// @Failure 404 {object} infra.ErrResponse
// @Failure 500 {object} infra.ErrResponse
// @Router /v2/chains/events [get]
func (s *Server) apiChainEvents(c echo.Context) error { // TODO: remove me to separate API
	offset, limit := s.getPageParams(c)
	layoutID := c.QueryParam("layout_id")
	storeID := c.QueryParam("store_id")
	key := c.QueryParam("key")
	kind := c.QueryParam("kind")
	severity := c.QueryParam("severity")
	from, to, err := s.getFromToParams(c)
	if err != nil {
		s.log.Warnf("getFromToParams error, %s", err)
		return c.JSON(http.StatusBadRequest, ErrInvalidRequest(err))
	}
	events, count, err := s.evRepo.FindChainEvents(from, to, layoutID, storeID, key, kind, severity, limit, offset)
	if err != nil {

		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	if len(events) == 0 {
		response := ChainEventsResponse{
			Data: events,
			Metadata: Metadata{
				ResultSet: ResultSet{
					Count:  0,
					Offset: offset,
					Limit:  limit,
					Total:  count,
				},
			},
		}
		return c.JSON(http.StatusOK, response)
	}
	fields := strings.Split(c.QueryParam("fields"), ",")
	if c.QueryParam("fields") == "" {
		fields = []string{"id", "key", "event_time", "message", "layout_id", "store_id"}
	}
	events.SetZeroValue(fields)
	response := ChainEventsResponse{
		Data: events,
		Metadata: Metadata{
			ResultSet: ResultSet{
				Count:  int64(len(events)),
				Offset: offset,
				Limit:  limit,
				Total:  count,
			},
		},
	}
	return c.JSON(http.StatusOK, response)
}

func (s *Server) serveChainEventsWS(c echo.Context) error {
	subscriber := c.Request().RemoteAddr
	query := c.Request().URL.Path
	s.log.Debugf("connected new client, remoteAddr: %s, agent: %s", c.Request().RemoteAddr, c.Request().UserAgent())
	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		s.log.Errorf("upgrade connection from %s (%s) failed, %v", c.Request().RemoteAddr, c.Request().Proto, err)
		return c.JSON(http.StatusInternalServerError, ErrServerInternal(err))
	}
	defer ws.Close()

	// get params
	p := &RequestEvent{}
LOOPPARAMS:
	for i := 0; i < 5; i++ {
		t, msg, err := ws.ReadMessage()
		if err != nil {
			s.log.Errorf("read message from ws error, %v, aborted...", err)
			return c.JSON(http.StatusInternalServerError,
				ErrServerInternal(fmt.Errorf("read message from ws error, %v, aborted", err)))
		}
		if t != websocket.TextMessage {
			s.log.Warnf("got not Text message [%d]=%s, need text params", t, string(msg))
			continue LOOPPARAMS
		}

		err = p.UnmarshalJSON(msg)
		if err != nil {
			s.log.Warnf("got wrong format message %s, need RequestEvent format", string(msg))
			// write error to ws
			err := ws.WriteMessage(websocket.TextMessage, []byte("wrong format"))
			if err != nil {
				log.Errorf("ws write error, %v abotring...", err)
				return c.JSON(http.StatusInternalServerError,
					ErrServerInternal(fmt.Errorf("ws write error, %v abotring", err)))
			}
		}
		break LOOPPARAMS
	}
	s.log.Debugf("got parameters: %v", p)
	s.mWS.WithLabelValues(query).Inc()
	defer s.mWS.WithLabelValues(query).Dec()
	// read request and run pinger...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// run writer of events
	t := time.Now()
	if p.From != nil {
		t = *p.From
	}
	chEvents := s.evRepo.FindConsumerChainEvents(subscriber,
		p.LayoutID, p.StoreID, p.Key, p.Kind, p.Severity, t, ctx.Done())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go s.writeToChainEventWS(chEvents, ws, wg, ctx.Done())

	wg.Wait()
	return nil
}

func (s *Server) writeToChainEventWS(in <-chan domain.Event,
	ws *websocket.Conn, wg *sync.WaitGroup, cancel <-chan struct{}) {
	//
	log.Debug("start read event channel and write to ws")
	defer wg.Done()
	defer log.Debug("stop read event channel and write to ws")
	tick := time.NewTicker(5 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-cancel:
			err := ws.WriteMessage(websocket.CloseMessage, []byte("cancelled"))
			if err != nil {
				s.log.Errorf("write ws message error, %s", err)
			}
			return
		case e, more := <-in:
			if !more {
				log.Warn("event channel closed, aborting...")
				return
			}
			err := ws.WriteJSON(e)
			if err != nil {
				log.Errorf("ws error, %v, aborted...", err)
				return
			}
		case t := <-tick.C:
			err := ws.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				log.Errorf("ws error, %v, aborted...", err)
				return
			}
		}

	}
}
