package main

import (
	"fmt"
)

type Base struct{}

func (self *Base) Magic() {
	fmt.Println("base magic")
}

func (self *Base) MoreMagic() {
	self.Magic()
	self.Magic()
}

type Voodoo struct {
	// 匿名item，表示的是包含和囊括，而不是集成
	// 可以通过Voodoo直接调用Base的任何方法
	// 类似 v.MoreMagic()
	Base

	// base Base
	// 如果是根据如上的形式进行定义，则表示Base为Voodoo的一个成员
	// 调用Base的方法则需要通过下列的形式
	// v.base.MoreMagic
}

func (self *Voodoo) Magic() {
	fmt.Println("voodoo magic")
}

func main() {
	v := new(Voodoo)
	v.Magic()
	v.MoreMagic()
}
