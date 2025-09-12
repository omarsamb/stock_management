package database

import (
	"database/sql"
	"fmt"
	"os"
	"stock_management/backend/models"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	*sql.DB
}

// Connect establishes a connection to the PostgreSQL database
func Connect() (*DB, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "stock_management"
	}

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// CreateTenant creates a new tenant
func (db *DB) CreateTenant(name string) (*models.Tenant, error) {
	tenant := &models.Tenant{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	query := `INSERT INTO tenants (id, name, created_at) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, tenant.ID, tenant.Name, tenant.CreatedAt)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

// CreateUser creates a new user
func (db *DB) CreateUser(email, password string, tenantID uuid.UUID, role string) (*models.User, error) {
	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(passwordHash),
		TenantID:     tenantID,
		Role:         role,
		CreatedAt:    time.Now(),
	}

	query := `INSERT INTO users (id, email, password_hash, tenant_id, role, created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(query, user.ID, user.Email, user.PasswordHash, user.TenantID, user.Role, user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, tenant_id, role, created_at FROM users WHERE email = $1`
	
	err := db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.TenantID, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ValidatePassword checks if the provided password matches the user's password
func (db *DB) ValidatePassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

// GetUsersByTenant retrieves all users for a given tenant
func (db *DB) GetUsersByTenant(tenantID uuid.UUID) ([]*models.User, error) {
	query := `SELECT id, email, password_hash, tenant_id, role, created_at FROM users WHERE tenant_id = $1 ORDER BY created_at DESC`
	
	rows, err := db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.TenantID, &user.Role, &user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}