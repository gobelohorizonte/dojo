# FileServer

FileServer, um protótipo que tem como objetivo apresentar como funciona um storage para uploads de arquivos. É um exemplo de como um server irá comportar para armazenar arquivos em seu storage ou em buckets.
Este projeto é puramente didâtico, todo construído com stdlib nativa de Go, usamos alguns libs externas para tratarmos JWT e tollbooth/limiter.

## Banco de Dados

Estamos usando Postgresql para simulções, uma tabela de login, pasta, file e ambiente de trabalho.
A versão do postgresql é psql (PostgreSQL) 9.6.9.

## Endpoints

login, / create / user, / ping, / hello, / upload, / download serão os principais endpoints que serão gerenciados.

O upload é feito sempre no servidor local, há uma tabela com as configurações para determinar como esse upload deve ocorrer.

Tudo está registrado no banco de dados Postgresql.

A autenticação para upload de upload só é feita com token de acesso, disponibilizada no login do aplicativo.

O download verifica se o arquivo é local ou se já foi enviado para a nuvem, se ainda estiver no servidor local, o sistema fará o download localmente, caso contrário, verificará qual servidor em nuvem deve ser baixado da nuvem.

## Dependencies

go get -u github.com/didip/tollbooth

go get -u github.com/lib/pq

go get -u golang.org/x/crypto/bcrypt

go get -u github.com/dgrijalva/jwt-go


## Estrutura do Programa

```go

- fileserver
	- certs
	 	- certs.go

	- handler
		- hello.go
		- login.go
		- ping.go
		- upload.go
		- upload_form_unic.go

	- pkg
		- auth
		- authi
		- cors
		- cryptf
		- gcolor
		- logf
		- pg
			
	- models
		- claim.go
		- responsetoken.go
		- user.go

	- route
		- apiconf.go
		- apishow.go
		- route.go	

	- src
		- fileserver
			- main.go	
			- Dockerfile

		- storage
```

# Create User e Update senha

```pg

$ createuser fileserver -U postgres -O fileserver -E UTF-8 -T template0

$ psql -d template1 -U postgres

template1=# alter user fileserver password '1234';
template1=# \q

```
# Create User e Update senha

```pg

$ createdb fileserver -U postgres -O fileserver -E UTF-8 -T template0

$ psql -d fileserver -f sql/fileserver.sql

```

## Install FileServer

```go

$ cd src/fileserver
$ go install
$ fileserver

```

## Docker FileServer

```docker

$ docker run -d -e PORT_SERVER=5009 -e DB_HOST_1=localhost -p 5009:5009 jeffotoni/fileserver

```

# A API retorna o status das solicitações

200 - OK. Successful.

400 - Bad request.

401 - Authorization required.

403 - Not allowed.

404 - Not found

409 - Already exists

420 - Rate limited


# Limites da API

Você pode usar a API como quiser. Criamos um limite de taxa de 20.000 solicitações por segundo. Você pode alterar este limite de taxa em api-route se a variável global é chamada de "NewLimiter", você também pode alterar o tamanho máximo de uploads que o servidor permitirá.

Para cada chamada de API, os cabeçalhos HTTP padrão são retornados em detalhes nas estatísticas de uso atuais, incluindo o limite por hora, o número restante de ações e a hora em que a contagem de tempo será reiniciada.

# Ping

O método é público e aberto não precisa de autenticação, para saber se o serviço está online

```sh

curl -X POST localhost:5000/v1/test/ping

OR

curl localhost:5000/v1/test/ping

```


```sh 

curl -X POST localhost:5000/v1/user \
-H "Content-Type: application/json" \
-H "X-Key: ZmlsZXNlcnZlcjIwMThnb2xhbmdiaA==" \
-d @example/createuser.json

```

# Login Form

O login é o manipulador responsável por autenticar a plataforma e retorna um token para que possa acessar os manipuladores privados.

```sh 

[POST] /v1/user/login

```

