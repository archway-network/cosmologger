package database

import (
	"container/list"
	"log"
	"sync"
)

var insertQueue *list.List
var insertQueueMutex sync.Mutex

func addToInsertQueue(db *Database, table string, row RowType) {

	go func() {

		insertQueueMutex.Lock()

		if insertQueue == nil {
			insertQueue = list.New()
		}
		insertQueue.PushBack(InsertQueueItem{
			Table: table,
			Row:   row,
			DB:    db,
		})

		insertQueueMutex.Unlock()
		executeInsertAsync()
	}()
}

func executeInsertAsync() {

	go func() {

		insertQueueMutex.Lock()
		defer insertQueueMutex.Unlock()

		if insertQueue == nil {
			return
		}
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
