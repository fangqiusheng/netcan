package proc

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type Mountinfo struct {
	MountId        string
	ParentId       string
	MajorMinor     string
	Root           string
	MountPoint     string
	MountOptions   string
	OptionalFields string
	FilesystemType string
	MountSource    string
	SuperOptions   string
}

func getMountPart(pieces []string, index int) string {
	if len(pieces) > index {
		return pieces[index]
	}
	return ""
}

func GetMountInfo(fd string) ([]Mountinfo, error) {
	file, err := os.Open(fd)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	return ParseMountInfo(file)
}

func ParseMountInfoString(tx string) *Mountinfo {
	pieces := strings.Split(tx, " ")
	count := len(pieces)
	if count < 1 {
		return nil
	}
	i := strings.Index(tx, " - ")
	postFields := strings.Fields(tx[i+3:])
	preFields := strings.Fields(tx[:i])
	return &Mountinfo{
		MountId:        getMountPart(preFields, 0),
		ParentId:       getMountPart(preFields, 1),
		MajorMinor:     getMountPart(preFields, 2),
		Root:           getMountPart(preFields, 3),
		MountPoint:     getMountPart(preFields, 4),
		MountOptions:   getMountPart(preFields, 5),
		OptionalFields: getMountPart(preFields, 6),
		FilesystemType: getMountPart(postFields, 0),
		MountSource:    getMountPart(postFields, 1),
		SuperOptions:   getMountPart(postFields, 2),
	}
}

func ParseMountInfo(buffer io.Reader) ([]Mountinfo, error) {
	info := []Mountinfo{}
	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		tx := scanner.Text()
		info = append(info, *ParseMountInfoString(tx))
	}

	if err := scanner.Err(); err != nil {
		return info, err
	}
	return info, nil
}
