package model

type StartTimerRequest struct {
	Duration       int    `json:"duration" ,validate:"required"`
	CallbackURL    string `json:"callback_url" ,validate:"required"`
	CallbackMethod string `json:"callback_method" ,validate:"required"`
}

type Response struct {
	TimerID string `json:"timer_id,omitempty"`
	Message string `json:"message"`
}

type QueueData struct {
	StartTimerRequest
	Response
}
