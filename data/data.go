package grdata

import "ajzat.tk/sectoredit/grm3"
// import "fmt"
import "io/ioutil"
import "strconv"

type grPackage struct{
  Game Map
  Textures []string
  Audio []string
  Scripts []string
}

type Map struct{
  SM SectorMesh
}

type SectorMesh struct{
  SM []Sector
  PT []grm3.Vec3
  }

type Sector struct{
  Faces []Face
  Portals []Portal
  Sprites []int
}

type Face struct{
  Points grm3.Poly3r
  Texture int
  Space grm3.ClippingSpace // used to accelerate collidision detection
  Normal grm3.Vec3 // calculated at runtime
  Shift grm3.Vec3 // calculated at runtime
  Udot grm3.Vec3 // calculated at runtime
  Vdot grm3.Vec3 // calculated at runtime
  TextureShift grm3.Vec2
  TextureLinear grm3.Vec2
  TextureRotate grm3.SinCos
  Render bool
  Physics bool
}

type Portal struct{
  Points grm3.Poly3r
  Normal grm3.Vec3 // calculated at runtime
  Space grm3.ClippingSpace // used to accelerate collidision detection
  RefranceSector int
  RefrancePortal int
}

type TriangleStack []grm3.Tri3
type TextureStack []int
type ModelStack []Model

type GameSettings struct{
  vlook float64
  move float64
  accel float64
  run float64
  raccel float64
  friction float64
  gravity float64
  jump float64
  air float64
}

type VertexArray []grm3.Vec3
type MappedIndex struct{ Index int; Tex grm3.Vec2 }
type IndexArray [][3]MappedIndex
type Model struct{ VA VertexArray; IA IndexArray; TA []int}

type Sprite struct{
  ModelIndex int
  Position grm3.Vec3
  Direction grm3.Vec3
  Velocity grm3.Vec3
  Scale float64
  Sector int
  Sectors []int
  Static bool
  Collider grm3.Box3
}

func (ts *TriangleStack) PushSprite(sp Sprite, stack *ModelStack, texture_index *[][]int) {
  modelSize := len((*stack)[sp.ModelIndex].VA)
  var newVertexArray []grm3.Vec3
  for i:=0; i<modelSize; i++{
    newVertexArray = append(newVertexArray, (*stack)[sp.ModelIndex].VA[i].SCALE(sp.Scale).EulerRotate(sp.Direction).ADD(sp.Position))
  }
  modelSize = len((*stack)[sp.ModelIndex].IA)
  for i:=0; i<modelSize; i++{
    *ts = append(*ts, grm3.Tri3{
        grm3.MappedVert{newVertexArray[(*stack)[sp.ModelIndex].IA[i][0].Index], (*stack)[sp.ModelIndex].IA[i][0].Tex},
        grm3.MappedVert{newVertexArray[(*stack)[sp.ModelIndex].IA[i][1].Index], (*stack)[sp.ModelIndex].IA[i][1].Tex},
        grm3.MappedVert{newVertexArray[(*stack)[sp.ModelIndex].IA[i][2].Index], (*stack)[sp.ModelIndex].IA[i][2].Tex}})
    (*texture_index)[(*stack)[sp.ModelIndex].TA[i]] = append((*texture_index)[(*stack)[sp.ModelIndex].TA[i]], len(*ts)-1)
  }
}

// push a polygon to the triangle stack
func (ts *TriangleStack) Push(pol grm3.MappedPoly3, texture_index *[][]int, texture int){
  var next grm3.MappedPoly3
  if len(pol) < 3{
    return
  }
  next = append(next, pol[0])
  i := 0
  for ; i+2 < len(pol); i+=2{
    *ts = append(*ts, grm3.Tri3{pol[i],pol[i+1],pol[i+2]})
    (*texture_index)[texture] = append((*texture_index)[texture], len(*ts)-1)
    next = append(next, pol[i+2])
  }
  i += 1
  for ; i < len(pol); i++{
    next = append(next, pol[i])
  }
  ts.Push(next, texture_index, texture)
}

// loading from file:

// panics on error
func checkError(e error) {
    if e != nil {
        panic(e)
    }
}

// opens a file
func OpenFormat(filePath string) string{
  dat, err := ioutil.ReadFile(filePath)
  checkError(err)
  return string(dat)
}

// reads a file, and returns a map
func ReadMap(data string) (SectorMesh, []string){
  var ( lSM SectorMesh; lSect Sector; lFace = Face{Render: true, Physics: true}; lPortal Portal; letter string; form string; value []string; txpath []string)
  for i:=0; i<len(data);{; letter = string(data[i])
    if letter == ";"{
      form = ReadBlock(&i, data)
      value = []string{}
      for {
        next := ReadBlock(&i, data)
        if next == "!" { i -= 1; break }
        value = append(value, next)
      }; if setForm(form, value, &lSM, &lSect, &lFace, &lPortal, &txpath) {return lSM, txpath}
    }; i++; }
  return lSM, txpath
}

// reads a data block
func ReadBlock(i *int, data string) string{
  var (contents string; letter string)
  for ; letter != "[" && letter != "!"; *i++ { letter = string(data[*i]) }
  if letter == "!" {return letter}
  for letter2 := ""; letter2 != "]"; *i++ {
    letter = string(data[*i])
    letter2 = string(data[*i+1])
    contents = contents + letter
  }
  return contents
}

