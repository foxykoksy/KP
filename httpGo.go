package main

import (
	memorycache"awesomeProject3/cache"
	"awesomeProject3/dbmethods"
	"fmt"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
	"time"
)

func main() {

	db, err := dbmethods.NewItemTable()
	if err != nil {
		panic(err)
	}
	/*tmp := ProductsInfo.Product{}
	tmp.Company = "JSandB"
	tmp.Amount = 77
	tmp.Price = 123
	tmp.Item = "Thing"
	db.AddItem(&tmp)
	tmp2 := ProductsInfo.Info{}
	tmp2.Rating = 10
	tmp2.Company = "HHeaven"
	tmp2.Information = "Come for us"
	db.AddInfo(&tmp2)*/
	str, _ := db.GetInfo("HHeaven")
	str2, _ := db.GetItem("Thing", "JSandB")
	fmt.Println(str)
	fmt.Println(str2)

fmt.Println("Работать ли с cache?\n да - введите 1\n нет - введите 0")
	var check int
	fmt.Fscan(os.Stdin, &check)
	if check == 1 {
		fmt.Println("Сохраним в кеше информацию о компании с бд. Время жизни контейнера одна минута," +
			" как и нашей записи.")
		cache := memorycache.New(time.Minute, 10*time.Minute)
		str, _ := db.GetInfo("HHeaven")
		cache.Set("myKey", str, time.Minute)
		fmt.Println("Хотим вывести информацию о компании?\n да - введите 1\n нет - введите 0")
		fmt.Fscan(os.Stdin, &check)
		if check == 1 {
			i, b := cache.Get("myKey")
			fmt.Printf("%s %t", "Информация существует?", b)
			fmt.Printf("%s %s", "\nИнформация: ", i)
		} else {
			fmt.Println("Хотим удалить из кеша или подождать пока закончиться время жизни?" +
				"\n Удалить - 1\n Ждать - 0")
			fmt.Fscan(os.Stdin, &check)
			if check == 0 {
				fmt.Println("Подождем минуту, чтобы отчистился кэш")
				time.Sleep(time.Minute)
				i, b := cache.Get("myKey")
				fmt.Printf("%s %t", "Информация существует?", b)
				fmt.Printf("%s %s", "\nИнформация: ", i)
				fmt.Println("\nПроверим действительно ли у объекта закончилось время жизни\n" +
					"выведем массив ключей у которых закончилось время")
				fmt.Println(cache.ExpiredKeys())
			} else if check == 1 {
				cache.Delete("myKey")
				i, b := cache.Get("myKey")
				fmt.Printf("%s %t", "Информация существует?", b)
				fmt.Printf("%s %s", "\nИнформация: ", i)
				fmt.Println("\nПроверим действительно ли удалили то что хранилось в кеше под нашим ключем или" +
					" у объекта закончилось время жизни\nвыведем массив ключей у которых закончилось время")
				fmt.Println(cache.ExpiredKeys())
			}
		}
	}
	s2, _ := db.GetInfo("HHeaven")
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world!\n"+s2)
	})
	e.Logger.Fatal(e.Start(":1323")) // http://localhost:1323
}