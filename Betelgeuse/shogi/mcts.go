package shogi

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Node struct {
	color        uint8
	ParentIndex  uint64
	ThisIndex    uint64
	TrialCount   int64
	PlayoutCount uint64
	WinCount     uint64
	DrawCount    uint64
	LostCount    uint64
	EvalCount    uint64
	WinRateSum   float32
	LostRateSum  float32
	IsLeaf       bool
	ChildIndexes []uint64
	move         uint32
	PolicyResult float32
}

type MCTSTree struct {
	BTree    BoardTree
	NodeList []Node
	//PlayoutCount    uint64
	RootOutput      []float32
	ValueLambda     float32
	Nthr            uint16
	TaskNumber      uint16
	TimeBuffer      int64
	SearchTimeLimit int64
	/*queue_from_main_thread_p []string
	queue_to_main_thread_p   []string
	queue_from_main_thread_v []string
	queue_to_main_thread_v   []string*/
	queue_from_main_thread_p chan string
	queue_to_main_thread_p   chan string
	queue_from_main_thread_v chan string
	queue_to_main_thread_v   chan string
	tt                       TT2
	is_abort                 bool
	is_finished              bool
	start_time               time.Time
}

// This function is created by Copilot and I changed in details.
func MaxInSlicei64(nums []int64) (int64, int64) {
	maxVal := nums[0]
	var index int64 = 0
	for _, v := range nums[1:] {
		if v >= maxVal {
			maxVal = v
			index += 1
		}
	}
	return maxVal, index
}

// This function is created by Copilot and I changed in details.
func MaxInSlicef32(nums []float32) (float32, uint64) {
	maxVal := nums[0]
	var index uint64 = 0
	for _, v := range nums[1:] {
		if v > maxVal {
			maxVal = v
			index += 1
		}
	}
	return maxVal, index
}

func TotalParam(num_tasks int, mcts_tree []*MCTSTree) ([]*MCTSTree, []uint32, []float32, []int64) {
	var limit = 3
	BTree := mcts_tree[0].BTree
	num_root_moves := len(BTree.RootMoves)
	if num_root_moves == 1 {
		limit = 1
	} else if num_root_moves == 2 {
		limit = 2
	}
	if num_tasks > 1 {
		for i := 1; i < num_tasks; i++ {
			temp_tree := mcts_tree[i]
			for j := 0; j < num_root_moves; j++ {
				mcts_tree[0].NodeList[j+1].TrialCount += temp_tree.NodeList[j+1].TrialCount
				mcts_tree[0].NodeList[j+1].WinCount += temp_tree.NodeList[j+1].WinCount
				mcts_tree[0].NodeList[j+1].DrawCount += temp_tree.NodeList[j+1].DrawCount
				mcts_tree[0].NodeList[j+1].LostCount += temp_tree.NodeList[j+1].LostCount
				mcts_tree[0].NodeList[j+1].EvalCount += temp_tree.NodeList[j+1].EvalCount
				mcts_tree[0].NodeList[j+1].WinRateSum += temp_tree.NodeList[j+1].WinRateSum
				mcts_tree[0].NodeList[j+1].LostRateSum += temp_tree.NodeList[j+1].LostRateSum
			}
		}
	}
	var trial_count_array []int64
	for i := 0; i < num_root_moves; i++ {
		trial_count_array = append(trial_count_array, mcts_tree[0].NodeList[i+1].TrialCount)
	}
	var win_rate_array []float32
	for i := 0; i < num_root_moves; i++ {
		if trial_count_array[i] == 0 {
			win_rate_array = append(win_rate_array, 0)
		} else {
			win_rate_array = append(win_rate_array, float32(mcts_tree[0].NodeList[i+1].WinCount)/float32(trial_count_array[i]))
		}
	}
	var max_node_index int64 = 0
	var return_trial_count_array []int64
	var return_moves []uint32
	var return_win_rate_array []float32
	for i := 0; i < limit; i++ {
		_, max_index := MaxInSlicei64(trial_count_array)
		if i == 0 {
			max_node_index = max_index + 1
		}
		return_trial_count_array = append(return_trial_count_array, trial_count_array[max_index])
		trial_count_array[max_index] = math.MinInt64
		return_moves = append(return_moves, BTree.RootMoves[max_index])
		return_win_rate_array = append(return_win_rate_array, win_rate_array[max_index])
	}
	//Move best_move = return_moves[0];
	best_move_node := mcts_tree[0].NodeList[max_node_index]
	var limit2 = len(best_move_node.ChildIndexes)
	var depth = 1
	if limit2 > 0 {
		var temp_trial_count_array []int64
		for i := 0; i < limit2; i++ {
			temp_trial_count_array = append(temp_trial_count_array, 0)
		}
		for i := 0; i < limit2; i++ {
			var index = best_move_node.ChildIndexes[i]
			var current_node = mcts_tree[0].NodeList[index]
			temp_trial_count_array[i] = current_node.TrialCount
		}
		_, mi := MaxInSlicei64(temp_trial_count_array)
		mi += 1
		var node = mcts_tree[0].NodeList[mi]
		for {
			var limit3 = len(node.ChildIndexes)
			if limit3 == 0 {
				break
			}
			for i := 0; i < limit2; i++ {
				temp_trial_count_array[i] = 0
			}
			for i := 0; i < limit3; i++ {
				var index = node.ChildIndexes[i]
				var current_node = mcts_tree[0].NodeList[index]
				temp_trial_count_array[i] = current_node.TrialCount
			}
			_, mi := MaxInSlicei64(temp_trial_count_array)
			node = mcts_tree[0].NodeList[node.ChildIndexes[mi]]
			depth += 1
		}
	}

	return mcts_tree, return_moves, return_win_rate_array, return_trial_count_array
}

