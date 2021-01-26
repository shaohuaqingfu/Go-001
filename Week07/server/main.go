package main

import (
	"Week07/constants"
	"Week07/errors"
	"Week07/model"
	"Week07/sync"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":11322")
	errors.PrintWrapError(err, "server: tcp监听失败")
	defer func() {
		err = listener.Close()
		errors.PrintWrapError(err, "server: tcp连接关闭失败")
	}()
	for {
		conn, err := listener.Accept()
		errors.PrintWrapError(err, "server: accept消息失败")
		bufChan := make(chan *bytes.Buffer)
		sync.Go(context.Background(), func(ctx context.Context) error {
			return readFromConn(conn, bufChan)
		})
		sync.Go(context.Background(), func(ctx context.Context) error {
			return handle(ctx, bufChan, conn)
		})
	}
}

func handle(ctx context.Context, bufChan <-chan *bytes.Buffer, conn net.Conn) error {
	for {
		select {
		case buf := <-bufChan:
			res := buf.Bytes()
			if len(res) == 0 {
				continue
			}
			msg := model.Message{}
			err := json.Unmarshal(res, &msg)
			errors.PrintWrapError(err, "server: 接收消息json序列化失败")
			fmt.Println("server: 接收消息 msg = ", msg)
			reply := []byte("server reply: ")
			_, err = conn.Write(append(reply, res...))
			errors.PrintWrapError(err, "server: 消息回复，写数据失败")
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
