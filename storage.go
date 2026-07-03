package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	DeleteAccount(int) error
	GetAccounts() ([]*Account, error)
	UpdateAccount(*Account) error
	GetAccountByNumber(int64) (*Account, error)
	Transfer(fromID, toID int, amount int64) (*TransferResult, error)
	RedeemCoupon(accountID int, couponCode string) error
	CreatePendingAccount(acc *PendingAccount) error
	GetPendingAccountByCode(code string) (*PendingAccount, error)
	DeletePendingAccount(id int) error
	GetTransactionHistory(accountID, limit, offset int, txType string, month int) ([]*TransactionRecord, error)
	CreatePendingProfileUpdate(ppu *PendingProfileUpdate) error
	GetPendingProfileUpdate(accountID int, newEmail, code string) (*PendingProfileUpdate, error)
	DeletePendingProfileUpdate(id int) error
	UpdateAccountName(accountID int, firstName, lastName string) error
	UpdateAccountEmail(accountID int, newEmail string) error
	UpdateAccountPassword(accountID int, encryptedPassword string) error
	GetAccountByEmail(email string) (*Account, error)
	CreateWebAuthnCredential(cred *WebAuthnCredential, accountID int) error
	GetWebAuthnCredentialsByAccountID(accountID int) ([]*WebAuthnCredential, error)
	CreateWebAuthnCredentialTable() error
}

