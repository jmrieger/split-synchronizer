package storage

import (
	"sync"
	"time"

	"github.com/splitio/go-split-commons/v4/storage"
	"github.com/splitio/go-split-commons/v4/storage/inmemory"
	"github.com/splitio/go-split-commons/v4/telemetry"
)

// Local telemetry constants
const (
	AuthEndpoint = iota
	SplitChangesEndpoint
	SegmentChangesEndpoint
	MySegmentsEndpoint
	ImpressionsBulkEndpoint
	ImpressionsBulkBeaconEndpoint
	ImpressionsCountEndpoint
	ImpressionsCountBeaconEndpoint
	EventsBulkEndpoint
	EventsBulkBeaconEndpoint
	TelemetryConfigEndpoint
	TelemetryRuntimeEndpoint
	LegacyTimeEndpoint
	LegacyTimesEndpoint
	LegacyCounterEndpoint
	LegacyCountersEndpoint
	LegacyGaugeEndpoint
)

type statusCodeMap struct {
	codes map[int]int64
	mutex sync.Mutex
}

func (s *statusCodeMap) incr(code int) {
	s.mutex.Lock()
	s.codes[code]++
	s.mutex.Unlock()
}

func (s *statusCodeMap) peek() map[int]int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	tmp := make(map[int]int64)
	for k, v := range s.codes {
		tmp[k] = v
	}

	return tmp
}

func newStatusCodeMap() statusCodeMap {
	return statusCodeMap{codes: make(map[int]int64)}
}

// EndpointStatusCodes keeps track of http status codes generated by the proxy endpoints
type EndpointStatusCodes struct {
	auth                   statusCodeMap
	splitChanges           statusCodeMap
	segmentChanges         statusCodeMap
	mySegments             statusCodeMap
	impressionsBulk        statusCodeMap
	impressionsBulkBeacon  statusCodeMap
	impressionsCount       statusCodeMap
	impressionsCountBeacon statusCodeMap
	eventsBulk             statusCodeMap
	eventsBulkBeacon       statusCodeMap
	telemetryConfig        statusCodeMap
	telemetryRuntime       statusCodeMap
	legacyTime             statusCodeMap
	legacyTimes            statusCodeMap
	legacyCounter          statusCodeMap
	legacyCounters         statusCodeMap
	legacyGauge            statusCodeMap
}

// IncrEndpointStatus increments the count of a specific status code for a specific endpoint
func (e *EndpointStatusCodes) IncrEndpointStatus(endpoint int, status int) {
	switch endpoint {
	case AuthEndpoint:
		e.auth.incr(status)
	case SplitChangesEndpoint:
		e.splitChanges.incr(status)
	case SegmentChangesEndpoint:
		e.segmentChanges.incr(status)
	case MySegmentsEndpoint:
		e.mySegments.incr(status)
	case ImpressionsBulkEndpoint:
		e.impressionsBulk.incr(status)
	case ImpressionsBulkBeaconEndpoint:
		e.impressionsBulkBeacon.incr(status)
	case ImpressionsCountEndpoint:
		e.impressionsCount.incr(status)
	case ImpressionsCountBeaconEndpoint:
		e.impressionsCountBeacon.incr(status)
	case EventsBulkEndpoint:
		e.eventsBulk.incr(status)
	case EventsBulkBeaconEndpoint:
		e.eventsBulkBeacon.incr(status)
	case TelemetryConfigEndpoint:
		e.telemetryConfig.incr(status)
	case TelemetryRuntimeEndpoint:
		e.telemetryRuntime.incr(status)
	case LegacyTimeEndpoint:
		e.legacyTime.incr(status)
	case LegacyTimesEndpoint:
		e.legacyTimes.incr(status)
	case LegacyCounterEndpoint:
		e.legacyCounter.incr(status)
	case LegacyCountersEndpoint:
		e.legacyCounters.incr(status)
	case LegacyGaugeEndpoint:
		e.legacyGauge.incr(status)
	}
}

// PeekEndpointStatus increments the count of a specific status code for a specific endpoint
func (e *EndpointStatusCodes) PeekEndpointStatus(endpoint int) map[int]int64 {
	switch endpoint {
	case AuthEndpoint:
		return e.auth.peek()
	case SplitChangesEndpoint:
		return e.splitChanges.peek()
	case SegmentChangesEndpoint:
		return e.segmentChanges.peek()
	case MySegmentsEndpoint:
		return e.mySegments.peek()
	case ImpressionsBulkEndpoint:
		return e.impressionsBulk.peek()
	case ImpressionsBulkBeaconEndpoint:
		return e.impressionsBulkBeacon.peek()
	case ImpressionsCountEndpoint:
		return e.impressionsCount.peek()
	case ImpressionsCountBeaconEndpoint:
		return e.impressionsCountBeacon.peek()
	case EventsBulkEndpoint:
		return e.eventsBulk.peek()
	case EventsBulkBeaconEndpoint:
		return e.eventsBulkBeacon.peek()
	case TelemetryConfigEndpoint:
		return e.telemetryConfig.peek()
	case TelemetryRuntimeEndpoint:
		return e.telemetryRuntime.peek()
	case LegacyTimeEndpoint:
		return e.legacyTime.peek()
	case LegacyTimesEndpoint:
		return e.legacyTimes.peek()
	case LegacyCounterEndpoint:
		return e.legacyCounter.peek()
	case LegacyCountersEndpoint:
		return e.legacyCounters.peek()
	case LegacyGaugeEndpoint:
		return e.legacyGauge.peek()
	}
	return nil
}

