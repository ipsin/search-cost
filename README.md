search-cost
===========
This package is used to explore the cost function for a "guessing game", 
where you are attempting to guess a non-negative integer between 
A and B inclusive.  The cost in this case is the sum of all numbers
guessed.

For example, if A=x,B=x+2, you would guess x+1 for a total cost of x+1.

The goal in this package will be to construct functions F(x,n) that 
determine the worst-case cost of searching A=x, B=x+n exactly, 
using min(), max() and linear functions of x.

F(x,n) = min(x + k + max(F(x,k-1), F(x+k+1, n-k-1)), **∀** ⌈n/2⌉ <= k < n)

Phrased intuitively, this is the cost of picking (x+k), plus the 
worst of the worst-case cost of the recursive picks on either size.  
Finding the k that minimizes the sum gives the value of F(x,n).

Here are some examples of F(x,n) for small n.  Notice how the regions where
the piecewise function gains more segments, before collapsing into a 
simpler solution.  This pattern continues for much larger values of n.

F(x,0) = 0

F(x,1) = x

F(x,2) = x+1

F(x,3) = 2x+2

F(x,4) = 2x+4

F(x,5) = 2x+6

F(x,6) = 2x+8 

F(x,7) = 2x+10 (1<=x<5)&44; 3x+6 (x>=5)

F(x,8) = 2x+12 (1<=x<4)&44; 3x+9 (x>=4)

F(x,9) = 2x+14 (1<=x<3)&44; 3x+12 (x>=3)

F(x,10) = 3x+15 

F(x,11) = 3x+18 (1<=x<5)\, 4x+14 (5<=x<9)\, 3x+22 (x>=9)

F(x,12) = 3x+21 (1<=x<4)\, 4x+18 (4<=x<8)\, 3x+25 (x>=8)

F(x,13) = 3x+24 (1<=x<3)\, 4x+22 (3<=x<7)\, 3x+28 (x>=7)

F(x,14) = 4x+26 (1<=x<6), 3x+31 (x>=6)

F(x,15) = 4x+30 (1<=x<5), 3x+34 (5<=x<21), 4x+14 (x>=21)

F(x,16) = 4x+34 (1<=x<4), 3x+37 (4<=x<20), 4x+18 (x>=20)

F(x,17) = 4x+38 (1<=x<3), 3x+40 (3<=x<19), 4x+22 (x>=19)

F(x,18) = 3x+43 (1<=x<18), 4x+26 (x>=18)

F(x,19) = 3x+46 (1<=x<13), 4x+34 (x>=13)

F(x,20) = 3x+49 (1<=x<12), 4x+38 (x>=12)

F(x,21) = 3x+52 (1<=x<11), 4x+42 (x>=11)

F(x,22) = 3x+55 (1<=x<10), 4x+46 (x>=10)

F(x,23) = 3x+58 (1<=x<9), 4x+50 (9<=x<21), 5x+30 (21<=x<37), 4x+66 (x>=37)

F(x,24) = 3x+61 (1<=x<8), 4x+54 (8<=x<20), 5x+35 (20<=x<36), 4x+70 (x>=36)

F(x,25) = 3x+64 (1<=x<7), 4x+58 (7<=x<19), 5x+40 (19<=x<35), 4x+74 (x>=35)

F(x,26) = 3x+67 (1<=x<6), 4x+62 (6<=x<18), 5x+45 (18<=x<34), 4x+78 (x>=34)

F(x,27) = 3x+70 (1<=x<5), 4x+66 (5<=x<13), 5x+54 (13<=x<29), 4x+82 (x>=29)

F(x,28) = 3x+73 (1<=x<4), 4x+70 (4<=x<12), 5x+59 (12<=x<28), 4x+86 (x>=28)

F(x,29) = 3x+76 (1<=x<3), 4x+74 (3<=x<11), 5x+64 (11<=x<27), 4x+90 (x>=27)

F(x,30) = 4x+78 (1<=x<10), 5x+69 (10<=x<26), 4x+94 (x>=26)

F(x,31) = 4x+82 (1<=x<9), 5x+74 (9<=x<21), 6x+54 (21<=x<23), 4x+98 (23<=x<69), 5x+30 (x>=69)
