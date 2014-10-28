package dht

import (
//  "fmt"
//  "net"
//  "bytes"
  "time"
//  "encoding/json"
//  "math/big"
//  "strconv"
  "strings"
//  "os"
  "sync"
)

type Table struct {
  head *Record
  tail *Record
  len   int
  lock  sync.RWMutex
}


type Record struct {
  Key     string
  Data    interface{}
  gen     int         // 0=orig, 1, 2 copies
  version int
  status  int         // 0=del, 1=invalidated, 2=up to date, 3=updated
  next   *Record
  prev   *Record
  table  *Table
}


//var table = createTable()

func (t *Table) Len() int {
  return t.len
}


func (t *Table) updGetRecords(msg *Message) {

  if t.head != nil {
      record := t.head
    for {
      mess := new (Message)
      mess.Key     = record.Key
      mess.Info    = record.Data
      mess.Gen     = record.gen - 1
      mess.Version = record.version
      mess.Status  = record.status
      mess.Dst     = msg.Src
      mess.Src     = msg.Dst
      doRemote(mess)
      time.Sleep(500 * time.Millisecond)
      record := record.next
      if record == nil {
        break;
      }
    }
  }
  mess := new (Message)
  mess.Info = nil
  mess.Dst  = msg.Src
  mess.Src  = msg.Dst
  doRemote(mess)
}


func (t *Table) findRecord(key string) *Record {

  if t.head == nil {
    return nil
  }
  first := t.First()
  for {
    if strings.EqualFold(key, first.Key) {
      return first
    }
    if first.next != nil {
      first = first.next
    } else {
      return nil
    }
  }
  return nil
}


func (t *Table) isInTable(key string) bool {

  if t.head == nil {
    return false
  }
  first := t.First()

  for {
    if strings.EqualFold(key, first.Key) {
      return true
    }
    if first.next != nil {
      first = first.next
    } else {
      return false
    }
  }
  return false
}


func (t *Table) Insert (gen, stat, rnd int,key string, val interface{}) int {

  if t.isInTable(key) {
    return 0     // Failure, duplicate
  }
  newRec := &Record{key, val, gen, rnd, stat, t.head, t.tail, t}
  t.lock.Lock()
  defer t.lock.Unlock()
  if t.head == nil {
    t.head      = newRec
    t.tail      = newRec
    newRec.prev = newRec
  } else {
    t.head = newRec
//    t.head.prev = nil
    t.head.next.prev = newRec
  }
  t.len++

//fmt.Println("Inserted", newRec.Data, newRec.status)
  return 1    // Success
}


func (t *Table) getRecord(key string) (interface{}, int) {

  record := t.findRecord(key)
  if record == nil {
    return nil, -1
  }
  if (record.status > 1) {
    return record.Data, record.gen
  }
  return nil, -1
}


func (t *Table) update(gen, stat, vers int, key string, val interface{}) int {

//fmt.Println("Update", key)
  record := t.findRecord(key)
  if record != nil {
    t.lock.Lock()
    if vers == 0 {
      tmp := record.version
      record.version = tmp + 1
    } else {
      record.version = vers
    }
    record.Data    = val
    record.version = vers
    record.status  = stat
    t.lock.Unlock()
  } else {
    return t.Insert(gen, stat, vers, key, val)
  }
  return 0
}


func (t *Table) Delete(key string) int {

//fmt.Println("Delete", key)
  record := t.findRecord(key)
  if record != nil {
    t.lock.Lock()
    record.status = 0
    t.lock.Unlock()
    return 1
  }
  return 0
}


func (t *Table) Remove(key string) int {

//fmt.Println("Remove", key)
  t.lock.RLock()
  record := t.findRecord(key)
  t.lock.RUnlock()
  if record == nil {
    return 0           // Fail, not found
  }

  t.lock.Lock()
  defer t.lock.Unlock()
  if record == t.head {
    t.head = record.next
    t.head.prev = record.prev
  } else {
    record.prev.next = record.next
    if record == t.tail {
      t.tail = record.prev
      t.head.prev = record.prev
    } else {
      record.next.prev = record.prev
    }
  }
  record.prev  = nil
  record.next  = nil
  record.table = nil
  record.Data  = nil
  record.Key   = ""
  t.len--
  return 1    // Success
}


func (t *Table) First() *Record {
  return t.head
}


func (t *Table) Tail() *Record {
  return t.tail
}


func (r *Record) Prev() *Record {
  return r.prev
}


func (r *Record) Next() *Record {
  return r.next
}


func createTable() *Table {

  tab     := &Table{}
  tab.len  = 0
  tab.head = nil
  tab.tail = nil
  return tab
}
