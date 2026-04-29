package shogi

import (
	//"Betelgeuse/common"
	"math/bits"
)

func PextAsm(src, mask uint64) uint64 // implemented in pext_amd64.s

func BBIni() BitBoard {
	var bb BitBoard
	bb.P[0] = 0
	bb.P[1] = 0
	bb.P[2] = 0
	return bb
}

func BBOr(bb0 BitBoard, bb1 BitBoard) BitBoard {
	var bb BitBoard
	bb.P[0] = bb0.P[0] | bb1.P[0]
	bb.P[1] = bb0.P[1] | bb1.P[1]
	bb.P[2] = bb0.P[2] | bb1.P[2]
	return bb
}

func BBAnd(bb0 BitBoard, bb1 BitBoard) BitBoard {
	var bb BitBoard
	bb.P[0] = bb0.P[0] & bb1.P[0]
	bb.P[1] = bb0.P[1] & bb1.P[1]
	bb.P[2] = bb0.P[2] & bb1.P[2]
	return bb
}

func BBAndOr(bb_base BitBoard, bb0 BitBoard, bb1 BitBoard) BitBoard {
	var bb BitBoard
	bb.P[0] = bb0.P[0] & bb1.P[0]
	bb.P[1] = bb0.P[1] & bb1.P[1]
	bb.P[2] = bb0.P[2] & bb1.P[2]
	return BBOr(bb_base, bb)
}

func BBNot(bb BitBoard) BitBoard {
	var bb_temp BitBoard
	bb_temp.P[0] = ^bb.P[0]
	bb_temp.P[1] = ^bb.P[1]
	bb_temp.P[2] = ^bb.P[2]
	return bb_temp
}

func BBXor(bb0 BitBoard, bb1 BitBoard) BitBoard {
	var bb BitBoard
	bb.P[0] = bb0.P[0] ^ bb1.P[0]
	bb.P[1] = bb0.P[1] ^ bb1.P[1]
	bb.P[2] = bb0.P[2] ^ bb1.P[2]
	return bb
}

func BBTest(bb BitBoard) uint32 {
	return bb.P[0] | bb.P[1] | bb.P[2]
}

func Xor(bb *BitBoard, sq int) {
	bb.P[0] ^= Atk.ABB_Mask[sq].P[0]
	bb.P[1] ^= Atk.ABB_Mask[sq].P[1]
	bb.P[2] ^= Atk.ABB_Mask[sq].P[2]
}

func Square(bb BitBoard) int {
	if bb.P[0] > 0 {
		return bits.LeadingZeros32(bb.P[0]) - 5
	}
	if bb.P[1] > 0 {
		return bits.LeadingZeros32(bb.P[1]) + 22
	}
	return bits.LeadingZeros32(bb.P[2]) + 49
}

func GetFileAttacks(bb_occupied BitBoard, sq int) BitBoard {
	var hash uint32
	var temp_hash [3]uint32
	temp_hash[0] = uint32(PextAsm(uint64(bb_occupied.P[0]), uint64(Atk.ABB_File_Mask_Ex[sq].P[0])))
	temp_hash[1] = uint32(PextAsm(uint64(bb_occupied.P[1]), uint64(Atk.ABB_File_Mask_Ex[sq].P[1])))
	temp_hash[2] = uint32(PextAsm(uint64(bb_occupied.P[2]), uint64(Atk.ABB_File_Mask_Ex[sq].P[2])))

	hash = ((temp_hash[0] << File_Shift_Table[sq][0]) | (temp_hash[1] << File_Shift_Table[sq][1]) | (temp_hash[2]))
	return Atk.ABB_File_Attacks[sq][hash]
}

func GetRankAttacks(bb_occupied BitBoard, sq int) BitBoard {
	var hash uint32
	switch RankTable[sq] {
	case Rank1, Rank2, Rank3:
		hash = uint32(PextAsm(uint64(bb_occupied.P[0]), uint64(Atk.ABB_Rank_Mask_Ex[sq].P[0])))
	case Rank4, Rank5, Rank6:
		hash = uint32(PextAsm(uint64(bb_occupied.P[1]), uint64(Atk.ABB_Rank_Mask_Ex[sq].P[1])))
	case Rank7, Rank8, Rank9:
		hash = uint32(PextAsm(uint64(bb_occupied.P[2]), uint64(Atk.ABB_Rank_Mask_Ex[sq].P[2])))
	}
	return Atk.ABB_Rank_Attacks[sq][hash]
}

