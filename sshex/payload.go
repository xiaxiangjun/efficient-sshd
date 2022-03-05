package sshex

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Payload struct {
	reader io.Reader
}

func NewPayload(buf []byte) Payload {
	return Payload{
		reader: bytes.NewReader(buf),
	}
}

func (self *Payload) ReadBytes() ([]byte, error) {
	var size [4]byte

	// 读取长度
	_, err := io.ReadFull(self.reader, size[:])
	if nil != err {
		return nil, err
	}

	buf := make([]byte, binary.BigEndian.Uint32(size[:]))
	_, err = io.ReadFull(self.reader, buf)
	if nil != err {
		return nil, err
	}

	return buf, nil
}

func (self *Payload) ReadString() (string, error) {
	buf, err := self.ReadBytes()
	if nil != err {
		return "", err
	}

	return string(buf), nil
}

func (self *Payload) ReadUint32() (uint32, error) {
	var size [4]byte

	// 读取长度
	_, err := io.ReadFull(self.reader, size[:])
	if nil != err {
		return 0, err
	}

	return binary.BigEndian.Uint32(size[:]), nil
}
