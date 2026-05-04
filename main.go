package main

import (
"encoding/json"
"fmt"
"log"
"os"
"strconv"
"sync"

"github.com/gofiber/fiber/v2"
"github.com/gofiber/fiber/v2/middleware/cors"
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
var mu sync.Mutex

func loadDB() error {
file, err := os.Open("data.json")
if err != nil {
data = DataBase{
Fam:      []string{},
Ima:      []string{},
Otch:     []string{},
Phone:    []string{},
Kolvo:    []int{},
Num:      []int{},
Capacity: []int{},
Status:   []bool{},
Date:     []string{},
Book:     []string{},
User:     []string{},
Admin:    []string{},
}
fmt.Println("Файл data.json не найден. Инициализация новой БД в памяти.")
return nil
}
defer file.Close()

decoder := json.NewDecoder(file)
if err := decoder.Decode(&data); err != nil {
fmt.Println("Ошибка чтения data.json. Инициализация новой БД.")
data = DataBase{
Fam:      []string{},
Ima:      []string{},
Otch:     []string{},
Phone:    []string{},
Kolvo:    []int{},
Num:      []int{},
Capacity: []int{},
Status:   []bool{},
Date:     []string{},
Book:     []string{},
User:     []string{},
Admin:    []string{},
}
return nil
}

length := len(data.Num)
if len(data.Capacity) != length || len(data.Status) != length || len(data.Date) != length || len(data.Book) != length {
fmt.Println("Обнаружена рассинхронизация данных комнат. Исправление структуры.")
minLen := length
if len(data.Capacity) < minLen { minLen = len(data.Capacity) }
if len(data.Status) < minLen { minLen = len(data.Status) }
if len(data.Date) < minLen { minLen = len(data.Date) }
if len(data.Book) < minLen { minLen = len(data.Book) }

data.Num = data.Num[:minLen]
data.Capacity = data.Capacity[:minLen]
data.Status = data.Status[:minLen]
data.Date = data.Date[:minLen]
data.Book = data.Book[:minLen]
}

fmt.Println("База данных успешно загружена.")
return nil
}

func saveDB() error {
outFile, err := os.Create("data.json")
if err != nil {
return fmt.Errorf("ошибка создания файла: %w", err)
}
defer outFile.Close()

encoder := json.NewEncoder(outFile)
encoder.SetIndent("", "  ")
if err := encoder.Encode(data); err != nil {
return fmt.Errorf("ошибка записи в файл: %w", err)
}
return nil
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

func findRoomIndexByNum(roomNum int) int {
for i, num := range data.Num {
if num == roomNum {
return i
}
}
return -1
}

func main() {
if err := loadDB(); err != nil {
log.Fatal("Критическая ошибка загрузки БД:", err)
}

app := fiber.New()

app.Use(cors.New())
app.Use(logger.New())

app.Get("/login", func(c *fiber.Ctx) error {
login := c.Query("login")
if login == "" {
return c.Status(400).JSON(fiber.Map{"error": "Логин обязателен"})
}

idx := findUser(login)
isAdmin := findAdmin(login)

if idx == -1 && !isAdmin {
return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
}

fio := login
if idx != -1 {
f := data.Fam[idx]
i := data.Ima[idx]
o := data.Otch[idx]
if f != "" || i != "" || o != "" {
fio = fmt.Sprintf("%s %s %s", f, i, o)
}
}

return c.JSON(fiber.Map{
"message":  "Вход выполнен успешно",
"login":    login,
"fio":      fio,
"is_admin": isAdmin,
})
})

app.Get("/register", func(c *fiber.Ctx) error {
fam := c.Query("fam")
ima := c.Query("ima")
otch := c.Query("otch")
phone := c.Query("phone")
login := c.Query("login")
kolvoStr := c.Query("kolvo")

if login == "" {
return c.Status(400).JSON(fiber.Map{"error": "Логин обязателен"})
}

kolvo, err := strconv.Atoi(kolvoStr)
if err != nil {
kolvo = 1
}

if isUserExists(login) {
return c.Status(409).JSON(fiber.Map{"error": "Пользователь с таким логином уже существует"})
}

mu.Lock()
data.Fam = append(data.Fam, fam)
data.Ima = append(data.Ima, ima)
data.Otch = append(data.Otch, otch)
data.Phone = append(data.Phone, phone)
data.Kolvo = append(data.Kolvo, kolvo)
data.User = append(data.User, login)
mu.Unlock()

if err := saveDB(); err != nil {
return c.Status(500).JSON(fiber.Map{"error": "Ошибка сохранения"})
}

return c.JSON(fiber.Map{"message": "Пользователь зарегистрирован", "fio": fmt.Sprintf("%s %s %s", fam, ima, otch)})
})

app.Get("/admin/create", func(c *fiber.Ctx) error {
login := c.Query("login")
if login == "" {
return c.Status(400).JSON(fiber.Map{"error": "Логин обязателен"})
}

if isUserExists(login) {
return c.Status(409).JSON(fiber.Map{"error": "Пользователь уже существует"})
}

mu.Lock()
data.Admin = append(data.Admin, login)
data.User = append(data.User, login)
data.Fam = append(data.Fam, "")
data.Ima = append(data.Ima, "")
data.Otch = append(data.Otch, "")
data.Phone = append(data.Phone, "")
data.Kolvo = append(data.Kolvo, 1)
mu.Unlock()

if err := saveDB(); err != nil {
return c.Status(500).JSON(fiber.Map{"error": "Ошибка сохранения"})
}

return c.JSON(fiber.Map{"message": "Администратор создан", "login": login})
})

app.Get("/room/create", func(c *fiber.Ctx) error {
login := c.Query("login")
if !findAdmin(login) {
return c.Status(403).JSON(fiber.Map{"error": "Требуется права администратора"})
}

numStr := c.Query("num")
capStr := c.Query("cap")

num, err := strconv.Atoi(numStr)
if err != nil {
return c.Status(400).JSON(fiber.Map{"error": "num должен быть числом"})
}

cap, err := strconv.Atoi(capStr)
if err != nil {
return c.Status(400).JSON(fiber.Map{"error": "cap должен быть числом"})
}

mu.Lock()
if findRoomIndexByNum(num) != -1 {
mu.Unlock()
return c.Status(409).JSON(fiber.Map{"error": "Комната с таким номером уже существует"})
}

data.Num = append(data.Num, num)
data.Capacity = append(data.Capacity, cap)
data.Status = append(data.Status, false)
data.Date = append(data.Date, "")
data.Book = append(data.Book, "")
mu.Unlock()

if err := saveDB(); err != nil {
return c.Status(500).JSON(fiber.Map{"error": "Ошибка сохранения"})
}

log.Printf("Комната %d создана. Статус: свободна", num)
return c.JSON(fiber.Map{"message": "Комната создана", "num": num, "status": "free"})
})

app.Get("/room/delete", func(c *fiber.Ctx) error {
login := c.Query("login")
if !findAdmin(login) {
return c.Status(403).JSON(fiber.Map{"error": "Требуется права администратора"})
}

roomNum, err := strconv.Atoi(c.Query("room_num"))
if err != nil {
return c.Status(400).JSON(fiber.Map{"error": "Неверный номер комнаты"})
}

mu.Lock()
idx := findRoomIndexByNum(roomNum)
if idx == -1 {
mu.Unlock()
return c.Status(404).JSON(fiber.Map{"error": "Комната не найдена"})
}

data.Num = append(data.Num[:idx], data.Num[idx+1:]...)
data.Capacity = append(data.Capacity[:idx], data.Capacity[idx+1:]...)
data.Status = append(data.Status[:idx], data.Status[idx+1:]...)
data.Date = append(data.Date[:idx], data.Date[idx+1:]...)
data.Book = append(data.Book[:idx], data.Book[idx+1:]...)
mu.Unlock()

if err := saveDB(); err != nil {
return c.Status(500).JSON(fiber.Map{"error": "Ошибка сохранения"})
}

return c.JSON(fiber.Map{"message": "Комната удалена"})
})

app.Get("/room/book", func(c *fiber.Ctx) error {
login := c.Query("login")
roomNumStr := c.Query("room_num")
date := c.Query("date")

roomNum, err := strconv.Atoi(roomNumStr)
if err != nil {
return c.Status(400).JSON(fiber.Map{"error": "room_num должен быть числом"})
}

mu.Lock()
defer mu.Unlock()

userIdx := findUser(login)
if userIdx == -1 {
return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
}

roomIdx := findRoomIndexByNum(roomNum)
if roomIdx == -1 {
return c.Status(404).JSON(fiber.Map{"error": "Комната не найдена"})
}

currentStatus := data.Status[roomIdx]
log.Printf("Попытка брони: Комната %d, Индекс %d, Текущий статус (true=занята): %v", roomNum, roomIdx, currentStatus)

if currentStatus {
return c.Status(409).JSON(fiber.Map{
"error":     "Комната уже забронирована",
"booked_by": data.Book[roomIdx],
"date":      data.Date[roomIdx],
})
}

if data.Capacity[roomIdx] < data.Kolvo[userIdx] {
return c.Status(400).JSON(fiber.Map{"error": "Недостаточно мест в комнате"})
}

data.Status[roomIdx] = true
data.Date[roomIdx] = date
data.Book[roomIdx] = fmt.Sprintf("%s %s %s", data.Fam[userIdx], data.Ima[userIdx], data.Otch[userIdx])

if err := saveDB(); err != nil {
return c.Status(500).JSON(fiber.Map{"error": "Ошибка сохранения"})
}

log.Printf("Комната %d успешно забронирована пользователем %s", roomNum, login)
return c.JSON(fiber.Map{
"message":   "Комната забронирована",
"room_num":  roomNum,
"date":      date,
"booked_by": data.Book[roomIdx],
})
})

app.Get("/room/unbook", func(c *fiber.Ctx) error {
login := c.Query("login")
roomNum, err := strconv.Atoi(c.Query("room_num"))
if err != nil {
return c.Status(400).JSON(fiber.Map{"error": "Неверный номер комнаты"})
}

mu.Lock()
defer mu.Unlock()

userIdx := findUser(login)
if userIdx == -1 {
return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
}

roomIdx := findRoomIndexByNum(roomNum)
if roomIdx == -1 {
return c.Status(404).JSON(fiber.Map{"error": "Комната не найдена"})
}

userName := fmt.Sprintf("%s %s %s", data.Fam[userIdx], data.Ima[userIdx], data.Otch[userIdx])
isAdmin := findAdmin(login)

if !isAdmin && data.Book[roomIdx] != userName {
return c.Status(403).JSON(fiber.Map{"error": "Нет прав для отмены этой брони"})
}

data.Status[roomIdx] = false
data.Date[roomIdx] = ""
data.Book[roomIdx] = ""

if err := saveDB(); err != nil {
return c.Status(500).JSON(fiber.Map{"error": "Ошибка сохранения"})
}

return c.JSON(fiber.Map{"message": "Бронь отменена"})
})

app.Get("/rooms", func(c *fiber.Ctx) error {
login := c.Query("login")
if findUser(login) == -1 && !findAdmin(login) {
return c.Status(401).JSON(fiber.Map{"error": "Необходимо войти"})
}

mu.Lock()
defer mu.Unlock()

rooms := make([]fiber.Map, 0, len(data.Num))
for i := range data.Num {
statusText := "свободна"
if data.Status[i] {
statusText = "забронирована"
}
rooms = append(rooms, fiber.Map{
"num":       data.Num[i],
"capacity":  data.Capacity[i],
"status":    statusText,
"is_busy":   data.Status[i],
"date":      data.Date[i],
"booked_by": data.Book[i],
})
}

return c.JSON(fiber.Map{"rooms": rooms})
})

app.Get("/rooms/available", func(c *fiber.Ctx) error {
login := c.Query("login")
userIdx := findUser(login)
if userIdx == -1 {
return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
}

mu.Lock()
defer mu.Unlock()

userKolvo := data.Kolvo[userIdx]
available := make([]fiber.Map, 0)

for i := range data.Num {
if !data.Status[i] && data.Capacity[i] >= userKolvo {
available = append(available, fiber.Map{
"num":      data.Num[i],
"capacity": data.Capacity[i],
})
}
}

return c.JSON(fiber.Map{"available_rooms": available})
})

app.Get("/user/profile", func(c *fiber.Ctx) error {
login := c.Query("login")
idx := findUser(login)
if idx == -1 {
return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
}

fio := fmt.Sprintf("%s %s %s", data.Fam[idx], data.Ima[idx], data.Otch[idx])
if fio == "  " {
fio = login
}

return c.JSON(fiber.Map{
"fio":      fio,
"phone":    data.Phone[idx],
"kolvo":    data.Kolvo[idx],
"is_admin": findAdmin(login),
})
})

app.Get("/user/bookings", func(c *fiber.Ctx) error {
login := c.Query("login")
idx := findUser(login)
if idx == -1 {
return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
}

userName := fmt.Sprintf("%s %s %s", data.Fam[idx], data.Ima[idx], data.Otch[idx])
bookings := make([]fiber.Map, 0)

for i := range data.Book {
if data.Book[i] == userName {
bookings = append(bookings, fiber.Map{
"room_num": data.Num[i],
"date":     data.Date[i],
"capacity": data.Capacity[i],
})
}
}

return c.JSON(fiber.Map{"bookings": bookings})
})

app.Get("/user/update_kolvo", func(c *fiber.Ctx) error {
login := c.Query("login")
kolvo, err := strconv.Atoi(c.Query("kolvo"))
if err != nil || kolvo <= 0 {
return c.Status(400).JSON(fiber.Map{"error": "Некорректное количество"})
}

idx := findUser(login)
if idx == -1 {
return c.Status(404).JSON(fiber.Map{"error": "Пользователь не найден"})
}

mu.Lock()
data.Kolvo[idx] = kolvo
mu.Unlock()

if err := saveDB(); err != nil {
return c.Status(500).JSON(fiber.Map{"error": "Ошибка сохранения"})
}

return c.JSON(fiber.Map{"message": "Обновлено", "kolvo": kolvo})
})

app.Get("/debug/status", func(c *fiber.Ctx) error {
login := c.Query("login")
if !findAdmin(login) {
return c.Status(403).JSON(fiber.Map{"error": "Только для админов"})
}

mu.Lock()
defer mu.Unlock()

return c.JSON(fiber.Map{
"users_count": len(data.User),
"rooms_count": len(data.Num),
"raw_rooms": data.Num,
"raw_status": data.Status,
"raw_book": data.Book,
"raw_capacity": data.Capacity,
})
})

log.Fatal(app.Listen(":3000"))
}
