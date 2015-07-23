package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/nmcclain/ldap"
	"github.com/vaughan0/go-ini"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	. "strings"
	. "time"
)

func read_ini(fname string, sec string) ini.Section {
	dir, _ := os.Getwd()
	file, _ := ini.LoadFile(dir + "\\" + fname)
	return file[sec]
}

func writef(text string, namefile string) string {
	d1 := []byte(text)
	ioutil.WriteFile(os.TempDir()+"\\"+namefile, d1, 0755)
	return os.TempDir() + "\\" + namefile
}

func ldap_grub_users() (map[string][]string, bool) {
	var m = map[string][]string{}
	var autoriz = true
	var ldapServer = "10.4.122.6"
	var ldapPort = uint16(389)
	var baseDN = "dc=mrg022,dc=mrg"
	var filter = []string{
		"(&(cn=*)(mail=*arg.nrg.org.ru*)(userPrincipalName=*mrg022*)(userAccountControl=512))"}
	var attributes = []string{
		"homeDirectory",
		"mail",
		"userPrincipalName"}
	l, err_ := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	bind_ := l.Bind("ldap", "ldapldap")
	if bind_ != nil {
		autoriz = false
	}
	if err_ != nil {
		println(err_.Error())
		//return err_.Error()
	}
	defer l.Close()
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		filter[0],
		attributes,
		nil)

	sr, err__ := l.Search(searchRequest)
	if err__ != nil {
		println(err__.Error())
	} else {
		for i := 0; i < len(sr.Entries); i++ {
			if Replace(sr.Entries[i].GetAttributeValue("homeDirectory"), " ", "", 100) != "" {
				var mas = []string{}
				mas = append(mas, sr.Entries[i].GetAttributeValue("homeDirectory"))
				mas = append(mas, sr.Entries[i].GetAttributeValue("mail"))

				m[Replace(sr.Entries[i].GetAttributeValue("userPrincipalName"), "@mrg022.mrg", "", 1)] = mas
			}
		}
	}

	return m, autoriz
}

func mail(rec_list []string, subject_ string, htmlbody string) {
	to_name := ""
	t_ := ""
	for i := 1; i <= len(rec_list); i++ {
		t_ += rec_list[i-1]
	}
	from := "barnaul@arg.nrg.org.ru"
	from_name := "AutoInformSystem"
	part1 := fmt.Sprintf("From: %s<%s>\r\nTo: %s <%s>\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed", from_name, from, to_name, t_, subject_)
	part2 := fmt.Sprintf("\r\nContent-Type: text/html\r\n\r\n%s\r\n", htmlbody)
	err := smtp.SendMail("H022-SRV-12.mrg022.mrg:25", nil, from, rec_list, []byte(part1+part2))
	if err != nil {
		log.Fatal(err)
	}

}

func mail_attach(rec_list []string, subject_ string, htmlbody string, file_location string, file_name string) {
	var buf bytes.Buffer
	to_name := ""
	t_ := ""
	for i := 1; i <= len(rec_list); i++ {
		t_ += rec_list[i-1] + "; "
	}
	from := "barnaul@arg.nrg.org.ru"
	from_name := "AutoInformSystem"
	marker := "{[!@!-!@!-!@!-!@!-!@!]}"
	part1 := fmt.Sprintf("From: %s<%s>\r\nTo: %s <%s>\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n--%s", from_name, from, to_name, t_, subject_, marker, marker)
	part2 := fmt.Sprintf("\r\nContent-Type: text/html\r\nContent-Transfer-Encoding:8bit\r\n\r\n%s\r\n--%s", htmlbody, marker)
	content, _ := ioutil.ReadFile(file_location)
	encoded := base64.StdEncoding.EncodeToString(content)
	lineMaxLength := 500
	nbrLines := len(encoded) / lineMaxLength
	for i := 0; i < nbrLines; i++ {
		buf.WriteString(encoded[i*lineMaxLength:(i+1)*lineMaxLength] + "\n")
	}
	buf.WriteString(encoded[nbrLines*lineMaxLength:])
	part3 := fmt.Sprintf("\r\nContent-Type: application/csv; name=\"%s\"\r\nContent-Transfer-Encoding:base64\r\nContent-Disposition: attachment; filename=\"%s\"\r\n\r\n%s\r\n--%s--", file_location, file_name, buf.String(), marker)
	err := smtp.SendMail("H022-SRV-12.mrg022.mrg:25", nil, from, rec_list, []byte(part1+part2+part3))
	if err != nil {
		log.Fatal(err)
	}
}

