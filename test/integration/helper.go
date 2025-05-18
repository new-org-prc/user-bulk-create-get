package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"sika/config"
	"sika/pkg/storage"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestDB struct {
	*gorm.DB
}

func SetupTestDB(t *testing.T) *TestDB {
	// Get the project root directory
	projectRoot, err := os.Getwd()
	require.NoError(t, err)
	projectRoot = filepath.Dir(filepath.Dir(projectRoot)) // Go up two levels from test/integration

	// Read test configuration
	testConfig, err := config.ReadStandard(filepath.Join(projectRoot, "test", "config.yaml"))
	require.NoError(t, err)

	// Connect to the database
	db, err := storage.NewPostgresGormConnection(testConfig.DB)
	require.NoError(t, err)

	// Run migrations
	err = storage.Migrate(db)
	require.NoError(t, err)

	return &TestDB{db}
}

func (db *TestDB) Cleanup(t *testing.T) {
	// Clear all data from the database
	err := db.Exec("DELETE FROM addresses").Error
	require.NoError(t, err)
	err = db.Exec("DELETE FROM users").Error
	require.NoError(t, err)
}

func (db *TestDB) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// WaitForDB waits for the database to be ready
func WaitForDB(t *testing.T, dbConfig config.DB) {
	maxRetries := 5
	retryInterval := time.Second * 2

	for i := 0; i < maxRetries; i++ {
		db, err := storage.NewPostgresGormConnection(dbConfig)
		if err == nil {
			sqlDB, err := db.DB()
			if err == nil {
				sqlDB.Close()
				return
			}
		}
		time.Sleep(retryInterval)
	}
	t.Fatal("Failed to connect to database after multiple retries")
}
