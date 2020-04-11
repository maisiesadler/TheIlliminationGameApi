go vet ./functions/setup/get
go vet ./functions/setup/p_id/get
go vet ./functions/setup/p_id/post
go vet ./functions/setup/put
go vet ./functions/user/get
go vet ./functions/user/post

go test ./functions/setup/get
go test ./functions/setup/p_id/get
go test ./functions/setup/p_id/post
go test ./functions/setup/put
go test ./functions/user/get
go test ./functions/user/post