package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/utils"
	"github.com/gomodule/redigo/redis"
	"math/rand"
	"pyg/models"
	"regexp"
	"strconv"
	"time"
)

type UserController struct {
	beego.Controller
}

func ErrResp(this *UserController, errmsg string, fileName string) {
	this.Data["errmsg"] = errmsg
	this.TplName = fileName
}

func (this *UserController) ShowRegister() {

	this.TplName = "register.html"
}
func (this *UserController) HandleRegister() {
	phone := this.GetString("phone")
	fmt.Println(phone)
	code := this.GetString("code")
	password := this.GetString("password")
	repassword := this.GetString("repassword")
	if phone == "" || code == "" || password == "" || repassword == "" {
		ErrResp(this, "输入不能为空", "register.html")
		fmt.Println("输入不能为空")
		return
	}
	if repassword != password {
		ErrResp(this, "两次密码输入不一致", "register.html")
		return
	}
	reg, _ := regexp.Compile(`^[1]+[3,8]+\d{9}$`)
	result := reg.FindString(phone)
	if result == "" {
		ErrResp(this, "手机号输入有误", "register.html")
		fmt.Println("手机号输入有无")
		return
	}
	conn, err := redis.Dial("tcp", "192.168.137.130:6379")
	if err != nil {
		resp := make(map[string]interface{})
		resp["errMsg"] = "redis数据库连接失败"
		resp["statusCode"] = 401
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}
	defer conn.Close()
	redisCode, err := redis.String(conn.Do("get", phone+"_code"))
	if redisCode != code {
		resp := make(map[string]interface{})
		resp["errMsg"] = "验证码错误"
		resp["statusCode"] = 401
		this.Data["json"] = resp
		this.ServeJSON()
		return
	}

	user := new(models.User)
	o := orm.NewOrm()
	user.Name = phone
	user.PassWord = password
	id, _ := o.Insert(user)

	this.Redirect("/active?id="+strconv.Itoa(int(id)), 302)
}
func (this *UserController) SendMsg() {
	phone := this.GetString("phone")
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(899999) + 100000
	strcode := strconv.Itoa(code)
	conn, err := redis.Dial("tcp", "192.168.137.130:6379")
	if err != nil {
		resp := make(map[string]interface{})
		resp["statusCode"] = 401
		resp["errmsg"] = "redis数据库连接失败"
		this.Data["json"] = resp
		this.ServeJSON()
		fmt.Println("redis.Dial err:", err)
		return
	}
	defer conn.Close()
	conn.Do("setex", phone+"_code", 60*5, strcode)

	fmt.Println(strcode)
	/*	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", "LTAI4FkgoSqad9KWK9sqcUfD", "ArTgvcOa2Rx0LPEEp4ewNJmCcLpxQY")
		request := dysmsapi.CreateSendSmsRequest()
		request.Method = "POST"
		request.Scheme = "https"
		request.Domain = "dysmsapi.aliyuncs.com"
		request.PhoneNumbers = phone
		request.SignName = "优尚之家"
		request.TemplateCode = "SMS_176530738"
		request.TemplateParam = `{"code":` + strcode + `}`
		response, err := client.SendSms(request)
		if err != nil {
			fmt.Print(err.Error())
		}
		fmt.Printf("response is %#v\n", response)*/

	resp := make(map[string]interface{})
	resp["statusCode"] = 200
	resp["errmsg"] = "OK"
	this.Data["json"] = resp
	this.ServeJSON()

	//conn,err:=redis.Dial("tcp","192.168.137.130:6379")
	//conn.Do("setex",code,)

}
func (this *UserController) ShowActive() {
	id := this.GetString("id")
	if id == "" {
		return
	}
	this.Data["id"] = id
	this.TplName = "register-email.html"
}
func (this *UserController) HandelActive() {
	email := this.GetString("email")
	id, err := this.GetInt("id")
	if err != nil || email == "" {
		this.Redirect("/active?id="+strconv.Itoa(id), 302)
		return
	}
	reg, err := regexp.Compile(`^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$`)
	if err != nil {
		fmt.Println("regexp.Compile err:", err)
		return
	}
	result := reg.FindString(email)
	fmt.Println(result)
	if result == "" {
		this.Redirect("/active?id="+strconv.Itoa(id), 302)
		return
	}
	config := `{"username":"799975844@qq.com","password":"bxplrxksimhnbfjb","host":"smtp.qq.com","port":587}`
	sendEmail := utils.NewEMail(config)
	sendEmail.From = "799975844@qq.com"
	sendEmail.To = []string{result}
	sendEmail.Subject = "品邮购用户激活"
	sendEmail.HTML = `<a href="http://192.168.137.130:8080/activeUser?id=` + strconv.Itoa(id) + "&email=" + email + `">点击激活用户</a>`
	err = sendEmail.Send()
	if err != nil {
		fmt.Println(err)
		return
	}
	this.Data["result"] = "邮件发送成功,请去目标邮箱激活"
	this.TplName = "email-result.html"
}
func (this *UserController) ActivateTheSuccess() {
	id, err := this.GetInt("id")
	email := this.GetString("email")
	if err != nil || email == "" {
		this.TplName = "register.html"
		return
	}
	user := new(models.User)
	o := orm.NewOrm()
	user.Id = id
	err = o.Read(user)
	if err != nil {
		fmt.Println("激活用户不存在")
		this.TplName = "register.html"
		return
	}
	user.Active = true
	user.Email = email
	_, err = o.Update(user)
	if err != nil {
		fmt.Println("激活用户失败")
		this.TplName = "register.html"
		return
	}
	this.Redirect("/login", 302)
}
func (this *UserController) ShowLogin() {
	this.TplName = "login.html"
}
func (this *UserController) HandleLogin() {
	loginName := this.GetString("loginname")
	password := this.GetString("password")
	if loginName == "" || password == "" {
		this.Redirect("/login", 302)
		return
	}
	user := new(models.User)
	user.Name = loginName
	o := orm.NewOrm()
	err := o.Read(user, "Name")
	if err != nil {
		this.Redirect("/login?errmsg=用户名或密码错误", 302)
		return
	}
	if user.PassWord != password {
		this.Redirect("/login?errmsg=用户名或密码错误", 302)
		return
	}
	if user.Active != true {
		this.Redirect("/login?errmsg=邮箱未激活,请先激活邮箱", 302)
		return
	}

	checked:=this.GetString("m1")
	fmt.Println(checked)



	this.Redirect("/", 302)

}