func newEndpointStatusCodes() EndpointStatusCodes {
	return EndpointStatusCodes{
		auth:                   newStatusCodeMap(),
		splitChanges:           newStatusCodeMap(),
		segmentChanges:         newStatusCodeMap(),
		mySegments:             newStatusCodeMap(),
		impressionsBulk:        newStatusCodeMap(),
		impressionsBulkBeacon:  newStatusCodeMap(),
		impressionsCount:       newStatusCodeMap(),
		impressionsCountBeacon: newStatusCodeMap(),
		eventsBulk:             newStatusCodeMap(),
		eventsBulkBeacon:       newStatusCodeMap(),
		telemetryConfig:        newStatusCodeMap(),
		telemetryRuntime:       newStatusCodeMap(),
		legacyTime:             newStatusCodeMap(),
		legacyTimes:            newStatusCodeMap(),
		legacyCounter:          newStatusCodeMap(),
		legacyCounters:         newStatusCodeMap(),
		legacyGauge:            newStatusCodeMap(),
	}
}

// ProxyEndpointLatencies defines an interface to access proxy server endpoint latencies numbers
type ProxyEndpointLatencies interface {
	PeekEndpointLatency(endpoint int) []int64
	RecordEndpointLatency(endpoint int, latency time.Duration)
}

// ProxyEndpointLatenciesImpl keep track of the latency introudiced by each proxy endpoint
type ProxyEndpointLatenciesImpl struct {
	auth                   inmemory.AtomicInt64Slice
	splitChanges           inmemory.AtomicInt64Slice
	segmentChanges         inmemory.AtomicInt64Slice
	mySegments             inmemory.AtomicInt64Slice
	impressionsBulk        inmemory.AtomicInt64Slice
	impressionsBulkBeacon  inmemory.AtomicInt64Slice
	impressionsCount       inmemory.AtomicInt64Slice
	impressionsCountBeacon inmemory.AtomicInt64Slice
	eventsBulk             inmemory.AtomicInt64Slice
	eventsBulkBeacon       inmemory.AtomicInt64Slice
	telemetryConfig        inmemory.AtomicInt64Slice
	telemetryRuntime       inmemory.AtomicInt64Slice
	legacyTime             inmemory.AtomicInt64Slice
	legacyTimes            inmemory.AtomicInt64Slice
	legacyCounter          inmemory.AtomicInt64Slice
	legacyCounters         inmemory.AtomicInt64Slice
	legacyGauge            inmemory.AtomicInt64Slice
}

// RecordEndpointLatency records a (bucketed) latency for a specific endpoint
func (p *ProxyEndpointLatenciesImpl) RecordEndpointLatency(endpoint int, latency time.Duration) {
	bucket := telemetry.Bucket(latency.Milliseconds())
	switch endpoint {
	case AuthEndpoint:
		p.auth.Incr(bucket)
	case SplitChangesEndpoint:
		p.splitChanges.Incr(bucket)
	case SegmentChangesEndpoint:
		p.segmentChanges.Incr(bucket)
	case MySegmentsEndpoint:
		p.mySegments.Incr(bucket)
	case ImpressionsBulkEndpoint:
		p.impressionsBulk.Incr(bucket)
	case ImpressionsBulkBeaconEndpoint:
		p.impressionsBulkBeacon.Incr(bucket)
	case ImpressionsCountEndpoint:
		p.impressionsCount.Incr(bucket)
	case ImpressionsCountBeaconEndpoint:
		p.impressionsCountBeacon.Incr(bucket)
	case EventsBulkEndpoint:
		p.eventsBulk.Incr(bucket)
	case EventsBulkBeaconEndpoint:
		p.eventsBulkBeacon.Incr(bucket)
	case TelemetryConfigEndpoint:
		p.telemetryRuntime.Incr(bucket)
	case TelemetryRuntimeEndpoint:
		p.telemetryConfig.Incr(bucket)
	case LegacyTimeEndpoint:
		p.legacyTime.Incr(bucket)
	case LegacyTimesEndpoint:
		p.legacyTimes.Incr(bucket)
	case LegacyCounterEndpoint:
		p.legacyCounter.Incr(bucket)
	case LegacyCountersEndpoint:
		p.legacyCounters.Incr(bucket)
	case LegacyGaugeEndpoint:
		p.legacyGauge.Incr(bucket)
	}
}

