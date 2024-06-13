package connection

import (
	"backend/server/env"
	"backend/server/structure"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type WhatsApp struct{}

var WhatsAppClient *whatsmeow.Client // save session connection
var WhatsAppQrCode string

var rabbitmq = RabbitMQ{}

var WhatsAppMessageTopic = "whatsapp-message"
var WhatsAppMessage RabbitMqStructure = rabbitmq.CreateModel(WhatsAppMessageTopic)

var dbPath = filepath.Join(env.GetPwd(), "whatsapp.db")

func (ref WhatsApp) Connect() {
	go func() {
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			file, err := os.Create(dbPath)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer file.Close()

			fmt.Println("✅ File whatsapp.db created successfully")
		}
		dbLog := waLog.Stdout("Database", "DEBUG", true)
		container, err := sqlstore.New("sqlite3", "file:"+dbPath+"?_foreign_keys=on", dbLog)
		if err != nil {
			panic(err)
		}
		// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
		deviceStore, err := container.GetFirstDevice()
		if err != nil {
			panic(err)
		}

		// create table if not exist...
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer db.Close()
		queries := []string{
			`CREATE TABLE IF NOT EXISTS send_messages (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					jid VARCHAR NOT NULL,
					is_success BOOLEAN,
					request TEXT NOT NULL,
					response TEXT,
					error_reason TEXT,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);`,
			`CREATE INDEX IF NOT EXISTS idx_jid ON send_messages (jid);`,
			`CREATE INDEX IF NOT EXISTS idx_jid_is_success ON send_messages (jid, is_success);`,
		}
		for _, query := range queries {
			_, err := db.Exec(query)
			if err != nil {
				fmt.Println("Error executing query:", err)
			}
		}

		// fmt.Println("AAA")
		// clientLog := waLog.Stdout("Client", "DEBUG", true)
		// WhatsAppClient = whatsmeow.NewClient(deviceStore, clientLog)

		WhatsAppClient = whatsmeow.NewClient(deviceStore, nil)
		// fmt.Println("BBB")

		is_connected := false
		for {
			qrChan, _ := WhatsAppClient.GetQRChannel(context.Background())
			// fmt.Println("CCC")
			err = WhatsAppClient.Connect()
			// fmt.Println("DDD")
			if err != nil {
				fmt.Printf("Error connecting: %s\n", err.Error())
				time.Sleep(time.Second * 10) // Sleep for 10 seconds before retrying
				continue
			}
			// fmt.Println("EEE")

			if qrChan == nil {
				is_connected = true
			} else {
				for evt := range qrChan {
					event := evt.Event
					// fmt.Printf("\n\n Event: %+v\n \n\n", event)
					if event == "code" {
						code := evt.Code

						tempDir := filepath.Join(env.GetPwd(), "temp")
						uuidv4 := uuid.NewString()
						tempPath := filepath.Join(tempDir, uuidv4+".png")
						err := qrcode.WriteFile(code, qrcode.Medium, 256, tempPath)
						if err != nil {
							fmt.Println("error on write file qr-code image:", err.Error())
						} else {
							qrFile, err := os.Open(tempPath)
							if err != nil {
								fmt.Println("error on open qr-code:", err.Error())
							}
							defer qrFile.Close()

							stat, _ := qrFile.Stat()
							size := stat.Size()
							qrBytes := make([]byte, size)
							_, err = qrFile.Read(qrBytes)
							if err != nil {
								fmt.Println("error on read file qr-code image:", err.Error())
							}
							WhatsAppQrCode = base64.StdEncoding.EncodeToString(qrBytes)

							os.Remove(tempPath)
						}

						// qrterminal.GenerateHalfBlock(code, qrterminal.L, os.Stdout) // on terminal
					} else if event == "timeout" {
						log.Println("🔃 Re-connecting QR Code...")
						break
					} else if event == "success" {
						is_connected = true
						break
					} else {
						fmt.Println("Log event:", event)
					}
				}
			}

			if is_connected {
				WhatsAppQrCode = ""
				break
			}
		}

		if is_connected {
			ref.Listener()
			ref.Consumer()
		}

		log.Println("✅ WhatsApp Success Connected...")
	}()
}

func (ref WhatsApp) Disconnect() {
	if WhatsAppClient.IsConnected() && WhatsAppClient.IsLoggedIn() {
		WhatsAppClient.Disconnect()
	}
	fmt.Println("WhatsApp Disconnect ✅")
}

