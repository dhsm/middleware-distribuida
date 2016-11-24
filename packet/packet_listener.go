package packet

type PacketListener interface {
    OnPacket (pkt Packet)
}
