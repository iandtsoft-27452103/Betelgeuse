package shogi

func MateIn1Ply(bt BoardTree, color int) uint32 {
	var mate_move uint32 = 0
	var null_move uint32 = 0
	var sq_can_check_by_drop [8]int = [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	var sq_can_check_by_move [8]int = [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	var pos_array [10]int = [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var pc_array [10]int = [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var sq_can_escape [8]int = [8]int{0, 0, 0, 0, 0, 0, 0, 0}
	var cnt_d = 0
	var cnt_m = 0
	var cnt_e = 0
	var opponent_color = color ^ 1
	var sq_opponent_king = bt.SQ_King[opponent_color]
	var bb_can_escape = BBAnd(BB_Full, BBNot(bt.BB_Occupied[opponent_color]))
	var hand = bt.Hand[color]
	var bb_opp_king_attacks = Atk.ABB_Piece_Attacks[opponent_color][King][sq_opponent_king]
	var flag = false
	for BBTest(bb_opp_king_attacks) > 0 {
		var sq = Square(bb_opp_king_attacks)
		Xor(&bb_opp_king_attacks, sq)
		var bb_myside_attacks = AttacksToPiece(bt, sq, opponent_color)
		var myside_attacks_count = PopCount(bb_myside_attacks)
		flag = false
		if myside_attacks_count >= 2 && bt.Board[sq] == int8(Empty) {
			//If there are attacks from opponent pieces except king, opponents can capture the checker.
			flag = true
		}
		if BBTest(BBAnd(bb_can_escape, Atk.ABB_Mask[sq])) > 0 {
			// If there are attacks from your pieces, you maybe generate escape move.
			if BBTest(IsAttacked(bt, sq, opponent_color)) == 0 {
				sq_can_escape[cnt_e] = sq
				cnt_e += 1
			}
		}
		if bt.Board[sq] == int8(Empty) && flag == false {
			sq_can_check_by_drop[cnt_d] = sq
			cnt_d += 1
		}
		var bb_enemy_attacks = IsAttacked(bt, sq, color^1)
		if bt.Board[sq] != int8(Empty) && BBTest(BBAnd(bt.BB_Occupied[opponent_color], Atk.ABB_Mask[sq])) > 0 && BBTest(bb_enemy_attacks) > 0 {
			sq_can_check_by_move[cnt_m] = sq
			cnt_m += 1
		}
		if myside_attacks_count < 2 && bt.Board[sq] == int8(Empty) && BBTest(bb_enemy_attacks) > 0 {
			sq_can_check_by_move[cnt_m] = sq
			cnt_m += 1
		}
	}
	for i := 0; i < cnt_d; i++ {
		var sq = sq_can_check_by_drop[i]
		var idirec = Adirec[sq][sq_opponent_king]
		var pt = Piece_Table[opponent_color]
		var bb = AttacksToPiece(bt, sq, opponent_color)
		var cnt_pos = 0
		var cnt_pc = 0
		for BBTest(bb) > 0 {
			var pos = Square(bb)
			Xor(&bb, pos)
			pos_array[cnt_pos] = pos
			cnt_pos += 1
			pc_array[cnt_pc] = int(bt.Board[pos])
			cnt_pc += 1
		}
		var pcs = pt[idirec]
		if hand > 0 {
			for j := 0; j < len(pt); j++ {
				var pc = pcs[j]
				if pc > int(Rook) {
					break
				}
				if pc != int(Pawn) && (hand&Hand_Mask[pc]) > 0 {
					if cnt_e == 0 {
						mate_move = Pack(uint32(Square_NB+pc-1), uint32(sq), uint32(pc), 0, 0)
						return mate_move
					}
					var counter = 0
					var mate_flag = true
					for k := 0; k < cnt_e; k++ {
						var sq_object = sq_can_escape[k]
						if sq == sq_object {
							counter += 1
						}
						if !IsCanEscape(bt, color, sq, Piece(pc), int(sq_opponent_king), sq_object, false) && !IsCanCapture(bt, color, opponent_color, sq, true, -1, Piece(pc)) {
							counter += 1
						} else {
							mate_flag = false
						}
					}
					if counter == cnt_e && mate_flag {
						mate_move = Pack(uint32(Square_NB+pc-1), uint32(sq), uint32(pc), 0, 0)
						return mate_move
					}
				}
			}
		}
	}
	for i := 0; i < cnt_m; i++ {
		var sq = sq_can_check_by_move[i]
		var idirec = Adirec[sq][sq_opponent_king]
		var pt = Piece_Table[opponent_color]
		var bb = AttacksToPiece(bt, sq, color)
		var attacks_count = PopCount(bb)
		if attacks_count < 2 && BBTest(bb) > 0 {
			var pos = Square(bb)
			var bb2 = AttacksToLongPiece(bt, pos, color)
			for BBTest(bb2) > 0 {
				var sq2 = Square(bb2)
				Xor(&bb2, sq2)
				var idirec2 = Adirec[sq2][sq_opponent_king]
				if idirec == idirec2 {
					if cnt_e == 0 {
						var pc = bt.Board[pos]
						if pc < 0 {
							pc = -pc
						}
						if !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
							var cap_pc = bt.Board[sq]
							if cap_pc < 0 {
								cap_pc = -cap_pc
							}
							mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 0)
							return mate_move
						}
					} else if cnt_e == 1 {
						var sq3 = sq_can_escape[0]
						var idirec3 = Adirec[sq3][sq_opponent_king]
						var pc = bt.Board[pos]
						if pc < 0 {
							pc = -pc
						}
						var cap_pc = bt.Board[sq]
						if cap_pc < 0 {
							cap_pc = -cap_pc
						}
						if !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
							if idirec < 0 {
								idirec = -idirec
							}
							if idirec3 < 0 {
								idirec3 = -idirec3
							}
							if idirec == idirec3 {
								switch Direction(idirec) {
								case Direc_File_U2d:
									if pc == int8(Lance) || pc == int8(Rook) || pc == int8(Dragon) {
										mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 0)
										return mate_move
									}
								case Direc_Rank_L2r:
									if pc == int8(Rook) || pc == int8(Dragon) {
										mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 0)
										return mate_move
									}
								case Direc_Diag1_U2d, Direc_Diag2_U2d:
									if pc == int8(Bishop) || pc == int8(Horse) {
										mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 0)
										return mate_move
									}
								}
							}
						}
					}
				} else {
					if cnt_e == 0 && BBTest(BBAnd(Atk.ABB_Piece_Attacks[color][Gold][sq], Atk.ABB_Piece_Attacks[opponent_color][King][sq_opponent_king])) > 0 && BBTest(BBAnd(Atk.ABB_Mask[sq], BB_Color_Position[opponent_color])) > 0 {
						var bb_myside_attacks = AttacksToPiece(bt, sq, opponent_color)
						var myside_attacks_count = PopCount(bb_myside_attacks)
						var idirec3 = Adirec[pos][bt.SQ_King[color]]
						Xor(&bt.BB_Occupied[color], pos)
						var bb = IsPinnedOnKing(&bt, pos, idirec3, color)
						Xor(&bt.BB_Occupied[color], pos)
						var pc = bt.Board[pos]
						if pc < 0 {
							pc = -pc
						}
						var cap_pc = bt.Board[sq]
						if cap_pc < 0 {
							cap_pc = -cap_pc
						}
						if myside_attacks_count < 2 && BBTest(bb) == 0 {
							mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 1)
							return mate_move
						}
					}
				}
			}
			continue
		}
		var cnt_pos = 0
		var cnt_pc = 0
		for BBTest(bb) > 0 {
			var pos = Square(bb)
			Xor(&bb, pos)
			pos_array[cnt_pos] = pos
			cnt_pos += 1
			pc_array[cnt_pc] = int(bt.Board[pos])
			cnt_pc += 1
		}
		var pcs = pt[idirec]
		// This maybe not make sense.
		if cnt_pos == 0 {
			continue
		}
		var index = 0
		for index < cnt_pos {
			var pos = pos_array[index]
			var pc = pc_array[index]
			if pc < 0 {
				pc = -pc
			}
			if pc == int(King) {
				index += 1
				continue
			}
			idirec = Adirec[pos][sq_opponent_king]
			if IsDiscoverKing2(&bt, pos, sq, color, pc) {
				index += 1
				continue
			}
			var flag = false
			for j := 0; j < len(pcs); j++ {
				if pc == pcs[j] {
					flag = true
					break
				}
			}
			if flag {
				var flag2 = false
				for j := 0; j < len(LongPieces2); j++ {
					if pc == int(LongPieces2[j]) {
						flag2 = true
						break
					}
				}
				if flag2 {
					if cnt_e == 0 {
						var flag3 = false
						for j := 0; j < len(LongPieces); j++ {
							if pc == int(LongPieces[j]) {
								flag3 = true
								break
							}
						}
						if flag3 && !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
							var flag_promo = 0
							if BBTest(BBAnd(Atk.ABB_Mask[sq], BB_Color_Position[opponent_color])) > 0 {
								flag_promo = 1
							}
							var cap_pc = bt.Board[sq]
							if cap_pc < 0 {
								cap_pc = -cap_pc
							}
							mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), uint32(flag_promo))
							return mate_move
						}
					}
					flag = false
					for j := 0; j < cnt_e; j++ {
						var sq_object = sq_can_escape[j]
						if sq == sq_object {
							continue
						}
						if !IsCanEscape(bt, color, sq, Piece(pc), int(sq_opponent_king), sq_object, false) && !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
							var flag_promo = 0
							if BBTest(BBAnd(Atk.ABB_Mask[pos], BB_Color_Position[opponent_color])) > 0 || BBTest(BBAnd(Atk.ABB_Mask[sq], BB_Color_Position[opponent_color])) > 0 {
								flag_promo = 1
							}
							var cap_pc = bt.Board[sq]
							if cap_pc < 0 {
								cap_pc = -cap_pc
							}
							mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), uint32(flag_promo))
							flag = true
						} else {
							flag = false
							mate_move = 0
							break
						}
					}
					if flag && mate_move != 0 {
						return mate_move
					}
				} else if pc == int(Dragon) || pc == int(Horse) {
					if cnt_e == 0 {
						if !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
							var cap_pc = bt.Board[sq]
							if cap_pc < 0 {
								cap_pc = -cap_pc
							}
							mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 0)
							return mate_move
						}
					}
					flag = false
					for j := 0; j < cnt_e; j++ {
						var sq_object = sq_can_escape[j]
						if sq == sq_object {
							continue
						}
						if !IsCanEscape(bt, color, sq, Piece(pc), int(sq_opponent_king), sq_object, false) && !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
							var cap_pc = bt.Board[sq]
							if cap_pc < 0 {
								cap_pc = -cap_pc
							}
							mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 0)
							flag = true
						} else {
							flag = false
							mate_move = 0
							break
						}
					}
					if flag && mate_move != 0 {
						return mate_move
					}
				} else {
					switch Piece(pc) {
					case Gold, Pro_Pawn, Pro_Lance, Pro_Knight, Pro_Silver, Silver:
						if cnt_e == 0 {
							if !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
								var cap_pc = bt.Board[sq]
								if cap_pc < 0 {
									cap_pc = -cap_pc
								}
								mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 0)
								return mate_move
							}
						}
						flag = false
						for j := 0; j < cnt_e; j++ {
							var sq_object = sq_can_escape[j]
							if sq == sq_object {
								continue
							}
							if !IsCanEscape(bt, color, sq, Piece(pc), int(sq_opponent_king), sq_object, false) && !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
								var cap_pc = bt.Board[sq]
								if cap_pc < 0 {
									cap_pc = -cap_pc
								}
								mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 0)
								flag = true
							} else {
								flag = false
								mate_move = 0
								break
							}
						}
						if flag && mate_move != 0 {
							return mate_move
						}
					}
				}
				// In this case, it generates pawn or lance move with no promoted.
			}
			if pc > int(Rook) {
				index += 1
				continue
			}
			var pc_promote = pc + Promote
			var flag4 = false
			// knight promote move
			// Knight cannnot mate opponent king from neighbour 8 Square.
			for j := 0; j < len(pcs); j++ {
				if pc_promote == pcs[j] {
					flag4 = true
					break
				}
			}
			if flag4 && pc == int(Knight) && BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[sq])) > 0 {
				if cnt_e == 0 {
					if !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) && BBTest(BBAnd(Atk.ABB_Piece_Attacks[color][Gold][sq], Atk.ABB_Mask[sq_opponent_king])) > 0 {
						var cap_pc = bt.Board[sq]
						if cap_pc < 0 {
							cap_pc = -cap_pc
						}
						mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 1)
						return mate_move
					}
				}
				flag = false
				for j := 0; j < cnt_e; j++ {
					var sq_object = sq_can_escape[j]
					if sq == sq_object {
						continue
					}
					if !IsCanEscape(bt, color, sq, Piece(pc), int(sq_opponent_king), sq_object, true) && !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
						flag = true
						var cap_pc = bt.Board[sq]
						if cap_pc < 0 {
							cap_pc = -cap_pc
						}
						mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 1)
					} else {
						flag = false
						mate_move = 0
						break
					}
				}
				if flag && mate_move != 0 {
					return mate_move
				}
			}
			// lance promote move or pawn promote move
			flag4 = false
			for j := 0; j < len(pcs); j++ {
				if pc_promote == pcs[j] {
					flag4 = true
					break
				}
			}
			var flag5 = false
			for j := 0; j < len(ShortPieces); j++ {
				if pc == int(ShortPieces[j]) {
					flag5 = true
					break
				}
			}
			if flag4 && flag5 && BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[sq])) > 0 {
				if cnt_e == 0 {
					if !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) && BBTest(BBAnd(Atk.ABB_Piece_Attacks[color][Gold][sq], Atk.ABB_Mask[sq_opponent_king])) > 0 {
						var cap_pc = bt.Board[sq]
						if cap_pc < 0 {
							cap_pc = -cap_pc
						}
						mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 1)
						return mate_move
					}
				}
				flag = false
				for j := 0; j < cnt_e; j++ {
					var sq_object = sq_can_escape[j]
					if sq == sq_object {
						continue
					}
					if !IsCanEscape(bt, color, sq, Piece(pc), int(sq_opponent_king), sq_object, true) && !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
						flag = true
						var cap_pc = bt.Board[sq]
						if cap_pc < 0 {
							cap_pc = -cap_pc
						}
						mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 1)
					} else {
						flag = false
						mate_move = 0
						break
					}
				}
				if flag && mate_move != 0 {
					var cap_pc = bt.Board[sq]
					if cap_pc < 0 {
						cap_pc = -cap_pc
					}
					mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 1)
					return mate_move
				}
			}
			// silver promote move
			if pc == int(Silver) {
				flag4 = false
				for j := 0; j < len(pcs); j++ {
					if pc_promote == pcs[j] {
						flag4 = true
						break
					}
				}
				if flag4 && BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[sq])) > 0 || BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[pos])) > 0 {
					if cnt_e == 0 {
						if !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) && BBTest(BBAnd(Atk.ABB_Piece_Attacks[color][Gold][sq], Atk.ABB_Mask[sq_opponent_king])) > 0 {
							var cap_pc = bt.Board[sq]
							if cap_pc < 0 {
								cap_pc = -cap_pc
							}
							mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 1)
							return mate_move
						}
					}
					flag = false
					for j := 0; j < cnt_e; j++ {
						var sq_object = sq_can_escape[j]
						if sq == sq_object {
							continue
						}
						if !IsCanEscape(bt, color, sq, Piece(pc), int(sq_opponent_king), sq_object, true) && !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
							var cap_pc = bt.Board[sq]
							if cap_pc < 0 {
								cap_pc = -cap_pc
							}
							mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 1)
							flag = true
						} else {
							flag = false
							mate_move = 0
							break
						}
					}
					if flag && mate_move != 0 {
						return mate_move
					}
				}

			}
			if pc < int(Bishop) {
				index += 1
				continue
			}
			// rook promote move or bishop promote move
			flag4 = false
			for j := 0; j < len(pcs); j++ {
				if pc_promote == pcs[j] {
					flag4 = true
					break
				}
			}
			flag5 = false
			for j := 0; j < len(LongPieces); j++ {
				if pc == int(LongPieces[j]) {
					flag5 = true
					break
				}
			}
			if flag4 && flag5 && BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[sq])) > 0 || BBTest(BBAnd(BB_Rev_Color_Position[color], Atk.ABB_Mask[pos])) > 0 {
				if cnt_e == 0 {
					if !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) && BBTest(BBAnd(Atk.ABB_Piece_Attacks[color][King][sq], Atk.ABB_Mask[sq_opponent_king])) > 0 {
						var flag_promo = 0
						if BBTest(BBAnd(Atk.ABB_Mask[pos], BB_Color_Position[opponent_color])) > 0 || BBTest(BBAnd(Atk.ABB_Mask[sq], BB_Color_Position[opponent_color])) > 0 {
							flag_promo = 1
						}
						var cap_pc = bt.Board[sq]
						if cap_pc < 0 {
							cap_pc = -cap_pc
						}
						mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), uint32(flag_promo))
						return mate_move
					}
				}
				flag = false
				for j := 0; j < cnt_e; j++ {
					var sq_object = sq_can_escape[j]
					if sq == sq_object {
						continue
					}
					if !IsCanEscape(bt, color, sq, Piece(pc), int(sq_opponent_king), sq_object, true) && !IsCanCapture(bt, color, opponent_color, sq, false, pos, Piece(pc)) {
						var flag_promo = 0
						if BBTest(BBAnd(Atk.ABB_Mask[pos], BB_Color_Position[opponent_color])) > 0 || BBTest(BBAnd(Atk.ABB_Mask[sq], BB_Color_Position[opponent_color])) > 0 {
							flag_promo = 1
						}
						var cap_pc = bt.Board[sq]
						if cap_pc < 0 {
							cap_pc = -cap_pc
						}
						mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), uint32(flag_promo))
						flag = true
					} else {
						flag = false
						mate_move = 0
						break
					}
				}
				if flag && mate_move != 0 {
					return mate_move
				}
			}
			index += 1
		}
	}
	// You cannot mate opponnent king from neighbour 8 square.
	// You maybe mate opponnent move using knight.
	var pc = Knight
	var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
	var bb = BBAnd(Atk.ABB_Piece_Attacks[opponent_color][pc][sq_opponent_king], BBOr(BBAnd(BBNot(bb_occupied), BB_Full), bt.BB_Occupied[opponent_color]))
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		var bb_opponent_attacks_to_sq = AttacksToPiece(bt, sq, opponent_color)
		if (hand&Hand_Mask[pc]) > 0 && bt.Board[sq] == int8(Empty) && cnt_e == 0 && BBTest(bb_opponent_attacks_to_sq) == 0 {
			// drop knight
			mate_move = Pack(uint32(Square_NB+pc-1), uint32(sq), uint32(pc), 0, 0)
			return mate_move
		}
		var bb_my_knight_attacks = BBAnd(Atk.ABB_Piece_Attacks[opponent_color][pc][sq], bt.BB_Piece[color][Knight])
		if BBTest(bb_my_knight_attacks) > 0 && cnt_e == 0 && BBTest(bb_opponent_attacks_to_sq) == 0 {
			var pos = Square(bb_my_knight_attacks)
			Xor(&bb_my_knight_attacks, pos)
			if IsDiscoverKing2(&bt, pos, sq, color, int(pc)) {
				continue
			}
			var cap_pc = bt.Board[sq]
			if cap_pc < 0 {
				cap_pc = -cap_pc
			}
			mate_move = Pack(uint32(pos), uint32(sq), uint32(pc), uint32(cap_pc), 0)
		}
	}
	if mate_move != 0 {
		return mate_move
	}
	return null_move
}
func IsCanEscape(bt BoardTree, color int, sq_checker int, pc_checker Piece, sq_opponent_king int, sq_object int, is_promo bool) bool {
	var bb_occupied = BBOr(bt.BB_Occupied[Black], bt.BB_Occupied[White])
	Xor(&bb_occupied, sq_opponent_king)
	Xor(&bb_occupied, sq_object)
	var bb_attacks = BBIni()
	switch pc_checker {
	case Rook:
		bb_attacks = GetRookAttacks(bb_occupied, sq_checker)
	case Dragon:
		bb_attacks = GetDragonAttacks(bb_occupied, sq_checker)
	case Bishop:
		bb_attacks = GetBishopAttacks(bb_occupied, sq_checker)
	case Horse:
		bb_attacks = GetHorseAttacks(bb_occupied, sq_checker)
	case Pawn, Knight, Silver:
		if is_promo {
			bb_attacks = Atk.ABB_Piece_Attacks[color][Gold][sq_checker]
		} else {
			bb_attacks = Atk.ABB_Piece_Attacks[color][pc_checker][sq_checker]
		}
	case Lance:
		if is_promo {
			bb_attacks = Atk.ABB_Piece_Attacks[color][Gold][sq_checker]
		} else {
			bb_attacks = GetLanceAttacks(bb_occupied, sq_checker, color)
		}
	default:
	}
	bb_attacks = BBAnd(bb_attacks, Atk.ABB_Mask[sq_object])
	if BBTest(bb_attacks) > 0 {
		return false
	}
	return true
}
func IsCanCapture(bt BoardTree, color int, opponent_color int, sq_object int, is_drop bool, ifrom int, ipiece Piece) bool {
	var bb_myside_attacks = AttacksToPiece(bt, sq_object, color)
	var myside_attacks_count = PopCount(bb_myside_attacks)
	var bb_opp_attacks = AttacksToPiece(bt, sq_object, opponent_color)
	var opp_attacks_count = PopCount(bb_opp_attacks)
	if opp_attacks_count > 1 {
		return true
	}
	if (opp_attacks_count == 1) && (myside_attacks_count == 0) {
		//In this case, there is only an opponent king's attacks, but there are no attacks from my pieces to objective square.
		return true
	}
	if opp_attacks_count >= myside_attacks_count {
		if opp_attacks_count == myside_attacks_count && is_drop {
			// In this case, there are two attacks. The one is from opponent king's and the other is from one of my piece,
			// but the check move is drop move.
			return false
		}
		if is_drop {
			return true
		}
		Xor(&bt.BB_Occupied[color], ifrom)
		Xor(&bt.BB_Piece[color][ipiece], ifrom)
		var bb = IsAttacked(bt, int(bt.SQ_King[opponent_color]), color)
		var bb2 = IsAttacked(bt, int(bt.SQ_King[color]), color)
		var bb3 = BBIni()
		switch ipiece {
		case Pawn, Lance, Rook, Dragon:
			var idirec = Adirec[ifrom][sq_object]
			switch Direction(idirec) {
			case Direc_File_U2d, Direc_File_D2u:
				bb3 = IsAttacked(bt, sq_object, color)
			}
		}
		Xor(&bt.BB_Piece[color][ipiece], ifrom)
		Xor(&bt.BB_Occupied[color], ifrom)
		if BBTest(bb2) > 0 {
			return true
		}
		if BBTest(bb) > 0 || BBTest(bb3) > 0 {
			return false
		}
		return true
	}
	return false
}
