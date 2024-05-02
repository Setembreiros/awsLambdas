update:
	go mod tidy

# Este build é para o environment de windows
# en windows é preciso executalo dende powershell estabelecendo primeiro as propiedades de contorna
# $env:GOOS="linux"
# $env:GOARCH="arm64"
# $env:CGO_ENABLED="0"
build-windows: update
	go build -gcflags "all=-N -l" -o deployments/preSignUp/newUserValidator/build/bootstrap cmd/preSignUp/newUserValidator/main.go
	go build -gcflags "all=-N -l" -o deployments/postConfirmation/newUserEventKafkaSender/build/bootstrap cmd/postConfirmation/newUserEventKafkaSender/main.go

# Este build é para o environment de linux
build-linux: update
	GOARCH=arm64 GOOS=linux go build -gcflags "all=-N -l" -o deployments/preSignUp/newUserValidator/build/bootstrap cmd/preSignUp/newUserValidator/main.go
	GOARCH=arm64 GOOS=linux go build -gcflags "all=-N -l" -o deployments/postConfirmation/newUserEventKafkaSender/build/bootstrap cmd/postConfirmation/newUserEventKafkaSender/main.go

test:
	go test ./tests/...

package-linux:
	zip -j deployments/preSignUp/newUserValidator/build/bootstrap.zip deployments/preSignUp/newUserValidator/build/bootstrap
	zip -j deployments/postConfirmation/newUserEventKafkaSender/build/bootstrap.zip deployments/postConfirmation/newUserEventKafkaSender/build/bootstrap

# O comando "make package-windows" non funciona, pero executar o comando "Compress-Archive..." dende powershell si funciona
package-windows:
	Compress-Archive -Path "deployments/preSignUp/newUserValidator/build/bootstrap" -DestinationPath deployments/preSignUp/newUserValidator/build/bootstrap.zip
	Compress-Archive -Path "deployments/postConfirmation/newUserEventKafkaSender/build/bootstrap" -DestinationPath deployments/postConfirmation/newUserEventKafkaSender/build/bootstrap.zip