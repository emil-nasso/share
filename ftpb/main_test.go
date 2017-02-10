package ftpb

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAcceptanceNonDelayed(t *testing.T) {
	fmt.Println("Non delayed output acceptance test (should update on same line all the way to 100/100):")
	pb := New(10)
	pb.Start()
	for i := 0; i <= 10; i++ {
		pb.Increase(1)
	}
	pb.Done()
}

func TestAcceptanceDelayed(t *testing.T) {
	fmt.Println("Delayed output acceptance test (should update on same line all the way to 100/100):")
	pb := New(10)
	pb.Start()
	for i := 0; i <= 10; i++ {
		pb.Increase(1)
		time.Sleep(time.Millisecond * 10)
	}
	pb.Done()
}

func TestChangeInstance(t *testing.T) {
	pb := New(1000)
	buffer := bytes.NewBufferString("")

	pb.SetByteBuffer(buffer)
	assert.Equal(t, "", buffer.String())

	buffer.Reset()
	pb.Increase(5)
	assert.Equal(t, "5 B / 1.0 kB\n", buffer.String())

	buffer.Reset()
	pb.Increase(10)
	assert.Equal(t, "15 B / 1.0 kB\n", buffer.String())

	buffer.Reset()
	pb.Set(500)
	assert.Equal(t, "500 B / 1.0 kB\n", buffer.String())

	buffer.Reset()
	pb.Done()
	assert.Equal(t, "1.0 kB / 1.0 kB\n", buffer.String())

}

func TestAboveOneHundredPercent(t *testing.T) {
	pb := New(1000)
	buffer := bytes.NewBufferString("")

	pb.SetByteBuffer(buffer)

	pb.Increase(1001)
	assert.Equal(t, "1.0 kB / 1.0 kB\n", buffer.String())
	pb.Done()
}

func testOutputFormat(t *testing.T) {
	assertProgressOutput(t, 2500, 10000, "2.5 kB / 10.0 kB")
	assertProgressOutput(t, 2500000, 1000000, "2.5 MB / 1.0 MB")
	assertProgressOutput(t, 25000000000, 50000000000, "250 GB / 500 GB")
}

func assertProgressOutput(t *testing.T, doneBytes, totalBytes int64, expected string) {
	pb := New(totalBytes)
	buffer := bytes.NewBufferString("")
	pb.SetByteBuffer(buffer)
	pb.Set(doneBytes)
	assert.Equal(t, expected, buffer.String())
}
