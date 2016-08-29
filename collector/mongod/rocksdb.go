package collector_mongod

import(
	"strings"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// byte size constants:
	kilobyte float64 = 1024
	megabyte float64 = kilobyte * 1024
	gigabyte float64 = megabyte * 1024
	terabyte float64 = gigabyte * 1024
	petabyte float64 = terabyte * 1024

	rocksDbWriteOps = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"writes_total",
		Help:		"The total number of write operations in RocksDB",
	})
	rocksDbWriteKeys = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"writes_key_total",
		Help:		"The total number of key write operations in RocksDB",
	})
	rocksDbWriteBatches = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"write_batches_total",
		Help:		"The total number of write batches in RocksDB",
	})
	rocksDbWriteBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"write_bytes_total",
		Help:		"The total number of data written by RocksDB",
	})
	rocksDbFlushedBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"flushed_bytes_total",
		Help:		"The total number of flushed bytes in RocksDB",
	})
	rocksDbWALOperations = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"write_ahead_log_operations_total",
		Help:		"The total number of Write-Ahead-Log operations in RocksDB",
	}, []string{"type"})
	rocksDbWALBytes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"write_ahead_log_bytes_total",
		Help:		"The total number of Write-Ahead-Log syncs in RocksDB",
	})
	rocksDbStalledSecs = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"stalled_seconds_total",
		Help:		"The total number of seconds RocksDB has spent stalled",
	})
	rocksDbStalls = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"stalls_total",
		Help:		"The total number of stalls in RocksDB",
	}, []string{"type"})
	rocksDbCompactionBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"compaction_bytes_total",
		Help:		"Total bytes processed during compaction between levels N and N+1 in RocksDB",
	}, []string{"level", "type"})
	rocksDbCompactionSecondsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"compaction_seconds_total",
		Help:		"The time spent doing compactions between levels N and N+1 in RocksDB",
	}, []string{"level"})
	rocksDbCompactionsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"compactions_total",
		Help:		"The total number of compactions between levels N and N+1 in RocksDB",
	}, []string{"level"})
	rocksDbCompactionKeys = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"compaction_keys_total",
		Help:		"The number of keys compared during compactions in RocksDB",
	}, []string{"level", "type"})
)