// interpret sea commands
func setForm(form string, value []string, lSM *SectorMesh, lSect *Sector, lFace *Face, lPortal *Portal, txpath *[]string) bool{
  switch form {
  case "apt": // append point
    (*lSM).PT = append((*lSM).PT, MakeStrArrayVec3(value))

  case "atx":
    (*txpath) = append(*txpath, value[0])

  case "awpt": // append indices
    (*lFace).Points = append((*lFace).Points, MakeStrArrayIndicesArray(value, lSM)...)

  case "wtex": // set wall texture
    (*lFace).Texture = MakeStrInt(value[0])

  case "wlin": // set texture linear
    (*lFace).TextureLinear = grm3.Vec2{1, 1}.DIV(MakeStrArrayVec2(value))

  case "wshf": // set texture shift
    (*lFace).TextureShift = MakeStrArrayVec2(value)

  case "wrot": // set texture rotation
    (*lFace).TextureRotate = grm3.Make_SinCos(grm3.Radians(180+MakeStrFloat(value[0])))

  case "norender":
    (*lFace).Render = false

  case "nophysics":
    (*lFace).Physics = false

  case "appt": // add portal indices
    (*lPortal).Points = append((*lFace).Points, MakeStrArrayIndicesArray(value, lSM)...)

  case "pref": // set reference sector and portal
    (*lPortal).RefranceSector = MakeStrInt(value[0])
    (*lPortal).RefrancePortal = MakeStrInt(value[1])

  case "push_wall":
    var right grm3.Vec3
    plane := grm3.MakePlane_Points(*(*lFace).Points[0],*(*lFace).Points[1],*(*lFace).Points[2])
    (*lFace).Normal = plane[grm3.Normal]
    (*lFace).Shift = plane[grm3.Position]
    if plane[grm3.Normal].Vec2().Magnitude() <= 0.00001{
      right = grm3.Vec3{0,1,0}
    } else {
      right = grm3.CROSS(grm3.Vec3{0,0,1}, plane[grm3.Normal]).UnitVector()
    }
    (*lFace).Udot = right
    (*lFace).Vdot = grm3.CROSS(plane[grm3.Normal], right)
    (*lFace).Space = grm3.Make_OrthClippingSpace((*lFace).Points.Pack(), plane[grm3.Normal])
    (*lSect).Faces = append((*lSect).Faces, *lFace)
    *lFace = Face{Render: true, Physics: true}

  case "push_portal":
    plane := grm3.MakePlane_Points(*(*lPortal).Points[0],*(*lPortal).Points[1],*(*lPortal).Points[2])
    (*lPortal).Normal = plane[grm3.Normal]
    (*lPortal).Space = grm3.Make_OrthClippingSpace((*lPortal).Points.Pack(), plane[grm3.Normal])
    (*lSect).Portals = append((*lSect).Portals, *lPortal)
    *lPortal = Portal{}

  case "push_sector":
    (*lSM).SM = append((*lSM).SM, *lSect)
    *lSect = Sector{}

  case "eof":
    return true
  }
  return false
}

// reading data from strings:
func MakeStrArrayVec3(value []string) [3]float64{
  var num [3]float64
  for i:=0; i<3; i++{
    val, err := strconv.ParseFloat(value[i], 64)
    checkError(err)
    num[i] = val
  }; return num
}

func MakeStrArrayVec2(value []string) [2]float64{
  var num [2]float64
  for i:=0; i<2; i++{
    val, err := strconv.ParseFloat(value[i], 64)
    checkError(err)
    num[i] = val
  }; return num
}

func MakeStrArrayIndicesArray(value []string, lSM *SectorMesh) []*grm3.Vec3{
  var num []*grm3.Vec3
  for i:=0; i<len(value); i++{
    val, err := strconv.ParseInt(value[i], 10, 64)
    checkError(err)
    num = append(num, &lSM.PT[val])
  }; return num
}

func MakeStrInt(str string) int{
  val, err := strconv.ParseInt(str, 10, 64)
  checkError(err)
  return int(val)
}

func MakeStrFloat(str string) float64{
  val, err := strconv.ParseFloat(str, 64)
  checkError(err)
  return val
}

func SearchInt(list *[]int, num int) bool{
  l := len(*list)
  for i:=0; i<l; i++{
    if (*list)[i] == num { return true }
  }
  return false
}

func FindInt(list *[]int, num int) (bool, int){
  l := len(*list)
  for i:=0; i<l; i++{
    if (*list)[i] == num { return true, i }
  }
  return false, 0
}

func AddSprite(drawlist *[]int, spriteID int) {
  if !SearchInt(drawlist, spriteID) { (*drawlist) = append((*drawlist), spriteID) }
}

func FindPopInt(list *[]int, num int) {
  test, ind := FindInt(list, num)
  if ind == len(*list) { *list = (*list)[:ind]; return }
  if test {*list = append((*list)[:ind], (*list)[ind+1:]...)}
}

func (sp *Sprite) SpriteUpdateSectorRegistration(sm *[]Sector, id int, insectors []int){
  for i:=0; i<len(sp.Sectors); i++{
    FindPopInt(&(*sm)[sp.Sectors[i]].Sprites, id)
  }
  for i:=0; i<len(insectors); i++{
    (*sm)[insectors[i]].Sprites = append((*sm)[insectors[i]].Sprites, id)
  }
  sp.Sectors = insectors
}
