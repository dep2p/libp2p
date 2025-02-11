// Package msgio 提供了带有长度前缀的消息读写功能的测试代码
package msgio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	str "strings"
	"sync"
	"testing"
	"time"
)

// randBuf 生成指定大小的随机字节切片
// 参数:
//   - r: 随机数生成器
//   - size: 要生成的字节切片大小
//
// 返回值:
//   - []byte: 生成的随机字节切片
func randBuf(r *rand.Rand, size int) []byte {
	buf := make([]byte, size)
	_, _ = r.Read(buf)
	return buf
}

// TestReadWrite 测试基本的读写功能
func TestReadWrite(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := NewWriter(buf)
	reader := NewReader(buf)
	SubtestReadWrite(t, writer, reader)
}

// TestReadWriteMsg 测试消息读写功能
func TestReadWriteMsg(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := NewWriter(buf)
	reader := NewReader(buf)
	SubtestReadWriteMsg(t, writer, reader)
}

// TestReadWriteMsgSync 测试并发消息读写功能
func TestReadWriteMsgSync(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := NewWriter(buf)
	reader := NewReader(buf)
	SubtestReadWriteMsgSync(t, writer, reader)
}

// TestReadClose 测试读取时关闭功能
func TestReadClose(t *testing.T) {
	r, w := io.Pipe()
	writer := NewWriter(w)
	reader := NewReader(r)
	SubtestReadClose(t, writer, reader)
}

// TestWriteClose 测试写入时关闭功能
func TestWriteClose(t *testing.T) {
	r, w := io.Pipe()
	writer := NewWriter(w)
	reader := NewReader(r)
	SubtestWriteClose(t, writer, reader)
}

// testIoReadWriter 实现了 io.Reader 和 io.Writer 接口的测试类型
type testIoReadWriter struct {
	io.Reader
	io.Writer
}

// TestReadWriterClose 测试读写器关闭功能
func TestReadWriterClose(t *testing.T) {
	r, w := io.Pipe()
	rw := NewReadWriter(testIoReadWriter{r, w})
	SubtestReaderWriterClose(t, rw)
}

// TestReadWriterCombine 测试读写器组合功能
func TestReadWriterCombine(t *testing.T) {
	r, w := io.Pipe()
	writer := NewWriter(w)
	reader := NewReader(r)
	rw := Combine(writer, reader)
	rw.Close()
}

// TestMultiError 测试多错误处理功能
func TestMultiError(t *testing.T) {
	emptyError := multiErr([]error{})
	if emptyError.Error() != "no errors" {
		t.Fatal("Expected no errors")
	}

	twoErrors := multiErr([]error{errors.New("one"), errors.New("two")})
	if eStr := twoErrors.Error(); !str.Contains(eStr, "one") && !str.Contains(eStr, "two") {
		t.Fatal("Expected error messages not included")
	}
}

// TestShortBufferError 测试缓冲区过短错误
func TestShortBufferError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := NewWriter(buf)
	reader := NewReader(buf)
	SubtestReadShortBuffer(t, writer, reader)
}

