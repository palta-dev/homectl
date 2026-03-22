package checks

import (
	"context"
	"log"
	"time"

	"github.com/palta-dev/homectl/apps/server/internal/config"
	"github.com/palta-dev/homectl/apps/server/internal/storage"
)

// Scheduler manages periodic health checks
type Scheduler struct {
	cfg      *config.Config
	executor *Executor
	db       *storage.DB
	interval time.Duration
	stopChan chan struct{}
}

// NewScheduler creates a new health check scheduler
func NewScheduler(cfg *config.Config, executor *Executor, db *storage.DB) *Scheduler {
	return &Scheduler{
		cfg:      cfg,
		executor: executor,
		db:       db,
		interval: 1 * time.Minute, // Default check interval
		stopChan: make(chan struct{}),
	}
}

// Start begins periodic checks in the background
func (s *Scheduler) Start(ctx context.Context) {
	log.Printf("Starting background health check worker (interval: %v)", s.interval)
	
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run first check immediately
	s.runAll(ctx)

	for {
		select {
		case <-ticker.C:
			s.runAll(ctx)
		case <-s.stopChan:
			log.Println("Stopping health check worker")
			return
		case <-ctx.Done():
			return
		}
	}
}

// Stop halts the scheduler
func (s *Scheduler) Stop() {
	close(s.stopChan)
}

// UpdateConfig updates the configuration used by the scheduler
func (s *Scheduler) UpdateConfig(newCfg *config.Config) {
	s.cfg = newCfg
}

func (s *Scheduler) runAll(ctx context.Context) {
	for _, group := range s.cfg.Groups {
		for _, service := range group.Services {
			// Skip services without checks
			if len(service.Checks) == 0 {
				continue
			}

			// For simplicity in the scheduler, we'll just run the first check for now.
			// In a more advanced implementation, we could run all and aggregate.
			check := service.Checks[0]
			
			// Use service name as serviceID for storage
			serviceID := service.Name

			go func(svcID string, c config.Check) {
				// Run check with a timeout
				checkCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()

				result, err := s.executor.Execute(checkCtx, c)
				if err != nil {
					log.Printf("[SCHEDULER ERROR] Failed to check %s: %v", svcID, err)
					return
				}

				// Record in SQLite if DB is available
				if s.db != nil {
					// Record check result for stats/sparklines
					err = s.db.RecordCheckResult(storage.CheckResult{
						ServiceID: svcID,
						Timestamp: result.Timestamp,
						State:     result.State,
						LatencyMs: result.LatencyMs,
						Error:     result.Error,
					})
					if err != nil {
						log.Printf("[SCHEDULER ERROR] Failed to record check for %s: %v", svcID, err)
					}

					// Record/Update incident status
					err = s.db.RecordIncident(svcID, result.State, result.Error)
					if err != nil {
						log.Printf("[SCHEDULER ERROR] Failed to record incident for %s: %v", svcID, err)
					}
				}
			}(serviceID, check)
		}
	}
}
