package control_packets

import (
	"bytes"
	"io"
)

type UnsubscribePacket struct {
	FixedHeader
	Message_id uint16
	Topics    []string
}



func (unsubscribe_packet *UnsubscribePacket) Boxing(writer io.Writer) error {
	var body bytes.Buffer
	var err error
	body.Write(encode_uint16(unsubscribe_packet.Message_id))
	for _, topic := range unsubscribe_packet.Topics {
		body.Write(encode_string(topic))
	}
	unsubscribe_packet.FixedHeader.Remaining_len = body.Len()
	packet := unsubscribe_packet.FixedHeader.Boxing_fh()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(writer)
	return err
}


func (unsubscribe_packet *UnsubscribePacket) Unboxing(bytes_reader io.Reader) error {
	var err error
	unsubscribe_packet.Message_id, err = decode_uint16(bytes_reader)
	if err != nil {
		return err
	}
	for topic, err := decode_string(bytes_reader); err == nil && topic != ""; topic, err = decode_string(bytes_reader) {
		unsubscribe_packet.Topics = append(unsubscribe_packet.Topics, topic)
	}
	return err
}