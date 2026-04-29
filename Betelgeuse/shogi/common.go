package shogi

import (
	"encoding/csv"
	"io"
	"math/bits"
	"os"
	"strconv"
)

const Square_NB = 81
const Color_NB = 2
const Piece_NB = 16
const No_Pro_Piece_NB = 8
const NFile = 9
const NRank = 9
const Label_NB = 32
const Moves_Max = 700
const Piece_Can_Drop_NB = 7
const Ply_Max = 1024
const Ply_Inc = 8
const Mate_Ply_Max = 128
const Value_Max = 32768
const Value_Min = -Value_Max
const Value_Mate = Value_Max - 1
const Value_Draw = 0
const Promote = 8
const KKP_END = 1091 // 倍にするかも？
const PP_END = 2182
const FV_SCALE = 32
const PV_LIMIT = 128
const Hand_Pawn = 1
const Hand_Lance = 1 << 5
const Hand_Knight = 1 << 8
const Hand_Silver = 1 << 11
const Hand_Gold = 1 << 14
const Hand_Bishop = 1 << 17
const Hand_Rook = 1 << 19
const Long_BB_Hash = 128

var Hand_Hash [Piece_Can_Drop_NB + 1]uint32 = [Piece_Can_Drop_NB + 1]uint32{0, Hand_Pawn, Hand_Lance, Hand_Knight, Hand_Silver, Hand_Gold, Hand_Bishop, Hand_Rook}
var Hand_Rev_Bit [Piece_Can_Drop_NB + 1]uint32 = [Piece_Can_Drop_NB + 1]uint32{0, 0, 5, 8, 11, 14, 17, 19}
var Hand_Mask [Piece_Can_Drop_NB + 1]uint32 = [Piece_Can_Drop_NB + 1]uint32{0, 31, 224, 1792, 14336, 114688, 393216, 1572864}

var Square_Edge [Square_NB]bool = [Square_NB]bool{true, true, true, true, true, true, true, true, true,
	true, false, false, false, false, false, false, false, true,
	true, false, false, false, false, false, false, false, true,
	true, false, false, false, false, false, false, false, true,
	true, false, false, false, false, false, false, false, true,
	true, false, false, false, false, false, false, false, true,
	true, false, false, false, false, false, false, false, true,
	true, false, false, false, false, false, false, false, true,
	true, true, true, true, true, true, true, true, true}
var Str_CSA [Square_NB + Piece_Can_Drop_NB]string = [Square_NB + Piece_Can_Drop_NB]string{"91", "81", "71", "61", "51", "41", "31", "21", "11",
	"92", "82", "72", "62", "52", "42", "32", "22", "12",
	"93", "83", "73", "63", "53", "43", "33", "23", "13",
	"94", "84", "74", "64", "54", "44", "34", "24", "14",
	"95", "85", "75", "65", "55", "45", "35", "25", "15",
	"96", "86", "76", "66", "56", "46", "36", "26", "16",
	"97", "87", "77", "67", "57", "47", "37", "27", "17",
	"98", "88", "78", "68", "58", "48", "38", "28", "18",
	"99", "89", "79", "69", "59", "49", "39", "29", "19",
	"00", "00", "00", "00", "00", "00", "00"}

var Str_USI [Square_NB + Piece_Can_Drop_NB]string = [Square_NB + Piece_Can_Drop_NB]string{"9a", "8a", "7a", "6a", "5a", "4a", "3a", "2a", "1a",
	"9b", "8b", "7b", "6b", "5b", "4b", "3b", "2b", "1b",
	"9c", "8c", "7c", "6c", "5c", "4c", "3c", "2c", "1c",
	"9d", "8d", "7d", "6d", "5d", "4d", "3d", "2d", "1d",
	"9e", "8e", "7e", "6e", "5e", "4e", "3e", "2e", "1e",
	"9f", "8f", "7f", "6f", "5f", "4f", "3f", "2f", "1f",
	"9g", "8g", "7g", "6g", "5g", "4g", "3g", "2g", "1g",
	"9h", "8h", "7h", "6h", "5h", "4h", "3h", "2h", "1h",
	"9i", "8i", "7i", "6i", "5i", "4i", "3i", "2i", "1i"}

