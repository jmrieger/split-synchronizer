package common

import "github.com/splitio/go-split-commons/v4/storage"

// Storages wraps storages in one struct
type Storages struct {
	SplitStorage          storage.SplitStorage
	SegmentStorage        storage.SegmentStorage
	LocalTelemetryStorage storage.TelemetryRuntimeConsumer
	EventStorage          storage.EventMultiSdkConsumer
	ImpressionStorage     storage.ImpressionMultiSdkConsumer
}
