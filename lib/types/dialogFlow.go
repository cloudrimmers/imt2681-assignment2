package types

// Response - A representation of the response from DialogFlow
type Response struct {
	Result struct {
		//NOTE: If need be, place ADDITIONAL PARAMETERS
		Parameters struct {
			CurrencyOut struct {
				CurrencyName string `json:"currency-name,omitempty"`
			} `json:"currency-out"`
			CurrencyIn struct {
				CurrencyName string `json:"currency-name,omitempty"`
			} `json:"currency-in"`
			Amount string `json:"amount,omitempty"`
		} `json:"parameters"`
	} `json:"result"`
	SessionID string `json:"sessionId"`
}

func (r *Response) GetSessionID() string {
	return r.SessionID
}
