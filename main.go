package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {
	test := NewHashmap[int, string](1234, 10)
	test.Set(111, "222")
	value, err := test.Get(111)
	if err != nil {
		fmt.Println(111)
		panic(err)
	}
	if value != "222" {
		panic("")
	}

	h := NewHashmap[int, string](1024, 1)
	for i := 0; i < 1e6; i++ {
		h.Set(i, fmt.Sprintf("%d", i))
	}
	for i := 0; i < 1e6; i++ {
		value, err = h.Get(i)
		if err != nil {
			println(i)
			panic(err)
		}
		if value != fmt.Sprintf("%d", i) {
			panic("")
		}
	}

	h = NewHashmap[int, string](1024, 0.3)
	start := time.Now()
	for i := 0; i < 1e3; i++ {
		h.Set(i, "")
	}
	for i := 0; i < 1e3; i++ {
		h.Get(i)
	}
	fmt.Printf("1e3: %dns\n", time.Since(start)/1e3/time.Nanosecond)

	h = NewHashmap[int, string](1024, 0.3)
	start = time.Now()
	for i := 0; i < 1e4; i++ {
		h.Set(i, "")
	}
	for i := 0; i < 1e4; i++ {
		h.Get(i)
	}
	fmt.Printf("1e4: %dns\n", time.Since(start)/1e4/time.Nanosecond)

	h = NewHashmap[int, string](1024, 0.3)
	start = time.Now()
	for i := 0; i < 1e5; i++ {
		h.Set(i, "")
	}
	for i := 0; i < 1e5; i++ {
		h.Get(i)
	}
	fmt.Printf("1e5: %dns\n", time.Since(start)/1e5/time.Nanosecond)

	h = NewHashmap[int, string](1024, 0.3)
	start = time.Now()
	for i := 0; i < 1e6; i++ {
		h.Set(i, "")
	}
	for i := 0; i < 1e6; i++ {
		h.Get(i)
	}
	fmt.Printf("1e6: %dns\n", time.Since(start)/1e6/time.Nanosecond)

	h = NewHashmap[int, string](1024, 0.3)
	start = time.Now()
	for i := 0; i < 1e7; i++ {
		h.Set(i, "")
	}
	for i := 0; i < 1e7; i++ {
		h.Get(i)
	}
	fmt.Printf("1e7: %dns\n", time.Since(start)/1e7/time.Nanosecond)

	nativeMap := make(map[int64]string, 1024)
	start = time.Now()
	for i := 0; i < 1e3; i++ {
		nativeMap[int64(i)] = ""
	}
	for i := 0; i < 1e3; i++ {
		_ = nativeMap[int64(i)]
	}
	fmt.Printf("native 1e3: %dns\n", time.Since(start)/1e3/time.Nanosecond)

	nativeMap = make(map[int64]string, 1024)
	start = time.Now()
	for i := 0; i < 1e4; i++ {
		nativeMap[int64(i)] = ""
	}
	for i := 0; i < 1e4; i++ {
		_ = nativeMap[int64(i)]
	}
	fmt.Printf("native 1e4: %dns\n", time.Since(start)/1e4/time.Nanosecond)

	nativeMap = make(map[int64]string, 1024)
	start = time.Now()
	for i := 0; i < 1e5; i++ {
		nativeMap[int64(i)] = ""
	}
	for i := 0; i < 1e5; i++ {
		_ = nativeMap[int64(i)]
	}
	fmt.Printf("native 1e5: %dns\n", time.Since(start)/1e5/time.Nanosecond)

	nativeMap = make(map[int64]string, 1024)
	start = time.Now()
	for i := 0; i < 1e6; i++ {
		nativeMap[int64(i)] = ""
	}
	for i := 0; i < 1e6; i++ {
		_ = nativeMap[int64(i)]
	}
	fmt.Printf("native 1e6: %dns\n", time.Since(start)/1e6/time.Nanosecond)

	nativeMap = make(map[int64]string, 1024)
	start = time.Now()
	for i := 0; i < 1e7; i++ {
		nativeMap[int64(i)] = ""
	}
	for i := 0; i < 1e7; i++ {
		_ = nativeMap[int64(i)]
	}
	fmt.Printf("native 1e7: %dns\n", time.Since(start)/1e7/time.Nanosecond)

}

func Test() {
}

type Node[K comparable, V any] struct {
	key   K
	value V
	hash  int64
	next  *Node[K, V]
}

func NewNode[K comparable, V any](key K, value V, hash int64) *Node[K, V] {
	return &Node[K, V]{
		key:   key,
		value: value,
		hash:  hash,
		next:  nil,
	}
}

type HashMap[K comparable, V any] struct {
	table      []*Node[K, V]
	capacity   int64
	size       int64
	loadFactor float64
	bound      int64
}

func NewHashmap[K comparable, V any](capacity int64, loadFactor float64) *HashMap[K, V] {
	return &HashMap[K, V]{
		table:      make([]*Node[K, V], capacity),
		capacity:   capacity,
		size:       0,
		loadFactor: loadFactor,
		bound:      int64(float64(capacity) * loadFactor),
	}
}

func (h *HashMap[K, V]) Get(key K) (value V, err error) {
	hash := hash(key)
	node := h.table[hash%h.capacity]
	for {
		if node == nil {
			return value, errors.New("not found")
		}
		if node.key == key {
			return node.value, nil
		}
		node = node.next
	}
}

func (h *HashMap[K, V]) Set(key K, value V) {
	hash := hash(key)
	h.setByHash(key, value, hash)
}

func (h *HashMap[K, V]) setByHash(key K, value V, hash int64) {
	node := h.table[hash%h.capacity]
	if node == nil {
		h.table[hash%h.capacity] = NewNode(key, value, hash)
		h.size++
		if h.size > h.bound {
			h.Enlarge()
		}
		return
	}
	for {
		if node.key == key {
			node.value = value
			h.size++
			if h.size > h.bound {
				h.Enlarge()
			}
			return
		}
		if node.next == nil {
			node.next = NewNode(key, value, hash)
			h.size++
			if h.size > h.bound {
				h.Enlarge()
			}
			return
		}
		node = node.next
	}
}

func (h *HashMap[K, V]) Delete(key K) (err error) {
	hash := hash(key)
	node := h.table[hash%h.capacity]
	if node == nil {
		return errors.New("not found")
	}
	if node.key == key {
		h.table[hash%h.capacity] = nil
		return nil
	}
	for {
		if node.next == nil {
			return errors.New("not found")
		}
		if node.next.key == key {
			node.next = node.next.next
			return nil
		}
		node = node.next
	}
}

func (h *HashMap[K, V]) Enlarge() {
	newMap := NewHashmap[K, V](2*h.capacity, h.loadFactor)

	for i := 0; i < int(h.capacity); i++ {
		node := h.table[i]
		for {
			if node == nil {
				break
			}
			newMap.setByHash(node.key, node.value, node.hash)
			node = node.next
		}
	}
	h.table = newMap.table
	h.capacity *= 2
	h.bound *= 2
}

func hash[K comparable](key K) (hash int64) {
	keyString := fmt.Sprint(key)
	hash = 0
	for _, ch := range keyString {
		hash += int64(ch)
		hash += hash << 10
		hash ^= hash >> 6
	}

	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15

	if hash >= 0 {
		return hash
	} else {
		return -hash
	}
}
