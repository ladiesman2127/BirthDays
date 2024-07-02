package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ladiesman2127/birthdays/internal/app/models"
	"github.com/ladiesman2127/birthdays/internal/database"
	"github.com/ladiesman2127/birthdays/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UsersController struct {
	usersDB *database.UserDB
}

var validate = validator.New()

func NewUsersController(db *database.UserDB) *UsersController {
	return &UsersController{usersDB: db}
}

func (controller *UsersController) SignUp(ctx *gin.Context) {
	user := models.User{}
	users := controller.usersDB.UsersCollection()
	// Bind request body to User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}
	if !utils.DateValid(user.BirthDay) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "date should in format dd/mm/yyyy",
		})
	}
	// Validate request fields
	if err := validate.Struct(user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Checking for duplicates
	if numberDuplicates, _ := users.CountDocuments(ctx.Request.Context(), bson.M{"phone": user.Phone}); numberDuplicates > 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "A user with a given phone number already exists",
		})
		return
	}
	if loginDuplicates, _ := users.CountDocuments(ctx.Request.Context(), bson.M{"login": user.Login}); loginDuplicates > 0 {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "A user with a given login already exists",
		})
		return
	}

	// Generating Object Id
	user.ID = primitive.NewObjectID()

	// Hashing password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash the password",
		})
		return
	}
	user.Password = hashedPassword

	// Inserting to DB
	if _, err := users.InsertOne(context.Background(), user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Result OK 200
	ctx.JSON(http.StatusOK, user.ID.Hex())
}

func (controller *UsersController) Auth(ctx *gin.Context) {
	// Checking request format
	credentials := models.Credentials{}
	if err := ctx.BindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Finding user by login
	user := models.User{}
	users := controller.usersDB.UsersCollection()
	if err := users.FindOne(ctx.Request.Context(), bson.M{"login": credentials.Login}).Decode(&user); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Checking if found user's password matches with passed one
	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*credentials.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Password is wrong",
		})
		return
	}
	// Replacing user with user that have session
	session, err := utils.NewSession(user.Login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	newUser := user
	newUser.Session = session
	_, err = users.ReplaceOne(ctx.Request.Context(), user, newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Success
	ctx.JSON(http.StatusOK, session)
}

func (controller *UsersController) GetUser(ctx *gin.Context) {
	user := models.User{}
	users := controller.usersDB.UsersCollection()
	userID, err := primitive.ObjectIDFromHex(ctx.Params.ByName("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalie id",
		})
	}
	err = users.FindOne(ctx.Request.Context(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (controller *UsersController) GetUsers(ctx *gin.Context) {
	var allUsers []models.User
	users := controller.usersDB.UsersCollection()
	cur, err := users.Find(ctx.Request.Context(), bson.D{})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for cur.Next(ctx.Request.Context()) {
		user := models.User{}
		if err := cur.Decode(&user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		allUsers = append(allUsers, user)
	}

	if err := cur.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	cur.Close(ctx.Request.Context())
	ctx.JSON(http.StatusOK, gin.H{
		"data": allUsers,
	})

}

func (controller *UsersController) GetBirthdayNotifications(ctx *gin.Context) {
	users := controller.usersDB.UsersCollection()
	user, ok := controller.getCurrentUser(ctx, users)
	if !ok {
		return
	}
	_, curMonth, curDay := time.Now().Date()
	var friendWithBirthday []*models.User
	for _, friendID := range user.Subscriptions {
		friend, ok := controller.getUser(friendID.Hex(), ctx, users)
		if !ok {
			continue
		}
		friendDay, friendMonth, _, ok := utils.ParseDate(friend.BirthDay)
		if !ok {
			continue
		}
		if *friendDay == curDay && *friendMonth == int(curMonth) {
			friendWithBirthday = append(friendWithBirthday, friend)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"result": friendWithBirthday,
	})
}

func (controller *UsersController) AddFriend(ctx *gin.Context) {
	controller.configureSubscriptions(ctx, false)
}

func (controller *UsersController) RemoveFriend(ctx *gin.Context) {
	controller.configureSubscriptions(ctx, true)
}

func (controller *UsersController) configureSubscriptions(ctx *gin.Context, remove bool) {
	users := controller.usersDB.UsersCollection()
	user, ok := controller.getCurrentUser(ctx, users)
	if !ok {
		return
	}
	friend, ok := controller.getUser(ctx.Params.ByName("id"), ctx, users)
	if !ok {
		return
	}
	newUser := user
	if remove {
		var newSubscribtions []primitive.ObjectID
		for _, curFriend := range user.Subscriptions {
			if friend.ID != curFriend {
				newSubscribtions = append(newSubscribtions, friend.ID)
			}
		}
		newUser.Subscriptions = newSubscribtions
	} else {
		newUser.Subscriptions = append(newUser.Subscriptions, friend.ID)
	}
	_, err := users.ReplaceOne(ctx.Request.Context(), bson.M{"_id": user.ID}, newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"result": "ok",
	})
}

func (controller *UsersController) getCurrentUser(ctx *gin.Context, users *mongo.Collection) (*models.User, bool) {
	userAccessToken := ctx.Request.Header.Get("token")
	parsedToken, err := utils.ParseToken(&userAccessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return nil, false
	}
	if !parsedToken.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "token has expired",
		})
		return nil, false
	}
	user := models.User{}
	userLogin, err := parsedToken.Claims.GetSubject()
	users.FindOne(ctx.Request.Context(), bson.M{"login": userLogin}).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return nil, false
	}
	return &user, true
}

func (controller *UsersController) getUser(id string, ctx *gin.Context, users *mongo.Collection) (*models.User, bool) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "invalid id",
		})
		return nil, false
	}
	user := models.User{}
	users.FindOne(ctx.Request.Context(), bson.M{"_id": userID}).Decode(&user)
	return &user, true
}
