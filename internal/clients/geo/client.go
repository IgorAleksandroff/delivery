package geo

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/IgorAleksandroff/delivery/internal/core/domain/model/kernel"
	"github.com/IgorAleksandroff/delivery/internal/pkg/errs"
	pb "github.com/IgorAleksandroff/delivery/pkg/clients/geo/geosrv/geopb"
)

type Client struct {
	conn     *grpc.ClientConn
	pbClient pb.GeoClient
	timeout  time.Duration
}

func NewClient(host string) (*Client, error) {
	if host == "" {
		return nil, errs.NewValueIsRequiredError("host")
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	pbClient := pb.NewGeoClient(conn)

	return &Client{
		conn:     conn,
		pbClient: pbClient,
		timeout:  5 * time.Second,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetGeolocation(ctx context.Context, street string) (kernel.Location, error) {
	// Формируем запрос
	req := &pb.GetGeolocationRequest{
		Street: street,
	}

	// Делаем запрос
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	resp, err := c.pbClient.GetGeolocation(ctx, req)
	if err != nil {
		return kernel.Location{}, err
	}

	// Создаем и возвращаем VO Geo
	location, err := kernel.NewLocation(int(resp.Location.X), int(resp.Location.Y))
	if err != nil {
		return kernel.Location{}, err
	}
	return location, nil
}
