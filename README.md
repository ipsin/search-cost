search-cost
===========
This package is used to explore the cost function for a "guessing game", 
where you are attempting to guess a non-negative integer between 
A and B inclusive.  The cost in this case is the sum of all numbers
guessed.

For example, if A=x,B=x+2, you would guess 11 for a total cost of 11.

The goal in this package will be to construct functions F(x,n) that 
determine the worst-case cost of searching A=x, B=x+n exactly, 
using min(), max() and linear functions of x.

Here are some examples for small n.

F(x,0) = 0

F(x,1) = x

F(x,2) = x + 1

F(x,3) = 2x + 2
