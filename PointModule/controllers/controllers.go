package controllers

import (
	"net/http"
	"strconv"
	"fmt"
	"encoding/json"
	"PointModule/models/mysql"
	"PointModule/logger"
	"time"
	"errors"
)

type response struct {
	ErrNo int `json:"errno"`
	ErrMsg string `json:"errmsg"`
	Data interface{} `json:"data"`
}

type points struct {
	UserId int `json:"user_id"`
	Total int `json:"total_points"`
}

type record struct {
	OrderId string `json:"order_id"`
	Operation int8 `json:"operation"`
	Points int `json:"points"`
	Time string `json:"time"`
}

type records struct {
	UserId int `json:"user_id"`
	Records []record `json:"records"`
}

const LEN_ORDERID  = 20

func handleErr(w http.ResponseWriter, msg string, err error) {
	var resp response

	if err != nil {
		logger.Info(err)

		resp.ErrNo = 1
		resp.ErrMsg = msg
		respB, err := json.Marshal(resp)
		if err != nil {
			logger.Error(err)
		} else {
			fmt.Fprintf(w, string(respB))
		}
	}
}

func GetTotalPoints(w http.ResponseWriter, r *http.Request) {
	var resp response
	var p points
	var err error
	var msg string

	resp.ErrNo = 0
	resp.ErrMsg = "OK"

	defer func() { handleErr(w, msg, err) }()

	p.UserId, err = strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		msg = " 用户ID错误"
		return
	}

	err = mysql.DB.QueryRow("SELECT total FROM total_points WHERE user_id=?", p.UserId).Scan(&p.Total)
	if err != nil {
		msg = "查询失败"
		return
	}

	resp.Data = p
	respB, err := json.Marshal(resp)
	if err != nil {
		msg = "系统操作失败"
		return
	} else {
		fmt.Fprintf(w, string(respB))
	}
}

func updatePoints(w http.ResponseWriter, r *http.Request, operation int8) {
	var resp response
	var msg string
	var err error

	resp.ErrNo = 0
	resp.ErrMsg = "OK"

	defer func() { handleErr(w, msg, err) }()

	userId, err := strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		msg = "用户ID错误"
		return
	}
	recordTableId := userId%100
	recordTable := "point_records_table_" + strconv.Itoa(recordTableId)

	orderId := r.FormValue("order_id")
	if orderId == "" {
		msg = "订单ID错误"
		err = errors.New("order_id is nil")
		return
	}
	if len(orderId) > LEN_ORDERID {
		msg = "订单ID长度超过" + string(LEN_ORDERID)
		err = errors.New("order_id length exceed")
		return
	}

	points, err := strconv.Atoi(r.FormValue("points"))
	if err != nil {
		msg = "积分错误"
		return
	}

	//考虑http幂等性，需要先查询该order_id对应的积分是否已更新
	var count int
	query := "SELECT COUNT(*) FROM " + recordTable + " WHERE user_id=? and order_id=?"
	err = mysql.DB.QueryRow(query, userId, orderId).Scan(&count)
	if err != nil {
		msg = "更新失败"
		return
	}

	if count == 0 {
		recordTime := time.Now().Unix()
		updateTime := time.Now().Format("2006-01-02 15:04:05")
		tx, _ := mysql.DB.Begin()
		query = "INSERT " + recordTable + " (user_id, order_id, operation, points, record_time) VALUES(?, ?, ?, ?, ?)"
		_, err = tx.Exec(query, userId, orderId, operation, points, recordTime)
		if err != nil {
			tx.Commit()
			msg = "更新失败"
			return
		}

		var query string
		if operation == 1 {
			query = "UPDATE total_points SET total=total+?, update_time=? WHERE user_id=?"
		} else if operation == 0 {
			query = "UPDATE total_points SET total=total-?, update_time=? WHERE user_id=?"
		}
		_, err = tx.Exec(query, points, updateTime, userId)
		if err != nil {
			tx.Rollback()
			msg = "更新失败"
			return
		}

		tx.Commit()
	} else {
		var p int
		var op int8
		query = "SELECT operation, points FROM " + recordTable + " WHERE user_id=? and order_id=?"
		err = mysql.DB.QueryRow(query, userId, orderId).Scan(&op, &p)
		if err != nil {
			msg = "更新失败"
			return
		}

		if (op != operation) || (p != points) {
			msg = "该订单已有积分记录"
			err = errors.New("points exception")
			return
		}
	}

	respB, err := json.Marshal(resp)
	if err != nil {
		msg = "系统操作失败"
		return
	} else {
		fmt.Fprintf(w, string(respB))
	}
}

func AddPoints(w http.ResponseWriter, r *http.Request) {
	updatePoints(w, r, 1)
}

func DeductPoints(w http.ResponseWriter, r *http.Request) {
	updatePoints(w, r, 0)
}

func GetPointsRecords(w http.ResponseWriter, r *http.Request) {
	var resp response
	var rec record
	var recs records
	var err error
	var msg string

	resp.ErrNo = 0
	resp.ErrMsg = "OK"

	defer func() { handleErr(w, msg, err) }()

	recs.UserId, err = strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		msg = "用户ID错误"
		return
	}
	recordTableId := recs.UserId%100
	recordTable := "point_records_table_" + strconv.Itoa(recordTableId)

	start := r.FormValue("start_date")
	if start == "" {
		start = "1970-01-01"
	}
	start = start + " 00:00:00"
	startT, err := time.Parse("2006-01-02 15:04:05", start)
	if err != nil {
		msg = "开始日期错误"
		return
	}
	startTime := startT.Unix()

	end := r.FormValue("end_date")
	if end == "" {
		end = time.Now().Format("2006-01-02")
	}
	end = end + " 23:59:59"
	endT, err := time.Parse("2006-01-02 15:04:05", end)
	if err != nil {
		msg = "结束日期错误"
		return
	}
	endTime := endT.Unix()

	query := "SELECT * FROM " + recordTable + " WHERE user_id=? AND record_time BETWEEN ? AND ?"
	rows, err := mysql.DB.Query(query, recs.UserId, startTime, endTime)
	if err != nil {
		msg = "查询失败"
		return
	}
	defer rows.Close()

	var id int
	var recordTime int64
	for rows.Next() {
		err = rows.Scan(&id, &recs.UserId, &rec.OrderId, &rec.Operation, &rec.Points, &recordTime)
		if err != nil {
			msg = "系统操作异常"
			return
		}
		rec.Time = time.Unix(recordTime, 0).Format("2006-01-02 15:04:05")
		recs.Records = append(recs.Records, rec)
	}
	resp.Data = recs

	respB, err := json.Marshal(resp)
	if err != nil {
		msg = "系统操作异常"
		return
	} else {
		fmt.Fprintf(w, string(respB))
	}
}
