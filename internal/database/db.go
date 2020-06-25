package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" //hfjfh
)

type Datastore interface {
	Close()

	GetProject(name string) *Project
	SaveProject(p *Project)

	GetFunction(name string, project uint) *Function
	SaveFunction(f *Function)

	GetCronFunctions() []*Function

	//AllProjects() ([]*Project, error)

	//CreateProject(project Project) (int64, error)
	//CreateFunction(function Function) (int64, error)
}

type DB struct {
	Con *gorm.DB
}

func NewDB(path string) (*DB, error) {
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Project{})
	db.AutoMigrate(&Function{})

	return &DB{Con: db}, nil
}

func (db *DB) Close() {
	db.Con.Close()
}

func (db *DB) GetProject(name string) *Project {
	p := &Project{}
	if db.Con.Find(p, "name = ?", name).RecordNotFound() {
		return nil
	}
	return p
}

func (db *DB) SaveProject(p *Project) {
	db.Con.Save(p)
}

func (db *DB) GetFunction(name string, project uint) *Function {
	fn := &Function{}
	if db.Con.Find(fn, "name = ? AND project_id = ?", name, project).RecordNotFound() {
		return nil
	}
	return fn
}

func (db *DB) SaveFunction(f *Function) {
	db.Con.Save(f)
}

func (db *DB) GetCronFunctions() []*Function {
	functions := []*Function{}
	db.Con.Find(&functions, "handler = ? ", CronJobHandler)
	return functions
}
