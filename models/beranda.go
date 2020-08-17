package models

type Beranda struct {
	ProgressStatus int     `json:"progress_status"`
	FinishStatus   int     `json:"finish_status"`
	CancelStatus   int     `json:"cancel_status"`
	IncomePrice    float32 `json:"income_price"`
}

type BerandaList struct {
	BarberID   int     `json:"barber_id"`
	BarberName string  `json:"barber_name"`
	FileID     int     `json:"file_id" `
	FileName   string  `json:"file_name"`
	FilePath   string  `json:"file_path"`
	FileType   string  `json:"file_type"`
	Price      float32 `json:"price" `
}
