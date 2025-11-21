package main

import (
	"database/sql"
	"hw6/handlers"
	"hw6/utils"
	"net/http"
)

/*
Менеджер MySQL-базы данных

Нельзя глобальные переменные, использовать поля структуры в замыкании

Запросы:
- GET / - возвращает список всех таблиц

- GET /$table?limit=5&offset=7 - возвращает список записей. limit по-умолчанию 5, offset 0
- PUT /$table - создаёт новую запись, данный по записи в теле запроса (POST-параметры)

- GET /$table/$id - возвращает информацию о самой записи или 404
- POST /$table/$id - обновляет запись, данные приходят в теле запроса (POST-параметры)
- DELETE /$table/$id - удаляет запись

Особенности:
- Роутинг запросов - руками
- В NewDBExplorer считываем из базы список таблиц, полей, работаем с ними при валидации
- спискок параметров - вы его подгружаете динамически при инициализации. Не задан строго
- Валидация "string - int - float - null".
- json в пустой интерфейс распаковывает как float, если не указаны спец. опции
- Вся работа через database/sql. Никаких orm.
- Все имена полей так как они в базе
- В случае ошибки -  возвращаем 500
- Не забывайте про SQL-инъекции. Неизвестные поля игнорируем

SHOW TABLES; SHOW FULL COLUMNS FROM `$table_name`;

Подсказки:
- Внутри row лежат значения полей и метаданные
- Активно применяться пустые интерфейсы
- Обработка null-значений
- Придётся вытаскивать неизвестное количество полей из row, подумайте как тут можно применить пустые интерфейсы
- Поднять mysql-базу локально через докер:

*/

func NewDBExplorer(db *sql.DB) (http.Handler, error) {
	// limit open connections to one to avoid unexpected extra connections
	// db.SetMaxOpenConns(1)
	// db.SetMaxIdleConns(1)
	tables, err := utils.LoadTables(db)
	if err != nil {
		return nil, err
	}
	handlers := &handlers.Handler{DB: db, Tables: tables}
	r := http.NewServeMux()
	r.HandleFunc("GET    /", handlers.ListTablesHandler)
	r.HandleFunc("GET    /health", handlers.HealthHandler)

	r.HandleFunc("GET    /{table}/{id}", handlers.GetRecordHandler)
	r.HandleFunc("POST   /{table}/{id}", handlers.UpdateRecordHandler)
	r.HandleFunc("DELETE /{table}/{id}", handlers.DeleteRecordHandler)
	r.HandleFunc("GET    /{table}", handlers.ListRecordsHandler)  // query params: limit, offset
	r.HandleFunc("GET    /{table}/", handlers.ListRecordsHandler) // support trailing slash
	r.HandleFunc("PUT    /{table}", handlers.CreateRecordHandler)
	r.HandleFunc("PUT    /{table}/", handlers.CreateRecordHandler)
	return r, nil
}
