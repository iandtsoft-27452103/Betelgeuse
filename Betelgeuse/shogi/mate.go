package shogi

import (
	"fmt"
)

func GenRootMoves(bt BoardTree) []uint32 {
	var checkMoves []uint32
	checkMoves = GenCheck(&bt, int(bt.RootColor), checkMoves)
	return checkMoves
}

func MateSearchWrapper(mst *MateSearchTree, depth_max int) string {
	var rest_depth = 1
	for rest_depth <= depth_max {
		mst.max_ply = rest_depth
		for i := 0; i < len(mst.move_cur); i++ {
			mst.move_cur[i] = 0
		}
		mst.mate_proc = make([][]uint32, 0)
		mst.no_mate_proc = make([][]uint32, 0)
		mst.first_move = 0
		mst.second_move = 0

		mst.is_mate_root = Offend(mst, int(mst.bt.RootColor), rest_depth, 1)

		if mst.is_mate_root {
			break
		}

		if mst.is_abort {
			mst.is_mate_root = false
			break
		}
		rest_depth += 2
	}
	if mst.is_mate_root && !mst.is_abort {
		mst.str_pv = OutResult(mst, rest_depth)
	}
	return mst.str_pv
}

func OutResult(mst *MateSearchTree, rest_depth int) string {
	var l = mst.mate_proc
	var nl = mst.no_mate_proc
	var b bool
	var str_pv = ""
	var str_color []string = []string{"+", "-"}
	b = false
	//var idxes map[int]int = map[int]int{}
	var nums []int
	var cnt = 0
	for i := 0; i < len(l); i++ {
		var s = fmt.Sprintf("%d", i+1) + " / " + fmt.Sprintf("%d", len(l))
		fmt.Println(s)
		for j := 0; j < len(nl); j++ {
			b = false
			for k := 0; k < len(nl[j])-1; k++ {
				if l[i][k] != nl[j][k] {
					b = true
					break
				}
			}
			if !b {
				//idxes[cnt] = i
				nums = append(nums, i)
				cnt += 1
			}
		}
	}
	for i := 0; i < len(l); i++ {
		// 問題点：マップにしたかったがうまく判定できていない。
		/*_, exists := idxes[i]
		if exists {
			continue
		}*/
		var b = false
		for j := 0; j < len(nums); j++ {
			if nums[j] == i {
				b = true
				break
			}
		}
		if b {
			continue
		}
		str_pv = ""
		var color = mst.bt.RootColor
		for j := 0; j < rest_depth; j++ {
			str_pv += str_color[color]
			str_pv += Move2CSA(l[i][j])
			if j != rest_depth-1 {
				str_pv += ", "
			}
			color ^= 1
		}
		fmt.Println(str_pv)
	}
	return str_pv
}

func Offend(mst *MateSearchTree, color int, rest_depth int, ply int) bool {
	var is_mate bool
	var checkMoves []uint32
	if mst.is_abort {
		return false
	}
	// generate check moves
	if ply != 1 {
		checkMoves = GenCheck(&mst.bt, color, checkMoves)
	} else {
		checkMoves = mst.root_check_moves
	}
	// If there are no check moves, it's not mate.
	if len(checkMoves) == 0 {
		return false
	}
	for i := 0; i < len(checkMoves); i++ {
		if mst.is_abort {
			return false
		}
		mst.move_cur[ply] = checkMoves[i]

		if ply == 3 && mst.move_cur[1] == 174746 && mst.move_cur[2] == 4291562639 && rest_depth >= 1 {
			fmt.Println("aa")
		}

		if PieceType(mst.move_cur[ply]) == uint32(Empty) {
			continue
		}
		Do(&mst.bt, mst.move_cur[ply], color)

		if mst.bt.SQ_King[0] != 0 {
			fmt.Println("aa")
		}

		// Sometimes, this problem occurs. I don't know what this causes.
		if BBTest(IsAttacked(mst.bt, int(mst.bt.SQ_King[color^1]), color^1)) == 0 {
			UnDo(&mst.bt, mst.move_cur[ply], color)
			continue
		}
		// case of discovered check
		if BBTest(IsAttacked(mst.bt, int(mst.bt.SQ_King[color]), color)) != 0 {
			UnDo(&mst.bt, mst.move_cur[ply], color)
			continue
		}
		is_mate = Defend(mst, color^1, rest_depth-1, ply+1)
		if is_mate {
			if ply == mst.max_ply {
				var moves []uint32 = make([]uint32, 0)
				for j := 1; j < ply+1; j++ {
					moves = append(moves, mst.move_cur[j])
				}
				mst.mate_proc = append(mst.mate_proc, moves)
			}
			UnDo(&mst.bt, mst.move_cur[ply], color)
			return true
		}
		UnDo(&mst.bt, mst.move_cur[ply], color)
	}
	return false
}

func Defend(mst *MateSearchTree, color int, rest_depth int, ply int) bool {
	var is_mate bool
	var mate_count int
	var evasionMoves []uint32 = make([]uint32, 0)

	if ply == 2 && mst.move_cur[1] == 174746 && rest_depth >= 1 {
		fmt.Println("aa")
	}

	if mst.is_abort {
		return false
	}
	evasionMoves = GenEvasion(&mst.bt, color, evasionMoves)
	// If rest depth equals to zero and there exists evasion moves, it isn't mate.
	if rest_depth == 0 && len(evasionMoves) > 0 {
		return false
	}
	// If there are no evasion moves, it's mate.
	if len(evasionMoves) == 0 {
		return true
	}
	for i := 0; i < len(evasionMoves); i++ {
		if mst.is_abort {
			return false
		}
		mst.move_cur[ply] = evasionMoves[i]

		if ply == 2 && mst.move_cur[1] == 174746 && mst.move_cur[2] == 4291562639 && rest_depth >= 0 {
			fmt.Println("aa")
		}

		Do(&mst.bt, mst.move_cur[ply], color)

		if mst.bt.SQ_King[0] != 0 {
			fmt.Println("aa")
		}

		// case of discovered check
		if BBTest(IsAttacked(mst.bt, int(mst.bt.SQ_King[color]), color)) != 0 {
			UnDo(&mst.bt, mst.move_cur[ply], color)
			continue
		}
		is_mate = Offend(mst, color^1, rest_depth-1, ply+1)
		if !is_mate {
			var moves []uint32 = make([]uint32, 0)
			for j := 1; j < ply+1; j++ {
				moves = append(moves, mst.move_cur[j])
			}
			mst.no_mate_proc = append(mst.no_mate_proc, moves)
			// If this is a no mate node, previous mate nodes are invalid.
			UnDo(&mst.bt, mst.move_cur[ply], color)
			return false
		} else {
			mate_count += 1
		}
		UnDo(&mst.bt, mst.move_cur[ply], color)
	}
	return true // All nodes are mate.
}
