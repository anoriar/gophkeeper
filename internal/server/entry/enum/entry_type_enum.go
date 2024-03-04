package enum

type EntryType string

const (
	Login EntryType = "login"
	Card  EntryType = "card"
)

var AllEntryTypes = []EntryType{Login, Card}

func IsEntryType(value EntryType) bool {
	for _, v := range AllEntryTypes {
		if v == value {
			return true
		}
	}
	return false
}
