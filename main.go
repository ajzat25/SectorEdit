package main

// fix:
//    frame rate dependent friction
//    walking on slopes
//    auto step
//    multi-texture effiecency
//    view bobbing

import (
  "ajzat.tk/sectoredit/grm3"
  "ajzat.tk/sectoredit/data"
  "ajzat.tk/sectoredit/pathfinder"
  "ajzat.tk/sectoredit/keymapper"

  "github.com/go-gl/gl/v2.1/gl"
  "github.com/go-gl/glfw/v3.2/glfw"

  "fmt"
  "log"
  "os"
  "time"
  "image"
  "image/draw"
  _ "image/png"
  "runtime"
  "math"
)

var (  // user settings
  screenWidth  = 1280
  screenHeight = 800
  fullscreen = false
  fov = 90
  look_speed = grm3.Vec2{0.5,0.5}
  camera_strafe_angle = 0.75
  camera_strafe_loss = 0.6
  strafing_tilt_max = 5.0
)

var (  // game settings
  vlookClamp = 90.0
  move_speed = 16.0
  accel_rate = 260.0
  run_speed = 20.0
  run_accel_rate = 280.0
  set_map_friction = 0.15
  gravity = grm3.Vec3{0,0,-42}
  jump = 14.0
  jetpack = 1.0
  air_control = 0.1
)

var ( // engine settings
  nearPlaneDistance = 0.1
  portalCullingDistance = -0.1
  glNear = 0.1
  glFar = 100
  max_frame_sectors = 32 // prevents infinite draw loops
  floorTest = grm3.Vec3{0,0,-0.25}
  jump_frames_reset = 0.15
  max_physics = 32
)

var ( // debug options
  frame_debug = false
  no_clip = false
  use_depth_test = true
  clip_walls = true
  wireframe = false
)

var (
  version = "18"
  makev = "0.7a"

  initialized = false

  grMap grdata.Map

  player pather.Player
  camera grm3.Tup3

  nPlane grm3.Plane
  pPlane grm3.Plane

  triangle_stack grdata.TriangleStack
  texture_index [][]int
  model_stack grdata.ModelStack
  sprite_stack []grdata.Sprite
  draw_sprite []int

  mslast grm3.Vec2
  deltaTime float64
  win *glfw.Window
  sectors_drawn int

  cameraSector int
  jump_frames float64
  standing = true

  aspect_ratio float64

  texture []uint32
  texturepath []string

  friction = 1.0-float64(set_map_friction)
)

func Start() *glfw.Window{

  var window *glfw.Window

  fmt.Println("")
  fmt.Println("SectorEdit"+version+" v"+makev+" By Andrew Johnson, 2018")
  if use_depth_test == false { fmt.Println("  -flag: no depth test")}
  if frame_debug == true { fmt.Println("  -flag: frame debugging")}
  if no_clip == true { fmt.Println("  -flag: no cliping")}
  if fullscreen == true { fmt.Println("  -flag: fullscreen")}
  fmt.Println()

  // start glfw
  if err := glfw.Init(); err != nil {;log.Fatalln("failed to initialize glfw:", err);}

  // dont let the window be bigger than the display
  mon := glfw.GetPrimaryMonitor();
  mode := mon.GetVideoMode();
  screenWidth = grm3.MinInt(mode.Width, screenWidth)
  screenHeight = grm3.MinInt(mode.Height, screenHeight)
  aspect_ratio = float64(screenWidth)/float64(screenHeight)

  // window settings
  glfw.WindowHint(glfw.Resizable, glfw.False)
  glfw.WindowHint(glfw.ContextVersionMajor, 2)
  glfw.WindowHint(glfw.ContextVersionMinor, 1)
  glfw.WindowHint(glfw.DoubleBuffer, glfw.True)

  // get a window
  window, err := glfw.CreateWindow(screenWidth, screenHeight, "SectorEdit"+version+" v"+makev, nil, nil)
  if err != nil {;panic(err);}
  window.MakeContextCurrent()

  if err := gl.Init(); err != nil {;panic(err);}
  // set up for mouse input
  window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
  window.SetInputMode(glfw.StickyKeysMode, glfw.False)
  // get first mouse position
  msx, msy := window.GetCursorPos()
  mslast = grm3.Vec2{msx,msy}
  initialized = true
  return window
}

