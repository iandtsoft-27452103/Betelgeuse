package shogi

func CSA2Move(bt BoardTree, str_csa string) uint32 {
	var ifrom = uint32(CSA_TO_SQ[str_csa[0:2]])
	var ito = uint32(CSA_TO_SQ[str_csa[2:4]])
	var flag_promo uint32 = 0
	var pc int
	if ifrom < Square_NB {
		pc = int(bt.Board[ifrom])
		if pc < 0 {
			pc = -pc
		}
	} else {
		pc = int(CSA_TO_PC[str_csa[4:6]])
		ifrom += uint32(pc) - 1
	}
	var cap_pc = int(bt.Board[ito])
	if cap_pc < 0 {
		cap_pc = -cap_pc
	}
	if pc < int(King) && CSA_TO_PC[str_csa[4:6]] > King {
		flag_promo = 1
	}
	return Pack(ifrom, ito, uint32(pc), uint32(cap_pc), flag_promo)
}
func Move2CSA(move uint32) string {
	var str string
	str = Str_CSA[From(move)] + Str_CSA[To(move)]
	if FlagPromo(move) == 0 {
		str += Str_Piece[PieceType(move)]
	} else {
		str += Str_Piece[PieceType(move)+8]
	}
	return str
}