```sh

curl -X POST localhost:5000/v1/user/login \
-H "X-Key: ZmlsZXNlcnZlcjIwMThnb2xhbmdiaA==" \
-d "user=email@server.com&password=1234"

OR

curl -X POST localhost:5000/v1/user/login \
-H "X-Key: ZmlsZXNlcnZlcjIwMThnb2xhbmdiaA==" \
-d "user=email@server.com" \
-d "password=1234"

OR

curl -X POST localhost:5000/v1/user/login \
-H "X-Key: ZmlsZXNlcnZlcjIwMThnb2xhbmdiaA==" \
-d '{"user":"email@server.com","password":"1234"}'

```

# Response

Após o login, você receberá um token para usar em todas as chamadas


```sh
msg
{	
	"status":"ok",
	
	"msg":"success",
	
	"token":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiamVmZi5vdG9uaUBnbWFpbC5jb20iLCJ1aWQiOiI0ZjQ0NzgxMTAzODkwNjA5ZmY1MDdmNjIzMTdlOGExMGFiMDc4ZjFmIiwidWlkd2tzIjoiMDE2NTcxNmNiNzBmMWM1N2ZhMzhhZGY5MGI1Y2QyMTdmNDE1NTJhNyIsImV4cCI6MTUzNzg0MzU5NCwiaXNzIjoiand0IEZpbGVTZXJ2ZXIifQ.XD6eCPDklwWw1WR7fViCSksRsMjsERZuHKRMIhQvgn63pwXTLBR6NvRAmc2JUyDmovEyDzo05LGDg3DadI-b5J3Su29OcZpKMF-aSqSBLzmuQ7Grl1z7EO-IircrkVkZ-Xuk8Ur3VBRhQt0Z_b7MwMgbClS2-SJGX5pPUi0SrZk",
	
	"expires":"2018-09-25"
}

```

# Hello 

O método Hello é um manipulador para testar seu token, ele receberá o token como um parâmetro, se o token estiver correto ele retornará um hello, caso contrário ele retornará um erro.


```sh 

[POST] /v1/test/hello

```

```sh

curl -X POST localhost:5000/v1/test/hello \
-H "Authorization: Bearer {token}" \

```

```sh

curl -X PUT localhost:5000/v1/user/enable \
-H "Authorization: Bearer {tokem}" \
-d "user=mail@mail.com"

```


# Upload to a File

O nome do campo do formulário é "file", para uploads.
Este manipulador permitirá que você faça upload usando o formulário para isso.

```sh 

[POST] /v1/file/upload

```

```sh

curl -X POST localhost:5000/v1/file/upload \
-H "Authorization: Bearer {token}" \
--form "file=@namefile.pdf"

```

# Upload to Multiple Files

O nome do campo do formulário é "file", para uploads.
Esse manipulador permitirá que você faça vários envios ao mesmo tempo usando o formulário para isso.

```sh 

[POST] /v1/file/upload

```

```sh

curl -X POST localhost:5000/v1/file/upload \
-H "Authorization: Bearer {token}" \
--form "file[]=@namefile1.pdf" \
--form "file[]=@namefile2.pdf" \
--form "file[]=@namefile3.pdf"

```

# Upload Binary

Este manipulador permitirá que você faça o upload do arquivo como um binário, para que isso funcione corretamente, você terá que passar pelo cabeçalho o nome do arquivo.

```sh 

[POST] /v1/file/upload

```

```sh

curl -X POST localhost:5000/v1/file/upload \
-H "Authorization: Bearer {token}" \
-H "Accept: binary/octet-stream" \
-H "Content-Type: binary/octet-stream" \
-H "Name-File: namefile.ext" \
--data-binary "@filetoupload"

```

# Download

Esse manipulador é responsável pelo download do arquivo que está no servidor.


```sh 

[GET] /v1/file/download

```

```sh

curl -X GET -o namefile.ext localhost:5000/v1/file/download \
-H "Content-Type: application/json" \ 
-H "Authorization: Bearer {token}" \ 
-d '{"id":"uuidfile"}'

OR

curl -X GET localhost:5000/v1/file/download \
-H "Authorization: Bearer {token}" \
-H "Content-Type: application/json" \ 
-d '{"id":"uuidfile"}' \ > namefile.ext

```