package orderrepo

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/orders"
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
	err = db.AutoMigrate(&OrderDTO{})
	require.NoError(t, err)

	// Очистка выполняется после завершения теста
	t.Cleanup(func() {
		postgresContainer.Terminate(ctx)
	})

	return ctx, db, nil
}

func Test_OrderRepositoryShouldCanAddOrder(t *testing.T) {
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	require.NoError(t, err)

	// Создаем репозиторий
	orderRepository, err := NewRepository(db)
	require.NoError(t, err)

	// Вызываем Add

	location, err := kernel.MaxLocation()
	require.NoError(t, err)
	orderAggregate, err := orders.NewOrder(uuid.New(), location)
	err = orderRepository.Add(ctx, orderAggregate)
	require.NoError(t, err)

	// Считываем данные из БД
	var orderFromDb OrderDTO
	err = db.First(&orderFromDb, "id = ?", orderAggregate.ID()).Error
	assert.NoError(t, err)

	// Проверяем эквивалентность
	require.Equal(t, orderAggregate.ID(), orderFromDb.ID)
	require.Equal(t, orderAggregate.Status(), orderFromDb.Status)
}