func GenRootMoves2(bt *BoardTree) []uint32 {
	var moves []uint32
	var legal_moves []uint32
	if BBTest(IsAttacked(*bt, int(bt.SQ_King[bt.RootColor]), int(bt.RootColor))) == uint32(0) {
		moves = GenCap(*bt, int(bt.RootColor), moves)
		moves = GenNoCap(*bt, int(bt.RootColor), moves)
		moves = GenDrop(bt, int(bt.RootColor), moves)
	} else {
		moves = GenEvasion(bt, int(bt.RootColor), moves)
	}
	var limit = len(moves)
	for i := 0; i < limit; i++ {
		var m = moves[i]
		if IsMoveValid(bt, m, int(bt.RootColor)) {
			legal_moves = append(legal_moves, m)
		}
	}
	return legal_moves
}

func SetRootOutput(cm *Communicate, mt *MCTSTree, root_moves []uint32) []uint32 {
	var str_msg = "p," + ToSFEN(mt.BTree, int(mt.BTree.RootColor))
	ThrowRequest(cm, str_msg)
	var str_response = ReceiveResponse(cm)
	var s = strings.Split(str_response, ",")
	var limit = len(s)
	// 探索中では8手返しているので、それに合わせるべきか？
	if limit > 8 {
		limit = 8
	}
	var moves []uint32
	for i := 0; i < limit; i++ {
		var ss = strings.Split(s[i], " ")
		var move = CSA2Move(mt.BTree, ss[0])
		f, _ := strconv.ParseFloat(ss[1], 64)
		moves = append(moves, move)
		mt.RootOutput = append(mt.RootOutput, float32(f))
	}
	return moves
}

func IsMoveValid(bt *BoardTree, move uint32, color int) bool {
	var ifrom = From(move)
	var ito = To(move)
	var cap_pc = CapPiece(move)
	// capture king
	if cap_pc == uint32(King) {
		return false
	}
	Do(bt, move, color)
	// discovered check
	if BBTest(IsAttacked(*bt, int(bt.SQ_King[color]), color)) > uint32(0) {
		UnDo(bt, move, color)
		return false
	}
	if ifrom < Square_NB {
		var idirec = Adirec[ifrom][ito]
		// moved pinned piece
		if BBTest(IsPinnedOnKing(bt, int(ifrom), idirec, color)) > uint32(0) {
			UnDo(bt, move, color)
			return false
		}
	} else {
		// drop move, but capture move
		if cap_pc > EMPTY {
			UnDo(bt, move, color)
			return false
		}
	}
	UnDo(bt, move, color)
	return true
}

