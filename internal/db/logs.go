package db

import "time"

type Log struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
	Date string `json:"date"`
}

func FindLogsByUserAndDate(userId int, date string) ([]Log, error) {
	result := []Log{}
	rows, err := DB.Query(`
	 SELECT id, text, date
	 FROM logs
	 WHERE user_id = ? AND date = ?
	 ORDER BY created_at ASC
	`, userId, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var r Log
		if err := rows.Scan(&r.Id, &r.Text, &r.Date); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func FindOneLogById(id int) *Log {
	var r Log

	err := DB.QueryRow(`
	SELECT id, text, date
	FROM logs
	WHERE id = ?
	`, id).Scan(&r.Id, &r.Text, &r.Date)
	if err != nil {
		return nil
	}

	return &r
}

func FindOneLogByIdAndUserId(id int, userId int) *Log {
	var r Log

	err := DB.QueryRow(`
	SELECT id, text, date
	FROM logs
	WHERE id = ? and user_id = ?
	`, id, userId).Scan(&r.Id, &r.Text, &r.Date)
	if err != nil {
		return nil
	}

	return &r
}

func AddLog(log Log, userId int) (int, error) {
	result, err := DB.Exec(`
	INSERT INTO logs (text, date, user_id, created_at)
	VALUES (?, ?, ?, ?)
	`, log.Text, log.Date, userId, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return -1, err
	}
	id, _ := result.LastInsertId()

	return int(id), nil
}

func UpdateLog(log Log, logId int, userId int) (int, error) {
	_, err := DB.Exec(`
	UPDATE logs
	SET text = ?, date = ?
	WHERE id = ? and user_id = ?
	`, log.Text, log.Date, logId, userId)
	if err != nil {
		return -1, err
	}

	return logId, nil
}

func RemoveLog(logId int, userId int) (int64, error) {
	res, err := DB.Exec(`
		DELETE FROM logs
		WHERE id = ? and user_id = ?
	`, logId, userId)

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
