package lsp

type WorkDoneProgressBeginRequest struct {
	Notification
	Params WorkDoneProgressBeginParams `json:"params"`
}

type WorkDoneProgressBeginParams struct {
	Token  string                           `json:"token"`
	Values WorkDoneProgressBeginParamsValue `json:"value"`
}

type WorkDoneProgressBeginParamsValue struct {
	Kind        string `json:"kind"`
	Title       string `json:"title"`
	Cancellable bool   `json:"cancellable"`
	Percentage  int    `json:"percentage"`
	Message     string `json:"message"`
}

type WorkDoneProgressReportRequest struct {
	Notification
	Params WorkDoneProgressReportParams `json:"params"`
}

type WorkDoneProgressReportParams struct {
	Token  string                            `json:"token"`
	Values WorkDoneProgressReportParamsValue `json:"value"`
}

type WorkDoneProgressReportParamsValue struct {
	Kind        string `json:"kind"`
	Cancellable bool   `json:"cancellable"`
	Percentage  int    `json:"percentage"`
	Message     string `json:"message"`
}

type WorkDoneProgressEndRequest struct {
	Notification
	Params WorkDoneProgressEndParams `json:"params"`
}

type WorkDoneProgressEndParams struct {
	Token  string                         `json:"token"`
	Values WorkDoneProgressEndParamsValue `json:"value"`
}

type WorkDoneProgressEndParamsValue struct {
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

func CreateProgressBeginRequest(token string, title string) WorkDoneProgressBeginRequest {
	return WorkDoneProgressBeginRequest{
		Notification: Notification{
			RPC:    "2.0",
			Method: "$/progress",
		},
		Params: WorkDoneProgressBeginParams{
			Token: token,
			Values: WorkDoneProgressBeginParamsValue{
				Kind:        "begin",
				Title:       title,
				Cancellable: false,
				Percentage:  0,
				Message:     "",
			},
		},
	}
}

func CreateProgressUpdateRequest(token string, message string, percentage int) WorkDoneProgressReportRequest {
	return WorkDoneProgressReportRequest{
		Notification: Notification{
			RPC:    "2.0",
			Method: "$/progress",
		},
		Params: WorkDoneProgressReportParams{
			Token: token,
			Values: WorkDoneProgressReportParamsValue{
				Kind:        "report",
				Cancellable: false,
				Percentage:  percentage,
				Message:     message,
			},
		},
	}
}

func CreateProgressEndRequest(token string, message string) WorkDoneProgressEndRequest {
	return WorkDoneProgressEndRequest{
		Notification: Notification{
			RPC:    "2.0",
			Method: "$/progress",
		},
		Params: WorkDoneProgressEndParams{
			Token: token,
			Values: WorkDoneProgressEndParamsValue{
				Kind:    "end",
				Message: message,
			},
		},
	}
}