func InitNode() Node {
	var n Node
	n.color = 0
	n.ParentIndex = 0
	n.ThisIndex = 0
	n.TrialCount = 0
	n.PlayoutCount = 0
	n.WinCount = 0
	n.DrawCount = 0
	n.LostCount = 0
	n.EvalCount = 0
	n.WinRateSum = 0.0
	n.LostRateSum = 0.0
	n.IsLeaf = true
	n.ChildIndexes = []uint64{}
	n.move = 0
	n.PolicyResult = 0.0
	return n
}

func Root(mt *MCTSTree) {
	var result_array []float32
	mt.NodeList = mt.NodeList[:0]
	var root_node Node
	mt.NodeList = append(mt.NodeList, root_node)
	var limit = len(mt.BTree.RootMoves)
	for i := 0; i < limit; i++ {
		var n Node = InitNode()
		n.color = mt.BTree.RootColor
		var m = mt.BTree.RootMoves[i]
		n.ParentIndex = 0
		n.ThisIndex = uint64(i) + 1
		n.move = m
		mt.NodeList[0].ChildIndexes = append(mt.NodeList[0].ChildIndexes, uint64(i)+1)
		n.PolicyResult = mt.RootOutput[i]
		result_array = append(result_array, mt.RootOutput[i])
		mt.NodeList = append(mt.NodeList, n)
	}

	for {
		elapsed := time.Since(mt.start_time)
		if elapsed.Milliseconds() > mt.SearchTimeLimit || mt.is_abort {
			break
		}
		var ucb1_array []float32
		var t uint64 = 0
		var i uint64 = 1
		var limit2 = uint64(limit) + 1
		for {
			if i >= limit2 {
				break
			}
			t += mt.NodeList[i].PlayoutCount
			i += 1
		}
		i = 1
		for {
			if i >= limit2 {
				break
			}
			var u float32 = mt.NodeList[i].PolicyResult * float32(math.Sqrt(float64(t))) / float32(mt.NodeList[i].PlayoutCount+1)
			var q float32 = 0
			if mt.NodeList[i].EvalCount > 0 && mt.NodeList[i].PlayoutCount > 0 {
				q = float32(((1-mt.ValueLambda)*(float32(mt.NodeList[i].WinRateSum)/float32(mt.NodeList[i].EvalCount)) + mt.ValueLambda*float32(float32(mt.NodeList[i].WinCount)/float32(mt.NodeList[i].PlayoutCount))))
			}
			ucb1_array = append(ucb1_array, u+q)
			i += 1
		}

		_, max_index := MaxInSlicef32(ucb1_array)
		max_index += 1
		Do(&mt.BTree, mt.BTree.RootMoves[max_index-1], int(mt.BTree.RootColor))
		if mt.NodeList[max_index].IsLeaf {
			elapsed = time.Since(mt.start_time)
			if (elapsed.Milliseconds()+mt.TimeBuffer) > mt.SearchTimeLimit || mt.is_abort {
				break
			}
			if mt.NodeList[max_index].TrialCount >= int64(mt.Nthr) {
				ExpandNode(mt, int(mt.BTree.RootColor^1), max_index, int(mt.BTree.Ply+1))
				if mt.is_abort {
					break
				}
			} else {
				var bt BoardTree
				bt = DeepCopy(mt.BTree)
				var result = PlayOut(&bt, int(bt.RootColor^1), int(mt.BTree.Ply)) // 引数plyの+1を取った。
				mt.NodeList[max_index].PlayoutCount++
				EvalNode(mt, max_index, int(bt.RootColor^1))
				if mt.is_abort {
					break
				}
				UpdateParam(mt, max_index, result)
			}
		} else {
			elapsed := time.Since(mt.start_time)
			if (elapsed.Milliseconds()+mt.TimeBuffer) > mt.SearchTimeLimit || mt.is_abort {
				break
			}
			DescendNode(mt, int(mt.BTree.RootColor^1), max_index, mt.Nthr, int(mt.BTree.Ply+1))
		}
		UnDo(&mt.BTree, mt.BTree.RootMoves[max_index-1], int(mt.BTree.RootColor))
	}
	mt.is_finished = true
}

