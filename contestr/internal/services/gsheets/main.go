package main

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

func main() {
	ctx := context.Background()

	b, err := os.ReadFile("secrets/contestr-90c9191419de.json")
	if err != nil {
		log.Fatalf("read credentials file: %v", err)
	}

	conf, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("parse credentials: %v", err)
	}

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		log.Fatalf("create Sheets service: %v", err)
	}

	spreadsheetId := "17Zi5YORxT-pHG_QQrcfs7_y55SqJz-ShGxYk_w5DzuM"
	readRange := "groups!A1:H20" // Например, диапазон для чтения

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Data:")
		for _, row := range resp.Values {
			fmt.Println(row)
		}
	}
}
