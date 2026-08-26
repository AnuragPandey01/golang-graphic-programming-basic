package camera

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraDirection int

const (
	DirectionForward CameraDirection = iota
	DirectionBackward
	DirectionLeft
	DirectionRight
)

const (
	Yaw         float32 = -90.0
	Pitch       float32 = 0.0
	Speed       float32 = 2.5
	Sensitivity float32 = 0.1
	Fov         float32 = 45.0
)

type Camera struct {

	// camera Attributes
	position mgl32.Vec3
	front    mgl32.Vec3
	up       mgl32.Vec3
	right    mgl32.Vec3
	worldUp  mgl32.Vec3

	// euler Angles
	yaw   float32
	pitch float32

	// camera options
	movementSpeed    float32
	mouseSensitivity float32
	fov              float32
}

func GetNewCamera(cameraPosition, worldUp mgl32.Vec3, yaw, pitch float32) *Camera {
	camera := Camera{
		position: cameraPosition,
		worldUp: worldUp,
		yaw:yaw,
		pitch: pitch,
		movementSpeed: Speed,
		mouseSensitivity: Sensitivity,
		fov: Fov,
	}

	// cameraFront and cameraPitch calculated using yaw and pitch
	camera.updateCameraVectors()
	return &camera
}

func (cam *Camera) ViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV( cam.position, cam.position.Add(cam.front), cam.up)
}

func (cam *Camera) Fov() float32 {
	return cam.fov
}

func (cam *Camera) ProcessKeyboard(direction CameraDirection, deltaTime float32){
	
	velocity := cam.movementSpeed * deltaTime
	
	switch direction{
	case DirectionForward:
		cam.position = cam.position.Add(cam.front.Mul(velocity))
	case DirectionBackward:
		cam.position = cam.position.Sub(cam.front.Mul(velocity))
	case DirectionLeft:
		cam.position = cam.position.Sub(cam.right.Mul(velocity))
	case DirectionRight:
		cam.position = cam.position.Add(cam.right.Mul(velocity))
	}
}

func (cam *Camera) ProcessMouseMovement(xoffset, yoffset float32, constrainPitch bool) {
	xoffset *= cam.mouseSensitivity
	yoffset *= cam.mouseSensitivity

	cam.yaw += xoffset
	cam.pitch += yoffset

	if constrainPitch {
		if cam.pitch > 89.0 {
			cam.pitch = 89.0
		}
		if cam.pitch < -89.0 {
			cam.pitch = -89.0
		}
	}

	cam.updateCameraVectors()
}

func (cam *Camera) ProcessMouseScroll(yoffset float32) {
	cam.fov -= yoffset
	if cam.fov < 1.0 {
		cam.fov = 1.0
	}
	if cam.fov > 45 {
		cam.fov = 45
	}
}

func (cam *Camera) updateCameraVectors() {
	var front mgl32.Vec3 = mgl32.Vec3{0, 0, 0}
	front[0] = float32(math.Cos(float64(mgl32.DegToRad(cam.yaw))) * math.Cos(float64(mgl32.DegToRad(cam.pitch))))
	front[1] = float32(math.Sin(float64(mgl32.DegToRad(cam.pitch))))
	front[2] = float32(math.Sin(float64(mgl32.DegToRad(cam.yaw))) * math.Cos(float64(mgl32.DegToRad(cam.pitch))))

	cam.front = front.Normalize()

	cam.right = cam.front.Cross(cam.worldUp).Normalize()
	cam.up = cam.right.Cross(cam.front).Normalize()
}
