package channels

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type VirtualMachine struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	CpuCount       int64  `json:"cpuCount"`
	TotalDiskSpace int64  `json:"totalDiskSpace"`
}

type VMStorage struct {
	Db *sql.DB
}

func NewStore(db *sql.DB) *VMStorage {
	return &VMStorage{Db: db}
}

func ParseDiscListString(str string) []int64 {
	var result []int64
	var lst = strings.Split(str, ",")

	for _, element := range lst {
		var num, err = strconv.ParseInt(element, 10, 64)
		if err != nil && num >= 0 {
			result = append(result, num)
		} else {
			fmt.Println("Warning: disk id is invalid (id won't be used):", element)
		}
	}

	return result
}

func (s *VMStorage) LoadTotalDiscSpace(discs []int64) int64 {
	var result int64
	for i, elem := range discs {
		rows, err := s.Db.Query("SELECT id, name, disk_space FROM discs WHERE id='" + strconv.FormatInt(elem, 10) + "'")
		if err == nil {
			var discSpace int64
			if err := rows.Scan(nil, nil, discSpace); err != nil {
				fmt.Println("Warning: error while scaning sql responce: ", err)
			}
			result += discSpace
		} else {
			fmt.Println("Warning: error while loading disc info (i:", i, ", id:", elem, ", err:", err)
		}
	}
	return result
}

func (s *VMStorage) ListVirtualMachines() ([]*VirtualMachine, error) {
	rows, err := s.Db.Query("SELECT id, name, cpu_count, connected_discs FROM virtual_machines LIMIT 200")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []*VirtualMachine
	for rows.Next() {
		var c VirtualMachine
		var connected_disks string
		if err := rows.Scan(&c.Id, &c.Name, &c.CpuCount, &connected_disks); err != nil {
			return nil, err
		}
		c.TotalDiskSpace = s.LoadTotalDiscSpace(ParseDiscListString(connected_disks))

		res = append(res, &c)
	}
	if res == nil {
		res = make([]*VirtualMachine, 0)
	}
	return res, nil
}

func (s *Store) CreateChannel(name string) error {
	if len(name) < 0 {
		return fmt.Errorf("channel name is not provided")
	}
	_, err := s.Db.Exec("INSERT INTO channels (name) VALUES ($1)", name)
	return err
}
