package entites

type RowItem struct {
	ID    int    `json:"id"`
	Order string `json:"order"`
}

type Webhook struct {
	TableId     int    `json:"table_id"`
	DatabaseId  int    `json:"database_id"`
	WorkspaceId int    `json:"workspace_id"`
	EventId     string `json:"event_id"`
	EventType   string `json:"event_type"`
	RowIDs      []int  `json:"row_ids,omitempty"`
}
