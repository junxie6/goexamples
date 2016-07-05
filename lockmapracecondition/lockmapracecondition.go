package main

import (
        "fmt"
        "sync"
)

// One approach would be to not export the Table member and provide Get/Set methods that lock appropriately.
// This way the only way to access the map's data is via the methods which guard access to the map. Each value is then guarded independently.
// Reference:
// http://stackoverflow.com/questions/25133116/locking-golang-recursive-map
// http://stackoverflow.com/questions/4498998/how-to-initialize-members-in-go-struct#
// https://play.golang.org/p/faO9six-Qx
type SyncMap struct {
        sync.RWMutex // embedded.  see http://golang.org/ref/spec#Struct_types
        hm           map[string]string
}

func NewSyncMap() *SyncMap {
        return &SyncMap{hm: make(map[string]string)}
}

func (m *SyncMap) Put(k string, v string) {
        m.Lock()
        defer m.Unlock()
        m.hm[k] = v
}

func (m *SyncMap) Get(k string) string {
        m.RLock()
        defer m.RUnlock()
        return m.hm[k]
}

func main() {
        sm := NewSyncMap()
        sm.Put("kTest", "vTest")
        fmt.Println(sm.Get("kTest"))
}
