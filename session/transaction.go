//封装事务的Begin、Commit、Rollback 三个接口
package session

import "EasyORM/log"

//调用 s.db.Begin() 得到 *sql.Tx 对象，赋值给 s.tx。
func (s *Session) Begin() (err error) {
	log.Info("事务开始")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Info("事务提交")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	log.Info("事务回滚")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}