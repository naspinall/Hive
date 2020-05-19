package models

import (
	"log"
)

type AuditAccessType string

func LogCreate(userID uint, resource string) {
	auditLog(userID, resource, AuditAccessType("CREATE"))
}

func LogGet(userID uint, resource string) {
	auditLog(userID, resource, AuditAccessType("GET"))
}

func LogDelete(userID uint, resource string) {
	auditLog(userID, resource, AuditAccessType("DELETE"))
}

func LogUpdate(userID uint, resource string) {
	auditLog(userID, resource, AuditAccessType("UPDATE"))
}

func auditLog(userID uint, resource string, accessType AuditAccessType) {
	log.Printf(`%s UserID: %d Resource: %s`, accessType, userID, resource)
}
