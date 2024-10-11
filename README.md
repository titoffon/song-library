# 🎶 Онлайн Библиотека Песен

Это приложение представляет собой RESTful API для управления библиотекой песен. Приложение позволяет получать информацию о песнях, управлять ими и обогащать данные через внешний API.

## 🔧 Конфигурация

Все настройки приложения управляются через переменные окружения, определенные в файле `.env`. Это включает параметры подключения к базе данных, URL внешнего API и порт сервера.

## 🗄 Миграции базы данных и тестовые данные

При запуске приложения автоматически выполняются миграции для создания необходимых таблиц. Если вы хотите воспользоваться тестовыми данными для наполнения базы данных, вам необходимо раскомментировать соответствующий блок в файле `database.go` внутри пакета `db`.

**Инструкция по добавлению тестовых данных:**

1. Откройте файл `db/database.go`.
2. Найдите блок кода, отвечающий за заполнение базы данных тестовыми данными:

    ```go
    // Выполняем SQL-файл для заполнения таблицы
    if err := migrations.ExecuteSQLFile(DB, "migrations/insert_to_songs_table.sql"); err != nil {
        log.Fatalf("Ошибка при выполнении SQL-файла: %v", err)
    }
    ```

3. Раскомментируйте этот блок кода, удалив `/*` и `*/` или символы комментариев `//`.
   
## 📋 Функциональные возможности

- **Получение списка песен** с фильтрацией по всем полям и поддержкой пагинации.
- **Получение текста песни** с пагинацией по куплетам.
- **Добавление новой песни**, обогащение данных через внешний API.
- **Обновление** и **удаление** песен.
- **Миграции базы данных** при запуске сервиса.
- **Конфигурация через файл `.env`**.
- **Автоматическая генерация Swagger-документации** для API.

## 🚀 **Настройте переменные окружения**

    Создайте файл `.env` в корневой директории со следующим содержимым:

    ```env
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=your_db_name
    EXTERNAL_API_URL=https://external.api/info
    SERVER_PORT=8000
    ```
## 📖 Документация API

Swagger-документация доступна по адресу: http://localhost:8000/swagger/index.html

## 🛠 Структура проекта

- `main.go` — точка входа в приложение.
- `controllers/` — обработчики HTTP-запросов.
- `models/` — определения моделей данных.
- `db/` — инициализация и подключение к базе данных.
- `migrations/` — файлы миграций SQL для создания и заполнения базы данных.
- `router.go` — настройка маршрутов HTTP.

## 📚 Основные возможности

### Получение списка песен

- **Endpoint:** `GET /songs`
- **Описание:** Возвращает список песен с возможностью фильтрации по всем полям и пагинации.
- **Параметры запроса:**
  - `music_group` (string): Фильтрация по названию группы.
  - `song` (string): Фильтрация по названию песни.
  - `page` (int): Номер страницы (по умолчанию 1).
  - `limit` (int): Количество элементов на странице (по умолчанию 1).

### Получение текста песни с пагинацией по куплетам

- **Endpoint:** `GET /songs/{id}/verse`
- **Описание:** Возвращает текст песни, разделенный на куплеты, с поддержкой пагинации.
- **Параметры запроса:**
  - `page` (int): Номер страницы (по умолчанию 1).
  - `limit` (int): Количество куплетов на странице (по умолчанию 1).

### Добавление новой песни

- **Endpoint:** `POST /songs`
- **Описание:** Добавляет новую песню и обогащает данные из внешнего API.
- **Тело запроса:**

    ```json
    {
      "group": "Muse",
      "song": "Supermassive Black Hole"
    }
    ```

### Обновление песни

- **Endpoint:** `PUT /songs/{id}`
- **Описание:** Обновляет информацию о песне по её ID.
- **Тело запроса:** Поля, которые необходимо обновить.

### Удаление песни

- **Endpoint:** `DELETE /songs/{id}`
- **Описание:** Удаляет песню из базы данных по её ID.
