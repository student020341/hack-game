package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	accountPkg "server/pkg/accounts"
	"server/pkg/models"
	testPkg "server/tests"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("account tests", func() {
	It("sanity test", func() {
		code, body, err := testPkg.SimpleGet(fmt.Sprintf("%s/status", httpServer.URL), nil)
		Expect(err).To(Succeed())
		Expect(*code).To(Equal(200))
		Expect(string(body)).To(Equal("ok"))
	})

	It("should create an account", func() {
		// create account
		input := struct {
			Username string
			Password string
		}{
			Username: "test",
			Password: "foo",
		}
		code, _, err := testPkg.SimplePost(fmt.Sprintf("%s/api/accounts", httpServer.URL), input, nil)
		Expect(err).To(Succeed())
		Expect(*code).To(Equal(200))

		// verify account exists
		var acc models.Account
		tx := s.DB.Take(&acc, "Username = ?", input.Username)
		Expect(tx.Error).To(Succeed())
		Expect(acc.Username).To(Equal(input.Username))
	})

	It("should log in to existing account", func() {
		input := struct {
			Username string
			Password string
		}{
			Username: "admin",
			Password: "admin",
		}

		url := fmt.Sprintf("%s/api/login", httpServer.URL)
		code, token, err := testPkg.SimplePost(url, input, nil)
		Expect(err).To(Succeed())
		Expect(*code).To(Equal(200))

		// verify log in token created
		var auth models.AuthSession
		tx := s.DB.Take(&auth, "Token = ?", string(token))
		Expect(tx.Error).To(Succeed())
		Expect(auth.Token).To(Equal(string(token)))
	})

	Context("logged in accounts", func() {
		var acc *models.Account
		var auth *models.AuthSession

		var adminAcc *models.Account
		var adminAuth *models.AuthSession

		var char *models.Character
		var adminChar *models.Character
		BeforeEach(func() {
			// create user and session
			var err error
			acc, err = accountPkg.CreateAccount("user", "pass")
			Expect(err).To(Succeed())

			auth, err = accountPkg.Login("pass", *acc)
			Expect(err).To(Succeed())

			// persist
			tx := s.DB.Create(&acc)
			Expect(tx.Error).To(Succeed())
			tx = s.DB.Create(&auth)
			Expect(tx.Error).To(Succeed())

			// create admin user
			adminAcc, err = accountPkg.CreateAccount("foo", "bar")
			Expect(err).To(Succeed())
			adminAuth, err = accountPkg.Login("bar", *adminAcc)
			Expect(err).To(Succeed())

			// persist
			tx = s.DB.Create(&adminAcc)
			Expect(tx.Error).To(Succeed())
			tx = s.DB.Create(&adminAuth)
			Expect(tx.Error).To(Succeed())

			// create characters for users
			char = &models.Character{
				ID:        uuid.New().String(),
				Name:      "testchar",
				AccountID: acc.ID,
			}
			tx = s.DB.Create(&char)
			Expect(tx.Error).To(Succeed())

			adminChar = &models.Character{
				ID:        uuid.New().String(),
				Name:      "adchar",
				AccountID: adminAcc.ID,
			}
			tx = s.DB.Create(&adminChar)
			Expect(tx.Error).To(Succeed())
		})

		It("user can log out", func() {
			code, _, err := testPkg.SimplePost(httpServer.URL+"/api/logout", nil, &auth.Token)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(200))

			// verify session deleted
			tx := s.DB.Take(&models.AuthSession{}, "Token = ?", &auth.Token)
			Expect(tx.Error).To(Equal(gorm.ErrRecordNotFound))
		})

		It("user can log out of all sessions", func() {
			// create a few more
			for i := 0; i < 5; i++ {
				a, err := accountPkg.Login("pass", *acc)
				Expect(err).To(Succeed())
				tx := s.DB.Create(&a)
				Expect(tx.Error).To(Succeed())
			}

			// verify there are several sessions
			var count int64
			tx := s.DB.Find(&models.AuthSession{}).Count(&count)
			Expect(tx.Error).To(Succeed())
			Expect(int(count)).To(Equal(7))

			code, _, err := testPkg.SimplePost(httpServer.URL+"/api/logout?all=1", nil, &auth.Token)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(200))

			// verify there are 0 sessions
			tx = s.DB.Find(&models.AuthSession{}).Count(&count)
			Expect(tx.Error).To(Succeed())
			Expect(int(count)).To(Equal(1))
		})

		It("regular user cannot access account list", func() {
			code, _, err := testPkg.SimpleGet(fmt.Sprintf("%s/api/accounts", httpServer.URL), nil)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(401))
		})

		It("admin can access account list", func() {
			code, _, err := testPkg.SimpleGet(fmt.Sprintf("%s/api/accounts", httpServer.URL), &adminAuth.Token)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(401))
		})

		It("user can list their characters", func() {
			code, body, err := testPkg.SimpleGet(httpServer.URL+"/api/characters", &auth.Token)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(200))

			// verify that users only get characters that belong to them
			var charList []models.Character
			err = json.Unmarshal(body, &charList)
			Expect(err).To(Succeed())
			Expect(len(charList)).To(Equal(1))
			Expect(charList[0].Name).To(Equal("testchar"))
		})

		It("user can create a character", func() {
			input := struct {
				Name string
			}{
				Name: "Foo",
			}
			code, _, err := testPkg.SimplePost(httpServer.URL+"/api/characters", input, &auth.Token)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(200))

			// verify
			var chars []models.Character
			tx := s.DB.Find(&chars, "account_id = ?", acc.ID)
			Expect(tx.Error).To(Succeed())
			Expect(len(chars)).To(Equal(2))
			Expect("Foo").To(BeElementOf([]string{chars[0].Name, chars[1].Name}))
		})

		It("user can delete a character", func() {
			input := struct {
				ID string
			}{
				ID: char.ID,
			}
			code, _, err := testPkg.SimpleDelete(httpServer.URL+"/api/characters", input, &auth.Token)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(200))

			//
			var chars []models.Character
			tx := s.DB.Find(&chars, "account_id = ?", acc.ID)
			Expect(tx.Error).To(Succeed())
			Expect(len(chars)).To(Equal(0))
		})

		It("user can list players in a server", func() {
			// put players in server
			s.Towns[0].Players = append(
				s.Towns[0].Players,
				models.Player{
					Account:   *adminAcc,
					Character: *adminChar,
					Socket:    nil,
				},
				models.Player{
					Account:   *acc,
					Character: *char,
					Socket:    nil,
				},
			)

			code, body, err := testPkg.SimpleGet(fmt.Sprintf("%s/api/servers/%s/players", httpServer.URL, s.Towns[0].ID), &auth.Token)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(200))

			var players []models.Player
			err = json.Unmarshal(body, &players)
			Expect(err).To(Succeed())
			Expect(len(players)).To(Equal(2))
		})

		// TODO move these tests to a different file
		It("user can join and leave a server", func() {
			var input struct {
				TownID      string
				CharacterID *string
			}
			input.CharacterID = &char.ID // player character
			input.TownID = s.Towns[0].ID // default server town

			// put player data in first server town
			code, _, err := testPkg.SimplePost(httpServer.URL+"/api/servers/join", input, &auth.Token)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(200))

			// verify player in server
			Expect(len(s.Towns[0].Players)).To(Equal(1))
			Expect(s.Towns[0].Players[0].Account.ID).To(Equal(acc.ID))
			Expect(s.Towns[0].Players[0].Socket).To(BeNil())

			// player got a positive response, client should create a socket connection
			uri := url.URL{Scheme: "ws", Host: httpServer.Listener.Addr().String(), Path: "/api/servers/socket"}
			c, _, err := websocket.DefaultDialer.Dial(uri.String(), http.Header{"token": []string{auth.Token}})
			Expect(err).To(Succeed())
			defer c.Close()

			Expect(s.Towns[0].Players[0].Socket).NotTo(BeNil())

			// read a message from the server
			_, message, err := c.ReadMessage()
			Expect(err).To(Succeed())
			Expect(string(message)).To(Equal("hello"))

			// leave
			code, _, err = testPkg.SimpleGet(httpServer.URL+"/api/servers/leave", &auth.Token)
			Expect(err).To(Succeed())
			Expect(*code).To(Equal(200))

			// verify
			for _, p := range s.Towns[0].Players {
				Expect(p.Account.ID).NotTo(Equal(acc.ID))
			}
		})

		It("users can chat", func() {
			// put players in server
			s.Towns[0].Players = append(
				s.Towns[0].Players,
				models.Player{
					Account:   *adminAcc,
					Character: *adminChar,
					Socket:    nil,
				},
				models.Player{
					Account:   *acc,
					Character: *char,
					Socket:    nil,
				},
			)

			// socket clients
			uri := url.URL{Scheme: "ws", Host: httpServer.Listener.Addr().String(), Path: "/api/servers/socket"}
			adminClient, _, err := websocket.DefaultDialer.Dial(uri.String(), http.Header{"token": []string{adminAuth.Token}})
			Expect(err).To(Succeed())
			defer adminClient.Close()

			userClient, _, err := websocket.DefaultDialer.Dial(uri.String(), http.Header{"token": []string{auth.Token}})
			Expect(err).To(Succeed())
			defer userClient.Close()

			// verify
			for _, t := range s.Towns {
				for _, p := range t.Players {
					Expect(p.Socket).NotTo(BeNil())
				}
			}

			// chat
			err = adminClient.WriteMessage(websocket.TextMessage, []byte("hello from admin"))
			Expect(err).To(Succeed())

			_, message, err := userClient.ReadMessage()
			Expect(err).To(Succeed())
			// read server welcome
			Expect(string(message)).To(Equal("hello"))

			// read message sent by other user
			_, message, err = userClient.ReadMessage()
			Expect(err).To(Succeed())
			Expect(string(message)).To(Equal("hello from admin"))
		})
	}) // logged in accounts
})
