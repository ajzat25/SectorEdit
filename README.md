# SectorEdit
SectorEdit is A CPU accelerated First Person 3D Graphics Engine in Golang.
## Tell Me More
SectorEdit uses optimisations similar to the techniques used some older game engines, like Doom, Quake, Build, and GoldSource. It works in OpenGL 2.1 immediate mode, and uses the fixed function pipeline. However, I am considering switching to the programmable pipeline. If you want to know more you may find the [wiki](https://github.com/ajzat25/SectorEdit/wiki) usefull.

### What does CPU acceleration mean?
It means that the gpu never see's a pixel of data (in the map) that won't end up visable on the screen, saving time futility rendering extra triangles.

# Screen Shots
![ScreenShot1](ScreenShot1.png)
![ScreenShot2](ScreenShot2.png)

### Why? VBO's, and VAO's have been around forever!
Yes, Im doing this for fun, not because it makes sense.

## Level Editor?
Not Yet. I'm making levels by hand untill I get around to that!

## How does it work?
Check out the [wiki](https://github.com/ajzat25/SectorEdit/wiki)

## Its Buggy!
Short: Yes. It is.

Long Version: Making a 3d game engine is lots of work, and ***I Work Alone*** so It will probably be pretty buggy at first. If you find a bug make sure its not already reported, and then let me know. I try to make fixing bugs my first priority, but it needs to be repeatable, well described, and include instructions on how to repeat the bug.

## Dependencies
* Grm3
* GO-GL 2.1
* GLFW 3.2
