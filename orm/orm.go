package orm

/*
orm思路：
1. diy一个log库，支持info和error类型，支持设置日志等级；
2. 封装数据库基本操作session，支持log日志打印；
3. engine支持连接关闭数据库，支持开启session进行数据库交互；
4. 考虑到不同数据库sql不同，需要增加适配器dialect，这也有利于扩展和解耦；
5. 建立数据库元素和go结构体的映射：
	 	- 表名(table name) —— 结构体名(struct name)
		- 字段名和字段类型 —— 成员变量和类型。
		- 额外的约束条件(例如非空、主键等) —— 成员变量的Tag
6. 通过reflect包将任意的对象解析为Schema实例；
7. 为复杂sql语句构建生成器clause
8. 
*/
import (
	"database/sql"

	"orm/dialect"
	"orm/log"
	"orm/session"
)

type Engine struct {
	db *sql.DB

	dialect dialect.Dialect
}

// driver是数据库类型，source是数据库地址
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	// 校验是否有对应数据库的适配器
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}

	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// 支持通过engine实例创建会话
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}