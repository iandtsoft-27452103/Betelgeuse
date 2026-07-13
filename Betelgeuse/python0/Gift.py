from os import error
import torch
import torch.nn as nn
import numpy
import bitop
import board
import color as C
#import attack
#import inout
import feature
import piece
import gencap
import gennocap
import gendrop
import genevasion
import sfen
import policy
import value
import csa
#import test3
#import inout
#import analyze
import sys
from pathlib import Path

def main():

    #bo = board.board()
    #cls_io = inout.inout()
    #cls_io.read_records('20220403_nhk_hai.txt')
    #t = test3.test3()
    #t.test4(cls_io.rec, bo)
    #a = analyze.analyze()
    #a.analyze_policy(cls_io.rec[0], bo)
    #return
    try:
        device = torch.device("cuda:0")
        model = policy.policy()
        current_path = Path.cwd()
        policy_path = current_path.joinpath("model.pth")
        value_path = current_path.joinpath("model_value.pth")
        print('aaa', flush=True)
        sys.stdout.flush()
        model.load_state_dict(torch.load(r"C:\test\model.pth", weights_only=True))
        #model.load_state_dict(torch.load(policy_path, weights_only=True))
        print('bbb', flush=True)
        sys.stdout.flush()
        model = model.to(device)
        model_value = value.value()
        model_value.load_state_dict(torch.load(r"C:\test\model_value.pth", weights_only=True))
        #model_value.load_state_dict(torch.load(value_path, weights_only=True))
        model_value = model_value.to(device)
        bo = board.board()
        ft = feature.feature(bo)
        bi = bitop.bitop(bo)
        pc = piece.piece()
        cls_sfen = sfen.sfen(bo, bi)
        cls_csa = csa.csa()
        at = ft.at
        c = C.color()
        cls_cap = gencap.gencap(bo, bi, pc, at, c)
        cls_nocap = gennocap.gennocap(bo, bi, pc, at, c)
        cls_drop = gendrop.gendrop(bo, bi, pc, at, c)
        cls_eva = genevasion.genevasion(bo, bi, pc, at, c)
        cm = ','

        temp1 = torch.zeros((1,63,9,9), dtype=torch.float, device=device)
        temp2 = torch.zeros((1,56,9,9), dtype=torch.float, device=device)
        temp3 = torch.zeros((1,1,9,9), dtype=torch.float, device=device)
        model.eval()

        #print('a')

        #1度実行しておくと初回の計算が速くなる。理由は不明。
        with torch.no_grad():
            _ = model.forward(temp1, temp2, temp3)
            model_value.batch_size = 1
            _ = model_value.forward(temp1, temp2, temp3)

        print('Hello World.', flush=True)
        sys.stdout.flush()
        cmd_line = ''
        #s0 = 'p'
        #s1 = "lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1";
        #s2 = "l2g3nl/2skgr3/1p2psbpp/2pp2p2/pn3p1P1/2PPP1P2/PP1GSP2P/LK1B3R1/1NSG3NL b - 1";
        #s3 = "7nl/2k1gl1+R1/+LG2Pr2p/2ppp1p2/1s3p3/1pPP2P2/2G1B1+b1P/2G6/1K7 w 2SNL3Ps2n3p 1";
        #cmd_line = s0 + ',' + s1 + ',' + s2 + ',' + s3
        #cmd_line = 'v,1r5nl/3gk1g2/2ns1s1p1/l1pp1ppbp/4P4/p1PP2P1P/1P1S1S3/PBG6/LN1K1+r2L b Pn3p 1,lnsgkgsnl/1r5b1/ppppppppp/9/9/9/PPPPPPPPP/1B5R1/LNSGKGSNL b - 1'
        #cmd_line = 'p,lnsgkgsnl/1r5b1/ppppppppp/9/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL w - 1,lnsgkgsnl/1r5b1/ppppppppp/9/9/2P6/PP1PPPPPP/1B5R1/LNSGKGSNL w - 1'

        while True:
            cmd_line = input()
            cmd = cmd_line.split(',')
            if cmd[0] == 'quit':
                break
            
            str_mode = cmd[0]
            limit = len(cmd) - 1
            li1 = []
            li2 = []
            li3 = []
            li_moves = []
            li_csa_moves = []
            li_directions = []
            #is_mate = False
            str_ret = ''
            cnt = 0
            for i in range(limit):
                bo.clear_board()
                cls_sfen.str_sfen = cmd[i + 1]
                cls_sfen.parse_sfen()
                if str_mode == 'p':
                    moves = []
                    csa_moves = []
                    directions = []
                    if cls_sfen.board.color == 0:
                        if at.is_attacked(cls_sfen.board, cls_sfen.board.sq_king[0], cls_sfen.board.color ^ 1) == 0:
                            cls_cap.b_gen_captures(moves)
                            cls_nocap.b_gen_nocaptures(moves)
                            cls_drop.b_gen_drop(moves)
                        else:
                            cls_eva.b_gen_evasion(moves)
                    else:
                        if at.is_attacked(cls_sfen.board, cls_sfen.board.sq_king[1], cls_sfen.board.color ^ 1) == 0:
                            cls_cap.w_gen_captures(moves)
                            cls_nocap.w_gen_nocaptures(moves)
                            cls_drop.w_gen_drop(moves)
                        else:
                            cls_eva.w_gen_evasion(moves)
                    if len(moves) == 0:
                        li_empty = []
                        li_moves.append(li_empty)
                        li_csa_moves.append(li_empty)
                        li_directions.append(li_empty)
                    else:
                        for j in range(len(moves)):
                            _, direc = ft.make_output_labels(cls_sfen.board, moves[j])
                            directions.append(direc)
                            str_csa = cls_csa.board_to_csa(cls_sfen.board, pc, moves[j])
                            csa_moves.append(str_csa)
                        li_moves.append(moves)
                        li_csa_moves.append(csa_moves)
                        li_directions.append(directions)
                        fe = ft.make_input_features1(cls_sfen.board)
                        li1.append(fe)
                        fe = ft.make_input_features2(cls_sfen.board)
                        li2.append(fe)
                        fe = ft.make_input_features3(cls_sfen.board, cls_sfen.board.color)
                        li3.append(fe)
                        cnt += 1
                elif str_mode == 'v':
                    fe = ft.make_input_features1(cls_sfen.board)
                    li1.append(fe)
                    fe = ft.make_input_features2(cls_sfen.board)
                    li2.append(fe)
                    fe = ft.make_input_features3(cls_sfen.board, cls_sfen.board.color)
                    li3.append(fe)
                    cnt += 1
            if str_mode == 'p':
                li1 = numpy.array(li1)
                li2 = numpy.array(li2)
                li3 = numpy.array(li3)
                x1 = torch.tensor(li1, dtype = torch.float)
                x1 = x1.to(device)
                x2 = torch.tensor(li2, dtype = torch.float)
                x2 = x2.to(device)
                x3 = torch.tensor(li3, dtype = torch.float)
                x3 = x3.to(device)                
                model.eval()
                with torch.no_grad():
                    y = model.forward(x1, x2, x3)
                    y = y.reshape(cnt, 256, 81)#アウトプット層は本来チャンネル数を32にすべきである。
                    y = y.tolist()
                    
                    k = 0
                    for i in range(len(li_moves)):
                        #print(i)
                        if len(li_moves[i]) == 0:
                            if i == 0 and len(li_moves) == 0:
                                str_ret += 'mate'
                            elif i == 0 and len(li_moves) > 0:
                                str_ret += 'mate:'
                            else:
                                if i == len(li_moves) - 1:
                                    str_ret += 'mate'
                                else:
                                    str_ret += 'mate:'
                            continue
                        #print('bbb')
                        #print(i)
                        l = []
                        for j in range(len(li_moves[i])):
                            v = y[k][li_directions[i][j]][li_moves[i][j].ito]
                            l.append(v)
                        #print('ccc')
                        length = len(l)
                        if length > 7:
                            length = 8
                        l = numpy.array(l)
                        l = torch.tensor(l, dtype = torch.float)
                        l = torch.softmax(l, 0, dtype = torch.float)
                        l_copy = l
                        indexes = []
                        for _ in range(length):
                            idx = int(torch.argmax(l_copy))
                            indexes.append(idx)
                            l_copy = l_copy.tolist()
                            l_copy[idx] = 0
                            l_copy = torch.tensor(l_copy, dtype = torch.float)

                        l = l.tolist()
                        vl = []
                        for j in range(length):
                            vl.append(l[indexes[j]])
                        s = ''
                        for j in range(length):
                            if j != 0:
                                s += cm
                            s += li_csa_moves[i][indexes[j]]
                            s += " " + str(vl[j])
                            #if i != cnt - 1 and j == length - 1:
                            #s += ":"
                        str_ret += s
                        if i != len(li_moves) - 1:
                            str_ret += ":"
                        k += 1
            elif str_mode == 'v':
                li1 = numpy.array(li1)
                li2 = numpy.array(li2)
                li3 = numpy.array(li3)
                x1 = torch.tensor(li1, dtype = torch.float)
                x1 = x1.to(device)
                x2 = torch.tensor(li2, dtype = torch.float)
                x2 = x2.to(device)
                x3 = torch.tensor(li3, dtype = torch.float)
                x3 = x3.to(device)                
                model_value.eval()
                with torch.no_grad():
                    model_value.batch_size = cnt
                    y = model_value.forward(x1, x2, x3)
                    y = torch.sigmoid(y)
                    y = y.tolist()
                    for i in range(cnt):
                        if i != 0:
                            str_ret += cm
                        str_ret += str(y[i])
            if str_ret != "":
                print(str_ret)
                #break
    except Exception:
        file = open('error_log.txt', 'w', 1, 'UTF-8')
        file.write(cmd_line)
        file.close()

if __name__ == "__main__":
    main()