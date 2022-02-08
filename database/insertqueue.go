package database

import (
	"container/list"
	"log"
	"sync"
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

var insertQueueMutex sync.Mutex

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
