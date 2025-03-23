package orderrepo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/courier"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
	"github.com/IgorAleksandroff/delivery/internal/pkg/testutil"
)

func setupTest(t *testing.T) (context.Context, *gorm.DB, error) {
	ctx := context.Background()
	postgresContainer, dsn, err := testutil.StartPostgresContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Подключаемся к БД через Gorm
	db, err := gorm.Open(postgresgorm.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Авто миграция (создаём таблицу)
	err = db.AutoMigrate(&CourierDTO{})
	require.NoError(t, err)

	err = db.AutoMigrate(&TransportDTO{})
	require.NoError(t, err)

	// Очистка выполняется после завершения теста
	t.Cleanup(func() {
		postgresContainer.Terminate(ctx)
	})

	return ctx, db, nil
}

func Test_CourierRepositoryShouldCanAddCourier(t *testing.T) {
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	require.NoError(t, err)

	// Создаем репозиторий
	courierRepository, err := NewRepository(db)
	require.NoError(t, err)

	// Вызываем Add
	location, err := kernel.MaxLocation()
	require.NoError(t, err)
	courierAggregate, err := courier.NewCourier("Велосипедист", "Велосипед", 2, location)
	err = courierRepository.Add(ctx, courierAggregate)
	require.NoError(t, err)

	// Считываем данные из БД
	var courierFromDb CourierDTO
	err = db.First(&courierFromDb, "id = ?", courierAggregate.ID()).Error
	assert.NoError(t, err)

	// Проверяем эквивалентность
	require.Equal(t, courierAggregate.ID(), courierFromDb.ID)
	require.Equal(t, courierAggregate.Status(), courierFromDb.Status)
}