func ExpandNode(mt *MCTSTree, color int, parent_index uint64, ply int) {
	var ulret = BBTest(IsAttacked(mt.BTree, int(mt.BTree.SQ_King[color]), color))
	if ulret == 0 && MateIn1Ply(mt.BTree, color) > 0 {
		// This case is mate.
		mt.NodeList[parent_index].WinRateSum = 0.0
		mt.NodeList[parent_index].LostRateSum = 1.0
		mt.NodeList[parent_index].EvalCount += 1
		return
	} else {
		if BBTest(IsAttacked(mt.BTree, int(mt.BTree.SQ_King[color^1]), color^1)) != 0 {
			var iret = IsDeclarationWin(mt.BTree)
			if iret == uint8(color+1) {
				mt.NodeList[parent_index].WinRateSum = 0.0
				mt.NodeList[parent_index].LostRateSum = 1.0
				mt.NodeList[parent_index].EvalCount += 1
				return
			} else if (color == 0 && iret == 2) || (color == 1 && iret == 1) {
				mt.NodeList[parent_index].WinRateSum = 1.0
				mt.NodeList[parent_index].LostRateSum = 0.0
				mt.NodeList[parent_index].EvalCount += 1
				return
			}
		}
	}

	// The color's king is checked.
	if ulret != 0 {
		var evasion_moves []uint32
		evasion_moves = GenEvasion(&mt.BTree, color, evasion_moves)
		// If there are no evasion moves, it's mate.
		if len(evasion_moves) == 0 {
			mt.NodeList[parent_index].WinRateSum = 1.0
			mt.NodeList[parent_index].LostRateSum = 0.0
			mt.NodeList[parent_index].EvalCount += 1
			return
		}
	}

	// Throw request to main thread.
	var str_sfen = "p," + ToSFEN(mt.BTree, color)
	mt.queue_to_main_thread_p <- str_sfen
	// The data format is like as "7776FU 0.65,2226FU 0.25,5556FU 0.1" .
	//close(mt.queue_to_main_thread_v) //後で精査する
	//time.Sleep(1 * time.Millisecond)
	str_receive := <-mt.queue_from_main_thread_p
	var str_temp0 = strings.Split(str_receive, ",")
	// The case of mate
	if str_receive == "mate" {
		mt.NodeList[parent_index].WinRateSum = 1.0
		mt.NodeList[parent_index].LostRateSum = 0.0
		mt.NodeList[parent_index].EvalCount += 1
		return
	}

	var limit = len(str_temp0)
	var moves []uint32
	var flag = false
	for i := 0; i < limit; i++ {
		var n Node = InitNode()
		n.color = uint8(color)
		var str_temp1 = strings.Split(str_temp0[i], " ")
		var m = CSA2Move(mt.BTree, str_temp1[0])
		var ifrom = From(m)
		var ito = To(m)
		if ifrom < Square_NB {
			// the case of discovered check
			if BBTest(IsPinnedOnKing(&mt.BTree, int(ifrom), Adirec[ifrom][ito], color)) != 0 {
				continue
			}
		}

		// the case of king capture
		// It's illegal.
		if CapPiece(m) == uint32(King) {
			continue
		}

		// the case of discovered check
		Do(&mt.BTree, m, color)
		if BBTest(IsAttacked(mt.BTree, int(mt.BTree.SQ_King[color]), color)) != 0 {
			UnDo(&mt.BTree, m, color)
			continue
		}
		n.ThisIndex = uint64(len(mt.NodeList))
		n.ParentIndex = mt.NodeList[parent_index].ThisIndex
		n.move = m
		moves = append(moves, m)
		mt.NodeList[parent_index].ChildIndexes = append(mt.NodeList[parent_index].ChildIndexes, n.ThisIndex)
		temp_v, _ := strconv.ParseFloat(str_temp1[1], 32) //already executed softmax function.
		n.PolicyResult = float32(temp_v)
		mt.NodeList = append(mt.NodeList, n)
		var current_index = n.ThisIndex
		var bt BoardTree = DeepCopy(mt.BTree)
		var result = PlayOut(&bt, color^1, ply)
		n.PlayoutCount = 1
		EvalNode(mt, current_index, color^1)
		UpdateParam(mt, current_index, result)
		UnDo(&mt.BTree, m, color)
		flag = true
	}
	if flag {
		mt.NodeList[parent_index].IsLeaf = false
	}
}

