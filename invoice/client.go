// package invoice provides the /invoices APIs
package invoice

import (
	"fmt"
	"net/url"
	"strconv"

	stripe "github.com/stripe/stripe-go"
)

// Client is the client used to invoke /invoices APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New POSTs new invoices.
// For more details see https://stripe.com/docs/api#create_invoice.
func New(params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return getC().New(params)
}

func (c Client) New(params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	body := &url.Values{
		"customer": {params.Customer},
	}

	if len(params.Desc) > 0 {
		body.Add("description", params.Desc)
	}

	if len(params.Statement) > 0 {
		body.Add("statement_description", params.Statement)
	}

	if len(params.Sub) > 0 {
		body.Add("subscription", params.Sub)
	}

	params.AppendTo(body)

	token := c.Key
	if params.Fee > 0 {
		body.Add("application_fee", strconv.FormatUint(params.Fee, 10))
	}

	invoice := &stripe.Invoice{}
	err := c.B.Call("POST", "/invoices", token, body, invoice)

	return invoice, err
}

// Get returns the details of an invoice.
// For more details see https://stripe.com/docs/api#retrieve_invoice.
func Get(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return getC().Get(id, params)
}

func (c Client) Get(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}
		params.AppendTo(body)
	}

	invoice := &stripe.Invoice{}
	err := c.B.Call("GET", "/invoices/"+id, c.Key, body, invoice)

	return invoice, err
}

// Pay pays an invoice.
// For more details see https://stripe.com/docs/api#pay_invoice.
func Pay(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return getC().Pay(id, params)
}

func (c Client) Pay(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}
		params.AppendTo(body)
	}

	invoice := &stripe.Invoice{}
	err := c.B.Call("POST", fmt.Sprintf("/invoices/%v/pay", id), c.Key, body, invoice)

	return invoice, err
}

// Update updates an invoice's properties.
// For more details see https://stripe.com/docs/api#update_invoice.
func Update(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return getC().Update(id, params)
}

func (c Client) Update(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	var body *url.Values
	token := c.Key

	if params != nil {
		body = &url.Values{}

		if len(params.Desc) > 0 {
			body.Add("description", params.Desc)
		}

		if len(params.Statement) > 0 {
			body.Add("statement_description", params.Statement)
		}

		if len(params.Sub) > 0 {
			body.Add("subscription", params.Sub)
		}

		if params.Closed {
			body.Add("closed", strconv.FormatBool(true))
		}

		if params.Forgive {
			body.Add("forgiven", strconv.FormatBool(true))
		}

		if params.Fee > 0 {
			body.Add("application_fee", strconv.FormatUint(params.Fee, 10))
		}

		params.AppendTo(body)
	}

	invoice := &stripe.Invoice{}
	err := c.B.Call("POST", "/invoices/"+id, token, body, invoice)

	return invoice, err
}

// GetNext returns the upcoming invoice's properties.
// For more details see https://stripe.com/docs/api#retrieve_customer_invoice.
func GetNext(params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return getC().GetNext(params)
}

func (c Client) GetNext(params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	body := &url.Values{
		"customer": {params.Customer},
	}

	if len(params.Sub) > 0 {
		body.Add("subscription", params.Sub)
	}

	params.AppendTo(body)

	invoice := &stripe.Invoice{}
	err := c.B.Call("GET", "/invoices", c.Key, body, invoice)

	return invoice, err
}

// List returns a list of invoices.
// For more details see https://stripe.com/docs/api#list_customer_invoices.
func List(params *stripe.InvoiceListParams) *stripe.InvoiceIter {
	return getC().List(params)
}

func (c Client) List(params *stripe.InvoiceListParams) *stripe.InvoiceIter {
	type invoiceList struct {
		stripe.ListMeta
		Values []*stripe.Invoice `json:"data"`
	}

	var body *url.Values
	var lp *stripe.ListParams

	if params != nil {
		body = &url.Values{}

		if len(params.Customer) > 0 {
			body.Add("customer", params.Customer)
		}

		if params.Date > 0 {
			body.Add("date", strconv.FormatInt(params.Date, 10))
		}

		params.AppendTo(body)
		lp = &params.ListParams
	}

	return &stripe.InvoiceIter{stripe.GetIter(lp, body, func(b url.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &invoiceList{}
		err := c.B.Call("GET", "/invoices", c.Key, &b, list)

		ret := make([]interface{}, len(list.Values))
		for i, v := range list.Values {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// ListLines returns a list of line items.
// For more details see https://stripe.com/docs/api#invoice_lines.
func ListLines(params *stripe.InvoiceLineListParams) *stripe.InvoiceLineIter {
	return getC().ListLines(params)
}

func (c Client) ListLines(params *stripe.InvoiceLineListParams) *stripe.InvoiceLineIter {
	body := &url.Values{}
	var lp *stripe.ListParams

	if len(params.Customer) > 0 {
		body.Add("customer", params.Customer)
	}

	if len(params.Sub) > 0 {
		body.Add("subscription", params.Sub)
	}

	params.AppendTo(body)
	lp = &params.ListParams

	return &stripe.InvoiceLineIter{stripe.GetIter(lp, body, func(b url.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.InvoiceLineList{}
		err := c.B.Call("GET", fmt.Sprintf("/invoices/%v/lines", params.Id), c.Key, &b, list)

		ret := make([]interface{}, len(list.Values))
		for i, v := range list.Values {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

func getC() Client {
	return Client{stripe.GetBackend(), stripe.Key}
}