package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/dumb-tech/brtool/entites"

	dynstruct "github.com/ompluscator/dynamic-struct"
)

const (
	eventTypeCreated = "rows.created"
	eventTypeUpdated = "rows.updated"
	eventTypeDeleted = "rows.deleted"
)

type CreateHandleFunc func(tableID int, items []dynstruct.Reader)
type UpdateHandleFunc func(tableID int, items []dynstruct.Reader, oldItems []dynstruct.Reader)
type DeleteHandleFunc func(tableID int, ids []int)

type EventHandler struct {
	CustomItemDef    dynstruct.DynamicStruct
	CustomWebhookDef dynstruct.DynamicStruct
	createHandler    CreateHandleFunc
	updateHandler    UpdateHandleFunc
	deleteHandler    DeleteHandleFunc
}

func NewEventHandler(definition any) *EventHandler {
	eh := &EventHandler{}
	eh.CustomItemDef = dynstruct.MergeStructs(entites.RowItem{}, definition).Build()
	eh.CustomWebhookDef = dynstruct.ExtendStruct(entites.Webhook{}).
		AddField("Items", eh.CustomItemDef.NewSliceOfStructs(), `json:"items,omitempty"`).
		AddField("OldItems", eh.CustomItemDef.NewSliceOfStructs(), `json:"old_items,omitempty"`).Build()

	return eh
}

func (eh *EventHandler) OnCreate(f CreateHandleFunc) { eh.createHandler = f }
func (eh *EventHandler) OnUpdate(f UpdateHandleFunc) { eh.updateHandler = f }
func (eh *EventHandler) OnDelete(f DeleteHandleFunc) { eh.deleteHandler = f }

func (eh *EventHandler) HandleRequest(req *http.Request) error {
	data := bytes.Buffer{}
	if _, err := data.ReadFrom(req.Body); err != nil {
		return err
	}
	defer req.Body.Close()

	return eh.Handle(data.Bytes())
}

func (eh *EventHandler) Handle(data []byte) error {
	webhook := eh.CustomWebhookDef.New()

	if err := json.Unmarshal(data, &webhook); err != nil {
		return err
	}

	whr := dynstruct.NewReader(webhook)
	evt := whr.GetField("EventType").String()

	switch evt {
	case eventTypeCreated:
		tableID := whr.GetField("TableId").Int()
		items := dynstruct.NewReader(whr.GetField("Items").Interface()).ToSliceOfReaders()
		eh.createHandler(tableID, items)
	case eventTypeUpdated:
		tableID := whr.GetField("TableId").Int()
		items := dynstruct.NewReader(whr.GetField("Items").Interface()).ToSliceOfReaders()
		oldItems := dynstruct.NewReader(whr.GetField("OldItems").Interface()).ToSliceOfReaders()
		eh.updateHandler(tableID, items, oldItems)
	case eventTypeDeleted:
		tableID := whr.GetField("TableId").Int()
		var ids []int
		cast, ok := whr.GetField("RowIDs").Interface().([]int)
		if ok {
			ids = cast
		}
		eh.deleteHandler(tableID, ids)
	default:
		return errors.New(fmt.Sprintf("unknown event type %q", evt))
	}

	return nil
}
