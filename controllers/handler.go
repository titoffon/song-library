package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"main/db" // Импортируем пакет db
	"main/models"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// GetSongs godoc
// @Summary Получить список песен с фильтрацией и пагинацией
// @Description Возвращает список песен с возможностью фильтрации по полям и пагинации. Для управления фильтрацией и пагинацией необходимо в теле запроса указать limit элементов на странице и page номер страницы
// @Tags songs
// @Accept json
// @Produce json
// @Param musicGroup query string false "Название музыкальной группы" example("Muse")
// @Param song query string false "Название песни" example("Supermassive Black Hole")
// @Param releaseDate query string false "Дата релиза песни" example("2006-06-19")
// @Param text query string false "Текст песни" example("Oh baby don”t \n\n you know I suffer?")
// @Param link query string false "Ссылка на источник" example("https://example.com/supermassive-black-hole")
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество элементов на странице" default(1)
// @Success 200 {array} models.Song
// @Failure 404 {object} models.ErrorResponse "Songs or page Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /songs [get]
func GetSongs(w http.ResponseWriter, r *http.Request) {
	var songs []models.Song
	musicGroup := r.URL.Query().Get("musicGroup")
	song := r.URL.Query().Get("song")
	releaseDate := r.URL.Query().Get("release_date")
	text := r.URL.Query().Get("text")
	link := r.URL.Query().Get("link")

	query := db.DB

	//проверка на параметры
	if musicGroup != "" {
		query = query.Where("musicGroup LIKE ?", "%"+musicGroup+"%")
	}
	if song != "" {
		query = query.Where("song LIKE ?", "%"+song+"%")
	}
	if releaseDate != "" {
		query = query.Where("releaseDate LIKE ?", "%"+releaseDate+"%")
	}
	if text != "" {
		query = query.Where("text LIKE ?", "%"+text+"%")
	}
	if link != "" {
		query = query.Where("link LIKE ?", "%"+link+"%")
	}

	// Получаем параметры пагинации
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	//Установка значений по умолчанию
	page := 1
	limit := 10 // Количество куплетов на странице по умолчанию

	if pageParam != "" {
		p, err := strconv.Atoi(pageParam)
		if err == nil && p > 0 {
			page = p
		} else {
			log.Printf("Установка значений по умолчанию: page = %v", page)
		}
	}

	if limitParam != "" {
		l, err := strconv.Atoi(limitParam)
		if err == nil && l > 0 {
			limit = l
		} else {
			log.Printf("Установка значений по умолчанию: limit = %v", limit)
		}
	}

	//реализация вывода с пагинацией
	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&songs).Error; err != nil {
		log.Printf("Ошибка при получении песен: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Новая проверка на отсутствие песен
	if len(songs) == 0 {
		log.Printf("Песни не найдены по заданным параметрам")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Песни не найдены"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

// GetSongVerse godoc
// @Summary Получить текст песни с пагинацией по куплетам
// @Description Возвращает текст песни, разделенный на куплеты, с поддержкой пагинации
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни" example(52)
// @Param page query int false "Номер страницы" default(1) example(1)
// @Param limit query int false "Количество куплетов на странице" default(1) example(2)
// @Success 200 {object} []string
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 404 {object} models.ErrorResponse "Song Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /songs/{id}/verse [get]
func GetSongVerse(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Неверный ID песни: %v", err)
		http.Error(w, "Неверный ID песни", http.StatusBadRequest)
		return
	}

	// Получаем параметры пагинации
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

	//Установка значений по умолчанию
	page := 1
	limit := 1 // Количество куплетов на странице по умолчанию

	if pageParam != "" {
		p, err := strconv.Atoi(pageParam)
		if err == nil && p > 0 {
			page = p
		}
	}

	if limitParam != "" {
		l, err := strconv.Atoi(limitParam)
		if err == nil && l > 0 {
			limit = l
		}
	}

	// Ищем песню в базе данных
	var song models.Song
	if err := db.DB.First(&song, id).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			http.Error(w, "Песня не найдена", http.StatusNotFound)
			log.Printf("Песня с ID %d не найдена", id)
			return
		}
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		log.Printf("Ошибка при получении песни: %v", err)
		return
	}

	// Разбиваем текст песни на куплеты
	verses := strings.Split(song.Text, "\n\n")
	//log.Printf("Срез строк: %v", verses)
	// Пагинация куплетов
	totalVerses := len(verses)
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	if startIndex >= totalVerses {
		http.Error(w, "Страница не найдена", http.StatusNotFound)
		log.Printf("Страница %d не найдена для песни ID %d", page, id)
		return
	}

	if endIndex > totalVerses {
		endIndex = totalVerses
	}

	paginatedVerses := verses[startIndex:endIndex]
	response_slice := make([]string, len(paginatedVerses))
	for i, str := range paginatedVerses {
		response_slice[i] = fmt.Sprintf("Куплет %d:\n %s", i, str)
	}
	// Формируем ответ
	/*response := models.VerseResponse{
		SongID:     song.ID,
		MusicGroup: song.MusicGroup,
		SongName:   song.Song,
		Page:       page,
		Limit:      limit,
		Total:      totalVerses,
		Verses:     paginatedVerses,
	}*/

	log.Printf("Получен текст песни ID %d, страница %d из %d", id, page, (totalVerses+limit-1)/limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response_slice)
}

