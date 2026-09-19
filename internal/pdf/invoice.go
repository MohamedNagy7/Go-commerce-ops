package pdf

import (
	"bytes"
	"fmt"

	"codeberg.org/go-pdf/fpdf"
	"github.com/MohamedNagy7/Go-commerce-ops/internal/invoice"
)

// ---------- Color palette ----------
// Matches the favicon: dark navy + coral accent.

type pdfColors struct {
	r, g, b int
}

var (
	colorPrimary = pdfColors{28, 30, 38}    // dark navy (header banner)
	colorAccent  = pdfColors{240, 90, 74}   // coral (favicon dot)
	colorDark    = pdfColors{44, 62, 80}    // dark slate (section labels)
	colorLight   = pdfColors{230, 232, 236} // hairlines
	colorRowAlt  = pdfColors{246, 247, 249} // faint row stripe
	colorText    = pdfColors{60, 60, 60}
	colorMuted   = pdfColors{130, 133, 140}
	colorWhite   = pdfColors{255, 255, 255}
)

func setFillColor(p *fpdf.Fpdf, c pdfColors) { p.SetFillColor(c.r, c.g, c.b) }
func setDrawColor(p *fpdf.Fpdf, c pdfColors) { p.SetDrawColor(c.r, c.g, c.b) }
func setTextColor(p *fpdf.Fpdf, c pdfColors) { p.SetTextColor(c.r, c.g, c.b) }

func money(v float64) string {
	return fmt.Sprintf("$%.2f", v)
}

// ---------- Invoice generation ----------