func Load(path string){
  if !initialized { panic("SectorEdit is not initialized") }
  fmt.Println("SectorEdit Loading...")
  fmt.Println("   User Settings")
  // load keymap
  grkey.MapKey(grkey.WalkForward, glfw.KeyW)
  grkey.MapKey(grkey.WalkBackward, glfw.KeyS)
  grkey.MapKey(grkey.StrafeRight, glfw.KeyD)
  grkey.MapKey(grkey.StrafeLeft, glfw.KeyA)
  grkey.MapKey(grkey.Jump, glfw.KeySpace)
  grkey.MapKey(grkey.Exit, glfw.KeyQ)
  grkey.MapKey(grkey.Run, glfw.KeyLeftShift)
  grkey.MapKey(grkey.Jetpack, glfw.KeyE)
  fmt.Println("   Built in Data")
  player.Position = grm3.Vec3{6,-1,4}
  player.Direction = grm3.Vec3{45,0,0}
  player.HitBox = grm3.Make_Box3(1,1,3.5,grm3.Vec3{0,0,1.25})
  player.Sector = 0
  pather.SetMaxPhysTicks(max_physics)
  fmt.Println("   SectorEdit Level Package")
  fmt.Println("    │")
  fmt.Println("    ├─┬─[Map]")
  fmt.Println("    │ ├─[Sector Mesh]")
  defer recover()
  grMap.SM, texturepath = grdata.ReadMap(grdata.OpenFormat(path+"/sm.sea"))
  fmt.Println("    │ └─[Data]")
  model_stack = append(model_stack, grdata.Model{VA: grdata.VertexArray{
    grm3.Vec3{-1,-1,-1},
    grm3.Vec3{-1,-1,1},
    grm3.Vec3{1,-1,1},
    grm3.Vec3{1,-1,-1},
    grm3.Vec3{-1,1,-1},
    grm3.Vec3{-1,1,1},
    grm3.Vec3{1,1,1},
    grm3.Vec3{1,1,-1},
    }, IA: grdata.IndexArray{
    [3]grdata.MappedIndex{grdata.MappedIndex{Index: 2, Tex: grm3.Vec2{0,0}}, grdata.MappedIndex{Index: 1, Tex: grm3.Vec2{0,1}}, grdata.MappedIndex{Index: 0, Tex: grm3.Vec2{1,1}}},
    [3]grdata.MappedIndex{grdata.MappedIndex{Index: 0, Tex: grm3.Vec2{1,1}}, grdata.MappedIndex{Index: 3, Tex: grm3.Vec2{1,0}}, grdata.MappedIndex{Index: 2, Tex: grm3.Vec2{0,0}}},
    [3]grdata.MappedIndex{grdata.MappedIndex{Index: 4, Tex: grm3.Vec2{1,1}}, grdata.MappedIndex{Index: 5, Tex: grm3.Vec2{1,0}}, grdata.MappedIndex{Index: 6, Tex: grm3.Vec2{0,0}}},
    [3]grdata.MappedIndex{grdata.MappedIndex{Index: 6, Tex: grm3.Vec2{1,1}}, grdata.MappedIndex{Index: 7, Tex: grm3.Vec2{1,0}}, grdata.MappedIndex{Index: 4, Tex: grm3.Vec2{0,0}}},
    }, TA: []int{2,2,2,2,2,2}} )
  sprite_stack = append(sprite_stack, grdata.Sprite{ ModelIndex: 0, Position: grm3.Vec3{4,0,4}, Direction: grm3.Vec3{0,0,180}, Scale: 1, Sector: 0, Sectors: []int{0}, Static: false, Collider: grm3.Make_Box3(2,2,2, grm3.Vec3{0,0,0}), Velocity: grm3.Vec3{-25,25,10} } )
  for i:=0; i<len(sprite_stack); i++{
    sprite_stack[i].SpriteUpdateSectorRegistration(&grMap.SM.SM, i, sprite_stack[i].Sectors)
  }
  fmt.Println("    │")
  fmt.Println("    ├─[Textures]")
  texture = make([]uint32, len(texturepath))
  for i:=0;i<len(texturepath);i++{ texture[i] = loadTexture(texturepath[i]) }
  fmt.Println("    ├─[Scripts]")
  fmt.Println("    ├─[Audio]")
  fmt.Println("    │")
  fmt.Println("   Done")
  fmt.Println()
  // os.Exit(0)
}

