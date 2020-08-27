mkdir bin

go build -o bin/game/finished/get ./functions/game/finished/get
go build -o bin/game/get ./functions/game/get
go build -o bin/game/p_id/get ./functions/game/p_id/get
go build -o bin/game/p_id/post ./functions/game/p_id/post
go build -o bin/game/put ./functions/game/put
go build -o bin/image/get ./functions/image/get
go build -o bin/setup/available/get ./functions/setup/available/get
go build -o bin/setup/get ./functions/setup/get
go build -o bin/setup/p_id/get ./functions/setup/p_id/get
go build -o bin/setup/p_id/post ./functions/setup/p_id/post
go build -o bin/setup/put ./functions/setup/put
go build -o bin/user/get ./functions/user/get
go build -o bin/user/post ./functions/user/post