package main

import (
	"fmt"
	"log"
	"path/filepath"
	"os"
	"strings"
)

var (
    Template string
)

func walkfn(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("处理路径失败 %q: %v\n", path, err)
		return err
	}
	if info.IsDir() {
		log.Printf("跳过目录 %q\n", path)
		return nil
	}

	if !strings.HasSuffix(path, "vpi") {
		log.Printf("跳过非vpi文件 %q\n", path)
		return nil
	}

	vpiFile, err := os.Open(path)
	if err != nil {
		log.Printf("打开文件失败: %s", path)
		return nil
	}
	defer vpiFile.Close()
	vpiContent := make([]byte, 4096)
	vpiFile.Read(vpiContent)
	fmt.Println(string(vpiContent))
	lastIndex := strings.Index(string(vpiContent), string("pdf"))
	fmt.Printf("lastIndex: %d", lastIndex)
	//firstIndex := lastIndex - 24
	//pdfContent := string(vpiContent)[firstIndex:lastIndex]
	//fmt.Println(string(pdfContent))
	return nil
	// Remove spaces in the string
	//oldName := filepath.Base(path)
	//newName := strings.ReplaceAll(oldName, " ", "")
	//newName = strings.ReplaceAll(newName, "_", "")
	//if strings.Contains(newName, "图纸目录") && strings.Contains(newName, "1") {
	//	newName = "001图纸目录.pdf"
	//} else if strings.Contains(newName, "图纸目录") && strings.Contains(newName, "2") {
	//	newName = "002图纸目录.pdf"
	//} else {
		// Replace "-电施-" to "DS_"
//		newName = strings.ReplaceAll(newName, "电施通-", "")
//		newName = strings.ReplaceAll(newName, "-电施-", "")
		// Replace "电施-" to "DS_" to handle situation when there
		// is only "电施-" in the string
//		newName = strings.ReplaceAll(newName, "电施-", "")
//		newName = strings.ReplaceAll(newName, "-电施", "")
		// re := regexp.MustCompile("(DS[[:digit:]]+)-")
		// newName = re.ReplaceAllString(newName, "${1}_")
	//}

	//if strings.Compare(oldName, newName) == 0 {
	//	log.Printf("skip already nice name %v\n", path)
	//	return nil
	//}
	//log.Printf("%q\n-->%q\n", oldName, newName)
	//err = os.Rename(path, filepath.Dir(path)+"/"+newName)
	//if err != nil {
	//	log.Printf("%q\n-->%q fail", oldName, newName)
	//	return err
	//}
	//return nil
}

func usage() {
	fmt.Printf("用法: %s [名字模板]\n", os.Args[0])
	fmt.Printf("\t\"名字模板\" 定义了目标文件的名字模板,模板为包含属性关键字的字符串\n")
	fmt.Printf("\t比如当名字模板为 \"图号#图名\"时:")
	fmt.Printf("对于图纸图名为\"电气抗震说明及大样图\", 图号为\"电施 - 03\"的图纸，产生文件的名字为:")
	fmt.Printf(" 电施-03#电气抗震说明及大样图.pdf\n")
	fmt.Printf("\t当前支持的属性字段:\n")
	fmt.Printf("\t\t-项目名称\n")
	fmt.Printf("\t\t-项目编号\n")
	fmt.Printf("\t\t-图名\n")
	fmt.Printf("\t\t-图号\n")
	fmt.Printf("\t注意：名字中间的空格将自动删除\n")
	fmt.Printf("\t如有额外需求，请联系姚雯雯\n")
}

func main() {
	if len(os.Args) != 2 {
		usage()
		fmt.Printf("请输入名字模板：")
	        fmt.Scan(&Template)
	} else {
		Template = os.Args[1]
	}

	f, err := os.OpenFile("whdatanalyzer.log", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("打开日志文件whdatanalyzer.log失败, %v\n", err)
		return
	}
	defer f.Close()

	log.SetOutput(f)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("获取工作路径失败 %v", err)
	}

	log.Printf("工作路径: %s, 名字模板: %s\n", dir, Template)
	err = filepath.Walk(dir, walkfn)

	if err != nil {
		log.Fatalf("处理路径失败 %s: %v\n", dir, err)
	}
}
