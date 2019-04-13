package grm3

// GRM3(3D & 2D Math for Go) formerly a go math library that has since been incoperated into sectorEdit
// Tagged sections start with "_" for ease of searching

// Sections:
//  _Constants
//  _Types
//  _Conversions
//  _Standard Opperations
//  _Vector Specials
//  _Collider Tests
//  _Low Level Interactions
//  _Clipping
//  _Constructors
//  _Polygon operations
//  _Euler Rotations
//  _Texture Mapping
//  _Convert to openGL
//  _Other

// vectors are stored as arrays, here is the order:
// _X_Y_Z_W_
// |0|1|2|3|
//  | | | |
//  | | | |-- Vec4
//  | | |-----Vec4 Vec3
//  | |-------Vec4 Vec3 Vec2
//  |---------Vec4 Vec3 Vec2



import "math"

// _Constants
const (
  x = 0 // vector 2/3/4
  y = 1 // vector 2/3/4
  z = 2 // vector 3/4
  w = 3 // vector 4
  position = 0 // line & plane
  normal = 1 // plane
  delta = 1 // line
  radius = 1 // sphere
  direction = 1 // Tup(2/3/4)
  min = 0 // boxes
  max = 1 // boxes
  sin = 0 // SinCos
  cos = 1 // SinCos
  yaw = 0
  pitch = 1
  roll = 2
  pi = 3.1415926535897932385
  toRad = pi/180
  toDeg = 180/pi
)

const ( //public
  X = 0
  Y = 1
  Z = 2
  W = 3
  Sin = 0
  Cos = 1
  Position = 0
  Direction = 1
  Radius = 1
  Normal = 1
  Yaw = 0
  Pitch = 1
  Roll = 2
  BoxMin = 0
  BoxMax = 1
)

func IdentityVector() Vec3{ // should be a constant but i cant figure out how to do that
  return Vec3{0,1,0}
}

// _Types
type Vec2 [2]float64
type Vec3 [3]float64
type Vec4 [4]float64
type Line2 Tup2
type Line3 Tup3
type Plane Tup3
type Sphere struct{Position Vec3; Radius float64}
type Box2 Tup2
type Box3 Tup3
type Box3ext [8]Vec3
type Poly2 []Vec2
type Poly3 []Vec3
type Poly3r []*Vec3
type SinCos [2]float64
type Tup2 [2]Vec2
type Tup3 [2]Vec3
type Tup4 [2]Vec4
type ClippingSpace []Plane
type Mat4 [16]float32
type MappedPoly3 []MappedVert
type Tri3 [3]MappedVert
type MappedVert struct{ Vert Vec3; Text Vec2}

// _Conversions
func (a Vec2) Vec3(z float64) Vec3{
  return Vec3{a[0], a[1], z}
}

func (a Vec2) Vec4(z, w float64) Vec4{
  return Vec4{a[0], a[1], z, w}
}

func (a Vec3) Vec4(z, w float64) Vec4{
  return Vec4{a[0], a[1], a[2], w}
}

func (a Vec4) Vec3(z, w float64) Vec3{
  return Vec3{a[0], a[1], a[2]}
}

func (a Vec3) Vec2() Vec2{
  return Vec2{a[0], a[1]}
}

func (a Poly3r) Pack() Poly3{
  var out Poly3
  for i := 0; i<len(a); i++{
    out = append(out, *a[i])
  }
  return out
}

// _Standard Opperations
// Vec2
func (a Vec2) ADD(b Vec2) Vec2{
  return Vec2{a[0] + b[0], a[1] + b[1]}
}

func (a Vec2) SUB(b Vec2) Vec2{
  return Vec2{a[0] - b[0], a[1] - b[1]}
}

func (a Vec2) MUL(b Vec2) Vec2{
  return Vec2{a[0] * b[0], a[1] * b[1]}
}

func (a Vec2) DIV(b Vec2) Vec2{
  return Vec2{a[0] / b[0], a[1] / b[1]}
}

