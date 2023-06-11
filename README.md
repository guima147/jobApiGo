# jobApiGo
api de crud simples de itens e itens de um pedido em Go<br>
O framework utilizado é o Gin , para instala-lo abra  o diretorio via termial e execute o seguinte comando:<br>
go get -u github.com/gin-gonic/gin <br>
ao rodar o projeto :  go run api.go  sera instalado dependecias do go-sqlite3

# metodos de exemplo
Criar um pedido:

Método: POST
URL: http://localhost:8080/api/v1/pedido
Cabeçalhos:
Content-Type: application/json]

Método: GET
URL: http://localhost:8080/api/v1/pedido/123

Obter todos os pedidos:
Método: GET
URL: http://localhost:8080/api/v1/pedido

Adicionar um item ao pedido:

Método: POST
URL: http://localhost:8080/api/v1/pedido/123/item
Cabeçalhos:Content-Type: application/json

Obter um item do pedido pelo número e índice:

Método: GET
URL: http://localhost:8080/api/v1/pedido/123/item/1

Obter todos os itens do pedido pelo número:

Método: GET
URL: http://localhost:8080/api/v1/pedido/123/item

Obter todos os itens do pedido pelo produto:

Método: GET
URL: http://localhost:8080/api/v1/pedido/item?produto=Camiseta