func PlayOut(bt *BoardTree, start_color int, temp_ply int) int {
	const ply_max int = 384
	var result = 2
	var ply = temp_ply
	var color = start_color
	for {
		if ply >= ply_max {
			break
		}
		var moves []uint32
		if BBTest(IsAttacked(*bt, int(bt.SQ_King[color]), color)) == 0 {
			var mate_move = MateIn1Ply(*bt, color)
			// It exists mate move.
			if mate_move != 0 {
				if color == int(bt.RootColor) {
					result = 0 // Root color wins.
				} else {
					result = 1 // Opponent color wins.
				}
				break
			}

			if BBTest(IsAttacked(*bt, int(bt.SQ_King[color^1]), color^1)) == 0 {
				// It's not checked in both colors.
				var iret = IsDeclarationWin(*bt)
				if iret == 1 && color == int(bt.RootColor) && color == 0 {
					result = 0 // Root color wins and black.
					break
				} else if iret == 2 && color == int(bt.RootColor) && color == 1 {
					result = 1 // Opponent color wins and root color is black.
					break
				}
			}
			moves = GenDrop(bt, color, moves)
			moves = GenNoCap(*bt, color, moves)
			moves = GenCap(*bt, color, moves)
		} else {
			moves = GenEvasion(bt, color, moves)
		}
		var legal_moves []uint32
		var limit = len(moves)
		for i := 0; i < limit; i++ {
			if IsMoveValid(bt, moves[i], color) {
				legal_moves = append(legal_moves, moves[i])
			}
		}
		if len(legal_moves) == 0 {
			if color != int(bt.RootColor) {
				result = 0 // Root color wins.
			} else {
				result = 1 // Opponent color wins.
			}
			break
		}
		var index = rand.Int31n(int32(len(legal_moves)))
		var move = legal_moves[index]
		Do(bt, move, color)
		color ^= 1
		ply += 1
	}
	return result
}

func EvalNode(mt *MCTSTree, node_index uint64, color int) {
	f, ok := mt.tt.value[mt.BTree.CurrentHash]
	if ok {
		// A datum exists in TT.
		mt.NodeList[node_index].WinRateSum = 1 - f
		mt.NodeList[node_index].LostRateSum = f
		mt.NodeList[node_index].EvalCount++
		return
	}

	// Throw request to main thread.
	var str_sfen = "v," + ToSFEN(mt.BTree, color)
	mt.queue_to_main_thread_v <- str_sfen
	//close(mt.queue_to_main_thread_v) //後で精査する
	//time.Sleep(1 * time.Millisecond)
	str_receive := <-mt.queue_from_main_thread_v
	temp_v, _ := strconv.ParseFloat(str_receive, 32)
	var v = float32(temp_v)

	mt.NodeList[node_index].WinRateSum = 1.0 - v
	mt.NodeList[node_index].LostRateSum = v
	mt.NodeList[node_index].EvalCount++
}

