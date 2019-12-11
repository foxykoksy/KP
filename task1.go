package main

import (
	"fmt"
	"math"
	"os"
)
func main() {
	var x,y,z,zero int
	//x=120150220
	fmt.Fscan(os.Stdin,&x)
	for x>0{
		z=x%10
		x=(x-z)/10
		if(z==0) {
			zero++
			continue
		}
		for i:=1;i<100;i++{
			if(y%int(math.Pow10(i))<=z*int(math.Pow10(i-1))){
				y=(y/int(math.Pow10(i-1))*10+z)*int(math.Pow10(i-1))+y%int(math.Pow10(i-1))
				break
			}
		}
	}
	for zero>0{
		fmt.Print("0")
		zero--
	}
	fmt.Print(y)
}
