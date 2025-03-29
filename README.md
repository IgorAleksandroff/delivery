# Шаблоны
```
https://gitlab.com/microarch-ru/ddd-in-practice/templates/
```
```
https://gitlab.com/microarch-ru/ddd-in-practice/infrastructure.git
```

# OpenApi (генерация HTTP сервера)
```
oapi-codegen -config configs/server.cfg.yaml https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/delivery/contracts/openapi.yml
```

# БД
```
https://pressly.github.io/goose/installation/
```

# Запросы к БД
```
-- Выборки
SELECT * FROM public.couriers;
SELECT * FROM public.transports;
SELECT * FROM public.orders;

SELECT * FROM public.outbox;

-- Очистка БД (все кроме справочников)
DELETE FROM public.couriers;
DELETE FROM public.transports;
DELETE FROM public.orders;
DELETE FROM public.outbox;

-- Добавить курьеров
    
-- Пеший
INSERT INTO public.transports(
    id, name, speed)
VALUES ('921e3d64-7c68-45ed-88fb-97ceb8148a7e', 'Пешком', 1);
INSERT INTO public.couriers(
    id, name, transport_id, location_x, location_y, status)
VALUES ('bf79a004-56d7-4e5f-a21c-0a9e5e08d10d', 'Пеший', '921e3d64-7c68-45ed-88fb-97ceb8148a7e', 1, 3, 'free');

-- Вело
INSERT INTO public.transports(
    id, name, speed)
VALUES ('b96a9d83-aefa-4d06-99fb-e630d17c3868', 'Велосипед', 2);
INSERT INTO public.couriers(
    id, name, transport_id, location_x, location_y, status)
VALUES ('db18375d-59a7-49d1-bd96-a1738adcee93', 'Вело', 'b96a9d83-aefa-4d06-99fb-e630d17c3868', 4,5, 'free');

-- Авто
INSERT INTO public.transports(
    id, name, speed)
VALUES ('c24d3116-a75c-4a4b-9b22-1a7dc95a8c79', 'Машина', 3);
INSERT INTO public.couriers(
    id, name, transport_id, location_x, location_y, status)
VALUES ('407f68be-5adf-4e72-81bc-b1d8e9574cf8', 'Авто', 'c24d3116-a75c-4a4b-9b22-1a7dc95a8c79', 7,9, 'free');     
```

# gRPC Client
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
export PATH="$PATH:$(go env GOPATH)/bin"
protoc --proto_path=src --go_out=out --go_opt=paths=source_relative foo.proto bar/baz.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helloworld/helloworld.proto
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative https://gitlab.com/microarch-ru/ddd-in-practice/system-design/-/raw/main/services/discount/contracts/grpc.proto?ref_type=heads
```

# Тестирование
```
mockery --all --case=underscore
```

# Документация используемых библилиотек
* [Goose] (https://github.com/pressly/goose)
* [Oapi-codegen] (https://github.com/oapi-codegen/oapi-codegen)
* [Protobuf] (https://protobuf.dev/reference/go/go-generated/)
* [gRPC] (https://grpc.io/docs/languages/go/)
* [Mockery] (https://vektra.github.io/mockery/latest/)