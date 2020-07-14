package auxiliary

import "parse/structs"

// GetOfficeNo
func GetOfficeNo(order *structs.BuyOrder) string {
	if order.OfficeNo == "" {
		return order.Pnr["OfficeNo"]
	}
	return order.OfficeNo
}
