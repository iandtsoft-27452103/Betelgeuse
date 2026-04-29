# About Pull Request

This repository is read-only, so Pull Request is not accepted. Thank you for your understanding.

# Betelgeuse

Betelgeuse is a shogi engine framework, written by Go.

Shogi is game like chess.

## Source Code Explanation

(1) atkop.go : Functions of piece attacks.

(2) bitop.go : Functions of bit operations.

(3) board.go : Functions of shogi board.

(4) common.go : Common constants and variables.

(5) csa.go : Functions of CSA format.

(6) genmoves.go : Functions for generating moves.

(7) hash.go : Functions of hash values.

(8) io.go : Functions for reading records.

(9) mate1ply.go : Function for mate in one ply.

(10) move.go : Functions for moves.

(11) main.go: The entry point of this software.

(12) pext_amd64.s: The function for pext bitboards for getting attacks from rook, dragon, bishop, horse and lance. This function is created by Copilot.

(13) sfen.go : Functions for SFEN.

(14) test.go : Functions for testing.

(15) tt.go :Functions for transposition table.

## Operating environment

(1) OS: Windows 11 Pro

(2) Go Version: 1.25.4

## How to build

Start console and execute the command below.

go build

## References

I developed this software referring to the softwares as below.

(1) Bonanza

(2) Apery

(3) YaneuraOu

(4) Gikou

(5) dlshogi

As far as I know, the source code for Bonanza and dlshogi is currently not publicly available.

## About the future

I think I'll add search functions and analyze records functions.
