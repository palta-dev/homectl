package storage

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "modernc.org/sqlite"
)

// DB wraps the SQLite connection
type DB struct {
	conn *sql.DB
}

// Incident represents a service incident record
type Incident struct {
	ID        int64     `json:"id"`
	ServiceID string    `json:"serviceId"`
	StartedAt time.Time `json:"startedAt"`
	EndedAt   *time.Time `json:"endedAt,omitempty"`
	State     string    `json:"state"` // down, degraded
	Error     string    `json:"error,omitempty"`
}

// UptimeStats represents uptime statistics
type UptimeStats struct {
	TotalChecks   int     `json:"totalChecks"`
	Successful    int     `json:"successful"`
	Failed        int     `json:"failed"`
	UptimePercent float64 `json:"uptimePercent"`
	AvgLatencyMs  float64 `json:"avgLatencyMs"`
}

// CheckResult represents a single check result for storage
type CheckResult struct {
	ServiceID string    `json:"serviceId"`
	Timestamp time.Time `json:"timestamp"`
	State     string    `json:"state"`
	LatencyMs int64     `json:"latencyMs"`
	Error     string    `json:"error,omitempty"`
}

// New creates a new SQLite database connection
func New(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

// migrate runs database migrations
func (d *DB) migrate() error {
	schema := `
	-- Incidents table
	CREATE TABLE IF NOT EXISTS incidents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service_id TEXT NOT NULL,
		started_at TIMESTAMP NOT NULL,
		ended_at TIMESTAMP,
		state TEXT NOT NULL,
		error TEXT,
		duration_seconds INTEGER
	);

	-- Check results table (for uptime stats)
	CREATE TABLE IF NOT EXISTS check_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service_id TEXT NOT NULL,
		timestamp TIMESTAMP NOT NULL,
		state TEXT NOT NULL,
		latency_ms INTEGER,
		error TEXT
	);

	-- User preferences table (for multi-user support)
	CREATE TABLE IF NOT EXISTS user_preferences (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL UNIQUE,
		layout_prefs TEXT,
		hidden_services TEXT,
		theme TEXT,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Indexes for performance
	CREATE INDEX IF NOT EXISTS idx_incidents_service ON incidents(service_id);
	CREATE INDEX IF NOT EXISTS idx_incidents_started ON incidents(started_at);
	CREATE INDEX IF NOT EXISTS idx_check_results_service ON check_results(service_id);
	CREATE INDEX IF NOT EXISTS idx_check_results_timestamp ON check_results(timestamp);
	`

	_, err := d.conn.Exec(schema)
	return err
}

// RecordIncident records a new incident or updates an existing one
func (d *DB) RecordIncident(serviceID, state, errorMsg string) error {
	// Check for active incident
	var activeIncidentID int64
	err := d.conn.QueryRow(`
		SELECT id FROM incidents 
		WHERE service_id = ? AND ended_at IS NULL 
		ORDER BY started_at DESC LIMIT 1
	`, serviceID).Scan(&activeIncidentID)

	if err == sql.ErrNoRows {
		// No active incident - start a new one if state is bad
		if state == "down" || state == "degraded" {
			_, err = d.conn.Exec(`
				INSERT INTO incidents (service_id, started_at, state, error)
				VALUES (?, ?, ?, ?)
			`, serviceID, time.Now(), state, errorMsg)
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	// Active incident exists
	if state == "up" {
		// Service recovered - close incident
		_, err = d.conn.Exec(`
			UPDATE incidents 
			SET ended_at = ?, duration_seconds = (strftime('%s', 'now') - strftime('%s', started_at))
			WHERE id = ?
		`, time.Now(), activeIncidentID)
		return err
	}

	// Update existing incident error if changed
	if errorMsg != "" {
		_, err = d.conn.Exec(`
			UPDATE incidents SET error = ? WHERE id = ?
		`, errorMsg, activeIncidentID)
	}
	return err
}

// RecordCheckResult stores a check result
func (d *DB) RecordCheckResult(result CheckResult) error {
	_, err := d.conn.Exec(`
		INSERT INTO check_results (service_id, timestamp, state, latency_ms, error)
		VALUES (?, ?, ?, ?, ?)
	`, result.ServiceID, result.Timestamp, result.State, result.LatencyMs, result.Error)
	return err
}

// GetUptimeStats returns uptime statistics for a service
func (d *DB) GetUptimeStats(serviceID string, duration time.Duration) (*UptimeStats, error) {
	since := time.Now().Add(-duration)

	row := d.conn.QueryRow(`
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN state = 'up' THEN 1 ELSE 0 END) as successful,
			SUM(CASE WHEN state != 'up' THEN 1 ELSE 0 END) as failed,
			AVG(CASE WHEN latency_ms > 0 THEN latency_ms ELSE NULL END) as avg_latency
		FROM check_results
		WHERE service_id = ? AND timestamp >= ?
	`, serviceID, since)

	var stats UptimeStats
	var avgLatency sql.NullFloat64

	err := row.Scan(&stats.TotalChecks, &stats.Successful, &stats.Failed, &avgLatency)
	if err != nil {
		return nil, err
	}

	if stats.TotalChecks > 0 {
		stats.UptimePercent = float64(stats.Successful) / float64(stats.TotalChecks) * 100
	}
	if avgLatency.Valid {
		stats.AvgLatencyMs = avgLatency.Float64
	}

	return &stats, nil
}

// GetRecentIncidents returns recent incidents for a service
func (d *DB) GetRecentIncidents(serviceID string, limit int) ([]Incident, error) {
	rows, err := d.conn.Query(`
		SELECT id, service_id, started_at, ended_at, state, error
		FROM incidents
		WHERE service_id = ?
		ORDER BY started_at DESC
		LIMIT ?
	`, serviceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []Incident
	for rows.Next() {
		var inc Incident
		err := rows.Scan(&inc.ID, &inc.ServiceID, &inc.StartedAt, &inc.EndedAt, &inc.State, &inc.Error)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}

	return incidents, rows.Err()
}

// GetActiveIncidents returns all currently active incidents
func (d *DB) GetActiveIncidents() ([]Incident, error) {
	rows, err := d.conn.Query(`
		SELECT id, service_id, started_at, ended_at, state, error
		FROM incidents
		WHERE ended_at IS NULL
		ORDER BY started_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var incidents []Incident
	for rows.Next() {
		var inc Incident
		err := rows.Scan(&inc.ID, &inc.ServiceID, &inc.StartedAt, &inc.EndedAt, &inc.State, &inc.Error)
		if err != nil {
			return nil, err
		}
		incidents = append(incidents, inc)
	}

	return incidents, rows.Err()
}

// SaveUserPreferences stores user preferences
func (d *DB) SaveUserPreferences(userID string, layoutPrefs, hiddenServices, theme string) error {
	layoutJSON, _ := json.Marshal(layoutPrefs)
	hiddenJSON, _ := json.Marshal(hiddenServices)

	_, err := d.conn.Exec(`
		INSERT INTO user_preferences (user_id, layout_prefs, hidden_services, theme, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(user_id) DO UPDATE SET
			layout_prefs = excluded.layout_prefs,
			hidden_services = excluded.hidden_services,
			theme = excluded.theme,
			updated_at = CURRENT_TIMESTAMP
	`, userID, string(layoutJSON), string(hiddenJSON), theme)
	return err
}

// GetUserPreferences retrieves user preferences
func (d *DB) GetUserPreferences(userID string) (layoutPrefs, hiddenServices, theme string, err error) {
	var layoutJSON, hiddenJSON sql.NullString

	err = d.conn.QueryRow(`
		SELECT layout_prefs, hidden_services, theme
		FROM user_preferences
		WHERE user_id = ?
	`, userID).Scan(&layoutJSON, &hiddenJSON, &theme)

	if err == sql.ErrNoRows {
		return "", "", "", nil
	}
	if err != nil {
		return "", "", "", err
	}

	return layoutJSON.String, hiddenJSON.String, theme, nil
}

// CleanupOldData removes old check results to prevent unbounded growth
func (d *DB) CleanupOldData(olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	_, err := d.conn.Exec(`
		DELETE FROM check_results WHERE timestamp < ?
	`, cutoff)
	return err
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.conn.Close()
}
