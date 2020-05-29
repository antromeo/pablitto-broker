package control_packets

import "io"

type UnsubackPacket struct {
	FixedHeader
	Message_id uint16
}


func (unsuback_packet *UnsubackPacket) Boxing(writer io.Writer) error {
	var err error
	unsuback_packet.FixedHeader.Remaining_len = 2
	packet := unsuback_packet.FixedHeader.Boxing_fh()
	packet.Write(encode_uint16(unsuback_packet.Message_id))
	_, err = packet.WriteTo(writer)
	return err
}

func (unsuback_packet *UnsubackPacket) Unboxing(bytes_reader io.Reader) error {
	var err error
	unsuback_packet.Message_id, err = decode_uint16(bytes_reader)
	return err
}