// GenerateInvoicePDF renders an invoice.InvoiceData into a PDF and returns
// the raw bytes. logoPNG is optional -- pass nil to render without a logo mark.
func GenerateInvoicePDF(data invoice.InvoiceData, logoPNG []byte) ([]byte, error) {
	pdfDoc := fpdf.New("P", "mm", "A4", "")
	pdfDoc.SetMargins(15, 15, 15)
	pdfDoc.SetAutoPageBreak(true, 20)
	pdfDoc.AddPage()

	pageWidth, _ := pdfDoc.GetPageSize()
	marginLeft, _, _, _ := pdfDoc.GetMargins()
	contentWidth := pageWidth - marginLeft*2

	// ---------- Header banner ----------
	const bannerH = 40.0
	setFillColor(pdfDoc, colorPrimary)
	pdfDoc.Rect(0, 0, pageWidth, bannerH, "F")

	// Logo badge: white circle with the favicon image inside, top-left.
	if len(logoPNG) > 0 {
		const badgeCX, badgeCY, badgeR = 25.0, bannerH / 2, 11.0
		setFillColor(pdfDoc, colorWhite)
		pdfDoc.Circle(badgeCX, badgeCY, badgeR, "F")
		opt := fpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
		pdfDoc.RegisterImageOptionsReader("logo", opt, bytes.NewReader(logoPNG))
		const logoSize = 15.0
		pdfDoc.ImageOptions("logo", badgeCX-logoSize/2, badgeCY-logoSize/2, logoSize, logoSize, false, opt, 0, "")
	}

	titleX := marginLeft + 24
	pdfDoc.SetXY(titleX, 12)
	setTextColor(pdfDoc, colorWhite)
	pdfDoc.SetFont("Arial", "B", 20)
	pdfDoc.Cell(100, 10, "INVOICE")

	pdfDoc.SetXY(titleX, 23)
	pdfDoc.SetFont("Arial", "", 10)
	setTextColor(pdfDoc, pdfColors{200, 202, 208})
	pdfDoc.Cell(120, 6, fmt.Sprintf("Order #%s", data.OrderID))

	pdfDoc.SetXY(pageWidth-marginLeft-70, 12)
	pdfDoc.SetFont("Arial", "", 10)
	setTextColor(pdfDoc, colorWhite)
	pdfDoc.CellFormat(70, 6, data.IssuedAt.Format("Jan 02, 2006"), "", 2, "R", false, 0, "")
	pdfDoc.SetX(pageWidth - marginLeft - 70)
	pdfDoc.SetFont("Arial", "", 8)
	setTextColor(pdfDoc, pdfColors{200, 202, 208})
	pdfDoc.CellFormat(70, 5, data.IssuedAt.Format("3:04 PM"), "", 2, "R", false, 0, "")

	// ---------- Greeting + status ----------
	y := bannerH + 14
	pdfDoc.SetXY(marginLeft, y)
	pdfDoc.SetFont("Arial", "", 11)
	setTextColor(pdfDoc, colorText)
	pdfDoc.Cell(contentWidth, 6, fmt.Sprintf("Hello %s,", data.CustomerEmail))
	pdfDoc.Ln(7)
	pdfDoc.SetX(marginLeft)
	pdfDoc.SetFont("Arial", "", 11)
	pdfDoc.Cell(contentWidth, 6, "Your order was successful. Here's your invoice.")

	// ---------- Amount charged callout ----------
	y = bannerH + 34
	setFillColor(pdfDoc, pdfColors{247, 248, 250})
	pdfDoc.RoundedRect(marginLeft, y, contentWidth, 22, 3, "1234", "F")
	pdfDoc.SetXY(marginLeft+8, y+4)
	pdfDoc.SetFont("Arial", "", 9)
	setTextColor(pdfDoc, colorMuted)
	pdfDoc.Cell(100, 5, "AMOUNT CHARGED")
	pdfDoc.SetXY(marginLeft+8, y+9)
	pdfDoc.SetFont("Arial", "B", 18)
	setTextColor(pdfDoc, colorPrimary)
	pdfDoc.Cell(100, 9, money(data.Total))

	setFillColor(pdfDoc, colorAccent)
	setTextColor(pdfDoc, colorWhite)
	pdfDoc.SetFont("Arial", "B", 9)
	badgeW := 20.0
	pdfDoc.SetXY(pageWidth-marginLeft-badgeW-8, y+7)
	pdfDoc.CellFormat(badgeW, 8, "PAID", "", 0, "C", true, 0, "")

	// ---------- Items table header ----------
	y = bannerH + 66
	pdfDoc.SetY(y)
	setFillColor(pdfDoc, colorDark)
	setTextColor(pdfDoc, colorWhite)
	pdfDoc.SetFont("Arial", "B", 10)

	colItem, colQty, colPrice, colSub := 90.0, 20.0, 35.0, 35.0
	pdfDoc.CellFormat(colItem, 9, "Item", "", 0, "L", true, 0, "")
	pdfDoc.CellFormat(colQty, 9, "Qty", "", 0, "C", true, 0, "")
	pdfDoc.CellFormat(colPrice, 9, "Price", "", 0, "R", true, 0, "")
	pdfDoc.CellFormat(colSub, 9, "Subtotal", "", 1, "R", true, 0, "")

	// ---------- Table rows with striping ----------
	setTextColor(pdfDoc, colorText)
	pdfDoc.SetFont("Arial", "", 10)

	for i, item := range data.Items {
		fill := i%2 == 0
		if fill {
			setFillColor(pdfDoc, colorRowAlt)
		}
		setDrawColor(pdfDoc, colorLight)

		pdfDoc.CellFormat(colItem, 9, item.ProductName, "", 0, "L", fill, 0, "")
		pdfDoc.CellFormat(colQty, 9, fmt.Sprintf("%d", item.Quantity), "", 0, "C", fill, 0, "")
		pdfDoc.CellFormat(colPrice, 9, money(item.Price), "", 0, "R", fill, 0, "")
		pdfDoc.CellFormat(colSub, 9, money(float64(item.Quantity)*item.Price), "", 1, "R", fill, 0, "")
	}

	setDrawColor(pdfDoc, colorLight)
	pdfDoc.SetLineWidth(0.2)
	pdfDoc.Line(marginLeft, pdfDoc.GetY(), pageWidth-marginLeft, pdfDoc.GetY())

	// ---------- Totals ----------
	pdfDoc.Ln(5)
	totalsX := marginLeft + contentWidth - 90

	drawTotalsRow := func(label string, value string, bold bool, fill bool, fontSize float64) {
		pdfDoc.SetX(totalsX)
		if bold {
			pdfDoc.SetFont("Arial", "B", fontSize)
		} else {
			pdfDoc.SetFont("Arial", "", fontSize)
		}
		if fill {
			setFillColor(pdfDoc, colorPrimary)
			setTextColor(pdfDoc, colorWhite)
		} else {
			setTextColor(pdfDoc, colorText)
		}
		pdfDoc.CellFormat(55, 8, label, "", 0, "R", fill, 0, "")
		pdfDoc.CellFormat(35, 8, value, "", 1, "R", fill, 0, "")
	}

	drawTotalsRow("Subtotal:", money(data.Subtotal), false, false, 10)
	pdfDoc.Ln(1)
	drawTotalsRow("Total:", money(data.Total), true, true, 12)

	// ---------- Footer ----------
	// Anchored a fixed distance below the totals, not pinned to the bottom
	// of the page -- avoids the large dead zone on short invoices.
	pdfDoc.Ln(16)
	setDrawColor(pdfDoc, colorLight)
	pdfDoc.SetLineWidth(0.2)
	pdfDoc.Line(marginLeft, pdfDoc.GetY(), pageWidth-marginLeft, pdfDoc.GetY())
	pdfDoc.Ln(4)

	pdfDoc.SetFont("Arial", "I", 8.5)
	setTextColor(pdfDoc, colorMuted)
	pdfDoc.CellFormat(contentWidth, 5, "Thank you for your purchase! This invoice was generated automatically.", "", 1, "C", false, 0, "")
	pdfDoc.SetFont("Arial", "", 8)
	pdfDoc.CellFormat(contentWidth, 5, fmt.Sprintf("Order #%s  |  %s", data.OrderID, data.IssuedAt.Format("Jan 02, 2006 3:04 PM")), "", 1, "C", false, 0, "")

	// ---------- Output ----------
	var buf bytes.Buffer
	if err := pdfDoc.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}
