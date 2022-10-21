package service

import (
	"NeuBot/model"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// HealthUpdateService 健康上报服务
type HealthUpdateService struct {
	db *sql.DB
}

func NewHealthUpdateService() (*HealthUpdateService, error) {
	db, err := sql.Open("mysql", "xgs:xgs1150840779@tcp(101.43.161.75:3306)/health?charset=utf8&parseTime=True")
	if err != nil {
		return nil, err
	}
	return &HealthUpdateService{db: db}, nil
}

// Insert 添加健康上报
func (h *HealthUpdateService) Insert(user *model.User) error {
	tx, _ := h.db.Begin()
	defer tx.Rollback()
	exec, err := tx.Exec("insert into User(username, password,invitation,status)value (?,?,'robot',1)", user.StdNumber, user.Password)
	if err != nil {
		return err
	}
	affected, _ := exec.RowsAffected()
	if affected == 1 {
		return tx.Commit()
	}
	return nil
}

// Cancel 取消健康上报
func (h *HealthUpdateService) Cancel(user *model.User) error {
	tx, _ := h.db.Begin()
	defer tx.Rollback()
	exec, err := tx.Exec("delete from User where username= ?", user.StdNumber)
	if err != nil {
		return err
	}
	affected, _ := exec.RowsAffected()
	if affected == 1 {
		return tx.Commit()
	}
	return nil
}

// GetUser 判断数据库中是否存在用户
func (h *HealthUpdateService) GetUser(number string) (bool, error) {
	row := h.db.QueryRow("select count(*) from User where username=?", number)
	var res int
	err := row.Scan(&res)
	if err != nil {
		return false, err
	}
	if res > 0 {
		return true, nil
	}
	return false, nil
}
