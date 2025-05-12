DB_CONTAINER_NAME = bank_system_db
DB_CONTAINER_PORT = 55432
DB_USER_NAME = user_db
DB_USER_PASS = 1213455

run: db_run build go_run

build:
	go mod tidy

db_restart: db_down db_run

go_run:
	go run main.go

go_update:
	go mod tidy

db_run:
	docker run -d \
	  --name ${DB_CONTAINER_NAME} \
	  -e POSTGRES_DB=${DB_CONTAINER_NAME} \
	  -e POSTGRES_USER=${DB_USER_NAME} \
	  -e POSTGRES_PASSWORD=${DB_USER_PASS} \
	  -p ${DB_CONTAINER_PORT}:5432 \
	  -v ./.db:/var/lib/postgresql/data \
	  --restart=always \
	  postgres:15.3

db_down:
	docker stop ${DB_CONTAINER_NAME}
	docker rm ${DB_CONTAINER_NAME}

unlock:
	sudo chown -R ${USER}:${USER} ./.db
	chmod 775 ./.db