func RunLoop(win *glfw.Window){
  if !initialized { panic("SectorEdit is not initialized") }; lasttime := time.Now()
  for !win.ShouldClose() {
    // need to update events at a faster rate
    for i:=0; i <= 4; i++{ glfw.PollEvents() }
    // player input & physics
    movement(win)
    physics()
    // build the triangle stack
    buildFrame()
    drawOBJ()
    // draw the triangle stack
    totalTri := glFrame()
		win.SwapBuffers()
    if frame_debug {
      fmt.Println("Finsihed Drawing", totalTri, "Triangles", "in", sectors_drawn, "Sectors @", int(1/deltaTime), "FPS"); fmt.Println("Sector:", player.Sector)
      if totalTri == 0 { fmt.Println("Oh No! Looks Like there was a issue... <empty tristack>")} // ; panic("Empty Triangle Stack") } // quit if there are no triangles to draw
      fmt.Println(""); fmt.Println("--New Frame")
    }
    deltaTime, lasttime = time.Since(lasttime).Seconds(), time.Now()
    // sprite_stack[0].Direction[grm3.Yaw] += 180*deltaTime
	}
}

func buildFrame(){
  // reset
  triangle_stack = grdata.TriangleStack{}
  texture_index = make([][]int, len(texture))
  draw_sprite = []int{}
  sectors_drawn = 0
  // get the position, and direction of the camera from the position of the player
  camera = player.GetCamera()
  // calculate a new near plane reletive to the player
  nPlane = grm3.Make_NearPlane(camera, nearPlaneDistance)
  pPlane = grm3.Make_NearPlane(camera, portalCullingDistance)
  // start drawing from the sector containing the player, a clipping space is also created from the players position, and direction
  SectorShader(player.Sector, -1, grm3.Make_Init_ClippingSpace(fov, screenWidth, screenHeight, camera[grm3.Position], camera[grm3.Direction]))
}

func SectorShader(sectorID, fromPoral int, cls grm3.ClippingSpace) bool{ /// returns false if render was sucessfull
  // iterate sectors draw, and fail if the limit is passed
  sectors_drawn++
  if sectors_drawn >= max_frame_sectors{ ; return false  ; }
  // draw all walls in sector
  for i := 0; i < len(grMap.SM.SM[sectorID].Faces); i++{ ; WallSub(grMap.SM.SM[sectorID].Faces[i], cls) ; }
  // run portal shader
  for i := 0; i < len(grMap.SM.SM[sectorID].Portals); i++{
    if i != fromPoral{ // skip portal that sector was drawn from
      PortalSub(sectorID, i, cls)
    }
  }
  for i:=0; i<len(grMap.SM.SM[sectorID].Sprites); i++ {
    grdata.AddSprite(&draw_sprite, grMap.SM.SM[sectorID].Sprites[i])
  }
  return true
}

func PortalSub(sectorID, portalID int, cls grm3.ClippingSpace){
  // pack the face into a polygon, and clip it to the near plane, then clip it to the parrent clipping space
  port_clipped := grMap.SM.SM[sectorID].Portals[portalID].Points.Pack().Close().Poly_Clip_Plane(pPlane).Clip_ClippingSpace(cls)
  // stop if there is no intersection
  if len(port_clipped) >= 3{
    // create a clipping space from the clipped portal
    clso := grm3.Make_ClippingSpace(port_clipped, camera[grm3.Position])
    // call the sector shader with the new clipping space
    SectorShader(grMap.SM.SM[sectorID].Portals[portalID].RefranceSector, grMap.SM.SM[sectorID].Portals[portalID].RefrancePortal, clso)
  }
}

func WallSub(face grdata.Face, cls grm3.ClippingSpace){
  if !face.Render {return}
  var ClippedPoly grm3.Poly3
  // clip the wall to the parrent clipping space
  if clip_walls{ ClippedPoly = face.Points.Pack().Close().Poly_Clip_Plane(nPlane).Clip_ClippingSpace(cls) } else { ClippedPoly = face.Points.Pack() }
  // if there isn't at least one visable triangle, then give up
  if len(ClippedPoly) < 3{return}
  // synthesize and inject texture coordinates
  mapper := ClippedPoly.InjectTexture(face.Normal, face.Shift, face.Udot, face.Vdot, face.TextureShift, face.TextureLinear, face.TextureRotate)
  // Create TriangleChain and add it to the stacked chain
  triangle_stack.Push(mapper, &texture_index, face.Texture)
}

func drawOBJ() {
  for i:=0; i<len(draw_sprite); i++ {
    triangle_stack.PushSprite(sprite_stack[draw_sprite[i]], &model_stack, &texture_index)
  }
}