func (s *PostgresStore) GetAccountByEmail(email string) (*Account, error) {
	row := s.db.QueryRow(`
		SELECT id, first_name, last_name, number, COALESCE(email, ''), COALESCE(encrypted_password, ''), balance, created_at
		FROM account
		WHERE email = $1`, email)

	acc := new(Account)
	err := row.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.Email,
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account with email [%s] not found", email)
	}
	if err != nil {
		return nil, err
	}

	return acc, nil
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(cfg *Config) (*PostgresStore, error) {
	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) CreateWebAuthnCredentialTable() error {
	query := `CREATE TABLE IF NOT EXISTS webauthn_credentials (
		id BYTEA PRIMARY KEY,
		public_key BYTEA NOT NULL,
		attestation_type TEXT NOT NULL,
		transport TEXT[] NOT NULL,
		flags JSONB NOT NULL,
		authenticator JSONB NOT NULL,
		account_id INTEGER NOT NULL REFERENCES account(id)
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateWebAuthnCredential(cred *WebAuthnCredential, accountID int) error {
	flagsJSON, err := json.Marshal(cred.Flags)
	if err != nil {
		return err
	}

	authenticatorJSON, err := json.Marshal(cred.Authenticator)
	if err != nil {
		return err
	}

	query := `INSERT INTO webauthn_credentials
	(id, public_key, attestation_type, transport, flags, authenticator, account_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = s.db.Exec(
		query,
		cred.ID,
		cred.PublicKey,
		cred.AttestationType,
		pq.Array(cred.Transport),
		flagsJSON,
		authenticatorJSON,
		accountID,
	)
	return err
}

func (s *PostgresStore) GetWebAuthnCredentialsByAccountID(accountID int) ([]*WebAuthnCredential, error) {
	rows, err := s.db.Query(`
		SELECT id, public_key, attestation_type, transport, flags, authenticator
		FROM webauthn_credentials
		WHERE account_id = $1`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	creds := []*WebAuthnCredential{}
	for rows.Next() {
		cred := new(WebAuthnCredential)
		var transport []string
		var flagsJSON []byte
		var authenticatorJSON []byte
		if err := rows.Scan(
			&cred.ID,
			&cred.PublicKey,
			&cred.AttestationType,
			pq.Array(&transport),
			&flagsJSON,
			&authenticatorJSON,
		); err != nil {
			return nil, err
		}
		cred.Transport = transport
		if len(flagsJSON) > 0 {
			if err := json.Unmarshal(flagsJSON, &cred.Flags); err != nil {
				return nil, err
			}
		}
		if len(authenticatorJSON) > 0 {
			if err := json.Unmarshal(authenticatorJSON, &cred.Authenticator); err != nil {
				return nil, err
			}
		}
		creds = append(creds, cred)
	}

	return creds, rows.Err()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		number BIGINT,
		encrypted_password VARCHAR(255) NOT NULL DEFAULT '',
		balance BIGINT,
		created_at TIMESTAMP
	);`

	if _, err := s.db.Exec(query); err != nil {
		return err
	}

	// Backward-compatible migration for older local table versions.
	if _, err := s.db.Exec(`ALTER TABLE account ALTER COLUMN encrypted_password SET NOT NULL;`); err != nil {
		return err
	}

	// Add email column if it does not exist yet.
	_, err := s.db.Exec(`ALTER TABLE account ADD COLUMN IF NOT EXISTS email VARCHAR(255);`)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	createdAt := time.Now().UTC()
	query := `INSERT INTO account
	(first_name, last_name, number, email, encrypted_password, balance, created_at)
	values
	($1, $2, $3, $4, $5, $6, $7)
	RETURNING id`

	err := s.db.QueryRow(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Email,
		acc.EncryptedPassword,
		acc.Balance,
		createdAt,
	).Scan(&acc.ID)

	if err != nil {
		return err
	}

	acc.CreatedAt = createdAt

	return nil
}

func (s *PostgresStore) UpdateAccount(acc *Account) error {
	_, err := s.db.Exec(
		`UPDATE account SET first_name = $1, last_name = $2, balance = $3 WHERE id = $4`,
		acc.FirstName, acc.LastName, acc.Balance, acc.ID,
	)
	return err
}

func (s *PostgresStore) Transfer(fromID, toID int, amount int64) (*TransferResult, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var balance int64
	if err := tx.QueryRow(
		"SELECT balance FROM account WHERE id = $1 FOR UPDATE", fromID,
	).Scan(&balance); err != nil {
		return nil, err
	}
	if balance < amount {
		return nil, fmt.Errorf("insufficient funds")
	}

	if _, err := tx.Exec(
		"UPDATE account SET balance = balance - $1 WHERE id = $2", amount, fromID,
	); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(
		"UPDATE account SET balance = balance + $1 WHERE id = $2", amount, toID,
	); err != nil {
		return nil, err
	}

	txID, err := generateTransactionID()
	if err != nil {
		return nil, err
	}

	transfer := &TransferResult{}
	err = tx.QueryRow(
		`INSERT INTO transfer (transaction_id, from_account_id, to_account_id, amount, transaction_type, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, transaction_id, from_account_id, to_account_id, amount, created_at`,
		txID,
		fromID,
		toID,
		amount,
		"transfer",
		"completed",
		time.Now().UTC(),
	).Scan(
		&transfer.ID,
		&transfer.TransactionID,
		&transfer.FromAccountID,
		&transfer.ToAccountID,
		&transfer.Amount,
		&transfer.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Exec("DELETE FROM account WHERE id = $1", id)
	return err
}

func (s *PostgresStore) GetAccountByNumber(number int64) (*Account, error) {
	row := s.db.QueryRow(`
		SELECT id, first_name, last_name, number, COALESCE(email, ''), COALESCE(encrypted_password, ''), balance, created_at
		FROM account
		WHERE number = $1`, number)

	acc := new(Account)
	err := row.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.Email,
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account with number [%d] not found", number)
	}
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	row := s.db.QueryRow(`
		SELECT id, first_name, last_name, number, COALESCE(email, ''), COALESCE(encrypted_password, ''), balance, created_at
		FROM account
		WHERE id = $1`, id)

	acc := new(Account)
	err := row.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.Email,
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account %d not found", id)
	}
	if err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`
		SELECT id, first_name, last_name, number, COALESCE(email, ''), COALESCE(encrypted_password, ''), balance, created_at
		FROM account`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*Account{}
	for rows.Next() {
		acc, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}

	return accounts, rows.Err()
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	acc := new(Account)
	err := rows.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Number,
		&acc.Email,
		&acc.EncryptedPassword,
		&acc.Balance,
		&acc.CreatedAt,
	)

	return acc, err
}

func (s *PostgresStore) CreateTransferTable() error {
	query := `CREATE TABLE IF NOT EXISTS transfer (
		id SERIAL PRIMARY KEY,
		transaction_id VARCHAR(32) NOT NULL UNIQUE,
		from_account_id INT NOT NULL REFERENCES account(id),
		to_account_id INT NOT NULL REFERENCES account(id),
		amount BIGINT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);`
	if _, err := s.db.Exec(query); err != nil {
		return err
	}

	// Backward-compatible migrations for new columns.
	if _, err := s.db.Exec(`ALTER TABLE transfer ADD COLUMN IF NOT EXISTS transaction_type VARCHAR(50) NOT NULL DEFAULT 'transfer';`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`ALTER TABLE transfer ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'completed';`); err != nil {
		return err
	}

	return nil
}

// CreatePendingAccountTable creates the table used to hold accounts pending email verification.
func (s *PostgresStore) CreatePendingAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS pending_account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		email VARCHAR(255) NOT NULL,
		number BIGINT NOT NULL,
		encrypted_password VARCHAR(255) NOT NULL,
		verification_code VARCHAR(6) NOT NULL,
		expires_at TIMESTAMP NOT NULL
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreatePendingAccount(acc *PendingAccount) error {
	query := `INSERT INTO pending_account
		(first_name, last_name, email, number, encrypted_password, verification_code, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`
	return s.db.QueryRow(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Email,
		acc.Number,
		acc.EncryptedPassword,
		acc.VerificationCode,
		acc.ExpiresAt,
	).Scan(&acc.ID)
}

func (s *PostgresStore) GetPendingAccountByCode(code string) (*PendingAccount, error) {
	row := s.db.QueryRow(`
		SELECT id, first_name, last_name, email, number, encrypted_password, verification_code, expires_at
		FROM pending_account
		WHERE verification_code = $1`, code)

	acc := new(PendingAccount)
	err := row.Scan(
		&acc.ID,
		&acc.FirstName,
		&acc.LastName,
		&acc.Email,
		&acc.Number,
		&acc.EncryptedPassword,
		&acc.VerificationCode,
		&acc.ExpiresAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid verification code")
	}
	return acc, err
}

func (s *PostgresStore) DeletePendingAccount(id int) error {
	_, err := s.db.Exec("DELETE FROM pending_account WHERE id = $1", id)
	return err
}

// CreateCouponRedemptionTable creates the table used to track single-use coupon redemptions.
func (s *PostgresStore) CreateCouponRedemptionTable() error {
	query := `CREATE TABLE IF NOT EXISTS coupon_redemption (
		id SERIAL PRIMARY KEY,
		coupon_code VARCHAR(255) NOT NULL,
		account_id INT NOT NULL,
		redeemed_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`
	if _, err := s.db.Exec(query); err != nil {
		return err
	}

	// Keep one redemption row per account before creating a unique index.
	if _, err := s.db.Exec(`
		DELETE FROM coupon_redemption cr
		USING coupon_redemption newer
		WHERE cr.account_id = newer.account_id
		AND cr.id < newer.id;
	`); err != nil {
		return err
	}

	if _, err := s.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS coupon_redemption_account_id_uidx ON coupon_redemption(account_id);`); err != nil {
		return err
	}

	return nil
}

// RedeemCoupon records a coupon redemption (enforcing global single-use) and
// adds 1000 to the account balance, all within a single transaction.
func (s *PostgresStore) RedeemCoupon(accountID int, couponCode string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert the redemption record. The UNIQUE index on account_id guarantees
	// each account can redeem only once.
	_, err = tx.Exec(
		`INSERT INTO coupon_redemption (coupon_code, account_id) VALUES ($1, $2)`,
		couponCode, accountID,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return fmt.Errorf("account already redeemed coupon")
		}
		return err
	}

	_, err = tx.Exec(
		`UPDATE account SET balance = balance + 1000 WHERE id = $1`,
		accountID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetTransactionHistory returns the paginated transaction history for an account.
// Results are filtered by transaction_type when txType is non-empty, and by month
// (1-12) when month is non-zero.  Results are sorted newest-first.
func (s *PostgresStore) GetTransactionHistory(accountID, limit, offset int, txType string, month int) ([]*TransactionRecord, error) {
	args := []any{accountID}
	argIdx := 2

	query := `
		SELECT id, transaction_id, from_account_id, to_account_id, amount, transaction_type, status, created_at
		FROM transfer
		WHERE (from_account_id = $1 OR to_account_id = $1)`

	if txType != "" {
		query += fmt.Sprintf(" AND transaction_type = $%d", argIdx)
		args = append(args, txType)
		argIdx++
	}

	if month != 0 {
		query += fmt.Sprintf(" AND EXTRACT(MONTH FROM created_at) = $%d", argIdx)
		args = append(args, month)
		argIdx++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := []*TransactionRecord{}
	for rows.Next() {
		rec := new(TransactionRecord)
		if err := rows.Scan(
			&rec.ID,
			&rec.TransactionID,
			&rec.FromAccountID,
			&rec.ToAccountID,
			&rec.Amount,
			&rec.TransactionType,
			&rec.Status,
			&rec.CreatedAt,
		); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, rows.Err()
}

// CreatePendingProfileUpdateTable creates the table used to hold pending
// email-change requests awaiting OTP confirmation.
func (s *PostgresStore) CreatePendingProfileUpdateTable() error {
	query := `CREATE TABLE IF NOT EXISTS pending_profile_update (
		id                SERIAL PRIMARY KEY,
		account_id        INT NOT NULL REFERENCES account(id),
		new_email         VARCHAR(255) NOT NULL,
		verification_code VARCHAR(6) NOT NULL,
		expires_at        TIMESTAMP NOT NULL,
		created_at        TIMESTAMP NOT NULL DEFAULT NOW(),
		UNIQUE (account_id, new_email)
	);`
	if _, err := s.db.Exec(query); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS pending_profile_update_code_idx ON pending_profile_update(verification_code);`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS pending_profile_update_expires_idx ON pending_profile_update(expires_at);`); err != nil {
		return err
	}
	return nil
}

// CreatePendingProfileUpdate saves a pending email-change request.  Only one
// active pending update per account is allowed; any existing row for the
// account is deleted first so a fresh request cleanly replaces the previous one.
func (s *PostgresStore) CreatePendingProfileUpdate(ppu *PendingProfileUpdate) error {
	// Delete any existing pending update for this account so a fresh request
	// replaces the previous one instead of failing on the unique constraint.
	if _, err := s.db.Exec(`DELETE FROM pending_profile_update WHERE account_id = $1`, ppu.AccountID); err != nil {
		return err
	}
	query := `INSERT INTO pending_profile_update
		(account_id, new_email, verification_code, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	return s.db.QueryRow(
		query,
		ppu.AccountID,
		ppu.NewEmail,
		ppu.VerificationCode,
		ppu.ExpiresAt,
		ppu.CreatedAt,
	).Scan(&ppu.ID)
}

func (s *PostgresStore) GetPendingProfileUpdate(accountID int, newEmail, code string) (*PendingProfileUpdate, error) {
	row := s.db.QueryRow(`
		SELECT id, account_id, new_email, verification_code, expires_at, created_at
		FROM pending_profile_update
		WHERE account_id = $1 AND new_email = $2 AND verification_code = $3`,
		accountID, newEmail, code)

	ppu := new(PendingProfileUpdate)
	err := row.Scan(&ppu.ID, &ppu.AccountID, &ppu.NewEmail, &ppu.VerificationCode, &ppu.ExpiresAt, &ppu.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("invalid or expired OTP")
	}
	return ppu, err
}

func (s *PostgresStore) DeletePendingProfileUpdate(id int) error {
	_, err := s.db.Exec("DELETE FROM pending_profile_update WHERE id = $1", id)
	return err
}

func (s *PostgresStore) UpdateAccountName(accountID int, firstName, lastName string) error {
	_, err := s.db.Exec(
		`UPDATE account SET first_name = $1, last_name = $2 WHERE id = $3`,
		firstName, lastName, accountID,
	)
	return err
}

func (s *PostgresStore) UpdateAccountEmail(accountID int, newEmail string) error {
	_, err := s.db.Exec(
		`UPDATE account SET email = $1 WHERE id = $2`,
		newEmail, accountID,
	)
	return err
}

func (s *PostgresStore) UpdateAccountPassword(accountID int, encryptedPassword string) error {
	_, err := s.db.Exec(
		`UPDATE account SET encrypted_password = $1 WHERE id = $2`,
		encryptedPassword, accountID,
	)
	return err
}