func UpdateParam(mt *MCTSTree, node_index uint64, result int) {
	mt.NodeList[node_index].TrialCount += 1
	if result == 0 {
		if mt.NodeList[node_index].color == mt.BTree.RootColor {
			mt.NodeList[node_index].WinCount += 1
		} else {
			mt.NodeList[node_index].LostCount += 1
		}
	} else if result == 1 {
		if mt.NodeList[node_index].color == mt.BTree.RootColor {
			mt.NodeList[node_index].LostCount += 1
		} else {
			mt.NodeList[node_index].WinCount += 1
		}
	} else {
		mt.NodeList[node_index].DrawCount += 1
	}
	var current_node = mt.NodeList[node_index]
	var delta float32
	var delta2 float32
	// Win rate is color's side.
	if current_node.color == mt.BTree.RootColor {
		delta = mt.NodeList[node_index].WinRateSum
		delta2 = mt.NodeList[node_index].LostRateSum
	} else {
		delta = mt.NodeList[node_index].LostRateSum
		delta2 = mt.NodeList[node_index].WinRateSum
	}
	if current_node.ParentIndex == 0 {
		return
	}
	for {
		var index = current_node.ParentIndex
		current_node = mt.NodeList[index]
		mt.NodeList[index].TrialCount += 1
		mt.NodeList[index].PlayoutCount += 1
		if result == 0 {
			if mt.NodeList[index].color == mt.BTree.RootColor {
				mt.NodeList[index].WinCount += 1
			} else {
				mt.NodeList[index].LostCount += 1
			}
		} else if result == 1 {
			if mt.NodeList[index].color == mt.BTree.RootColor {
				mt.NodeList[index].LostCount += 1
			} else {
				mt.NodeList[index].WinCount += 1
			}
		} else {
			mt.NodeList[index].DrawCount += 1
		}

		if current_node.color == mt.BTree.RootColor {
			mt.NodeList[index].WinRateSum += delta
			mt.NodeList[index].LostRateSum += delta2
		} else {
			mt.NodeList[index].WinRateSum += delta2
			mt.NodeList[index].LostRateSum += delta
		}

		mt.NodeList[index].EvalCount += 1
		if current_node.ParentIndex == 0 {
			break
		}
	}
}

func DescendNode(mt *MCTSTree, color int, node_index uint64, nthr uint16, ply int) {
	if len(mt.NodeList[node_index].ChildIndexes) == 0 {
		return
	}

	var idx uint64 = 0
	var flag = false
	for {
		elapsed := time.Since(mt.start_time)
		if (elapsed.Milliseconds()+mt.TimeBuffer) > mt.SearchTimeLimit || mt.is_abort {
			break
		}
		var ucb1_array []float32
		var t uint64 = 0
		var i uint64 = 0
		var limit = uint64(len(mt.NodeList[node_index].ChildIndexes))
		for {
			if i >= limit {
				break
			}
			idx = mt.NodeList[node_index].ChildIndexes[i]
			t += mt.NodeList[idx].PlayoutCount
			i += 1
		}
		i = 0
		for {
			if i >= limit {
				break
			}
			idx = mt.NodeList[node_index].ChildIndexes[i]
			var u float32 = mt.NodeList[idx].PolicyResult * float32(math.Sqrt(float64(t))) / float32(mt.NodeList[idx].PlayoutCount+1)
			var q float32 = 0
			if mt.NodeList[idx].EvalCount > 0 && mt.NodeList[idx].PlayoutCount > 0 {
				q = float32(((1-mt.ValueLambda)*(float32(mt.NodeList[idx].WinRateSum)/float32(mt.NodeList[idx].EvalCount)) + mt.ValueLambda*float32(float32(mt.NodeList[idx].WinCount)/float32(mt.NodeList[idx].PlayoutCount))))
			}
			ucb1_array = append(ucb1_array, u+q)
			i += 1
		}
		_, max_index := MaxInSlicef32(ucb1_array)
		idx = mt.NodeList[node_index].ChildIndexes[max_index]
		Do(&mt.BTree, mt.NodeList[idx].move, color)
		if mt.NodeList[idx].IsLeaf {
			elapsed = time.Since(mt.start_time)
			if (elapsed.Milliseconds() + mt.TimeBuffer) > mt.SearchTimeLimit {
				flag = true
				goto end
			}
			if mt.is_abort {
				break
			}
			if mt.NodeList[idx].TrialCount >= int64(mt.Nthr) {
				ExpandNode(mt, color^1, idx, ply+1)
				if mt.is_abort {
					break
				} else {
					var bt BoardTree
					bt = DeepCopy(mt.BTree)
					var result = PlayOut(&bt, color^1, ply)
					mt.NodeList[idx].PlayoutCount++
					EvalNode(mt, idx, color^1)
					if mt.is_abort {
						break
					}
					UnDo(&mt.BTree, mt.NodeList[idx].move, color)
					AscendNode(mt, color, idx, result)
					return
				}
			}
		} else {
			elapsed = time.Since(mt.start_time)
			if (elapsed.Milliseconds() + mt.TimeBuffer) > mt.SearchTimeLimit {
				flag = true
				goto end
			}
			if mt.is_abort {
				break
			}

			DescendNode(mt, color^1, idx, nthr, ply+1)
			return
		}
		UnDo(&mt.BTree, mt.NodeList[idx].move, color)
	end:
		if flag {
			var current_node = mt.NodeList[idx]
			var temp_color = color ^ 1
			for {
				var index = current_node.ParentIndex
				current_node = mt.NodeList[index]
				if current_node.ParentIndex == 0 {
					break
				}
				UnDo(&mt.BTree, current_node.move, temp_color)
				temp_color ^= 1
			}
			break
		}
	}
}

