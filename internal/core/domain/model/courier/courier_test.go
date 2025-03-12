package courier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
)

func TestCourier_StepsToOrder(t *testing.T) {
	type fields struct {
		name      string
		transport *Transport
		location  kernel.Location
	}

	type args struct {
		orderLocation kernel.Location
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantSteps int
	}{
		{
			name: "Same location - zero steps",
			fields: fields{
				name:      "Тестовый курьер",
				transport: MustNewTransport("Велосипед", 2),
				location:  kernel.MustNewLocation(5, 5),
			},
			args: args{
				orderLocation: kernel.MustNewLocation(5, 5),
			},
			wantSteps: 0,
		},
		{
			name: "Horizontal movement - bike speed 2",
			fields: fields{
				name:      "Тестовый курьер",
				transport: MustNewTransport("Велосипед", 2),
				location:  kernel.MustNewLocation(1, 1),
			},
			args: args{
				orderLocation: kernel.MustNewLocation(1, 9),
			},
			wantSteps: 4,
		},
		{
			name: "Vertical movement - car speed 3",
			fields: fields{
				name:      "Тестовый курьер",
				transport: MustNewTransport("Авто", 3),
				location:  kernel.MustNewLocation(1, 1),
			},
			args: args{
				orderLocation: kernel.MustNewLocation(10, 1),
			},
			wantSteps: 3,
		},
		{
			name: "Diagonal movement - foot speed 1",
			fields: fields{
				name:      "Тестовый курьер",
				transport: MustNewTransport("Пешком", 1),
				location:  kernel.MustNewLocation(1, 1),
			},
			args: args{
				orderLocation: kernel.MustNewLocation(4, 4),
			},
			wantSteps: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := MustNewCourier(tt.fields.name, tt.fields.transport, tt.fields.location)

			gotSteps, err := c.StepsToOrder(tt.args.orderLocation)
			require.NoError(t, err)
			assert.Equal(t, tt.wantSteps, gotSteps)
		})
	}
}
