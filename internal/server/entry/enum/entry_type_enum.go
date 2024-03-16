package enum

type EntryType string

const (
	Login EntryType = "login"
	Card  EntryType = "card"
	Text  EntryType = "text"
	Bin   EntryType = "bin"
)

var AllEntryTypes = []EntryType{Login, Card, Text, Bin}

func IsEntryType(value string) bool {
	for _, v := range AllEntryTypes {
		if string(v) == value {
			return true
		}
	}
	return false
}
