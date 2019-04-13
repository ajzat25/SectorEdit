package grkey

import "github.com/go-gl/glfw/v3.2/glfw"

const (
  WalkForward = 0
  WalkBackward = 1
  StrafeLeft = 2
  StrafeRight = 3
  Jump = 4
  Exit = 5
  Menu = 6
  Run = 7
  Jetpack = 8
)

type keys [10]glfw.Key

var (
  KeyMap keys
)

func MapKey(key int, value glfw.Key){
  KeyMap[key] = value
}

func KeyValue(key int) glfw.Key{
  return KeyMap[key]
}
