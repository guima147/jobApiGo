package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Pedido struct {
	ID      int    `json:"id"`
	Numero  string `json:"numero"`
	Cliente string `json:"cliente"`
}

type ItemPedido struct {
	ID      int     `json:"id"`
	Numero  string  `json:"numero"`
	Indice  int     `json:"indice"`
	SKU     string  `json:"sku"`
	Produto string  `json:"produto"`
	Preco   float64 `json:"preco"`
	Qtd     int     `json:"qtd"`
}

func main() {
	db, err := sql.Open("sqlite3", "./pedidos.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable(db)
	//definindo os  endpoints
	router := gin.Default()

	router.POST("/api/v1/pedido", createPedido(db))
	router.GET("/api/v1/pedido/:numero", getPedidoByNumero(db))
	router.GET("/api/v1/pedido", getAllPedidos(db))
	router.POST("/api/v1/pedido/:numero/item", createItemPedido(db))
	router.GET("/api/v1/pedido/:numero/item/:indice", getItemPedidoByIndice(db))
	router.GET("/api/v1/pedido/:numero/item", getAllItemPedidos(db))
	router.GET("/api/v1/pedido/item", getItemPedidosByProduto(db))

	router.Run(":8080")
}

// criando tabelas no sqlite
func createTable(db *sql.DB) {
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS pedidos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			numero TEXT NOT NULL UNIQUE,
			cliente TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS itens_pedidos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			numero TEXT NOT NULL,
			indice INTEGER NOT NULL,
			sku TEXT NOT NULL,
			produto TEXT NOT NULL,
			preco REAL NOT NULL,
			qtd INTEGER NOT NULL,
			FOREIGN KEY (numero) REFERENCES pedidos(numero)
		);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

// insere pedidos
func createPedido(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var pedido Pedido
		if err := c.ShouldBindJSON(&pedido); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stmt, err := db.Prepare("INSERT INTO pedidos(numero, cliente) VALUES(?, ?)")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		result, err := stmt.Exec(pedido.Numero, pedido.Cliente)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		pedido.ID = int(lastInsertID)

		c.JSON(http.StatusCreated, pedido)
	}
}

// lista pedidos por numero do pedido
func getPedidoByNumero(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		numero := c.Param("numero")

		var pedido Pedido
		err := db.QueryRow("SELECT id, numero, cliente FROM pedidos WHERE numero = ?", numero).Scan(&pedido.ID, &pedido.Numero, &pedido.Cliente)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Pedido não encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, pedido)
	}
}

// lista todos os pedidos
func getAllPedidos(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT id, numero, cliente FROM pedidos")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var pedidos []Pedido
		for rows.Next() {
			var pedido Pedido
			err := rows.Scan(&pedido.ID, &pedido.Numero, &pedido.Cliente)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			pedidos = append(pedidos, pedido)
		}

		c.JSON(http.StatusOK, pedidos)
	}
}

// insert dos itens do pedido
func createItemPedido(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		numero := c.Param("numero")

		var itemPedido ItemPedido
		if err := c.ShouldBindJSON(&itemPedido); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		stmt, err := db.Prepare("INSERT INTO itens_pedidos(numero, indice, sku, produto, preco, qtd) VALUES(?, ?, ?, ?, ?, ?)")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer stmt.Close()

		result, err := stmt.Exec(numero, itemPedido.Indice, itemPedido.SKU, itemPedido.Produto, itemPedido.Preco, itemPedido.Qtd)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		lastInsertID, err := result.LastInsertId()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		itemPedido.ID = int(lastInsertID)
		itemPedido.Numero = numero

		c.JSON(http.StatusCreated, itemPedido)
	}
}

// lista itens do pedido pelo numero e indice
func getItemPedidoByIndice(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		numero := c.Param("numero")
		indice, err := strconv.Atoi(c.Param("indice"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var itemPedido ItemPedido
		err = db.QueryRow("SELECT id, numero, indice, sku, produto, preco, qtd FROM itens_pedidos WHERE numero = ? AND indice = ?", numero, indice).Scan(&itemPedido.ID, &itemPedido.Numero, &itemPedido.Indice, &itemPedido.SKU, &itemPedido.Produto, &itemPedido.Preco, &itemPedido.Qtd)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Item de pedido não encontrado"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, itemPedido)
	}
}

// lista todos os itens do pedido por numero do pedido
func getAllItemPedidos(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		numero := c.Param("numero")

		rows, err := db.Query("SELECT id, numero, indice, sku, produto, preco, qtd FROM itens_pedidos WHERE numero = ?", numero)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var itemPedidos []ItemPedido
		for rows.Next() {
			var itemPedido ItemPedido
			err := rows.Scan(&itemPedido.ID, &itemPedido.Numero, &itemPedido.Indice, &itemPedido.SKU, &itemPedido.Produto, &itemPedido.Preco, &itemPedido.Qtd)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			itemPedidos = append(itemPedidos, itemPedido)
		}

		c.JSON(http.StatusOK, itemPedidos)
	}
}

// lita instens do pedido por  nome do produto
func getItemPedidosByProduto(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		produto := c.Query("produto")

		rows, err := db.Query("SELECT id, numero, indice, sku, produto, preco, qtd FROM itens_pedidos WHERE produto = ?", produto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var itemPedidos []ItemPedido
		for rows.Next() {
			var itemPedido ItemPedido
			err := rows.Scan(&itemPedido.ID, &itemPedido.Numero, &itemPedido.Indice, &itemPedido.SKU, &itemPedido.Produto, &itemPedido.Preco, &itemPedido.Qtd)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			itemPedidos = append(itemPedidos, itemPedido)
		}

		c.JSON(http.StatusOK, itemPedidos)
	}
}
