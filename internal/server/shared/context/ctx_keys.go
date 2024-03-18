package context

const UserIDContextKey = UserIDContextType("userID")
const TransactionKey = TransactionContextType("transaction")

type UserIDContextType string
type TransactionContextType string