var (
	rocksDbNumImmutableMemTable = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"immutable_memtables",
		Help:		"The total number of immutable MemTables in RocksDB",
	})
	rocksDbMemTableFlushPending = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"pending_memtable_flushes",
		Help:		"The total number of MemTable flushes pending in RocksDB",
	}) 
	rocksDbCompactionPending = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"pending_compactions",
		Help:		"The total number of compactions pending in RocksDB",
	}) 
	rocksDbBackgroundErrors = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"background_errors",
		Help:		"The total number of background errors in RocksDB",
	}) 
	rocksDbMemTableBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"memtable_bytes",
		Help:		"The current number of MemTable bytes in RocksDB",
	}, []string{"type"}) 
	rocksDbNumEntriesMemTableActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"memtable_active_entries",
		Help:		"The current number of cctive MemTable entries in RocksDB",
	}) 
	rocksDbNumEntriesImmMemTable = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"immutable_memtable_entries",
		Help:		"The current number of immutable MemTable entries in RocksDB",
	}) 
	rocksDbEstimateTableReadersMem = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"estimate_table_readers_memory_bytes",
		Help:		"The estimate RocksDB table-reader memory bytes",
	}) 
	rocksDbNumSnapshots = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"snapshots",
		Help:		"The current number of snapshots in RocksDB",
	}) 
	rocksDbOldestSnapshotTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"oldest_snapshot_timestamp",
		Help:		"The timestamp of the oldest snapshot in RocksDB",
	}) 
	rocksDbNumLiveVersions = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"live_versions",
		Help:		"The current number of live versions in RocksDB",
	}) 
	rocksDbTotalLiveRecoveryUnits = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"total_live_recovery_units",
		Help:		"The total number of live recovery units in RocksDB",
	}) 
	rocksDbBlockCacheUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"block_cache_bytes",
		Help:		"The bytes used in the RocksDB Block Cache",
	}) 
	rocksDbTransactionEngineKeys = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"transaction_engine_keys",
		Help:		"The current number of transaction engine keys in RocksDB",
	}) 
	rocksDbTransactionEngineSnapshots = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"transaction_engine_snapshots",
		Help:		"The current number of transaction engine snapshots in RocksDB",
	}) 
	rocksDbWritesPerBatch = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"writes_per_batch",
		Help:		"The number of writes per batch in RocksDB",
	}) 
	rocksDbWritesPerSec = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"writes_per_second",
		Help:		"The number of writes per second in RocksDB",
	}) 
	rocksDbStallPercent = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"stall_percent",
		Help:		"The percentage of time RocksDB has been stalled",
	}) 
	rocksDbWALWritesPerSync = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"write_ahead_log_writes_per_sync",
		Help:		"The number of writes per Write-Ahead-Log sync in RocksDB",
	}) 
	rocksDbWALBytesPerSecs = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"write_ahead_log_bytes_per_second",
		Help:		"The number of bytes written per second by the Write-Ahead-Log in RocksDB",
	}) 
	rocksDbNumLevels = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"num_levels",
		Help:		"The number of compaction levels in RocksDB",
	})
	rocksDbLevelNumFiles = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"num_files",
		Help:		"The number of files in a RocksDB level",
	}, []string{"level"})
	rocksDbCompactionThreads = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"compaction_file_threads",
		Help:		"The number of threads currently doing compaction for levels in RocksDB",
	}, []string{"level"})
	rocksDbLevelScore = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"compaction_score",
		Help:		"The compaction score of RocksDB levels",
	}, []string{"level"})
	rocksDbLevelSizeBytes = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"size_bytes",
		Help:		"The total byte size of levels in RocksDB",
	}, []string{"level"})
	rocksDbCompactionBytesPerSec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"compaction_bytes_per_second",
		Help:		"The rate at which data is processed during compaction between levels N and N+1 in RocksDB",
	}, []string{"level", "type"})
	rocksDbCompactionWriteAmplification = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"compaction_write_amplification",
		Help:		"The write amplification factor from compaction between levels N and N+1 in RocksDB",
	}, []string{"level"})
	rocksDbCompactionAvgSeconds = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:	Namespace,
		Subsystem:	"rocksdb",
		Name:		"compaction_avg_seconds",
		Help:		"The average time per compaction between levels N and N+1 in RocksDB",
	}, []string{"level"})
)

type RocksDbStats struct {
	NumImmutableMemTable		string		`bson:"num-immutable-mem-table"`
	MemTableFlushPending		string		`bson:"mem-table-flush-pending"`
	CompactionPending		string		`bson:"compaction-pending"`
	BackgroundErrors		string		`bson:"background-errors"`
	CurSizeMemTableActive		string		`bson:"cur-size-active-mem-table"`
	CurSizeAllMemTables		string		`bson:"cur-size-all-mem-tables"`
	NumEntriesMemTableActive	string		`bson:"num-entries-active-mem-table"`
	NumEntriesImmMemTables		string		`bson:"num-entries-imm-mem-tables"`
	EstimateTableReadersMem		string		`bson:"estimate-table-readers-mem"`
	NumSnapshots			string		`bson:"num-snapshots"`
	OldestSnapshotTime		string		`bson:"oldest-snapshot-time"`
	NumLiveVersions			string		`bson:"num-live-versions"`
	BlockCacheUsage			string		`bson:"block-cache-usage"`
	TotalLiveRecoveryUnits		float64		`bson:"total-live-recovery-units"`
	TransactionEngineKeys		float64		`bson:"transaction-engine-keys"`
	TransactionEngineSnapshots	float64		`bson:"transaction-engine-snapshots"`
	Stats				[]string	`bson:"stats"`
	ThreadStatus			[]string	`bson:"thread-status"`
}

type RocksDbLevelStatsFiles struct {
	Num		float64
	CompThreads	float64
}

type RocksDbLevelStats struct {
	Level		string
	Files		*RocksDbLevelStatsFiles
	Score		float64
	SizeMB		float64
	ReadGB		float64
	RnGB		float64
	Rnp1GB		float64
	WriteGB		float64
	WnewGB		float64
	MovedGB		float64
	WAmp		float64
	RdMBPSec	float64
	WrMBPSec	float64
	CompSec		float64
	CompCnt		float64
	AvgSec		float64
	KeyIn		float64
	KeyDrop		float64
}

