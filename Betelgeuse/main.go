package main

import (
	"Betelgeuse/shogi"
	"fmt"
	"os"
	"strconv"
)

func main() {

	if len(os.Args) != 11 {
		fmt.Println("Usage: go run main.go <num1> <num2>")
		return
	}

	temp_n, _ := strconv.Atoi(os.Args[1])
	var num_tasks = uint16(temp_n)
	temp_n, _ = strconv.Atoi(os.Args[2])
	var num_mate_tasks = uint8(temp_n)
	temp_n, _ = strconv.Atoi(os.Args[3])
	var thinking_time = uint16(temp_n)
	temp_n, _ = strconv.Atoi(os.Args[4])
	var mate_search_depth = uint16(temp_n)
	analyze_file_name := os.Args[5]
	record_file_name := os.Args[6]
	str_game_date := os.Args[7]
	str_match_name := os.Args[8]
	str_black_player := os.Args[9]
	str_white_player := os.Args[10]

	shogi.Init()
	shogi.IniRand(5489)
	shogi.IniRandomTable()
	fmt.Printf("Hello, World!\n")

	shogi.AnalyzeRecord(num_tasks, num_mate_tasks, thinking_time, mate_search_depth,
		analyze_file_name, record_file_name, str_game_date, str_match_name, str_black_player, str_white_player)

	//shogi.AnalyzeRecord(4, 1, 3, 9, "analyze_result.txt", "20220403_nhk_hai.txt", "2022/4/3", "第72回NHK杯１回戦", "木村一基九段", "黒田尭之五段")

	//shogi.TestCommunicatePython()

	/*var cm = shogi.InitCommunicate()
	shogi.Boot(&cm)
	shogi.ThrowRequest(&cm, "p,lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1")
	var s = shogi.ReceiveResponse(&cm)*/
	//fmt.Printf(s)
	//fmt.Printf("処理時間: %s ms\n", time.Now())
}
