package sshex

// https://datatracker.ietf.org/doc/html/rfc4254
//byte      SSH_MSG_CHANNEL_REQUEST
//uint32    recipient channel
//string    "pty-req"
//boolean   want_reply
//string    TERM environment variable value (e.g., vt100)
//uint32    terminal width, characters (e.g., 80)
//uint32    terminal height, rows (e.g., 24)
//uint32    terminal width, pixels (e.g., 640)
//uint32    terminal height, pixels (e.g., 480)
//string    encoded terminal modes

type PtyReqPayload struct {
	Env         string
	CharWidth   uint32
	CharHeight  uint32
	PixelWidth  uint32
	PixelHeight uint32
	Mode        []byte
}

func ParsePtyReqPayload(buf []byte) (*PtyReqPayload, error) {
	var err error
	pty := &PtyReqPayload{}
	pl := NewPayload(buf)

	pty.Env, err = pl.ReadString()
	if nil != err {
		return nil, err
	}

	pty.CharWidth, err = pl.ReadUint32()
	if nil != err {
		return nil, err
	}

	pty.CharHeight, err = pl.ReadUint32()
	if nil != err {
		return nil, err
	}

	pty.PixelWidth, err = pl.ReadUint32()
	if nil != err {
		return nil, err
	}

	pty.PixelHeight, err = pl.ReadUint32()
	if nil != err {
		return nil, err
	}

	pty.Mode, err = pl.ReadBytes()
	if nil != err {
		return nil, err
	}

	return pty, nil
}
