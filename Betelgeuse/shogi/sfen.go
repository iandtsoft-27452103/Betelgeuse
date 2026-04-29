package shogi

import (
	"fmt"
	"strings"
)

func ToSFEN(bt BoardTree, color int) string {
	var flag = false
	var i = 0
	var empty_count = 0
	var str_sfen = ""
	for i < Square_NB {
		var str_piece = Str_SFEN_Pc[Piece(bt.Board[i])]
		if str_piece == "" {
			empty_count += 1
			flag = true
		} else {
			if flag == true {
				flag = false
				str_sfen += fmt.Sprintf("%d", empty_count)
				empty_count = 0
			}
			str_sfen += str_piece
		}
		if i != (Square_NB-1) && FileTable[i] == File9 {
			if empty_count > 0 {
				flag = false
				str_sfen += fmt.Sprintf("%d", empty_count)
				empty_count = 0
			}
			str_sfen += "/"
		}
		i += 1
	}
	str_sfen += " "
	str_sfen += Str_Color[Color(color)]
	str_sfen += " "
	var k = 0
	if bt.Hand[Black] == 0 && bt.Hand[White] == 0 {
		str_sfen += "-"
	} else {
		for i = int(Black); i < Color_NB; i++ {
			for j := Rook; j >= Pawn; j-- {
				var num = (bt.Hand[i] & Hand_Mask[j]) >> Hand_Rev_Bit[j]
				if num == 0 {
					continue
				}
				if num > 0 {
					if num == 1 {
						k = -Sign_Table[i] * int(j)
						str_sfen += Str_SFEN_Pc[Piece(k)]
					} else if num > 1 {
						k = -Sign_Table[i] * int(j)
						str_sfen += fmt.Sprintf("%d", num) + Str_SFEN_Pc[Piece(k)]
					}
				}
			}
		}
	}
	str_sfen += " 1"
	return str_sfen
}
func ToBoard(str_sfen string) BoardTree {
	var flag = false
	var bt BoardTree = InitBoard()
	Clear(&bt)
	var int_pc = 0
	var str_temp = strings.Split(str_sfen, " ")
	var str_board = str_temp[0]
	var limit = len(str_board)
	var sq = 0
	for j := 0; j < limit; j++ {
		var s string = str_board[j : j+1]
		if s == "+" {
			flag = true
		} else if s == "/" {
			continue
		} else {
			var limit2 = len(Set_Empty_Num)
			var flag2 = false
			for k := 0; k < limit2; k++ {
				if s == Set_Empty_Num[k] {
					flag2 = true
					break
				}
			}
			if flag2 {
				var empty_num = Int_Empty_Num[s]
				var m = 0
				for m < int(empty_num) {
					bt.Board[sq] = int8(Empty)
					sq += 1
					m += 1
				}
			} else {
				var int_pc = Int_Pc[s]
				if int_pc > 0 {
					if flag {
						int_pc += Promote
						flag = false
					}
					bt.BB_Piece[Black][int_pc] = BBOr(bt.BB_Piece[Black][int_pc], Atk.ABB_Mask[sq])
					bt.BB_Occupied[Black] = BBOr(bt.BB_Occupied[Black], Atk.ABB_Mask[sq])
					if int_pc == King {
						bt.SQ_King[Black] = uint8(sq)
					}
				} else {
					if flag {
						int_pc -= Promote
						flag = false
					}
					bt.BB_Piece[White][-int_pc] = BBOr(bt.BB_Piece[White][-int_pc], Atk.ABB_Mask[sq])
					bt.BB_Occupied[White] = BBOr(bt.BB_Occupied[White], Atk.ABB_Mask[sq])
					if int_pc == -King {
						bt.SQ_King[White] = uint8(sq)
					}
				}
				bt.Board[sq] = int8(int_pc)
				sq += 1
			}
		}
	}
	var str_color = str_temp[1]
	bt.RootColor = uint8(Num_Color[str_color])
	var str_hand = str_temp[2]
	limit = len(str_hand)
	flag = false
	var num = 1
	for j := 0; j < limit; j++ {
		var s = str_hand[j : j+1]
		if s == "-" {
			break
		}
		if s == "1" && !flag {
			flag = true
		} else {
			if flag {
				num = 10 + int(Int_Hand_Num[s])
				flag = false
			} else {
				var limit2 = len(Set_Hand_Num)
				var flag2 = false
				for k := 0; k < limit2; k++ {
					if s == Set_Hand_Num[k] {
						flag2 = true
						break
					}
				}
				if flag2 {
					num = int(Int_Hand_Num[s])
				} else {
					int_pc = int(Int_Pc[s])
					var color Color
					if int_pc > 0 {
						color = Black
					} else {
						color = White
						int_pc = -int_pc
					}
					var k = 0
					for k < num {
						bt.Hand[color] += Hand_Hash[int_pc]
						k += 1
					}
					num = 1
				}
			}
		}
	}
	bt.CurrentHash = HashFunc(bt)
	bt.Hash[0] = bt.PrevHash
	bt.Hash[1] = bt.CurrentHash
	bt.Ply = 1
	return bt
}
