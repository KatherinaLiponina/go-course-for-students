package httpfiber

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/KatherinaLiponina/validation"

	"homework6/internal/ads"
	"homework6/internal/app"
)

type validationStruct struct {
	Title string `validate:"title"`
	Text string `validate:"text"`
}

func newValidationStruct(title string, text string) validationStruct {
	return validationStruct{Title: title, Text: text}
}

// Метод для создания объявления (ad)
func createAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody createAdRequest
		err := c.BodyParser(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		var ad * ads.Ad
		ad, err = a.CreateAd(reqBody.Title, reqBody.Text, reqBody.UserID)

		if errors.Is(err, app.ErrBadRequest) {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		} else if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}
		return c.JSON(AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody changeAdStatusRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		var ad * ads.Ad
		ad, err = a.ChangeAdStatus(int64(adID), reqBody.UserID, reqBody.Published)

		if errors.Is(err, app.ErrForbidden) {
			c.Status(http.StatusForbidden)
			return c.JSON(AdErrorResponse(err))
		} else if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody updateAdRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		var ad * ads.Ad
		ad, err = a.UpdateAd(int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)

		if errors.Is(err, app.ErrForbidden) {
			c.Status(http.StatusForbidden)
			return c.JSON(AdErrorResponse(err))
		} else if errors.Is(err, app.ErrBadRequest) {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		} else if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}
