package dbprovider

import (
	"errors"
	"os"
	"testing"
	"time"

	"dbprovider/helpers"
	"github.com/google/go-cmp/cmp"

	"consts"
	"dbprovider/models"
)

func TestMain(m *testing.M) {
	err := os.Setenv(consts.DbEnvironmentKey, consts.DbTestName)
	if err != nil {
		panic(err)
	}
	InitConnection()
	os.Exit(m.Run())
}

func TestInitDB(t *testing.T) {
	err := Manager.InitDB()
	if err != nil {
		t.Fatal(err)
	}

	total := int64(0)
	Manager.getDB().Unscoped().Where("1 = 1").Find(&models.Timings{}).Count(&total)

	if total != 1 {
		if total == 0 {
			t.Error("no timings data")
		} else {
			t.Error("too many timings data")
		}
	}

	Manager.getDB().Unscoped().Where("1 = 1").Find(&models.Query{}).Count(&total)
	if total != 0 {
		t.Errorf("got %d query entities, expected 0", total)
	}

	Manager.getDB().Unscoped().Where("1 = 1").Find(&models.Task{}).Count(&total)
	if total != 0 {
		t.Errorf("got %d task entities, expected 0", total)
	}

	Manager.getDB().Unscoped().Where("1 = 1").Find(&models.Worker{}).Count(&total)
	if total != 0 {
		t.Errorf("got %d worker entities, expected 0", total)
	}
}

func TestNewQuery(t *testing.T) {
	t1 := &models.Task{
		Operation: "",
		Index:     0,
		Result:    1.0,
		IsDone:    true,
	}

	t2 := &models.Task{
		Operation: "",
		Index:     1,
		Result:    2.0,
		IsDone:    true,
	}

	t3 := &models.Task{
		Operation: "+",
		Duration:  2 * time.Second,
		Subtasks:  []*models.Task{t1, t2},
	}

	q := &models.Query{
		Expression: "1+2",
		BadMessage: "",
		Tasks:      []*models.Task{t1, t2, t3},
	}

	if err := Manager.NewQuery(q); err != nil {
		t.Fatal(err)
	}
	if q.ID == 0 {
		t.Fatal("query entity wasn't created")
	}
	if t1.ID == 0 {
		t.Fatal("first child task entity wasn't created")
	}
	if t2.ID == 0 {
		t.Fatal("second task entity wasn't created")
	}
	if t3.ID == 0 {
		t.Fatal("main task entity wasn't created")
	}
	if t1.ParentID != t3.ID {
		t.Fatal("first child task entity isn't associated with main task entity")
	}
	if t2.ParentID != t3.ID {
		t.Fatal("first child task entity isn't associated with main task entity")
	}
	if t1.TargetID != q.ID || t2.TargetID != q.ID || t3.TargetID != q.ID {
		t.Fatal("task entities isn't associated with query entity")
	}
	if !t1.IsDone || !t2.IsDone {
		t.Fatal("child tasks are a plain numbers, but not marked as done")
	}
	if !t3.IsReady {
		t.Fatal("main task entity isn't marked as ready")
	}

	{
		q.Tasks = nil // clear associations
		qx := &models.Query{}
		Manager.getDB().First(qx, q.ID)
		if !cmp.Equal(q, qx) {
			t.Error("query entities aren't equal")
			t.Error(">>>", q)
			t.Error(">>>", qx)
		}
	}

	{
		t1.Subtasks = nil // clear associations
		t1.Parent = nil   // clear associations
		t1.Target = nil   // clear associations
		tx := &models.Task{}
		Manager.getDB().First(tx, t1.ID)
		if !cmp.Equal(t1, tx) {
			t.Error("task entities aren't equal")
			t.Error(">>>", t1)
			t.Error(">>>", tx)
		}

		t2.Subtasks = nil // clear associations
		t2.Parent = nil   // clear associations
		t2.Target = nil   // clear associations
		tx = &models.Task{}
		Manager.getDB().First(tx, t2.ID)
		if !cmp.Equal(t2, tx) {
			t.Error("task entities aren't equal")
			t.Error(">>>", t2)
			t.Error(">>>", tx)
		}

		t3.Subtasks = nil // clear associations
		t3.Parent = nil   // clear associations
		t3.Target = nil   // clear associations
		tx = &models.Task{}
		Manager.getDB().First(tx, t3.ID)
		if !cmp.Equal(t3, tx) {
			t.Error("task entities aren't equal")
			t.Error(">>>", t2)
			t.Error(">>>", tx)
		}
	}
}

func TestBrokenQuery(t *testing.T) {
	q := &models.Query{}
	if err := Manager.NewQuery(q); !errors.Is(err, helpers.ErrQueryContractBothMissed) {
		t.Errorf("expected %v got %v", helpers.ErrQueryContractBothMissed, err)
	}

	q = &models.Query{
		BadMessage: "kekeke",
		Tasks:      []*models.Task{{}, {}, {}},
	}
	if err := Manager.NewQuery(q); !errors.Is(err, helpers.ErrQueryContractBothPresented) {
		t.Errorf("expected %v got %v", helpers.ErrQueryContractBothPresented, err)
	}
}

func TestUpdateQuery(t *testing.T) {
	q := &models.Query{
		BadMessage: "kekeke",
	}
	if err := Manager.NewQuery(q); err != nil {
		t.Fatal(err)
	}

	q.BadMessage = "blah-blah-blah"
	if err := Manager.UpdateQuery(q); err != nil {
		t.Fatal(err)
	}

	q.BadMessage = ""
	q.Tasks = []*models.Task{{}, {}, {}}
	if err := Manager.UpdateQuery(q); err != nil {
		t.Fatal(err)
	}

	if err := Manager.UpdateQuery(q); !errors.Is(err, ErrGoodQueryCannotBeUpdated) {
		t.Errorf("expected %v got %v", ErrGoodQueryCannotBeUpdated, err)
	}

	q.BadMessage = "became dark side"
	q.Tasks = nil
	if err := Manager.UpdateQuery(q); !errors.Is(err, ErrGoodQueryCannotBeUpdated) {
		t.Errorf("expected %v got %v", ErrGoodQueryCannotBeUpdated, err)
	}

	if err := Manager.UpdateQuery(&models.Query{}); !errors.Is(err, ErrUpdateQueryWithoutID) {
		t.Errorf("expected %v got %v", ErrUpdateQueryWithoutID, err)
	}
}
