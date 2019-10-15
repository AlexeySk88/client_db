package orders

type Order struct {
	OrderID int
	DistrictID int
	Price float64
	EntryIDs []int
}