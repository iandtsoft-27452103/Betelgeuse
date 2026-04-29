package shogi

//"Betelgeuse/common"
//"math/bits"

func InitBoard() BoardTree {
	var bt BoardTree
	bt.BB_Occupied[0].P[2] = 134022655
	bt.BB_Occupied[1].P[0] = 134022655
	//bt.BB_Rotated[0] = BBIni()
	//bt.BB_Rotated[0] = BBOr(BB_File[0], BB_File[2])
	//bt.BB_Rotated[0] = BBOr(bt.BB_Rotated[0], Atk.ABB_Mask[10])
	//bt.BB_Rotated[0] = BBOr(bt.BB_Rotated[0], Atk.ABB_Mask[64])
	//bt.BB_Rotated[1] = BBIni()
	//bt.BB_Rotated[1] = BBOr(BB_File[6], BB_File[8])
	//bt.BB_Rotated[1] = BBOr(bt.BB_Rotated[1], Atk.ABB_Mask[16])
	//bt.BB_Rotated[1] = BBOr(bt.BB_Rotated[1], Atk.ABB_Mask[70])
	bt.BB_Piece[0][Pawn].P[2] = 133955584
	bt.BB_Piece[1][Pawn].P[0] = 511
	bt.BB_Piece[0][Lance].P[2] = 257
	bt.BB_Piece[1][Lance].P[0] = 67371008
	bt.BB_Piece[0][Knight].P[2] = 130
	bt.BB_Piece[1][Knight].P[0] = 34078720
	bt.BB_Piece[0][Silver].P[2] = 68
	bt.BB_Piece[1][Silver].P[0] = 17825792
	bt.BB_Piece[0][Gold].P[2] = 40
	bt.BB_Piece[1][Gold].P[0] = 10485760
	bt.BB_Piece[0][Bishop].P[2] = 65536
	bt.BB_Piece[1][Bishop].P[0] = 1024
	bt.BB_Piece[0][Rook].P[2] = 1024
	bt.BB_Piece[1][Rook].P[0] = 65536
	bt.BB_Piece[0][King].P[2] = 16
	bt.BB_Piece[1][King].P[0] = 4194304
	for c := Black; c <= White; c++ {
		for pc := Pro_Pawn; pc <= Dragon; pc++ {
			bt.BB_Piece[c][pc] = BBIni()
		}
	}
	bt.Board[0] = int8(-Lance)
	bt.Board[1] = int8(-Knight)
	bt.Board[2] = int8(-Silver)
	bt.Board[3] = int8(-Gold)
	bt.Board[4] = int8(-King)
	bt.Board[5] = int8(-Gold)
	bt.Board[6] = int8(-Silver)
	bt.Board[7] = int8(-Knight)
	bt.Board[8] = int8(-Lance)
	bt.Board[10] = int8(-Rook)
	bt.Board[16] = int8(-Bishop)
	for sq := 18; sq <= 26; sq++ {
		bt.Board[sq] = int8(-Pawn)
	}
	for sq := 54; sq <= 62; sq++ {
		bt.Board[sq] = int8(Pawn)
	}
	bt.Board[64] = int8(Bishop)
	bt.Board[70] = int8(Rook)
	bt.Board[72] = int8(Lance)
	bt.Board[73] = int8(Knight)
	bt.Board[74] = int8(Silver)
	bt.Board[75] = int8(Gold)
	bt.Board[76] = int8(King)
	bt.Board[77] = int8(Gold)
	bt.Board[78] = int8(Silver)
	bt.Board[79] = int8(Knight)
	bt.Board[80] = int8(Lance)
	bt.Hand[0] = uint32(0)
	bt.Hand[1] = uint32(0)
	bt.CurrentHash = HashFunc(bt)
	bt.RootColor = uint8(Black)
	bt.SQ_King[Black] = uint8(76)
	bt.SQ_King[White] = uint8(4)
	bt.Ply = 1
	bt.PrevHash = 0
	bt.Hash[1] = bt.CurrentHash
	// EvalArray is not initialized.
	return bt
}