func (a Vec2) SCALE(b float64) Vec2{
  return Vec2{a[0] * b, a[1] * b}
}

func (a Vec2) SCALEDIV(b float64) Vec2{
  return Vec2{a[0] / b, a[1] / b}
}

func (a Vec2) QuickRot(by SinCos) Vec2{
  return Vec2{(a[x]*by[cos]) - (a[y]*by[sin]), (a[y]*by[cos]) + (a[x]*by[sin])}
}

func (a Vec2) QuickRot90() Vec2{
  return Vec2{-1*a[y], a[x]}
}

func (a Vec2) Invert() Vec2{
  return Vec2{0-a[x], 0-a[y]}
}

// Vec3
func (a Vec3) ADD(b Vec3) Vec3{
  return Vec3{a[0] + b[0], a[1] + b[1], a[2] + b[2]}
}

func (a Vec3) SUB(b Vec3) Vec3{
  return Vec3{a[0] - b[0], a[1] - b[1], a[2] - b[2]}
}

func (a Vec3) MUL(b Vec3) Vec3{
  return Vec3{a[0] * b[0], a[1] * b[1], a[2] * b[2]}
}

func (a Vec3) DIV(b Vec3) Vec3{
  return Vec3{a[0] / b[0], a[1] / b[1], a[2] / b[2]}
}

func (a Vec3) SCALE(b float64) Vec3{
  return Vec3{a[0] * b, a[1] * b, a[2] * b}
}

func (a Vec3) SCALEDIV(b float64) Vec3{
  return Vec3{a[0] / b, a[1] / b, a[2] / b}
}

func (a Vec3) Invert() Vec3{
  return Vec3{0-a[x], 0-a[y], 0-a[z]}
}

// Vec4
func (a Vec4) ADD(b Vec4) Vec4{
  return Vec4{a[0] + b[0], a[1] + b[1], a[2] + b[2], a[3] + b[3]}
}

func (a Vec4) SUB(b Vec4) Vec4{
  return Vec4{a[0] - b[0], a[1] - b[1], a[2] - b[2], a[3] - b[3]}
}

func (a Vec4) MUL(b Vec4) Vec4{
  return Vec4{a[0] * b[0], a[1] * b[1], a[2] * b[2], a[3] * b[3]}
}

func (a Vec4) DIV(b Vec4) Vec4{
  return Vec4{a[0] / b[0], a[1] / b[1], a[2] / b[2], a[3] / b[3]}
}

func (a Vec4) SCALE(b float64) Vec4{
  return Vec4{a[0] * b, a[1] * b, a[2] * b, a[3] * b}
}

func (a Vec4) SCALEDIV(b float64) Vec4{
  return Vec4{a[0] / b, a[1] / b, a[2] / b, a[3] * b}
}

// _Vector Specials
// Vec2
func (a Vec2) DOT(b Vec2) float64{
  return a[0]*b[0] + a[1]*b[1]
}

func (a Vec2) Magnitude() float64{
  return math.Sqrt((a[x] * a[x]) + (a[y] * a[y]))
}

func CROSSZ(a, b Vec2) float64{
  return (a[x]*b[y]) - (a[y]*b[x])
}

func (a Vec2) UnitVector() Vec2{
  return a.SCALEDIV(a.Magnitude())
}

func DIST2(a,b Vec2) float64{
  return (a.SUB(b).Magnitude())
}

// Vec3
func (a Vec3) DOT(b Vec3) float64{
  return a[0]*b[0] + a[1]*b[1] + a[2]*b[2]
}

func CROSS(a, b Vec3) Vec3{
  return Vec3{(a[y]*b[z]) - (a[z]*b[y]), (a[z]*b[x]) - (a[x]*b[z]), (a[x]*b[y]) - (a[y]*b[x])}
}

func (a Vec3) Magnitude() float64{
  return math.Sqrt((a[x] * a[x]) + (a[y] * a[y]) + (a[z] * a[z]))
}

