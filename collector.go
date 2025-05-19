package main

import (
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	restic_check_success           *prometheus.Desc
	restic_locks_total             *prometheus.Desc
	restic_snapshots_total         *prometheus.Desc
	restic_scrape_duration_seconds *prometheus.Desc
	restic_backup_timestamp        *prometheus.Desc
	restic_backup_files_total      *prometheus.Desc
	restic_backup_files_new        *prometheus.Desc
	restic_backup_files_changed    *prometheus.Desc
	restic_backup_size_total       *prometheus.Desc
	restic_backup_runtime          *prometheus.Desc
}

var labels = []string{
	"client_hostname",
	"client_username",
	"client_version",
	"snapshot_id",
	"snapshot_tags",
	"snapshot_paths",
}

// Initialize Metrics
// New metrics also have to be appended to the Collector struct, the Describe and the Collect functions separately
func newCollector() *Collector {
	// Add repo_path to labels, if specified in config
	if Config.USE_REPO_PATH {
		labels = append(labels, "repo_path")
	}

	return &Collector{
		restic_check_success: prometheus.NewDesc("restic_check_success",
			"Shows whether a check was sucessful",
			nil, nil,
		),
		restic_locks_total: prometheus.NewDesc("restic_locks_total",
			"Shows the amount of locks on the repository",
			nil, nil,
		),
		restic_snapshots_total: prometheus.NewDesc("restic_snapshots_total",
			"Shows the total amount of snapshots in the repository",
			nil, nil,
		),
		restic_scrape_duration_seconds: prometheus.NewDesc("restic_scrape_duration_seconds",
			"Shows the duration of the scrape",
			nil, nil,
		),
		restic_backup_timestamp: prometheus.NewDesc("restic_backup_timestamp",
			"Shows the start time of the snapshot",
			labels, nil,
		),
		restic_backup_files_total: prometheus.NewDesc("restic_backup_files_total",
			"Shows the total amount of files in the snapshot",
			labels, nil,
		),
		restic_backup_files_new: prometheus.NewDesc("restic_backup_files_new",
			"Shows the amount of new files in the snapshot",
			labels, nil,
		),
		restic_backup_files_changed: prometheus.NewDesc("restic_backup_files_changed",
			"Shows the amount of changed files in the snapshot",
			labels, nil,
		),
		restic_backup_size_total: prometheus.NewDesc("restic_backup_size_total",
			"Shows the amount of bytes in the snapshot",
			labels, nil,
		),
		restic_backup_runtime: prometheus.NewDesc("restic_backup_runtime",
			"Shows the time the snapshot took",
			labels, nil,
		),
	}
}

// Pass all Descriptions to Prometheus
func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {

	ch <- collector.restic_check_success
	ch <- collector.restic_locks_total
	ch <- collector.restic_snapshots_total
	ch <- collector.restic_scrape_duration_seconds
	ch <- collector.restic_backup_timestamp
	ch <- collector.restic_backup_files_total
	ch <- collector.restic_backup_files_new
	ch <- collector.restic_backup_files_changed
	ch <- collector.restic_backup_size_total
	ch <- collector.restic_backup_runtime
}

// Pass Metric values to prometheus
func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	startTime := time.Now()

	snapshots, err := getSnapshots()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewInvalidDesc(err), err)
	}

	snapshot_count, err := getSnapshotCount()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewInvalidDesc(err), err)
	}

	// Per snapshot metrics
	for _, snapshot := range snapshots {
		labelValues := []string{
			snapshot.Hostname,
			snapshot.Username,
			snapshot.ProgramVersion,
			snapshot.ID,
			strings.Join(snapshot.Tags, ","),
			strings.Join(snapshot.Paths, ","),
		}

		if Config.USE_REPO_PATH {
			labelValues = append(labelValues, os.Getenv("RESTIC_REPOSITORY"))
		}

		ch <- prometheus.MustNewConstMetric(collector.restic_backup_timestamp, prometheus.GaugeValue, float64(snapshot.Summary.BackupStart.Unix()), labelValues...)
		ch <- prometheus.MustNewConstMetric(collector.restic_backup_files_total, prometheus.GaugeValue, float64(snapshot.Summary.TotalFilesProcessed), labelValues...)
		ch <- prometheus.MustNewConstMetric(collector.restic_backup_files_new, prometheus.GaugeValue, float64(snapshot.Summary.FilesNew), labelValues...)
		ch <- prometheus.MustNewConstMetric(collector.restic_backup_files_changed, prometheus.GaugeValue, float64(snapshot.Summary.FilesChanged), labelValues...)
		ch <- prometheus.MustNewConstMetric(collector.restic_backup_size_total, prometheus.GaugeValue, float64(snapshot.Summary.TotalBytesProcessed), labelValues...)
		ch <- prometheus.MustNewConstMetric(collector.restic_backup_runtime, prometheus.GaugeValue, float64(snapshot.Summary.BackupEnd.Unix()-snapshot.Summary.BackupStart.Unix()), labelValues...)
	}

	// Last check's status (Don't run check on every scrape, that would be too much load, just get from the variable, that we set in the ticker)
	ch <- prometheus.MustNewConstMetric(collector.restic_check_success, prometheus.GaugeValue, float64(checkResult))

	// Get locks on the repo, just discard metric, if the call fails
	if locks, err := getLocks(); err == nil {
		ch <- prometheus.MustNewConstMetric(collector.restic_locks_total, prometheus.GaugeValue, float64(locks))
	}

	// Get total amount of snapshots in the repo
	ch <- prometheus.MustNewConstMetric(collector.restic_snapshots_total, prometheus.GaugeValue, float64(snapshot_count))

	// Get how long our metric retrieval took
	duration := time.Since(startTime)
	ch <- prometheus.MustNewConstMetric(collector.restic_scrape_duration_seconds, prometheus.GaugeValue, duration.Seconds())
}
