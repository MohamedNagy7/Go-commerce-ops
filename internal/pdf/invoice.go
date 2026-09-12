package pdf

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"codeberg.org/go-pdf/fpdf"
	"github.com/MohamedNagy7/Go-commerce-ops/internal/rabbitmq"
)

// ---------- Color palette ----------

type pdfColors struct {
	r, g, b int
}

var (
	colorPrimary = pdfColors{41, 128, 185}  // blue
	colorDark    = pdfColors{44, 62, 80}    // dark slate
	colorLight   = pdfColors{236, 240, 241} // light gray
	colorRowAlt  = pdfColors{245, 247, 250} // faint blue
	colorText    = pdfColors{60, 60, 60}
	colorWhite   = pdfColors{255, 255, 255}
	colorGray    = pdfColors{130, 130, 130}
)

func setFillColor(p *fpdf.Fpdf, c pdfColors) { p.SetFillColor(c.r, c.g, c.b) }
func setDrawColor(p *fpdf.Fpdf, c pdfColors) { p.SetDrawColor(c.r, c.g, c.b) }
func setTextColor(p *fpdf.Fpdf, c pdfColors) { p.SetTextColor(c.r, c.g, c.b) }

func money(v float64) string {
	if v != float64(0) {
		return strconv.FormatFloat(v, 'f', 2, 64)
	}
	return fmt.Sprintf("$%.2f", v)
}

// ---------- Invoice generation ----------

func GenerateInvoicePDF(payload rabbitmq.InvoiceRequestedPayload) ([]byte, error) {
	pdfDoc := fpdf.New("P", "mm", "A4", "")
	pdfDoc.SetMargins(15, 15, 15)
	pdfDoc.SetAutoPageBreak(true, 25)
	pdfDoc.AddPage()

	pageWidth, _ := pdfDoc.GetPageSize()
	marginLeft, _, _, _ := pdfDoc.GetMargins()
	contentWidth := pageWidth - marginLeft*2

	// ---------- Header banner ----------
	setFillColor(pdfDoc, colorPrimary)
	pdfDoc.Rect(0, 0, pageWidth, 38, "F")

	pdfDoc.SetY(10)
	pdfDoc.SetX(marginLeft)
	setTextColor(pdfDoc, colorWhite)
	pdfDoc.SetFont("Arial", "B", 22)
	pdfDoc.Cell(100, 12, "INVOICE")

	pdfDoc.SetX(pageWidth - marginLeft - 60)
	pdfDoc.SetFont("Arial", "", 11)
	pdfDoc.Cell(60, 12, fmt.Sprintf("Date: %s", time.Now().Format("Jan 02, 2006")))

	setTextColor(pdfDoc, colorText)

	// ---------- Invoice meta ----------
	pdfDoc.SetY(48)
	pdfDoc.SetFont("Arial", "", 11)

	// Left: billed to
	pdfDoc.SetFont("Arial", "B", 11)
	setTextColor(pdfDoc, colorDark)
	pdfDoc.Cell(90, 7, "Billed To")
	pdfDoc.Ln(8)

	pdfDoc.SetFont("Arial", "", 11)
	setTextColor(pdfDoc, colorText)
	pdfDoc.Cell(90, 6, payload.Email)
	pdfDoc.Ln(6)

	// Right: invoice details
	pdfDoc.SetY(48)
	pdfDoc.SetX(marginLeft + 100)
	pdfDoc.SetFont("Arial", "", 11)
	pdfDoc.Cell(30, 7, "Invoice No:")
	pdfDoc.Cell(60, 7, payload.OrderID)
	pdfDoc.Ln(7)
	pdfDoc.SetX(marginLeft + 100)
	pdfDoc.Cell(30, 7, "Order ID:")
	pdfDoc.Cell(60, 7, payload.OrderID)

	pdfDoc.Ln(14)

	// ---------- Items table header ----------
	setFillColor(pdfDoc, colorDark)
	setTextColor(pdfDoc, colorWhite)
	pdfDoc.SetFont("Arial", "B", 11)

	colItem, colQty, colPrice, colSub := 90.0, 20.0, 35.0, 35.0
	pdfDoc.CellFormat(colItem, 10, "Item", "", 0, "L", true, 0, "")
	pdfDoc.CellFormat(colQty, 10, "Qty", "", 0, "C", true, 0, "")
	pdfDoc.CellFormat(colPrice, 10, "Price", "", 0, "R", true, 0, "")
	pdfDoc.CellFormat(colSub, 10, "Subtotal", "", 1, "R", true, 0, "")

	// ---------- Table rows with striping ----------
	setTextColor(pdfDoc, colorText)
	pdfDoc.SetFont("Arial", "", 10)

	for i, item := range payload.Items {
		fill := false
		if i%2 == 0 {
			setFillColor(pdfDoc, colorRowAlt)
			fill = true
		}
		setDrawColor(pdfDoc, colorLight)

		pdfDoc.CellFormat(colItem, 9, item.ProductName, "", 0, "L", fill, 0, "")
		pdfDoc.CellFormat(colQty, 9, fmt.Sprintf("%d", item.Quantity), "", 0, "C", fill, 0, "")
		pdfDoc.CellFormat(colPrice, 9, money(item.Price), "", 0, "R", fill, 0, "")
		pdfDoc.CellFormat(colSub, 9, money(float64(item.Quantity)*item.Price), "", 1, "R", fill, 0, "")
	}

	// ---------- Totals ----------
	pdfDoc.Ln(6)

	totalsX := marginLeft + contentWidth - 90

	// Compute subtotal from items
	var subtotal float64
	for _, item := range payload.Items {
		subtotal += float64(item.Quantity) * item.Price
	}

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
		pdfDoc.CellFormat(55, 9, label, "", 0, "R", fill, 0, "")
		pdfDoc.CellFormat(35, 9, value, "", 1, "R", fill, 0, "")
	}

	drawTotalsRow("Subtotal:", money(subtotal), false, false, 11)
	drawTotalsRow("Total:", money(payload.TotalAmount), true, true, 13)

	// ---------- Footer ----------
	pdfDoc.SetY(-30)
	setDrawColor(pdfDoc, colorLight)
	pdfDoc.SetLineWidth(0.3)
	pdfDoc.Line(marginLeft, pdfDoc.GetY(), pageWidth-marginLeft, pdfDoc.GetY())
	pdfDoc.Ln(4)

	pdfDoc.SetFont("Arial", "I", 9)
	setTextColor(pdfDoc, colorGray)
	pdfDoc.CellFormat(contentWidth, 6, "Thank you for your purchase! This invoice was generated automatically.", "", 1, "C", false, 0, "")

	// ---------- Output ----------
	var buf bytes.Buffer
	if err := pdfDoc.Output(&buf); err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}
	return buf.Bytes(), nil
}