// DeleteSong godoc
// @Summary Удалить песню по ID
// @Description Удаляет песню из базы данных по ее идентификатору
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни" example(52)
// @Success 200 {object} models.MessageResponse "Песня успешно удалена"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 404 {object} models.ErrorResponse "Song Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /songs/{id} [delete]
func DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Неверный ID песни: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный ID песни"})
		return
	}

	// Ищем и удаляем песню в базе данных
	var song models.Song
	if err := db.DB.First(&song, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Песня с ID %d не найдена", id)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Песня не найдена"})
			return
		}
		log.Printf("Ошибка при поиске песни: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка сервера"})
		return
	}

	if err := db.DB.Delete(&song).Error; err != nil {
		log.Printf("Ошибка при удалении песни: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка сервера"})
		return
	}

	log.Printf("Песня с ID %d успешно удалена", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.MessageResponse{Message: "Песня успешно удалена"})
}

// UpdateSong godoc
// @Summary Изменить данные песни
// @Description Обновляет информацию о песне по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни" example(52)
// @Param song body models.UpdateSongRequest true "Данные для обновления"
// @Success 200 {object} models.Song "Обновленная песня"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 404 {object} models.ErrorResponse "Song Not Found"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /songs/{id} [put]
func UpdateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Printf("Неверный ID песни: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный ID песни"})
		return
	}

	// Парсим данные из тела запроса
	var updateData models.UpdateSongRequest
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		log.Printf("Ошибка при парсинге данных запроса: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	// Ищем песню в базе данных
	var song models.Song
	if err := db.DB.First(&song, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Песня с ID %d не найдена", id)
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Песня не найдена"})
			return
		}
		log.Printf("Ошибка при поиске песни: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка сервера"})
		return
	}

	// Обновляем поля песни, если они переданы в запросе
	if updateData.MusicGroup != nil {
		song.MusicGroup = *updateData.MusicGroup
	}
	if updateData.Song != nil {
		song.Song = *updateData.Song
	}
	if updateData.ReleaseDate != nil {
		song.ReleaseDate = *updateData.ReleaseDate
	}
	if updateData.Text != nil {
		song.Text = *updateData.Text
	}
	if updateData.Link != nil {
		song.Link = *updateData.Link
	}

	// Сохраняем изменения в базе данных
	log.Printf("SONG: %v", song)
	if err := db.DB.Save(&song).Error; err != nil {
		log.Printf("Ошибка при обновлении песни: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка при обновлении песни"})
		return
	}

	log.Printf("Песня с ID %d успешно обновлена", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(song)
}

// AddSong godoc
// @Summary Добавить новую песню
// @Description Добавляет новую песню и обогащает данные из внешнего API
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.AddSongRequest true "Данные песни"
// @Success 200 {object} models.Song
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /songs [post]
func AddSong(w http.ResponseWriter, r *http.Request) {
	//Парсим входящий запрос
	var addSongReq models.AddSongRequest
	if err := json.NewDecoder(r.Body).Decode(&addSongReq); err != nil {
		log.Printf("Ошибка при парсинге запроса: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Неверный формат данных"})
		return
	}

	if addSongReq.Group == "" || addSongReq.Song == "" {
		log.Printf("Отсутствуют обязательные поля в запросе")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Поля 'group' и 'song' обязательны"})
		return
	}

	//Формируем запрос к внешнему API
	params := url.Values{}
	params.Add("group", addSongReq.Group)
	params.Add("song", addSongReq.Song)

	ExternalAPIURL := os.Getenv("EXTERNAL_API_URL")
	//log.Printf("ExternalAPIURL: %s", ExternalAPIURL)

	apiURL, err := url.Parse(ExternalAPIURL)
	if err != nil {
		log.Printf("Некорректный URL внешнего API: %v", err)
		w.WriteHeader(http.StatusInternalServerError) //500
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка сервера"})
		return
	}

	apiURL.RawQuery = params.Encode()

	//Логирование конечного URL запроса
	log.Printf("Запрос к внешнему API: %s", apiURL.String())

	// Создаем HTTP клиент с тайм-аутом для защиты от вечного ожидания
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	//Делаем запрос к внешнему API
	resp, err := client.Get(apiURL.String())
	if err != nil {
		log.Printf("Ошибка при запросе к внешнему API: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка внешнего сервиса"})
		return
	}
	defer resp.Body.Close()

	//Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		log.Printf("Внешний API вернул статус: %d", resp.StatusCode)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка внешнего сервиса"})
		return
	}

	//Читаем и парсим ответ от внешнего API
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа внешнего API: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка сервера"})
		return
	}

	var songDetail models.SongDetail
	if err := json.Unmarshal(body, &songDetail); err != nil {
		log.Printf("Ошибка при парсинге ответа внешнего API: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка сервера"})
		return
	}

	// Проверяем наличие обязательных полей в ответе внешнего API
	if songDetail.ReleaseDate == "" || songDetail.Text == "" || songDetail.Link == "" {
		log.Printf("Ответ внешнего API не содержит необходимых данных")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Некорректный ответ внешнего сервиса"})
		return
	}

	//Создаем новую песню с обогащенными данными
	newSong := models.Song{
		MusicGroup:  addSongReq.Group,
		Song:        addSongReq.Song,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	//Сохраняем песню в базу данных
	if err := db.DB.Create(&newSong).Error; err != nil {
		log.Printf("Ошибка при сохранении песни в базу данных: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Ошибка при сохранении данных"})
		return
	}

	log.Printf("Песня '%s' группы '%s' успешно добавлена", newSong.Song, newSong.MusicGroup)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newSong)
}
