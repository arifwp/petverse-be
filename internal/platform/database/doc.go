// Package database owns SQL connection setup, transaction helpers, and migrations.
//
// Keep database-specific code here. Domain modules should depend on repository
// interfaces, not concrete drivers.
package database
