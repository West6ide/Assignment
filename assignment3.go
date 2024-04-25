package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
	"time"
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       int
}

var (
	ctx         context.Context
	redisClient *redis.Client
	db          *sqlx.DB
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "20030625"
	dbname   = "first_db"
)

func init() {
	ctx = context.Background()

	// Подключение к Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // нет установленного пароля
		DB:       0,  // использовать базу данных по умолчанию
	})

	// Подключение к PostgreSQL
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sqlx.Connect("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
}

func getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	productIDStr := params["id"]
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Попытка получить данные о продукте из кэша Redis
	cachedProduct, err := redisClient.Get(ctx, "product:"+productIDStr).Result()
	if err == nil {
		// Возвращаем данные о продукте из кэша Redis
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cachedProduct))
		return
	}

	// Если данные о продукте не найдены в кэше Redis, получаем их из базы данных PostgreSQL
	var product Product
	err = db.Get(&product, "SELECT * FROM products WHERE id = $1", productID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Преобразование продукта в JSON
	productJSON, _ := json.Marshal(product)

	// Сохранение данных о продукте в кэше Redis с TTL 24 часа
	redisClient.Set(ctx, "product:"+productIDStr, string(productJSON), 24*time.Hour)

	// Возвращаем данные о продукте клиенту
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(productJSON)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/products/{id}", getProductByIDHandler).Methods("GET")

	fmt.Println("Server is running...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