// rocksdb time-format string parser: returns float64 of seconds:
func ParseTime(str string) float64 {
	time_str := strings.Split(str, " ")[0]
	time_split := strings.Split(time_str, ":")
	seconds_hour, err := strconv.ParseFloat(time_split[0], 64)
	seconds_min, err := strconv.ParseFloat(time_split[1], 64)
	seconds, err := strconv.ParseFloat(time_split[2], 64)
	if err != nil {
		return float64(-1)
	}
	return (seconds_hour * 3600) + (seconds_min * 60) + seconds
}

// rocksdb metric string parser: converts string-numbers to float64s and parses metric units (MB, KB, etc):
func ParseStr(str string) float64 {
	var multiply float64 = 1
	var str_remove string = ""
	if strings.Contains(str, " KB") || strings.HasSuffix(str, "KB") {
		multiply = kilobyte 
		str_remove = "KB"
	} else if strings.Contains(str, " MB") || strings.HasSuffix(str, "MB") {
		multiply = megabyte
		str_remove = "MB"
	} else if strings.Contains(str, " GB") || strings.HasSuffix(str, "GB") {
		multiply = gigabyte
		str_remove = "GB"
	} else if strings.Contains(str, " TB") || strings.HasSuffix(str, "TB") {
		multiply = terabyte
		str_remove = "TB"
	} else if strings.Contains(str, " PB") || strings.HasSuffix(str, "PB") {
		multiply = petabyte
		str_remove = "PB"
	} else if strings.Contains(str, "K") {
		first_field := strings.Split(str, " ")[0]
		if strings.HasSuffix(first_field, "K") {
			multiply = 1000
			str_remove = "K"
		}
	} else if strings.HasSuffix(str, "B") {
		str_remove = "B"
	} else if strings.HasSuffix(str, "H:M:S") {
		return ParseTime(str)
	}

	if str_remove != "" {
		str = strings.Replace(str, str_remove, "", 1)
	}

	// use the first thing that is a parseable number:
	for _, word := range strings.Split(str, " ") {
		float, err := strconv.ParseFloat(word, 64)
		if err == nil {
			return float * multiply
		}
	}
	return float64(-1)
}

// splits strings with multi-whitespace delimeters into a slice:
func SplitByWs(str string) []string {
	var fields []string
	for _, field := range strings.Split(str, " ") {
		if field != "" {
			fields = append(fields, field)
		}
	}
	return fields
}

func ProcessLevelStatsLineFiles(str string) *RocksDbLevelStatsFiles {
	split := strings.Split(str, "/")
	numFiles, err := strconv.ParseFloat(split[0], 64)
	compThreads, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return &RocksDbLevelStatsFiles{}
	}
	return &RocksDbLevelStatsFiles{
		Num: numFiles,
		CompThreads: compThreads,
	}
}

func ProcessLevelStatsLine(line string) *RocksDbLevelStats {
	var stats *RocksDbLevelStats
	if strings.HasPrefix(line, " ") {
		fields := SplitByWs(line)
		stats = &RocksDbLevelStats{
			Level: fields[0],
			Files: ProcessLevelStatsLineFiles(fields[1]),
			SizeMB: ParseStr(fields[2]),
			Score: ParseStr(fields[3]),
			ReadGB: ParseStr(fields[4]),
			RnGB: ParseStr(fields[5]),
			Rnp1GB: ParseStr(fields[6]),
			WriteGB: ParseStr(fields[7]),
			WnewGB: ParseStr(fields[8]),
			MovedGB: ParseStr(fields[9]),
			WAmp: ParseStr(fields[10]),
			RdMBPSec: ParseStr(fields[11]),
			WrMBPSec: ParseStr(fields[12]),
			CompSec: ParseStr(fields[13]),
			CompCnt: ParseStr(fields[14]),
			AvgSec: ParseStr(fields[15]),
			KeyIn: ParseStr(fields[16]),
			KeyDrop: ParseStr(fields[17]),
		}
	}
	return stats
}

func (stats *RocksDbStats) GetStatsSection(section_prefix string) []string {
	var lines []string
	var is_section bool
	for _, line := range stats.Stats {
		if is_section {
			if strings.HasPrefix(line, "** ") && strings.HasSuffix(line, " **") {
				break
			} else if line != "" {
				lines = append(lines, line)
			}
		} else if strings.HasPrefix(line, section_prefix) {
			is_section = true
		}
	}
	return lines
}

