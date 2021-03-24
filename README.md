# take4

This is a small proof of concept demonstrating how to implement a simple AI for a game
without providing the computer any guidance on how to play the game.

This is achieved by utilizing a Monte Carlo Tree Search (MCTS), which is a trivial algorithm
to implement. That said, the MCTS can provide competitive human like play behavior when run
at a sufficient depth.

The numSimulations constant at the top of the code dictates how many games to simulate before
making a move. Setting this to 1,000 seems to provide a nice performance to accuracy tradeoff,
but it can be fun to play around with different values.

It's worth noting that none of this code has been optimized, as the implementation is designed
to serve as a reference (optimized for simplicity over performance).

```
[P1] Enter a move (1-7): 6
1 2 3 4 5 6 7
_ _ _ _ _ _ _
_ _ _ _ _ _ _
_ _ X O _ _ _
_ _ O X O X _
_ _ O X X O _
_ X O X X O _
CPU chooses 3

Game over! P2 wins.
1 2 3 4 5 6 7
_ _ _ _ _ _ _
_ _ O _ _ _ _
_ _ X O _ _ _
_ _ O X O X _
_ _ O X X O _
_ X O X X O _
```
