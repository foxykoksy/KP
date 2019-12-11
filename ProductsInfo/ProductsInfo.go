package ProductsInfo

import "fmt"

type Product struct {
	Id      int
	Item    string
	Company string
	Price   int
	Amount  int
}
type Info struct {
	Id          int
	Company     string
	Information string
	Rating      int
}

func ShowItem(p Product) {
	fmt.Printf("%s %d\n %s %s\n %s %s\n %s %d\n %s %d\n\n",
		"item id:", p.Id,
		"item:", p.Item,
		"company:", p.Company,
		"price(BLR):", p.Price,
		"amount:", p.Amount)
}