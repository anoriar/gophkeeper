package enum

type EntryType string

const (
	Login EntryType = "login"
	Card  EntryType = "card"
)

var AllEntryTypes = []EntryType{Login, Card}

func IsEntryType(value string) bool {
	for _, v := range AllEntryTypes {
		if string(v) == value {
			return true
		}
	}
	return false
}
