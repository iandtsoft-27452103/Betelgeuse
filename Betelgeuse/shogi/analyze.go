package shogi

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func AnalyzeRecord(num_tasks uint16, num_mate_tasks uint8, thinking_time uint16, mate_search_depth uint16,
	analyze_file_name string, record_file_name string, str_game_date string, str_match_name string,
	str_black_player string, str_white_player string) {
	fp, err := os.Create(analyze_file_name)
	if err != nil {
		panic(err)
	}
	//defer fp.Close()
	writer := bufio.NewWriter(fp)
	_, err = writer.WriteString("対局日：" + str_game_date + "\n\n")
	_, err = writer.WriteString("棋戦名：" + str_match_name + "\n\n")
	_, err = writer.WriteString("先手：" + str_black_player + "\n\n")
	_, err = writer.WriteString("後手：" + str_white_player + "\n\n")

	var str_result = ""
	var str_mate_pv = ""
	//var str_first_move = ""
	var mate_first_move uint32

	var cm = InitCommunicate()
	Boot(&cm)
	var temp_r = ReadRecords(record_file_name)
	var r = temp_r[0]
	var denomi = len(r.str_moves)
	var denomi_black = 0
	var denomi_white = 0
	var move_first_accuracy [Color_NB]int = [Color_NB]int{0, 0}
	var move_second_accuracy [Color_NB]int = [Color_NB]int{0, 0}
	var move_third_accuracy [Color_NB]int = [Color_NB]int{0, 0}
	var bt BoardTree = InitBoard()
	var color = 0
	for i := 0; i < len(r.str_moves); i++ { // len(r.str_moves)

		/*if i != (len(r.str_moves) - 1) {
			var m = CSA2Move(bt, r.str_moves[i])
			Do(&bt, m, color)
			color ^= 1
			continue
		}*/

		bt.RootColor = uint8(color)
		str_result = ""
		str_mate_pv = ""
		//str_first_move = ""
		mate_first_move = 0
		//strconv.Itoa(i)
		fmt.Printf("ply = " + strconv.Itoa(i+1) + "\n")
		/*if i == 5 {
			var a = "0"
			fmt.Printf(a)
		}*/
		str_result, str_mate_pv, mate_first_move, _ = SearchWrapper(&cm, &bt, num_tasks, num_mate_tasks, thinking_time, mate_search_depth, r.str_moves[i], &move_first_accuracy, &move_second_accuracy, &move_third_accuracy)
		if str_mate_pv != "" {
			var str_accurate = "×"
			var str_color = "+"
			if color == int(White) {
				str_color = "-"
			}
			if r.str_moves[i] == Move2CSA(mate_first_move) {
				str_accurate = "○"
				move_first_accuracy[color] += 1
			} else {
				var s = str_mate_pv[1:6]
				if r.str_moves[i] == s {
					str_accurate = "○"
					move_first_accuracy[color] += 1
				}
			}
			str_result = "ply=" + strconv.Itoa(i+1) + ", 棋譜の手: " + str_color + r.str_moves[i] + ", result= " + str_accurate + ", 詰みあり： " + str_mate_pv
			_, err = writer.WriteString(str_result + "\n")
			writer.Flush()
		} else {
			var str_color = "+"
			if color == int(White) {
				str_color = "-"
			}
			str_result = "ply=" + strconv.Itoa(i+1) + ", 棋譜の手: " + str_color + r.str_moves[i] + ", " + str_result
			_, err = writer.WriteString(str_result + "\n")
			writer.Flush() // 注）これを入れておかないと訪問回数が文字化けする率が高くなる模様。ただし入れておいても一定数文字化けする。
		}
		var m = CSA2Move(bt, r.str_moves[i])
		Do(&bt, m, color)
		color ^= 1
		if color == int(Black) {
			denomi_black += 1
		} else {
			denomi_white += 1
		}
	}
	_, err = writer.WriteString("\n")
	var temp_n = move_first_accuracy[0]
	var br float64 = float64(temp_n) / float64(denomi_black)
	br = RoundToN(br, 2)
	var temp_s0 = strconv.FormatFloat(br, 'f', 2, 64)
	var temp_s1 = "先手一致率：" + strconv.Itoa(temp_n) + " / " + strconv.Itoa(denomi_black) + " = " + temp_s0
	_, err = writer.WriteString(temp_s1 + "\n")
	temp_n = move_first_accuracy[1]
	br = float64(temp_n) / float64(denomi_white)
	temp_s0 = strconv.FormatFloat(br, 'f', 2, 64)
	temp_s1 = "後手一致率：" + strconv.Itoa(temp_n) + " / " + strconv.Itoa(denomi_white) + " = " + temp_s0
	_, err = writer.WriteString(temp_s1 + "\n")
	temp_n = move_first_accuracy[0] + move_first_accuracy[1]
	br = float64(temp_n) / float64(denomi)
	temp_s0 = strconv.FormatFloat(br, 'f', 2, 64)
	temp_s1 = "全体一致率：" + strconv.Itoa(temp_n) + " / " + strconv.Itoa(denomi) + " = " + temp_s0
	_, err = writer.WriteString(temp_s1 + "\n")
	temp_n = move_first_accuracy[0] + move_second_accuracy[0] + move_third_accuracy[0] + move_first_accuracy[0] + move_second_accuracy[0] + move_third_accuracy[0]
	br = float64(temp_n) / float64(denomi)
	temp_s0 = strconv.FormatFloat(br, 'f', 2, 64)
	temp_s1 = "3位以内の確率：" + strconv.Itoa(temp_n) + " / " + strconv.Itoa(denomi) + " = " + temp_s0
	_, err = writer.WriteString(temp_s1 + "\n")
	_, err = writer.WriteString("解析エンジン名：Betelgeuse Ver.1.0.0\n")
	//ThrowRequest(&cm, "p,lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1")
	//var s = shogi.ReceiveResponse(&cm)
	//cm.cmd.SysProcAttr.Token.Close()
	writer.Flush()
	fp.Close()
}

