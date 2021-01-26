package main

import (
	"Week09/constants"
	"Week09/errors"
	"Week09/model"
	"Week09/sync"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {

	conn, err := net.Dial("tcp", ":11322")
	errors.PrintWrapError(err, "client: 连接server失败")
	m1 := &model.Message{
		ID:      1000,
		Type:    model.ReplyType,
		Kind:    model.ReplyKind,
		Title:   "message title",
		Content: "message content",
	}
	m2 := &model.Message{
		ID:      1001,
		Type:    model.ReplyType,
		Kind:    model.ReplyKind,
		Title:   "message title 2",
		Content: "message content 2",
	}
	buf, err := json.Marshal(m1)
	errors.PrintWrapError(err, "client: 发送消息json序列化失败")
	buf2, err := json.Marshal(m2)
	errors.PrintWrapError(err, "client: 发送消息json序列化失败")
	buf = append(buf, constants.NByte)
	buf = append(buf, buf2...)
	buf = append(buf, constants.NByte)
	_, err = conn.Write(buf)
	errors.PrintWrapError(err, "client: 发送消息失败")

	bufChan := make(chan *bytes.Buffer)
	sync.Go(context.Background(), func(ctx context.Context) error {
		return readFromConn(conn, bufChan)
	})
	brk := false
	for {
		select {
		case buffer := <-bufChan:
			fmt.Print("client: ", buffer.String())
		case <-time.After(5 * time.Second):
			brk = true
			break
		}
		if brk {
			break
		}
	}
}

func readFromConn(conn net.Conn, bufChan chan *bytes.Buffer) error {
	reader := bufio.NewReader(conn)
	buf := make([]byte, 2048)

	for {
		localBuf := new(bytes.Buffer)
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				bufChan <- localBuf
				return nil
			}
			return err
		}

		var last = 0
		for i := 0; i < n; i++ {
			if buf[i] == constants.NByte {
				_, err = localBuf.Write(buf[last : i+1])
				errors.PrintWrapError(err, "server: 写缓冲失败")
				bufChan <- localBuf
				localBuf = new(bytes.Buffer)
				last = i + 1
			}
		}
		if last != n {
			_, err = localBuf.Write(buf[last:])
			errors.PrintWrapError(err, "server: 写缓冲失败")
		}
	}
}
