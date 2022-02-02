package database

import (
	"container/list"
	"log"
)

var insertQueue *list.List

func addToInsertQueue(db *Database, table string, row RowType) {

	if insertQueue == nil {
		insertQueue = list.New()
	}

	insertQueue.PushBack(InsertQueueItem{
		Table: table,
		Row:   row,
		DB:    db,
	})

	executeInsertAsync()
}

func executeInsertAsync() {

	go func() {

		for e := insertQueue.Front(); e != nil; e = e.Next() {
			insertQueue.Remove(e)

			ins := e.Value.(InsertQueueItem)

			_, err := ins.DB.Insert(ins.Table, ins.Row)
			if err != nil {
				log.Printf("Error in Async Insert: %v\n", err)
			}
		}

	}()
}
