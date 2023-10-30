package main

import (
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nfnt/resize"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("_")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			if update.Message.Photo != nil {
				// Download the file
				file, err := bot.GetFile(tgbotapi.FileConfig{FileID: update.Message.Photo[len(update.Message.Photo)-1].FileID})
				if err != nil {
					log.Printf("Failed to download file: %v", err)
				}

				msg.Text = file.FilePath

				downloadFile(file.FilePath, file.FileID, fmt.Sprint(update.Message.From.ID), ".png")
				compressImage("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".png")

				fileP, err := os.ReadFile("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".png")
				if err != nil {
					log.Fatal(err)
				}

				v := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{
					Name:  file.FilePath,
					Bytes: fileP,
				})

				bot.Send(v)
				os.Remove("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".png")
			}

			if update.Message.Voice != nil {
				// Download the file
				file, err := bot.GetFile(tgbotapi.FileConfig{FileID: update.Message.Voice.FileID})
				if err != nil {
					log.Printf("Failed to download file: %v", err)
				}

				msg.Text = file.FilePath

				downloadFile(file.FilePath, file.FileID, fmt.Sprint(update.Message.From.ID), ".ogg")
				compressAudio("media/"+fmt.Sprint(update.Message.From.ID)+"/"+file.FileID+".ogg", ".ogg")

				fileP, err := os.ReadFile("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".ogg.output.ogg")
				if err != nil {
					log.Fatal(err)
				}

				v := tgbotapi.NewVoice(update.Message.Chat.ID, tgbotapi.FileBytes{
					Name:  file.FilePath,
					Bytes: fileP,
				})

				bot.Send(v)

				os.Remove("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".ogg")
				os.Remove("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".ogg.output.ogg")
			}

			if update.Message.Audio != nil {
				// Download the file
				file, err := bot.GetFile(tgbotapi.FileConfig{FileID: update.Message.Audio.FileID})
				if err != nil {
					log.Printf("Failed to download file: %v", err)
				}

				msg.Text = file.FilePath

				downloadFile(file.FilePath, file.FileID, fmt.Sprint(update.Message.From.ID), ".mp3")
				compressAudio("media/"+fmt.Sprint(update.Message.From.ID)+"/"+file.FileID+".mp3", ".mp3")

				fileP, err := os.ReadFile("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".mp3.output.mp3")
				if err != nil {
					log.Fatal(err)
				}

				v := tgbotapi.NewAudio(update.Message.Chat.ID, tgbotapi.FileBytes{
					Name:  file.FilePath,
					Bytes: fileP,
				})

				bot.Send(v)

				os.Remove("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".mp3")
				os.Remove("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".mp3.output.mp3")
			}

			if update.Message.VideoNote != nil {
				// Download the file
				file, err := bot.GetFile(tgbotapi.FileConfig{FileID: update.Message.VideoNote.FileID})
				if err != nil {
					log.Printf("Failed to download file: %v", err)
				}

				msg.Text = file.FilePath

				downloadFile(file.FilePath, file.FileID, fmt.Sprint(update.Message.From.ID), ".mp4")
				compressVideo("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".mp4")

				fileP, err := os.ReadFile("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".mp4.output.mp4")
				if err != nil {
					log.Fatal(err)
				}

				v := tgbotapi.NewVideoNote(update.Message.Chat.ID, 1, tgbotapi.FileBytes{
					Name:  file.FilePath,
					Bytes: fileP,
				})

				bot.Send(v)

				os.Remove("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".mp4")
				os.Remove("media/" + fmt.Sprint(update.Message.From.ID) + "/" + file.FileID + ".mp4.output.mp4")

			}
		}
	}
}

func downloadFile(filePath string, fileId string, userId string, fileExt string) {
	os.MkdirAll("media/"+userId, os.ModePerm)
	url := "https://api.telegram.org/file/bot6766292283:AAH4yr8zQz6OjCClTCxyTC6RZ0yc0y7kLYU/" + filePath
	err := DownloadFile("media/"+userId+"/"+fileId+fileExt, url)
	if err != nil {
		fmt.Println("Error downloading file: ", err)
		return
	}

	fmt.Println("Downloaded: " + url)

}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func compressImage(files string) {
	// Open the image file
	file, err := os.Open(files)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Decode the image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// Resize the image to the desired dimensions
	resizedImg := resize.Resize(uint(img.Bounds().Dx()/2), uint(img.Bounds().Dy()/2), img, resize.Bicubic)

	os.Remove(files)

	// Create a new output file
	out, err := os.Create(files)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Compress and save the image to the output file
	jpeg.Encode(out, resizedImg, &jpeg.Options{Quality: 6})
}

func compressAudio(files string, fileExt string) {
	if fileExt == ".ogg" {
		ffmpeg.Input(files).
			Output(files+".output"+fileExt, ffmpeg.KwArgs{"c:v": "libvorbis", "ab": "32k", "ar": "8000"}).
			OverWriteOutput().ErrorToStdOut().Run()
	} else {
		ffmpeg.Input(files).
			Output(files+".output"+fileExt, ffmpeg.KwArgs{"b:a": "8k", "map": "a"}).
			OverWriteOutput().ErrorToStdOut().Run()
	}
}

func compressVideo(files string) {
	ffmpeg.Input(files).
		Output(files+".output"+".mp4", ffmpeg.KwArgs{"b": "24k", "b:a": "8k"}).
		OverWriteOutput().ErrorToStdOut().Run()
}
