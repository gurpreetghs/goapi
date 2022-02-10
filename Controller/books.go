package Controller

import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"main.go/Models"
	"os"
	"time"
	"strings"
	"fmt"
	//"strconv"

)

func FindBooks(c *gin.Context) {
	var books Models.Book
	Models.DB.Find(&books)

	c.JSON(200, gin.H{"data": books})
}

type CreateBookInput struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Username string `json "username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateBook(c *gin.Context) {

	var input CreateBookInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Create book
	book := Models.Book{Title: input.Title, Author: input.Author, Username: input.Username, Password: input.Password}
	Models.DB.Create(&book)

	c.JSON(200, gin.H{"data": book})
}

func FindBook(c *gin.Context) {
	var book Models.Book

	if err := Models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(400, gin.H{"error": "Record not found"})
		return
	}
	c.JSON(200, gin.H{"data": book})
}

type UpdateBookInput struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

func UpdateBook(c *gin.Context) {
	// Get model if exist
	var book Models.Book
	if err := Models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	Models.DB.Table("books").Where("id = ?", c.Param("id")).Updates(&input)

	c.JSON(http.StatusOK, gin.H{"data": input})
}

//Delete data
func DeleteBook(c *gin.Context) {
	var book Models.Book
	if err := Models.DB.Where("id = ?", c.Param("id")).First(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error Not found !"})
		return
	}
	Models.DB.Delete(&book)

	c.JSON(200, gin.H{"data": true})
}


//Login cradintal
type CheckBookInput struct {
	Username string `json "username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {

	var input CheckBookInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	var book Models.Book
	
	if err := Models.DB.Where("username = ? AND password = ?",input.Username,input.Password).First(&book).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please provide valid login detail"})
		return 
	}

	token, err := CreateToken(book.ID)
  		if err != nil {
     	c.JSON(http.StatusUnprocessableEntity, err.Error())
     	return
  	}
	  
	  c.JSON(http.StatusOK, gin.H{"access token":token})
}





func CreateToken(userId uint64) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") 
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
	   return "", err
	}
	return token, nil
  }




func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}


func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("token")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
	   return strArr[1]
	}
	return ""
  }

  func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
	   return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
	   return err
	}
	return nil
  }
  