func (a Vec3) UnitVector() Vec3{
  return a.SCALEDIV(a.Magnitude())
}

func DIST3(a,b Vec3) float64{
  return (a.SUB(b).Magnitude())
}

// Vec4
func (a Vec4) DOT(b Vec4) float64{
  return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] + a[3]*b[3]
}

func (a Vec4) Magnitude() float64{
  return math.Sqrt((a[x] * a[x]) + (a[y] * a[y]) + (a[z] * a[z]) + (a[w] * a[w]))
}

func (a Vec4) UnitVector() Vec4{
  return a.SCALEDIV(a.Magnitude())
}

// line3 operations
func (a Line3) AtT(t float64) Vec3{
  return a[position].ADD(a[delta].SCALE(t))
}

// --------------------------------------------------------------------------------------------------------------------------------------
// _More Complicated Math below ----------------------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------------------------------------------------

// _Collider Tests

// Sphere - Plane
func (sph Sphere) Test_Intersect_Plane(pla Plane) bool{
  return pla.Tof_Intersect_Line(Line3{sph.Position, pla[normal].UnitVector().SCALE(sph.Radius)}) < 1.000001
}

// RayT - Plane
func (ray Line3) Test_Intersect_Polygon(a Poly3, cls ClippingSpace, n Vec3) bool{
  var pla Plane
  pt := ray[Position].ADD(ray[Direction])
  pla[normal] = n
  pla[position] = a[1]
  if !pt.Test_PlaneSide2(pla){
    pti := pla.Intersect_Line(ray)
    if pti.Test_InClippingSpace(cls) { return true }
  }
  return false
}

// _Low Level Interactions
func (a Vec3) Test_PlaneSide2(b Plane) bool{
  return b[normal].DOT(b[position].SUB(a)) <= 0
}

func (a Vec3) Test_PlaneSide(b Plane) bool{
  return b[normal].DOT(b[position].SUB(a)) <= -0.00001
}

func (a Plane) Intersect_Line(b Line3) Vec3{
  return b.AtT(a.Tof_Intersect_Line(b))
}

func (a Plane) Test_RayIntersect(b Line3) bool{
  t := a.Tof_Intersect_Line(b)
  return (t >= 0 && t <= 1)
}

func (a Vec3) Test_InClippingSpace(cls ClippingSpace) bool{
  for i:=0; i<len(cls); i++{
    if !a.Test_PlaneSide2(cls[i]){ return false }
  }
  return true
}

func (plane Plane) Tof_Intersect_Line(line Line3) float64{ // finds the intersection in terms of parametric t on the line
  denominator := ( line[delta].DOT(plane[normal]) )
  if AproxZero(denominator){
    return 0
  }
  return (plane[position].SUB(line[position])).DOT(plane[normal])/denominator
}

// _Clipping
func (pol Poly3) Clip_ClippingSpace(spa ClippingSpace) Poly3{
  if len(pol) < 3{
    return pol
  }
  pol2 := append(pol, pol[0])
  for i := 0; i < len(spa); i++{
    pol2 = pol2.Poly_Clip_Plane(spa[i])
    if len(pol2) < 3 {
      return pol2
    }
    pol2 = append(pol2, pol2[0])
  }
  return pol2[:len(pol2)-1]
}

func (pol Poly3) Poly_Clip_Plane(pla Plane) Poly3{ // clip a polygon to plane, this function needs a closed polygon, and returns an open polygon
  var out Poly3
  var ain bool
  var bin bool
  if len(pol) < 1{
    return pol
  }
  size := len(pol)-1
  bin = pol[0].Test_PlaneSide(pla)
  for i := 0; i < size; i++{
    ain = bin
    bin = pol[i+1].Test_PlaneSide(pla)
    switch {
      case ain && bin:
        out = append(out, pol[i+1])
      case !ain && bin:
        out = append(out, pla.Intersect_Line(MakeLine3_Points(pol[i], pol[i+1])), pol[i+1])
      case ain && !bin:
        out = append(out, pla.Intersect_Line(MakeLine3_Points(pol[i], pol[i+1])))
    }
  }
  return out
}

