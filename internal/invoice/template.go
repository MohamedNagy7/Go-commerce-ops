package invoice

import (
	"fmt"
	"strings"
)

func money(v float64) string {
	return fmt.Sprintf("$%.2f", v)
}

// BuildInvoiceHTML renders the invoice notification email. The logo is
// referenced as cid:logo -- pair this with an inline Attachment whose
// ContentID is "logo" (see Processor.Process).
func BuildInvoiceHTML(data InvoiceData) string {
	var rows strings.Builder
	for _, item := range data.Items {
		rows.WriteString(fmt.Sprintf(`
              <tr>
                <td style="padding:10px 0;border-bottom:1px solid #eceef1;color:#3c3c3c;font-size:14px;">%s</td>
                <td style="padding:10px 0;border-bottom:1px solid #eceef1;color:#3c3c3c;font-size:14px;text-align:center;">%d</td>
                <td style="padding:10px 0;border-bottom:1px solid #eceef1;color:#3c3c3c;font-size:14px;text-align:right;">%s</td>
              </tr>`, item.ProductName, item.Quantity, money(float64(item.Quantity)*item.Price)))
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="margin:0;padding:0;background-color:#f4f5f7;font-family:Arial,Helvetica,sans-serif;">
  <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f5f7;padding:32px 0;">
    <tr>
      <td align="center">
        <table role="presentation" width="560" cellpadding="0" cellspacing="0" style="background-color:#ffffff;border-radius:10px;overflow:hidden;">
          <tr>
            <td style="background-color:#1c1e26;padding:24px 32px;">
              <table role="presentation" cellpadding="0" cellspacing="0">
                <tr>
                  <td style="width:44px;">
                    <img src="cid:logo" width="36" height="36" style="display:block;border-radius:50%%;" alt="logo" />
                  </td>
                  <td style="padding-left:12px;">
                    <div style="color:#ffffff;font-size:20px;font-weight:bold;">Invoice</div>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:28px 32px 0 32px;">
              <p style="margin:0 0 4px 0;color:#3c3c3c;font-size:14px;">Hello %s,</p>
              <p style="margin:0;color:#3c3c3c;font-size:14px;">Your order <strong>#%s</strong> was successful on %s</p>
            </td>
          </tr>
          <tr>
            <td style="padding:20px 32px 0 32px;">
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f7f8fa;border-radius:8px;">
                <tr>
                  <td style="padding:16px 20px;">
                    <div style="color:#82858c;font-size:11px;letter-spacing:.04em;text-transform:uppercase;">Amount Charged</div>
                    <div style="color:#1c1e26;font-size:24px;font-weight:bold;margin-top:4px;">%s</div>
                  </td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:24px 32px 0 32px;">
              <div style="color:#1c1e26;font-size:13px;font-weight:bold;margin-bottom:8px;">Order Details</div>
              <table role="presentation" width="100%%" cellpadding="0" cellspacing="0">
                <tr>
                  <td style="padding-bottom:8px;border-bottom:2px solid #1c1e26;color:#82858c;font-size:12px;text-transform:uppercase;">Item</td>
                  <td style="padding-bottom:8px;border-bottom:2px solid #1c1e26;color:#82858c;font-size:12px;text-transform:uppercase;text-align:center;">Qty</td>
                  <td style="padding-bottom:8px;border-bottom:2px solid #1c1e26;color:#82858c;font-size:12px;text-transform:uppercase;text-align:right;">Subtotal</td>
                </tr>%s
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:16px 32px 0 32px;" align="right">
              <table role="presentation" cellpadding="0" cellspacing="0">
                <tr>
                  <td style="padding:4px 12px;color:#3c3c3c;font-size:13px;">Subtotal</td>
                  <td style="padding:4px 0;color:#3c3c3c;font-size:13px;text-align:right;">%s</td>
                </tr>
                <tr>
                  <td style="padding:6px 12px;color:#1c1e26;font-size:14px;font-weight:bold;">Total</td>
                  <td style="padding:6px 0;color:#1c1e26;font-size:14px;font-weight:bold;text-align:right;">%s</td>
                </tr>
              </table>
            </td>
          </tr>
          <tr>
            <td style="padding:28px 32px 32px 32px;">
              <p style="margin:0;color:#82858c;font-size:12px;text-align:center;">Your invoice is attached as a PDF. Thank you for shopping with us.</p>
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, data.CustomerEmail, data.OrderID, data.IssuedAt.Format("Monday, January 2, 2006 3:04 PM"),
		money(data.Total), rows.String(), money(data.Subtotal), money(data.Total))
}