// This function is created by Copilot
func RoundToN(val float64, n int) float64 {
	if n < 0 {
		panic("n must be >= 0")
	}
	// math.Round は float64 専用なので変換
	pow := math.Pow(10, float64(n))
	rounded := math.Round(float64(val)*pow) / pow
	return rounded
}

func SearchWrapper(cm *Communicate, bt *BoardTree, num_tasks uint16, num_mate_tasks uint8, thinking_time uint16, mate_search_depth uint16, str_record_move string,
	move_first_accuracy *[2]int, move_second_accuracy *[2]int, move_third_accurasy *[2]int) (string, string, uint32, string) {

	var mate_search_trees []*MateSearchTree
	var checkMoves = GenRootMoves(*bt)
	if len(checkMoves) > 0 {
		for i := 0; i < int(num_mate_tasks); i++ {
			var mst MateSearchTree
			mst.bt = DeepCopy(*bt)
			for j := 0; j < len(checkMoves); j++ {
				mst.root_check_moves = append(mst.root_check_moves, checkMoves[j])
			}
			mst.max_ply = int(mate_search_depth)
			mate_search_trees = append(mate_search_trees, &mst)
		}
	}

	done_mate := make(chan int)

	if len(checkMoves) > 0 {
		for i := 0; i < int(num_mate_tasks); i++ {
			go func(id int) {
				MateSearchWrapper(mate_search_trees[i], int(mate_search_depth))
				done_mate <- id
			}(i)
		}
	}

	var mcts_trees []*MCTSTree
	for i := 0; i < int(num_tasks); i++ {
		var mcts MCTSTree
		mcts.TaskNumber = uint16(i)
		mcts.BTree = DeepCopy(*bt)
		mcts.ValueLambda = 0.5 // ※引数で設定できるようにするか？
		mcts.Nthr = 30         // ※引数で設定できるようにするか？
		mcts.TimeBuffer = 500  // ミリ秒
		//mcts.BTree.RootMoves = GenRootMoves2(bt)
		var moves []uint32 = GenRootMoves2(bt)
		mcts.BTree.RootMoves = SetRootOutput(cm, &mcts, moves)
		mcts.SearchTimeLimit = int64(thinking_time) * 1000
		mcts.queue_from_main_thread_p = make(chan string)
		mcts.queue_to_main_thread_p = make(chan string, 1)
		mcts.queue_from_main_thread_v = make(chan string, 1)
		mcts.queue_to_main_thread_v = make(chan string, 1)
		mcts.tt = TT2{}
		mcts.is_abort = false
		mcts.is_finished = false
		mcts_trees = append(mcts_trees, &mcts)
	}

	done := make(chan int)
	var st time.Time

	for i := 0; i < int(num_tasks); i++ {
		go func(id int) {
			mcts_trees[i].start_time = time.Now()
			Root(mcts_trees[i])
			done <- id
		}(i)
	}

	st = time.Now()

	var task_index = 0
	var task_finished = 0
	var mate_task_number = -1
	var str_value_request []string
	var str_policy_request []string
	var num_value_request []int
	var num_policy_request []int
	for i := 0; i < int(num_tasks); i++ {
		str_value_request = append(str_value_request, "")
		str_policy_request = append(str_policy_request, "")
	}

	var flag_v = false
	var flag_p = false

	for {
		/*if len(mcts_trees[0].queue_to_main_thread_v) > 0 {
			s := <-mcts_trees[0].queue_to_main_thread_v
			ThrowRequest(cm, s)
			str_response := ReceiveResponse(cm)
			str_response = strings.Replace(str_response, "\r\n", "", -1)
			mcts_trees[0].queue_from_main_thread_v <- str_response
		}*/
		if str_value_request[task_index] == "" && len(mcts_trees[task_index].queue_to_main_thread_v) > 0 {
			str_value_request[task_index] = <-mcts_trees[task_index].queue_to_main_thread_v
			num_value_request = append(num_value_request, task_index)
		}

		/*if len(mcts_trees[0].queue_to_main_thread_p) > 0 {
			s := <-mcts_trees[0].queue_to_main_thread_p
			ThrowRequest(cm, s)
			str_response := ReceiveResponse(cm)
			str_response = strings.Replace(str_response, "\r\n", "", -1)
			mcts_trees[0].queue_from_main_thread_p <- str_response
		}*/

		if str_policy_request[task_index] == "" && len(mcts_trees[task_index].queue_to_main_thread_p) > 0 {
			str_policy_request[task_index] = <-mcts_trees[task_index].queue_to_main_thread_p
			num_policy_request = append(num_policy_request, task_index)
		}

		if len(num_value_request) > 0 {
			if num_tasks >= 4 && len(num_value_request) < 2 && !flag_v {
				time.Sleep(1 * time.Millisecond)
				flag_v = true
			} else {
				var s = "v,"
				for i := 0; i < len(num_value_request); i++ {
					var temp0 = str_value_request[num_value_request[i]]
					var temp1 = strings.Split(temp0, ",")
					s += temp1[1]
					if i != len(num_value_request)-1 {
						s += ","
					}
				}
				ThrowRequest(cm, s)
				str_response := ReceiveResponse(cm)
				str_response = strings.Replace(str_response, "\r\n", "", -1)
				s2 := strings.Split(str_response, ",")
				for i := 0; i < len(num_value_request); i++ {
					mcts_trees[num_value_request[i]].queue_from_main_thread_v <- s2[i]
					str_value_request[num_value_request[i]] = ""
				}
				num_value_request = num_value_request[:0]
				flag_v = false
			}
		}

		if len(num_policy_request) > 0 {
			if num_tasks >= 4 && len(num_policy_request) < 2 && !flag_p {
				time.Sleep(1 * time.Millisecond)
				flag_p = true
			} else {
				var s = "p,"
				for i := 0; i < len(num_policy_request); i++ {
					var temp0 = str_policy_request[num_policy_request[i]]
					var temp1 = strings.Split(temp0, ",")
					s += temp1[1]
					if i != len(num_policy_request)-1 {
						s += ","
					}
				}
				ThrowRequest(cm, s)
				str_response := ReceiveResponse(cm)
				str_response = strings.Replace(str_response, "\r\n", "", -1)
				for i := 0; i < len(num_policy_request); i++ {
					mcts_trees[num_policy_request[i]].queue_from_main_thread_p <- str_response
					str_policy_request[num_policy_request[i]] = ""
				}
				num_policy_request = num_policy_request[:0]
				flag_p = false
			}
		}

		for i := 0; i < int(num_tasks); i++ {
			if mcts_trees[i].is_finished {
				task_finished++
			}
		}

		if task_finished == int(num_tasks) {
			break
		} else {
			task_finished = 0
		}

		task_index++

		if task_index == int(num_tasks) {
			task_index = 0
		}

		if len(checkMoves) > 0 {
			for i := 0; i < int(num_mate_tasks); i++ {
				if mate_search_trees[i].is_mate_root {
					mate_task_number = i
					break
				}
			}
		}

		if mate_task_number >= 0 {
			break
		}

		// 時間制限で脱出した際でもgo routineは終了できているようだ。
		elapsed := time.Since(st)
		if elapsed.Milliseconds() > int64(thinking_time)*1000 {
			//break
		}

		/*id := <-done
		if id == 0 {
			fmt.Printf("BBB")
			break
		}*/
	}

	// 詰みを見つけた際の処理を書く
	str_mate_pv := ""
	str_result := ""
	str_first_move := ""
	var mate_first_move uint32 = 0
	if mate_task_number >= 0 {
		mate_first_move = mate_search_trees[mate_task_number].first_move
		str_mate_pv = mate_search_trees[mate_task_number].str_pv
	} else {
		_, moves, f, t := TotalParam(int(num_tasks), mcts_trees)
		str_first_move = Move2CSA(moves[0])
		for i := 0; i < len(moves); i++ {
			var sr = ""
			if Move2CSA(moves[i]) == str_record_move {
				sr = "result= ○"
				switch i {
				case 0:
					move_first_accuracy[bt.RootColor] += 1
				case 1:
					move_second_accuracy[bt.RootColor] += 1
				case 2:
					move_third_accurasy[bt.RootColor] += 1
				}
			} else {
				sr = "result= ×"
			}
			var str_color string
			if bt.RootColor == uint8(Black) {
				str_color = "+"
			} else {
				str_color = "-"
			}

			var s = "候補手" + strconv.Itoa(i+1) + ": " + str_color + Move2CSA(moves[i]) + ", " + sr + ", 勝率：" + strconv.FormatFloat(float64(f[i]), 'f', 2, 64) + ", 訪問回数" + strconv.FormatInt(t[i], 10) + ", "
			str_result += s
		}
	}

	/*mcts_trees, moves, f, t := TotalParam(int(num_tasks), mcts_trees)
	str_result := ""
	str_first_move := Move2CSA(moves[0])
	for i := 0; i < len(moves); i++ {
		var sr = ""
		if Move2CSA(moves[i]) == str_record_move {
			sr = "result= ○"
			switch i {
			case 0:
				move_first_accuracy[bt.RootColor] += 1
			case 1:
				move_second_accuracy[bt.RootColor] += 1
			case 2:
				move_third_accurasy[bt.RootColor] += 1
			}
		} else {
			sr = "result= ×"
		}
		var str_color string
		if bt.RootColor == uint8(Black) {
			str_color = "+"
		} else {
			str_color = "-"
		}

		var s = "候補手" + strconv.Itoa(i+1) + ": " + str_color + Move2CSA(moves[i]) + ", " + sr + ", 勝率：" + strconv.FormatFloat(float64(f[i]), 'f', 2, 64) + ", 訪問回数" + strconv.Itoa(int(t[i])) + ", "
		str_result += s
	}*/
	return str_result, str_mate_pv, mate_first_move, str_first_move
}