// _Constructors
func Make_SinCos(rad float64) SinCos{
  return SinCos{math.Sin(rad), math.Cos(rad)}
}

func MakeLine2_Points(a,b Vec2) Line2{
  return Line2{a,b.SUB(a)}
}

func MakeLine3_Points(a,b Vec3) Line3{
  return Line3{a,b.SUB(a)}
}

func MakePlane_Points(a,b,c Vec3) (out Plane){
  out[normal] = CROSS(a.SUB(b), c.SUB(b)).UnitVector()
  out[position] = b
  return out
}

func Make_ClippingSpace(in Poly3, camera Vec3) ClippingSpace{ //may need to rearange points to fix normal direction
  var out ClippingSpace
  size := len(in)
  in2 := append(in, in[0])
  for i := 0; i < size; i++{
    out = append(out, MakePlane_Points(in2[i],camera,in2[i+1]))
  }
  return out
}

func Make_OrthClippingSpace(in Poly3, normal Vec3) ClippingSpace{ //may need to rearange points to fix normal direction
  var out ClippingSpace
  size := len(in)
  in2 := append(in, in[0])
  for i := 0; i < size; i++{
    pnormal := CROSS(in2[i+1].SUB(in2[i]), normal)
    out = append(out, Plane{in2[i].ADD(pnormal.SCALE(-0.005)), pnormal})
  }
  return out
}

func Make_Init_ClippingSpace(fov, width, height int, cameraPos, cameraRot Vec3) ClippingSpace{
  aspect := float64(width) / float64(height); rfov := Radians(float64(fov)/2); square := math.Sin(rfov)
  sx := square*aspect; sz := square; sy := math.Cos(rfov)
  return Make_ClippingSpace(Poly3{
  Vec3{-sx,sy,sz}.UnitVector().EulerRotate(cameraRot).UnitVector().ADD(cameraPos),
  Vec3{sx,sy,sz}.UnitVector().EulerRotate(cameraRot).UnitVector().ADD(cameraPos),
  Vec3{sx,sy,-sz}.UnitVector().EulerRotate(cameraRot).UnitVector().ADD(cameraPos),
  Vec3{-sx,sy,-sz}.UnitVector().EulerRotate(cameraRot).UnitVector().ADD(cameraPos)}, cameraPos)
}

func Make_NearPlane(player Tup3, dist float64) Plane{
  var out Plane
  out[normal] = IdentityVector().EulerRotate(player[direction])
  out[position] = player[position].ADD(out[normal].SCALE(dist))
  return out
}

func Make_Box3ext(a Box3) (out Box3ext){
  out[0] = Vec3{a[0][x], a[0][y], a[0][z]}
  out[1] = Vec3{a[1][x], a[0][y], a[0][z]}
  out[2] = Vec3{a[1][x], a[1][y], a[0][z]}
  out[3] = Vec3{a[0][x], a[1][y], a[0][z]}
  out[4] = Vec3{a[0][x], a[0][y], a[1][z]}
  out[5] = Vec3{a[1][x], a[0][y], a[1][z]}
  out[6] = Vec3{a[1][x], a[1][y], a[1][z]}
  out[7] = Vec3{a[0][x], a[1][y], a[1][z]}
  return out
}

// Makes a Box from length width height, and position
func Make_Box3(xs, ys, zs float64, move Vec3) Box3{
  var out Box3
  out[0] = Vec3{xs/-2, ys/-2, zs/-2}.SUB(move)
  out[1] = Vec3{xs/2, ys/2, zs/2}.SUB(move)
  return out
}

// _Polygon operations
func (a Poly3) Close() Poly3{
  if len(a) < 0{
    return a
  } else {
    return append(a, a[0])
  }
}

func (a Poly3) Open() Poly3{
  if len(a) < 0{
    return a
  } else {
    return a[:len(a)-1]
  }
}

