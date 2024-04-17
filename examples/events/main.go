package main

import (
	"brtool/entites"
	"brtool/handler"
	"fmt"
	dynstruct "github.com/ompluscator/dynamic-struct"
	"log"
	"time"
)

const createDataSample = `{
    "table_id": 531,
    "database_id": 101,
    "workspace_id": 100,
    "event_id": "c6516a4f-47fd-4609-8b21-f81386266e46",
    "event_type": "rows.created",
    "items": [
        {
            "id": 6,
            "order": "2.00000000000000000000",
            "Number": null,
            "Updated": "2024-04-15T17:27:05.492550Z",
            "Url": {
                "url": "(\u043d\u0435 \u043e\u0431\u0440\u0430\u0431\u043e\u0442\u0430\u043d\u043e)",
                "label": null
            },
            "Report": "(\u043d\u0435 \u043e\u0431\u0440\u0430\u0431\u043e\u0442\u0430\u043d\u043e)"
        }
    ]
}`

const updateDataSample = `{
    "table_id": 531,
    "database_id": 101,
    "workspace_id": 100,
    "event_id": "cf2b25b9-09c1-47d7-80bd-93cb984c6720",
    "event_type": "rows.updated",
    "items": [
        {
            "id": 8,
            "order": "1.00000000000000000000",
            "Number": "14846168484",
            "Updated": "2024-04-15T17:34:16.070200Z",
            "Url": {
                "url": "(\u043d\u0435 \u043e\u0431\u0440\u0430\u0431\u043e\u0442\u0430\u043d\u043e)",
                "label": null
            },
            "Report": "(\u043d\u0435 \u043e\u0431\u0440\u0430\u0431\u043e\u0442\u0430\u043d\u043e)"
        }
    ],
    "old_items": [
        {
            "id": 8,
            "order": "1.00000000000000000000",
            "Number": "4846168484",
            "Updated": "2024-04-15T17:32:41.250376Z",
            "Url": {
                "url": "(\u043d\u0435 \u043e\u0431\u0440\u0430\u0431\u043e\u0442\u0430\u043d\u043e)",
                "label": null
            },
            "Report": "(\u043d\u0435 \u043e\u0431\u0440\u0430\u0431\u043e\u0442\u0430\u043d\u043e)"
        }
    ]
}`

const deleteDataSample = `{
    "table_id": 531,
    "database_id": 101,
    "workspace_id": 100,
    "event_id": "17295694-2b1d-433a-aa5e-fec9ddb579b7",
    "event_type": "rows.deleted",
    "row_ids": [7, 6, 1]
}`

type WorkingStruct struct {
	entites.RowItem
	Number  string    `json:"Number"`
	Updated time.Time `json:"Updated"`
	URL     struct {
		URL   string `json:"url"`
		Label string `json:"label"`
	} `json:"Url"`
	Report string `json:"Report"`
}

func main() {
	h := handler.NewEventHandler(WorkingStruct{})

	h.OnCreate(func(itemsSet []dynstruct.Reader) {
		fmt.Println("--- event:create ---")
		for _, item := range itemsSet {
			fmt.Println(item.GetField("Number").String())
			fmt.Println(item.GetField("Updated").Time().String())
		}
	})
	h.OnUpdate(func(itemsSet []dynstruct.Reader, oldItemsSet []dynstruct.Reader) {
		fmt.Println("--- event:update ---")
		fmt.Println("--- event:update:items ---")
		for _, item := range itemsSet {
			fmt.Println(item.GetField("Number").String())
			fmt.Println(item.GetField("Updated").Time().String())
		}
		fmt.Println("--- event:update:old-items ---")
		for _, item := range oldItemsSet {
			fmt.Println(item.GetField("Number").String())
			fmt.Println(item.GetField("Updated").Time().String())
		}
	})
	h.OnDelete(func(ids []int) {
		fmt.Println("--- event:delete ---")
		for _, did := range ids {
			fmt.Println("-- Deleted row ID:", did)
		}
	})

	if err := h.Handle([]byte(createDataSample)); err != nil {
		log.Fatal(err)
	}
	if err := h.Handle([]byte(updateDataSample)); err != nil {
		log.Fatal(err)
	}
	if err := h.Handle([]byte(deleteDataSample)); err != nil {
		log.Fatal(err)
	}
}
