package model

type CreateReportDTO struct {
  Title string `json:"title"`
  Text  string `json:"text"`
}

type CreateFeedbackDTO struct {
  Text string `json:"feedback_text"`
}

type CloseReportDTO struct {
  ComplaintId string  `json:"complaint_id"`
  Verdict     string  `json:"verdict"`
  Rating      float64 `json:"rating"`
}
