package Demo

/* ============================================================================ */
/* ======================= Session, ShoppingCart & Item ======================= */
/* ============================================================================ */

/* simply a shopping session per customer level */
type Session struct {
	UserId string
	Cart *ShoppingCart
}

/* a real shopping cart containing what a customer interested to buy */
type ShoppingCart struct {
	Items map[string]Item
}

/* item (name + price) */
type Item struct {
	// a more practical field should be item_id instead of Name
	Name string
	Quantity int
	Price float32
}

// ctor
func NewItem(name string, quantity int, price float32) *Item {
	return &Item{
		Name: name,
		Quantity: quantity,
		Price: price,
	}
}
func NewShoppingCart() *ShoppingCart {
	return &ShoppingCart{
		Items: make(map[string]Item),
	}
}

func (cart *ShoppingCart) AddItemToCart(name string, quantity int, price float32) {
	// ignore checkings of existing items here (should actually do a check and then increment count instead
	ptrItem := NewItem(name, quantity, price)
	cart.Items[name] = *ptrItem
}

func (cart *ShoppingCart) CalculateTrxAmount() float32 {
	fTrx := float32(0.0)

	for _, ptrItem := range cart.Items {
		fTrx += ptrItem.Price * float32(ptrItem.Quantity)
	}
	// apply VAT 20% (check the Rules)
	fTrx *= 1.20

	// shipping fees (check the Rules)
	if fTrx > 10.0 {
		fTrx += 2.0
	} else {
		fTrx += 3.0
	}
	return fTrx
}


/* ========================================================================= */
/* ======================= Inventory & InventoryItem ======================= */
/* ========================================================================= */

/* inherited the values of Item PLUS an inventoryCount */
type InventoryItem struct {
	item *Item
	inventoryCount int
	ready bool
}

type Inventory struct {
	items map[string]InventoryItem
}

// ctor
func NewInventory() *Inventory {
	return &Inventory{
		items: make(map[string]InventoryItem),
	}
}
func NewInventoryItem(name string, price float32, inventoryCount int) *InventoryItem {
	// create the Item
	ptrItem := NewItem(name, 0, price)
	return &InventoryItem{
		item: ptrItem,
		inventoryCount: inventoryCount,
		ready: true,
	}
}

func (inv *InventoryItem) SetInventoryCount(count int) *InventoryItem {
	inv.inventoryCount = count
	return inv
}
func (inv *Inventory) FindInventoryItem(name string) *InventoryItem {
	item := inv.items[name]
	return &item
}
func (inv *Inventory) AddToInventory(name string, price float32) {
	ptrInventoryItem := inv.FindInventoryItem(name)
	// which means nil
	if ptrInventoryItem.ready == false {
		ptrIItem := NewInventoryItem(name, price, 1)
		inv.items[name] = *ptrIItem
	} else {
		ptrInventoryItem.inventoryCount++
		// ** need to set back the updated Object pointer to the inventory map.... Jesus (less convenient than Java)
		inv.items[name] = *ptrInventoryItem
	}
}
func (inv *InventoryItem) GetItemFromInventoryItem() *Item {
	return inv.item
}


