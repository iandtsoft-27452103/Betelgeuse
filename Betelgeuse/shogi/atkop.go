package shogi

func IsPinnedOnKing(bt *BoardTree, sq int, idirec int, color int) BitBoard {
	var bb_ret = BBIni()
	var bb_occupied = BBOr(bt.BB_Occupied[0], bt.BB_Occupied[1])
	var bb_attacks = BBIni()
	var bb_object = BBIni()
	var adirec = idirec
	if adirec < 0 {
		adirec = -idirec
	}
	switch adirec {
	case Direc_File_U2d:
		bb_attacks = GetFileAttacks(bb_occupied, sq)
		//bb_attacks = BBAnd(bb_attacks, Atk.ABB_Mask[bt.SQ_King[color]])
		if BBTest(BBAnd(bb_attacks, Atk.ABB_Mask[bt.SQ_King[color]])) > 0 {
			bb_object = BBOr(bt.BB_Piece[color^1][Rook], bt.BB_Piece[color^1][Dragon])
			bb_object = BBOr(bb_object, bt.BB_Piece[color^1][Lance])
			return BBAnd(bb_attacks, bb_object)
		}
	case Direc_Rank_L2r:
		bb_attacks = GetRankAttacks(bb_occupied, sq)
		//bb_attacks = BBAnd(bb_attacks, Atk.ABB_Mask[bt.SQ_King[color]])
		if BBTest(BBAnd(bb_attacks, Atk.ABB_Mask[bt.SQ_King[color]])) > 0 {
			bb_object = BBOr(bt.BB_Piece[color^1][Rook], bt.BB_Piece[color^1][Dragon])
			return BBAnd(bb_attacks, bb_object)
		}
	case Direc_Diag1_U2d:
		bb_attacks = GetDiag1Attacks(bb_occupied, sq)
		//bb_attacks = BBAnd(bb_attacks, Atk.ABB_Mask[bt.SQ_King[color]])
		if BBTest(BBAnd(bb_attacks, Atk.ABB_Mask[bt.SQ_King[color]])) > 0 {
			bb_object = BBOr(bt.BB_Piece[color^1][Bishop], bt.BB_Piece[color^1][Horse])
			return BBAnd(bb_attacks, bb_object)
		}
	case Direc_Diag2_U2d:
		bb_attacks = GetDiag2Attacks(bb_occupied, sq)
		//bb_attacks = BBAnd(bb_attacks, Atk.ABB_Mask[bt.SQ_King[color]])
		if BBTest(BBAnd(bb_attacks, Atk.ABB_Mask[bt.SQ_King[color]])) > 0 {
			bb_object = BBOr(bt.BB_Piece[color^1][Bishop], bt.BB_Piece[color^1][Horse])
			return BBAnd(bb_attacks, bb_object)
		}
	}
	return bb_ret
}

func AttacksToPiece(bt BoardTree, sq int, color int) BitBoard {
	var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
	var bb_ret = BBAnd(bt.BB_Piece[color][Pawn], Atk.ABB_Piece_Attacks[color^1][Pawn][sq])
	bb_ret = BBAndOr(bb_ret, bt.BB_Piece[color][Knight], Atk.ABB_Piece_Attacks[color^1][Knight][sq])
	bb_ret = BBAndOr(bb_ret, bt.BB_Piece[color][Silver], Atk.ABB_Piece_Attacks[color^1][Silver][sq])
	var bb_total_gold = BBOr(bt.BB_Piece[color][Gold], bt.BB_Piece[color][Pro_Pawn])
	bb_total_gold = BBOr(bb_total_gold, bt.BB_Piece[color][Pro_Lance])
	bb_total_gold = BBOr(bb_total_gold, bt.BB_Piece[color][Pro_Knight])
	bb_total_gold = BBOr(bb_total_gold, bt.BB_Piece[color][Pro_Silver])
	bb_ret = BBAndOr(bb_ret, bb_total_gold, Atk.ABB_Piece_Attacks[color^1][Gold][sq])
	var bb_hdk = BBOr(bt.BB_Piece[color][Horse], bt.BB_Piece[color][Dragon])
	bb_hdk = BBOr(bb_hdk, bt.BB_Piece[color][King])
	bb_ret = BBAndOr(bb_ret, bb_hdk, Atk.ABB_Piece_Attacks[color^1][King][sq])
	var bb_bh = BBOr(bt.BB_Piece[color][Bishop], bt.BB_Piece[color][Horse])
	bb_ret = BBAndOr(bb_ret, bb_bh, GetBishopAttacks(bb_occupied, sq))
	var bb_rd = BBOr(bt.BB_Piece[color][Rook], bt.BB_Piece[color][Dragon])
	bb_ret = BBAndOr(bb_ret, bb_rd, GetRookAttacks(bb_occupied, sq))
	bb_ret = BBAndOr(bb_ret, bt.BB_Piece[color][Lance], GetLanceAttacks(bb_occupied, sq, color^1))
	return bb_ret
}

