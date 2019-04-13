package pather

import "ajzat.tk/sectoredit/grm3"
import "ajzat.tk/sectoredit/data"
import "fmt"

var (
  max_ticks = 32
)

func SetMaxPhysTicks(a int){
  max_ticks = a
}

type NPC struct{
  Position grm3.Vec3
  Velocity grm3.Vec3
  Direction float64
  Sector int
  HitBox grm3.Box3
}

type Player struct{
  Position grm3.Vec3
  Velocity grm3.Vec3
  Direction grm3.Vec3
  Sector int
  NextSector int
  HitBox grm3.Box3
}

func (a Player) GetCamera() grm3.Tup3{
  return grm3.Tup3{a.Position, a.Direction}
}

func RayTSectorTraversal(grMap *grdata.SectorMesh, ray grm3.Line3, sectorID int, fromPortal int) int{
  for i := 0; i < len((*grMap).SM[sectorID].Portals); i++{ // for each portal in the starting sector
    // skip the portal that would move back to the starting sector, then test if the ray intersects the portal
    if i != fromPortal{ if ray.Test_Intersect_Polygon((*grMap).SM[sectorID].Portals[i].Points.Pack(), (*grMap).SM[sectorID].Portals[i].Space, (*grMap).SM[sectorID].Portals[i].Normal){// skip portal that sector was drawn from
      // if the ray intersects the portal update to that sector, and see if it crosses any other portals in that sector
      return RayTSectorTraversal(grMap, ray, (*grMap).SM[sectorID].Portals[i].RefranceSector, (*grMap).SM[sectorID].Portals[i].RefrancePortal)
    };}
  }
  // if the ray doesn't intersect any new portals, return the final sector
  return sectorID
}

func LargeBoxCollider(grMap *grdata.SectorMesh, pts grm3.Box3ext, position, velocity grm3.Vec3, deltaTime, friction float64, sectorID, click int) (int, grm3.Vec3, grm3.Vec3){
  var plane grm3.Plane
  newDelta := velocity
  if click > max_ticks {fmt.Println("Failed create collision response in", max_ticks, "physics ticks or less"); return sectorID, position, velocity.SCALE(0.25)}
  for p:=0; p<8; p++{
    ptSect := RayTSectorTraversal(grMap, grm3.Line3{position, pts[p].ADD(newDelta.SCALE(deltaTime))}, sectorID, -1)
    faceCount := len((*grMap).SM[ptSect].Faces)
    for i:=0; i<faceCount; i++{
      if (*grMap).SM[ptSect].Faces[i].Physics{
      plane[grm3.Normal] = (*grMap).SM[ptSect].Faces[i].Normal
      plane[grm3.Position] = *(*grMap).SM[ptSect].Faces[i].Points[0]
      if pts[p].ADD(position).ADD(newDelta.SCALE(deltaTime)).Test_PlaneSide2(plane){
        newDelta = NewVelocity(newDelta, plane[grm3.Normal])
        newDelta = newDelta.Vec2().SCALE(friction).Vec3(newDelta[grm3.Z])
        if newDelta.Magnitude() < 0.01 { return sectorID, position, grm3.Vec3{0,0,0} }
        return LargeBoxCollider(grMap, pts, position, newDelta, deltaTime, friction, sectorID, click+1)
      };}
    }
  }
  newSect := RayTSectorTraversal(grMap, grm3.Line3{position, velocity.SCALE(deltaTime)}, sectorID, -1)
  return newSect, position.ADD(newDelta.SCALE(deltaTime)), newDelta
}

func NewVelocity(velocity, normal grm3.Vec3) grm3.Vec3{
  perpLine := grm3.Line3{velocity, normal}
  plane := grm3.Plane{grm3.Vec3{0,0,0}, normal}
  return plane.Intersect_Line(perpLine)
}

func TestRayTSectorWalls(grMap *grdata.SectorMesh, ray grm3.Line3, sectorID int) bool{
  // if the ray crosses sectors
  ptSect := RayTSectorTraversal(grMap, ray, sectorID, -1)
  for i := 0; i < len((*grMap).SM[ptSect].Faces); i++{
    if ray.Test_Intersect_Polygon((*grMap).SM[ptSect].Faces[i].Points.Pack(), (*grMap).SM[ptSect].Faces[i].Space, (*grMap).SM[ptSect].Faces[i].Normal.Invert()){
      return true
    };}
  // if the ray doesn't intersect any faces
  return false
}

// // project the players velocity onto the plane
// perpLine := grm3.Line3{pts[i].ADD(newDelta), pla[grm3.Normal]}
// // update the velocity
// newDelta = pla.Intersect_Line(perpLine).SUB(pts[i])
// // update the box for new walls
// translate = position.ADD(newDelta.SCALE(deltaTime))
// pts = grm3.Make_Box3ext(box, translate)
// }

// impliment max tracks per frame
// impliment max sector steps
// impliment max distance
// impliment slow turning
// impliment ceiling rejection
// keep track of top x paths, and calculate total distance to pick path

// impiment generic movement physics, using a external function CheckMotion(motion line2, sector sector)

// walk to(self location, self sector, dest location, dest sector, step size, max distance, max sectors):

// moving object = you
// desination = dest

// get sectors of you, and dest. (keep track of these to know when route must be recalculated)

// keep track of sector jumps
// store sector jumps in tree
// found = false
// start at sector you{
//    Add Sector To Tree, with Portal I came from
//    Is sector = dest sector
//    if yes:
//      found = true
//      Exit Loop
//    else:
//      Jump to all sectors connected to current sector, and do the same
//}
// if !found:
//  return fail
// Start at Node of Dest Sector in the tree (save this when finding sector path)
// Copy back into a sector plan (Dest Sector)::(Sector, Portal, Coord),(Sector, Portal, Coord) ect.
// if len(plan) == 1:
//    if Dist to plan[0] <= step_size:
//      return
//    else:
//      Walk(Dest, step size)
// if Dist to plan[0] <= step_size:
//    Position = plan[0]
//    delete Position 1
//    update sector to plan[0]
//    walk(plan[0], step_size)
// else:
//    walk(plan[0], step_size)

// def walk(coord, step_size):
//
// Turn in direction of coord
// Move in direction step_size amt


// can shoot(Location, Destination, Path):
// define line location -> destiatnion
// for portal in path:
//    if !intersect (portal, path):
//        return no
// return yes
