// TCP2 project main.go
package work

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"time"
)

type Receiver struct {
	Account string
}
type Collector struct {
	Title      string
	CommType   string
	RemoteIP   string
	RemotePort int
}
type Info struct {
	Receivers  []Receiver
	Collectors []Collector
}

const (
	SOI = byte(0x7E)
	CR  = byte(0x0D)
)

func getconn(ip string, deal string) (err error, conn net.Conn) {

	conn, err = net.DialTimeout(deal, ip, time.Second*3)

	return err, conn
}
func goWork(title string, ip string, deal string, receiver []Receiver) {
	//	conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	fmt.Println("IP is", ip)
	fmt.Print("work")
	state := 1 //状态
	var count int
	//	var err error
	for {
		if count >= 5 && state == 1 {
			fmt.Print("故障，发送邮件")
			for _, account := range receiver {
				fmt.Println("这里空的吗？", account.Account)
				body := "IP地址为" + ip + "的" + title + "发生了故障"
				go SendMail(body, account.Account) //采集程序故障
			}
			state = 2 //表示是发送邮件之后的状态
			//			time.Sleep(time.Second * 180)
			count = 0
		}
		err, conn := getconn(ip, deal)
		if err != nil {
			fmt.Println("握手失败")
			count++
			time.Sleep(time.Second * 3)
			continue
		}
		fmt.Println("握手成功")
		conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		sendData(conn)
		fmt.Println("发送数据了")
		_, err = getData(conn)
		if err != nil {
			count++
		} else {
			count = 0
		}
		if err == nil && state == 2 {
			fmt.Print("恢复了")
			for _, account := range receiver {
				body := "IP地址为" + ip + "的" + title + "恢复正常了"
				go SendMail(body, account.Account) //采集程序故障
			}
			state = 1

		}
		conn.Close()
		time.Sleep(time.Second * 3)
	}

}

//发送数据包
func sendData(conn net.Conn) {
	//	writer := bufio.NewWriterSize(conn, 8)
	data := []byte{}
	data = append([]byte{SOI})
	data = append(data, CR)
	fmt.Println(data)
	conn.Write(data)
}

//接收数据包
func getData(conn net.Conn) (data string, err error) {
	byte := make([]byte, 20)
	_, err = conn.Read(byte)
	//	fmt.Println(string(byte))
	return string(byte), err
}

//发送邮件
func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.Sendee951Mail(host, auth, user, send_to, msg)
	return err

}
func SendMail(body string, email string) {
	user := "nanjian@modoutech.com"
	password := "WZmodou@123"
	host := "smtp.exmail.qq.com:25"
	fmt.Println("email", email)
	to := email

	subject := "采集问题报告"

	fmt.Println("send email")
	err := SendToMail(user, password, host, to, subject, body, "html")
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}
}

func Start() {
	var info Info
	err := ReadFile("conf.json", &info)
	if err != nil {
		fmt.Println("readFile: ", err.Error())
		return
	}
	collectors := info.Collectors
	emails := info.Receivers
	for _, collector := range collectors {
		ip := collector.RemoteIP + ":" + fmt.Sprintf("%d", collector.RemotePort)2				

		title := collector.Title
		go goWork(title, ip, collector.CommType, emails)
	}

	//		 goWork(name string, ip string, deal string, emails []string)
	for {
		time.Sleep(time.Second * 5)
	}
}