func Clear(bt *BoardTree) {
	for c := Black; c <= White; c++ {
		for i := 0; i <= 2; i++ {
			bt.BB_Occupied[c].P[i] = 0
			//bt.BB_Rotated[c].P[i] = 0
			bt.Hand[c] = 0
			bt.SQ_King[c] = 0
			for pc := Pawn; pc <= Dragon; pc++ {
				bt.BB_Piece[c][pc].P[i] = 0
			}
		}
	}
	for sq := 0; sq < Square_NB; sq++ {
		bt.Board[sq] = EMPTY
	}
	bt.CurrentHash = HashFunc(*bt)
	bt.RootColor = uint8(Black)
	bt.Ply = 1
	bt.PrevHash = 0
	for i := 0; i < Ply_Max; i++ {
		bt.EvalArray[i] = 0
	}
}

func DeepCopy(bt_base BoardTree) BoardTree {
	var bt BoardTree
	for c := Black; c < Color_NB; c++ {
		bt.SQ_King[c] = bt_base.SQ_King[c]
		bt.Hand[c] = bt_base.Hand[c]
		for i := 0; i < 3; i++ {
			bt.BB_Occupied[c].P[i] = bt_base.BB_Occupied[c].P[i]
			//bt.BB_Rotated[c].P[i] = bt_base.BB_Rotated[c].P[i]
		}
	}
	for c := Black; c < Color_NB; c++ {
		for pc := Pawn; pc <= Dragon; pc++ {
			for i := 0; i < 3; i++ {
				bt.BB_Piece[c][pc].P[i] = bt_base.BB_Piece[c][pc].P[i]
			}
		}
	}
	for sq := 0; sq < Square_NB; sq++ {
		bt.Board[sq] = bt_base.Board[sq]
	}
	for ply := 0; ply < Ply_Max; ply++ {
		bt.Hash[ply] = bt_base.Hash[ply]
		bt.EvalArray[ply] = bt_base.EvalArray[ply]
	}
	bt.RootColor = bt_base.RootColor
	bt.CurrentHash = bt_base.CurrentHash
	bt.PrevHash = bt_base.PrevHash
	bt.Ply = bt_base.Ply
	return bt
}

