package apis

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"simpleRest/app"
	"simpleRest/models"
	"simpleRest/utils"
	"strconv"
	"time"
)

type userService interface {
	New(rs app.RequestScope, user *models.User) error
	Verify(rs app.RequestScope, code string, userId int) error
	ForgotPassword(rs app.RequestScope, email string) (*models.User, error)
	ResetPassword(rs app.RequestScope, code, pass, confirmPass string, userId int) error
	Login(rs app.RequestScope, email, password string) (*models.User, error)
}

type userResource struct {
	service userService
}

func ServeUserResource(rg *gin.RouterGroup, service userService) {
	r := &userResource{service}
	rg.POST("/users/login", r.login)
	rg.POST("/users/new", r.new)
	rg.GET("/users/verify/:code", r.verify)
	rg.POST("/users/forgot_password", r.forgotPassword)
	rg.POST("/users/reset/:code", r.resetPassword)
}

func (r *userResource) new(c *gin.Context) {
	var user models.User
	dec := json.NewDecoder(c.Request.Body)
	err := dec.Decode(&user)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "check the request body"})
		return
	}

	rs := app.GetRequestScope(c)
	err = r.service.New(rs, &user)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "check the request body"})
		return
	}

	SendEmailCode(rs, user.BusinessEmail, user.EmailCode, user.Id)
	c.JSON(http.StatusOK, &user)
}

func (r *userResource) verify(c *gin.Context) {
	code := c.Param("code")
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "user_id must be a number"})
		return
	}

	err = r.service.Verify(app.GetRequestScope(c), code, userId)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"code": "BAD_REQUEST", "message": err})
		return
	}
	c.Status(http.StatusOK)
}

func (r *userResource) forgotPassword(c *gin.Context) {
	var req struct {
		BusinessEmail string `json:"business_email"`
	}
	dec := json.NewDecoder(c.Request.Body)
	err := dec.Decode(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "check the request body"})
		return
	}
	rs := app.GetRequestScope(c)
	user, err := r.service.ForgotPassword(rs, req.BusinessEmail)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "check the request body"})
		return
	}

	SendPasswordCode(rs, user.BusinessEmail, user.PasswordCode, user.Id)
	c.Status(http.StatusOK)
}

func (r *userResource) resetPassword(c *gin.Context) {
	code := c.Param("code")
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "user_id must be a number"})
		return
	}
	var req struct {
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	dec := json.NewDecoder(c.Request.Body)
	err = dec.Decode(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "check the request body"})
		return
	}
	rs := app.GetRequestScope(c)
	err = r.service.ResetPassword(rs, code, req.Password, req.ConfirmPassword, userId)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "check the request body"})
		return
	}
	c.Status(http.StatusOK)
}

func (r *userResource) login(c *gin.Context) {
	var req struct {
		BusinessEmail string `json:"business_email"`
		Password      string `json:"password"`
	}
	dec := json.NewDecoder(c.Request.Body)
	err := dec.Decode(&req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "check the request body"})
		return
	}
	rs := app.GetRequestScope(c)
	user, err := r.service.Login(rs, req.BusinessEmail, req.Password)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "BAD_REQUEST", "message": "check the request body"})
		return
	}
	token := utils.CreateToken(app.Config.JWTSigningMethod, app.Config.JWTSigningKey, user.Id, time.Hour*72)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func JWT(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.Abort()
		c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "token is not valid"})
		return
	}
	userId, err := utils.ParseToken(app.Config.JWTVerificationKey, token)
	if err != nil {
		c.Abort()
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"code": "UNAUTHORIZED", "message": "token is not valid"})
		return
	}
	app.GetRequestScope(c).SetUserId(userId)
}
