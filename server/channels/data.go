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

type ConnectionRequest struct {
	DiskId int64 `json:"diskId"`
	VmId   int64 `json:"vmId"`
}

func ParseDiscListString(str string) []int64 {
	var result []int64
	var lst = strings.Split(str, ",")

	for _, element := range lst {
		var num, err = strconv.ParseInt(strings.Trim(element, " "), 10, 64)
		if err == nil && num >= 0 {
			result = append(result, num)
		} else {
			fmt.Println("Warning: disk id is invalid (id won't be used):", element, err)
		}
	}

	return result
}

func removeDuplicateValues(intSlice []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func ContainsInSlice(s []int64, elem int64) bool {
	for _, v := range s {
		if v == elem {
			return true
		}
	}
	return false
}

func (s *VMStorage) LoadTotalDiscSpace(discs []int64) int64 {
	var result int64
	for i, elem := range discs {
		var request = "SELECT disk_space FROM discs WHERE id=" + strconv.FormatInt(elem, 10) + ";"
		rows, err := s.Db.Query(request)
		if err == nil {
			defer rows.Close()
			var discSpace int64
			if rows.Next() {
				if err := rows.Scan(&discSpace); err != nil {
					fmt.Println("Warning: error while scaning sql responce: ", err, "( id: ", elem, ")")
				}
				result += discSpace
			}
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

type VMNotFoundError struct{}

func (e *VMNotFoundError) Error() string {
	return "VM not found"
}

func (s *VMStorage) ConnectDisk(arg *ConnectionRequest) error {
	rows, err := s.Db.Query("SELECT connected_discs FROM virtual_machines WHERE id=$1", arg.VmId)
	if err != nil {
		return err
	}
	var cd string
	if rows.Next() {
		err = rows.Scan(&cd)
		if err != nil {
			return err
		}
	} else {
		return &VMNotFoundError{}
	}
	var parsed_cd = removeDuplicateValues(ParseDiscListString(cd))
	if ContainsInSlice(parsed_cd, arg.DiskId) {
		return nil
	}

	parsed_cd = append(parsed_cd, arg.DiskId)
	var new_data string
	for i, elem := range parsed_cd {
		new_data += strconv.FormatInt(elem, 10)
		if i < len(parsed_cd)-1 {
			new_data += ", "
		}
	}
	var request = "UPDATE virtual_machines SET connected_discs = '" + new_data + "' WHERE id=" + strconv.FormatInt(arg.VmId, 10)
	_, err = s.Db.Exec(request)

	return err
}