func Do(bt *BoardTree, m uint32, color int) {
	bt.PrevHash = bt.CurrentHash
	var ifrom = From(m)
	var ito = To(m)
	var ipiece = PieceType(m)
	var is_promote = FlagPromo(m)
	if ifrom >= Square_NB {
		Xor(&bt.BB_Piece[color][ipiece], int(ito))
		bt.CurrentHash ^= PieceRand[color][ipiece][ito]
		bt.Hand[color] -= Hand_Hash[ipiece]
		bt.Board[ito] = int8(-Sign_Table[color]) * int8(ipiece)
		Xor(&bt.BB_Occupied[color], int(ito))
		//Xor(&bt.BB_Rotated[color], int(Index_Rotated[ito]))
	} else {
		var bb_set_clear = BBOr(Atk.ABB_Mask[ifrom], Atk.ABB_Mask[ito])
		//var bb_set_clear2 = BBOr(Atk.ABB_Mask[Index_Rotated[ifrom]], Atk.ABB_Mask[Index_Rotated[ito]])
		bt.BB_Occupied[color] = BBXor(bt.BB_Occupied[color], bb_set_clear)
		//bt.BB_Rotated[color] = BBXor(bt.BB_Rotated[color], bb_set_clear2)
		bt.Board[ifrom] = int8(Empty)
		if is_promote > 0 {
			Xor(&bt.BB_Piece[color][ipiece], int(ifrom))
			Xor(&bt.BB_Piece[color][ipiece+Promote], int(ito))
			bt.CurrentHash ^= PieceRand[color][ipiece][ifrom] ^ PieceRand[color][ipiece+Promote][ito]
			bt.Board[ito] = int8(-Sign_Table[color]) * int8(ipiece+Promote)
		} else {
			if ipiece == uint32(King) {
				bt.SQ_King[color] = uint8(ito)
			}
			bt.BB_Piece[color][ipiece] = BBXor(bt.BB_Piece[color][ipiece], bb_set_clear)
			bt.CurrentHash ^= PieceRand[color][ipiece][ifrom] ^ PieceRand[color][ipiece][ito]
			bt.Board[ito] = int8(-Sign_Table[color]) * int8(ipiece)
		}
		var icap_piece = CapPiece(m)
		var index = icap_piece
		if icap_piece > 0 {
			if icap_piece > uint32(King) {
				index -= Promote
			}
			bt.Hand[color] += Hand_Hash[index]
			Xor(&bt.BB_Piece[color^1][icap_piece], int(ito))
			bt.CurrentHash ^= PieceRand[color^1][icap_piece][ito]
			Xor(&bt.BB_Occupied[color^1], int(ito))
			//Xor(&bt.BB_Rotated[color^1], int(Index_Rotated[ito]))
		}
	}
	bt.Hash[bt.Ply] = bt.PrevHash
	bt.Hash[bt.Ply+1] = bt.CurrentHash
	bt.Ply += 1
}
func UnDo(bt *BoardTree, m uint32, color int) {
	bt.CurrentHash = bt.PrevHash
	var ifrom = From(m)
	var ito = To(m)
	var ipiece = PieceType(m)
	var is_promote = FlagPromo(m)
	if ifrom >= Square_NB {
		Xor(&bt.BB_Piece[color][ipiece], int(ito))
		bt.Hand[color] += Hand_Hash[ipiece]
		bt.Board[ito] = int8(Empty)
		Xor(&bt.BB_Occupied[color], int(ito))
		//Xor(&bt.BB_Rotated[color], int(Index_Rotated[ito]))
	} else {
		var bb_set_clear = BBOr(Atk.ABB_Mask[ifrom], Atk.ABB_Mask[ito])
		//var bb_set_clear2 = BBOr(Atk.ABB_Mask[Index_Rotated[ifrom]], Atk.ABB_Mask[Index_Rotated[ito]])
		bt.BB_Occupied[color] = BBXor(bt.BB_Occupied[color], bb_set_clear)
		//bt.BB_Rotated[color] = BBXor(bt.BB_Rotated[color], bb_set_clear2)
		bt.Board[ifrom] = int8(-Sign_Table[color]) * int8(ipiece)
		if is_promote > 0 {
			Xor(&bt.BB_Piece[color][ipiece], int(ifrom))
			Xor(&bt.BB_Piece[color][ipiece+Promote], int(ito))
		} else {
			if ipiece == uint32(King) {
				bt.SQ_King[color] = uint8(ifrom)
			}
			bt.BB_Piece[color][ipiece] = BBXor(bt.BB_Piece[color][ipiece], bb_set_clear)
		}
		var icap_piece = CapPiece(m)
		var index = icap_piece
		if icap_piece > 0 {
			if icap_piece > uint32(King) {
				index -= Promote
			}
			bt.Hand[color] -= Hand_Hash[index]
			Xor(&bt.BB_Piece[color^1][icap_piece], int(ito))
			Xor(&bt.BB_Occupied[color^1], int(ito))
			//Xor(&bt.BB_Rotated[color^1], int(Index_Rotated[ito]))
			bt.Board[ito] = int8(Sign_Table[color]) * int8(icap_piece)
		} else {
			bt.Board[ito] = int8(Empty)
		}
	}
	bt.PrevHash = bt.Hash[bt.Ply-2]
	bt.Hash[bt.Ply] = 0
	bt.Ply -= 1
}

