package rpc

type byteBuffer struct {
	buffers   []byte
	length    int
	readCount int
}

func newByteBuffer(length int) byteBuffer {
	return byteBuffer{
		buffers:   make([]byte, length, length),
		length:    length,
		readCount: 0,
	}
}

func (b *byteBuffer) Read(buffer []byte, count int, offset int) int {
	remain := b.length - b.readCount
	reading := count - offset
	//if reading >= remain {
	//    copy(b.buffers[b.readCount:b.readCount+remain], buffer[offset:offset+remain])
	//    b.readCount += remain
	//    return remain
	//}
	//copy(b.buffers[b.readCount:b.readCount+reading], buffer[offset:offset+reading])
	//b.readCount += reading
	//return reading

	readLength := 0
	if reading >= remain {
		readLength = remain
	} else {
		readLength = reading
	}
	copy(b.buffers[b.readCount:b.readCount+readLength], buffer[offset:offset+readLength])
    b.readCount += readLength
	return readLength
}

func (b byteBuffer) Done() bool {
	return b.readCount == b.length
}

func (b *byteBuffer) Dispose() {
	b.buffers = nil
}
