package tests

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

type Builder struct {
	windowType string
	doorType   string
	floor      int
}

// Builder pattern
func newBuilder() Builder {
	return Builder{}
}

func builderWithFields() Builder {
	return Builder{
		windowType: "Wooden type",
		doorType:   "Wooden type",
		floor:      1,
	}
}

func (b Builder) setWindowType(s string) Builder {
	b.windowType = s
	return b
}

func (b Builder) setDoorType(s string) Builder {
	b.doorType = s
	return b
}

func (b Builder) setNumFloor(n int) Builder {
	b.floor = n
	return b
}

func TestBuilderPattern(t *testing.T) {
	build := newBuilder().setWindowType("Wooden type").setDoorType("Plastic type").setNumFloor(2)
	house := builderWithFields()
	log.Println(house)
	log.Println(build)
}

// Factory pattern
type CarFactory struct{}

type iCar interface {
	getName() string
	getPrice() int
}

type Car struct {
	name  string
	price int
}

type Audi struct {
	Car
}

type Bmw struct {
	Car
}

func (c *Car) getName() string {
	return c.name
}

func (c *Car) getPrice() int {
	return c.price
}

func newAudi() iCar {
	return &Audi{
		Car: Car{
			name:  "Audi",
			price: 12000,
		},
	}
}

func newBmw() iCar {
	return &Bmw{
		Car: Car{
			name:  "Bmw",
			price: 23000,
		},
	}
}

func (cf *CarFactory) createCar(carName string) (iCar, error) {
	if carName == "Audi" {
		return newAudi(), nil
	}
	if carName == "Bmw" {
		return newBmw(), nil
	}

	return nil, fmt.Errorf("Invalid car name")
}

func TestFactoryPattern(t *testing.T) {
	car := CarFactory{}
	audi, _ := car.createCar("Audi")
	bmw, _ := car.createCar("Bmw")

	log.Println(audi, "Price :", audi.getPrice())
	log.Println(bmw, "Price :", bmw.getPrice())
}

// Prototype pattern
type Employee struct {
	Name    string
	Gender  string
	Address *Address
}

type Address struct {
	StreetName string
	City       string
}

func (add *Address) DeepCopy() *Address {
	return &Address{
		StreetName: add.StreetName,
		City:       add.City,
	}
}

func TestPro(t *testing.T) {

	employee1 := Employee{
		Name:   "John",
		Gender: "Male",
		Address: &Address{
			StreetName: "Baker Street",
			City:       "London",
		},
	}

	employee2 := employee1
	employee2.Name = "Surya"
	employee2.Address = employee1.Address.DeepCopy()
	employee2.Address.StreetName = "Marine Drive"
	employee2.Address.City = "Mumbai"

	fmt.Println(employee1, employee1.Address)
	fmt.Println(employee2, employee2.Address)

}

// Singleton Pattern
var lock = &sync.Mutex{}

type single struct {
}

var singleInstance *single

func getInstance() *single {

	if singleInstance == nil {

		lock.Lock()
		defer lock.Unlock()

		if singleInstance == nil {
			fmt.Println("Creating Single Instance Now")
			singleInstance = &single{}

		} else {
			fmt.Println("Single Instance already created-1")
		}
	}

	return singleInstance
}

func TestSingletonPattern(t *testing.T) {
	for i := 0; i < 30; i++ {
		go getInstance()
	}

	time.Sleep(time.Millisecond * 200)
}

// Abstract factory
type iShoe interface {
	getLogo() string
	getSize() string
}

type shoe struct {
	logo string
	size string
}

func (s *shoe) getLogo() string {
	return s.logo
}

func (s *shoe) getSize() string {
	return s.size
}

type iShort interface {
	getLogo() string
	getSize() string
}

type short struct {
	logo string
	size string
}

func (s *short) getLogo() string {
	return s.logo
}

func (s *short) getSize() string {
	return s.size
}

type adidasShirt struct {
	shoe
}

type adidasShort struct {
	short
}

type adidas struct{}

func (a *adidas) makeShort() iShoe {
	return &adidasShirt{
		shoe: shoe{
			logo: "adidas",
			size: "M",
		},
	}
}

func (a *adidas) makeShirt() iShort {
	return &adidasShort{
		short: short{
			logo: "adidas",
			size: "M",
		},
	}
}

type nikeShirt struct {
	shoe
}

type nikeShort struct {
	short
}

type nike struct{}

func (b *nike) makeShort() iShoe {
	return &nikeShirt{
		shoe: shoe{
			logo: "nike",
			size: "S",
		},
	}
}

func (b *nike) makeShirt() iShort {
	return &nikeShort{
		short: short{
			logo: "nike",
			size: "S",
		},
	}
}

type iSportsAttireFactory interface {
	makeShort() iShoe
	makeShirt() iShort
}

func getSportsAttireFactory(brand string) (iSportsAttireFactory, error) {
	if brand == "adidas" {
		return &adidas{}, nil
	}
	if brand == "nike" {
		return &nike{}, nil
	}
	return nil, fmt.Errorf("Wrong brand type passed")
}

func TestAbstract(t *testing.T) {
	adidasFactory, _ := getSportsAttireFactory("adidas")
	nikeFactory, _ := getSportsAttireFactory("nike")

	nikeShoe := nikeFactory.makeShort()
	nikeShort := nikeFactory.makeShirt()

	adidasShoe := adidasFactory.makeShort()
	adidasShort := adidasFactory.makeShirt()

	log.Println(nikeShoe.getLogo())
	log.Println(nikeShort.getLogo())
	log.Println(adidasShoe.getLogo())
	log.Println(adidasShort.getLogo())

}

type computer interface {
	insertInSquarePort()
}

type mac struct{}

func (m *mac) insertInSquarePort() {
	log.Println("Insert square port into mac machine")
}

type windows struct{}

func (w *windows) insertInCirclePort() {
	log.Println("Insert circle port into windows machine")
}

type windowsAdapter struct {
	windowsMachine *windows
}

func (w *windowsAdapter) insertInSquarePort() {
	log.Println("Windows port with adapter")
	w.windowsMachine.insertInCirclePort()
}

type client struct {
}

func (c *client) insertUsbInComputer(com computer) {
	com.insertInSquarePort()
}

func TestAdapterPattern(t *testing.T) {
	client := &client{}
	mac := &mac{}

	client.insertUsbInComputer(mac)
	windowsMachine := &windows{}
	windowsMachineAdapter := &windowsAdapter{
		windowsMachine: windowsMachine,
	}
	client.insertUsbInComputer(windowsMachineAdapter)
}

type printer interface {
	printFile()
}

type comp interface {
	print()
	setPrinter(printerValue printer)
}
