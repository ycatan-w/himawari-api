package db

import (
	"time"
)

type Event struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Start       int    `json:"start"`
	End         int    `json:"end"`
}

func FindEventsByUserAndDate(userId int, date string) ([]Event, error) {
	result := []Event{}
	rows, err := DB.Query(`
	 SELECT id, title, description, date, start, end
	 FROM events
	 WHERE user_id = ? AND date = ?
	 ORDER BY start ASC
	`, userId, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.Id, &e.Title, &e.Description, &e.Date, &e.Start, &e.End); err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func FindOneEventById(eventId int) *Event {
	var e Event

	err := DB.QueryRow(`
	SELECT id, title, description, date, start, end
	FROM events
	WHERE id = ?`, eventId).Scan(&e.Id, &e.Title, &e.Description, &e.Date, &e.Start, &e.End)

	if err != nil {
		return nil
	}

	return &e
}

func FindOneEventByIdAndUserId(eventId int, userId int) *Event {
	var e Event

	err := DB.QueryRow(`
	SELECT id, title, description, date, start, end
	FROM events
	WHERE id = ? and user_id = ?`, eventId, userId).Scan(&e.Id, &e.Title, &e.Description, &e.Date, &e.Start, &e.End)
	if err != nil {
		return nil
	}

	return &e
}

func AddEvent(event Event, userId int) (int, error) {
	result, err := DB.Exec(`
	INSERT INTO events (title, description, date, start, end, user_id, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`, event.Title, event.Description, event.Date, event.Start, event.End, userId, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return -1, err
	}
	id, _ := result.LastInsertId()

	return int(id), nil
}

func UpdateEvent(event Event, eventId int, userId int) (int, error) {
	_, err := DB.Exec(`
	UPDATE events
	SET title = ?, description = ?, date = ?, start = ?, end = ?
	WHERE id = ? and user_id = ?
	`, event.Title, event.Description, event.Date, event.Start, event.End, eventId, userId)
	if err != nil {
		return -1, err
	}

	return eventId, nil
}

func RemoveEvent(eventId int, userId int) (int64, error) {
	res, err := DB.Exec(`
		DELETE FROM events
		WHERE id = ? and user_id = ?
	`, eventId, userId)

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
