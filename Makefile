build:
	go build -o ./bin/ws_server cmd/ws_server/hub.go cmd/ws_server/server.go
	go build -o ./bin/findgame_server cmd/findGame_server/find.go 	
	go build -o ./bin/login_server cmd/login_server/login.go 
