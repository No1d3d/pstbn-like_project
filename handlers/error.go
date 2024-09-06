package handlers

type Result struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}

func ResultFromError(err error) Result {
	return Result{
		Success:      false,
		ErrorMessage: err.Error(),
	}
}

func BadResult(message string) Result {
	return Result{
		Success:      false,
		ErrorMessage: message,
	}
}
