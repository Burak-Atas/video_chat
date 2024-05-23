package main

/*
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

// Room yapısı, her özel sohbet odası için bağlantıları ve kilitlemeyi yönetir

	type Room struct {
		connections map[*websocket.Conn]bool
		sync.RWMutex
	}

	type User struct {
		Name string
		Conn map[*websocket.Conn]bool
	}

var rooms = make(map[string]*Room)
var roomsMutex sync.Mutex

	func main() {
		router := gin.Default()

		router.POST("/login", func(c *gin.Context) {
			fmt.Println("Login request received")
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		router.POST("/validate", func(c *gin.Context) {
			fmt.Println("Validation request received")
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		router.GET("/ws/:roomID", func(c *gin.Context) {
			roomID := c.Param("roomID")
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				fmt.Println("WebSocket connection upgrade failed:", err)
				return
			}
			defer conn.Close()

			room := getOrCreateRoom(roomID)
			room.Lock()
			room.connections[conn] = true
			room.Unlock()

			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("Error reading message:", err)
					break
				}
				room.RLock()
				for conn := range room.connections {
					if conn != nil {
						err := conn.WriteMessage(websocket.TextMessage, msg)
						if err != nil {
							fmt.Println("Error writing message:", err)
						}
					}
				}
				room.RUnlock()
			}

			room.Lock()
			delete(room.connections, conn)
			room.Unlock()
		})

		router.Run()
	}

	func getOrCreateRoom(roomID string) *Room {
		roomsMutex.Lock()
		defer roomsMutex.Unlock()
		if room, ok := rooms[roomID]; ok {
			return room
		}
		room := &Room{
			connections: make(map[*websocket.Conn]bool),
		}
		rooms[roomID] = room
		return room
	}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Room struct {
	usr []*User
	sync.RWMutex
	msg []Message
}

type User struct {
	Name string
	Conn *websocket.Conn
}

type Message struct {
	msg    []byte
	isView bool
}

var rooms = make(map[string]*Room, 1)
var roomsMutex sync.Mutex

func main() {
	router := gin.Default()

	router.POST("/login", func(c *gin.Context) {
		fmt.Println("Login request received")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/validate", func(c *gin.Context) {
		fmt.Println("Validation request received")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/ws/:roomID/:userName", func(c *gin.Context) {
		roomID := c.Param("roomID")
		userName := c.Param("userName")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("WebSocket connection upgrade failed:", err)
			return
		}
		defer conn.Close()

		room := getOrCreateRoom(roomID)
		room.Lock()
		user := User{
			Name: userName,
			Conn: conn,
		}
		room.usr = append(room.usr, &user)
		room.Unlock()

		// Önceki mesajları yeni kullanıcıya gönder
		room.RLock()
		for _, msg := range room.msg {
			err := conn.WriteMessage(websocket.TextMessage, msg.msg)
			if err != nil {
				fmt.Println("Error sending previous messages:", err)
			}
		}
		room.RUnlock()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				break
			}

			room.Lock()
			message := Message{
				msg:    msg,
				isView: false,
			}
			room.msg = append(room.msg, message)

			for _, user := range room.usr {
				if user != nil {
					err := user.Conn.WriteMessage(websocket.TextMessage, msg)
					if err != nil {
						fmt.Println("Error writing message:", err)
					}
				}
			}
			room.Unlock()
		}
	})

	router.Run()
}

func getOrCreateRoom(roomID string) *Room {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()
	if room, ok := rooms[roomID]; ok {
		return room
	}
	room := &Room{
		usr: []*User{},
	}
	rooms[roomID] = room
	return room
}
*/
import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Rooms struct {
	cotegoryName string
	r            map[string]Room
}

type Room struct {
	users []*User
	name  string
	sync.RWMutex
	msg []Message
}

type User struct {
	Name string
	Conn *websocket.Conn
}
type Message struct {
	msg    []byte
	isView bool
}

var rooms = make(map[string]*Room, 10)
var roomsMutex sync.Mutex

func main() {
	router := gin.Default()

	router.Use(cors.Default()) // CORS middleware'ini ekleyin

	router.POST("/login", func(c *gin.Context) {
		fmt.Println("Login request received")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/validate", func(c *gin.Context) {
		fmt.Println("Validation request received")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/room", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"oda": rooms,
		})
	})

	router.GET("/ws/:roomID/:userName", func(c *gin.Context) {
		roomID := c.Param("roomID")
		userName := c.Param("userName")
		fmt.Println(userName + " odaya katıldı")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("WebSocket connection upgrade failed:", err)
			return
		}
		defer conn.Close()

		room := getOrCreateRoom(roomID)
		room.Lock()
		user := &User{
			Name: userName,
			Conn: conn,
		}
		room.users = append(room.users, user)
		room.Unlock()

		// Yeni kullanıcıyı diğer kullanıcılara bildir
		room.RLock()
		for _, u := range room.users {
			if u != user {
				err := u.Conn.WriteJSON(map[string]string{
					"type":     "new_user",
					"userName": userName,
				})
				for _, msg := range room.msg {
					err := conn.WriteMessage(websocket.TextMessage, msg.msg)
					if err != nil {
						fmt.Println("Error sending previous messages:", err)
					}
				}

				if err != nil {
					fmt.Println("Error notifying new user:", err)
				}
			}
		}
		room.RUnlock()

		// Mesaj dinleme ve iletme
		for {

			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				break
			}

			message := Message{
				msg:    msg,
				isView: false,
			}
			room.msg = append(room.msg, message)

			room.RLock()
			for _, u := range room.users {
				if u != user {
					err := u.Conn.WriteMessage(websocket.TextMessage, msg)
					if err != nil {
						fmt.Println("Error writing message:", err)
					}
				}
			}
			room.RUnlock()
		}

		// Kullanıcı ayrıldığında bağlantıyı kapat
		room.Lock()
		for i, u := range room.users {
			if u == user {
				room.users = append(room.users[:i], room.users[i+1:]...)
				break
			}
		}
		room.Unlock()
	})

	router.Run()
}

func getOrCreateRoom(roomID string) *Room {
	roomsMutex.Lock()
	defer roomsMutex.Unlock()
	if room, ok := rooms[roomID]; ok {
		return room
	}
	room := &Room{
		users: []*User{},
		name:  roomID,
	}
	rooms[roomID] = room
	return room
}
