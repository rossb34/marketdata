package marketdata

import "time"

type MDEntryType int

const (
	BID MDEntryType = iota
	OFFER
	TRADE
)

type MDUpdateAction int

const (
	NEW MDUpdateAction = iota
	CHANGE
	DELETE
	NONE
)

type Message int

const (
	MD_SNAPSHOT_FULL_REFRESH Message = iota
	MD_INCREMENTAL_REFRESH
)

type MDEntry struct {
	Action            MDUpdateAction `json:"md_update_action"`
	Type              MDEntryType    `json:"md_entry_type"`
	Symbol            string         `json:"symbol"`
	RptSequenceNumber uint64         `json:"rpt_sequence_number"`
	Price             Price          `json:"price"`
	Size              Quantity       `json:"size"`
	NumberOfOrders    int32          `json:"number_of_orders"`
	PriceLevelIndex   int32          `json:"price_level"`
}

type MDSnapshotFullRefresh struct {
	MessageType       Message   `json:"message_type"`
	Timestamp         time.Time `json:"timestamp"`
	EndpointSendTime  time.Time `json:"endpoint_send_time"`
	TransactTime      time.Time `json:"transact_time"`
	Symbol            string    `json:"symbol"`
	EndpointName      string    `json:"endpoint_name"`
	MsgSequenceNumber uint64    `json:"msg_sequence_number"`
	RptSequenceNumber uint64    `json:"rpt_sequence_number"`
	Entries           []MDEntry `json:"entries"`
}

// Follow the sbe pattern to initialize default values with a *Init function

/// Initializes a MDSnapshotFullRefresh
func MDSnapshotFullRefreshInit(m *MDSnapshotFullRefresh) {
	m.MessageType = MD_SNAPSHOT_FULL_REFRESH
}

type MDIncrementalRefresh struct {
	MessageType       Message   `json:"message_type"`
	Timestamp         time.Time `json:"timestamp"`
	EndpointSendTime  time.Time `json:"endpoint_send_time"`
	TransactTime      time.Time `json:"transact_time"`
	EndpointName      string    `json:"endpoint_name"`
	MsgSequenceNumber uint64    `json:"msg_sequence_number"`
	Entries           []MDEntry `json:"entries"`
}

/// Initializes a MDIncrementalRefresh
func MDIncrementalRefreshInit(m *MDIncrementalRefresh) {
	m.MessageType = MD_INCREMENTAL_REFRESH
}
