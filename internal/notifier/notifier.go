package notifier

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"notification-bot/internal/models"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramNotifier(bot *tgbotapi.BotAPI, chatID int64) *TelegramNotifier {
	return &TelegramNotifier{
		bot:    bot,
		chatID: chatID,
	}
}
// func (t *TelegramNotifier) Notify(news models.News) error {
// 	caption := formatMessage(news)

// 	// если картинок нет → обычное сообщение
// 	if len(news.ImagesLinks) == 0 {
// 		msg := tgbotapi.NewMessage(t.chatID, caption)
// 		msg.ParseMode = "Markdown"

// 		_, err := t.bot.Send(msg)
// 		return err
// 	}

// 	var media []interface{}

// 	for i, imgURL := range news.ImagesLinks {
// 		log.Printf("downloading image #%d: %s", i+1, imgURL)

// 		resp, err := http.Get(imgURL)
// 		if err != nil {
// 			log.Printf("failed download image #%d: %v", i+1, err)
// 			continue
// 		}

// 		if resp.StatusCode != http.StatusOK {
// 			log.Printf("bad status for image #%d: %d", i+1, resp.StatusCode)
// 			resp.Body.Close()
// 			continue
// 		}

// 		contentType := resp.Header.Get("Content-Type")
// 		if !strings.HasPrefix(contentType, "image/") {
// 			log.Printf("bad content-type for image #%d: %s", i+1, contentType)
// 			resp.Body.Close()
// 			continue
// 		}

// 		photo := tgbotapi.NewInputMediaPhoto(
// 			tgbotapi.FileReader{
// 				Name:   fmt.Sprintf("image_%d.jpg", i+1),
// 				Reader: resp.Body,
// 			},
// 		)

// 		// caption только у первой валидной картинки
// 		if len(media) == 0 {
// 			photo.Caption = caption
// 			photo.ParseMode = "Markdown"
// 		}

// 		media = append(media, photo)

// 		// Telegram максимум 10 файлов
// 		if len(media) == 10 {
// 			break
// 		}
// 	}

// 	// если все картинки битые то отправляем просто текст
// 	if len(media) == 0 {
// 		msg := tgbotapi.NewMessage(t.chatID, caption)
// 		msg.ParseMode = "Markdown"

// 		_, err := t.bot.Send(msg)
// 		return err
// 	}

// 	msg := tgbotapi.NewMediaGroup(t.chatID, media)

// 	_, err := t.bot.SendMediaGroup(msg)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
func (t *TelegramNotifier) Notify(news models.News) error {
    if len(news.ImagesLinks) == 0 {
		caption := formatMessage(news, false)
        msg := tgbotapi.NewMessage(t.chatID, caption)
        msg.ParseMode = "Markdown"
        _, err := t.bot.Send(msg)
        return err
    }
	caption := formatMessage(news, true)
    var media []interface{}

    for i, imgURL := range news.ImagesLinks {
        log.Printf("downloading image #%d: %s", i+1, imgURL)

        resp, err := http.Get(imgURL)
        if err != nil {
            log.Printf("failed download image #%d: %v", i+1, err)
            continue
        }

        if resp.StatusCode != http.StatusOK {
            log.Printf("bad status for image #%d: %d", i+1, resp.StatusCode)
            resp.Body.Close()
            continue
        }

        contentType := resp.Header.Get("Content-Type")
        if !strings.HasPrefix(contentType, "image/") {
            log.Printf("bad content-type for image #%d: %s", i+1, contentType)
            resp.Body.Close()
            continue
        }

        // ✅ Читаем всё тело в буфер — это и есть фикс
        imgData, err := io.ReadAll(resp.Body)
        resp.Body.Close()
        if err != nil {
            log.Printf("failed to read image #%d body: %v", i+1, err)
            continue
        }

        photo := tgbotapi.NewInputMediaPhoto(
            tgbotapi.FileReader{
                Name:   fmt.Sprintf("image_%d.jpg", i+1),
                Reader: bytes.NewReader(imgData),
            },
        )

        if len(media) == 0 {
            photo.Caption = caption
            photo.ParseMode = "Markdown"
        }

        media = append(media, photo)

        if len(media) == 10 {
            break
        }
    }

    if len(media) == 0 {
        msg := tgbotapi.NewMessage(t.chatID, caption)
        msg.ParseMode = "Markdown"
        _, err := t.bot.Send(msg)
        return err
    }

    msg := tgbotapi.NewMediaGroup(t.chatID, media)
    _, err := t.bot.SendMediaGroup(msg)
    return err
}
func formatMessage(n models.News, hasImages bool) string {
    maxCaption := 4096
    if hasImages {
        maxCaption = 1024
    }

    // сам текст занимает не весь лимит — нужно учесть остальные поля
    overhead := len(fmt.Sprintf("🏛️ *%s*\n📰 *%s*\n\n📅 %s\n\n\n\n📎 %s",
        n.Source, n.Title, n.Date, n.Link))
    textLimit := maxCaption - overhead

    text := truncateByParagraphs(n.Text, textLimit)

    return fmt.Sprintf(
        "🏛️ *%s*\n📰 *%s*\n\n📅 %s\n\n%s\n\n📎 %s",
        n.Source,
        n.Title,
        n.Date,
        text,
        n.Link,
    )
}

func truncateByParagraphs(text string, max int) string {
    if len(text) <= max {
        return text
    }

    paragraphs := strings.Split(text, "\n\n")
    var result strings.Builder

    for _, para := range paragraphs {
        para = strings.TrimSpace(para)
        if para == "" {
            continue
        }

        // +3 для "..."
        if result.Len()+len(para)+3 > max {
            break
        }

        if result.Len() > 0 {
            result.WriteString("\n\n")
        }
        result.WriteString(para)
    }

    // если даже первый абзац не влез — обрываем по предложению
    if result.Len() == 0 {
        return truncateBySentence(text, max)
    }

    return result.String() + "..."
}

func truncateBySentence(text string, max int) string {
    if len(text) <= max {
        return text
    }

    // ищем последнюю точку/восклицательный/вопросительный знак до лимита
    cutAt := -1
    for i, ch := range text[:max] {
        if ch == '.' || ch == '!' || ch == '?' {
            cutAt = i + 1
        }
    }

    if cutAt == -1 {
        // совсем нет знаков препинания — обрываем по слову
        cutAt = strings.LastIndex(text[:max], " ")
    }

    if cutAt <= 0 {
        return text[:max] + "..."
    }

    return strings.TrimSpace(text[:cutAt]) + "..."
}