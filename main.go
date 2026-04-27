package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors" // 1. Импорт CORS middleware
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type DataBase struct {
	Fam      []string `json:"fam"`
	Ima      []string `json:"ima"`
	Otch     []string `json:"otch"`
	Phone    []string `json:"phone"`
	Kolvo    []int    `json:"kolvo"`
	Num      []int    `json:"num"`
	Capacity []int    `json:"cap"`
	Status   []bool   `json:"status"`
	Date     []string `json:"date"`
	Book     []string `json:"book"`
	User     []string `json:"user"`
	Admin    []string `json:"admin"`
}

var data DataBase

func loadDB() error {
	file, err := os.Open("data.json")
	if err != nil {
		// Создаем файл с пустой структурой
		data = DataBase{}
		saveDB()
		return nil
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(&data)
}

func saveDB() error {
	outFile, err := os.Create("data.json")
	if err != nil {
		return err
	}
	defer outFile.Close()
	encoder := json.NewEncoder(outFile)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func findUser(login string) int {
	for i, userLogin := range data.User {
		if userLogin == login {
			return i
		}
	}
	return -1
}

func findAdmin(login string) bool {
	for _, adminLogin := range data.Admin {
		if adminLogin == login {
			return true
		}
	}
	return false
}

func isUserExists(login string) bool {
	return findUser(login) != -1 || findAdmin(login)
}

func main() {
	// Загружаем базу данных
	if err := loadDB(); err != nil {
		log.Fatal("Ошибка загрузки базы данных:", err)
	}

	app := fiber.New()

	// 2. Подключение CORS middleware
	// Разрешает запросы с любого источника (для разработки).
	// В продакшене лучше указать конкретный домен: AllowOrigins: "http://mysite.com"
	app.Use(cors.New())

	// Логирование запросов
	app.Use(logger.New())

	// Регистрация пользователя
	app.Get("/register", func(c *fiber.Ctx) error {
		fam := c.Query("fam")
		ima := c.Query("ima")
		otch := c.Query("otch")
		phone := c.Query("phone")
		login := c.Query("login")

		kolvoStr := c.Query("kolvo")
		kolvo, err := strconv.Atoi(kolvoStr)
		if err != nil {
			return c.JSON(fiber.Map{"error": "kolvo должен быть числом"})
		}

		if isUserExists(login) {
			return c.JSON(fiber.Map{"error": "Пользователь с таким логином уже существует"})
		}

		data.Fam = append(data.Fam, fam)
		data.Ima = append(data.Ima, ima)
		data.Otch = append(data.Otch, otch)
		data.Phone = append(data.Phone, phone)
		data.Kolvo = append(data.Kolvo, kolvo)
		data.User = append(data.User, login)

		if err := saveDB(); err != nil {
			return c.JSON(fiber.Map{"error": "Ошибка сохранения"})
		}

		return c.JSON(fiber.Map{"message": "Пользователь зарегистрирован", "fio": fmt.Sprintf("%s %s %s", fam, ima, otch)})
	})

	// Создание администратора
	app.Get("/admin/create", func(c *fiber.Ctx) error {
		login := c.Query("login")

		if isUserExists(login) {
			return c.JSON(fiber.Map{"error": "Пользователь с таким логином уже существует"})
		}

		// Добавляем в список админов
		data.Admin = append(data.Admin, login)
		// Также добавляем в пользователи для возможности входа
		data.User = append(data.User, login)
		// Добавляем пустые данные пользователя
		data.Fam = append(data.Fam, "")
		data.Ima = append(data.Ima, "")
		data.Otch = append(data.Otch, "")
		data.Phone = append(data.Phone, "")
		data.Kolvo = append(data.Kolvo, 0)

		if err := saveDB(); err != nil {
			return c.JSON(fiber.Map{"error": "Ошибка сохранения"})
		}

		return c.JSON(fiber.Map{"message": "Администратор создан", "login": login})
	})

	// Создание комнаты (только админ)
	app.Get("/room/create", func(c *fiber.Ctx) error {
		login := c.Query("login")

		if !findAdmin(login) {
			return c.JSON(fiber.Map{"error": "Требуется права администратора"})
		}

		numStr := c.Query("num")
		capStr := c.Query("cap")

		num, err := strconv.Atoi(numStr)
		if err != nil {
			return c.JSON(fiber.Map{"error": "num должен быть числом"})
		}

		cap, err := strconv.Atoi(capStr)
		if err != nil {
			return c.JSON(fiber.Map{"error": "cap должен быть числом"})
		}

		data.Num = append(data.Num, num)
		data.Capacity = append(data.Capacity, cap)
		data.Status = append(data.Status, false)
		data.Date = append(data.Date, "неизвестно")
		data.Book = append(data.Book, "неизвестно")

		if err := saveDB(); err != nil {
			return c.JSON(fiber.Map{"error": "Ошибка сохранения"})
		}

		return c.JSON(fiber.Map{"message": "Создана комната", "num": num})
	})

	// Удаление пользователя (только админ)
	app.Get("/user/delete", func(c *fiber.Ctx) error {
		login := c.Query("login")

		if !findAdmin(login) {
			return c.JSON(fiber.Map{"error": "Требуется права администратора"})
		}

		userLogin := c.Query("user_login")
		userIndex := findUser(userLogin)

		if userIndex == -1 {
			return c.JSON(fiber.Map{"error": "Пользователь не найден"})
		}

		// Удаляем пользователя из всех срезов
		data.Fam = append(data.Fam[:userIndex], data.Fam[userIndex+1:]...)
		data.Ima = append(data.Ima[:userIndex], data.Ima[userIndex+1:]...)
		data.Otch = append(data.Otch[:userIndex], data.Otch[userIndex+1:]...)
		data.Phone = append(data.Phone[:userIndex], data.Phone[userIndex+1:]...)
		data.Kolvo = append(data.Kolvo[:userIndex], data.Kolvo[userIndex+1:]...)
		data.User = append(data.User[:userIndex], data.User[userIndex+1:]...)

		// Удаляем из админов, если был
		for i, adminLogin := range data.Admin {
			if adminLogin == userLogin {
				data.Admin = append(data.Admin[:i], data.Admin[i+1:]...)
				break
			}
		}

		if err := saveDB(); err != nil {
			return c.JSON(fiber.Map{"error": "Ошибка сохранения"})
		}

		return c.JSON(fiber.Map{"message": "Пользователь удален"})
	})

	// Бронирование комнаты (админ или пользователь)
	app.Get("/room/book", func(c *fiber.Ctx) error {
		login := c.Query("login")
		roomNumStr := c.Query("room_num")
		date := c.Query("date")

		roomNum, err := strconv.Atoi(roomNumStr)
		if err != nil {
			return c.JSON(fiber.Map{"error": "room_num должен быть числом"})
		}

		userIndex := findUser(login)
		if userIndex == -1 {
			return c.JSON(fiber.Map{"error": "Пользователь не найден"})
		}

		// Находим индекс комнаты по номеру
		roomIndex := -1
		for i, num := range data.Num {
			if num == roomNum {
				roomIndex = i
				break
			}
		}

		if roomIndex == -1 {
			return c.JSON(fiber.Map{"error": "Комната не найдена"})
		}

		if data.Status[roomIndex] {
			return c.JSON(fiber.Map{"error": "Комната уже забронирована"})
		}

		if data.Capacity[roomIndex] < data.Kolvo[userIndex] {
			return c.JSON(fiber.Map{"error": "Комната не вместит указанное количество персон"})
		}

		data.Status[roomIndex] = true
		data.Date[roomIndex] = date
		data.Book[roomIndex] = fmt.Sprintf("%s %s %s", data.Fam[userIndex], data.Ima[userIndex], data.Otch[userIndex])

		if err := saveDB(); err != nil {
			return c.JSON(fiber.Map{"error": "Ошибка сохранения"})
		}

		return c.JSON(fiber.Map{
			"message":   "Комната успешно забронирована",
			"room_num":  roomNum,
			"date":      date,
			"booked_by": data.Book[roomIndex],
		})
	})

	// Отмена бронирования (админ или пользователь, если это его бронь)
	app.Get("/room/unbook", func(c *fiber.Ctx) error {
		login := c.Query("login")
		roomNumStr := c.Query("room_num")

		roomNum, err := strconv.Atoi(roomNumStr)
		if err != nil {
			return c.JSON(fiber.Map{"error": "room_num должен быть числом"})
		}

		userIndex := findUser(login)
		if userIndex == -1 {
			return c.JSON(fiber.Map{"error": "Пользователь не найден"})
		}

		userName := fmt.Sprintf("%s %s %s", data.Fam[userIndex], data.Ima[userIndex], data.Otch[userIndex])

		// Находим индекс комнаты по номеру
		roomIndex := -1
		for i, num := range data.Num {
			if num == roomNum {
				roomIndex = i
				break
			}
		}

		if roomIndex == -1 {
			return c.JSON(fiber.Map{"error": "Комната не найдена"})
		}

		// Проверяем права: админ или владелец брони
		if !findAdmin(login) && data.Book[roomIndex] != userName {
			return c.JSON(fiber.Map{"error": "Нет прав для отмены этой брони"})
		}

		data.Status[roomIndex] = false
		data.Date[roomIndex] = "неизвестно"
		data.Book[roomIndex] = "неизвестно"

		if err := saveDB(); err != nil {
			return c.JSON(fiber.Map{"error": "Ошибка сохранения"})
		}

		return c.JSON(fiber.Map{"message": "Бронирование отменено"})
	})

	// Удаление комнаты (только админ)
	app.Get("/room/delete", func(c *fiber.Ctx) error {
		login := c.Query("login")
		roomNumStr := c.Query("room_num")

		if !findAdmin(login) {
			return c.JSON(fiber.Map{"error": "Требуется права администратора"})
		}

		roomNum, err := strconv.Atoi(roomNumStr)
		if err != nil {
			return c.JSON(fiber.Map{"error": "room_num должен быть числом"})
		}

		roomIndex := -1
		for i, num := range data.Num {
			if num == roomNum {
				roomIndex = i
				break
			}
		}

		if roomIndex == -1 {
			return c.JSON(fiber.Map{"error": "Комната не найдена"})
		}

		// Удаляем комнату из всех срезов
		data.Num = append(data.Num[:roomIndex], data.Num[roomIndex+1:]...)
		data.Capacity = append(data.Capacity[:roomIndex], data.Capacity[roomIndex+1:]...)
		data.Status = append(data.Status[:roomIndex], data.Status[roomIndex+1:]...)
		data.Date = append(data.Date[:roomIndex], data.Date[roomIndex+1:]...)
		data.Book = append(data.Book[:roomIndex], data.Book[roomIndex+1:]...)

		if err := saveDB(); err != nil {
			return c.JSON(fiber.Map{"error": "Ошибка сохранения"})
		}

		return c.JSON(fiber.Map{"message": "Комната удалена"})
	})

	// Изменение количества персон (только пользователь для себя)
	app.Get("/user/update_kolvo", func(c *fiber.Ctx) error {
		login := c.Query("login")
		kolvoStr := c.Query("kolvo")

		kolvo, err := strconv.Atoi(kolvoStr)
		if err != nil {
			return c.JSON(fiber.Map{"error": "kolvo должен быть числом"})
		}

		userIndex := findUser(login)
		if userIndex == -1 {
			return c.JSON(fiber.Map{"error": "Пользователь не найден"})
		}

		if kolvo <= 0 {
			return c.JSON(fiber.Map{"error": "Количество персон должно быть положительным"})
		}

		data.Kolvo[userIndex] = kolvo

		if err := saveDB(); err != nil {
			return c.JSON(fiber.Map{"error": "Ошибка сохранения"})
		}

		return c.JSON(fiber.Map{"message": "Количество персон обновлено", "new_kolvo": kolvo})
	})

	// Изменение даты бронирования (только пользователь для своей брони)
	app.Get("/room/update_date", func(c *fiber.Ctx) error {
		login := c.Query("login")
		roomNumStr := c.Query("room_num")
		newDate := c.Query("new_date")

		roomNum, err := strconv.Atoi(roomNumStr)
		if err != nil {
			return c.JSON(fiber.Map{"error": "room_num должен быть числом"})
		}

		userIndex := findUser(login)
		if userIndex == -1 {
			return c.JSON(fiber.Map{"error": "Пользователь не найден"})
		}

		userName := fmt.Sprintf("%s %s %s", data.Fam[userIndex], data.Ima[userIndex], data.Otch[userIndex])

		// Находим индекс комнаты по номеру
		roomIndex := -1
		for i, num := range data.Num {
			if num == roomNum {
				roomIndex = i
				break
			}
		}

		if roomIndex == -1 {
			return c.JSON(fiber.Map{"error": "Комната не найдена"})
		}

		// Проверяем, что это его бронь или он админ
		if !findAdmin(login) && data.Book[roomIndex] != userName {
			return c.JSON(fiber.Map{"error": "Нет прав для изменения этой брони"})
		}

		data.Date[roomIndex] = newDate

		if err := saveDB(); err != nil {
			return c.JSON(fiber.Map{"error": "Ошибка сохранения"})
		}

		return c.JSON(fiber.Map{"message": "Дата бронирования обновлена", "new_date": newDate})
	})

	// Просмотр профиля пользователя
	app.Get("/user/profile", func(c *fiber.Ctx) error {
		login := c.Query("login")

		userIndex := findUser(login)
		if userIndex == -1 {
			return c.JSON(fiber.Map{"error": "Пользователь не найден"})
		}

		return c.JSON(fiber.Map{
			"fio":      fmt.Sprintf("%s %s %s", data.Fam[userIndex], data.Ima[userIndex], data.Otch[userIndex]),
			"phone":    data.Phone[userIndex],
			"kolvo":    data.Kolvo[userIndex],
			"is_admin": findAdmin(login),
		})
	})

	// Просмотр всех комнат
	app.Get("/rooms", func(c *fiber.Ctx) error {
		login := c.Query("login")

		if !findAdmin(login) && findUser(login) == -1 {
			return c.JSON(fiber.Map{"error": "Необходимо войти"})
		}

		rooms := make([]fiber.Map, len(data.Num))
		for i := range data.Num {
			status := "свободна"
			if data.Status[i] {
				status = "забронирована"
			}

			rooms[i] = fiber.Map{
				"num":       data.Num[i],
				"capacity":  data.Capacity[i],
				"status":    status,
				"date":      data.Date[i],
				"booked_by": data.Book[i],
			}
		}

		return c.JSON(fiber.Map{"rooms": rooms})
	})

	// Просмотр бронирований пользователя
	app.Get("/user/bookings", func(c *fiber.Ctx) error {
		login := c.Query("login")

		userIndex := findUser(login)
		if userIndex == -1 {
			return c.JSON(fiber.Map{"error": "Пользователь не найден"})
		}

		userName := fmt.Sprintf("%s %s %s", data.Fam[userIndex], data.Ima[userIndex], data.Otch[userIndex])

		bookings := []fiber.Map{}
		for i := range data.Book {
			if data.Book[i] == userName {
				bookings = append(bookings, fiber.Map{
					"room_num": data.Num[i],
					"capacity": data.Capacity[i],
					"date":     data.Date[i],
				})
			}
		}

		return c.JSON(fiber.Map{"bookings": bookings})
	})

	// Проверка доступных комнат для бронирования
	app.Get("/rooms/available", func(c *fiber.Ctx) error {
		login := c.Query("login")

		userIndex := findUser(login)
		if userIndex == -1 {
			return c.JSON(fiber.Map{"error": "Пользователь не найден"})
		}

		availableRooms := []fiber.Map{}
		for i := range data.Num {
			if !data.Status[i] && data.Capacity[i] >= data.Kolvo[userIndex] {
				availableRooms = append(availableRooms, fiber.Map{
					"num":      data.Num[i],
					"capacity": data.Capacity[i],
				})
			}
		}

		return c.JSON(fiber.Map{"available_rooms": availableRooms})
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		login := c.Query("login")

		if login == "" {
			return c.JSON(fiber.Map{"error": "Логин обязателен"})
		}

		if !isUserExists(login) {
			return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
		}

		idx := findUser(login)
		fio := ""
		if idx != -1 {
			fio = fmt.Sprintf("%s %s %s", data.Fam[idx], data.Ima[idx], data.Otch[idx])
		}

		return c.JSON(fiber.Map{
			"message":  "Вход выполнен успешно",
			"login":    login,
			"fio":      fio,
			"is_admin": findAdmin(login),
		})
	})

	log.Fatal(app.Listen(":3000"))
}
