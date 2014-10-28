package dht

import (
  "fmt"
  "net"
  "bytes"
  "time"
  "encoding/json"
//  "math/big"
//  "strconv"
//  "strings"
//  "os"
//  "sync"
)


func doStartupReplication(t *Table, nod *Node) {

  time.Sleep(25000 * time.Millisecond)
  for nod.nextNode == nil {
    time.Sleep(25000 * time.Millisecond)
  }

fmt.Println("Starting startup replication", nod.nodeId)
  msg := new (Message)
  msg.Idx    = 21
  msg.Key    = ""
  msg.Gen    = numCopies
  msg.Info   = nil
  port := answPort.getPort() // += 210
  service := fmt.Sprintf("%s:%d", nod.localAddress, port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err, 201)
  msg.Src = localAddr
  doCall(nod, msg)
  service = fmt.Sprintf(":%d", port)
  localAddr, err = net.ResolveUDPAddr("udp", service)
  checkError(err, 202)
  conn, err := net.ListenUDP("udp", localAddr)
  checkError(err, 203)
  for {
    answ := waitForRec(conn)
    if answ.Info == nil {
      conn.Close()
      return
    }
    if answ.Gen >= 0 {
      _ = t.Insert(answ.Gen, 2, answ.Version, answ.Key, answ.Info)
    }
  }
  conn.Close()
}


func waitForRec(conn *net.UDPConn) *Message {

  var buf [512]byte
  _, err := conn.Read(buf[0:])
  checkError(err, 204)
  dec  := json.NewDecoder(bytes.NewReader([]byte(buf[0:])))
  answ := new (Message)
  answ.Info = nil
  err = dec.Decode(&answ)
  checkError(err, 205)
  return answ
}


func scanReplication(t *Table, nod *Node) {

  record := t.head
  for {
    if record == nil {
//fmt.Println("Scan ready")
      return
    }
    t.lock.Lock()
    status  := record.status
    key     := record.Key
    gen     := record.gen
    version := record.version
    data    := record.Data
fmt.Println("Testrecord: ", record)
    t.lock.Unlock()

    switch status {
    case 0:    // Deleted, check if copy still existing, if not - remove
      if gen < numCopies {
fmt.Println("Check if remote is deleted")
        stat, _ := checkCopyStatus(nod, key)
        if stat < 0 {
          t.Remove(key)
        }
      } else {
        t.Remove(key)
      }
    case 2:    // Up-to-date, validate anyway
      if gen < numCopies {
        stat, vers := checkCopyStatus(nod, key)
        if stat < 0 {
fmt.Println("Insert remote")
          _= updateCopy(nod, gen, version, key, data)
        } else {
          if version != vers {
fmt.Println("Flag as updated")
            t.lock.Lock()  // Something is not up to date, flag as updated
            record.status = 3
            t.lock.Unlock()
          }
        }
      }
    case 3:    // Updated, populate
      if gen < numCopies {
fmt.Println("Check remote status for update")
        stat, vers := checkCopyStatus(nod, key)
fmt.Println("Stat, vers, version", stat, vers, version)
        if (stat == 2) && (version == vers) {
fmt.Println("Remote is up to date")
          t.lock.Lock()  // Copy is up to date, flag as up to date
          record.status = 2
          t.lock.Unlock()
        } else {
          if vers > version {   // Remote is newer
            answ := getRemote(nod, key)
fmt.Println("Update local", answ)
            t.lock.Lock()
            record.version = answ.Version
            record.Data    = answ.Info
            record.status = 2
            t.lock.Unlock()
          } else {
fmt.Println("Update remote", record)
            _= updateCopy(nod, gen, version, key, data)
          }
        }
      }
    case 1:    // Invalidated, wait for update, do nothing
    }
  record = record.next
  }
}


func getRemote(nod *Node, key string) *Message {

  msg := new (Message)
  msg.Idx = 14
  msg.Key = key
  msg.Dst = nod.nextNode
  port := answPort.getPort() // += 210
  service := fmt.Sprintf("%s:%d", nod.localAddress, port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err, 206)
  msg.Src = localAddr
  doRemote(msg)
  answ := waitForResult(port)
  return answ
}


func checkCopyStatus(nod *Node, key string) (int, int) {

  msg := new (Message)
  msg.Idx = 19
  msg.Key = key
  msg.Src = nod.localNode
  port := answPort.getPort() // += 210
  service := fmt.Sprintf("%s:%d", nod.localAddress, port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err, 207)
  msg.Src  = localAddr
  msg.Info = nil
  doCall(nod, msg)
  answ := waitForResult(port)
  return answ.Status, answ.Version
}


func updateCopy(nod *Node, gen, vers int, key string, info interface{}) int {

fmt.Println("Update copy", nod.nodeId, nod.nextKey, nod.nextNode)
  msg := new (Message)
  msg.Idx     = 20
  msg.Key     = key
  msg.Gen     = gen + 1
  msg.Version = vers
  msg.Info    = info
  port := answPort.getPort() // += 210
  service := fmt.Sprintf("%s:%d", nod.localAddress, port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err, 208)
  msg.Src = localAddr
  doCall(nod, msg)
  answ := waitForResult(port)
  return answ.Status
}


func doReplication(t *Table, nod *Node) {

  time.Sleep(25000 * time.Millisecond)
  for {
fmt.Println("Starting periodic replication", nod.nodeId)
    scanReplication(t, nod)
    time.Sleep(15000 * time.Millisecond)
  }
}


func invalidateCopy(nod *Node, msg *Message) {

  mess := new (Message)
  mess.Idx  = 18
  mess.Key  = msg.Key
  mess.Gen  = msg.Gen + 1
  mess.Src  = nod.nextNode
  mess.Info = nil
  doCall(nod, mess)
}


func doInvalidateCopy(t *Table, key string) {

  record := t.findRecord(key)
  if record != nil {
    t.lock.Lock()
    record.status = 1
    t.lock.Unlock()
  }
}


func markCopyAsDeleted(nod *Node, msg *Message) {

  mess := new (Message)
  mess.Idx  = 17
  mess.Key  = msg.Key
  mess.Gen  = msg.Gen + 1
  mess.Src  = nod.localNode
  mess.Dst  = nod.nextNode
  mess.Info = nil
  doCall(nod, mess)
}


func doCall(nod *Node, msg *Message) {

  msg.Dst  = nod.nextNode
//  msg.Info = nil
  conn, err := net.DialUDP("udp", nil, msg.Dst)
  checkError(err, 209)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err, 210)
  _, err = conn.Write(buffer)
  checkError(err, 211)
}


func getRecordStatus(t *Table, key string) (int, int) {

  record := t.findRecord(key)
  if record == nil {
    return -1, 0
  }
  return record.status, record.version
}
