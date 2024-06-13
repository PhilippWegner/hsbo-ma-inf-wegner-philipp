package database

type Repository interface {
	InsertLog(entry LogEntry) error
}
