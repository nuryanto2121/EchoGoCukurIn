package models

type SsSequenceNo struct {
	SequenceID int    `json:"sequence_id" gorm:"primary_key;auto_increment:true"`
	SequenceCd string `json:"sequence_cd" gorm:"type:varchar(20)"`
	Prefix     string `json:"prefix" gorm:"type:varchar(10)"`
	SeqNo      int    `json:"seq_no" gorm:"type:integer"`
	LogUserTable
}
