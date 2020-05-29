package control_packets

import (
	"bytes"
	"io"
)

type ConnackPacket struct {
	FixedHeader
	Session_present bool
	Return_code     byte
}

func (connack_packet *ConnackPacket) Boxing(writer io.Writer) error {
	var body_bytes bytes.Buffer
	var err error

	body_bytes.WriteByte(bool_to_byte(connack_packet.Session_present))
	body_bytes.WriteByte(connack_packet.Return_code)
	connack_packet.FixedHeader.Remaining_len = 2
	control_packet := connack_packet.FixedHeader.Boxing_fh()
	control_packet.Write(body_bytes.Bytes())
	_, err = control_packet.WriteTo(writer)
	return err
}

func (connack_packet *ConnackPacket) Unboxing(bytes_reader io.Reader) error {
	first_byte, err := decode_byte(bytes_reader)
	if err != nil {
		return err
	}
	connack_packet.Session_present = 1&first_byte > 0
	connack_packet.Return_code, err = decode_byte(bytes_reader)
	return err
}