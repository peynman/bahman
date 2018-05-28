package database_test

import (
	"testing"
	application "github.com/peyman-abdi/avalanche/app/core/app"
	"github.com/peyman-abdi/avalanche/app/core/logger"
	"github.com/peyman-abdi/testil"
	"github.com/peyman-abdi/avalanche/app/core/config"
	"github.com/peyman-abdi/avalanche/app/interfaces"
	"os"
	"time"
	"github.com/peyman-abdi/avalanche/app/core/database"
)

var app interfaces.Application
var conf interfaces.Config
var log interfaces.Logger
var repo interfaces.Repository
var mig interfaces.Migrator

var envs = map[string]string {
}
var configs = map[string]interface{} {
	"app.hjson": map[string]interface{} {
		"debug": true,
	},
	"database.hjson": map[string]interface{} {
		"app": "sqlite3",
		"runtime": map[string]interface{} {
			"migrations": "migrations",
			"connection": "sqlite3",
		},
		"connections": map[string]interface{} {
			"sqlite3": map[string]interface{} {
				"driver": "sqlite3",
				"file": "storage(\"test.db\")",
			},
		},
	},
}
func init()  {
	app = application.Initialize(0)
	os.MkdirAll(app.StoragePath(""), 0700)

	testil.CreateConfigFiles(app, configs)

	conf = config.Initialize(app)
	log = logger.Initialize(conf)
	log.LoadConsole()

	repo, mig = database.Initialize(conf, log)
}

