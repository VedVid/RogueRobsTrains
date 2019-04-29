1. Description.

2. Controls

3. Source code

4. License

5. Contact







1. Description

"Rogue Robs Trains" is simple roguelike made in 7 days for 7DRL2019 game jam, hosted by Slash,

Jeff Lait, and Darren Grey on itch.io.

Game is made with use of RAWIG roguelike template (https://github.com/VedVid/RAWIG),

and uses BearLibTerminal (display-focused library made by Cfyz)

and Deferral-Square font (made by Brian Bruggeman).



In "Rogue Robs Trains" player is bandit that boarded train transport of gold. You must kill

everyone and stop train to loot chests.



The core of the game is ranged combat mechanics. The most important factor is distance from

target. There is range indicator for every equipped weapon, showed as three bars:

for short (<=10), medium (>10, <=20) and long (>20) distance.


Additionally, player has to remember that some weapons needs to be cocked before every shot,

and sometimes reloading takes much more time than expected...



Animation option is off by default as experimental feature that does not blend with

RAWIG architecture well.





2. Controls

Game supports 8-direction movement:



q w e        y k u        home  up  pgup          kp7 kp8 kp9

 \|/          \|/             \  |  /                \ | /

a-s-d   or   h-.-l   or   left-space-right   or   kp4-kp5-kp6

 /|\          /|\             /  |  \                / | \

z x c        b j n         end  down pgdn         kp1-kp2-kp3



Special actions:

1     - change weapon to rifle

2     - change weapon to revolver

3     - change weapon to melee

f     - target / fire / cock weapon

[tab] - switch target

r     - reload / uncock weapon

i     - inspect

g     - pick up

p     - pull lever



<SHIFT> + S - save game

<SHIFT> + Q - quit game



This control scheme is the default one for QWERTY keyboard layout. However, since release

0.0.4, RRT supports QWERTZ, AZERTY and (experimental) Dvorak keyboards. It is possible to

customize scheme as well.

To change controls related options, edit options_controls.cfg file.




3. Source code

Source code is available at https://github.com/VedVid/RogueRobsTrains under permissive

FreeBSD license.
 Note that code is very messy due to time pressure. It is the very reason

why this game is not very moddable, even if uses json files to store data.

But if you don't like colors, feel free to change them, it won't break anything :)





4. Game is provided free of charge, and source code is available under FreeBSD license.

However, RRT uses third-party libraries and assets that may be licensed differently - read

LICENSE and LICENSES-NOTICE to learn more.




5. You may find me on Discord (#5352).

I am roguelikes discord and /r/roguelikes subreddit (as VedVid) regular.

You may find me on twitter as well - @Ved_RL.