func movement(win *glfw.Window){
  var walkdir grm3.Vec3
  var speed float64
  var accel float64

  // mouse look
  msx, msy := win.GetCursorPos()
  ms := grm3.Vec2{msx,msy}
  msd := ms.SUB(mslast).MUL(look_speed)
  mslast = ms
  player.Direction[grm3.Yaw] -= msd[grm3.X]
  player.Direction[grm3.Pitch] -= msd[grm3.Y]

  // clamp vertial look
  player.Direction[grm3.Pitch] = grm3.AbsClamp(player.Direction[grm3.Pitch], vlookClamp)
  // copy an euler angle from player direction (could add pitch or roll, but why)
  walkdir[grm3.Yaw] = player.Direction[grm3.Yaw]

  netMove := grm3.Vec2{0,0}
  run := false
  jumpb := false
  jetpackb := false
  if win.GetKey(grkey.KeyValue(grkey.Exit))         ==  glfw.Press{ os.Exit(0) }
  if win.GetKey(grkey.KeyValue(grkey.Run))          ==  glfw.Press{ run = true }
  if win.GetKey(grkey.KeyValue(grkey.Jump))         ==  glfw.Press{ jumpb = true }
  if win.GetKey(grkey.KeyValue(grkey.WalkForward))  ==  glfw.Press{ netMove[grm3.Y] += 1 }
  if win.GetKey(grkey.KeyValue(grkey.WalkBackward)) ==  glfw.Press{ netMove[grm3.Y] -= 1 }
  if win.GetKey(grkey.KeyValue(grkey.StrafeLeft))   ==  glfw.Press{ netMove[grm3.X] += -1 }
  if win.GetKey(grkey.KeyValue(grkey.StrafeRight))  ==  glfw.Press{ netMove[grm3.X] += 1 }
  if win.GetKey(grkey.KeyValue(grkey.Jetpack))      ==  glfw.Press{ jetpackb = true }

  // running yealds different acceleration rate
  switch run && standing {case false: accel = accel_rate; case true: accel = run_accel_rate}
  // only walk if the amount of net movement is greater than 0
  if !standing { accel = accel*air_control }
  if netMove.Magnitude() > 0{
    // find the direction of the net movement in degrees
    walkdir[grm3.Yaw] += grm3.Degrees(math.Atan2(netMove[grm3.Y],netMove[grm3.X]))-90
    // move the player by a rotated idenity vector, scaled the product of deltaTime and acceleration rate
    // scaling by deltaTime insures a consistant movement speed at any frame rate
    player.Velocity = player.Velocity.ADD(grm3.IdentityVector().EulerRotate(walkdir).SCALE(deltaTime*accel))
  }
  if jump_frames > 0 { jumpb = false; jump_frames+=deltaTime }
  if jump_frames > jump_frames_reset { jump_frames = 0 }
  if jumpb && standing { player.Velocity[grm3.Z] += jump; jump_frames+=deltaTime}

  // camera tilt with strafe
  player.Direction[grm3.Roll] = player.Direction[grm3.Roll]*camera_strafe_loss
  if netMove[grm3.X] != 0{ player.Direction[grm3.Roll] += camera_strafe_angle*netMove[grm3.X] }
  player.Direction[grm3.Roll] = grm3.AbsClamp(player.Direction[grm3.Roll], strafing_tilt_max)

  // break out the X,Y compnents of the Velocity vector of player
  move_component := player.Velocity.Vec2()
  // limit the players speed to move_speed
  switch run{case false: speed = move_speed; case true: speed = run_speed}
  if move_component.Magnitude() > speed{ move_component = move_component.UnitVector().SCALE(speed) }
  // combine the players Z velocity and new XY velocity, and replace player velocity with this
  player.Velocity = move_component.Vec3(player.Velocity[grm3.Z])
  if jetpackb { player.Velocity[grm3.Z] += jetpack }
}

func physics(){
  var box grm3.Box3ext
  // player <-> map physics
  // gravity by frame
  player.Velocity = player.Velocity.ADD(gravity.SCALE(deltaTime))
  // get the full boundingbox
  box = grm3.Make_Box3ext(player.HitBox)
  // detect standing
  standing = false
  for i:=0; i<4; i++{
    pt := box[i].Vec2().SCALE(0.9).Vec3(box[i][grm3.Z]).ADD(player.Position)
    if pather.TestRayTSectorWalls(&grMap.SM, grm3.Line3{pt, floorTest}, player.Sector) { standing = true; break }
  }
  // collision detection, and sector Sector Traversal
  player.Sector, player.Position, player.Velocity = pather.LargeBoxCollider(&grMap.SM, box, player.Position, player.Velocity, deltaTime, friction, player.Sector, 0)

  // sprite <-> map physics
  size := len(sprite_stack)
  for i:=0; i<size; i++{
    if !sprite_stack[i].Static{
      sprite_stack[i].Velocity = sprite_stack[i].Velocity.ADD(gravity.SCALE(deltaTime))
      box = grm3.Make_Box3ext(sprite_stack[i].Collider)
      sprite_stack[i].Sector, sprite_stack[i].Position, sprite_stack[i].Velocity = pather.LargeBoxCollider(&grMap.SM, box, sprite_stack[i].Position, sprite_stack[i].Velocity, deltaTime, friction, sprite_stack[i].Sector, 0)
      sprite_stack[i].SpriteUpdateSectorRegistration(&grMap.SM.SM, i, []int{sprite_stack[i].Sector})
    }
  }
}