// _Euler Rotations
func (a Vec3) EulerZ(angle float64) Vec3{
  var b Vec3
  ang := Make_SinCos(Radians(angle))
  b[x] = (a[x]*ang[cos]) - (a[y]*ang[sin])
  b[y] = (a[x]*ang[sin]) + (a[y]*ang[cos])
  b[z] = a[z]
  return b
}

func (a Vec3) EulerX(angle float64) Vec3{
  var b Vec3
  ang := Make_SinCos(Radians(angle))
  b[x] = a[x]
  b[y] = (a[y]*ang[cos]) - (a[z]*ang[sin])
  b[z] = (a[y]*ang[sin]) + (a[z]*ang[cos])
  return b
}

func (a Vec3) EulerY(angle float64) Vec3{
  var b Vec3
  ang := Make_SinCos(Radians(angle))
  b[x] = (a[z]*ang[sin]) + (a[x]*ang[cos])
  b[y] = a[y]
  b[z] = (a[z]*ang[cos]) - (a[x]*ang[sin])
  return b
}

func (a Vec3) EulerRotate(angle Vec3) Vec3{
  var s1 Vec3
  s1 = a
  s1 = s1.EulerY(angle[roll])
  s1 = s1.EulerX(angle[pitch])
  s1 = s1.EulerZ(angle[yaw])
  return s1
}

// _Texture Mapping
func (a Poly3) InjectTexture(norm, pos, Udot, Vdot Vec3, teShift, teLinear Vec2, teRotate SinCos) MappedPoly3{
  var out MappedPoly3
  for i:=0; i<len(a);i++{
    shift := a[i].SUB(pos)
    tex := Vec2{shift.DOT(Udot),shift.DOT(Vdot)}.ADD(teShift).MUL(teLinear).QuickRot(teRotate)
    out = append(out, MappedVert{Vert: a[i], Text: tex})
  }
  return out
}

// _Convert to openGL
func (a Vec2) GLTexture() (float32, float32){
  return float32(a[0]), float32(a[1])
}

func (a Vec3) GLPosition3() (float32, float32, float32){
  return float32(a[x]), float32(a[y]), float32(a[z])
}

func (a Vec3) GLPositionSwap4() (float32, float32, float32, float32){
  return float32(a[x]), float32(a[z]), float32(0-a[y]), 1
}

func (a Vec3) GLPositionSwap3() (float32, float32, float32){
  return float32(a[x]), float32(a[z]), float32(0-a[y])
}

func PerspectiveMat4(fovy, aspect, near, far float32) Mat4 {
	// fovy = (fovy * math.Pi) / 180.0 // convert from degrees to radians
	nmf, f := near-far, float32(1./math.Tan(float64(fovy)/2.0))

	return Mat4{float32(f / aspect), 0, 0, 0, 0, float32(f), 0, 0, 0, 0, float32((near + far) / nmf), -1, 0, 0, float32((2. * far * near) / nmf), 0}
}


// _Other
func AproxEq(a,b float64) bool{
  return math.Abs(a - b) < 0.000001
}

func AproxZero(a float64) bool{
  return math.Abs(a) < 0.000001
}

func AproxZero2(a,b float64) bool{
  return math.Abs(a) < b
}

func Radians(deg float64) float64{
  return toRad*deg
}

func Degrees(rad float64) float64{
  return toDeg*rad
}

func AbsClamp(value, clamp float64) float64{
  if value > clamp{
    return clamp
  }
  if value < -clamp{
    return -clamp
  }
  return value
}

func Min(a, b float64) float64{
    if a < b {
        return a
    }
    return b
}

func MinInt(a, b int) int{
    if a < b {
        return a
    }
    return b
}

func (a Vec3) StealAxisFrom(b Vec3, c int) Vec3{
  out := a
  out[c] = b[c]
  return out
}

func (a Vec2) StealAxisFrom(b Vec2, c int) Vec2{
  out := a
  out[c] = b[c]
  return out
}
