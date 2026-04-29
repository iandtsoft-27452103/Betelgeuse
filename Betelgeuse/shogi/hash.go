package shogi

//"Betelgeuse/bitop"
//"Betelgeuse/common"
//"fmt"

var PieceRand [Color_NB][Piece_NB][Square_NB]uint64

type RandWorkT struct {
	Count int
	Cnst  [2]uint
	Vec   [RandN]uint
}

var RandWork RandWorkT

func IniRand(u uint) {
	RandWork.Count = RandN
	RandWork.Cnst[0] = 0
	RandWork.Cnst[1] = 0x9908b0df
	var i uint
	for i = 1; i < RandN; i++ {
		u = uint(i + 1812433253*(u^(u>>30)))
		u &= Mask32
		RandWork.Vec[i] = u
	}
	//fmt.Println("Mr. McMahon!")
}

func Rand32() uint {
	var u uint
	var u0 uint
	var u1 uint
	var u2 uint
	var i int
	if RandWork.Count == RandN {
		RandWork.Count = 0
		for i = 0; i < RandN-RandM; i++ {
			u = RandWork.Vec[i] & MaskU
			u |= RandWork.Vec[i+1] & MaskL
			u0 = RandWork.Vec[i+RandM]
			u1 = u >> 1
			u2 = RandWork.Cnst[u&1]
			RandWork.Vec[i] = u0 ^ u1 ^ u2
		}
		for ; i < RandN-1; i++ {
			u = RandWork.Vec[i] & MaskU
			u |= RandWork.Vec[i+1] & MaskL
			u0 = RandWork.Vec[i+RandM-RandN]
			u1 = u >> 1
			u2 = RandWork.Cnst[u&1]
			RandWork.Vec[i] = u0 ^ u1 ^ u2
		}
		u = RandWork.Vec[RandN-1] & MaskU
		u |= RandWork.Vec[0] & MaskL
		u0 = RandWork.Vec[RandM-1]
		u1 = u >> 1
		u2 = RandWork.Cnst[u&1]
		RandWork.Vec[RandN-1] = u0 ^ u1 ^ u2
	}
	u = RandWork.Vec[RandWork.Count]
	RandWork.Count++
	u ^= (u >> 11)
	u ^= (u << 7) & 0x9d2c5680
	u ^= (u << 15) & 0xefc60000
	u ^= (u >> 18)
	return u
}

func Rand64() uint {
	var h = Rand32()
	var l = Rand32()
	return l | (h << 32)
}

func IniRandomTable() {
	for c := 0; c < Color_NB; c++ {
		for pc := Pawn; pc < Dragon; pc++ {
			for sq := 0; sq < Square_NB; sq++ {
				PieceRand[c][pc][sq] = uint64(Rand64())
			}
		}
	}
}

func HashFunc(bt BoardTree) uint64 {
	var key uint64 = 0
	for c := 0; c < Color_NB; c++ {
		for pc := Pawn; pc <= Dragon; pc++ {
			var bb = bt.BB_Piece[c][pc]
			if pc == None {
				continue
			}
			for BBTest(bb) > 0 {
				//var sq = uint(80) - bb.TrailingZeroBits()
				var sq = Square(bb)
				Xor(&bb, sq)
				//fmt.Println(bb)
				//fmt.Println(Atk.ABB_Mask[sq])
				key ^= PieceRand[c][pc][sq]
			}
		}
	}
	return key
}