func (stats *RocksDbStats) GetStatsLine(section_prefix string, line_prefix string) []string {
	var fields []string
	for _, line := range stats.GetStatsSection(section_prefix) {
		if strings.HasPrefix(line, line_prefix) {
			line = strings.Replace(line, line_prefix, "", 1)
			fields = strings.Split(line, ", ")
		}
	}
	return fields
}

func (stats *RocksDbStats) GetStatsLineField(section_prefix string, line_prefix string, idx int) float64 {
	var field float64 = -1
	stats_line := stats.GetStatsLine(section_prefix, line_prefix)
	if len(stats_line) > idx {
		field = ParseStr(stats_line[idx])
	}
	return field
}

func (stats *RocksDbStats) ProcessLevelStats() {
	var maxLvl float64 = 0
	var levels []*RocksDbLevelStats
	var is_section bool
	for _, line := range stats.Stats {
		if is_section {
			if strings.HasPrefix(line, " Int") {
				break
			} else if line != "" {
				levels = append(levels, ProcessLevelStatsLine(line))
			}
		} else if strings.HasPrefix(line, "------") {
			is_section = true
		}
	}
	for _, level := range levels {
		lvlNum := ParseStr(strings.Replace(level.Level, "L", "", 1))
		if lvlNum > maxLvl {
			maxLvl = lvlNum
		}
		if level.Level != "L0" {
			rocksDbCompactionBytes.With(prometheus.Labels{"level": level.Level, "type": "read"}).Set(level.ReadGB * gigabyte)
			rocksDbCompactionBytes.With(prometheus.Labels{"level": level.Level, "type": "read_n"}).Set(level.RnGB * gigabyte)
			rocksDbCompactionBytes.With(prometheus.Labels{"level": level.Level, "type": "read_np1"}).Set(level.Rnp1GB * gigabyte)
			rocksDbCompactionBytes.With(prometheus.Labels{"level": level.Level, "type": "moved"}).Set(level.MovedGB * gigabyte)
			rocksDbCompactionBytesPerSec.With(prometheus.Labels{"level": level.Level, "type": "read"}).Set(level.RdMBPSec * megabyte)
			rocksDbCompactionWriteAmplification.WithLabelValues(level.Level).Set(level.WAmp)
			rocksDbCompactionKeys.With(prometheus.Labels{"level": level.Level, "type": "in"}).Set(level.KeyIn)
			rocksDbCompactionKeys.With(prometheus.Labels{"level": level.Level, "type": "drop"}).Set(level.KeyDrop)
		}
		rocksDbLevelScore.WithLabelValues(level.Level).Set(level.Score)
		rocksDbLevelNumFiles.WithLabelValues(level.Level).Set(level.Files.Num)
		rocksDbCompactionThreads.WithLabelValues(level.Level).Set(level.Files.CompThreads)
		rocksDbLevelSizeBytes.WithLabelValues(level.Level).Set(level.SizeMB * megabyte)
		rocksDbCompactionSecondsTotal.WithLabelValues(level.Level).Set(level.CompSec)
		rocksDbCompactionAvgSeconds.WithLabelValues(level.Level).Set(level.AvgSec)
		rocksDbCompactionBytes.With(prometheus.Labels{"level": level.Level, "type": "write"}).Set(level.WriteGB * gigabyte)
		rocksDbCompactionBytes.With(prometheus.Labels{"level": level.Level, "type": "write_new_np1"}).Set(level.WriteGB * gigabyte)
		rocksDbCompactionBytesPerSec.With(prometheus.Labels{"level": level.Level, "type": "write"}).Set(level.WrMBPSec * megabyte)
		rocksDbCompactionsTotal.WithLabelValues(level.Level).Set(level.CompCnt)
	}
	rocksDbNumLevels.Set(maxLvl + 1)
}

func (stats *RocksDbStats) ProcessStalls() {
	for _, stall_line := range stats.GetStatsLine("** Compaction Stats [default] **", "Stalls(count): ") {
		stall_split := strings.Split(stall_line, " ")
		if len(stall_split) == 2 {
			stall_type := stall_split[1]
			stall_count := stall_split[0]
			rocksDbStalls.WithLabelValues(stall_type).Set(ParseStr(stall_count))
		}
	}
}

