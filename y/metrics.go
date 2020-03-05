/*
 * Copyright (C) 2017 Dgraph Labs, Inc. and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package y

import "expvar"

var (
	// LSMSize has size of the LSM in bytes
	LSMSize *expvar.Map
	// VlogSize has size of the value log in bytes
	VlogSize *expvar.Map
	// PendingWrites tracks the number of pending writes.
	PendingWrites *expvar.Map

	// These are cumulative

	// NumReads has cumulative number of reads
	NumReads *expvar.Int
	// NumWrites has cumulative number of writes
	NumWrites *expvar.Int
	// NumBytesRead has cumulative number of bytes read
	NumBytesRead *expvar.Int
	// NumBytesWritten has cumulative number of bytes written
	NumBytesWritten *expvar.Int
	// NumLSMGets is number of LMS gets
	NumLSMGets *expvar.Map
	// NumLSMBloomHits is number of LMS bloom hits
	NumLSMBloomHits *expvar.Map
	// NumGets is number of gets
	NumGets *expvar.Int
	// NumPuts is number of puts
	NumPuts *expvar.Int
	// NumBlockedPuts is number of blocked puts
	NumBlockedPuts *expvar.Int
	// NumMemtableGets is number of memtable gets
	NumMemtableGets *expvar.Int
)

// Temporary Lightstep extras.
type lightstepMetrics struct {
	numLSMReads             *expvar.Int
	numLSMLogicalBytesRead  *expvar.Int
	numLSMPhysicalBytesRead *expvar.Int // Assuming 4KB block size.

	numVLogReads             *expvar.Int
	numVLogLogicalBytesRead  *expvar.Int
	numVLogPhysicalBytesRead *expvar.Int // Assuming 4KB block size.

	NumGCReadsForProbing   *expvar.Int
	NumGCReadsForRewriting *expvar.Int
}

var LightstepMetrics *lightstepMetrics

const blockSize = 4 * 1024

func (metrics *lightstepMetrics) RecordLSMRead(offset int64, size int64) {
	metrics.numLSMReads.Add(1)
	metrics.numLSMLogicalBytesRead.Add(size)

	start := (offset / blockSize)
	limit := (offset + size + blockSize - 1) / blockSize
	metrics.numLSMPhysicalBytesRead.Add((limit - start + 1) * blockSize)
}

func (metrics *lightstepMetrics) RecordVLogRead(offset int64, size int64) {
	metrics.numVLogReads.Add(1)
	metrics.numVLogLogicalBytesRead.Add(size)

	start := (offset / blockSize)
	limit := (offset + size + blockSize - 1) / blockSize
	metrics.numVLogPhysicalBytesRead.Add((limit - start + 1) * blockSize)
}

// These variables are global and have cumulative values for all kv stores.
func init() {
	NumReads = expvar.NewInt("badger_disk_reads_total")
	NumWrites = expvar.NewInt("badger_disk_writes_total")
	NumBytesRead = expvar.NewInt("badger_read_bytes")
	NumBytesWritten = expvar.NewInt("badger_written_bytes")
	NumLSMGets = expvar.NewMap("badger_lsm_level_gets_total")
	NumLSMBloomHits = expvar.NewMap("badger_lsm_bloom_hits_total")
	NumGets = expvar.NewInt("badger_gets_total")
	NumPuts = expvar.NewInt("badger_puts_total")
	NumBlockedPuts = expvar.NewInt("badger_blocked_puts_total")
	NumMemtableGets = expvar.NewInt("badger_memtable_gets_total")
	LSMSize = expvar.NewMap("badger_lsm_size_bytes")
	VlogSize = expvar.NewMap("badger_vlog_size_bytes")
	PendingWrites = expvar.NewMap("badger_pending_writes_total")

	LightstepMetrics = &lightstepMetrics{
		numLSMReads:              expvar.NewInt("lightstep_badger_num_lsm_reads"),
		numLSMLogicalBytesRead:   expvar.NewInt("lightstep_badger_num_lsm_logical_bytes_read"),
		numLSMPhysicalBytesRead:  expvar.NewInt("lightstep_badger_num_lsm_physical_bytes_read"),
		numVLogReads:             expvar.NewInt("lightstep_badger_num_vlog_reads"),
		numVLogLogicalBytesRead:  expvar.NewInt("lightstep_badger_num_vlog_logical_bytes_read"),
		numVLogPhysicalBytesRead: expvar.NewInt("lightstep_badger_num_vlog_physical_bytes_read"),
		NumGCReadsForProbing:     expvar.NewInt("lightstep_badger_num_gc_reads_for_probing"),
		NumGCReadsForRewriting:   expvar.NewInt("lightstep_badger_num_gc_reads_for_rewriting"),
	}
}
