package dht

import (
  "fmt"
//  "net"
//  "bytes"
  "time"
//  "encoding/json"
//  "math/big"
  "math/rand"
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

  var record *Record
fmt.Println("Head: ", &t.head)
fmt.Println("Tail: ", &t.tail)
  if t.head != nil {
      record = t.head
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
fmt.Println("Server loop")
      time.Sleep(500 * time.Millisecond)
fmt.Println("Record: ", &record)
fmt.Println("Next: ",   &record.next)
fmt.Println("Prev: ",   &record.prev)
      record = record.next
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
  record := t.head
  for {
    if strings.EqualFold(key, record.Key) {
      return record
    }
    if record.next != nil {
      record = record.next
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
  first := t.Head()

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
fmt.Println("Duplicate", key)
    return 0     // Failure, duplicate
  }
  newRec := &Record{key, val, gen, rnd, stat, nil, nil, t}
  t.lock.Lock()
  if t.head == nil {
    t.head      = newRec
    t.tail      = newRec
    newRec.prev = newRec
  } else {
    t.head.prev = newRec
    newRec.prev = t.tail
    t.tail.next = newRec
    t.tail      = newRec
//    t.head = newRec
////    t.head.prev = nil
//    t.head.next.prev = newRec
  }
  t.len++
  t.lock.Unlock()

fmt.Println("Inserted", newRec.Data, newRec.status, newRec.gen)
  return 1    // Success
}


func (t *Table) getRecord(key string) (interface{}, int) {

fmt.Println("Get record: ", key)
  record := t.findRecord(key)
  if record == nil {
fmt.Println("Found nothing 1")
    return nil, -1
  }
fmt.Println("Rec stat: ", record.status)
  if (record.status > 1) {
fmt.Println("Get data: ", record.Data, record.gen)
    return record.Data, record.gen
  }
fmt.Println("Found nothing 2")
  return nil, -1
}


func (t *Table) update(gen, stat, vers int, key string, val interface{}) int {

fmt.Println("Update", key)
  record := t.findRecord(key)
fmt.Println("Found: ", record)
  if record != nil {
    t.lock.Lock()
    if vers == 0 {
      tmp := record.version
      record.version = tmp + 1
fmt.Println("Oldvers: ", record.version)
      record.status  = 3
    } else {
      record.version = vers
      record.status  = stat
    }
    record.Data    = val
    t.lock.Unlock()
    return 1
  } else {
    rnd  := rand.New(rand.NewSource(99))
    vers := rnd.Intn(25000)
    return t.Insert(gen, stat, vers, key, val)
  }
  return 0
}


func (t *Table) Delete(key string) int {

fmt.Println("Delete", key)
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

fmt.Println("Remove", key)
  t.lock.RLock()
  record := t.findRecord(key)
  t.lock.RUnlock()
  if record == nil {
fmt.Println("Remove not found")
    return 0           // Fail, not found
  }

  t.lock.Lock()
  if record == t.head {
    if t.len > 1 {
fmt.Println("Remove single")
      t.head = record.next
      t.head.prev = record.prev
    } else {
fmt.Println("Remove first")
      t.head = nil
      t.tail = nil
    }
  } else {
    if record == t.tail {
fmt.Println("Remove last")
      t.tail      = record.prev
      t.tail.next = nil
      t.head.prev = t.tail
    } else {
fmt.Println("Remove in list")
      record.next.prev = record.prev
      record.prev.next = record.next
    }
  }
  t.len--
  t.lock.Unlock()

  record.prev  = nil
  record.next  = nil
  record.table = nil
  record.Data  = nil
  record.Key   = ""
  t.len--
  return 1    // Success
}


func (t *Table) Head() *Record {
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
