set CGO_ENABLED=0 
set GOOS=linux 
go build -a -installsuffix cgo -o client-connector .
docker build -t client-connector -f Dockerfile .
docker run -it client-connector
