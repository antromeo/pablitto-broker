package control_packets

import "io"

type PingrespPacket struct {
	FixedHeader
}

func (pingresp_packet *PingrespPacket) Boxing(writer io.Writer) error {
	packet := pingresp_packet.FixedHeader.Boxing_fh()
	_, err := packet.WriteTo(writer)
	return err
}

func (pingresp_packet *PingrespPacket) Unboxing(bytes_reader io.Reader) error {
	return nil
}