func GetLanceAttacks(bb_occupied BitBoard, sq int, color int) BitBoard {
	return BBAnd(GetFileAttacks(bb_occupied, sq), Atk.ABB_Lance_Mask_Ex[color][sq])
}

func GetRookAttacks(bb_occupied BitBoard, sq int) BitBoard {
	return BBOr(GetFileAttacks(bb_occupied, sq), GetRankAttacks(bb_occupied, sq))
}

func GetDragonAttacks(bb_occupied BitBoard, sq int) BitBoard {
	return BBOr(Atk.ABB_Piece_Attacks[0][King][sq], BBOr(GetFileAttacks(bb_occupied, sq), GetRankAttacks(bb_occupied, sq)))
}

func GetDiag1Attacks(bb_occupied BitBoard, sq int) BitBoard {
	var hash uint32
	var temp_hash [3]uint32
	temp_hash[0] = uint32(PextAsm(uint64(bb_occupied.P[0]), uint64(Atk.ABB_Diag1_Mask_Ex[sq].P[0])))
	temp_hash[1] = uint32(PextAsm(uint64(bb_occupied.P[1]), uint64(Atk.ABB_Diag1_Mask_Ex[sq].P[1])))
	temp_hash[2] = uint32(PextAsm(uint64(bb_occupied.P[2]), uint64(Atk.ABB_Diag1_Mask_Ex[sq].P[2])))

	hash = ((temp_hash[0] << Diag1_Shift_Table[sq][0]) | (temp_hash[1] << Diag1_Shift_Table[sq][1]) | (temp_hash[2]))
	return Atk.ABB_Diag1_Attacks[sq][hash]
}

func GetDiag2Attacks(bb_occupied BitBoard, sq int) BitBoard {
	var hash uint32
	var temp_hash [3]uint32
	temp_hash[0] = uint32(PextAsm(uint64(bb_occupied.P[0]), uint64(Atk.ABB_Diag2_Mask_Ex[sq].P[0])))
	temp_hash[1] = uint32(PextAsm(uint64(bb_occupied.P[1]), uint64(Atk.ABB_Diag2_Mask_Ex[sq].P[1])))
	temp_hash[2] = uint32(PextAsm(uint64(bb_occupied.P[2]), uint64(Atk.ABB_Diag2_Mask_Ex[sq].P[2])))

	hash = ((temp_hash[0] << Diag2_Shift_Table[sq][0]) | (temp_hash[1] << Diag2_Shift_Table[sq][1]) | (temp_hash[2]))
	return Atk.ABB_Diag2_Attacks[sq][hash]
}

func GetBishopAttacks(bb_occupied BitBoard, sq int) BitBoard {
	return BBOr(GetDiag1Attacks(bb_occupied, sq), GetDiag2Attacks(bb_occupied, sq))
}

func GetHorseAttacks(bb_occupied BitBoard, sq int) BitBoard {
	return BBOr(Atk.ABB_Piece_Attacks[0][King][sq], BBOr(GetDiag1Attacks(bb_occupied, sq), GetDiag2Attacks(bb_occupied, sq)))
}

/*func GetHash(bb BitBoard, index int8) uint32 {
	var hash uint32
	var b = bb.P[0] | bb.P[1] | bb.P[2]
	switch index {
	case 0:
		hash = ((b & BB_Rank[0].P[index]) >> 18) | ((b & BB_Rank[1].P[index]) >> 9) | (b & BB_Rank[2].P[index])
	case 1:
		hash = ((b & BB_Rank[3].P[index]) >> 18) | ((b & BB_Rank[4].P[index]) >> 9) | (b & BB_Rank[5].P[index])
	case 2:
		hash = ((b & BB_Rank[6].P[index]) >> 18) | ((b & BB_Rank[7].P[index]) >> 9) | (b & BB_Rank[8].P[index])
	}
	return (hash >> 1)
}*/

func PopCount(bb BitBoard) int {
	var n = bits.OnesCount32(bb.P[0])
	n += bits.OnesCount32(bb.P[1])
	n += bits.OnesCount32(bb.P[2])
	return n
}
