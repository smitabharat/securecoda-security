package scheduler

import (
	"context"
	"time"

	"securecoda/internal/detector"

	"github.com/rs/zerolog/log"
)

type Scheduler struct {
	det      *detector.Detector
	interval time.Duration
}

func NewScheduler(det *detector.Detector, interval time.Duration, logger interface{}) *Scheduler {
	return &Scheduler{det: det, interval: interval}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Info().Msg("running detection cycle")
			if err := s.det.RunCycle(); err != nil {
				log.Error().Err(err).Msg("detector failed")
			}
		}
	}
}
