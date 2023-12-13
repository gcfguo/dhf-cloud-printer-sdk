package model

type (
	PrinterBindReq struct {
		PrinterName string `json:"printer_name"`
		HostID      string `json:"host_id"`
		PrinterType string `json:"printer_type"`
	}

	PrinterBindRes struct {
		PrinterID string `json:"printer_id"`
	}
)

type (
	PrinterUnBindReq struct {
		PrinterID string `json:"printer_id"`
	}
	PrinterUnBindRes struct {
	}
)

type (
	ListPrinterReq struct {
		PageNo   int `json:"page_no"`
		PageSize int `json:"page_size"`
	}
	ListPrinterRes struct {
		List  []*PrinterInf `json:"list"`
		Total int           `json:"total"`
	}
	PrinterInf struct {
		PrinterName string `json:"printer_name" `
		HostID      string `json:"host_id"     `
		PrinterId   string `json:"printer_id"   `
		PrinterType string `json:"printer_type"`
		Status      string `json:"status"       `
		CreateTime  string `json:"create_time"  `
	}
)

type (
	PrintReq struct {
		PrinterID   string `json:"printer_id"`
		ContentType string `json:"content_type"`
		MetaType    string `json:"meta_type"`
		MetaValue   string `json:"meta_value"`
	}
	PrintRes struct {
		PrinterName string `json:"printer_name"`
		MetaType    string `json:"meta_type"`
		MetaValue   string `json:"meta_value"`
	}
)
