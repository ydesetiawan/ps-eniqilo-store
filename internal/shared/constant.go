package shared

import "os"

const (
	StockAdjustmentAccountName = "Inventory Adjustments"
)

var (
	ServiceName = os.Getenv("SERVICE_NAME_HTTP")
	DatadogEnv  = os.Getenv("DD_ENV")
)
