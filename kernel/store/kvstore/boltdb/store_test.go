//  Copyright (c) 2014 Couchbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package boltdb

import (
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/tiglabs/baud/kernel/store/kvstore"
	"github.com/tiglabs/baud/kernel/store/kvstore/test"
)

func open(t *testing.T) kvstore.KVStore {
	rv, err := New(&StoreConfig{Path: "test"})
	if err != nil {
		t.Fatal(err)
	}
	return rv
}

func cleanup(t *testing.T, s kvstore.KVStore) {
	err := s.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = os.RemoveAll("test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestBoltDBKVCrud(t *testing.T) {
	s := open(t)
	defer cleanup(t, s)
	test.CommonTestKVCrud(t, s)
}

func TestBoltDBReaderIsolation(t *testing.T) {
	s := open(t)
	defer cleanup(t, s)
	test.CommonTestReaderIsolation(t, s)
}

func TestBoltDBReaderOwnsGetBytes(t *testing.T) {
	s := open(t)
	defer cleanup(t, s)
	test.CommonTestReaderOwnsGetBytes(t, s)
}

func TestBoltDBWriterOwnsBytes(t *testing.T) {
	s := open(t)
	defer cleanup(t, s)
	test.CommonTestWriterOwnsBytes(t, s)
}

func TestBoltDBPrefixIterator(t *testing.T) {
	s := open(t)
	defer cleanup(t, s)
	test.CommonTestPrefixIterator(t, s)
}

func TestBoltDBPrefixIteratorSeek(t *testing.T) {
	s := open(t)
	defer cleanup(t, s)
	test.CommonTestPrefixIteratorSeek(t, s)
}

func TestBoltDBRangeIterator(t *testing.T) {
	s := open(t)
	defer cleanup(t, s)
	test.CommonTestRangeIterator(t, s)
}

func TestBoltDBRangeIteratorSeek(t *testing.T) {
	s := open(t)
	defer cleanup(t, s)
	test.CommonTestRangeIteratorSeek(t, s)
}

func TestBoltDBConfig(t *testing.T) {
	var tests = []struct {
		in          *StoreConfig
		path        string
		bucket      string
		noSync      bool
		fillPercent float64
	}{
		{
			&StoreConfig{Path: "test", Bucket: "mybucket", NoSync: true, FillPercent: 0.75},
			"test",
			"mybucket",
			true,
			0.75,
		},
		{
			&StoreConfig{Path: "test"},
			"test",
			"baud",
			false,
			bolt.DefaultFillPercent,
		},
	}

	for _, test := range tests {
		kv, err := New(test.in)
		if err != nil {
			t.Fatal(err)
		}
		bs, ok := kv.(*Store)
		if !ok {
			t.Fatal("failed type assertion to *boltdb.Store")
		}
		if bs.path != test.path {
			t.Fatalf("path: expected %q, got %q", test.path, bs.path)
		}
		if string(bs.bucket) != test.bucket {
			t.Fatalf("bucket: expected %q, got %q", test.bucket, bs.bucket)
		}
		if bs.noSync != test.noSync {
			t.Fatalf("noSync: expected %t, got %t", test.noSync, bs.noSync)
		}
		if bs.fillPercent != test.fillPercent {
			t.Fatalf("fillPercent: expected %f, got %f", test.fillPercent, bs.fillPercent)
		}
		cleanup(t, kv)
	}
}
