package tests

import (
	"log"
	"testing"
)

type ICamera interface {
	MakePicture() string
}

type IPhone interface {
	MakeCall() string
}

type camera struct {
}

func (c *camera) MakePicture() string {
	return "Picture Done"
}

type phone struct {
}

func (c *phone) MakeCall() string {
	return "Call Done"
}

type CameraPhone struct {
	camera camera
	phone  phone
}

func (c CameraPhone) MakePicture() string {
	return c.camera.MakePicture()
}

func (c CameraPhone) MakeCall() string {
	return c.phone.MakeCall()
}

func TestAbstraction(t *testing.T) {
	cameraPhone := CameraPhone{
		camera: camera{},
		phone:  phone{},
	}

	log.Println(cameraPhone.MakePicture())
	log.Println(cameraPhone.MakeCall())

}
