package control_packets

import "io"

type PubackPacket struct {
	FixedHeader
	Message_id uint16
}


func (puback_packet *PubackPacket) Boxing(writer io.Writer) error {
	var err error
	puback_packet.FixedHeader.Remaining_len = 2
	packet := puback_packet.FixedHeader.Boxing_fh()
	packet.Write(encode_uint16(puback_packet.Message_id))
	_, err = packet.WriteTo(writer)
	return err
}

func (puback_packet *PubackPacket) Unboxing(bytes_reader io.Reader) error {
	var err error
	puback_packet.Message_id, err = decode_uint16(bytes_reader)
	return err
}