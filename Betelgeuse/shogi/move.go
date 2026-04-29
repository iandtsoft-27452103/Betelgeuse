package shogi

func To(m uint32) uint32 {
	return m & 0x007f
}

func From(m uint32) uint32 {
	return (m >> 7) & 0x007f
}

func FlagPromo(m uint32) uint32 {
	return (m >> 14) & 1
}

func PieceType(m uint32) uint32 {
	return (m >> 15) & 0x000f
}

func CapPiece(m uint32) uint32 {
	return (m >> 19) & 0x000f
}

	/*
	xxxxxxxx xxxxxxxx x1111111 To
	xxxxxxxx xx111111 1xxxxxxx From
	xxxxxxxx x1xxxxxx xxxxxxxx FlagPromo
	xxxxx111 1xxxxxxx xxxxxxxx PieceType
	x1111xxx xxxxxxxx xxxxxxxx CapPiece
	*/

func Pack(from uint32, to uint32, pc uint32, cap_pc uint32, flag_promo uint32) uint32 {
	return (cap_pc << 19) | (pc << 15) | (flag_promo << 14) | (from << 7) | to
}

func SetNullMove() uint32{
	return (1 << 23)
}
