package database

import (
	"fmt"
	"log"
	"sync"
)

type InsertQueue struct {
	insert          chan InsertQueueItem
	db              *Database
	closeProtection sync.Once
	closed          bool
}

func NewInsertQueue(db *Database) *InsertQueue {
	return &InsertQueue{
		insert:          make(chan InsertQueueItem, 100),
		closeProtection: sync.Once{},
		closed:          false,
		db:              db,
	}
}

func (i *InsertQueue) AddToInsertQueue(table string, row ...RowType) {
	i.insert <- InsertQueueItem{
		Table: table,
		Rows:  row,
		DB:    i.db,
	}
}

func (i *InsertQueue) Start() error {
	if i.closed {
		return fmt.Errorf("queue is already closed")
	}
	go func() {
	Exit:
		for {
			select {
			case item, ok := <-i.insert:
				if !ok {
					log.Printf("Insert queue channel was closed, breaking out")
					break Exit
				}
				if len(item.Rows) == 0 {
					continue
				} else if len(item.Rows) == 1 {
					_, err := i.db.Insert(item.Table, item.Rows[0])
					if err != nil {
						log.Printf("Error in Async Insert: %v\n", err)
					}
				} else {
					_, err := i.db.BatchInsert(item.Table, item.Rows...)
					if err != nil {
						log.Printf("Error in Async Batch Insert: %v\n", err)
					}
				}
			}
		}
	}()
	return nil
}

func (i *InsertQueue) Stop() {
	i.closeProtection.Do(func() {
		close(i.insert)
		i.closed = true
	})
}
