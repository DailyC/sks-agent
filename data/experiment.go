package data

import (
	"time"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
)

type ExperimentModel struct {
	Uid        string
	Command    string
	SubCommand string
	Flag       string
	Status     string
	Error      string
	CreateTime string
	UpdateTime string
}

const expTableDDL = `CREATE TABLE IF NOT EXISTS experiment (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	uid VARCHAR(32) UNIQUE,
	command VARCHAR NOT NULL,
	sub_command VARCHAR,
	flag VARCHAR,
	status VARCHAR,
	error VARCHAR,
	create_time VARCHAR,
	update_time VARCHAR
)`

var expIndexDDL = []string{
	`CREATE INDEX exp_uid_uidx ON experiment (uid)`,
	`CREATE INDEX exp_command_idx ON experiment (command)`,
	`CREATE INDEX exp_status_idx ON experiment (status)`,
}

var insertExpDML = `INSERT INTO 
	experiment (uid, command, sub_command, flag, status, error, create_time, update_time)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
`

func (s *Source) checkAndInitExperimentTable() {
	exists, err := s.experimentTableExists()
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	if !exists {
		err = s.initExperimentTable()
		if err != nil {
			logrus.Fatalf(err.Error())
		}
	}
}

func (s *Source) experimentTableExists() (bool, error) {
	stmt, err := s.DB.Prepare(tableExistsDQL)
	if err != nil {
		return false, fmt.Errorf("select experiment table exists err when invoke db prepare, %s", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query("experiment")
	if err != nil {
		return false, fmt.Errorf("select experiment table exists or not err, %s", err)
	}
	defer rows.Close()
	var c int
	for rows.Next() {
		rows.Scan(&c)
		break
	}
	return c != 0, nil
}

func (s *Source) initExperimentTable() error {
	_, err := s.DB.Exec(expTableDDL)
	if err != nil {
		return fmt.Errorf("create experiment table err, %s", err)
	}
	for _, sql := range expIndexDDL {
		s.DB.Exec(sql)
	}
	return nil
}

func (s *Source) InsertExperimentModel(model *ExperimentModel) error {
	stmt, err := s.DB.Prepare(insertExpDML)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		model.Uid,
		model.Command,
		model.SubCommand,
		model.Flag,
		model.Status,
		model.Error,
		model.CreateTime,
		model.UpdateTime,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Source) UpdateExperimentModelByUid(uid, status, errMsg string) error {
	stmt, err := s.DB.Prepare(`UPDATE experiment
	SET status = ?, error = ?, update_time = ?
	WHERE uid = ?
`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(status, errMsg, time.Now().Format(time.RFC3339Nano), uid)
	if err != nil {
		return err
	}
	return nil
}

func (s *Source) QueryExperimentModelByUid(uid string) (*ExperimentModel, error) {
	stmt, err := s.DB.Prepare(`SELECT * FROM experiment WHERE uid = ?`)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	models, err := getExperimentModelsFrom(rows)
	if err != nil {
		return nil, err
	}
	if len(models) == 0 {
		return nil, nil
	}
	return models[0], nil
}

func (s *Source) ListExperimentModels() ([]*ExperimentModel, error) {
	stmt, err := s.DB.Prepare(`SELECT * FROM experiment`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getExperimentModelsFrom(rows)
}

func (s *Source) QueryExperimentModelsByCommand(target string) ([]*ExperimentModel, error) {
	stmt, err := s.DB.Prepare(`SELECT * FROM experiment where command = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(target)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return getExperimentModelsFrom(rows)
}

func getExperimentModelsFrom(rows *sql.Rows) ([]*ExperimentModel, error) {
	models := make([]*ExperimentModel, 0)
	for rows.Next() {
		var id int
		var uid, command, subCommand, flag, status, error, createTime, updateTime string
		err := rows.Scan(&id, &uid, &command, &subCommand, &flag, &status, &error, &createTime, &updateTime)
		if err != nil {
			return nil, err
		}
		model := &ExperimentModel{
			Uid:        uid,
			Command:    command,
			SubCommand: subCommand,
			Flag:       flag,
			Status:     status,
			Error:      error,
			CreateTime: createTime,
			UpdateTime: updateTime,
		}
		models = append(models, model)
	}
	return models, nil
}