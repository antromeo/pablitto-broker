package control_packets

import (
	"bytes"
	"fmt"
	"io"
)

type PublishPacket struct {
	FixedHeader
	Topic_name string
	Message_id uint16
	Payload    []byte
}

func (publish_packet *PublishPacket) Boxing(writer io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(encode_string(publish_packet.Topic_name))
	if publish_packet.Qos > 0 {
		body.Write(encode_uint16(publish_packet.Message_id))
	}
	publish_packet.FixedHeader.Remaining_len = body.Len()+len(publish_packet.Payload)
	packet := publish_packet.FixedHeader.Boxing_fh()
	packet.Write(body.Bytes())
	packet.Write(publish_packet.Payload)
	_, err = writer.Write(packet.Bytes())
	return err
}


func (publish_packet *PublishPacket) Unboxing(bytes_reader io.Reader) error {
	var payload_len = publish_packet.FixedHeader.Remaining_len
	var err error
	publish_packet.Topic_name, err = decode_string(bytes_reader)
	if err != nil {
		return err
	}

	if publish_packet.Qos > 0 {
		publish_packet.Message_id, err = decode_uint16(bytes_reader)
		if err != nil {
			return err
		}
		payload_len -= len(publish_packet.Topic_name) + 4
	} else {
		payload_len -= len(publish_packet.Topic_name) + 2}
	if payload_len < 0 {
		return fmt.Errorf("error unpacking publish, payload length < 0")
	}
	publish_packet.Payload = make([]byte, payload_len)
	_, err = bytes_reader.Read(publish_packet.Payload)
	return err
}