// make use of textures more efficent
func glFrame() int{
  numTex := len(texture_index)
  total := 0

  // if depth testing is on clear the buffers
  if use_depth_test{;gl.Clear(gl.DEPTH_BUFFER_BIT);}
  // clear color buffer
  gl.Clear(gl.COLOR_BUFFER_BIT)

  // create the camera matrix
  gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
  gl.Rotatef(float32(camera[grm3.Direction][grm3.Roll]),    0, 0, 1)
  gl.Rotatef(float32(0-camera[grm3.Direction][grm3.Pitch]), 1, 0, 0)
  gl.Rotatef(float32(0-camera[grm3.Direction][grm3.Yaw]),   0, 1, 0)
  gl.Translatef(camera[grm3.Position].Invert().GLPositionSwap3())
  gl.Color4f(1, 1, 1, 1)

  // draw the triangles
  for i:=0; i<numTex; i++ { // i in range of all textures
    size := len(texture_index[i])
    if size > 0 { // if there is at least one triangle to draw
      total += size
      gl.BindTexture(gl.TEXTURE_2D, texture[i])
      gl.Begin(gl.TRIANGLES)
      for j:=0; j<size; j++{ // for each triangle, j, of texture i
        GLdraw(triangle_stack[texture_index[i][j]]) // draws a triangle
      }
      gl.End()
    };}
    return total
}

// draw a Tri3 useing OpenGL
func GLdraw(a grm3.Tri3){
  gl.TexCoord2f(a[0].Text.GLTexture()) ; gl.Vertex4f(a[0].Vert.GLPositionSwap4())
  gl.TexCoord2f(a[1].Text.GLTexture()) ; gl.Vertex4f(a[1].Vert.GLPositionSwap4())
  gl.TexCoord2f(a[2].Text.GLTexture()) ; gl.Vertex4f(a[2].Vert.GLPositionSwap4())
  // gl.Color3f(grm3.Vec3{0.3,0.3,0.3}.GLPosition3()) ; gl.TexCoord2f(a[0].Text.GLTexture()) ; gl.Vertex4f(a[0].Vert.GLPositionSwap4())
  // gl.Color3f(grm3.Vec3{0.6,0.6,0.6}.GLPosition3()) ; gl.TexCoord2f(a[1].Text.GLTexture()) ; gl.Vertex4f(a[1].Vert.GLPositionSwap4())
  // gl.Color3f(grm3.Vec3{1,1,1}.GLPosition3()) ; gl.TexCoord2f(a[2].Text.GLTexture()) ; gl.Vertex4f(a[2].Vert.GLPositionSwap4())
}

func GLsetup(){
  gl.ClearColor(0, 0, 0, 1.0)
  if use_depth_test{;gl.ClearDepth(1);gl.Enable(gl.DEPTH_TEST);}
  gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
  tmpMat := grm3.PerspectiveMat4(float32(grm3.Radians(float64(fov))), float32(aspect_ratio), float32(glNear), float32(glFar))
	gl.MultMatrixf(&tmpMat[0])
	gl.MatrixMode(gl.MODELVIEW)
  if wireframe { gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE); gl.LineWidth(12) } else { gl.PolygonMode(gl.FRONT, gl.FILL);  gl.Enable(gl.CULL_FACE);  gl.CullFace(gl.BACK) }
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func loadTexture(file string) uint32 {
  // open the file
	imgFile, err := os.Open(file)
	if err != nil {;log.Fatalf("texture %q not found on disk: %v\n", file, err);}
  // decode image
	img, _, err := image.Decode(imgFile)
	if err != nil {;panic(err);}
  // draw the file to an image
	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {;panic("unsupported stride");}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
  // make an opengl texture from the image
	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))
	return texture
}

func main(){
  // when the main function ends exit glfw
  defer glfw.Terminate()
  // initialize SectorEdit
  win = Start()
  Load("assets/map/dev/dev1") // load a SectorEdit Map Package
  GLsetup() // set up OpenGL
  RunLoop(win) // run the main loop
}
