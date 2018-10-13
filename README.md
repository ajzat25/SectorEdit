# SectorEdit
SectorEdit is A CPU accelerated First Person 3D Graphics Engine in Golang.
## Tell Me More
SectorEdit uses optimisations similar to the techniques used some older game engines, like Doom, Quake, Build, and GoldSource. It works in OpenGL 2.1 immediate mode, and uses the fixed function pipeline. However, I am considering switching to the programmable pipeline.

### What does CPU acceleration mean?
It means that the gpu never see's a pixel of data (in the map) that won't end up visable on the screen, saving time futility rendering extra triangles.

### Why? VBO's, and VAO's have been around forever!
Yes, Im doing this for fun, not because it makes sense.

## Level Editor?
Not Yet. I'm making levels by hand untill I get around to that!

## How does it work?
Check out this neat [wiki](https://github.com/ajzat25/SectorEdit/wiki)

### Dependencies
* Grm3
* GO-GL 2.1
* GLFW 3.2