// PeekEndpointLatency records a (bucketed) latency for a specific endpoint
func (p *ProxyEndpointLatenciesImpl) PeekEndpointLatency(endpoint int) []int64 {
	switch endpoint {
	case AuthEndpoint:
		return p.auth.ReadAll()
	case SplitChangesEndpoint:
		return p.splitChanges.ReadAll()
	case SegmentChangesEndpoint:
		return p.segmentChanges.ReadAll()
	case MySegmentsEndpoint:
		return p.mySegments.ReadAll()
	case ImpressionsBulkEndpoint:
		return p.impressionsBulk.ReadAll()
	case ImpressionsBulkBeaconEndpoint:
		return p.impressionsBulkBeacon.ReadAll()
	case ImpressionsCountEndpoint:
		return p.impressionsCount.ReadAll()
	case ImpressionsCountBeaconEndpoint:
		return p.impressionsCountBeacon.ReadAll()
	case EventsBulkEndpoint:
		return p.eventsBulk.ReadAll()
	case EventsBulkBeaconEndpoint:
		return p.eventsBulkBeacon.ReadAll()
	case TelemetryConfigEndpoint:
		return p.telemetryRuntime.ReadAll()
	case TelemetryRuntimeEndpoint:
		return p.telemetryConfig.ReadAll()
	case LegacyTimeEndpoint:
		return p.legacyTime.ReadAll()
	case LegacyTimesEndpoint:
		return p.legacyTimes.ReadAll()
	case LegacyCounterEndpoint:
		return p.legacyCounter.ReadAll()
	case LegacyCountersEndpoint:
		return p.legacyCounters.ReadAll()
	case LegacyGaugeEndpoint:
		return p.legacyGauge.ReadAll()
	}
	return nil
}

// newProxyEndpointLatenciesImpl creates a new latency tracker
func newProxyEndpointLatenciesImpl() ProxyEndpointLatenciesImpl {
	init := func() inmemory.AtomicInt64Slice {
		toRet, _ := inmemory.NewAtomicInt64Slice(telemetry.LatencyBucketCount)
		return toRet
	}

	return ProxyEndpointLatenciesImpl{
		auth:                   init(),
		splitChanges:           init(),
		segmentChanges:         init(),
		mySegments:             init(),
		impressionsBulk:        init(),
		impressionsBulkBeacon:  init(),
		impressionsCount:       init(),
		impressionsCountBeacon: init(),
		eventsBulk:             init(),
		eventsBulkBeacon:       init(),
		telemetryConfig:        init(),
		telemetryRuntime:       init(),
		legacyTime:             init(),
		legacyTimes:            init(),
		legacyCounter:          init(),
		legacyCounters:         init(),
		legacyGauge:            init(),
	}
}

// ProxyTelemetryPeeker is able to peek at locally captured metrics
type ProxyTelemetryPeeker interface {
	PeekEndpointLatency(resource int) []int64
	PeekEndpointStatus(resource int) map[int]int64
}

// ProxyEndpointTelemetry defines the interface that endpoints use to capture latency & status codes
type ProxyEndpointTelemetry interface {
	ProxyTelemetryPeeker
	RecordEndpointLatency(endpoint int, latency time.Duration)
	IncrEndpointStatus(endpoint int, status int)
}

// ProxyTelemetryFacade defines the set of methods required to accept local telemetry as well as runtime telemetry
type ProxyTelemetryFacade interface {
	storage.TelemetryStorage
	ProxyEndpointTelemetry
}

// ProxyTelemetryFacadeImpl exposes local telemetry functionality
type ProxyTelemetryFacadeImpl struct {
	ProxyEndpointLatenciesImpl
	EndpointStatusCodes
	*inmemory.TelemetryStorage
}

// NewProxyTelemetryFacade instantiates a local telemetry facade
func NewProxyTelemetryFacade() *ProxyTelemetryFacadeImpl {
	ts, _ := inmemory.NewTelemetryStorage()
	return &ProxyTelemetryFacadeImpl{
		ProxyEndpointLatenciesImpl: newProxyEndpointLatenciesImpl(),
		EndpointStatusCodes:        newEndpointStatusCodes(),
		TelemetryStorage:           ts,
	}
}

// Ensure interface compliance
var _ ProxyTelemetryFacade = (*ProxyTelemetryFacadeImpl)(nil)
var _ storage.TelemetryStorage = (*ProxyTelemetryFacadeImpl)(nil)
