package chopper

type TcpBinaryHander struct {
	buffer       byteBuffer
	callBackFunc func([]byte)
	parse        func([]byte, int, int) int
}

const headerSize = 4

func NewTcpBinaryHander(readDoneFunc func([]byte)) *TcpBinaryHander {
	obj := new(TcpBinaryHander)
	obj.restart()
	obj.callBackFunc = readDoneFunc
	return obj
}

// Wrap the data length into same package
func (TcpBinaryHander) Wrap(data []byte) []byte {
	totalLen := len(data)
	size := []byte{byte(totalLen >> 24), byte(totalLen >> 16), byte(totalLen >> 8), byte(totalLen)}
	return append(size, data...)
}

// Parse the received data
func (t *TcpBinaryHander) Parse(buffer []byte, count int) {
	for offset := 0; offset < count; {
		offset += t.parse(buffer, count, offset)
	}
}

func (t *TcpBinaryHander) Dispose() {
	t.parse = nil
	t.callBackFunc = nil
	t.buffer.Dispose()
}

func (t *TcpBinaryHander) restart() {
	t.buffer = newByteBuffer(headerSize)
	t.parse = t.parsePacketLength
}

func (t *TcpBinaryHander) parsePacketLength(buffer []byte, count int, offset int) int {
	var index = t.buffer.Read(buffer, count, offset)
	if !t.buffer.Done() {
		return index
	}
	size := int(t.buffer.buffers[0])<<24 | int(t.buffer.buffers[1])<<16 | int(t.buffer.buffers[2])<<8 | int(t.buffer.buffers[3])
	t.buffer = newByteBuffer(size)
	t.parse = t.parsePacketBody
	return index
}

func (t *TcpBinaryHander) parsePacketBody(buffer []byte, count int, offset int) int {
	index := t.buffer.Read(buffer, count, offset)
	if !t.buffer.Done() {
		return index
	}
	t.callBackFunc(t.buffer.buffers)
	t.restart()
	return index
}
