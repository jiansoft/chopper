package chopper

import (
	"sync"

	"github.com/jiansoft/robin"
)

// This struct creates a single large buffer which can be divided up
// and assigned to peer struct for use with each
// socket I/O operation.
// This enables bufffers to be easily reused and guards against
// fragmenting heap memory.
type BufferManager struct {
	buffers      []byte
	bufferSize   int
	totalBytes   int64
	currentIndex int64
	freeIndex    *robin.ConcurrentQueue
	lock         *sync.Mutex
}

func (b *BufferManager) init(maxConnectCount int, bufferSize int) *BufferManager {
	//Size of the underlying Byte array.
	b.bufferSize = bufferSize
	//The total number of bytes controlled by the buffer pool.
	b.totalBytes = int64(maxConnectCount) * int64(b.bufferSize)
	//The underlying Byte array maintained by the Buffer Manager.
	b.buffers = make([]byte, b.totalBytes, b.totalBytes)
	//Current index of the underlying Byte array.
	b.currentIndex = 0
	//Pool of indexes for the Buffer Manager.
	b.freeIndex = robin.NewConcurrentQueue()
	b.lock = new(sync.Mutex)
	return b
}

func NewBufferManager(maxConnectCount int, bufferSize int) *BufferManager {
	return new(BufferManager).init(maxConnectCount, bufferSize)
}

// Assigns a buffer from the buffer pool to the
// specified peer struct
func (b *BufferManager) SetBuffer(peer *PeerSession) {
	b.lock.Lock()
	defer b.lock.Unlock()
	offset, ok := b.freeIndex.TryDequeue()
	if ok {
		peer.bufferOffst = offset.(int64)
		peer.buffers = b.buffers[peer.bufferOffst : peer.bufferOffst+int64(b.bufferSize)]
	} else {
		if b.totalBytes-int64(b.bufferSize) < b.currentIndex {
			peer.buffers = make([]byte, b.bufferSize)
			peer.bufferOffst = -1
			//The buffer pool is empty.
			//return false
		} else {
			peer.bufferOffst = b.currentIndex
			peer.buffers = b.buffers[peer.bufferOffst : peer.bufferOffst+int64(b.bufferSize)]
			b.currentIndex += int64(b.bufferSize)
		}
	}
	//return true
}

// Removes the buffer from a peer struct.
// This frees the buffer back to the buffer pool
func (b *BufferManager) FreeBuffer(peer *PeerSession) {
	peer.buffers = nil
	b.freeIndex.Enqueue(peer.bufferOffst)
}
