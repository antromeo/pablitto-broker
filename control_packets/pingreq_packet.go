package control_packets

import "io"

type PingreqPacket struct {
	FixedHeader
}

func (pingreq_packet *PingreqPacket) Boxing(writer io.Writer) error {
	packet := pingreq_packet.FixedHeader.Boxing_fh()
	_, err := packet.WriteTo(writer)
	return err
}
func (pr *PingreqPacket) Unboxing(bytes_reader io.Reader) error {
	return nil
}
