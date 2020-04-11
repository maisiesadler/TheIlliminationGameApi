mkdir bin

go build -o bin/setup/get ./functions/setup/get
go build -o bin/setup/p_id/get ./functions/setup/p_id/get
go build -o bin/setup/p_id/post ./functions/setup/p_id/post
go build -o bin/setup/put ./functions/setup/put
go build -o bin/user/get ./functions/user/get
go build -o bin/user/post ./functions/user/post