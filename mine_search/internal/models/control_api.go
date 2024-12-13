package models

type ControlApiTb struct {
	ApiName    string `json:"api_name"`
	Active     int    `json:"active"`
	Conditions string `json:"conditions"`
}

func (ControlApiTb) TableName() string { return "control_api_tb" }

func (ControlApiTb) ColumnApiName() string {
	return "api_name"
}

func (ControlApiTb) ColumnActive() string {
	return "active"
}
