package validation

import (
	"errors"
	"fmt"
	"receipt_processor/models"
	"regexp"
	"strings"
	"time"
)

func SanitizeAndValidate(receipt *models.Receipt) error {
	// Trim whitespace
	receipt.Retailer = strings.TrimSpace(receipt.Retailer)
	receipt.PurchaseDate = strings.TrimSpace(receipt.PurchaseDate)
	receipt.PurchaseTime = strings.TrimSpace(receipt.PurchaseTime)
	receipt.Total = strings.TrimSpace(receipt.Total)

	// Define a generic receipt error
	receiptError := errors.New("the receipt is invalid")

	// Validate retailer name
	retailerPattern := regexp.MustCompile(`^[\w\s\-&]+$`)
	if !retailerPattern.MatchString(receipt.Retailer) {
		fmt.Println("Retailer name invalid format:", receipt.Retailer)
		return receiptError
	}

	// Validate that there is at least one item
	if len(receipt.Items) < 1 {
		fmt.Println("Purchase items should be more than 1")
		return receiptError
	}

	// Validate items
	shortDescPattern := regexp.MustCompile(`^[\w\s\-]+$`)
	pricePattern := regexp.MustCompile(`^\d+\.\d{2}$`)
	for i := range receipt.Items {
		receipt.Items[i].ShortDescription = strings.TrimSpace(receipt.Items[i].ShortDescription)
		receipt.Items[i].Price = strings.TrimSpace(receipt.Items[i].Price)

		// Validate item short description
		if !shortDescPattern.MatchString(receipt.Items[i].ShortDescription) {
			fmt.Println("Purchase item description invalid format:", receipt.Items[i].ShortDescription)
			return receiptError
		}

		// Validate item price
		if !pricePattern.MatchString(receipt.Items[i].Price) {
			fmt.Println("Purchase item price invalid format:", receipt.Items[i].Price)
			return receiptError
		}
	}

	// Validate total
	if !pricePattern.MatchString(receipt.Total) {
		fmt.Println("Purchase total price invalid format:", receipt.Total)
		return receiptError
	}

	// Validate purchase date
	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		fmt.Println("Purchase date invalid format:", receipt.PurchaseDate)
		return receiptError
	}

	// Validate purchase time
	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		fmt.Println("Purchase time invalid format:", receipt.PurchaseTime)
		return receiptError
	}

	return nil
}
