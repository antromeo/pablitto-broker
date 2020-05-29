package control_packets

import (
	"bytes"
	"io"
)

type SubackPacket struct {
	FixedHeader
	Message_id   uint16
	Return_codes_topic []byte
}

func (suback_packet *SubackPacket) Boxing(writer io.Writer) error {
	var body bytes.Buffer
	var err error
	body.Write(encode_uint16(suback_packet.Message_id))
	body.Write(suback_packet.Return_codes_topic)
	suback_packet.FixedHeader.Remaining_len = body.Len()
	packet := suback_packet.FixedHeader.Boxing_fh()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(writer)
	return err
}

func (suback_packet *SubackPacket) Unboxing(bytes_reader io.Reader) error {
	var qos_topics bytes.Buffer
	var err error
	suback_packet.Message_id, err = decode_uint16(bytes_reader)
	if err != nil {
		return err
	}
	_, err = qos_topics.ReadFrom(bytes_reader)
	if err != nil {
		return err
	}
	suback_packet.Return_codes_topic = qos_topics.Bytes()
	return nil
}