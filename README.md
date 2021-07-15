## ray_engine

#### A toy raycasting engine built with Go + Ebiten v2 2D library

Heavily based on [3DSage fantastic Youtube videos of a C/OpenGL raycasting engine](https://www.youtube.com/watch?v=gYRrGTC7GtA)

_Gun mode !_

![img](engine.gif) 

_2D map rendered in 3D_

![img](screenie.png)


#### Build & run

Build with Go 1.16 and Ebiten v2 on Linux

     go build
    ./ray_engine

#### Keymaps

Arrows keys or ZSQD (Azerty) : Move

'i' : debug info toogle

'f' : fullscreen toogle

'm' : Gun mode

'k' : quit

#### Features, todos and idea box

- [X] Port to Ebiten v2
- [ ] Proper Collisions
- [X] Scale map to 16x16
- [X] 2D minimap for gun mode
- [X] Add basic floor/ceiling - **Just a matching png for now**
- [ ] Binary textures
- [ ] Up/down parallax/Y-Shearing
- [ ] Ebiten Audio mp3 sound support
- [ ] Weapon swap/shield
- [ ] Proper ballistics - **Very very rough prototype, to improve**
- [ ] Cube destruction/basic enemies - **Very very rough prototype, to improve**

#### Known bugs

- [ ] Mouse support is very dodgy (blocks keyboard left/right movements)
- [ ] "scaling issue" -> cubes becomes rectangles from afar