package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)
type treeInfo struct {
	Id string //Путь до файла-папки и именем
	Pid string //Путь до родительской папки
	Name string // Имя файла-папки
	IsDir bool // Ключ файла или папки
	Level int // Уровень вложенности
	Size int64 // Размер файла
	IsEnd int // 1-последний . 0 - Не последний на уровне
}
func (p *treeInfo) SetIsEnd(i int) {
	p.IsEnd = i
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error{
	map_i :=0
	m := make(map[int]treeInfo)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		name := ""
		pname := ""
		levelSubDir := 0
		if info.IsDir(){
			name = fmt.Sprintf("%s", info.Name())
			pname = fmt.Sprintf("%s", path)
			levelSubDir = strings.Count(pname, "\\")
			m[map_i] = treeInfo{
				Id: pname,
				Pid: strings.Replace(pname,"\\"+name, "",1),
				Name: name,
				IsDir: true,
				Level: levelSubDir,
				Size: 0,
				IsEnd: 1, // незнаем какой и помечаем как последний
			}
			map_i++
		} else {
			if printFiles {
				name = fmt.Sprintf("%s", info.Name())
				pname = fmt.Sprintf("%s", path)
				levelSubDir = strings.Count(pname, "\\")

				m[map_i] = treeInfo{
					Id: pname,
					Pid: strings.Replace(pname,"\\"+name, "",1),
					Name: name,
					IsDir: false,
					Level: levelSubDir,
					Size: info.Size(),
					IsEnd: 1, // незнаем какой и помечаем как последний
				}
				map_i++
			}
		}
		return nil
	})
	for i:=0;i<len(m);i++{
		for j:=i+1; j<len(m);j++{
			if j != i {
				leveli := m[i].Level
				levelj := m[j].Level
				if leveli == levelj {
					if m[i].Pid == m[j].Pid{
						mTmp := m[i]
						mTmp.SetIsEnd(0) // помечаем как не последний
						m[i] = mTmp
						j = len(m)
						break
					}
				}
			}
		}
	}
	mTmp := m[0]
	prefix :=""
	name := ""
	for i:=0;i<len(m);i++ {
		mTmp = m[i]
		if mTmp.Level > 0 {
			prefix = ""
			tmp_i :=mTmp.Level-1
			for j:=i-1; j>0;j--{
				if m[j].IsDir && m[j].Level == tmp_i && tmp_i > 0{
					if m[j].IsEnd == 1{
						prefix = "\t"+ prefix
					} else {
						prefix = "│\t" + prefix
					}
					tmp_i = tmp_i -1
				}
			}
			if mTmp.IsEnd == 1{
				prefix = prefix + "└───"
			} else {
				prefix = prefix + "├───"
			}
			name = ""
			if m[i].IsDir {
				name = m[i].Name
			} else {
				if m[i].Size > 0 {
					name = m[i].Name+" ("+strconv.Itoa(int(m[i].Size))+"b)"
				} else {
					name = m[i].Name+" (empty)"
				}
			}
			fmt.Fprintf(out,"%s%s\n",prefix,name)
		}
	}
	return err
}