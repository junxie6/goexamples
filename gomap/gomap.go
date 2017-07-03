package gomap

import (
	//"fmt"
	"hash/fnv"
)

const (
	// Maximum number of key/value pairs a bucket can hold.
	bucketCntBits = 3
	bucketCnt     = 1 << bucketCntBits
)

func NewGoMap() *GoMap {
	defaultB := uint8(8)
	numOfBuckets := 1 << defaultB

	buckets := make([]*bucket, numOfBuckets)

	for i := 0; i < numOfBuckets; i++ {
		buckets[i] = &bucket{}
	}

	return &GoMap{
		B:                defaultB,
		lowOrderBitsMask: 1<<bucketCnt - 1,
		buckets:          buckets,
	}
}

type GoMap struct {
	B                uint8 // log_2 of # of buckets
	lowOrderBitsMask uint64
	buckets          []*bucket
}

type bucket struct {
	count    uint8
	topHash  [bucketCnt]uint8
	keyArr   [bucketCnt]string
	valueArr [bucketCnt]string
}

func genHash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func (m *GoMap) Add(key string, value string) {
	hash := genHash(key)
	topHash := uint8(hash >> (64 - 8))
	bucket := m.buckets[hash&m.lowOrderBitsMask]

	var count uint8

	if i, ok := m.KeyExists(key); !ok {
		count = (*bucket).count
		(*bucket).count++
	} else {
		count = i
	}

	(*bucket).topHash[count] = topHash
	(*bucket).keyArr[count] = key
	(*bucket).valueArr[count] = value
}

func (m *GoMap) KeyExists(key string) (uint8, bool) {
	hash := genHash(key)
	topHash := uint8(hash >> (64 - 8))
	bucket := m.buckets[hash&m.lowOrderBitsMask]

	for k, tophash := range (*bucket).topHash {
		if topHash != tophash {
			continue
		}

		if key != (*bucket).keyArr[k] {
			continue
		}

		return uint8(k), true
	}

	return 0, false
}

func (m *GoMap) Get(key string) string {
	hash := genHash(key)
	topHash := uint8(hash >> (64 - 8))
	bucket := m.buckets[hash&m.lowOrderBitsMask]

	for k, tophash := range (*bucket).topHash {
		if topHash != tophash {
			continue
		}

		if key != (*bucket).keyArr[k] {
			continue
		}

		return (*bucket).valueArr[k]
	}

	return ""
}
