package Demo

import (
	"fmt"
	"github.com/DATA-DOG/godog"
)

/* "github.com/stretchr/testify/assert" - if you need Assert API */

var inventory *Inventory
var cart *ShoppingCart

func thereIsAWhichCosts(name string, price float32) error {
	inventory.AddToInventory(name, price)
	//fmt.Println("inventory till now ...", inventory.items)
	return nil
}

func iAddTheToTheBasket(nameOfProduct string) error {
	ptrIItem := inventory.FindInventoryItem(nameOfProduct)
	if ptrIItem.ready == true {
		ptrItem := ptrIItem.GetItemFromInventoryItem()
		cart.AddItemToCart(ptrItem.Name, 1, ptrItem.Price)

	} else {
		// else handling
		fmt.Println("NO such item yet... ignore for now though, could Panic()")
	}
	//fmt.Println("** cart contents => ", cart)

	return nil
}

func iShouldHaveProductsInTheBasket(counts int) error {
	itemLen := len(cart.Items)
	if counts != itemLen {
		//fmt.Printf("the cart's item count is NOT equals to the given count~ Expected %d, Instead got %d\n", itemLen, counts)
		//return godog.ErrPending
		return fmt.Errorf("the cart's item count is NOT equals to the given count~ Expected %d, Instead got %d\n", counts, itemLen)
	}
	return nil
}

func theOverallBasketPriceShouldBe(estimatedTrx float32) error {
	calTrx := cart.CalculateTrxAmount()
	if estimatedTrx != calTrx {
		//fmt.Printf("the estimated trx amount does NOT match with the calculated trx amount~ Expected %v, Instead got %v\n", calTrx, estimatedTrx)
		//return godog.ErrPending
		return fmt.Errorf("the estimated trx amount does NOT match with the calculated trx amount~ Expected %v, Instead got %v\n", estimatedTrx, calTrx )
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
	// sort of like ctor
	s.BeforeScenario(func(interface{}) {
		inventory = NewInventory()

		// should even include a session (anonymous for example); make it easy, ignore the session for now
		cart = NewShoppingCart()
	})


	s.Step(`^there is a "([^"]*)", which costs £([0-9]+?\.[0-9]+)$`, thereIsAWhichCosts)
	s.Step(`^I add the "([^"]*)" to the basket$`, iAddTheToTheBasket)
	s.Step(`^I should have (\d+) products in the basket$`, iShouldHaveProductsInTheBasket)
	s.Step(`^the overall basket price should be £([0-9]+?\.[0-9]+)$`, theOverallBasketPriceShouldBe)
}