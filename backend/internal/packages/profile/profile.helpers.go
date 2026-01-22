package profile

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strings"
	"time"
)

func scanForVirus(r io.Reader) error {
	const (
		clamdAddr = "localhost:3310"
		timeout   = 30 * time.Second
	)

	conn, err := net.DialTimeout("tcp", clamdAddr, timeout)
	if err != nil {
		return errors.New("virus scanner unavailable")
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(timeout))

	writer := bufio.NewWriter(conn)

	// INSTREAM command (z = null-terminated)
	if _, err := writer.WriteString("zINSTREAM\000"); err != nil {
		return err
	}

	buf := make([]byte, 32*1024) // 32KB chunks
	for {
		n, err := r.Read(buf)
		if n > 0 {
			// write chunk length
			if err := binary.Write(writer, binary.BigEndian, uint32(n)); err != nil {
				return err
			}

			// write chunk data
			if _, err := writer.Write(buf[:n]); err != nil {
				return err
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	// zero-length chunk = end of stream
	if err := binary.Write(writer, binary.BigEndian, uint32(0)); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	// read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// clean response
	response = strings.TrimSpace(response)

	switch {
	case strings.HasSuffix(response, "OK"):
		return nil

	case strings.Contains(response, "FOUND"):
		return errors.New("virus detected")

	default:
		return errors.New("virus scan error")
	}
}
