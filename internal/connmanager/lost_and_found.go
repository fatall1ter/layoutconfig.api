package connmanager

import (
	"context"
	"time"

	"git.countmax.ru/countmax/layoutconfig.api/repos"
	"git.countmax.ru/countmax/pkg/logging"
)

func (m *Manager) addLF(ldb lostDB) {
	(*m.lostAndFound) = append((*m.lostAndFound), ldb)
}

func (m *Manager) delLF(i int) {
	copy((*m.lostAndFound)[i:], (*m.lostAndFound)[i+1:])
}

func (m *Manager) retryFill(ctx context.Context) {
	log := logging.FromContext(ctx)
	for i, lf := range *m.lostAndFound {
		cmr, err := repos.NewLayoutRepo(ctx, lf.scope, lf.cs, lf.timeout)
		if err != nil {
			log.Errorf("repos.NewLayoutRepo for %s failed: %s", getSrvPortDB(lf.cs), err)
		}
		err = m.RegisterRepo(cmr)
		if err != nil {
			log.Errorf("lostAndFound retryFill, registerRepo failed: %s", err)
			continue
		}
		m.delLF(i)
		m.retryFill(ctx)
	}
}

func (m *Manager) lostRetrier(ctx context.Context, period time.Duration) {
	log := logging.FromContext(ctx)
	log.Debugf("start lostRetrier with period %v", period)
	defer log.Debug("stop lostRetrier")
	t := time.NewTicker(period)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			log.Debug("time to retryFill")
			m.retryFill(ctx)
		}
	}
}
