package email

import (
	"fmt"
	"github.com/go-gomail/gomail"
	"log"
	"math/rand"
	"os"
	"time"
	"../redis"
)
/*
func main() {

	result := CreateCaptcha()
	fmt.Println(result)
}

 */

func SendConfirm(reciver string)  {
	conn := redisconfirm.InitRedis("csl")
	captcha := CreateCaptcha()
	redisconfirm.SetCaptcha(&conn,reciver,captcha)
	sendMail(reciver,captcha)
}

func EmailConfirm(reciver string,captcha string) bool {
	conn := redisconfirm.InitRedis("csl")
	return redisconfirm.CaptchaConfirm(&conn,reciver,captcha)
}

func sendMail(reciver string,captcha string){
	m := gomail.NewMessage()
	//passwd := "jdSpider@126.com"
	//passwd := "QMDACUYRXYFHRDUJ"
	hostMail := os.Getenv("hostMail")
	passwd := os.Getenv("mailPasswd")
	m.SetAddressHeader("From", hostMail, "jdSpider") // 发件人
	m.SetHeader("To",m.FormatAddress(reciver, "收件人")) // 收件人
	m.SetHeader("Subject", "来自jdSpier的验证码")     // 主题
	body := fmt.Sprintf("验证码：%s",captcha)
	m.SetBody("text/plain",body) // 正文

	d := gomail.NewDialer("smtp.126.com", 465, hostMail, passwd)
	//fmt.Println("send email!")
	if err := d.DialAndSend(m); err != nil {
		log.Println("发送失败", err)
		return
	}
}


func CreateCaptcha() string {
	captcha := rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000)
	return fmt.Sprintf("%06v",captcha)
}