var USI_TO_SQ = map[string]uint8{"9a": 0, "8a": 1, "7a": 2, "6a": 3, "5a": 4, "4a": 5, "3a": 6, "2a": 7, "1a": 8,
	"9b": 9, "8b": 10, "7b": 11, "6b": 12, "5b": 13, "4b": 14, "3b": 15, "2b": 16, "1b": 17,
	"9c": 18, "8c": 19, "7c": 20, "6c": 21, "5c": 22, "4c": 23, "3c": 24, "2c": 25, "1c": 26,
	"9d": 27, "8d": 28, "7d": 29, "6d": 30, "5d": 31, "4d": 32, "3d": 33, "2d": 34, "1d": 35,
	"9e": 36, "8e": 37, "7e": 38, "6e": 39, "5e": 40, "4e": 41, "3e": 42, "2e": 43, "1e": 44,
	"9f": 45, "8f": 46, "7f": 47, "6f": 48, "5f": 49, "4f": 50, "3f": 51, "2f": 52, "1f": 53,
	"9g": 54, "8g": 55, "7g": 56, "6g": 57, "5g": 58, "4g": 59, "3g": 60, "2g": 61, "1g": 62,
	"9h": 63, "8h": 64, "7h": 65, "6h": 66, "5h": 67, "4h": 68, "3h": 69, "2h": 70, "1h": 71,
	"9i": 72, "8i": 73, "7i": 74, "6i": 75, "5i": 76, "4i": 77, "3i": 78, "2i": 79, "1i": 80}

var Index_Bitboard [Square_NB]int8 = [Square_NB]int8{0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1,
	2, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2}

var Index_Rotated [Square_NB]int8 = [Square_NB]int8{8, 17, 26, 35, 44, 53, 62, 71, 80,
	7, 16, 25, 34, 43, 52, 61, 70, 79,
	6, 15, 24, 33, 42, 51, 60, 69, 78,
	5, 14, 23, 32, 41, 50, 59, 68, 77,
	4, 13, 22, 31, 40, 49, 58, 67, 76,
	3, 12, 21, 30, 39, 48, 57, 66, 75,
	2, 11, 20, 29, 38, 47, 56, 65, 74,
	1, 10, 19, 28, 37, 46, 55, 64, 73,
	0, 9, 18, 27, 36, 45, 54, 63, 72}

