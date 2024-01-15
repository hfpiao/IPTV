package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type DisplayNameStruct struct {
	Lang string `xml:"lang,attr"`
	Text string `xml:",innerxml"`
}

type TitleStruct struct {
	Lang string `xml:"lang,attr"`
	Text string `xml:",innerxml"`
}

type IconStruct struct {
	Src string `xml:"src,attr"`
}

type DescStruct struct {
	Lang string `xml:"lang,attr"`
	Text string `xml:",innerxml"`
}

type CategoryStruct struct {
	Lang string `xml:"lang,attr"`
	Text string `xml:",innerxml"`
}

type ProgrammeStruct struct {
	ChannelId string      `xml:"channel,attr"`
	Start     string      `xml:"start,attr"`
	Stop      string      `xml:"stop,attr"`
	Title     TitleStruct `xml:"title"`
	//Desc      DescStruct     `xml:"desc,omitempty"`
	//Category  CategoryStruct `xml:"category,omitempty"`
	//Icon      IconStruct     `xml:"icon,omitempty"`
}

type ChannelStruct struct {
	ID          string            `xml:"id,attr"`
	DisplayName DisplayNameStruct `xml:"display-name"`
	//Icon        IconStruct        `xml:"icon,omitempty"`
	//Url         string            `xml:"url,omitempty"`
}

type Tv struct {
	XMLName   xml.Name          `xml:"tv"`
	InfoName  string            `xml:"info-name,attr"`
	InfoUrl   string            `xml:"info-url,attr"`
	Channel   []ChannelStruct   `xml:"channel"`
	Programme []ProgrammeStruct `xml:"programme"`
}

func main() {
	// 读取文本文件中的ID列表
	file, err := os.Open("channel_list.log")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("关闭channel_list.log失败:", err)
			return
		}
	}(file)

	ids := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id := scanner.Text()
		ids = append(ids, id)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件时出错:", err)
		return
	}

	// 解析XML字符串为tv结构体实例
	// 读取文本文件中的ID列表
	epgData, err := ioutil.ReadFile("EPG.xml")
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	var root, tempRoot Tv
	xml.Unmarshal(epgData, &root)

	// 遍历所有channel标签，检查其ID是否在要删除的列表中，如果是，则删除该标签
	channelMap := map[string]ChannelStruct{}
	for _, channel := range root.Channel {
		channelMap[channel.ID] = channel
	}

	for _, v := range ids {
		if v, ok := channelMap[v]; ok {
			tempRoot.Channel = append(tempRoot.Channel, v)
		}

		for _, programme := range root.Programme {
			if programme.ChannelId == v {
				tempRoot.Programme = append(tempRoot.Programme, programme)
			}
		}
	}

	tempRoot.InfoName = "Modify by hfpiao"
	tempRoot.InfoUrl = "https://epg.112114.xyz"

	// 创建一个缓冲区
	buffer := &bytes.Buffer{}
	buffer.WriteString(xml.Header)
	encoder := xml.NewEncoder(buffer)
	encoder.Indent("", "   ")
	err = encoder.Encode(tempRoot)
	if err != nil {
		fmt.Println("生成xml失败:", err)
		return
	}

	// 将缓冲区的内容写入文件
	err = ioutil.WriteFile("hfpiao.xml", buffer.Bytes(), 0644)
	if err != nil {
		fmt.Println("无法保存文件:", err)
		return
	}
}