func IsMatePawnDrop(bt BoardTree, sq_drop int, color int) bool {
	if color == int(White) {
		if (sq_drop-9) >= 0 && bt.Board[sq_drop-9] != -int8(King) {
			return false
		}
	} else {
		if (sq_drop+9) < Square_NB && bt.Board[sq_drop+9] != int8(King) {
			return false
		}
	}
	var bb_sum = BBAnd(bt.BB_Piece[color][Knight], Atk.ABB_Piece_Attacks[color^1][Knight][sq_drop])
	bb_sum = BBAndOr(bb_sum, bt.BB_Piece[color][Silver], Atk.ABB_Piece_Attacks[color^1][Silver][sq_drop])
	var bb_total_gold = BBOr(bt.BB_Piece[color][Gold], bt.BB_Piece[color][Pro_Pawn])
	bb_total_gold = BBOr(bb_total_gold, bt.BB_Piece[color][Pro_Lance])
	bb_total_gold = BBOr(bb_total_gold, bt.BB_Piece[color][Pro_Knight])
	bb_total_gold = BBOr(bb_total_gold, bt.BB_Piece[color][Pro_Silver])
	bb_sum = BBAndOr(bb_sum, bb_total_gold, Atk.ABB_Piece_Attacks[color^1][Gold][sq_drop])
	var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
	var bb_bh = BBOr(bt.BB_Piece[color][Bishop], bt.BB_Piece[color][Horse])
	bb_sum = BBAndOr(bb_sum, bb_bh, GetBishopAttacks(bb_occupied, sq_drop))
	var bb_rd = BBOr(bt.BB_Piece[color][Rook], bt.BB_Piece[color][Dragon])
	bb_sum = BBAndOr(bb_sum, bb_rd, GetRookAttacks(bb_occupied, sq_drop))
	var bb_hd = BBOr(bt.BB_Piece[color][Horse], bt.BB_Piece[color][Dragon])
	bb_sum = BBAndOr(bb_sum, bb_hd, Atk.ABB_Piece_Attacks[color][King][sq_drop]) // Black king attacks are as same as whites.
	for BBTest(bb_sum) > 0 {
		var ifrom = Square(bb_sum)
		Xor(&bb_sum, ifrom)
		if IsDiscoverKing(&bt, ifrom, sq_drop, color) {
			continue
		}
		return false
	}
	var iking = bt.SQ_King[color]
	var bret = true
	Xor(&bt.BB_Occupied[color^1], sq_drop)
	var bb_temp2 = BBNot(bt.BB_Occupied[color])
	bb_temp2 = BBAnd(bb_temp2, BB_Full)
	var bb_move = BBAnd(Atk.ABB_Piece_Attacks[color][King][iking], bb_temp2)
	for BBTest(bb_move) > 0 {
		var ito = Square(bb_move)
		if BBTest(IsAttacked(bt, ito, color)) == 0 {
			bret = false
			break
		}
		Xor(&bb_move, ito)
	}
	Xor(&bt.BB_Occupied[color^1], sq_drop)
	return bret
}