// Explanation for return value.
// 0 : You can not declare winning in this position.
// 1 : The winner is black player.
// 2 : The winner is white player.
func IsDeclarationWin(bt BoardTree) uint8 {
	var black_score uint8 = 0
	var white_score uint8 = 0
	var b_tekijin_piece_count uint8 = 0
	var w_tekijin_piece_count uint8 = 0
	var b_hand_piece_count [Rook + 1]uint8 = [Rook + 1]uint8{0, 0, 0, 0, 0, 0, 0, 0}
	var w_hand_piece_count [Rook + 1]uint8 = [Rook + 1]uint8{0, 0, 0, 0, 0, 0, 0, 0}
	var b_board_piece_count [Piece_NB]uint8 = [Piece_NB]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var w_board_piece_count [Piece_NB]uint8 = [Piece_NB]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var bb0 = BBAnd(bt.BB_Piece[Black][King], BB_White_Position)
	var bb1 = BBAnd(bt.BB_Piece[White][King], BB_Black_Position)
	if BBTest(bb0) == uint32(0) && BBTest(bb1) == uint32(0) {
		return 0
	}
	if BBTest(bb0) > 0 {
		for i := Pawn; i <= Rook; i++ {
			b_hand_piece_count[i] = uint8((bt.Hand[Black] & Hand_Mask[i]) >> Hand_Rev_Bit[i])
			if i >= Bishop {
				black_score += uint8(5) * b_hand_piece_count[i]
			} else {
				black_score += b_hand_piece_count[i]
			}
		}
		for i := Pawn; i <= Dragon; i++ {
			if i == None {
				continue
			}
			var bb_object = BBAnd(bt.BB_Piece[Black][i], BB_Rev_Color_Position[Black])
			b_board_piece_count[i] = uint8(PopCount(bb_object))
			b_tekijin_piece_count += b_board_piece_count[i]
			var bb_temp = BBOr(BB_DMZ, BB_Rev_Color_Position[White])
			bb_object = BBAnd(bb_temp, bt.BB_Piece[Black][i])
			b_board_piece_count[i] += uint8(PopCount(bb_object))
			if i == King {
				continue
			}
			if i == Bishop || i == Rook || i >= Horse {
				black_score += uint8(5) * b_board_piece_count[i]
			} else {
				black_score += b_board_piece_count[i]
			}
		}
	}
	if BBTest(bb1) > 0 {
		for i := Pawn; i <= Rook; i++ {
			w_hand_piece_count[i] = uint8((bt.Hand[White] & Hand_Mask[i]) >> Hand_Rev_Bit[i])
			if i >= Bishop {
				white_score += uint8(5) * w_hand_piece_count[i]
			} else {
				white_score += w_hand_piece_count[i]
			}
		}
		for i := Pawn; i <= Dragon; i++ {
			if i == None {
				continue
			}
			var bb_object = BBAnd(bt.BB_Piece[White][i], BB_Rev_Color_Position[White])
			w_board_piece_count[i] = uint8(PopCount(bb_object))
			w_tekijin_piece_count += w_board_piece_count[i]
			var bb_temp = BBOr(BB_DMZ, BB_Rev_Color_Position[Black])
			bb_object = BBAnd(bb_temp, bt.BB_Piece[White][i])
			w_board_piece_count[i] += uint8(PopCount(bb_object))
			if i == King {
				continue
			}
			if i == Bishop || i == Rook || i >= Horse {
				white_score += uint8(5) * w_board_piece_count[i]
			} else {
				white_score += w_board_piece_count[i]
			}
		}
	}
	if BBTest(bb0) > 0 && black_score >= uint8(28) && b_tekijin_piece_count >= uint8(10) {
		return 1
	}
	if BBTest(bb1) > 0 && white_score >= uint8(27) && w_tekijin_piece_count >= uint8(10) {
		return 2
	}
	return 0
}
func IsRepetition(bt BoardTree, tt TT) uint8 {
	var limit uint16 = bt.Ply - 12
	if limit < 1 {
		return 0
	}
	var counter = 0
	var i = bt.Ply
	for i >= limit {
		if bt.CurrentHash == bt.Hash[i] {
			counter += 1
		}
		i -= 1
	}
	// If the same hash value is detected three times, the repetition is realizing.
	if counter > 2 {
		if _, ok := tt.is_check[bt.CurrentHash]; ok {
			var b = tt.is_check[bt.CurrentHash]
			if !b {
				return 1 // normal repetition
			} else {
				return 2 // check succeeding repetition -> the turn of doing check move is a loser.
			}
		}
	}
	return 0 // The repetition is not detected.
}
