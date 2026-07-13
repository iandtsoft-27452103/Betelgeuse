# About Pull Request

This repository is read-only, so Pull Request is not accepted. Thank you for your understanding.

# Caution

If you use this software for commercial, you must pay fees to NVIDIA.

# Betelgeuse

Betelgeuse is a shogi engine framework, written by Go and Python.

Shogi is game like chess.

## Source Code Explanation For Go

(1) analyze.go : Functions for analyzing records.

(2) atkop.go : Functions of piece attacks.

(3) bitop.go : Functions of bit operations.

(4) board.go : Functions of shogi board.

(5) common.go : Common constants and variables.

(6) communicate.go : Functions of communicating with Python.

(7) csa.go : Functions of CSA format.

(8) genmoves.go : Functions for generating moves.

(9) hash.go : Functions of hash values.

(10) io.go : Functions for reading records.

(11) mate.go : Functions for mate search.

(12) mate1ply.go : Function for mate in one ply.

(13) mcts.go : Functions for Monte Carlo Tree Search.

(14) move.go : Functions for moves.

(15) main.go: The entry point of this software.

(16) pext_amd64.s: The function for pext bitboards for getting attacks from rook, dragon, bishop, horse and lance. This function is created by Copilot.

(17) sfen.go : Functions for SFEN.

(18) test.go : Functions for testing.

(19) tt.go :Functions for transposition table.

## Source Code Explanation For Python

The source codes for predictions are written by Python. These are reused from Asklepios. So there exists the same bugs as Asklepios.

## Convolutional Neural Network

Convolutional Neural Network is composed of Policy Network and Value Network. Policy Network predicts best move of current position. Value Network evaluates current position.

### Input Features of Neural Network

(1) position of pieces

(2) effects of pieces

(3) difference of effects for every square

(4) turn

### Output Labels of Policy Neural Network

(1) position of move to

(2) move direction

The number of output labels are 32x81.

### Output Labels of Value Neural Network

win or lost

## Learning functions

Learning functions are written by PyTorch. Policy Neural Network learns multi task classification. Value Neural Network learns binary classification.

### Record format for learning

Record format for learning is as below.

B,119,2726FU,3334FU,7776FU,4344FU, ...

First column is game result. Second Column is game ply. The following columns are moves. When learning Value Neural Network, you use first column. And learning Policy Neural Network, you use columns from third column to last column.

This software use the records created by Gikou-2.0.1.

## Operating environment

(1) OS: Windows 11 Pro

(2) Go Version: 1.25.4

(3) Memory: 16GB or more. About 32GB is recommended.

(4) Memory usage on Go side: Sorry, I didn't measured.

(6) Memory usage on the PyTorch side: Sorry, I didn't measured.

(7) Python's version: 3.13.14

(8) PyTorch's version: 2.13.0

(9) It is necessary that CUDA and cuDNN corresponding to PyTorch are installed.

## How to build

Start console and execute the command below.

go build

## How to use

If you navigate to the cnn folder and execute the start.bat, the specified game record will be analyzed. The command-line arguments are as follows:

(1) The number of tasks for MCTS   max = 8, min = 1

(2) The number of tasks for mate search   max = 2, min = 1

(3) Thinking seconds per one move

(4) Mate search depth

(5) Output file name

(6) Record file name for analyzing

(7) Game date

(8) Match name

(9) Black Player name

(10) White Player name

About (1) and (2), there is no upper limit to the number of tasks, in theory.
See analyze.go

## Known bugs

(1) The result file by analyzing record output trial counts, but this counts are minimum int value, occasionally.

(2) The label table in Python source code includes incorrect class numbers.

## References

I developed this software referring to the softwares as below.

(1) Bonanza

(2) Apery

(3) YaneuraOu

(4) Gikou

(5) dlshogi

As far as I know, the source code for Bonanza and dlshogi is currently not publicly available.

I developed this software referring to the books as below. All books are written in Japanese, so I write the name of the books in Japanese.

(1) 山岡忠夫(2018),『将棋AIで学ぶディープラーニング』マイナビ出版 

(2) 山岡忠夫、加納邦彦(2021), 『強い将棋ソフトの創りかた　Pythonで実装するディープラーニング将棋AI』マイナビ出版

(3) 大槻知史(著)、三宅陽一郎(監修)(2018), 『最強囲碁AI アルファ碁解体新書　増補改訂版』翔泳社

(4) 原田達也(2017), 機械学習プロフェッショナルシリーズ『画像認識』講談社

## About the future

I think I'll add alpha-beta search functions with goroutine.
