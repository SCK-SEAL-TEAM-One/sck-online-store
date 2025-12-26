package order

import (
	"fmt"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type PDFHelper interface {
	GenerateOrderSummaryPDF(orderSummary OrderSummary) ([]byte, error)
}

type OrderSummaryPDFGenerator struct{}

type LineItem struct {
	No    int
	Brand string
	Title string
	Price float64
	Unit  int
}

func (orderSummaryPDF OrderSummaryPDFGenerator) GenerateOrderSummaryPDF(orderSummary OrderSummary) ([]byte, error) {
	company := "SCK Shopping Mall"
	formatter := message.NewPrinter(language.English)

	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(20).
		WithTopMargin(10).
		WithRightMargin(20).
		Build()

	page := maroto.New(cfg)

	defaultFontProp := props.Text{
		Family: fontfamily.Helvetica,
		Size:   10,
	}

	// --- HEADER ---
	page.AddRow(10,
		text.NewCol(12, "Order Summary",
			props.Text{
				Family: fontfamily.Helvetica, Style: fontstyle.Bold, Align: align.Center, Size: 16, Top: 2,
			}),
	)

	// Company Box
	page.AddRow(12,
		text.NewCol(12, company, props.Text{
			Family: fontfamily.Helvetica, Style: fontstyle.Bold, Align: align.Center, Size: 14, Top: 2,
		}),
	).WithStyle(&props.Cell{BorderType: border.Full})

	page.AddRow(5, col.New(12))

	// --- CUSTOMER INFO ---
	orderIDText := fmt.Sprintf("Order No.: %d", orderSummary.OrderID)
	fullNameText := fmt.Sprintf("Full Name: %s %s", orderSummary.FirstName, orderSummary.LastName)
	paymentMethodText := fmt.Sprintf("Payment Method: %s", orderSummary.PaymentMethod)
	trackingNumberText := fmt.Sprintf("Tracking Number: %s", orderSummary.TrackingNumber)

	addInfoRow(page, orderIDText, defaultFontProp)
	addInfoRow(page, fullNameText, defaultFontProp)
	addInfoRow(page, paymentMethodText, defaultFontProp)
	addInfoRow(page, trackingNumberText, defaultFontProp)

	page.AddRow(5, col.New(12))

	// --- TABLE HEADER
	headerStyle := &props.Cell{BorderType: border.Full}
	headerText := props.Text{Family: fontfamily.Helvetica, Style: fontstyle.Bold, Align: align.Center, Top: 2}

	page.AddRow(10,
		text.NewCol(1, "No", headerText).WithStyle(headerStyle),
		text.NewCol(7, "Product Brand - Product Title", headerText).WithStyle(headerStyle),
		text.NewCol(3, "Price", headerText).WithStyle(headerStyle),
		text.NewCol(1, "Unit", headerText).WithStyle(headerStyle),
	)

	// --- TABLE CONTENT ---
	rowNormalNumber := props.Text{Family: fontfamily.Helvetica, Align: align.Center, Top: 2}
	rowText := props.Text{Family: fontfamily.Helvetica, Align: align.Left, Top: 2, Left: 2}
	textCurrency := props.Text{Family: fontfamily.Helvetica, Align: align.Left, Right: 2, Top: 2, Left: 2}
	valueCurrency := props.Text{Family: fontfamily.Helvetica, Align: align.Right, Right: 2, Top: 2}
	boxStyle := &props.Cell{BorderType: border.Full}

	for index, product := range orderSummary.OrderProductList {
		page.AddRow(8,
			text.NewCol(1, fmt.Sprintf("%d", index+1), rowNormalNumber).WithStyle(boxStyle),
			text.NewCol(7, fmt.Sprintf("%s - %s", product.ProductBrand, product.ProductName), rowText).WithStyle(boxStyle),
			col.New(3).Add(
				text.New("THB", textCurrency),
				text.New(formatter.Sprintf("%.2f", product.PriceTHB), valueCurrency),
			).WithStyle(boxStyle),
			text.NewCol(1, fmt.Sprintf("%d", product.Quantity), rowNormalNumber).WithStyle(boxStyle),
		)
	}

	page.AddRow(5, col.New(12))

	// --- TOTALS ---
	rowTextBold := props.Text{Family: fontfamily.Helvetica, Style: fontstyle.Bold, Size: 10, Align: align.Left, Top: 2, Left: 2}
	addTotalRow(page, "Merchandise Subtotal", formatter.Sprintf("%.2f", orderSummary.SubTotalPrice), rowTextBold)
	addTotalRow(page, "Shipping Fee", formatter.Sprintf("%.2f", orderSummary.ShippingFee), rowTextBold)
	addTotalRow(page, "Total Price", formatter.Sprintf("%.2f", orderSummary.TotalPrice), rowTextBold)

	page.AddRow(5, col.New(12))

	// --- RECEIVING POINT ---
	page.AddRow(8,
		col.New(7), // Spacer
		text.NewCol(3, "Receiving Point", rowTextBold).WithStyle(boxStyle),
		col.New(2).Add(
			text.New(fmt.Sprintf("%d", orderSummary.ReceivingPoint), props.Text{
				Family: fontfamily.Helvetica, Style: fontstyle.Bold, Align: align.Right, Right: 2, Top: 2,
			}),
		).WithStyle(boxStyle),
	)

	document, err := page.Generate()
	if err != nil {
		return nil, err
	}
	return document.GetBytes(), nil

}

func addInfoRow(page core.Maroto, content string, props props.Text) {
	page.AddRow(6, text.NewCol(12, content, props))
}

func addTotalRow(page core.Maroto, label string, value string, textProps props.Text) {
	// Align Right
	propsLabel := textProps
	propsLabel.Align = align.Left
	propsLabel.Top = 2

	propsVal := textProps
	propsVal.Align = align.Right
	propsVal.Top = 2
	propsVal.Right = 2

	// Currency Props (Left aligned for "THB")
	propsCurr := textProps
	propsCurr.Align = align.Left
	propsCurr.Top = 2
	propsCurr.Left = 2 // Padding from left edge

	boxStyle := &props.Cell{BorderType: border.Full}

	page.AddRow(8,
		col.New(7), // Spacer
		text.NewCol(3, label, propsLabel).WithStyle(boxStyle),
		col.New(2).Add(
			text.New("THB", propsCurr),
			text.New(value, propsVal),
		).WithStyle(boxStyle),
	)
}