func IsAttacked(bt BoardTree, sq int, color int) BitBoard {
	var bb_ret = BBIni()
	var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
	if sq+Delta_Table[color] >= 0 && (sq+Delta_Table[color]) < Square_NB {
		if bt.Board[sq+Delta_Table[color]] == int8(Sign_Table[color]*int(Pawn)) {
			bb_ret = Atk.ABB_Mask[sq+Delta_Table[color]]
		}
	}
	bb_ret = BBAndOr(bb_ret, bt.BB_Piece[color^1][Knight], Atk.ABB_Piece_Attacks[color][Knight][sq])
	bb_ret = BBAndOr(bb_ret, bt.BB_Piece[color^1][Silver], Atk.ABB_Piece_Attacks[color][Silver][sq])
	var bb_total_gold = BBOr(bt.BB_Piece[color^1][Gold], bt.BB_Piece[color^1][Pro_Pawn])
	bb_total_gold = BBOr(bb_total_gold, bt.BB_Piece[color^1][Pro_Lance])
	bb_total_gold = BBOr(bb_total_gold, bt.BB_Piece[color^1][Pro_Knight])
	bb_total_gold = BBOr(bb_total_gold, bt.BB_Piece[color^1][Pro_Silver])
	bb_ret = BBAndOr(bb_ret, bb_total_gold, Atk.ABB_Piece_Attacks[color][Gold][sq])
	var bb_hdk = BBOr(bt.BB_Piece[color^1][Horse], bt.BB_Piece[color^1][Dragon])
	bb_hdk = BBOr(bb_hdk, bt.BB_Piece[color^1][King])
	bb_ret = BBAndOr(bb_ret, bb_hdk, Atk.ABB_Piece_Attacks[color][King][sq])
	var bb_bh = BBOr(bt.BB_Piece[color^1][Bishop], bt.BB_Piece[color^1][Horse])
	bb_ret = BBAndOr(bb_ret, bb_bh, GetBishopAttacks(bb_occupied, sq))
	var bb_rd = BBOr(bt.BB_Piece[color^1][Rook], bt.BB_Piece[color^1][Dragon])
	bb_ret = BBAndOr(bb_ret, bb_rd, GetRookAttacks(bb_occupied, sq))
	bb_ret = BBAndOr(bb_ret, bt.BB_Piece[color^1][Lance], GetLanceAttacks(bb_occupied, sq, color))
	return bb_ret
}

func IsAttackedByLongPieces(bt BoardTree, sq int, color int) BitBoard {
	var bb_ret = BBIni()
	var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
	var bb_bh = BBOr(bt.BB_Piece[color^1][Bishop], bt.BB_Piece[color^1][Horse])
	bb_ret = BBAndOr(bb_ret, bb_bh, GetBishopAttacks(bb_occupied, sq))
	var bb_rd = BBOr(bt.BB_Piece[color^1][Rook], bt.BB_Piece[color^1][Dragon])
	bb_ret = BBAndOr(bb_ret, bb_rd, GetRookAttacks(bb_occupied, sq))
	bb_ret = BBAndOr(bb_ret, bt.BB_Piece[color^1][Lance], GetLanceAttacks(bb_occupied, sq, color))
	return bb_ret
}

func AttacksToLongPiece(bt BoardTree, sq int, color int) BitBoard {
	var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
	var bb_bh = BBOr(bt.BB_Piece[color][Bishop], bt.BB_Piece[color][Horse])
	var bb_ret = BBAnd(bb_bh, GetBishopAttacks(bb_occupied, sq))
	var bb_rd = BBOr(bt.BB_Piece[color][Rook], bt.BB_Piece[color][Dragon])
	bb_ret = BBAndOr(bb_ret, bb_rd, GetRookAttacks(bb_occupied, sq))
	bb_ret = BBAndOr(bb_ret, bt.BB_Piece[color][Lance], GetLanceAttacks(bb_occupied, sq, color^1))
	return bb_ret
}

func IsDiscoverKing(bt *BoardTree, ifrom int, ito int, color int) bool {
	var idirec = Adirec[bt.SQ_King[color]][ifrom]
	if idirec == Direc_Misc && idirec != Adirec[bt.SQ_King[color]][ito] && BBTest(IsPinnedOnKing(bt, ifrom, idirec, color)) != 0 {
		return true
	} else {
		return false
	}
}
func IsDiscoverKing2(bt *BoardTree, ifrom int, ito int, color int, ipiece int) bool {
	var idirec = Adirec[bt.SQ_King[color]][ifrom]
	Xor(&bt.BB_Piece[color][ipiece], ifrom)
	Xor(&bt.BB_Occupied[color], ifrom)
	if idirec != Direc_Misc && idirec != Adirec[bt.SQ_King[color]][ito] && BBTest(IsPinnedOnKing(bt, ifrom, idirec, color)) != 0 {
		Xor(&bt.BB_Piece[color][ipiece], ifrom)
		Xor(&bt.BB_Occupied[color], ifrom)
		return true
	} else {
		Xor(&bt.BB_Piece[color][ipiece], ifrom)
		Xor(&bt.BB_Occupied[color], ifrom)
		return false
	}
}
