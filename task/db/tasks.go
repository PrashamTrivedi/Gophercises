package taskdb

import (
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var db *bolt.DB

type Task struct {
	Key   int
	Value string
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		return err
	})
}

func StoreTask(task string) (int, error) {
	var id int
	err := db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := intToByteSlice(id)
		return b.Put(key, []byte(task))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func ListTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(t *bolt.Tx) error {
		b := t.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{Key: byteSliceToint(k), Value: string(v)})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func DeleteTask(key int) error {
	return db.Update(func(t *bolt.Tx) error {
		b := t.Bucket(taskBucket)
		return b.Delete(intToByteSlice(key))
	})

}
func intToByteSlice(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func byteSliceToint(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