// SubtestReadWrite 测试读写功能的子测试
// 参数:
//   - t: 测试对象
//   - writer: 写入器
//   - reader: 读取器
func SubtestReadWrite(t *testing.T, writer WriteCloser, reader ReadCloser) {
	msgs := [1000][]byte{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range msgs {
		msgs[i] = randBuf(r, r.Intn(1000))
		n, err := writer.Write(msgs[i])
		if err != nil {
			t.Fatal(err)
		}
		if n != len(msgs[i]) {
			t.Fatal("wrong length:", n, len(msgs[i]))
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	for i := 0; ; i++ {
		msg2 := make([]byte, 1000)
		n, err := reader.Read(msg2)
		if err != nil {
			if err == io.EOF {
				if i < len(msg2) {
					t.Error("failed to read all messages", len(msgs), i)
				}
				break
			}
			t.Error("unexpected error", err)
		}

		msg1 := msgs[i]
		msg2 = msg2[:n]
		if !bytes.Equal(msg1, msg2) {
			t.Fatal("message retrieved not equal\n", msg1, "\n\n", msg2)
		}
	}

	if err := reader.Close(); err != nil {
		t.Error(err)
	}
}

// SubtestReadWriteMsg 测试消息读写功能的子测试
// 参数:
//   - t: 测试对象
//   - writer: 写入器
//   - reader: 读取器
func SubtestReadWriteMsg(t *testing.T, writer WriteCloser, reader ReadCloser) {
	msgs := [1000][]byte{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range msgs {
		msgs[i] = randBuf(r, r.Intn(1000))
		err := writer.WriteMsg(msgs[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	for i := 0; ; i++ {
		msg2, err := reader.ReadMsg()
		if err != nil {
			if err == io.EOF {
				if i < len(msg2) {
					t.Error("failed to read all messages", len(msgs), i)
				}
				break
			}
			t.Error("unexpected error", err)
		}

		msg1 := msgs[i]
		if !bytes.Equal(msg1, msg2) {
			t.Fatal("message retrieved not equal\n", msg1, "\n\n", msg2)
		}
	}

	if err := reader.Close(); err != nil {
		t.Error(err)
	}
}

// SubtestReadWriteMsgSync 测试并发消息读写功能的子测试
// 参数:
//   - t: 测试对象
//   - writer: 写入器
//   - reader: 读取器
func SubtestReadWriteMsgSync(t *testing.T, writer WriteCloser, reader ReadCloser) {
	msgs := [1000][]byte{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range msgs {
		msgs[i] = randBuf(r, r.Intn(1000)+4)
		NBO.PutUint32(msgs[i][:4], uint32(i))
	}

	var wg1 sync.WaitGroup
	var wg2 sync.WaitGroup

	errs := make(chan error, 10000)
	for i := range msgs {
		wg1.Add(1)
		go func(i int) {
			defer wg1.Done()

			err := writer.WriteMsg(msgs[i])
			if err != nil {
				errs <- err
			}
		}(i)
	}

	wg1.Wait()
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(msgs)+1; i++ {
		wg2.Add(1)
		go func(i int) {
			defer wg2.Done()

			msg2, err := reader.ReadMsg()
			if err != nil {
				if err == io.EOF {
					if i < len(msg2) {
						errs <- fmt.Errorf("failed to read all messages %d %d", len(msgs), i)
					}
					return
				}
				errs <- fmt.Errorf("unexpected error: %s", err)
			}

			mi := NBO.Uint32(msg2[:4])
			msg1 := msgs[mi]
			if !bytes.Equal(msg1, msg2) {
				errs <- fmt.Errorf("message retrieved not equal\n%s\n\n%s", msg1, msg2)
			}
		}(i)
	}

	wg2.Wait()
	close(errs)

	if err := reader.Close(); err != nil {
		t.Error(err)
	}

	for e := range errs {
		t.Error(e)
	}
}

// TestBadSizes 测试错误大小的处理
func TestBadSizes(t *testing.T) {
	data := make([]byte, 4)

	// 在64位系统上,这会因为太大而失败
	// 在32位系统上,这会因为太小而失败
	NBO.PutUint32(data, 4000000000)
	buf := bytes.NewReader(data)
	read := NewReader(buf)
	msg, err := read.ReadMsg()
	if err == nil {
		t.Fatal(err)
	}
	_ = msg
}

// SubtestReadClose 测试读取时关闭功能的子测试
// 参数:
//   - t: 测试对象
//   - writer: 写入器
//   - reader: 读取器
func SubtestReadClose(t *testing.T, writer WriteCloser, reader ReadCloser) {
	defer writer.Close()

	buf := [10]byte{}
	done := make(chan struct{})
	go func() {
		defer close(done)
		time.Sleep(10 * time.Millisecond)
		reader.Close()
	}()
	n, err := reader.Read(buf[:])
	if n != 0 || err == nil {
		t.Error("expected to read nothing")
	}
	<-done
}

// SubtestWriteClose 测试写入时关闭功能的子测试
// 参数:
//   - t: 测试对象
//   - writer: 写入器
//   - reader: 读取器
func SubtestWriteClose(t *testing.T, writer WriteCloser, reader ReadCloser) {
	defer reader.Close()

	buf := [10]byte{}
	done := make(chan struct{})
	go func() {
		defer close(done)
		time.Sleep(10 * time.Millisecond)
		writer.Close()
	}()
	n, err := writer.Write(buf[:])
	if n != 0 || err == nil {
		t.Error("expected to write nothing")
	}
	<-done
}

// SubtestReaderWriterClose 测试读写器关闭功能的子测试
// 参数:
//   - t: 测试对象
//   - rw: 读写器
func SubtestReaderWriterClose(t *testing.T, rw ReadWriteCloser) {
	buf := [10]byte{}
	done := make(chan struct{})
	go func() {
		defer close(done)
		time.Sleep(10 * time.Millisecond)
		buf := [10]byte{}
		rw.Read(buf[:])
		rw.Close()
	}()
	n, err := rw.Write(buf[:])
	if n != 10 || err != nil {
		t.Error("Expected to write 10 bytes")
	}
	<-done
}

// SubtestReadShortBuffer 测试缓冲区过短错误的子测试
// 参数:
//   - t: 测试对象
//   - writer: 写入器
//   - reader: 读取器
func SubtestReadShortBuffer(t *testing.T, writer WriteCloser, reader ReadCloser) {
	defer reader.Close()
	shortReadBuf := [1]byte{}
	done := make(chan struct{})

	go func() {
		defer writer.Close()
		defer close(done)
		time.Sleep(10 * time.Millisecond)
		largeWriteBuf := [10]byte{}
		writer.Write(largeWriteBuf[:])
	}()
	<-done
	n, _ := reader.NextMsgLen()
	if n != 10 {
		t.Fatal("Expected next message to have length of 10")
	}
	_, err := reader.Read(shortReadBuf[:])
	if err != io.ErrShortBuffer {
		t.Fatal("Expected short buffer error")
	}
}
