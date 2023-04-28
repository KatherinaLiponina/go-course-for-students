package httpgin

import (
	"encoding/json"
	"homework9/internal/ads"
	"homework9/internal/app"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAdByID(a app.App) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		ad, err := a.GetAdByID(int64(id))
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, adResponse{*ad})
	}

	return gin.HandlerFunc(fn)
}

func CreateAd(a app.App) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		var data createAdRequest
		err = json.Unmarshal(body, &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ad, err := a.CreateAd(data.Title, data.Text, data.UserID)
		if err != nil {
			if err == app.ErrBadRequest {
				c.Status(http.StatusBadRequest)
			} else {
				c.Status(http.StatusNotFound)
			}
			return
		}
		c.JSON(http.StatusOK, adResponse{*ad})
	}
	return gin.HandlerFunc(fn)
}

func ChangeAdStatus(a app.App) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		var data changeAdStatusRequest
		err = json.Unmarshal(body, &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ad, err := a.ChangeAdStatus(int64(id), data.UserID, data.Published)
		if err != nil {
			if err == app.ErrForbidden {
				c.Status(http.StatusForbidden)
			} else {
				c.Status(http.StatusNotFound)
			}
			return
		}
		c.JSON(http.StatusOK, adResponse{*ad})
	}
	return gin.HandlerFunc(fn)
}

func UpdateAd(a app.App) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		var data updateAdRequest
		err = json.Unmarshal(body, &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		ad, err := a.UpdateAd(int64(id), data.UserID, data.Title, data.Text)
		if err != nil {
			if err == app.ErrForbidden {
				c.Status(http.StatusForbidden)
			} else if err == app.ErrBadRequest {
				c.Status(http.StatusBadRequest)
			} else {
				c.Status(http.StatusNotFound)
			}
			return
		}
		c.JSON(http.StatusOK, adResponse{*ad})
	}
	return gin.HandlerFunc(fn)
}

func Select(a app.App) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusOK, adsResponse{a.Select()})
			return
		}
		var data selectAdRequest
		err = json.Unmarshal(body, &data)
		if err != nil {
			c.JSON(http.StatusOK, adsResponse{a.Select()})
			return
		}
		var arr []ads.Ad
		if data.ByAuthor {
			arr, err = a.SelectByAuthor(data.AuthorID)
		} else if data.ByCreation {
			arr = a.SelectByCreation(data.CreationTime)
		} else if data.All {
			arr = a.SelectAll()
		} else {
			arr = a.Select()
		}
		if err != nil {
			c.Status(http.StatusBadRequest)
		}
		c.JSON(http.StatusOK, adsResponse{arr})
	}
	return gin.HandlerFunc(fn)
}

func CreateUser(a app.App) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		var data createOrUpdateUser
		err = json.Unmarshal(body, &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		usr := a.CreateUser(data.Nickname, data.Email)
		c.JSON(http.StatusOK, userResponse{*usr})
	}
	return gin.HandlerFunc(fn)
}

func UpdateUser(a app.App) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		var data createOrUpdateUser
		err = json.Unmarshal(body, &data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
		usr, err := a.UpdateUser(int64(id), data.Nickname, data.Email)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, userResponse{*usr})
	}
	return gin.HandlerFunc(fn)
}

func FindAdByTitle(a app.App) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		title := c.Query("title")
		if title == "" {
			c.Status(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, adsResponse{a.FindByTitle(title)})
	}

	return gin.HandlerFunc(fn)
}

func GetUserByID(a app.App) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		usr, err := a.GetUserByID(int64(id))
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, userResponse{*usr})
	}

	return gin.HandlerFunc(fn)
}

func DeleteUserByID(a app.App) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		usr, err := a.DeleteUser(int64(id))
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusOK, userResponse{*usr})
	}

	return gin.HandlerFunc(fn)
}

func DeleteAdByID(a app.App) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		author := c.Query("author")
		if author == "" {
			c.Status(http.StatusBadRequest)
			return
		}
		authorID, err := strconv.Atoi(author)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		ad, err := a.DeleteAd(int64(id), int64(authorID))
		if err != nil {
			if err == app.ErrForbidden {
				c.Status(http.StatusForbidden)
			} else if err == app.ErrBadRequest {
				c.Status(http.StatusBadRequest)
			} else {
				c.Status(http.StatusNotFound)
			}
			return
		}
		c.JSON(http.StatusOK, adResponse{*ad})
	}
	return gin.HandlerFunc(fn)
}