func (stats *RocksDbStats) Describe(ch chan<- *prometheus.Desc) {
	rocksDbWriteOps.Describe(ch)
	rocksDbWriteKeys.Describe(ch)
	rocksDbWriteBatches.Describe(ch)
	rocksDbWriteBytes.Describe(ch)
	rocksDbWritesPerBatch.Describe(ch)
	rocksDbWritesPerSec.Describe(ch)
	rocksDbFlushedBytes.Describe(ch)
	rocksDbWALOperations.Describe(ch)
	rocksDbWALBytes.Describe(ch)
	rocksDbWALBytesPerSecs.Describe(ch)
	rocksDbWALWritesPerSync.Describe(ch)
	rocksDbStallPercent.Describe(ch)
	rocksDbStalledSecs.Describe(ch)
	rocksDbNumLevels.Describe(ch)
	rocksDbLevelNumFiles.Describe(ch)
	rocksDbCompactionThreads.Describe(ch)
	rocksDbLevelSizeBytes.Describe(ch)
	rocksDbLevelScore.Describe(ch)
	rocksDbCompactionBytes.Describe(ch)
	rocksDbCompactionBytesPerSec.Describe(ch)
	rocksDbCompactionWriteAmplification.Describe(ch)
	rocksDbCompactionSecondsTotal.Describe(ch)
	rocksDbCompactionAvgSeconds.Describe(ch)
	rocksDbCompactionsTotal.Describe(ch)
	rocksDbCompactionKeys.Describe(ch)
	rocksDbNumImmutableMemTable.Describe(ch)
	rocksDbMemTableFlushPending.Describe(ch)
	rocksDbCompactionPending.Describe(ch)
	rocksDbBackgroundErrors.Describe(ch)
	rocksDbMemTableBytes.Describe(ch)
	rocksDbNumEntriesMemTableActive.Describe(ch)
	rocksDbNumEntriesImmMemTable.Describe(ch)
	rocksDbEstimateTableReadersMem.Describe(ch)
	rocksDbNumSnapshots.Describe(ch)
	rocksDbOldestSnapshotTimestamp.Describe(ch)
	rocksDbNumLiveVersions.Describe(ch)
	rocksDbBlockCacheUsage.Describe(ch)
	rocksDbTotalLiveRecoveryUnits.Describe(ch)
	rocksDbTransactionEngineKeys.Describe(ch)
	rocksDbTransactionEngineSnapshots.Describe(ch)
}

