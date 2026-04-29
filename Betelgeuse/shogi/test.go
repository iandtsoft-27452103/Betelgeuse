package shogi

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// 詰み探索を複数のgoroutineで動かすテスト用コードは後ほど書く。

func ReadTestData(file_name string) ([]string, []string) {
	var comments = make([]string, 0)
	var data = make([]string, 0)
	fp, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	reader := csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	var flag = 0
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			if flag == 0 {
				comments = append(comments, line[0])
			} else {
				data = append(data, line[0])
			}
		}
		flag ^= 1
	}
	return comments, data
}

// テスト用のメソッドはCopilot先生に回答いただいたものを改変した。
func TestGenDrop() {
	var comments, data = ReadTestData("test_data_drop.txt")
	fp, err := os.Create("debug_log.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	// bufio.Writerで効率的に書き込み
	writer := bufio.NewWriter(fp)
	fmt.Println(comments)
	fmt.Println(data)
	var bt = InitBoard()
	var moves []uint32 = make([]uint32, 0)
	for i := 0; i < len(data); i++ {
		bt = ToBoard(data[i])
		moves = make([]uint32, 0)
		moves = GenDrop(&bt, int(bt.RootColor), moves)

		_, err = writer.WriteString(comments[i] + "\n")

		// 書き込み処理
		str_moves := []string{}
		for j := 0; j < len(moves); j++ {
			str_moves = append(str_moves, Move2CSA(moves[j]))
		}

		for j := 0; j < len(moves); j++ {
			if j != len(moves)-1 {
				_, err = writer.WriteString(str_moves[j] + ",")
			} else {
				_, err = writer.WriteString(str_moves[j] + "\n")
			}

			if err != nil {
				fmt.Println("書き込みエラー:", err)
				return
			}
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("フラッシュエラー:", err)
			return
		}
		fmt.Println(bt.CurrentHash)
		fmt.Println(moves)
	}
}
func TestGenNoCap() {
	var comments, data = ReadTestData("test_data_gennocap.txt")
	fp, err := os.Create("debug_log.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	// bufio.Writerで効率的に書き込み
	writer := bufio.NewWriter(fp)
	fmt.Println(comments)
	fmt.Println(data)
	var bt = InitBoard()
	var moves []uint32 = make([]uint32, 0)
	for i := 0; i < len(data); i++ {
		bt = ToBoard(data[i])
		moves = make([]uint32, 0)
		moves = GenNoCap(bt, int(bt.RootColor), moves)

		_, err = writer.WriteString(comments[i] + "\n")

		// 書き込み処理
		str_moves := []string{}
		for j := 0; j < len(moves); j++ {
			str_moves = append(str_moves, Move2CSA(moves[j]))
		}

		for j := 0; j < len(moves); j++ {
			if j != len(moves)-1 {
				_, err = writer.WriteString(str_moves[j] + ",")
			} else {
				_, err = writer.WriteString(str_moves[j] + "\n")
			}

			if err != nil {
				fmt.Println("書き込みエラー:", err)
				return
			}
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("フラッシュエラー:", err)
			return
		}
		fmt.Println(bt.CurrentHash)
		fmt.Println(moves)
	}
}
func TestGenCap() {
	var comments, data = ReadTestData("test_data_gencap.txt")
	fp, err := os.Create("debug_log.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	// bufio.Writerで効率的に書き込み
	writer := bufio.NewWriter(fp)
	fmt.Println(comments)
	fmt.Println(data)
	var bt = InitBoard()
	var moves []uint32 = make([]uint32, 0)
	for i := 0; i < len(data); i++ {
		bt = ToBoard(data[i])
		moves = make([]uint32, 0)
		moves = GenCap(bt, int(bt.RootColor), moves)

		_, err = writer.WriteString(comments[i] + "\n")

		// 書き込み処理
		str_moves := []string{}
		for j := 0; j < len(moves); j++ {
			str_moves = append(str_moves, Move2CSA(moves[j]))
		}

		for j := 0; j < len(moves); j++ {
			if j != len(moves)-1 {
				_, err = writer.WriteString(str_moves[j] + ",")
			} else {
				_, err = writer.WriteString(str_moves[j] + "\n")
			}

			if err != nil {
				fmt.Println("書き込みエラー:", err)
				return
			}
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("フラッシュエラー:", err)
			return
		}
		fmt.Println(bt.CurrentHash)
		fmt.Println(moves)
	}
}
func TestGenEvasion() {
	var comments, data = ReadTestData("test_data_evasion.txt")
	fp, err := os.Create("debug_log.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	// bufio.Writerで効率的に書き込み
	writer := bufio.NewWriter(fp)
	fmt.Println(comments)
	fmt.Println(data)
	var bt = InitBoard()
	var moves []uint32 = make([]uint32, 0)
	for i := 0; i < len(data); i++ {
		bt = ToBoard(data[i])
		moves = make([]uint32, 0)
		fmt.Println("baka")
		fmt.Println(i)
		fmt.Println(comments[i])
		moves = GenEvasion(&bt, int(bt.RootColor), moves)

		_, err = writer.WriteString(comments[i] + "\n")

		// 書き込み処理
		str_moves := []string{}
		for j := 0; j < len(moves); j++ {
			str_moves = append(str_moves, Move2CSA(moves[j]))
		}

		for j := 0; j < len(moves); j++ {
			if j != len(moves)-1 {
				_, err = writer.WriteString(str_moves[j] + ",")
			} else {
				_, err = writer.WriteString(str_moves[j] + "\n")
			}

			if err != nil {
				fmt.Println("書き込みエラー:", err)
				return
			}
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("フラッシュエラー:", err)
			return
		}
		fmt.Println(bt.CurrentHash)
		fmt.Println(moves)
	}
}
func TestGenCheck() {
	var comments, data = ReadTestData("test_data_check.txt")
	fp, err := os.Create("debug_log.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	// bufio.Writerで効率的に書き込み
	writer := bufio.NewWriter(fp)
	fmt.Println(comments)
	fmt.Println(data)
	var bt = InitBoard()
	var moves []uint32 = make([]uint32, 0)
	for i := 0; i < len(data); i++ {
		bt = ToBoard(data[i])
		moves = make([]uint32, 0)
		fmt.Println("baka")
		fmt.Println(i)
		moves = GenCheck(&bt, int(bt.RootColor), moves)

		_, err = writer.WriteString(comments[i] + "\n")

		// 書き込み処理
		str_moves := []string{}
		for j := 0; j < len(moves); j++ {
			str_moves = append(str_moves, Move2CSA(moves[j]))
		}

		for j := 0; j < len(moves); j++ {
			if j != len(moves)-1 {
				_, err = writer.WriteString(str_moves[j] + ",")
			} else {
				_, err = writer.WriteString(str_moves[j] + "\n")
			}

			if err != nil {
				fmt.Println("書き込みエラー:", err)
				return
			}
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("フラッシュエラー:", err)
			return
		}
		fmt.Println(bt.CurrentHash)
		fmt.Println(moves)
	}
}
func TestGenCheck2() {
	var comments, data = ReadTestData("test_data_b_check_additional.txt")
	fp, err := os.Create("debug_log.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	// bufio.Writerで効率的に書き込み
	writer := bufio.NewWriter(fp)
	fmt.Println(comments)
	fmt.Println(data)
	var bt = InitBoard()
	var moves []uint32 = make([]uint32, 0)
	for i := 0; i < len(data); i++ {
		bt = ToBoard(data[i])
		moves = make([]uint32, 0)
		fmt.Println("baka")
		fmt.Println(i)
		moves = GenCheck(&bt, int(bt.RootColor), moves)

		_, err = writer.WriteString(comments[i] + "\n")

		// 書き込み処理
		str_moves := []string{}
		for j := 0; j < len(moves); j++ {
			str_moves = append(str_moves, Move2CSA(moves[j]))
		}

		for j := 0; j < len(moves); j++ {
			if j != len(moves)-1 {
				_, err = writer.WriteString(str_moves[j] + ",")
			} else {
				_, err = writer.WriteString(str_moves[j] + "\n")
			}

			if err != nil {
				fmt.Println("書き込みエラー:", err)
				return
			}
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("フラッシュエラー:", err)
			return
		}
		fmt.Println(bt.CurrentHash)
		fmt.Println(moves)
	}
}
func TestGenCheck3() {
	var comments, data = ReadTestData("test_data_w_check_additional.txt")
	fp, err := os.Create("debug_log.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	// bufio.Writerで効率的に書き込み
	writer := bufio.NewWriter(fp)
	fmt.Println(comments)
	fmt.Println(data)
	var bt = InitBoard()
	var moves []uint32 = make([]uint32, 0)
	for i := 0; i < len(data); i++ {
		bt = ToBoard(data[i])
		moves = make([]uint32, 0)
		fmt.Println("baka")
		fmt.Println(i)
		moves = GenCheck(&bt, int(bt.RootColor), moves)

		_, err = writer.WriteString(comments[i] + "\n")

		// 書き込み処理
		str_moves := []string{}
		for j := 0; j < len(moves); j++ {
			str_moves = append(str_moves, Move2CSA(moves[j]))
		}

		for j := 0; j < len(moves); j++ {
			if j != len(moves)-1 {
				_, err = writer.WriteString(str_moves[j] + ",")
			} else {
				_, err = writer.WriteString(str_moves[j] + "\n")
			}

			if err != nil {
				fmt.Println("書き込みエラー:", err)
				return
			}
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("フラッシュエラー:", err)
			return
		}
		fmt.Println(bt.CurrentHash)
		fmt.Println(moves)
	}
}
func TestMate1Ply() {
	var comments, data = ReadTestData("test_data_b_mate1ply.txt")
	fp, err := os.Create("debug_log.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	// bufio.Writerで効率的に書き込み
	writer := bufio.NewWriter(fp)
	fmt.Println(comments)
	fmt.Println(data)
	var bt = InitBoard()
	var mate_move uint32 = 0
	for i := 0; i < len(data); i++ {
		bt = ToBoard(data[i])
		fmt.Println("baka")
		fmt.Println(i)
		mate_move = MateIn1Ply(bt, int(bt.RootColor))

		_, err = writer.WriteString(comments[i] + "\n")

		if mate_move != 0 {
			var str_move = Move2CSA(mate_move)
			_, err = writer.WriteString(str_move + "\n")
		} else {
			_, err = writer.WriteString("\n")
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("フラッシュエラー:", err)
			return
		}
		fmt.Println(bt.CurrentHash)
		fmt.Println(mate_move)
	}
}
func TestMate1Ply2() {
	var comments, data = ReadTestData("test_data_w_mate1ply.txt")
	fp, err := os.Create("debug_log.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	// bufio.Writerで効率的に書き込み
	writer := bufio.NewWriter(fp)
	fmt.Println(comments)
	fmt.Println(data)
	var bt = InitBoard()
	var mate_move uint32 = 0
	for i := 0; i < len(data); i++ {
		bt = ToBoard(data[i])
		fmt.Println("baka")
		fmt.Println(i)
		mate_move = MateIn1Ply(bt, int(bt.RootColor))

		_, err = writer.WriteString(comments[i] + "\n")

		if mate_move != 0 {
			var str_move = Move2CSA(mate_move)
			_, err = writer.WriteString(str_move + "\n")
		} else {
			_, err = writer.WriteString("\n")
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("フラッシュエラー:", err)
			return
		}
		fmt.Println(bt.CurrentHash)
		fmt.Println(mate_move)
	}
}

func TestMateSearch() {
	var str_sfen [13]string
	str_sfen[0] = "6s2/6R2/6Bk1/6p2/7N1/9/9/9/9 b GN 1"
	str_sfen[1] = "9/9/9/9/1n7/2P6/1Kb6/2r6/2S6 w gn 1"
	str_sfen[2] = "5s1nl/5s1k1/4+Pp1pp/6R2/7N1/9/9/9/9 b 2GN 1"
	str_sfen[3] = "9/9/9/9/1n7/2r6/PP1P+p4/1K1S5/LN1S5 w 2gn 1"
	str_sfen[4] = "8l/7R1/5Sg1k/6N1s/6p2/9/9/9/9 b BGS 1"
	str_sfen[5] = "9/9/9/9/2P6/S1n6/K1Gs5/1r7/L8 w bgs 1"
	str_sfen[6] = "5p2l/4Br2k/4s1gpp/6N+R1/9/9/9/9/9 b 2S 1"
	str_sfen[7] = "9/9/9/9/9/1+rn6/PPG1S4/K2Rb4/L2P5 w 2s 1"
	str_sfen[8] = "7nl/4kl1+R1/1+L6p/2pppPp2/1s3p3/1pPP2P2/n1G1B3P/2G6/2K1+b4 b RG2SNL2Pgsn4p 1"
	str_sfen[9] = "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"
	str_sfen[10] = "lnsgkgsnl/1r5b1/ppppppppp/9/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL w - 1"              // not tested.
	str_sfen[11] = "lnsgk2nl/1r4gs1/p1pppp1pp/6p2/1p5P1/2P6/PPSPPPP1P/7R1/LN1GKGSNL b Bb 1"         // not tested.
	str_sfen[12] = "lr5nl/3g1kg2/2n1pssp1/p1p2pp1p/1p1PP2P1/P1P2PP1P/1PSS2N2/1KG2G3/LN2R3L w BPb 1" // not tested.
	var mst MateSearchTree
	mst.bt = InitBoard()
	Clear(&mst.bt)
	mst.bt = ToBoard(str_sfen[0])
	mst.root_check_moves = GenRootMoves(mst.bt)
	var str_pv = MateSearchWrapper(&mst, 5)
	fmt.Println(str_pv)
}

func TestRepetition() {
	var rs = ReadRecords("test_repetition.txt")
	var r = rs[0]
	var bt = InitBoard()
	var tt TT
	IniTT(&tt)
	var color = 0
	for i := 0; i < len(r.str_moves); i++ {
		var m = CSA2Move(bt, r.str_moves[i])
		fmt.Println(r.str_moves[i])
		Do(&bt, m, color)
		tt.is_check[bt.CurrentHash] = false
		color ^= 1
		var result = IsRepetition(bt, tt)
		fmt.Println(result)
	}
}

func TestDeclarationWin() {
	var str_sfen [11]string
	var bt = InitBoard()
	str_sfen[0] = "+L+NSGKGS+N+L/1+R5+B1/+P+P+P+P+P+P+P+P+P/9/9/9/+p+p+p+p+p+p+p+p+p/1+r5+b1/+l+nsgkgs+n+l b - 1"  //後手勝ち
	str_sfen[1] = "+L+NSGKGS+N+L/+P+R5+B1/+P+P+P+P+P+P+P+P+P/9/9/9/+p+p+p+p+p+p+p+p1/1+r5+b1/+l+nsgkgs+n+l b - 1"  //先手勝ち
	str_sfen[2] = "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1"                                //どちらの勝ちでもない → 初期局面
	str_sfen[3] = "+L+NSGK4/1+R7/+P+P+P+P+P4/9/9/9/+p+p+p+p+p4/7+b1/+l+nsgk4 b BGSNL4Prgsnl4p 1"                   //後手勝ち
	str_sfen[4] = "+L+NSGK4/1+R7/+P+P+P+P+P4/9/9/9/+p+p+p+p+p4/7+b1/+l+nsgk4 b BGSNL5Prgsnl3p 1"                   //先手勝ち
	str_sfen[5] = "4k4/9/9/9/9/9/9/4p4/4K4 b RB2G2S2N2L9Prb2g2s2n2l8p 1"                                           //先手玉が王手
	str_sfen[6] = "4k4/4P4/9/9/9/9/9/9/4K4 b RB2G2S2N2L8Prb2g2s2n2l9p 1"                                           //後手玉が王手
	str_sfen[7] = "+L+NSGKGS+N+L/1+R5+B1/+P+P+P+P+P+P+P+P+P/9/9/8k/+p+p+p+p+p+p+p+p+p/1+r5+b1/+l+nsg1gs+n+l b - 1" //どちらの勝ちでもない → 後手玉が宣言勝ちの位置にいない
	str_sfen[8] = "+L+NSG1GS+N+L/+P+R5+B1/+P+P+P+P+P+P+P+P+P/K8/9/9/+p+p+p+p+p+p+p+p1/1+r5+b1/+l+nsgkgs+n+l b - 1" //どちらの勝ちでもない → 先手玉が宣言勝ちの位置にいない
	str_sfen[9] = "+L+NSGK4/1+R7/+P+P+P+P+P4/9/9/8k/+p+p+p+p+p4/7+b1/+l+nsg5 b BGSNL4Prgsnl4p 1"                   //どちらの勝ちでもない → 後手玉が宣言勝ちの位置にいない
	str_sfen[10] = "+L+NSG5/1+R7/+P+P+P+P+P4/K8/9/9/+p+p+p+p+p4/7+b1/+l+nsgk4 b BGSNL5Prgsnl3p 1"                  //どちらの勝ちでもない → 先手玉が宣言勝ちの位置にいない
	for i := 0; i < len(str_sfen); i++ {
		Clear(&bt)
		bt = ToBoard(str_sfen[i])
		iret := IsDeclarationWin(bt)
		switch iret {
		case 0:
			fmt.Println("宣言勝ちの局面ではない。")
		case 1:
			fmt.Println("先手の勝ち。")
		case 2:
			fmt.Println("後手の勝ち。")
		}
	}
}

func TestDoMove() {
	var rs = ReadRecords("20220403_nhk_hai.txt")
	var r = rs[0]
	var bt = InitBoard()
	var color = 0
	for i := 0; i < len(r.str_moves); i++ {
		var m = CSA2Move(bt, r.str_moves[i])
		Do(&bt, m, color)
		color ^= 1
		if i == 34 {
			continue
			//break
		}
	}
	OutBoard(bt)
}

func TestUnDoMove() {
	var rs = ReadRecords("20220403_nhk_hai.txt")
	var r = rs[0]
	var bt = InitBoard()
	var color = 0
	for i := 0; i < len(r.str_moves); i++ {
		var m = CSA2Move(bt, r.str_moves[i])
		var bt_before = DeepCopy(bt)
		Do(&bt, m, color)
		UnDo(&bt, m, color)
		var bt_after = DeepCopy(bt)
		if VerifyBoard(bt_before, bt_after) == true {
			var s = "i=" + strconv.Itoa(i)
			fmt.Println(s)
			fmt.Println("エラーが発生しました。")
			return
		}
		Do(&bt, m, color)
		color ^= 1
	}
	OutBoard(bt)
}

func OutBoard(bt BoardTree) {
	var s = "盤面の配列"
	fmt.Println(s)
	s = " 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 "
	fmt.Println(s)
	for i := 0; i < Square_NB; i++ {
		var pc = bt.Board[i]
		if FileTable[i] == File1 {
			s = ""
		}
		if pc != 0 {
			if pc > 0 {
				s += " " + Str_Piece_JP[pc]
			} else {
				s += "v" + Str_Piece_JP[-pc]
			}
		} else {
			s += "   "
		}
		if FileTable[i] == File9 {
			s += "|\n"
			fmt.Println(s)
		} else {
			s += "|"
		}
	}
	fmt.Println("")
	s = "ビットボードの駒"
	fmt.Println(s)
	var c = Black
	var bb = bt.BB_Piece[c][Pawn]
	s = "先手の歩："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Lance]
	s = "先手の香："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Knight]
	s = "先手の桂："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Silver]
	s = "先手の銀："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Gold]
	s = "先手の金："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Bishop]
	s = "先手の角："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Rook]
	s = "先手の飛："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][King]
	s = "先手の玉："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Pro_Pawn]
	s = "先手のと："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Pro_Lance]
	s = "先手の成香："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Pro_Knight]
	s = "先手の成桂："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Pro_Silver]
	s = "先手の成銀："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Horse]
	s = "先手の馬："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Dragon]
	s = "先手の龍："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	c = White
	bb = bt.BB_Piece[c][Pawn]
	s = "後手の歩："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Lance]
	s = "後手の香："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Knight]
	s = "後手の桂："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Silver]
	s = "後手の銀："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Gold]
	s = "後手の金："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Bishop]
	s = "後手の角："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Rook]
	s = "後手の飛："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][King]
	s = "後手の玉："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Pro_Pawn]
	s = "後手のと："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Pro_Lance]
	s = "後手の成香："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Pro_Knight]
	s = "後手の成桂："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Pro_Silver]
	s = "後手の成銀："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Horse]
	s = "後手の馬："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	bb = bt.BB_Piece[c][Dragon]
	s = "後手の龍："
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	fmt.Println()
	s = "先手玉の位置：" + strconv.Itoa(int(bt.SQ_King[Black]))
	fmt.Println(s)
	s = "後手玉の位置：" + strconv.Itoa(int(bt.SQ_King[White]))
	fmt.Println(s)
	fmt.Println()
	s = "先手のOccupied："
	fmt.Println(s)
	bb = bt.BB_Occupied[Black]
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	s = "後手のOccupied："
	fmt.Println(s)
	bb = bt.BB_Occupied[White]
	for BBTest(bb) > 0 {
		var sq = Square(bb)
		Xor(&bb, sq)
		s += strconv.Itoa(sq) + ","
	}
	fmt.Println(s)
	fmt.Println()
	s = "先手の持ち駒："
	fmt.Println(s)
	for i := Pawn; i <= Rook; i++ {
		var n = (bt.Hand[0] & Hand_Mask[i]) >> Hand_Rev_Bit[i]
		s = Str_Piece_JP[i] + "：" + strconv.Itoa(int(n))
		fmt.Println(s)
	}
	fmt.Println()
	s = "後手の持ち駒："
	fmt.Println(s)
	for i := Pawn; i <= Rook; i++ {
		var n = (bt.Hand[1] & Hand_Mask[i]) >> Hand_Rev_Bit[i]
		s = Str_Piece_JP[i] + "：" + strconv.Itoa(int(n))
		fmt.Println(s)
	}
}

