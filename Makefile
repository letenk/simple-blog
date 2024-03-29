# Run App
server:
	go run main.go

# Running test
test:
	go test -v -cover ./...

# Migrate up
migrateup:
	migrate -path migration -database "mysql://secret:secret@tcp(127.0.0.1:3306)/simple_blog" --verbose up

# Migrate up for table test
migrateuptest:
	migrate -path migration -database "mysql://secret:secret@tcp(127.0.0.1:3306)/simple_blog_test" --verbose up

# migrate down
migratedown:
	migrate -path migration -database "mysql://secret:secret@tcp(127.0.0.1:3306)/simple_blog" --verbose down

# migrate down for table test
migratedowntest:
	migrate -path migration -database "mysql://secret:secret@tcp(127.0.0.1:3306)/simple_blog_test" --verbose down

.PHONY: server test migrateup migratedown migrateuptest
