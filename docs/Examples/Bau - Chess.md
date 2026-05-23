#example 

[chass.bau](https://github.com/thomasmueller/bau-lang/blob/main/src/test/resources/org/bau/converter/chess.bau)

```js
terminal : org.bau.os.terminal  
keyCode  : org.bau.os.keyCode  
StringBuilder: org.bau.StringBuilder  
  
main: {   
  
    terminal.is_terminal() ? {  
       print("Not a terminal")  
       return  
    }    
    terminal.enable_raw_mode(refresh_screen)  
    init()  
    loop {  
       refresh_screen()  
       key : terminal.read_editor_key()  
       key < 0 ? { break }  
       key == 0 ? { continue }  
       match  
        { key == ord('q') => break },  
         { key == ord(' ') =>   
          state == 0 ? {  
             p : xx + 8 * yy  
             board[p] == 0 || { is_black(board[p]) != black_turn }  ? {  
                continue                 
             }  
             from : p  
             state : 1  
             moves : get_possible_moves(p, 0)  
             while({ moves != 0 },{  
                target : ints.number_of_trailing_zeros(moves)  
                board[target] |= 16  
                moves ^= 1 << target  
             })  
             refresh_screen()  
          },{ state == 1 ? {  
             p : xx + 8 * yy  
             allowed : (board[p] & 16) != 0  
             1.to(64).each({  
                i Int  
                board[i] &= ~16  
             })  
              allowed ? {  
                last_white = move(from, p)  
                show_cursor = 0  
                refresh_screen()  
                move : negamax(1, 5, not black_turn, -ints.MAX_INT, ints.MAX_INT)  
                move != 0 ? {  
                   last_black = move  
                   move(move)  
                }  
             }  
             }  
          }  
             show_cursor = 1  
             refresh_screen()  
             state = 0  
       }, { key ==  ord('c') =>   
          p : xx + 8 * yy  
          state == 0 ? {  
             from : p  
             state = 1  
          },{  
             last_white = move(from, p)  
             refresh_screen()  
             state = 0  
          }  
       }, { key ==  ord('u') =>   
          last_black != 0 &&  {last_white != 0} ? {  
             undo(last_black)  
             undo(last_white)  
             refresh_screen()  
          }  
          last_black = 0  
          last_white = 0  
       }, { key == ord('s') =>   
          black_turn = not black_turn  
          show_cursor = 0  
          refresh_screen()  
          move : negamax(1, 5, not black_turn, -ints.MAX_INT, ints.MAX_INT)  
          move != 0 ? {  
             move(move)  
          }  
          show_cursor = 1  
          refresh_screen()  
          state = 0  
       },  
       { key == keyCode.ARROW_RIGHT =>  xx = ints.min(7, xx + 1) },  
       { key == keyCode.ARROW_LEFT  =>  xx = ints.max(0, xx - 1) },  
       { key == keyCode.ARROW_UP    =>  yy = ints.max(0, yy - 1) },   
       { key == keyCode.ARROW_DOWN  =>  yy = ints.min(7, yy + 1) }  
  }  
}  
  
// UI  
xx : 3  
yy : 6  
state : 0  
show_cursor : 1  
from : 3  
black_turn : 0  
last_black : 0  
last_white : 0  
  
  
KING   : 1  
QUEEN  : 2  
ROOK   : 3  
BISHOP : 4  
KNIGHT : 5  
PAWN   : 6  
BLACK  : 6  
board  :[Int]()  
castling_flags : 0   
pawn_moved_2 : 0  
turn : 0  
  
  
refresh_screen : {  
   buff : StringBuilder()  
   buff.data = [Int]()   
   // hide cursor, go home  
    buff.append('\x1b[?25l\x1b[H\x1b[0m')  
    buff.append('   a  b  c  d  e  f  g  h  \r\n')  
    8.times().each({  
        y Int  
        buff.append(ints.to_str(8 - y))  
        buff.append(' ')  
        8.times().each({  
         x Int  
         b : board[x + 8 * y]  
         allowed : (b & 16) != 0  
         b = b & 16  
         match   
           { allowed != 0     =>  buff.append('\x1b[30;100m') },   
           { (x + y) % 2 == 1 =>  buff.append('\x1b[30;107m') }  
         
         is_black(b) ? {  
          // red  
          buff.append('\x1b[31m')  
         }, {  
          // blue  
          buff.append('\x1b[94m')  
         }  
         buff.append(' ')  
         b == 0 ? {  
          buff.append(' ')  
         }, {  
          b2 : [Int]()  
          b2[0] = 0xe2  
          b2[1] = 0x99  
          b2[2] = 0x94 + (b - 1)  
          buff.append(b2)  
         }  
         buff.append(' ')  
         // reset all attributes  
         buff.append('\x1b[0m')  
         buff.append(' \r\n')  
      })  
    })  
    buff.append('  arrows:select space:move\r\n')  
    buff.append('  u:undo  s:switch  q:quit')  
    buff.append('\x1b[')  
    buff.append(ints.to_str(yy + 2))  
    buff.append(';')  
    buff.append(ints.to_str(1 + 3 * xx + 3))  
    buff.append('H')  
    show_cursor ? {  
        buff.append('\x1b[?25h')  
    }  
    terminal.write_to_terminal(buff.data, buff.len)       
  }  
    
init : {  
    board[0] = ROOK  
    board[1] = KNIGHT  
    board[2] = BISHOP  
    board[3] = QUEEN  
    board[4] = KING  
    board[5] = BISHOP  
    board[6] = KNIGHT  
    board[7] = ROOK  
    1.to(8).each({  
        i Int  
        board[i + 56] = board[i]  
        board[i] += BLACK  
        board[i + 8] = PAWN + BLACK  
        board[i + 48] = PAWN  
    })  
}  
negamax #(top Int, depth Int, black Int, alpha Int, beta Int, Int) {  
    best : -ints.MAX  
    best_move : 0  
    depth <= 0 ? {  
        best = evaluate_board(black)  
        best >= beta || { depth < -1 } ? {  
            return best  
        }  
    }  
    1.to(2).each({  
       phase Int  
        depth <= 0 && {phase == 1} ? {  
            return best  
        }  
       1.to(64).each({  
            board[i] == 0 || {is_black(board[i]) != black} ? {  
                continue  
            }  
            attack_only : 0  
            phase == 0 ? {  
                attack_only = 1  
            }  
            moves : get_possible_moves(i, attack_only)  
            while({ moves != 0 }, {  
                target : ints.number_of_trailing_zeros(moves)  
                moves ^= 1 << target  
                capture : 0  
                board[target] != 0 ? {  
                    capture = 1  
                }  
                capture != attack_only ? {  
                  continue  
                }  
                move : move(i, target)  
                score : -negamax(0, depth - 1, 1 - black, -beta, -alpha)  
               
                score > best ? {  
                    best_move = move  
                    best = score  
                    alpha = ints.max(alpha, score)  
                }  
                undo(move)  
             top == 0 && { best >= beta } ? {   
                break  
             }  
            })  
       })  
    })  
    top != 0 ? {  
        best_move  
    }  
    best  
}  
  
is_field_attacked #(black Int, pos Int, Int) {  
    1.to(64).each({  
       i Int    
        b : board[i]  
        b == 0 || {is_black(b) == black} ? {  
          continue  
        }  
        moves : get_possible_moves(i, 1)  
        ((moves >> pos) & 1) == 1 ? {  
            return 1  
        }  
    })  
    0  
}   
evaluate_board #(black Int, Int) {  
    sum : 0  
    1.to(64).each({  
       i Int    
        b : board[i]  
        sc : 0  
        p : get_piece(b)  
        p == KING   ? { sc = 10000000 }  
        p == QUEEN  ? { sc = 1000 }  
        p == ROOK   ? { sc = 500 }  
        p == KNIGHT ? {  
            sc = ints.bit_count(get_possible_moves(i, 0))  
            sc += 320  
        }  
        p == BISHOP ? {  
            sc = ints.bit_count(get_possible_moves(i, 0))  
            sc += 330  
        }  
        p == PAWN ? {  
            sc = 100  
            turn > 40 ? {  
                sc += 1 << is_black(b) ?  { (i / 8) - 1} ,{ 6 - (i / 8) }  
            }  
        }  
         
        is_black(b) != black ? {  
            sc = -sc  
        }  
        sum += sc  
      })  
    sum  
}  
  
get_possible_moves #(from Int, attacks_only Int, Int) {  
    b : board[from]  
    p : get_piece(b)  
    black : is_black(b)  
    max_dist : 1  
    result : 0  
    p == QUEEN || { p == ROOK } || { p == BISHOP } ? { max_dist = 7 }  
      
    p == KING || { p == QUEEN } || { p == BISHOP } ? {  
        result |= slide(from, max_dist,  1,  1)  
        result |= slide(from, max_dist, -1, -1)  
        result |= slide(from, max_dist,  1, -1)  
        result |= slide(from, max_dist, -1,  1)  
    }  
      
    p == KING || { p == QUEEN } || { p == ROOK } ? {  
        result |= slide(from, max_dist, 1, 0)  
        result |= slide(from, max_dist, -1, 0)  
        result |= slide(from, max_dist, 0, 1)  
        result |= slide(from, max_dist, 0, -1)  
    }  
    p == KNIGHT ? {  
        result |= slide(from, max_dist, 1, 2)  
        result |= slide(from, max_dist, 2, 1)  
        result |= slide(from, max_dist, -1, -2)  
        result |= slide(from, max_dist, -2, -1)  
        result |= slide(from, max_dist, 1, -2)  
        result |= slide(from, max_dist, 2, -1)  
        result |= slide(from, max_dist, -1, 2)  
        result |= slide(from, max_dist, -2, 1)  
    }  
    p == PAWN ? {  
        dir : is_black(b)? {1},{-1}  
        // straight  
        dist : 1  
        attacks_only == false ? {  
            
            dir == 1 && { from / 8 == 1 } && { board[from + 16] == 0 } ? {  
                dist = 2  
            }  
            dir == -1 && { from / 8 == 6 } && { board[from - 16] == 0 } ? {  
                dist = 2  
            }  
            result |= slide(from, dist, 0, dir)  
            result != 0 && { board[from + dir * 8] != 0 } ? {  
                result = 0  
            }  
        }  
        pawn_moved_2 / 8 == from / 8 ? {  
            // en passant  
            pawn_moved_2 == from - 1 ? {  
                result |= 1 << (from - 1)  
            }, {  
             pawn_moved_2 == from + 1 ? {  
                result |= 1 << (from + 1)  
             }  
           }  
        }  
        capture : slide(from, max_dist, 1, dir)  
        capture != 0 && { board[from + 1 + dir * 8] != 0 } ? {  
            result |= capture  
        }  
        capture = slide(from, max_dist, -1, dir)  
        capture != 0 && { board[from - 1 + dir * 8] != 0 } ? {  
            result |= capture  
        }  
    }  
    attacks_only == 0 && p == KING ? {  
       // castling  
       r : castling_flags >> black? {0}, {2}  
       (r & 3) != 3 && { is_field_attacked(black, from) == false } ? {  
          (r & 1) == 0 && { is_field_attacked(black, from - 1) == false } ? {  
             rook : slide(from - 4, 8, 1, 0)  
             rook >>= from - 3  
             (rook & 7) == 7 ? {  
                result |= 1 << (from - 2)  
             }  
          }  
          (r & 2) == 0 && { is_field_attacked(black, from + 1) == false } ? {  
             rook : slide(from + 3, 8, -1, 0)  
             rook >>= from + 1  
             (rook & 3) == 3 ? {  
                result |= 1 << (from + 2)  
             }  
          }  
       }  
    }  
    result  
}  
  
slide #(from Int, max_dist Int, xo Int, yo Int, Int) {  
    x : from & 7  
    y : from / 8  
    is_black : is_black(board[from])  
    result : 0  
    i : 1  
    loop {  
        x += xo  
        y += yo  
        i > max_dist || { x < 0 } || { x > 7 } || { y < 0 } || { y > 7 } ? {  
          break  
        }  
        p : x + y * 8  
        b : board[p]  
        b != 0 ? {  
          is_black(b) != is_black ? {  
            result |= 1 << p  
          }  
          break  
        }  
        result |= 1 << p  
        i += 1  
    }  
    result  
}  
  
// translation line --------------  
update_castling_rights #(pos Int) {  
    p : board[pos]  
    get_piece(p) == ROOK ? {  
        (pos & 7) == 0 || { (pos & 7) == 7 } ? {  
            which : (pos & 7) == 0? {1}, {2}  
            castling_flags |= which << is_black(p)? {0}, {2}  
        }  
    }  
}  
  
move #(move Int, Int) {  
    source : (move >> 16) & 0xff  
    target : (move >> 8) & 0xff  
    move(source, target)  
}  
  
move #(source Int, target Int, Int) {  
    turn += 1  
    captured : board[target]  
    old_castling_flags : castling_flags  
    update_castling_rights(source)  
    update_castling_rights(target)  
    old : board[source]  
    board[target] = old  
    p : get_piece(old)  
    is_black : is_black(old)  
    board[source] = 0  
    old_pawn_moved : pawn_moved_2  
    pawn_moved_2 = 0  
    p == PAWN ? {  
        shift : (target & 7) - (source & 7)  
        shift != 0 && { captured == 0 } ?{  
            // en passant capture  
            board[source + shift] = 0  
        }  
        target <= 7 || { target >= 56 } ? {  
            // promotion  
            board[target] = QUEEN + is_black ? { BLACK }  , { 0 }  
        }  
        ints.abs(source - target) == 16 ? {  
            pawn_moved_2 = target  
        }  
    }  
            
    p == KING ? {  
        castling_flags |= 3 << is_black ? {0}, {2}  
        ints.abs((source & 7) - (target & 7)) > 1 ? {  
            // castling  
            target > source ? {  
                board[target - 1] = board[target + 1]  
                board[target + 1] = 0  
            }, {  
                board[target + 1] = board[target - 2]  
                board[target - 2] = 0  
            }  
        }  
    }  
  
    (old << 40) | (old_pawn_moved << 32) | (captured << 24) | (source << 16) | (target << 8) | old_castling_flags  
  
}   
undo #(move Int) {  
    turn -= 1  
    old : (move >> 40) & 0xff  
    pawn_moved_2 = (move >> 32) & 0xff  
    captured : (move >> 24) & 0xff  
    source : (move >> 16) & 0xff  
    target : (move >> 8) & 0xff  
    castling_flags = move & 0xff  
    board[target] = captured  
    board[source] = old  
    get_piece(old) == KING ? {  
        ints.abs((source & 7) - (target & 7)) > 1 ? {  
            // undo castling  
            target > source ? {  
                board[target + 1] = board[target - 1]  
                board[target - 1] = 0  
            }, {  
                board[target - 2] = board[target + 1]  
                board[target + 1] = 0  
          }  
        }  
    }  
    get_piece(old) == PAWN ? {  
        shift : (target & 7) - (source & 7)  
        shift != 0 && { captured == 0 } ? {  
            // en passant capture  
            board[source + shift] = PAWN + is_black(old) ? {0}, {BLACK}  
        }  
    }  
}  
is_black #(p Int, Bool) {  
    p > 6  
}  
  
get_piece #(p Int, Bool) {  
    p > 6 ? { p - 6 },{ p }  
}
```
