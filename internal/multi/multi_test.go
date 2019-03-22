package multi

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func TestNewMapWriter(t *testing.T) {
	writer := NewMapWriter()
	if writer == nil {
		t.Error("Unexpected: nil")
	}
	_, ok := interface{}(writer).(io.Writer)
	if !ok {
		t.Error("Unexpected: not an io.Writer")
	}
	_, ok = interface{}(writer).(*MapWriter)
	if !ok {
		t.Error("Unexpected: not a *MapWriter")
	}
}

func TestSet(t *testing.T) {
	writer := NewMapWriter()
	var (
		id1     = "id1"
		id2     = "id2"
		buffer1 bytes.Buffer
		buffer2 bytes.Buffer
	)
	size := writer.Set(id1, &buffer1)
	if size != 1 {
		t.Errorf("Unexpected map size: %d. Expected: %d", size, 1)
	}
	size = writer.Set(id2, &buffer2)
	if size != 2 {
		t.Errorf("Unexpected map size: %d. Expected: %d", size, 2)
	}
}

func TestDelete(t *testing.T) {
	writer := NewMapWriter()
	var (
		id1     = "id1"
		id2     = "id2"
		buffer1 bytes.Buffer
		buffer2 bytes.Buffer
	)
	writer.Set(id1, &buffer1)
	writer.Set(id2, &buffer2)
	size := writer.Delete(id1)
	if size != 1 {
		t.Errorf("Unexpected map size: %d. Expected: %d", size, 1)
	}
	size = writer.Delete(id2)
	if size != 0 {
		t.Errorf("Unexpected map size: %d. Expected: %d", size, 0)
	}
}

func TestSize(t *testing.T) {
	writer := NewMapWriter()
	var (
		id1     = "id1"
		id2     = "id2"
		buffer1 bytes.Buffer
		buffer2 bytes.Buffer
	)
	size := writer.Size()
	if size != 0 {
		t.Errorf("Unexpected map size: %d. Expected: %d", size, 0)
	}
	writer.Set(id1, &buffer1)
	size = writer.Size()
	if size != 1 {
		t.Errorf("Unexpected map size: %d. Expected: %d", size, 1)
	}
	writer.Set(id2, &buffer2)
	size = writer.Size()
	if size != 2 {
		t.Errorf("Unexpected map size: %d. Expected: %d", size, 2)
	}
	writer.Delete(id1)
	size = writer.Size()
	if size != 1 {
		t.Errorf("Unexpected map size: %d. Expected: %d", size, 1)
	}
	writer.Delete(id2)
	size = writer.Size()
	if size != 0 {
		t.Errorf("Unexpected map size: %d. Expected: %d", size, 0)
	}
}

func TestWrite(t *testing.T) {
	writer := NewMapWriter()
	var (
		id1     = "id1"
		id2     = "id2"
		buffer1 bytes.Buffer
		buffer2 bytes.Buffer
	)
	writer.Set(id1, &buffer1)
	writer.Set(id2, &buffer2)
	writer.Write([]byte("banana"))
	output1, _ := ioutil.ReadAll(&buffer1)
	if string(output1) != "banana" {
		t.Errorf("Unexpected output: %s. Expected: %s", output1, "banana")
	}
	output2, _ := ioutil.ReadAll(&buffer2)
	if string(output2) != "banana" {
		t.Errorf("Unexpected output: %s. Expected: %s", output2, "banana")
	}
	writer.Delete(id1)
	writer.Write([]byte("banana"))
	output1, _ = ioutil.ReadAll(&buffer1)
	if string(output1) != "" {
		t.Errorf("Unexpected output: %s. Expected: %s", output1, "")
	}
	output2, _ = ioutil.ReadAll(&buffer2)
	if string(output2) != "banana" {
		t.Errorf("Unexpected output: %s. Expected: %s", output2, "banana")
	}
}

func benchmarkWrite(b *testing.B, numWriters int, numBytes int) {
	writer := NewMapWriter()
	for n := 0; n < numWriters; n++ {
		writer.Set("id"+string(n), ioutil.Discard)
	}
	data := make([]byte, numBytes)
	for n := 0; n < b.N; n++ {
		writer.Write(data)
	}
}

func BenchmarkWrite_64_1K(b *testing.B)  { benchmarkWrite(b, 64, 1024) }
func BenchmarkWrite_64_1M(b *testing.B)  { benchmarkWrite(b, 64, 1024^2) }
func BenchmarkWrite_64_1G(b *testing.B)  { benchmarkWrite(b, 64, 1024^3) }
func BenchmarkWrite_256_1K(b *testing.B) { benchmarkWrite(b, 256, 1024) }
func BenchmarkWrite_256_1M(b *testing.B) { benchmarkWrite(b, 256, 1024^2) }
func BenchmarkWrite_256_1G(b *testing.B) { benchmarkWrite(b, 256, 1024^3) }
func BenchmarkWrite_1K_1K(b *testing.B)  { benchmarkWrite(b, 1024, 1024) }
func BenchmarkWrite_1K_1M(b *testing.B)  { benchmarkWrite(b, 1024, 1024^2) }
func BenchmarkWrite_1K_1G(b *testing.B)  { benchmarkWrite(b, 1024, 1024^3) }
