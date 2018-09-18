@REM git clone https:/github.com/go-swagger/go-swagger
@REM cd go-swagger\go-swagger\cmd\swagger
@REM go build
@REM go install

@REM swagger generate spec -o ./swagger.json
swagger generate server -A client-connector
