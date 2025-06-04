package utils

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var googleSheetID string = os.Getenv("GoogleSheets")

func SaveToGoogleSheets(client Client) error {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("credentials.json"))
	if err != nil {
		return fmt.Errorf("Unable to create Sheets service: %v", err)
	}

	// Prepare data for Google Sheets with all fields in the correct order
	values := [][]interface{}{
		{
			client.Name,
			client.Phone,
			client.Address,
			client.Message,
			client.CallWeek,
			client.HighValue,
			client.BookingDone,
			client.BookingDate,
			client.Signed,
			client.SignedPrice,
			client.WorkDate,
			client.LeadUser,
		},
	}

	// Update the range to reflect all columns used (A to I here)
	sheetRange := "Sheet3!A:I"
	valueRange := &sheets.ValueRange{
		Values: values,
	}

	_, err = srv.Spreadsheets.Values.Append(googleSheetID, sheetRange, valueRange).
		ValueInputOption("RAW").
		Do()

	if err != nil {
		return fmt.Errorf("Unable to append data to Google Sheets: %v", err)
	}

	fmt.Println("Data successfully written to Google Sheets")
	return nil
}
