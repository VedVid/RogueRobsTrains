1. Description.
2. Controls
3. Source code
4. License
5. Contact



1. Description
"Rogue Robs Trains" is simple roguelike made in 7 days for 7DRL2019 game jam, hosted by Slash, Jeff Lait, and Darren Grey on itch.io.
Game is made with use of RAWIG roguelike template (https://github.com/VedVid/RAWIG), and uses BearLibTerminal (display-focused library made by Cfyz) and Deferral-Square font (made by Brian Bruggeman). 

In "Rogue Robs Trains" player is bandit that boarded train transport of gold. You must kill everyone and stop train.

The core of the game is ranged combat mechanics. Two important factors are distance from target (rifles are most efficient on long (>20 tiles) range, revolvers on medium (>10 and <= 20 tiles) range.
Additionally, player has to remember that some weapons needs to be cocked before every shot, and sometimes reloading takes much more time than expected...


2. Controls
Game supports 8-direction movement:

q w e        y k u        home  up  pgup          kp7 kp8 kp9
 \|/          \|/             \  |  /                \ | /
a-s-d   or   h-.-l   or   left-space-right   or   kp4-kp5-kp6
 /|\          /|\             /  |  \                / | \
z x c        b j n         end  down pgdn         kp1-kp2-kp3

Special actions:
1 - change weapon to rifle
2 - change weapon to revolver
3 - change weapon to melee
f - target / fire / cock weapon
r - reload / uncock weapon
i - inspect
g - pick up
p - pull lever


3. Source code
Source code is available https://github.com/VedVid/RogueRobsTrains under permissive FreeBSD license.
Note that code is much like spaghetti.

4. Game is provided free of charge, and source code is available under FreeBSD license. However, RRT uses third-party libraries and assets that may be licensed differently - read LICENSE and LICENSES-NOTICE to learn more.

5. You may find my on Discord (#5352). I am roguelikes discord and /r/roguelikes subreddit (as VedVid) regular. You may find me on twitter as well - @Ved_RL.