func (stats *RocksDbStats) Export(ch chan<- prometheus.Metric) {
	// cumulative stats from db.serverStatus().rocksdb.stats (parsed):
	rocksDbFlushedBytes.Set(stats.GetStatsLineField("** Compaction Stats [default] **", "Flush(GB): ", 0) * gigabyte)
	rocksDbWriteOps.Set(stats.GetStatsLineField("** DB Stats **", "Cumulative writes: ", 0))
	rocksDbWriteKeys.Set(stats.GetStatsLineField("** DB Stats **", "Cumulative writes: ", 1))
	rocksDbWriteBatches.Set(stats.GetStatsLineField("** DB Stats **", "Cumulative writes: ", 2))
	rocksDbWritesPerBatch.Set(stats.GetStatsLineField("** DB Stats **", "Cumulative writes: ", 3))
	rocksDbWriteBytes.Set(stats.GetStatsLineField("** DB Stats **", "Cumulative writes: ", 4))
	rocksDbWritesPerSec.Set(stats.GetStatsLineField("** DB Stats **", "Cumulative writes: ", 5))
	rocksDbWALOperations.WithLabelValues("write").Set(stats.GetStatsLineField("** DB Stats **", "Cumulative WAL: ", 0))
	rocksDbWALOperations.WithLabelValues("sync").Set(stats.GetStatsLineField("** DB Stats **", "Cumulative WAL: ", 1))
	rocksDbWALBytes.Set(stats.GetStatsLineField("** DB Stats **", "Cumulative WAL: ", 3))
	rocksDbWALBytesPerSecs.Set(stats.GetStatsLineField("** DB Stats **", "Cumulative WAL: ", 4))
	rocksDbWALWritesPerSync.Set(stats.GetStatsLineField("** DB Stats **", "Cumulative WAL: ", 2))
	rocksDbStalledSecs.Set(stats.GetStatsLineField("** DB Stats **", "Cumulative stall: ", 0))
	rocksDbStallPercent.Set(stats.GetStatsLineField("** DB Stats **", "Cumulative stall: ", 1))

	// stats from db.serverStatus().rocksdb (parsed):
	rocksDbNumImmutableMemTable.Set(ParseStr(stats.NumImmutableMemTable))
	rocksDbMemTableFlushPending.Set(ParseStr(stats.MemTableFlushPending))
	rocksDbCompactionPending.Set(ParseStr(stats.CompactionPending))
	rocksDbBackgroundErrors.Set(ParseStr(stats.BackgroundErrors))
	rocksDbNumEntriesMemTableActive.Set(ParseStr(stats.NumEntriesMemTableActive))
	rocksDbNumEntriesImmMemTable.Set(ParseStr(stats.NumEntriesImmMemTables))
	rocksDbNumSnapshots.Set(ParseStr(stats.NumSnapshots))
	rocksDbOldestSnapshotTimestamp.Set(ParseStr(stats.OldestSnapshotTime))
	rocksDbNumLiveVersions.Set(ParseStr(stats.NumLiveVersions))
	rocksDbBlockCacheUsage.Set(ParseStr(stats.BlockCacheUsage))
	rocksDbEstimateTableReadersMem.Set(ParseStr(stats.EstimateTableReadersMem))
	rocksDbBlockCacheUsage.Set(ParseStr(stats.BlockCacheUsage))
	rocksDbMemTableBytes.WithLabelValues("active").Set(ParseStr(stats.CurSizeMemTableActive))
	rocksDbMemTableBytes.WithLabelValues("total").Set(ParseStr(stats.CurSizeAllMemTables))

	// stats from db.serverStatus().rocksdb (unparsed - somehow these aren't real types!):
	rocksDbTotalLiveRecoveryUnits.Set(stats.TotalLiveRecoveryUnits)
	rocksDbTransactionEngineKeys.Set(stats.TransactionEngineKeys)
	rocksDbTransactionEngineSnapshots.Set(stats.TransactionEngineSnapshots)

	// process per-level stats in to vectors:
	stats.ProcessLevelStats()

	// process stall counts into a vector:
	stats.ProcessStalls()

	rocksDbWriteOps.Collect(ch)
	rocksDbWriteKeys.Collect(ch)
	rocksDbWriteBatches.Collect(ch)
	rocksDbWriteBytes.Collect(ch)
	rocksDbWritesPerBatch.Collect(ch)
	rocksDbWritesPerSec.Collect(ch)
	rocksDbFlushedBytes.Collect(ch)
	rocksDbWALOperations.Collect(ch)
	rocksDbWALBytes.Collect(ch)
	rocksDbWALBytesPerSecs.Collect(ch)
	rocksDbWALWritesPerSync.Collect(ch)
	rocksDbStallPercent.Collect(ch)
	rocksDbStalledSecs.Collect(ch)
	rocksDbNumLevels.Collect(ch)
	rocksDbLevelNumFiles.Collect(ch)
	rocksDbCompactionThreads.Collect(ch)
	rocksDbLevelSizeBytes.Collect(ch)
	rocksDbLevelScore.Collect(ch)
	rocksDbCompactionBytes.Collect(ch)
	rocksDbCompactionBytesPerSec.Collect(ch)
	rocksDbCompactionWriteAmplification.Collect(ch)
	rocksDbCompactionSecondsTotal.Collect(ch)
	rocksDbCompactionAvgSeconds.Collect(ch)
	rocksDbCompactionsTotal.Collect(ch)
	rocksDbCompactionKeys.Collect(ch)
	rocksDbNumImmutableMemTable.Collect(ch)
	rocksDbMemTableFlushPending.Collect(ch)
	rocksDbCompactionPending.Collect(ch)
	rocksDbBackgroundErrors.Collect(ch)
	rocksDbNumEntriesMemTableActive.Collect(ch)
	rocksDbNumEntriesImmMemTable.Collect(ch)
	rocksDbNumSnapshots.Collect(ch)
	rocksDbOldestSnapshotTimestamp.Collect(ch)
	rocksDbNumLiveVersions.Collect(ch)
	rocksDbTotalLiveRecoveryUnits.Collect(ch)
	rocksDbTransactionEngineKeys.Collect(ch)
	rocksDbTransactionEngineSnapshots.Collect(ch)
	rocksDbMemTableBytes.Collect(ch)
	rocksDbEstimateTableReadersMem.Collect(ch)
	rocksDbBlockCacheUsage.Collect(ch)
	rocksDbStalls.Collect(ch)
}