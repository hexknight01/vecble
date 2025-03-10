// Copyright 2023 The LevelDB-Go and Pebble Authors. All rights reserved. Use
// of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package block

import (
	"encoding/binary"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/cockroachdb/crlib/testutils/leaktest"
	"github.com/cockroachdb/pebble/v2/internal/cache"
	"github.com/stretchr/testify/require"
)

func TestCompressionRoundtrip(t *testing.T) {
	defer leaktest.AfterTest(t)()

	seed := uint64(time.Now().UnixNano())
	t.Logf("seed %d", seed)
	rng := rand.New(rand.NewPCG(0, seed))

	for compression := DefaultCompression + 1; compression < NCompression; compression++ {
		t.Run(compression.String(), func(t *testing.T) {
			payload := make([]byte, 1+rng.IntN(10<<10 /* 10 KiB */))
			for i := range payload {
				payload[i] = byte(rng.Uint32())
			}
			// Create a randomly-sized buffer to house the compressed output. If it's
			// not sufficient, Compress should allocate one that is.
			compressedBuf := make([]byte, 1+rng.IntN(1<<10 /* 1 KiB */))

			btyp, compressed := compress(compression, payload, compressedBuf)
			v, err := decompress(btyp, compressed)
			require.NoError(t, err)
			got := payload
			if v != nil {
				got = v.Buf()
				require.Equal(t, payload, got)
				cache.Free(v)
			}
		})
	}
}

// TestDecompressionError tests that a decompressing a value that does not
// decompress returns an error.
func TestDecompressionError(t *testing.T) {
	defer leaktest.AfterTest(t)()
	rng := rand.New(rand.NewPCG(0, 1 /* fixed seed */))

	// Create a buffer to represent a faux zstd compressed block. It's prefixed
	// with a uvarint of the appropriate length, followed by garabge.
	fauxCompressed := make([]byte, rng.IntN(10<<10 /* 10 KiB */))
	compressedPayloadLen := len(fauxCompressed) - binary.MaxVarintLen64
	n := binary.PutUvarint(fauxCompressed, uint64(compressedPayloadLen))
	fauxCompressed = fauxCompressed[:n+compressedPayloadLen]
	for i := range fauxCompressed[:n] {
		fauxCompressed[i] = byte(rng.Uint32())
	}

	v, err := decompress(ZstdCompressionIndicator, fauxCompressed)
	t.Log(err)
	require.Error(t, err)
	require.Nil(t, v)
}

// decompress decompresses an sstable block into memory manually allocated with
// `cache.Alloc`.  NB: If Decompress returns (nil, nil), no decompression was
// necessary and the caller may use `b` directly.
func decompress(algo CompressionIndicator, b []byte) (*cache.Value, error) {
	if algo == NoCompressionIndicator {
		return nil, nil
	}
	// first obtain the decoded length.
	decodedLen, prefixLen, err := DecompressedLen(algo, b)
	if err != nil {
		return nil, err
	}
	b = b[prefixLen:]
	// Allocate sufficient space from the cache.
	decoded := cache.Alloc(decodedLen)
	decodedBuf := decoded.Buf()
	if err := DecompressInto(algo, b, decodedBuf); err != nil {
		cache.Free(decoded)
		return nil, err
	}
	return decoded, nil
}