var Str_Piece [Piece_NB]string = [Piece_NB]string{"None", "FU", "KY", "KE", "GI", "KI", "KA", "HI", "OU", "TO", "NY", "NK", "NG", "None", "UM", "RY"}
var Str_Piece_JP [Piece_NB]string = [Piece_NB]string{"None", "歩", "香", "桂", "銀", "金", "角", "飛", "玉", "と", "杏", "圭", "全", "None", "馬", "龍"}
var Set_Empty_Num [10]string = [10]string{"", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var Set_Hand_Num [10]string = [10]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
var Usi_Drop_Piece = map[string]Piece{"P": Pawn, "L": Lance, "N": Knight, "S": Silver, "G": Gold, "B": Bishop, "R": Rook}

type Color uint8

const (
	Black Color = iota
	White
)

type Result uint8

const (
	BlackWin Result = iota
	WhiteWin
	Draw
)

type Piece int8

const (
	Empty Piece = iota
	Pawn
	Lance
	Knight
	Silver
	Gold
	Bishop
	Rook
	King
	Pro_Pawn
	Pro_Lance
	Pro_Knight
	Pro_Silver
	None
	Horse
	Dragon
)

var CSA_TO_SQ = map[string]uint8{"91": 0, "81": 1, "71": 2, "61": 3, "51": 4, "41": 5, "31": 6, "21": 7, "11": 8,
	"92": 9, "82": 10, "72": 11, "62": 12, "52": 13, "42": 14, "32": 15, "22": 16, "12": 17,
	"93": 18, "83": 19, "73": 20, "63": 21, "53": 22, "43": 23, "33": 24, "23": 25, "13": 26,
	"94": 27, "84": 28, "74": 29, "64": 30, "54": 31, "44": 32, "34": 33, "24": 34, "14": 35,
	"95": 36, "85": 37, "75": 38, "65": 39, "55": 40, "45": 41, "35": 42, "25": 43, "15": 44,
	"96": 45, "86": 46, "76": 47, "66": 48, "56": 49, "46": 50, "36": 51, "26": 52, "16": 53,
	"97": 54, "87": 55, "77": 56, "67": 57, "57": 58, "47": 59, "37": 60, "27": 61, "17": 62,
	"98": 63, "88": 64, "78": 65, "68": 66, "58": 67, "48": 68, "38": 69, "28": 70, "18": 71,
	"99": 72, "89": 73, "79": 74, "69": 75, "59": 76, "49": 77, "39": 78, "29": 79, "19": 80,
	"00": 81}

var CSA_TO_PC = map[string]Piece{"FU": Pawn, "KY": Lance, "KE": Knight, "GI": Silver, "KI": Gold, "KA": Bishop, "HI": Rook, "OU": King,
	"TO": Pro_Pawn, "NY": Pro_Lance, "NK": Pro_Knight, "NG": Pro_Silver, "None": None, "UM": Horse, "RY": Dragon}

var Str_SFEN_Pc = map[Piece]string{Pawn: "P", Lance: "L", Knight: "N", Silver: "S",
	Gold: "G", Bishop: "B", Rook: "R", King: "K",
	Pro_Pawn: "+P", Pro_Lance: "+L", Pro_Knight: "+N", Pro_Silver: "+S",
	Horse: "+B", Dragon: "+R", -Pawn: "p", -Lance: "l", -Knight: "n", -Silver: "s",
	-Gold: "g", -Bishop: "b", -Rook: "r", -King: "k",
	-Pro_Pawn: "+p", -Pro_Lance: "+l", -Pro_Knight: "+n", -Pro_Silver: "+s",
	-Horse: "+b", -Dragon: "+r"}

var Int_Pc = map[string]Piece{"P": Pawn, "L": Lance, "N": Knight, "S": Silver,
	"G": Gold, "B": Bishop, "R": Rook, "K": King,
	"p": -Pawn, "l": -Lance, "n": -Knight, "s": -Silver,
	"g": -Gold, "b": -Bishop, "r": -Rook, "k": -King}

var Int_Empty_Num = map[string]uint8{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9}

var Int_Hand_Num = map[string]uint8{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9}

var Str_Color = map[Color]string{Black: "b", White: "w"}

var Num_Color = map[string]Color{"b": Black, "w": White}

var Piece_Table [Color_NB]map[int][]int

var File_Shift_Table [Square_NB][Color_NB]uint

var Diag1_Shift_Table [Square_NB][Color_NB]uint

var Diag2_Shift_Table [Square_NB][Color_NB]uint

type Direction int8

const (
	Direc_Misc         = 0
	Direc_Diag1_U2d    = 8
	Direc_Diag1_D2u    = -8
	Direc_Diag2_U2d    = 10
	Direc_Diag2_D2u    = -10
	Direc_File_U2d     = 9
	Direc_File_D2u     = -9
	Direc_Rank_L2r     = 1
	Direc_Rank_R2l     = -1
	Direc_Knight_L_U2d = 19
	Direc_Knight_R_U2d = 17
	Direc_Knight_L_D2u = -17
	Direc_Knight_R_D2u = -19
)

type File uint8

const (
	File1 File = iota
	File2
	File3
	File4
	File5
	File6
	File7
	File8
	File9
)

type Rank uint8

const (
	Rank1 Rank = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
	Rank9
)

var FileTable [Square_NB]File = [Square_NB]File{File1, File2, File3, File4, File5, File6, File7, File8, File9,
	File1, File2, File3, File4, File5, File6, File7, File8, File9,
	File1, File2, File3, File4, File5, File6, File7, File8, File9,
	File1, File2, File3, File4, File5, File6, File7, File8, File9,
	File1, File2, File3, File4, File5, File6, File7, File8, File9,
	File1, File2, File3, File4, File5, File6, File7, File8, File9,
	File1, File2, File3, File4, File5, File6, File7, File8, File9,
	File1, File2, File3, File4, File5, File6, File7, File8, File9,
	File1, File2, File3, File4, File5, File6, File7, File8, File9}

var RankTable [Square_NB]Rank = [Square_NB]Rank{Rank1, Rank1, Rank1, Rank1, Rank1, Rank1, Rank1, Rank1, Rank1,
	Rank2, Rank2, Rank2, Rank2, Rank2, Rank2, Rank2, Rank2, Rank2,
	Rank3, Rank3, Rank3, Rank3, Rank3, Rank3, Rank3, Rank3, Rank3,
	Rank4, Rank4, Rank4, Rank4, Rank4, Rank4, Rank4, Rank4, Rank4,
	Rank5, Rank5, Rank5, Rank5, Rank5, Rank5, Rank5, Rank5, Rank5,
	Rank6, Rank6, Rank6, Rank6, Rank6, Rank6, Rank6, Rank6, Rank6,
	Rank7, Rank7, Rank7, Rank7, Rank7, Rank7, Rank7, Rank7, Rank7,
	Rank8, Rank8, Rank8, Rank8, Rank8, Rank8, Rank8, Rank8, Rank8,
	Rank9, Rank9, Rank9, Rank9, Rank9, Rank9, Rank9, Rank9, Rank9}

type Label []uint8

const (
	EMPTY = iota
	UP_LEFT
	UP
	UP_RIGHT
	LEFT
	RIGHT
	DOWN_LEFT
	DOWN
	DOWN_RIGHT
	UP_LEFT_KNIGHT
	UP_RIGHT_KNIGHT
	DOWN_LEFT_KNIGHT
	DOWN_RIGHT_KNIGHT
	UP_LEFT_PRO
	UP_PRO
	UP_RIGHT_PRO
	LEFT_PRO
	RIGHT_PRO
	DOWN_LEFT_PRO
	DOWN_PRO
	DOWN_RIGHT_PRO
	UP_LEFT_KNIGHT_PRO
	UP_RIGHT_KNIGHT_PRO
	DOWN_LEFT_KNIGHT_PRO
	DOWN_RIGHT_KNIGHT_PRO
	DROP_PAWN
	DROP_LANCE
	DROP_KNIGHT
	DROP_SILVER
	DROP_GOLD
	DROP_BISHOP
	DROP_ROOK
)

type BitBoard struct {
	P [3]uint32
}

var BB_Black_Position BitBoard
var BB_White_Position BitBoard
var BB_Color_Position [Color_NB]BitBoard
var BB_Rev_Color_Position [Color_NB]BitBoard
var BB_DMZ BitBoard
var BB_Full BitBoard
var BB_File [NFile]BitBoard
var BB_Rank [NRank]BitBoard
var BB_Knight_Must_Promote [Color_NB]BitBoard

// var BB_Lance_Mask [Color_NB][NRank]BitBoard
var BB_Pawn_Lance_Can_Drop [Color_NB]BitBoard
var BB_Knight_Can_Drop [Color_NB]BitBoard
var BB_Others_Can_Drop BitBoard
var BB_Pawn_Mask [Color_NB]BitBoard

const RandM = 397
const RandN = 624
const MaskU uint = 0x80000000
const MaskL uint = 0x7fffffff
const Mask32 uint = 0xffffffff

var Delta_Table [Color_NB]int = [Color_NB]int{-9, 9}
var Sign_Table [Color_NB]int = [Color_NB]int{-1, 1}
var Set_Long_Attack_Pieces [10]Piece = [10]Piece{Lance, Bishop, Rook, Horse, Dragon, -Lance, -Bishop, -Rook, -Horse, -Dragon}

var Set_Piece_Can_Promote0 [3]Piece = [3]Piece{Pawn, Lance, Knight}
var Set_Piece_Can_Promote1 [3]Piece = [3]Piece{Silver, Bishop, Rook}
var ShortPieces [2]Piece = [2]Piece{Pawn, Lance}
var LongPieces [2]Piece = [2]Piece{Bishop, Rook}
var LongPieces2 [3]Piece = [3]Piece{Bishop, Rook, Lance}
var Piece_RD [2]Piece = [2]Piece{Rook, Dragon}
var Piece_BH [2]Piece = [2]Piece{Bishop, Horse}
var Piece_RDBHL [5]Piece = [5]Piece{Rook, Dragon, Bishop, Horse, Lance}
var Adirec [Square_NB][Square_NB]int

type BoardTree struct {
	BB_Piece    [Color_NB][Piece_NB]BitBoard
	BB_Occupied [Color_NB]BitBoard
	//BB_Rotated  [Color_NB]BitBoard
	SQ_King     [Color_NB]uint8
	Hand        [Color_NB]uint32
	Board       [Square_NB]int8
	RootColor   uint8
	Hash        [Ply_Max]uint64
	CurrentHash uint64
	PrevHash    uint64
	Ply         uint16
	EvalArray   [Ply_Max]int32
}

type Attacks struct {
	ABB_Mask          [Square_NB]BitBoard
	ABB_Piece_Attacks [Color_NB][Piece_NB][Square_NB]BitBoard
	ABB_File_Mask_Ex  [Square_NB]BitBoard
	ABB_Rank_Mask_Ex  [Square_NB]BitBoard
	ABB_Diag1_Mask_Ex [Square_NB]BitBoard
	ABB_Diag2_Mask_Ex [Square_NB]BitBoard
	//ABB_Diagonal_Mask_Ex [Square_NB]BitBoard
	//ABB_Cross_Mask_Ex    [Square_NB]BitBoard
	ABB_Rank_Attacks  [Square_NB][Long_BB_Hash]BitBoard
	ABB_File_Attacks  [Square_NB][Long_BB_Hash]BitBoard
	ABB_Diag1_Attacks [Square_NB][Long_BB_Hash]BitBoard
	ABB_Diag2_Attacks [Square_NB][Long_BB_Hash]BitBoard
	//ABB_Cross_Attacks    [Square_NB][Long_BB_Hash]BitBoard
	ABB_Lance_Mask_Ex [Color_NB][Square_NB]BitBoard
	ABB_Obstacles     [Square_NB][Square_NB]BitBoard
	//ABB_Lance_Attacks [Color_NB][Square_NB][Long_BB_Hash]BitBoard
}

var Atk Attacks

type TT struct {
	value    map[uint64]int16
	color    map[uint64]Color
	is_check map[uint64]bool
	move     map[uint64]uint32
	ply      map[uint64]uint16
}

type MateSearchTree struct {
	move_cur         [Mate_Ply_Max]uint32
	mate_proc        [][]uint32
	no_mate_proc     [][]uint32
	first_move       uint32
	second_move      uint32
	max_ply          int
	is_abort         bool
	is_mate_root     bool
	bt               BoardTree
	root_check_moves []uint32
	str_pv           string
}

type Record struct {
	str_moves []string
	moves     []uint32
	ply       int
	result    string
}

func Init() {
	//fmt.Println(BB_Pawn_Mask)
	Piece_Table[0] = make(map[int][]int)
	Piece_Table[1] = make(map[int][]int)
	var li []int
	li = append(li, int(Silver))
	li = append(li, int(Gold))
	li = append(li, int(Bishop))
	li = append(li, int(Pro_Pawn))
	li = append(li, int(Pro_Lance))
	li = append(li, int(Pro_Knight))
	li = append(li, int(Pro_Silver))
	li = append(li, int(Horse))
	li = append(li, int(Dragon))
	Piece_Table[0][int(Direc_Diag1_U2d)] = li
	Piece_Table[0][int(Direc_Diag2_U2d)] = li
	Piece_Table[1][int(Direc_Diag1_D2u)] = li
	Piece_Table[1][int(Direc_Diag2_D2u)] = li
	var li2 []int
	li2 = append(li2, int(Pawn))
	li2 = append(li2, int(Lance))
	li2 = append(li2, int(Silver))
	li2 = append(li2, int(Gold))
	li2 = append(li2, int(Rook))
	li2 = append(li2, int(Pro_Pawn))
	li2 = append(li2, int(Pro_Lance))
	li2 = append(li2, int(Pro_Knight))
	li2 = append(li2, int(Pro_Silver))
	li2 = append(li2, int(Horse))
	li2 = append(li2, int(Dragon))
	Piece_Table[0][int(Direc_File_U2d)] = li2
	Piece_Table[1][int(Direc_File_D2u)] = li2
	var li3 []int
	li3 = append(li3, int(Gold))
	li3 = append(li3, int(Rook))
	li3 = append(li3, int(Pro_Pawn))
	li3 = append(li3, int(Pro_Lance))
	li3 = append(li3, int(Pro_Knight))
	li3 = append(li3, int(Pro_Silver))
	li3 = append(li3, int(Horse))
	li3 = append(li3, int(Dragon))
	Piece_Table[0][int(Direc_Rank_L2r)] = li3
	Piece_Table[0][int(Direc_Rank_R2l)] = li3
	Piece_Table[0][int(Direc_File_D2u)] = li3
	Piece_Table[1][int(Direc_Rank_L2r)] = li3
	Piece_Table[1][int(Direc_Rank_R2l)] = li3
	Piece_Table[1][int(Direc_File_U2d)] = li3
	var li4 []int
	li4 = append(li4, int(Silver))
	li4 = append(li4, int(Bishop))
	li4 = append(li4, int(Horse))
	li4 = append(li4, int(Dragon))
	Piece_Table[0][int(Direc_Diag1_D2u)] = li4
	Piece_Table[0][int(Direc_Diag2_D2u)] = li4
	Piece_Table[1][int(Direc_Diag1_U2d)] = li4
	Piece_Table[1][int(Direc_Diag2_U2d)] = li4
	var v uint32 = 1
	for i := 80; i >= 54; i-- {
		Atk.ABB_Mask[i].P[2] = v
		v = v << 1
	}
	v = 1
	for i := 53; i >= 27; i-- {
		Atk.ABB_Mask[i].P[1] = v
		v = v << 1
	}
	v = 1
	for i := 26; i >= 0; i-- {
		Atk.ABB_Mask[i].P[0] = v
		v = v << 1
	}

	fp, err := os.Open("tables/abb_file_mask_ex.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			d, _ := strconv.Atoi(line[3])
			//fmt.Println((temp_big))
			Atk.ABB_File_Mask_Ex[uint(a)].P[0] = uint32(b)
			Atk.ABB_File_Mask_Ex[uint(a)].P[1] = uint32(c)
			Atk.ABB_File_Mask_Ex[uint(a)].P[2] = uint32(d)
		}
	}

	fp, err = os.Open("tables/abb_file_attacks.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader = csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			d, _ := strconv.Atoi(line[3])
			e, _ := strconv.Atoi(line[4])
			//fmt.Println((temp_big))
			Atk.ABB_File_Attacks[uint(a)][uint(b)].P[0] = uint32(c)
			Atk.ABB_File_Attacks[uint(a)][uint(b)].P[1] = uint32(d)
			Atk.ABB_File_Attacks[uint(a)][uint(b)].P[2] = uint32(e)
		}
	}

	fp, err = os.Open("tables/abb_rank_mask_ex.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader = csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			d, _ := strconv.Atoi(line[3])
			//fmt.Println((temp_big))
			Atk.ABB_Rank_Mask_Ex[uint(a)].P[0] = uint32(b)
			Atk.ABB_Rank_Mask_Ex[uint(a)].P[1] = uint32(c)
			Atk.ABB_Rank_Mask_Ex[uint(a)].P[2] = uint32(d)
		}
	}

	fp, err = os.Open("tables/abb_rank_attacks.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader = csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			d, _ := strconv.Atoi(line[3])
			e, _ := strconv.Atoi(line[4])
			//fmt.Println((temp_big))
			Atk.ABB_Rank_Attacks[uint(a)][uint(b)].P[0] = uint32(c)
			Atk.ABB_Rank_Attacks[uint(a)][uint(b)].P[1] = uint32(d)
			Atk.ABB_Rank_Attacks[uint(a)][uint(b)].P[2] = uint32(e)
		}
	}

	fp, err = os.Open("tables/abb_diag1_mask_ex.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader = csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			d, _ := strconv.Atoi(line[3])
			//fmt.Println((temp_big))
			Atk.ABB_Diag1_Mask_Ex[uint(a)].P[0] = uint32(b)
			Atk.ABB_Diag1_Mask_Ex[uint(a)].P[1] = uint32(c)
			Atk.ABB_Diag1_Mask_Ex[uint(a)].P[2] = uint32(d)
		}
	}

	fp, err = os.Open("tables/abb_diag1_attacks.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader = csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			d, _ := strconv.Atoi(line[3])
			e, _ := strconv.Atoi(line[4])
			//fmt.Println((temp_big))
			Atk.ABB_Diag1_Attacks[uint(a)][uint(b)].P[0] = uint32(c)
			Atk.ABB_Diag1_Attacks[uint(a)][uint(b)].P[1] = uint32(d)
			Atk.ABB_Diag1_Attacks[uint(a)][uint(b)].P[2] = uint32(e)
		}
	}

	fp, err = os.Open("tables/abb_diag2_mask_ex.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader = csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			d, _ := strconv.Atoi(line[3])
			//fmt.Println((temp_big))
			Atk.ABB_Diag2_Mask_Ex[uint(a)].P[0] = uint32(b)
			Atk.ABB_Diag2_Mask_Ex[uint(a)].P[1] = uint32(c)
			Atk.ABB_Diag2_Mask_Ex[uint(a)].P[2] = uint32(d)
		}
	}

	fp, err = os.Open("tables/abb_diag2_attacks.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader = csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			d, _ := strconv.Atoi(line[3])
			e, _ := strconv.Atoi(line[4])
			//fmt.Println((temp_big))
			Atk.ABB_Diag2_Attacks[uint(a)][uint(b)].P[0] = uint32(c)
			Atk.ABB_Diag2_Attacks[uint(a)][uint(b)].P[1] = uint32(d)
			Atk.ABB_Diag2_Attacks[uint(a)][uint(b)].P[2] = uint32(e)
		}
	}

	fp, err = os.Open("tables/abb_minus_rays.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader = csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			d, _ := strconv.Atoi(line[3])
			//fmt.Println((temp_big))
			Atk.ABB_Lance_Mask_Ex[0][uint(a)].P[0] = uint32(b)
			Atk.ABB_Lance_Mask_Ex[0][uint(a)].P[1] = uint32(c)
			Atk.ABB_Lance_Mask_Ex[0][uint(a)].P[2] = uint32(d)
		}
	}

	fp, err = os.Open("tables/abb_plus_rays.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader = csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			d, _ := strconv.Atoi(line[3])
			//fmt.Println((temp_big))
			Atk.ABB_Lance_Mask_Ex[1][uint(a)].P[0] = uint32(b)
			Atk.ABB_Lance_Mask_Ex[1][uint(a)].P[1] = uint32(c)
			Atk.ABB_Lance_Mask_Ex[1][uint(a)].P[2] = uint32(d)
		}
	}

	fp, err = os.Open("tables/abb_obstacle.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader = csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			d, _ := strconv.Atoi(line[3])
			e, _ := strconv.Atoi(line[4])
			//fmt.Println((temp_big))
			Atk.ABB_Obstacles[uint(a)][uint(b)].P[0] = uint32(c)
			Atk.ABB_Obstacles[uint(a)][uint(b)].P[1] = uint32(d)
			Atk.ABB_Obstacles[uint(a)][uint(b)].P[2] = uint32(e)
		}
	}

	for c := Black; c < Color_NB; c++ {
		for pc := Pawn; pc <= Dragon; pc++ {
			for sq := 0; sq < Square_NB; sq++ {
				switch pc {
				case Pawn:
					var pos = sq + Sign_Table[c]*NRank
					if pos >= 0 && pos < Square_NB {
						Atk.ABB_Piece_Attacks[c][pc][sq] = Atk.ABB_Mask[pos]
					}
				case Knight:
					if c == Black {
						if sq >= 18 {
							var pos = sq + Direc_Knight_L_D2u
							if pos >= 0 && pos < Square_NB {
								if RankTable[sq] == (RankTable[pos] + 2) {
									Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
								}
							}
							pos = sq + Direc_Knight_R_D2u
							if pos >= 0 && pos < Square_NB {
								if RankTable[sq] == (RankTable[pos] + 2) {
									Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
								}
							}
						}
					} else {
						if sq <= 62 {
							var pos = sq + Direc_Knight_L_U2d
							if pos >= 0 && pos < Square_NB {
								if RankTable[sq] == (RankTable[pos] - 2) {
									Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
								}
							}
							pos = sq + Direc_Knight_R_U2d
							if pos >= 0 && pos < Square_NB {
								if RankTable[sq] == (RankTable[pos] - 2) {
									Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
								}
							}
						}
					}
				case Silver:
					var pos = sq + Direc_Diag1_D2u
					if pos >= 0 && pos < Square_NB && RankTable[sq] != RankTable[pos] {
						Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
					}
					pos = sq + Direc_Diag1_U2d
					if pos >= 0 && pos < Square_NB && RankTable[sq] != RankTable[pos] {
						Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
					}
					pos = sq + Direc_Diag2_D2u
					if pos >= 0 && pos < Square_NB {
						if RankTable[pos] == RankTable[sq]-1 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
					pos = sq + Direc_Diag2_U2d
					if pos >= 0 && pos < Square_NB {
						if RankTable[pos] == RankTable[sq]+1 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
					pos = sq + Sign_Table[c]*NRank
					if pos >= 0 && pos < Square_NB {
						Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
					}
				case Gold, Pro_Pawn, Pro_Lance, Pro_Knight, Pro_Silver:
					var pos = sq + Direc_Rank_L2r
					if pos >= 0 && pos < Square_NB {
						if FileTable[sq] != File9 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
					pos = sq + Direc_Rank_R2l
					if pos >= 0 && pos < Square_NB {
						if FileTable[sq] != File1 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
					pos = sq + Direc_File_D2u
					if pos >= 0 && pos < Square_NB {
						if RankTable[sq] != Rank1 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
					pos = sq + Direc_File_U2d
					if pos >= 0 && pos < Square_NB {
						if RankTable[sq] != Rank9 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
					if c == Black {
						pos = sq + Direc_Diag1_D2u
						if pos >= 0 && pos < Square_NB && RankTable[sq] != RankTable[pos] {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
						pos = sq + Direc_Diag2_D2u
						if pos >= 0 && pos < Square_NB && RankTable[pos] == RankTable[sq]-1 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					} else {
						pos = sq + Direc_Diag1_U2d
						if pos >= 0 && pos < Square_NB && RankTable[sq] != RankTable[pos] {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
						pos = sq + Direc_Diag2_U2d
						if pos >= 0 && pos < Square_NB && RankTable[pos] == RankTable[sq]+1 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
				case King:
					var pos = sq + Direc_Rank_L2r
					if pos >= 0 && pos < Square_NB {
						if FileTable[sq] != File9 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
					pos = sq + Direc_Rank_R2l
					if pos >= 0 && pos < Square_NB {
						if FileTable[sq] != File1 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
					pos = sq + Direc_File_D2u
					if pos >= 0 && pos < Square_NB {
						if RankTable[sq] != Rank1 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
					pos = sq + Direc_File_U2d
					if pos >= 0 && pos < Square_NB {
						if RankTable[sq] != Rank9 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
					pos = sq + Direc_Diag1_D2u
					if pos >= 0 && pos < Square_NB && RankTable[sq] != RankTable[pos] {
						Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
					}
					pos = sq + Direc_Diag1_U2d
					if pos >= 0 && pos < Square_NB && RankTable[sq] != RankTable[pos] {
						Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
					}
					pos = sq + Direc_Diag2_D2u
					if pos >= 0 && pos < Square_NB {
						if RankTable[pos] == RankTable[sq]-1 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
					pos = sq + Direc_Diag2_U2d
					if pos >= 0 && pos < Square_NB {
						if RankTable[pos] == RankTable[sq]+1 {
							Atk.ABB_Piece_Attacks[c][pc][sq] = BBOr(Atk.ABB_Piece_Attacks[c][pc][sq], Atk.ABB_Mask[pos])
						}
					}
				}
			}
		}
	}

	fp, err = os.Open("tables/adirec.txt")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader = csv.NewReader(fp)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			a, _ := strconv.Atoi(line[0])
			b, _ := strconv.Atoi(line[1])
			c, _ := strconv.Atoi(line[2])
			//fmt.Println((temp_big))
			Adirec[uint(a)][uint(b)] = c
		}
	}

	BB_Black_Position = BBIni()
	BB_Black_Position.P[2] = 134217727
	BB_White_Position = BBIni()
	BB_White_Position.P[0] = 134217727
	BB_Color_Position[0] = BB_Black_Position
	BB_Color_Position[1] = BB_White_Position
	BB_Rev_Color_Position[0] = BBIni()
	BB_Rev_Color_Position[1] = BBIni()
	BB_Rev_Color_Position[0] = BB_White_Position
	BB_Rev_Color_Position[1] = BB_Black_Position
	BB_DMZ = BBIni()
	BB_DMZ.P[1] = 134217727
	BB_Full = BBIni()
	BB_Full.P[0] = 134217727
	BB_Full.P[1] = 134217727
	BB_Full.P[2] = 134217727
	for i := 0; i < 9; i++ {
		BB_File[i] = BBIni()
		BB_Rank[i] = BBIni()
	}
	BB_File[0].P[0] = 67240192
	BB_File[0].P[1] = 67240192
	BB_File[0].P[2] = 67240192
	BB_File[1].P[0] = 33620096
	BB_File[1].P[1] = 33620096
	BB_File[1].P[2] = 33620096
	BB_File[2].P[0] = 16810048
	BB_File[2].P[1] = 16810048
	BB_File[2].P[2] = 16810048
	BB_File[3].P[0] = 8405024
	BB_File[3].P[1] = 8405024
	BB_File[3].P[2] = 8405024
	BB_File[4].P[0] = 4202512
	BB_File[4].P[1] = 4202512
	BB_File[4].P[2] = 4202512
	BB_File[5].P[0] = 2101256
	BB_File[5].P[1] = 2101256
	BB_File[5].P[2] = 2101256
	BB_File[6].P[0] = 1050628
	BB_File[6].P[1] = 1050628
	BB_File[6].P[2] = 1050628
	BB_File[7].P[0] = 525314
	BB_File[7].P[1] = 525314
	BB_File[7].P[2] = 525314
	BB_File[8].P[0] = 262657
	BB_File[8].P[1] = 262657
	BB_File[8].P[2] = 262657
	BB_Rank[0].P[0] = 133955584
	BB_Rank[1].P[0] = 261632
	BB_Rank[2].P[0] = 511
	BB_Rank[3].P[1] = 133955584
	BB_Rank[4].P[1] = 261632
	BB_Rank[5].P[1] = 511
	BB_Rank[6].P[2] = 133955584
	BB_Rank[7].P[2] = 261632
	BB_Rank[8].P[2] = 511
	BB_Knight_Must_Promote[0] = BBIni()
	BB_Knight_Must_Promote[1] = BBIni()
	BB_Knight_Must_Promote[0].P[0] = 134217216
	BB_Knight_Must_Promote[1].P[2] = 262143
	BB_Pawn_Lance_Can_Drop[0] = BBXor(BB_Full, BB_Rank[0])
	BB_Pawn_Lance_Can_Drop[1] = BBXor(BB_Full, BB_Rank[8])
	var bb_temp = BBOr(BB_Rank[0], BB_Rank[1])
	BB_Knight_Can_Drop[0] = BBXor(BB_Full, bb_temp)
	bb_temp = BBOr(BB_Rank[7], BB_Rank[8])
	BB_Knight_Can_Drop[1] = BBXor(BB_Full, bb_temp)
	BB_Others_Can_Drop = BB_Full
	bb_temp = BBOr(BB_Rank[4], BB_Rank[5])
	bb_temp = BBOr(bb_temp, BB_Rank[6])
	bb_temp = BBOr(bb_temp, BB_Rank[7])
	bb_temp = BBOr(bb_temp, BB_Rank[8])
	BB_Pawn_Mask[0] = bb_temp
	bb_temp = BBOr(BB_Rank[0], BB_Rank[1])
	bb_temp = BBOr(bb_temp, BB_Rank[2])
	bb_temp = BBOr(bb_temp, BB_Rank[3])
	bb_temp = BBOr(bb_temp, BB_Rank[4])
	BB_Pawn_Mask[1] = bb_temp

	var shift [2]int
	for i := 0; i < Square_NB; i++ {
		shift[1] = bits.OnesCount32(Atk.ABB_Diag1_Mask_Ex[i].P[2])
		shift[0] = shift[1] + bits.OnesCount32(Atk.ABB_Diag1_Mask_Ex[i].P[1])
		for j := 0; j < Color_NB; j++ {
			Diag1_Shift_Table[i][j] = uint(shift[j])
		}
		shift[1] = bits.OnesCount32(Atk.ABB_Diag2_Mask_Ex[i].P[2])
		shift[0] = shift[1] + bits.OnesCount32(Atk.ABB_Diag2_Mask_Ex[i].P[1])
		for j := 0; j < Color_NB; j++ {
			Diag2_Shift_Table[i][j] = uint(shift[j])
		}
		shift[1] = bits.OnesCount32(Atk.ABB_File_Mask_Ex[i].P[2])
		shift[0] = shift[1] + bits.OnesCount32(Atk.ABB_File_Mask_Ex[i].P[1])
		for j := 0; j < Color_NB; j++ {
			File_Shift_Table[i][j] = uint(shift[j])
		}
	}
}
