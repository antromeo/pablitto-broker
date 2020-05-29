package control_packets

import "io"

type DisconnectPacket struct {
	FixedHeader
}

func (disconnect_packet *DisconnectPacket) Boxing(writer io.Writer) error {
	packet := disconnect_packet.FixedHeader.Boxing_fh()
	_, err := packet.WriteTo(writer)
	return err
}
func (disconnect_packet *DisconnectPacket) Unboxing(bytes_reader io.Reader) error {
	return nil
}