func TestAttackTables() {
	for sq := 0; sq < Square_NB; sq++ {
		var bb = Atk.ABB_Piece_Attacks[0][King][sq]
		fmt.Println(bb)
	}
}

func VerifyBoard(bt_before BoardTree, bt_after BoardTree) bool {
	for i := Pawn; i <= Dragon; i++ {
		for j := 0; j < 3; j++ {
			if bt_before.BB_Piece[Black][i].P[j] != bt_after.BB_Piece[Black][i].P[j] {
				fmt.Println("先手の" + Str_Piece_JP[i] + "のビットボードが一致しません。")
				return true
			}
			if bt_before.BB_Piece[White][i].P[j] != bt_after.BB_Piece[White][i].P[j] {
				fmt.Println("後手の" + Str_Piece_JP[i] + "のビットボードが一致しません。")
				return true
			}
		}
	}
	for i := 0; i < 3; i++ {
		if bt_before.BB_Occupied[Black].P[i] != bt_after.BB_Occupied[Black].P[i] {
			fmt.Println("先手のOccupiedのビットボードが一致しません。")
			return true
		}
		if bt_before.BB_Occupied[White].P[i] != bt_after.BB_Occupied[White].P[i] {
			fmt.Println("後手のOccupiedのビットボードが一致しません。")
			return true
		}
	}
	if bt_before.SQ_King[Black] != bt_after.SQ_King[Black] {
		fmt.Println("先手玉の位置が一致しません。")
		return true
	}
	if bt_before.SQ_King[White] != bt_after.SQ_King[White] {
		fmt.Println("後手玉の位置が一致しません。")
		return true
	}
	if bt_before.Hand[Black] != bt_after.Hand[Black] {
		fmt.Println("先手の持ち駒が一致しません。")
		return true
	}
	if bt_before.Hand[White] != bt_after.Hand[White] {
		fmt.Println("後手の持ち駒が一致しません。")
		return true
	}
	if bt_before.RootColor != bt_after.RootColor {
		fmt.Println("手番が一致しません。")
		return true
	}
	return false
}
