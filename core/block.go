package core

// import "github.com/ethereum/go-ethereum/core/types"

type (
	Block struct {
		BlockNumber       int64               `gorm:"primaryKey;column:id"` // Block Number
		Timestamp         uint64              `gorm:"column:timestamp"`
		AccountOperations []*AccountOperation `gorm:"foreign:block_num"`
	}
)

func (Block) TableName() string {
	return "blocks"
}

func (b *Block) AddAccountOperation(accountOperation *AccountOperation) {
	b.AccountOperations = append(b.AccountOperations, accountOperation)
}
