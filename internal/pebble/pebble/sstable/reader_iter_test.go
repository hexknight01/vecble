// Copyright 2024 The LevelDB-Go and Pebble Authors. All rights reserved. Use
// of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package sstable

import (
	"context"
	"math/rand/v2"
	"testing"

	"github.com/cockroachdb/crlib/testutils/leaktest"
	"github.com/cockroachdb/pebble/v2/internal/base"
	"github.com/cockroachdb/pebble/v2/internal/testutils"
	"github.com/cockroachdb/pebble/v2/objstorage/objstorageprovider"
	"github.com/cockroachdb/pebble/v2/sstable/block"
	"github.com/cockroachdb/pebble/v2/vfs"
	"github.com/stretchr/testify/require"
)

// TestIteratorErrorOnInit tests the path where creation of an iterator fails
// when reading the index block.
func TestIteratorErrorOnInit(t *testing.T) {
	defer leaktest.AfterTest(t)()
	mem := vfs.NewMem()

	f0, err := mem.Create("test.sst", vfs.WriteCategoryUnspecified)
	require.NoError(t, err)
	writerOpts := WriterOptions{
		Comparer:   base.DefaultComparer,
		MergerName: base.DefaultMerger.Name,
	}
	w := NewWriter(objstorageprovider.NewFileWritable(f0), writerOpts)
	require.NoError(t, w.Set([]byte("test"), nil))
	require.NoError(t, w.Close())

	f1 := testutils.CheckErr(mem.Open("test.sst"))
	r, err := newReader(f1, ReaderOptions{
		Comparer: base.DefaultComparer,
		Merger:   base.DefaultMerger,
	})
	require.NoError(t, err)
	defer r.Close()

	// Swap the readable in the reader.
	bad := testutils.CheckErr(mem.Create("bad.sst", vfs.WriteCategoryUnspecified))
	require.NoError(t, bad.Close())
	bad = testutils.CheckErr(mem.Open("bad.sst"))
	saveReadable := r.readable
	r.readable = testutils.CheckErr(NewSimpleReadable(bad))

	var pool block.BufferPool
	pool.Init(5)
	defer pool.Release()

	var stats base.InternalIteratorStats
	for k := 0; k < 20; k++ {
		if rand.IntN(2) == 0 {
			_, err := newRowBlockSingleLevelIterator(
				context.Background(),
				r,
				nil, /* v */
				NoTransforms,
				nil /* lower */, nil, /* upper */
				nil /* filterer */, NeverUseFilterBlock,
				&stats,
				nil, /* statsAccum */
				MakeTrivialReaderProvider(r),
				&pool,
			)
			require.Error(t, err)
		} else {
			_, err := newRowBlockTwoLevelIterator(
				context.Background(),
				r,
				nil, /* v */
				NoTransforms,
				nil /* lower */, nil, /* upper */
				nil /* filterer */, NeverUseFilterBlock,
				&stats,
				nil, /* statsAccum */
				MakeTrivialReaderProvider(r),
				&pool,
			)
			require.Error(t, err)
		}
	}
	require.NoError(t, r.readable.Close())
	r.readable = saveReadable
}