func AscendNode(mt *MCTSTree, color int, node_index uint64, result int) {
	mt.NodeList[node_index].TrialCount += 1
	if result == 0 {
		if mt.NodeList[node_index].color == mt.BTree.RootColor {
			mt.NodeList[node_index].WinCount += 1
		} else {
			mt.NodeList[node_index].LostCount += 1
		}
	} else if result == 1 {
		if mt.NodeList[node_index].color == mt.BTree.RootColor {
			mt.NodeList[node_index].LostCount += 1
		} else {
			mt.NodeList[node_index].WinCount += 1
		}
	} else {
		mt.NodeList[node_index].DrawCount += 1
	}
	var current_node = mt.NodeList[node_index]
	var delta float32 = 0.0
	var delta2 float32 = 0.0
	// Win rate is color's side.
	if current_node.color == mt.BTree.RootColor {
		delta = mt.NodeList[node_index].WinRateSum
		delta2 = mt.NodeList[node_index].LostRateSum
	} else {
		delta = mt.NodeList[node_index].LostRateSum
		delta2 = mt.NodeList[node_index].WinRateSum
	}
	if current_node.ParentIndex == 0 {
		return
	}
	var temp_color = color
	for {
		var index = current_node.ParentIndex
		current_node = mt.NodeList[index]
		mt.NodeList[index].TrialCount += 1
		mt.NodeList[index].PlayoutCount += 1
		if result == 0 {
			if mt.NodeList[index].color == mt.BTree.RootColor {
				mt.NodeList[index].WinCount += 1
			} else {
				mt.NodeList[index].LostCount += 1
			}
		} else if result == 1 {
			if mt.NodeList[index].color == mt.BTree.RootColor {
				mt.NodeList[index].LostCount += 1
			} else {
				mt.NodeList[index].WinCount += 1
			}
		} else {
			mt.NodeList[index].DrawCount += 1
		}

		if current_node.color == mt.BTree.RootColor {
			mt.NodeList[index].WinRateSum += delta
			mt.NodeList[index].LostRateSum += delta2
		} else {
			mt.NodeList[index].WinRateSum += delta2
			mt.NodeList[index].LostRateSum += delta
		}

		mt.NodeList[index].EvalCount += 1
		if current_node.ParentIndex == 0 {
			break
		}
		UnDo(&mt.BTree, current_node.move, temp_color)
		temp_color ^= 1
	}
}
