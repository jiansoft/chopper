package rpc

import (
	"testing"
)

// go test -v -bench=".*"
func TestAdd(t *testing.T) {
	source := make([]byte, 10)
	for q := 0; q < len(source); q++ {
		source[q] = byte(q)
	}
	t.Logf("原始資料 %d %v", len(source), source)
	totalLen := 4 + len(source)
	sendData := make([]byte, totalLen)
	sendData[0] = byte((totalLen >> 24 & 255))
	sendData[1] = byte(totalLen >> 16 & 255)
	sendData[2] = byte(totalLen >> 8 & 255)
	sendData[3] = byte(totalLen & 255)
	copy(sendData[4:], source)
	t.Logf("發送資料 %d %v", len(sendData), sendData)

	reader := newTcpBinaryHander()
	reader.onComplete = onReceived
	reader.Parse(sendData, len(sendData))

	for q := 4; q < len(sendData); q++ {
		sendData[q] = byte(q + 100)
	}
	t.Logf("發送資料 %d %v", len(sendData), sendData)
	t.Logf("第一段發送資料 %d %v", len(sendData[0:13]), sendData[0:13])
	reader.Parse(sendData[0:13], len(sendData[0:13]))

	source = make([]byte, 15)
	source[0] = 113
	copy(source[1:5], sendData[0:4])
	reverse := reverse(sendData[4:14])
	copy(source[5:], reverse)

	t.Logf("第二段發送資料 %d %v", len(source), source)
	reader.Parse(source, len(source))

	//if Add(1, 2) == 3 {
	//    t.Log("mymath.Add PASS")
	//} else {
	//    t.Error("mymath.Add FAIL")
	//}
}

func BenchmarkTcpBinaryHander(b *testing.B) {
	reader := newTcpBinaryHander()
	reader.onComplete = onReceived

	b.ResetTimer()
	//b.Log("Benchmark_TcpBinaryHander run:",b.N)
	for i := 0; i < b.N; i++ {
		source := make([]byte, 10)
		for q := 0; q < len(source); q++ {
			source[q] = byte(q)
		}
		//b.Logf("原始資料 %d %v", size(source), source)
		totalLen := 4 + len(source)
		sendData := make([]byte, totalLen)
		sendData[0] = byte((totalLen >> 24 & 255))
		sendData[1] = byte(totalLen >> 16 & 255)
		sendData[2] = byte(totalLen >> 8 & 255)
		sendData[3] = byte(totalLen & 255)
		copy(sendData[4:], source)
		//b.Logf("發送資料 %d %v", size(sendData), sendData)

		reader.Parse(sendData, len(sendData))

		for q := 4; q < len(sendData); q++ {
			sendData[q] = byte(q + 100)
		}
		//b.Logf("發送資料 %d %v", size(sendData), sendData)
		//b.Logf("第一段發送資料 %d %v", size(sendData[0:13]), sendData[0:13])
		reader.Parse(sendData[0:13], len(sendData[0:13]))

		source = make([]byte, 15)
		source[0] = 113
		copy(source[1:5], sendData[0:4])
		reverse := reverse(sendData[4:14])
		copy(source[5:], reverse)

		//b.Logf("第二段發送資料 %d %v", size(source), source)
		reader.Parse(source, len(source))
	}
	//b.Log("Benchmark_TcpBinaryHander PASS")
}

func onReceived(data []byte) {
	//log.Infof("資料收完 PASS %d %v", size(data), data)
}

func reverse(numbers []byte) []byte {
	newNumbers := make([]byte, len(numbers))
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		newNumbers[i], newNumbers[j] = numbers[j], numbers[i]
	}
	return newNumbers
}
