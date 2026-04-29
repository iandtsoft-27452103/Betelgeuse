package shogi

func GenDrop(bt *BoardTree, color int, moves []uint32) []uint32 {
	var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
	var bb_piece_can_drop [8]BitBoard
	var bb_empty = BBNot(bb_occupied)
	bb_empty = BBAnd(bb_empty, BB_Full)
	if (bt.Hand[color] & Hand_Mask[Pawn]) > 0 {
		var bb = BBIni()
		for i := File1; i < NFile; i++ {
			bb = BBAnd(BB_File[i], bt.BB_Piece[color][Pawn])
			if BBTest(bb) == 0 {
				var bb2 = BBNot(bt.BB_Piece[color][Pawn])
				bb2 = BBAnd(bb2, BB_Full)
				bb2 = BBAnd(bb2, bb_empty)
				bb2 = BBAnd(bb2, BB_File[i])
				bb_piece_can_drop[Pawn] = BBOr(bb_piece_can_drop[Pawn], bb2)
			}
		}
		bb_piece_can_drop[Pawn] = BBAnd(bb_piece_can_drop[Pawn], BB_Pawn_Lance_Can_Drop[color])
		var sq = int(bt.SQ_King[color^1]) + Delta_Table[color^1]
		bb = BBIni()
		if sq >= 0 && sq < Square_NB {
			bb = BBAnd(bb_piece_can_drop[Pawn], Atk.ABB_Mask[sq])
			if bt.Board[sq] == int8(Empty) && BBTest(bb) > 0 {
				if IsMatePawnDrop(*bt, sq, color^1) {
					Xor(&bb_piece_can_drop[Pawn], sq)
				}
			}
		}
	}
	bb_piece_can_drop[Lance] = BBAnd(BB_Pawn_Lance_Can_Drop[color], bb_empty)
	bb_piece_can_drop[Knight] = BBAnd(BB_Knight_Can_Drop[color], bb_empty)
	bb_piece_can_drop[Silver] = BBAnd(BB_Others_Can_Drop, bb_empty)
	bb_piece_can_drop[Gold] = bb_piece_can_drop[Silver]
	bb_piece_can_drop[Bishop] = bb_piece_can_drop[Silver]
	bb_piece_can_drop[Rook] = bb_piece_can_drop[Silver]
	for i := Pawn; i <= Rook; i++ {
		if (bt.Hand[color] & Hand_Mask[i]) > 0 {
			var bb = bb_piece_can_drop[i]
			for BBTest(bb) > 0 {
				var ifrom = uint32(Square_NB + i - 1)
				var ito = Square(bb)
				Xor(&bb, ito)
				var move = Pack(ifrom, uint32(ito), uint32(i), 0, 0)
				moves = append(moves, move)
			}
		}
	}
	return moves
}
func GenNoCap(bt BoardTree, color int, moves []uint32) []uint32 {
	var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
	var bb_empty = BBNot(bb_occupied)
	bb_empty = BBAnd(bb_empty, BB_Full)
	var bb_from = bt.BB_Piece[color][Pawn]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][Pawn][ifrom], bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_temp = BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])
			var flag_promo = 0
			if BBTest(bb_temp) > 0 {
				flag_promo = 1
			}
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Pawn), 0, uint32(flag_promo))
			moves = append(moves, move)
		}
	}
	bb_from = bt.BB_Piece[color][Knight]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][Knight][ifrom], bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_can_promote = BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])
			if BBTest(bb_can_promote) > 0 {
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Knight), 0, 1)
				moves = append(moves, move)
			}
			var bb_temp = BBAnd(BB_Knight_Must_Promote[color], Atk.ABB_Mask[ito])
			if BBTest(bb_temp) == 0 {
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Knight), 0, 0)
				moves = append(moves, move)
			}
		}
	}
	bb_from = bt.BB_Piece[color][Silver]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][Silver][ifrom], bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_from_to = BBOr(Atk.ABB_Mask[ifrom], Atk.ABB_Mask[ito])
			var bb_can_promote = BBAnd(BB_Rev_Color_Position[color], bb_from_to)
			if BBTest(bb_can_promote) > 0 {
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Silver), 0, 1)
				moves = append(moves, move)
			}
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Silver), 0, 0)
			moves = append(moves, move)
		}
	}
	var piece_list [6]Piece = [6]Piece{Gold, King, Pro_Pawn, Pro_Lance, Pro_Knight, Pro_Silver}
	var l = len(piece_list)
	for i := 0; i < l; i++ {
		bb_from = bt.BB_Piece[color][piece_list[i]]
		for BBTest(bb_from) > 0 {
			var ifrom = Square(bb_from)
			Xor(&bb_from, ifrom)
			var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][piece_list[i]][ifrom], bb_empty)
			for BBTest(bb_to) > 0 {
				var ito = Square(bb_to)
				Xor(&bb_to, ito)
				var move = Pack(uint32(ifrom), uint32(ito), uint32(piece_list[i]), 0, 0)
				moves = append(moves, move)
			}
		}
	}
	bb_from = bt.BB_Piece[color][Lance]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = GetLanceAttacks(bb_occupied, ifrom, color)
		bb_to = BBAnd(bb_to, bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_can_promote = BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])
			if BBTest(bb_can_promote) > 0 {
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Lance), 0, 1)
				moves = append(moves, move)
			}
			var bb_temp2 = BBAnd(BB_Knight_Must_Promote[color], Atk.ABB_Mask[ito])
			if BBTest(bb_temp2) == 0 {
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Lance), 0, 0)
				moves = append(moves, move)
			}
		}
	}
	bb_from = bt.BB_Piece[color][Bishop]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(GetBishopAttacks(bb_occupied, ifrom), bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_from_to = BBOr(Atk.ABB_Mask[ifrom], Atk.ABB_Mask[ito])
			var bb_temp4 = BBAnd(BB_Rev_Color_Position[color], bb_from_to)
			var flag_promo = 0
			if BBTest(bb_temp4) > 0 {
				flag_promo = 1
			}
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Bishop), 0, uint32(flag_promo))
			moves = append(moves, move)
		}
	}
	bb_from = bt.BB_Piece[color][Horse]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(GetHorseAttacks(bb_occupied, ifrom), bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Horse), 0, 0)
			moves = append(moves, move)
		}
	}
	bb_from = bt.BB_Piece[color][Rook]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(GetRookAttacks(bb_occupied, ifrom), bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_from_to = BBOr(Atk.ABB_Mask[ifrom], Atk.ABB_Mask[ito])
			var bb_temp5 = BBAnd(BB_Rev_Color_Position[color], bb_from_to)
			var flag_promo = 0
			if BBTest(bb_temp5) > 0 {
				flag_promo = 1
			}
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Rook), 0, uint32(flag_promo))
			moves = append(moves, move)
		}
	}
	bb_from = bt.BB_Piece[color][Dragon]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(GetDragonAttacks(bb_occupied, ifrom), bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Dragon), 0, 0)
			moves = append(moves, move)
		}
	}
	return moves
}
func GenCap(bt BoardTree, color int, moves []uint32) []uint32 {
	var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
	var bb_can_cap = bt.BB_Occupied[color^1]
	var bb_from = bt.BB_Piece[color][Pawn]
	var sgn = -Sign_Table[color]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][Pawn][ifrom], bb_can_cap)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_temp = BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])
			var flag_promo = 0
			if BBTest(bb_temp) > 0 {
				flag_promo = 1
			}
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Pawn), uint32(cap_pc), uint32(flag_promo))
			moves = append(moves, move)
		}
	}
	bb_from = bt.BB_Piece[color][Knight]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][Knight][ifrom], bb_can_cap)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_can_promote = BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])
			if BBTest(bb_can_promote) > 0 {
				var cap_pc = -sgn * int(bt.Board[ito])
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Knight), uint32(cap_pc), 1)
				moves = append(moves, move)
			}
			var bb_temp = BBAnd(BB_Knight_Must_Promote[color], Atk.ABB_Mask[ito])
			if BBTest(bb_temp) == 0 {
				var cap_pc = -sgn * int(bt.Board[ito])
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Knight), uint32(cap_pc), 0)
				moves = append(moves, move)
			}
		}
	}
	bb_from = bt.BB_Piece[color][Silver]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][Silver][ifrom], bb_can_cap)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_from_to = BBOr(Atk.ABB_Mask[ifrom], Atk.ABB_Mask[ito])
			var bb_can_promote = BBAnd(BB_Rev_Color_Position[color], bb_from_to)
			var cap_pc = -sgn * int(bt.Board[ito])
			if BBTest(bb_can_promote) > 0 {
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Silver), uint32(cap_pc), 1)
				moves = append(moves, move)
			}
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Silver), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	var piece_list [6]Piece = [6]Piece{Gold, King, Pro_Pawn, Pro_Lance, Pro_Knight, Pro_Silver}
	var l = len(piece_list)
	for i := 0; i < l; i++ {
		bb_from = bt.BB_Piece[color][piece_list[i]]
		for BBTest(bb_from) > 0 {
			var ifrom = Square(bb_from)
			Xor(&bb_from, ifrom)
			var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][piece_list[i]][ifrom], bb_can_cap)
			for BBTest(bb_to) > 0 {
				var ito = Square(bb_to)
				Xor(&bb_to, ito)
				var cap_pc = -sgn * int(bt.Board[ito])
				var move = Pack(uint32(ifrom), uint32(ito), uint32(piece_list[i]), uint32(cap_pc), 0)
				moves = append(moves, move)
			}
		}
	}
	bb_from = bt.BB_Piece[color][Lance]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(GetLanceAttacks(bb_occupied, ifrom, color), Atk.ABB_Lance_Mask_Ex[color][ifrom])
		bb_to = BBAnd(bb_to, bb_can_cap)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_can_promote = BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])
			if BBTest(bb_can_promote) > 0 {
				var cap_pc = -sgn * int(bt.Board[ito])
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Lance), uint32(cap_pc), 1)
				moves = append(moves, move)
			}
			var bb_temp2 = BBAnd(BB_Knight_Must_Promote[color], Atk.ABB_Mask[ito])
			if BBTest(bb_temp2) == 0 {
				var cap_pc = -sgn * int(bt.Board[ito])
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Lance), uint32(cap_pc), 0)
				moves = append(moves, move)
			}
		}
	}
	bb_from = bt.BB_Piece[color][Bishop]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(GetBishopAttacks(bb_occupied, ifrom), bb_can_cap)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_from_to = BBOr(Atk.ABB_Mask[ifrom], Atk.ABB_Mask[ito])
			var bb_temp4 = BBAnd(BB_Rev_Color_Position[color], bb_from_to)
			var flag_promo = 0
			if BBTest(bb_temp4) > 0 {
				flag_promo = 1
			}
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Bishop), uint32(cap_pc), uint32(flag_promo))
			moves = append(moves, move)
		}
	}
	bb_from = bt.BB_Piece[color][Horse]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(GetHorseAttacks(bb_occupied, ifrom), bb_can_cap)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Horse), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	bb_from = bt.BB_Piece[color][Rook]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(GetRookAttacks(bb_occupied, ifrom), bb_can_cap)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var bb_from_to = BBOr(Atk.ABB_Mask[ifrom], Atk.ABB_Mask[ito])
			var bb_temp5 = BBAnd(BB_Rev_Color_Position[color], bb_from_to)
			var flag_promo = 0
			if BBTest(bb_temp5) > 0 {
				flag_promo = 1
			}
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Rook), uint32(cap_pc), uint32(flag_promo))
			moves = append(moves, move)
		}
	}
	bb_from = bt.BB_Piece[color][Dragon]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(GetDragonAttacks(bb_occupied, ifrom), bb_can_cap)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Dragon), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	return moves
}
func GenEvasion(bt *BoardTree, color int, moves []uint32) []uint32 {
	var sgn = -Sign_Table[color]
	var sq_king = int(bt.SQ_King[color])
	var ifrom = sq_king
	Xor(&bt.BB_Occupied[color], ifrom)
	var bb_not_color = BBNot(bt.BB_Occupied[color])
	bb_not_color = BBAnd(bb_not_color, BB_Full)
	var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][King][sq_king], bb_not_color)
	for BBTest(bb_to) > 0 {
		var ito = Square(bb_to)
		if BBTest(IsAttacked(*bt, ito, color)) == 0 {
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(King), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
		Xor(&bb_to, ito)
	}
	Xor(&bt.BB_Occupied[color], ifrom)
	var bb_checker = AttacksToPiece(*bt, sq_king, color^1)
	var checker_num = PopCount(bb_checker)
	if checker_num == 2 {
		return moves
	}
	var sq_checker = Square(bb_checker)
	var bb_cap_checker = AttacksToPiece(*bt, sq_checker, color)
	var ito = sq_checker
	var is_contain = false
	for BBTest(bb_cap_checker) > 0 {
		ifrom = Square(bb_cap_checker)
		Xor(&bb_cap_checker, ifrom)
		if ifrom == sq_king {
			continue
		}
		var ipiece = sgn * int(bt.Board[ifrom])
		var idirec = Adirec[ifrom][ito]
		var flag = false

		if BBTest(IsPinnedOnKing(bt, ifrom, idirec, color)) == 0 {
			is_contain = false
			for i := 0; i < len(Set_Piece_Can_Promote0); i++ {
				if ipiece == int(Set_Piece_Can_Promote0[i]) {
					is_contain = true
					break
				}
			}
			if is_contain {
				if BBTest(BBAnd(Atk.ABB_Piece_Attacks[color][ipiece][ifrom], Atk.ABB_Mask[sq_checker])) > 0 && BBTest(BBAnd(Atk.ABB_Piece_Attacks[color][ipiece][ifrom], BB_Rev_Color_Position[color])) > 0 {
					var cap_pc = -sgn * int(bt.Board[ito])
					var move = Pack(uint32(ifrom), uint32(ito), uint32(ipiece), uint32(cap_pc), 1)
					Do(bt, move, color)
					if BBTest(IsAttacked(*bt, sq_king, color)) == 0 {
						moves = append(moves, move)
					}
					UnDo(bt, move, color)
					if ipiece != int(Pawn) {
						flag = true
					}
				}
			}
			is_contain = false
			for i := 0; i < len(Set_Piece_Can_Promote1); i++ {
				if ipiece == int(Set_Piece_Can_Promote1[i]) {
					is_contain = true
					break
				}
			}
			if is_contain {
				if BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ifrom])) > 0 || BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])) > 0 {
					var cap_pc = -sgn * int(bt.Board[ito])
					var move = Pack(uint32(ifrom), uint32(ito), uint32(ipiece), uint32(cap_pc), 1)
					Do(bt, move, color)
					if BBTest(IsAttacked(*bt, sq_king, color)) == 0 {
						moves = append(moves, move)
					}
					UnDo(bt, move, color)
					if ipiece != int(Silver) {
						flag = true
					}
				}
			}
			if !flag {
				var cap_pc = -sgn * int(bt.Board[ito])
				var move = Pack(uint32(ifrom), uint32(ito), uint32(ipiece), uint32(cap_pc), 0)
				Do(bt, move, color)
				if BBTest(IsAttacked(*bt, sq_king, color)) == 0 {
					if ipiece == int(Pawn) && BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])) > 0 {
					} else {
						moves = append(moves, move)
					}
				}
				UnDo(bt, move, color)
			}
		}
	}
	var checker = sgn * int(bt.Board[sq_checker])
	is_contain = false
	for i := 0; i < len(Set_Long_Attack_Pieces); i++ {
		if checker == int(Set_Long_Attack_Pieces[i]) {
			is_contain = true
			break
		}
	}
	if !is_contain {
		return moves
	}
	if BBTest(BBAnd(bb_checker, Atk.ABB_Piece_Attacks[color][King][sq_king])) > 0 {
		return moves
	} else {
		if is_contain {
			var bb_inter = Atk.ABB_Obstacles[sq_king][sq_checker]
			for BBTest(bb_inter) > 0 {
				ito = Square(bb_inter)
				Xor(&bb_inter, ito)
				var bb_defender = AttacksToPiece(*bt, ito, color)
				for BBTest(bb_defender) > 0 {
					ifrom = Square(bb_defender)
					Xor(&bb_defender, ifrom)
					if ifrom == sq_king {
						continue
					}
					var ipiece = sgn * int(bt.Board[ifrom])
					var idirec = Adirec[sq_king][ifrom]
					var flag = false
					if idirec == Direc_Misc || BBTest(IsPinnedOnKing(bt, ifrom, idirec, color)) == 0 {
						is_contain = false
						for i := 0; i < len(Set_Piece_Can_Promote0); i++ {
							if ipiece == int(Set_Piece_Can_Promote0[i]) {
								is_contain = true
								break
							}
						}
						if is_contain {
							if ipiece != int(Lance) && BBTest(BBAnd(Atk.ABB_Piece_Attacks[color][ipiece][ifrom], BB_Rev_Color_Position[color])) > 0 {
								var cap_pc = -sgn * int(bt.Board[ito])
								var move = Pack(uint32(ifrom), uint32(ito), uint32(ipiece), uint32(cap_pc), 1)
								moves = append(moves, move)
								if ipiece == int(Pawn) {
									flag = true
								}
							} else if ipiece == int(Lance) {
								var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
								var bb0 = GetLanceAttacks(bb_occupied, ifrom, color)
								var bb1 = BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])
								if BBTest(bb0) > 0 && BBTest(bb1) > 0 {
									var cap_pc = -sgn * int(bt.Board[ito])
									var move = Pack(uint32(ifrom), uint32(ito), uint32(ipiece), uint32(cap_pc), 1)
									moves = append(moves, move)
								}
							}
						}
						is_contain = false
						for i := 0; i < len(Set_Piece_Can_Promote1); i++ {
							if ipiece == int(Set_Piece_Can_Promote1[i]) {
								is_contain = true
								break
							}
						}
						if is_contain {
							if BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ifrom])) > 0 || BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])) > 0 {
								var cap_pc = -sgn * int(bt.Board[ito])
								var move = Pack(uint32(ifrom), uint32(ito), uint32(ipiece), uint32(cap_pc), 1)
								moves = append(moves, move)
								if ipiece != int(Silver) {
									flag = true
								}
							}
						}
						if !flag {
							if (ipiece == int(Knight) || ipiece == int(Lance)) && BBTest(BBAnd(BB_Knight_Must_Promote[color], Atk.ABB_Mask[ito])) > 0 {
								continue
							}
							var cap_pc = -sgn * int(bt.Board[ito])
							var move = Pack(uint32(ifrom), uint32(ito), uint32(ipiece), uint32(cap_pc), 0)
							Do(bt, move, color)
							if BBTest(IsAttacked(*bt, sq_king, color)) == 0 {
								moves = append(moves, move)
							}
							UnDo(bt, move, color)
						}
					}
				}
			}
		}
		var bb_empty = Atk.ABB_Obstacles[sq_king][sq_checker]
		var bb_piece_can_drop [8]BitBoard
		for i := 0; i < 8; i++ {
			bb_piece_can_drop[i] = BBIni()
		}
		if bt.Hand[color]&Hand_Mask[Pawn] > 0 {
			for i := File1; i < NFile; i++ {
				if BBTest(BBAnd(BB_File[i], bt.BB_Piece[color][Pawn])) == 0 {
					var bb = BBAnd(BBNot(bt.BB_Piece[color][Pawn]), BB_Full)
					bb = BBAnd(bb, BB_Pawn_Lance_Can_Drop[color])
					bb = BBAnd(bb, bb_empty)
					bb = BBAnd(bb, BB_File[i])
					bb_piece_can_drop[Pawn] = BBOr(bb_piece_can_drop[Pawn], bb)
				}
			}
			var sq = int(bt.SQ_King[color]) + Delta_Table[color]
			if (sq >= 0 && sq < Square_NB) && bt.Board[sq] == int8(Empty) && BBTest(BBAnd(bb_piece_can_drop[Pawn], Atk.ABB_Mask[sq])) == 0 {
				if IsMatePawnDrop(*bt, sq, color) {
					Xor(&bb_piece_can_drop[Pawn], sq)
				}
			}
		}
		bb_piece_can_drop[Lance] = BBAnd(BB_Pawn_Lance_Can_Drop[color], bb_empty)
		bb_piece_can_drop[Knight] = BBAnd(BB_Knight_Can_Drop[color], bb_empty)
		bb_piece_can_drop[Silver] = BBAnd(BB_Others_Can_Drop, bb_empty)
		bb_piece_can_drop[Gold] = bb_piece_can_drop[Silver]
		bb_piece_can_drop[Bishop] = bb_piece_can_drop[Silver]
		bb_piece_can_drop[Rook] = bb_piece_can_drop[Silver]
		for i := int(Pawn); i <= int(Rook); i++ {
			if (bt.Hand[color] & Hand_Mask[i]) > 0 {
				var bb_object = bb_piece_can_drop[i]
				for BBTest(bb_object) > 0 {
					ifrom = Square_NB + i - 1
					ito = Square(bb_object)
					Xor(&bb_object, ito)
					var move = Pack(uint32(ifrom), uint32(ito), uint32(i), 0, 0)
					moves = append(moves, move)
				}
			}
		}
	}
	return moves
}
func GenCheck(bt *BoardTree, color int, moves []uint32) []uint32 {
	var sgn = -Sign_Table[color]
	var opponent_color = color ^ 1
	var sq_opponent_king = bt.SQ_King[opponent_color]
	var sq_object = int(sq_opponent_king) + Delta_Table[opponent_color]
	var sq_pawn = int(sq_opponent_king) + (2 * Delta_Table[opponent_color])
	// pawn move
	var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
	var bb_empty = BBAnd(BBNot(bb_occupied), BB_Full)
	var bb_move_to = BBAnd(BBOr(bt.BB_Occupied[color^1], bb_empty), BB_Full)
	// generate no promote move
	if sq_pawn >= 0 && sq_pawn < Square_NB && (int(bt.Board[sq_pawn]) == Sign_Table[opponent_color]*int(Pawn)) && (BBTest(BBAnd(Atk.ABB_Mask[sq_pawn], BB_Pawn_Mask[color])) > 0) && BBTest(BBAnd(Atk.ABB_Mask[sq_object], bb_move_to)) > 0 {
		var cap_pc = -sgn * int(bt.Board[sq_pawn])
		var move = Pack(uint32(sq_pawn), uint32(sq_object), uint32(Pawn), uint32(cap_pc), 0)
		moves = append(moves, move)
	}
	// pawn promote move
	var bb_from = bt.BB_Piece[color][Pawn]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Piece_Attacks[color][Pawn][ifrom])
		bb_to = BBAnd(bb_to, Atk.ABB_Piece_Attacks[opponent_color][King][sq_opponent_king])
		bb_to = BBAnd(bb_to, bb_move_to)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Pawn), uint32(cap_pc), 1)
			moves = append(moves, move)
		}
	}
	// pawn move with discovered check, using dragon or rook attacks.
	bb_from = BBAnd(Atk.ABB_Rank_Attacks[sq_opponent_king][0], bt.BB_Piece[color][Pawn])
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		if BBTest(BBAnd(BBOr(bt.BB_Piece[color][Rook], bt.BB_Piece[color][Dragon]), Atk.ABB_Rank_Attacks[ifrom][0])) > 0 && BBTest(BBAnd(Atk.ABB_Piece_Attacks[color][Pawn][ifrom], bb_move_to)) > 0 {
			var flag_promo uint32 = 0
			if BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Piece_Attacks[color][Pawn][ifrom])) > 0 {
				flag_promo = 1
			}
			var ito = ifrom + Delta_Table[color]
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Pawn), uint32(cap_pc), flag_promo)
			moves = append(moves, move)
		}
	}

	// pawn move with discovered check, using horse or bishop attacks.
	bb_from = BBAnd(Atk.ABB_Diag1_Attacks[sq_opponent_king][0], bt.BB_Piece[color][Pawn])
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_bh = BBOr(bt.BB_Piece[color][Bishop], bt.BB_Piece[color][Horse])
		if BBTest(BBAnd(Atk.ABB_Diag1_Attacks[ifrom][0], bb_bh)) > 0 && BBTest(BBAnd(Atk.ABB_Piece_Attacks[color][Pawn][ifrom], bb_move_to)) > 0 {
			var flag_promo uint32 = 0
			if BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Piece_Attacks[color][Pawn][ifrom])) > 0 {
				flag_promo = 1
			}
			var ito = ifrom + Delta_Table[color]
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Pawn), uint32(cap_pc), flag_promo)
			moves = append(moves, move)
		}
	}
	bb_from = BBAnd(Atk.ABB_Diag2_Attacks[sq_opponent_king][0], bt.BB_Piece[color][Pawn])
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_bh = BBOr(bt.BB_Piece[color][Bishop], bt.BB_Piece[color][Horse])
		if BBTest(BBAnd(Atk.ABB_Diag2_Attacks[ifrom][0], bb_bh)) > 0 && BBTest(BBAnd(Atk.ABB_Piece_Attacks[color][Pawn][ifrom], bb_move_to)) > 0 {
			var flag_promo uint32 = 0
			if BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Piece_Attacks[color][Pawn][ifrom])) > 0 {
				flag_promo = 1
			}
			var ito = ifrom + Delta_Table[color]
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Pawn), uint32(cap_pc), flag_promo)
			moves = append(moves, move)
		}
	}
	var bb_temp = BBIni()
	// drop pawn
	if sq_object >= 0 && sq_object < Square_NB {
		bb_temp = BBAnd(BB_File[FileTable[sq_object]], bt.BB_Piece[color][Pawn])
	}
	if BBTest(bb_temp) == 0 && (sq_object >= 0 && sq_object < Square_NB) && (bt.Hand[color]&Hand_Mask[Pawn]) > 0 && bt.Board[sq_object] == int8(Empty) && !IsMatePawnDrop(*bt, sq_object, color^1) {
		var move = Pack(uint32(Square_NB+Pawn-1), uint32(sq_object), uint32(Pawn), 0, 0)
		moves = append(moves, move)
	}
	// silver move
	bb_from = bt.BB_Piece[color][Silver]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var idirec = Adirec[sq_opponent_king][ifrom]
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][Silver][ifrom], Atk.ABB_Piece_Attacks[opponent_color][Silver][sq_opponent_king])
		bb_to = BBAnd(bb_to, bb_move_to)
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBIni()
			var bb_temp2 = BBIni()
			bb_temp2 = AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king))
			bb_temp2 = BBAnd(bb_temp2, Atk.ABB_Piece_Attacks[color][Silver][ifrom])
			bb_temp2 = BBAnd(bb_temp2, bb_move_to)
			bb_to = BBOr(bb_to, bb_temp2)
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Silver), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	// silver promote move
	bb_from = bt.BB_Piece[color][Silver]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var idirec = Adirec[sq_opponent_king][ifrom]
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][Silver][ifrom], Atk.ABB_Piece_Attacks[opponent_color][Gold][sq_opponent_king])
		bb_to = BBAnd(bb_to, bb_move_to)
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBIni()
			var bb_temp2 = BBIni()
			bb_temp2 = AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king))
			bb_temp2 = BBAnd(bb_temp2, Atk.ABB_Piece_Attacks[color][Silver][ifrom])
			bb_temp2 = BBAnd(bb_temp2, bb_move_to)
			bb_to = BBOr(bb_to, bb_temp2)
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			if BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ifrom])) > 0 || BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])) > 0 {
				var cap_pc = -sgn * int(bt.Board[ito])
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Silver), uint32(cap_pc), 1)
				moves = append(moves, move)
			}
		}
	}
	// drop silver
	if (bt.Hand[color] & Hand_Mask[Silver]) > 0 {
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[opponent_color][Silver][sq_opponent_king], bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var move = Pack(uint32(Square_NB+Silver-1), uint32(ito), uint32(Silver), 0, 0)
			moves = append(moves, move)
		}
	}
	// gold move or promoted gold move
	bb_from = BBOr(bt.BB_Piece[color][Gold], bt.BB_Piece[color][Pro_Pawn])
	bb_from = BBOr(bb_from, bt.BB_Piece[color][Pro_Lance])
	bb_from = BBOr(bb_from, bt.BB_Piece[color][Pro_Knight])
	bb_from = BBOr(bb_from, bt.BB_Piece[color][Pro_Silver])
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var idirec = Adirec[sq_opponent_king][ifrom]
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][Gold][ifrom], Atk.ABB_Piece_Attacks[opponent_color][Gold][sq_opponent_king])
		bb_to = BBAnd(bb_to, bb_move_to)
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBIni()
			var bb_temp2 = BBIni()
			bb_temp2 = AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king))
			bb_temp2 = BBAnd(bb_temp2, Atk.ABB_Piece_Attacks[color][Gold][ifrom])
			bb_temp2 = BBAnd(bb_temp2, bb_move_to)
			bb_to = BBOr(bb_to, bb_temp2)
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var pc = sgn * int(bt.Board[ifrom])
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(pc), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	// drop gold
	if (bt.Hand[color] & Hand_Mask[Gold]) > 0 {
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[opponent_color][Gold][sq_opponent_king], bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var move = Pack(uint32(Square_NB+Gold-1), uint32(ito), uint32(Gold), 0, 0)
			moves = append(moves, move)
		}
	}
	// knight move
	bb_from = bt.BB_Piece[color][Knight]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var idirec = Adirec[sq_opponent_king][ifrom]
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][Knight][ifrom], Atk.ABB_Piece_Attacks[opponent_color][Knight][sq_opponent_king])
		bb_to = BBAnd(bb_to, bb_move_to)
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBIni()
			var bb_temp2 = BBIni()
			bb_temp2 = AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king))
			bb_temp2 = BBAnd(bb_temp2, Atk.ABB_Piece_Attacks[color][Knight][ifrom])
			bb_temp2 = BBAnd(bb_temp2, bb_move_to)
			bb_to = BBOr(bb_to, bb_temp2)
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Knight), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	// knight promote move
	bb_from = bt.BB_Piece[color][Knight]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var idirec = Adirec[sq_opponent_king][ifrom]
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][Knight][ifrom], Atk.ABB_Piece_Attacks[opponent_color][Gold][sq_opponent_king])
		bb_to = BBAnd(bb_to, bb_move_to)
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBIni()
			var bb_temp2 = BBIni()
			bb_temp2 = AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king))
			bb_temp2 = BBAnd(bb_temp2, Atk.ABB_Piece_Attacks[color][Knight][ifrom])
			bb_temp2 = BBAnd(bb_temp2, bb_move_to)
			bb_to = BBOr(bb_to, bb_temp2)
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			if BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ifrom])) > 0 || BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ito])) > 0 {
				var cap_pc = -sgn * int(bt.Board[ito])
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Knight), uint32(cap_pc), 1)
				moves = append(moves, move)
			}
		}
	}
	// drop knight
	if (bt.Hand[color] & Hand_Mask[Knight]) > 0 {
		var bb_to = BBAnd(Atk.ABB_Piece_Attacks[opponent_color][Knight][sq_opponent_king], bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var move = Pack(uint32(Square_NB+Knight-1), uint32(ito), uint32(Knight), 0, 0)
			moves = append(moves, move)
		}
	}
	// king move
	var sq_king = bt.SQ_King[color]
	var temp_idirec = Adirec[sq_opponent_king][sq_king]
	if temp_idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, int(sq_king), temp_idirec, opponent_color)) > 0 {
		var bb_temp = BBIni()
		var bb_to = BBAnd(AddBehindAttacks(bb_temp, Direction(temp_idirec), int(sq_opponent_king)), Atk.ABB_Piece_Attacks[color][King][sq_king])
		bb_to = BBAnd(bb_to, bb_move_to)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(sq_king), uint32(ito), uint32(King), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	// lance move => except discovered check move, there must be capture move.
	bb_from = bt.BB_Piece[color][Lance]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = GetLanceAttacks(bb_occupied, ifrom, color)
		var bb_temp2 = BBNot(BB_Knight_Must_Promote[color])
		bb_temp2 = BBAnd(bb_temp2, BB_Full)
		bb_temp2 = BBAnd(bb_temp2, bt.BB_Occupied[opponent_color])
		bb_temp2 = BBAnd(bb_temp2, bb_move_to)
		bb_to = BBAnd(bb_to, bb_temp2)
		var bb_attacks = bb_to
		bb_to = BBAnd(bb_to, GetLanceAttacks(bb_occupied, int(sq_opponent_king), color^1))
		var idirec = Adirec[sq_opponent_king][ifrom]
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			//bb_temp = BBAnd(bb_attacks, AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king)))
			bb_temp = BBOr(bb_attacks, AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king)))
			bb_to = BBOr(bb_to, bb_temp)
			if color == int(Black) {
				var bb_temp3 = BB_File[FileTable[ifrom]]
				var bb_temp4 = BBOr(BB_Rank[2], BB_Rank[3])
				bb_temp3 = BBAnd(bb_temp3, bb_temp4)
				bb_to = BBAnd(bb_to, bb_temp3)
			} else {
				var bb_temp3 = BB_File[FileTable[ifrom]]
				var bb_temp4 = BBOr(BB_Rank[6], BB_Rank[5])
				bb_temp3 = BBAnd(bb_temp3, bb_temp4)
				bb_to = BBAnd(bb_to, bb_temp3)
			}
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Lance), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	// lance promote move
	bb_from = bt.BB_Piece[color][Lance]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = GetLanceAttacks(bb_occupied, ifrom, color)
		var bb_attacks = bb_to
		var bb_temp2 = BB_Rev_Color_Position[color]
		bb_temp2 = BBAnd(bb_temp2, BB_Full)
		bb_temp2 = BBAnd(bb_temp2, Atk.ABB_Piece_Attacks[opponent_color][Gold][sq_opponent_king])
		bb_temp2 = BBAnd(bb_temp2, bb_move_to)
		bb_to = BBAnd(bb_to, bb_temp2)
		bb_to = BBAnd(bb_to, GetLanceAttacks(bb_occupied, int(sq_opponent_king), color^1))
		var idirec = Adirec[sq_opponent_king][ifrom]
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBAnd(bb_attacks, AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king)))
			bb_to = BBOr(bb_to, bb_temp)
			var bb_temp3 = BBOr(BB_Color_Position[Black], BB_Color_Position[White])
			bb_to = BBAnd(bb_to, bb_temp3)
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Lance), uint32(cap_pc), 1)
			moves = append(moves, move)
		}
	}
	// drop lance
	if (bt.Hand[color] & Hand_Mask[Lance]) > 0 {
		var bb_to = BBAnd(GetLanceAttacks(bb_occupied, int(sq_opponent_king), color^1), bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var move = Pack(uint32(Square_NB+Lance-1), uint32(ito), uint32(Lance), 0, 0)
			moves = append(moves, move)
		}
	}
	// rook move
	bb_from = BBAnd(bt.BB_Piece[color][Rook], BBOr(BB_Color_Position[color], BB_DMZ))
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = GetRookAttacks(bb_occupied, ifrom)
		var bb_attacks = bb_to
		bb_to = BBAnd(bb_to, bb_move_to)
		var idirec = Adirec[sq_opponent_king][ifrom]
		bb_to = BBAnd(bb_to, GetRookAttacks(bb_occupied, int(sq_opponent_king)))
		bb_to = BBAnd(bb_to, BBOr(BB_Color_Position[color], BB_DMZ))
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBAnd(bb_attacks, AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king)))
			bb_to = BBOr(bb_to, bb_temp)
			bb_to = BBAnd(bb_to, BBOr(BB_Color_Position[color], BB_DMZ))
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Rook), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	// rook promote move
	bb_from = bt.BB_Piece[color][Rook]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = GetRookAttacks(bb_occupied, ifrom)
		var bb_attacks = bb_to
		bb_to = BBAnd(bb_to, bb_move_to)
		var idirec = Adirec[sq_opponent_king][ifrom]
		bb_to = BBAnd(bb_to, GetDragonAttacks(bb_occupied, int(sq_opponent_king)))
		bb_to = BBAnd(bb_to, BBOr(BB_Rev_Color_Position[color], BB_DMZ))
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBAnd(bb_attacks, AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king)))
			bb_to = BBOr(bb_to, bb_temp)
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			if BBTest(BBAnd(Atk.ABB_Mask[ifrom], BB_Rev_Color_Position[color])) > 0 || BBTest(BBAnd(Atk.ABB_Mask[ito], BB_Rev_Color_Position[color])) > 0 {
				var cap_pc = -sgn * int(bt.Board[ito])
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Rook), uint32(cap_pc), 1)
				moves = append(moves, move)
			}
		}
	}
	// drop rook
	if (bt.Hand[color] & Hand_Mask[Rook]) > 0 {
		var bb_to = BBAnd(GetRookAttacks(bb_occupied, int(sq_opponent_king)), bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var move = Pack(uint32(Square_NB+Rook-1), uint32(ito), uint32(Rook), 0, 0)
			moves = append(moves, move)
		}
	}
	// bishop move
	bb_from = BBAnd(bt.BB_Piece[color][Bishop], BBOr(BB_Color_Position[color], BB_DMZ))
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = GetBishopAttacks(bb_occupied, ifrom)
		var bb_attacks = bb_to
		bb_to = BBAnd(bb_to, bb_move_to)
		var idirec = Adirec[sq_opponent_king][ifrom]
		bb_to = BBAnd(bb_to, GetBishopAttacks(bb_occupied, int(sq_opponent_king)))
		bb_to = BBAnd(bb_to, BBOr(BB_Color_Position[color], BB_DMZ))
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBAnd(bb_attacks, AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king)))
			bb_to = BBOr(bb_to, bb_temp)
			bb_to = BBAnd(bb_to, BBOr(BB_Color_Position[color], BB_DMZ))
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Bishop), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	// bishop promote move
	bb_from = bt.BB_Piece[color][Bishop]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = GetBishopAttacks(bb_occupied, ifrom)
		var bb_attacks = bb_to
		bb_to = BBAnd(bb_to, bb_move_to)
		var idirec = Adirec[sq_opponent_king][ifrom]
		bb_to = BBAnd(bb_to, GetHorseAttacks(bb_occupied, int(sq_opponent_king)))
		bb_to = BBAnd(bb_to, BBOr(BB_Rev_Color_Position[color], BB_DMZ))
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBAnd(bb_attacks, AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king)))
			bb_to = BBOr(bb_to, bb_temp)
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			if BBTest(BBAnd(Atk.ABB_Mask[ifrom], BB_Rev_Color_Position[color])) > 0 || BBTest(BBAnd(Atk.ABB_Mask[ito], BB_Rev_Color_Position[color])) > 0 {
				var cap_pc = -sgn * int(bt.Board[ito])
				var move = Pack(uint32(ifrom), uint32(ito), uint32(Bishop), uint32(cap_pc), 1)
				moves = append(moves, move)
			}
		}
	}
	// drop bishop
	if (bt.Hand[color] & Hand_Mask[Bishop]) > 0 {
		var bb_to = BBAnd(GetBishopAttacks(bb_occupied, int(sq_opponent_king)), bb_empty)
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var move = Pack(uint32(Square_NB+Bishop-1), uint32(ito), uint32(Bishop), 0, 0)
			moves = append(moves, move)
		}
	}
	// dragon move
	bb_from = bt.BB_Piece[color][Dragon]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = GetDragonAttacks(bb_occupied, ifrom)
		var bb_attacks = bb_to
		bb_to = BBAnd(bb_to, bb_move_to)
		var idirec = Adirec[sq_opponent_king][ifrom]
		bb_to = BBAnd(bb_to, GetDragonAttacks(bb_occupied, int(sq_opponent_king)))
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBIni()
			bb_to = BBOr(bb_to, BBAnd(bb_attacks, AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king))))
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Dragon), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	// horse move
	bb_from = bt.BB_Piece[color][Horse]
	for BBTest(bb_from) > 0 {
		var ifrom = Square(bb_from)
		Xor(&bb_from, ifrom)
		var bb_to = GetHorseAttacks(bb_occupied, ifrom)
		var bb_attacks = bb_to
		bb_to = BBAnd(bb_to, bb_move_to)
		var idirec = Adirec[sq_opponent_king][ifrom]
		bb_to = BBAnd(bb_to, GetHorseAttacks(bb_occupied, int(sq_opponent_king)))
		if idirec != int(Direc_Misc) && BBTest(IsPinnedOnKing(bt, ifrom, idirec, opponent_color)) > 0 {
			bb_temp = BBIni()
			bb_to = BBOr(bb_to, BBAnd(bb_attacks, AddBehindAttacks(bb_temp, Direction(idirec), int(sq_opponent_king))))
		}
		for BBTest(bb_to) > 0 {
			var ito = Square(bb_to)
			Xor(&bb_to, ito)
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(Horse), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
	}
	return moves
}
func GenCapTreatPiece(bt BoardTree, color int, moves []uint32, threat_move uint32) {
	var sgn = -Sign_Table[color]
	var sq_attacks_to = From(threat_move)
	// If mate by drop move, defense side do not capture threat piece.
	if sq_attacks_to >= Square_NB {
		return
	}
	var bb_attacks_from = AttacksToPiece(bt, int(sq_attacks_to), color)
	for BBTest(bb_attacks_from) > 0 {
		var ifrom = Square(bb_attacks_from)
		Xor(&bb_attacks_from, ifrom)
		var ipiece = sgn * int(bt.Board[ifrom])
		var cap_pc = -sgn * int(bt.Board[sq_attacks_to])
		var move = Pack(uint32(ifrom), uint32(sq_attacks_to), uint32(ipiece), uint32(cap_pc), 0)
		moves = append(moves, move)
		if BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[sq_attacks_to])) > 0 && ipiece < int(King) {
			cap_pc = -sgn * int(bt.Board[sq_attacks_to])
			move = Pack(uint32(ifrom), uint32(sq_attacks_to), uint32(ipiece), uint32(cap_pc), 1)
			moves = append(moves, move)
		}
		var flag = false
		for i := 0; i < len(Set_Piece_Can_Promote1); i++ {
			var j = Set_Piece_Can_Promote1[i]
			if j == Silver || j == Bishop || j == Rook {
				flag = true
				break
			}
		}
		if (BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ifrom]))) > 0 && flag {
			cap_pc = -sgn * int(bt.Board[sq_attacks_to])
			move = Pack(uint32(ifrom), uint32(sq_attacks_to), uint32(ipiece), uint32(cap_pc), 1)
			moves = append(moves, move)
		}
	}
}
func GenKingMove(bt *BoardTree, color int, moves []uint32) {
	var sgn = -Sign_Table[color]
	var sq_king = int(bt.SQ_King[color])
	var ifrom = sq_king
	Xor(&bt.BB_Occupied[color], ifrom)
	var bb_not_color = BBNot(bt.BB_Occupied[color])
	bb_not_color = BBAnd(bb_not_color, BB_Full)
	var bb_to = BBAnd(Atk.ABB_Piece_Attacks[color][King][sq_king], bb_not_color)
	for BBTest(bb_to) > 0 {
		var ito = Square(bb_to)
		if BBTest(IsAttacked(*bt, ito, color)) == 0 {
			var cap_pc = -sgn * int(bt.Board[ito])
			var move = Pack(uint32(ifrom), uint32(ito), uint32(King), uint32(cap_pc), 0)
			moves = append(moves, move)
		}
		Xor(&bb_to, ito)
	}
	Xor(&bt.BB_Occupied[color], ifrom)
}
func GenInterfere(bt BoardTree, color int, moves []uint32, threat_move uint32) {
	var sgn = -Sign_Table[color]
	var sq_threat_to = To(threat_move)
	var bb_empty = BBAnd(BBNot(BBOr(bt.BB_Occupied[0], bt.BB_Occupied[1])), BB_Full)
	var bb_object = BBOr(bb_empty, bt.BB_Occupied[color^1])
	var bb_attacks_from = AttacksToPiece(bt, int(sq_threat_to), color)
	if BBTest(BBAnd(bb_object, Atk.ABB_Mask[sq_threat_to])) > 0 {
		for BBTest(bb_attacks_from) > 0 {
			var ifrom = Square(bb_attacks_from)
			Xor(&bb_attacks_from, ifrom)
			var ipiece = bt.Board[ifrom]
			if ipiece < 0 {
				ipiece = -ipiece
			}
			if ipiece == int8(King) {
				continue
			}
			var cap_pc = -sgn * int(bt.Board[sq_threat_to])
			var move = Pack(uint32(ifrom), uint32(sq_threat_to), uint32(ipiece), uint32(cap_pc), 0)
			moves = append(moves, move)
			if BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[sq_threat_to])) > 0 && ipiece < int8(King) {
				var move = Pack(uint32(ifrom), uint32(sq_threat_to), uint32(ipiece), uint32(cap_pc), 1)
				moves = append(moves, move)
			}
			var flag = false
			for i := 0; i < len(Set_Piece_Can_Promote1); i++ {
				var j = Set_Piece_Can_Promote1[i]
				if j == Silver || j == Bishop || j == Rook {
					flag = true
					break
				}
			}
			if BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[ifrom])) > 0 && flag {
				var move = Pack(uint32(ifrom), uint32(sq_threat_to), uint32(ipiece), uint32(cap_pc), 1)
				moves = append(moves, move)
			}
		}
	}
}
func AddBehindAttacks(bb BitBoard, idirec Direction, ik int) BitBoard {
	var bb_tmp = BBIni()
	var idirec2 = idirec
	if idirec < 0 {
		idirec2 = -idirec2
	}
	switch idirec2 {
	case Direc_Diag1_U2d:
		bb_tmp = Atk.ABB_Diag1_Attacks[ik][0]
	case Direc_Diag2_U2d:
		bb_tmp = Atk.ABB_Diag2_Attacks[ik][0]
	case Direc_File_U2d:
		bb_tmp = Atk.ABB_File_Attacks[ik][0]
	case Direc_Rank_L2r:
		bb_tmp = Atk.ABB_Rank_Attacks[ik][0]
	}
	bb_tmp = BBAnd(BB_Full, BBNot(bb_tmp))
	return BBOr(bb_tmp, bb)
}