func (ref WhatsApp) Listener() {
	WhatsAppClient.AddEventHandler(func(evt interface{}) {
		// switch v := evt.(type) {
		// case *events.Message:
		// 	fmt.Printf("\n\n Received a Message! %+v\n\n", v)
		// // case *events.Receipt:
		// // 	fmt.Printf("\n\n Received a Receipt! %+v\n\n", v)
		// case *events.ConnectFailure:
		// 	fmt.Printf("\n\n Received a ConnectFailure! %+v\n\n", v)
		// case *events.Disconnected:
		// 	fmt.Printf("\n\n Received a Disconnected! %+v\n\n", v)
		// case *events.Picture:
		// 	fmt.Printf("\n\n Received a Picture! %+v\n\n", v)
		// case *events.Presence:
		// 	fmt.Printf("\n\n Received a Presence! %+v\n\n", v)
		// }
	})
}

func (ref WhatsApp) Consumer() {

	go func() {
		RabbitMQ := RabbitMQ{}
		msgs, Connection, Channel := RabbitMQ.CreateConsumer(WhatsAppMessageTopic)
		defer Connection.Close()
		defer Channel.Close()

		for d := range msgs {
			var err error
			body := d.Body

			var data structure.IWhatsAppSendQueueRabbitMQ
			err = json.Unmarshal(body, &data)
			if err != nil {
				log.Printf("Error deserializing message: %s", err)
				d.Ack(false) // acknowledge the message because data is not object
				continue
			}

			_type := data.Type
			target_number := data.TargetNumber
			if _type == "text" {
				message := *data.Message
				recipient, err := types.ParseJID(target_number)
				if err != nil {
					log.Printf("Error parsing JID: %s", err)
					d.Ack(false) // acknowledge the message
					continue
				}

				resp, err := WhatsAppClient.SendMessage(context.Background(), recipient, &waProto.Message{
					Conversation: proto.String(message),
				})
				if err != nil {
					log.Printf("Error sending message: %s", err)
					err = d.Nack(false, false) // multiple set to false, requeue set to false
					if err != nil {
						log.Printf("Error sending Nack: %s", err)
					}
					continue
				}

				fmt.Printf("\nresp: %+v\nmessage: %+v\n\n", resp, message) // debug...
			} else if _type == "image" || _type == "file" {
				if data.FileName == nil {
					log.Println("FileName not found")
					d.Ack(false) // acknowledge the message
					continue
				}

				message := ""
				if data.Message != nil {
					message = *data.Message
				}
				filename := *data.FileName
				recipient, err := types.ParseJID(target_number)
				if err != nil {
					log.Printf("Error parsing JID: %s", err)
					d.Ack(false) // acknowledge the message
					continue
				}
				fmt.Println(
					"type:", _type,
					"| target_number:", target_number,
					"| filename:", filename,
					"| message:", message,
				)

				tempFile := filepath.Join(env.GetPwd(), "temp", filename)
				data, err := os.ReadFile(tempFile) // data
				if err != nil {
					log.Printf("Error ReadFile: %s", err)
					d.Ack(false) // acknowledge the message
					continue
				}
				fmt.Printf("tempFile: %+v\n", tempFile)

				uploaded, err := WhatsAppClient.Upload(context.Background(), data, whatsmeow.MediaImage)
				if err != nil {
					log.Printf("failed upload %s to whatsapp server", _type)
					d.Ack(false) // acknowledge the message
					continue
				}

				msg := &waProto.Message{ImageMessage: &waProto.ImageMessage{
					Caption:       proto.String(message),
					URL:           proto.String(uploaded.URL),
					DirectPath:    proto.String(uploaded.DirectPath),
					MediaKey:      uploaded.MediaKey,
					Mimetype:      proto.String(http.DetectContentType(data)),
					FileEncSHA256: uploaded.FileEncSHA256,
					FileSHA256:    uploaded.FileSHA256,
					FileLength:    proto.Uint64(uint64(len(data))),
				}}
				resp, err := WhatsAppClient.SendMessage(context.Background(), recipient, msg)
				if err != nil {
					log.Printf("Error sending image message: %v", err)
					err = d.Nack(false, false) // multiple set to false, requeue set to false
					if err != nil {
						log.Printf("Error sending Nack: %s", err)
					}
					continue
				} else {
					log.Printf("Image message sent (server timestamp: %s)", resp.Timestamp)
				}
			} else {
				fmt.Println(
					"type:", _type,
					"| target_number:", target_number,
				)
			}

			d.Ack(true) // finished...
			time.Sleep(3 * time.Second)
		}
	}()

}
