/*
Copyright (c) 2015, Simon J Mudd
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// This is to hold a cached set of digest of queries
package querycache

import (
	"log"
	"sync"

	"github.com/sjmudd/queryprofiler/querykey"
)

var global *QueryCache

func init() {
	global = NewQueryCache()
}

type QueryCache struct {
	sync.RWMutex
	cache map[string]string
}

func NewQueryCache() *QueryCache {
	n := new(QueryCache)

	n.cache = make(map[string]string)

	return n
}

// Put saves a query to cache.
func (qkc *QueryCache) Put(aQueryKey querykey.QueryKey, query string) {
	qkc.RLock()
	old, found := qkc.cache[aQueryKey.String()]
	qkc.RUnlock()
	if !found {
		qkc.Lock()
		qkc.cache[aQueryKey.String()] = query
		qkc.Unlock()
	} else {
		// validate we don't try to save a _different_ value
		if old != query {
			log.Println("qks.Put()")
			log.Println("aQueryKey:", aQueryKey)
			log.Println("query:", query)
			log.Println("old value for existing digest:", old)
			log.Fatal("FATAL")
		}
	}
}

// Get a query by by the digest from the cache
func (qkc QueryCache) Get(aQueryKey querykey.QueryKey) (string, bool) {
	qkc.RLock()
	result, found := qkc.cache[aQueryKey.String()]
	qkc.RUnlock()

	return result, found
}

// saves to global cache
func Put(aQueryKey querykey.QueryKey, query string) {
	if global == nil {
		log.Panic("Put() global == nil")
	}
	global.Put(aQueryKey, query)
}

// Get a query by by the digest from the cache
func Get(aQueryKey querykey.QueryKey) (string, bool) {
	if global == nil {
		log.Panic("Get() global == nil")
	}
	return global.Get(aQueryKey)
}