type TestModel struct {
	ID int64
	TextValue string `gorm:"size:199;unique_index;not null;default:''"`
	NullText *string `gorm:"size:199;"`
	NullInt *int
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
func (t *TestModel) PrimaryKey() string {
	return "id"
}
func (t *TestModel) TableName() string {
	return "tests"
}
type TestModelBefore struct {
	ID int64
	NullInt *int
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
func (t *TestModelBefore) PrimaryKey() string {
	return "id"
}
func (t *TestModelBefore) TableName() string {
	return "tests"
}

type TestMigratable1 struct {
}
func (t *TestMigratable1) Up(migrator interfaces.Migrator) error {
	var err error
	if err = mig.AutoMigrate(&TestModel{}); err != nil {
		return err
	}
	return nil
}
func (t *TestMigratable1) Down(migrator interfaces.Migrator) error {
	var err error
	if err = mig.DropTableIfExists(&TestModel{}); err != nil {
		return err
	}
	return nil
}
type TestMigratable0 struct {
}
func (t *TestMigratable0) Up(migrator interfaces.Migrator) error {
	var err error
	if err = mig.CreateTable(&TestModelBefore{}); err != nil {
		return err
	}
	return nil
}
func (t *TestMigratable0) Down(migrator interfaces.Migrator) error {
	var err error
	if err = mig.DropTableIfExists(&TestModelBefore{}); err != nil {
		return err
	}
	return nil
}

// Tests
func TestMigrations(t *testing.T) {
	var err error

	if err = mig.Migrate([]interfaces.Migratable{
		new(TestMigratable0),
	}); err != nil {
		t.Error(err)
	}

	if !mig.HasTable("tests") {
		t.Errorf("Failed creating tests table!")
	}


	var migrations []*database.MigrationModel
	if err = repo.Query(&database.MigrationModel{}).Get(&migrations); err != nil {
		t.Error(err)
	}

	if len(migrations) != 1 {
		t.Errorf("Migrations in not inserted")
	}

	if err = mig.Migrate([]interfaces.Migratable{
		new(TestMigratable1),
	}); err != nil {
		t.Error(err)
	}

	if err = repo.Query(&database.MigrationModel{}).Get(&migrations); err != nil {
		t.Error(err)
	}

	if len(migrations) != 2 {
		t.Errorf("Migrations in not inserted")
	}

	if err = mig.Rollback([]interfaces.Migratable {
		new(TestMigratable1),
	}); err != nil {
		t.Error(err)
	}

	if err = repo.Query(&database.MigrationModel{}).Get(&migrations); err != nil {
		t.Error(err)
	}

	if len(migrations) != 1 {
		t.Errorf("Migrations in not rolledback")
	}
}
func TestQueries(t *testing.T) {
	var err error
	if err = mig.Migrate([]interfaces.Migratable{
		new(TestMigratable1),
	}); err != nil {
		t.Error(err)
	}

	InsertTest(t)
	QueryTest(t)
	UpdateTest(t)
	SoftDeleteTest(t)
}

// Internal test functions
func UpdateTest(t *testing.T) {
	var err error
	var tQuery = repo.Query(&TestModel{})

	var ms *TestModel
	if err = tQuery.Where("null_text is not null").GetFirst(&ms); err != nil {
		t.Error(err)
	}

	if ms == nil {
		t.Errorf("Object not found!")
		return
	}

	ms.NullInt = intRef(33)
	ms.NullText = strRef("updated field")
	repo.UpdateEntity(ms)

	queryTestNullText("null_text not null", "updated field", t)

	if err = tQuery.Where("null_text is null").Updates(map[string]interface{} {
		"null_text": "new value",
	}); err != nil {
		t.Error(err)
	}

	ms = nil
	if err = tQuery.Where("null_text = ?", "new value").GetLast(&ms); err != nil {
		t.Error(err)
	}
	if ms == nil {
		t.Errorf("Object not found!")
	}

	if err = tQuery.Update(&ms, map[string]interface{} {
		"null_int": 144,
		"null_text": "better than anything",
	}); err != nil {
		t.Error(err)
	}
	if err = tQuery.Where("null_int = ?", 144).GetLast(&ms); err != nil {
		t.Error(err)
	}
	if ms == nil {
		t.Errorf("Object not found!")
	}
}
func SoftDeleteTest(t *testing.T) {
	var err error
	var tQuery = repo.Query(&TestModel{})

	var ms *TestModel
	if err = tQuery.Where("null_text is not null").GetFirst(&ms); err != nil {
		t.Error(err)
	}
	if ms == nil {
		t.Errorf("Object not found!")
		return
	}

	if err = repo.DeleteEntity(ms); err != nil {
		t.Error(err)
	}

	var count int
	if err = tQuery.Select("count(*)").GetValue(&count); err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Errorf("Expecting remaining 1 but found: %d", count)
	}

	ms = nil
	if err = tQuery.Where("null_text is not null").GetLast(&ms); err != nil {
		t.Error(err)
	}
	if ms == nil {
		t.Errorf("Object not found!")
		return
	}

	if err = repo.DeleteEntity(ms); err != nil {
		t.Error(err)
	}
	if err = tQuery.Select("count(*)").GetValue(&count); err != nil {
		t.Error(err)
	}
	if count > 0 {
		t.Errorf("All entities are deleted but count is bigger than zero!")
	}

	//t.Logf("all deleted: %d", count)

	if err = tQuery.IncludeDeleted().Select("count(*)").GetValue(&count); err != nil {
		t.Error(err)
	}
	if count != 2 {
		t.Errorf("Deleted records are not found with soft delete: %d", count)
	}

}
func QueryTest(t *testing.T) {
	var err error
	var tQuery = repo.Query(&TestModel{})

	var ms []*TestModel
	if err = tQuery.Get(&ms); err != nil {
		t.Error(err)
	}
	if len(ms) != 2 {
		t.Error("could not query test objects")
	}

	//for _, m := range ms {
	//	t.Logf("Object: %v", m)
	//}

	QueryValuesTest(t)
}
func QueryValuesTest(t *testing.T) {
	var err error
	var tQuery = repo.Query(&TestModel{})

	var lastCreate time.Time
	if err = tQuery.Select("created_at").Limit(1).Offset(1).GetValue(&lastCreate); err != nil {
		t.Error(err)
	}
	//t.Logf("Last created_at: %v", lastCreate)

	QueryIntValuesTest(t)
	QueryStrValuesTest(t)
}
func QueryIntValuesTest(t *testing.T) {
	var err error
	var tQuery = repo.Query(&TestModel{})

	var maxId = 0
	if err = tQuery.Select("max(id)").GetValue(&maxId); err != nil {
		t.Error(err)
	}
	if maxId == 0 {
		t.Errorf("Max ID was not found: %d", maxId)
	}

	var idsList []int64
	if err = tQuery.Select("id").Order("id ASC").GetValues(&idsList); err != nil {
		t.Error(err)
	}

	if len(idsList) == 2 && idsList[0] == 1 && idsList[1] == 2 {
		//t.Logf("idsList: %v", idsList)
	} else {
		t.Errorf("Ids List error: %v", idsList)
	}

	var idsRefList []*int64
	if err = tQuery.Select("null_int").Order("id ASC").GetValues(&idsRefList); err != nil {
		t.Error(err)
	}

	if len(idsRefList) == 2 && *idsRefList[0] == 12 && *idsRefList[1] == 3 {
		//t.Logf("idsRefList: %v", idsRefList)
	} else {
		t.Errorf("Ids Ref List error: %v", idsRefList)
	}

}
func QueryStrValuesTest(t *testing.T) {
	var err error
	var tQuery = repo.Query(&TestModel{})

	var strList []string
	if err = tQuery.Select("text_value").GetValues(&strList); err != nil {
		t.Error(err)
	}

	if len(strList) == 2 && (
		(strList[0] == "test string 2" && strList[1] == "test string 1") ||
			(strList[1] == "test string 2" && strList[0] == "test string 1")) {
		//t.Logf("strList: %v", strList)
	} else {
		t.Errorf("Strings List error: %v", strList)
	}

	var strsRefList []*string
	if err = tQuery.Select("null_text").Order("id ASC").GetValues(&strsRefList); err != nil {
		t.Error(err)
	}

	if len(strsRefList) == 2 && *strsRefList[1] == "null able text" {
		//t.Logf("strRefList[1]: %v", *strsRefList[1])
	} else {
		t.Errorf("Strs Ref List error: %v", strsRefList)
	}
}
func InsertTest(t *testing.T) {
	var err error

	m1 := &TestModel{ TextValue: "test string 2", NullInt: intRef(12) }
	if err = repo.Insert(m1); err != nil {
		t.Error(err)
	}
	if m1.ID <= 0 {
		t.Error("test object was not inserted")
	}

	m2 := &TestModel{ TextValue: "test string 1", NullText: strRef("null able text"), NullInt: intRef(3)}
	if err = repo.Insert(m2); err != nil {
		t.Error(err)
	}
	if m2.ID <= 0 {
		t.Errorf("test object was not inserted: %d", m2.ID)
	}
}
// Helper functions
func queryTestNullText(where string, expected string, t *testing.T) *TestModel {
	var err error
	var tQuery = repo.Query(&TestModel{})

	var ms *TestModel
	if err = tQuery.Where(where).GetFirst(&ms); err != nil {
		t.Error(err)
	}

	if ms == nil {
		t.Errorf("Object not found!")
		return nil
	}

	if ms.NullText != nil && *ms.NullText != expected {
		t.Errorf("Null text is not as expected: %v", ms)
	}

	return ms
}
func strRef(str string) *string {
	return &str
}
func intRef(i int) *int {
	return &i
}