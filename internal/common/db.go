package common

import (
	"eduVix_backend/internal/config"
	"eduVix_backend/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error

	DB, err = gorm.Open(postgres.Open(config.Cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connect error:", err)
	}

	// Создание ENUM типов PostgreSQL, если ещё нет
	DB.Exec(`
        DO $$ BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
                CREATE TYPE user_role AS ENUM ('user','admin');
            END IF;
        END$$;
    `)

	DB.Exec(`
        DO $$ BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'product_type') THEN
                CREATE TYPE product_type AS ENUM ('P','FD');
            END IF;
        END$$;
    `)

	DB.Exec(`
        DO $$ BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
                CREATE TYPE order_status AS ENUM ('in_progress','cancelled','failed','paid','completed');
            END IF;
        END$$;
    `)

	// AutoMigrate всех моделей
	if err := DB.AutoMigrate(
		// Пользователи и токены
		&models.User{},
		&models.RefreshToken{},
		&models.EmailConfirmation{},

		// Товары
		&models.Processor{},
		&models.FlashDrive{},
		&models.Image{},

		// Корзина и заказы
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
	); err != nil {
		log.Fatal("DB migrate error:", err)
	}

	log.Println("✅ DB initialized and migrated!")
}
