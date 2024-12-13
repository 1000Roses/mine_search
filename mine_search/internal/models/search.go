package models

type (
	ConditionSearching struct {
		ConditionName  string `json:"name"`
		ConditionValue string `json:"value"`
	}
	OrderBySearching struct {
		OrderByName  string `json:"name"`
		OrderByValue string `json:"value"`
	}
	Search struct {
		Text       string               `json:"text" validate:"required"`
		Conditions []ConditionSearching `json:"conditions"`
		OrderBys   []OrderBySearching   `json:"order_bys"`
	}
)