func runone(nameini string, sec string) {
	nef_list := []string{".jpg", ".bmp", ".mp3", ".exe", ".iso", ".ogg", ".ape", ".jpeg", ".mp4", ".avi", ".mkv", ".3gp", ".mov"}
	block_list := []string{".dll", ".tmp", ".dat", ".alias", ".multi", ".rws", "iconset\\emoticons", "ViPNet"}
	bloci := 0
	iii := 1
	iiii := 1
	text__ := ""
	text___ := ""
	text__nef := ""
	text___nef := ""
	var newel string
	flag.Parse()
	////file_ini := read_ini(flag.Arg(0))
	file_ini := read_ini(nameini, sec)
	source_dir := Replace(file_ini["path"], "\\\\H022-SRV-02\\PUBLIC\\Personal\\", "\\\\?\\k:\\Personal\\", 50) + "\\"

	src, err_ := os.Stat(source_dir)
	if err_ != nil {
		panic(err_)
	}
	if !src.IsDir() {
		fmt.Println("Source is not a directory")
		os.Exit(1)
	}
	fileList := []string{}
	err := filepath.Walk(source_dir, func(path string, f os.FileInfo, err error) error {
		for i := range fileList {
			newel = strconv.FormatInt(f.Size(), 10) + "_" + Replace(Replace(f.ModTime().String(), " ", "", 10), "+0600RTZ", "", 1)
			if newel == Split(fileList[i], "@@@")[1] {
				src_0, _ := os.Stat(path)
				src_1, _ := os.Stat(Split(fileList[i], "@@@")[0])
				if (!src_0.IsDir()) && (!src_1.IsDir()) {
					for block_inc := range block_list {
						if Contains(path, block_list[block_inc]) {
							bloci++
						}
					}
					if bloci == 0 {
						text__ += strconv.Itoa(iii) + ") \"" + path + "\"   =   \"" + Split(fileList[i], "@@@")[0] + "\" \r\n"
						fn0 := Split(path, "\\")
						fnl := len(fn0) - 1
						fn := fn0[fnl]
						dirn := Split(path, fn)[0]

						fn00 := Split(fileList[i], "@@@")[0]
						fn0_ := Split(fn00, "\\")
						fnl_ := len(fn0_) - 1
						fn_ := fn0_[fnl_]
						dirn_ := Split(fn00, fn_)[0]

						text___ += "<b>" + strconv.Itoa(iii) + ")</b> <a href='" + path + "'>" + fn + "</a> <a style='color:green' href='" + dirn + "'>*</a>  <b style='background-color:#CED8F6'>И</b>  <a href='" + Split(fileList[i], "@@@")[0] + "'>" + fn_ + "</a> <a style='color:green' href='" + dirn_ + "'>*</a><br>"
						iii++
					}
				}
			}
		}

		for nef_inc := range nef_list {
			if Contains(path, nef_list[nef_inc]) {
				text__nef += strconv.Itoa(iiii) + ") \"" + path + "\" \r\n"

				text___nef += "<b>" + strconv.Itoa(iiii) + ")</b> <a href='" + path + "'>" + Split(Replace(path, nef_list[nef_inc], "<b style='color:red'>"+nef_list[nef_inc]+"</b>", 1), "\\")[len(Split(path, "\\"))-1] + "</a> <a style='color:green' href='" + Replace(path, Split(path, "\\")[len(Split(path, "\\"))-1], "", 1) + "'>*</a><br>"
				iiii++
			}
		}
		fileList = append(fileList, path+"@@@"+strconv.FormatInt(f.Size(), 10)+"_"+Replace(Replace(f.ModTime().String(), " ", "", 10), "+0600RTZ", "", 10))
		return nil
	})
	print(err)
	d1 := []byte("Дубли: \r\n \r\n" + text__ + "файлы, которые возможно стоит убрать: \r\n \r\n" + text__nef)
	ioutil.WriteFile(os.TempDir()+"\\"+"dubl.txt", d1, 0755)

	var body, file_location, file_name, subject string

	subject = "Сводка по дисковому хранилищу"
	body = "<h3>Здравствуйте!</h3><br><b style='color:red'>" + file_ini["greetings"] + "</b><br><br>Согласно регламента организации и для ежедневной архивации данных на диске <b style='color:red'>K</b> и <b style='color:red'>Личной папке</b> не должны храниться личные данные(видео, фотографии, исполняемые файлы формата *.exe и др. личные файлы), поэтому убедительная просьба происпектировать доступные Вам папки во избежание автоматического удаления.<br>Также просим проинспектировать представленные ниже(или во вложении) файлы, которые возможно являются множественными копиями. Если файлы являются одинаковыми, просьба устранить дубликаты.<br><br><h3>Возможные дубликаты:</h3>" + text___ + "<br><h3>Файлы, которые возможно стоит убрать:</h3> <br>" + text___nef
	file_name = "dubl.txt"
	file_location = os.TempDir() + "\\" + file_name
	if iii > 1 || iiii > 1 {
		mail_attach([]string{file_ini["post"]}, subject, body, file_location, file_name)
	}
	//mail(list_post, subject, body)
	os.Remove(file_location)
}

func main() {
	dir, _ := os.Getwd()
	_, err := os.Stat(dir + "\\" + "conf.ini")
	text_ini := ""
	text_, _ := ldap_grub_users()
	if err != nil {

		for nnn, p := range text_ {
			text_ini += `
;` + nnn + `
[` + nnn + `]
post=iurii@arg.nrg.org.ru
path=` + p[0] + `
greetings="* для пользователя с логином ` + nnn + `"
`
		}
		d1 := []byte(text_ini)
		ioutil.WriteFile(dir+"\\"+"conf.ini", d1, 0755)
		const delay = 2 * Second
		Sleep(delay)
	}
	for nnnn, _ := range text_ {
		println(nnnn)
		go runone("conf.ini", nnnn)
	}

}
