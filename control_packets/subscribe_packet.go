package control_packets

import (
	"bytes"
	"io"
)

type SubscribePacket struct {
	FixedHeader
	Message_id      uint16
	Topics          []string
	Qos_topics      []byte
}


func (subscribe_packet *SubscribePacket) Boxing(writer io.Writer) error {
	var body bytes.Buffer
	var err error

	body.Write(encode_uint16(subscribe_packet.Message_id))
	for index, topic := range subscribe_packet.Topics {
		body.Write(encode_string(topic))
		body.WriteByte(subscribe_packet.Qos_topics[index])
	}
	subscribe_packet.FixedHeader.Remaining_len = body.Len()
	packet := subscribe_packet.FixedHeader.Boxing_fh()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(writer)
	return err
}

func (subscribe_packet *SubscribePacket) Unboxing(bytes_reader io.Reader) error {
	var err error
	subscribe_packet.Message_id, err = decode_uint16(bytes_reader)
	if err != nil {
		return err
	}
	payload_len := subscribe_packet.FixedHeader.Remaining_len - 2
	for payload_len > 0 {
		topic, err := decode_string(bytes_reader)
		if err != nil {
			return err
		}
		subscribe_packet.Topics = append(subscribe_packet.Topics, topic)
		qos, err := decode_byte(bytes_reader)
		if err != nil {
			return err
		}
		subscribe_packet.Qos_topics = append(subscribe_packet.Qos_topics, qos)
		payload_len -= 2 + len(topic) + 1
	}
	return nil
}