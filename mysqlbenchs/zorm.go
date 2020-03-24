package benchs

import (
	"context"
	"fmt"

	"gitee.com/chunanyong/zorm"
)

func init() {
	st := NewSuite("zorm")
	st.InitF = func() {
		st.AddBenchmark("Insert", 2000*ORM_MULTI, 0, ZormInsert)
		st.AddBenchmark("BulkInsert 100 row", 2000*ORM_MULTI, 0, ZormInsertMulti)
		st.AddBenchmark("Update", 2000*ORM_MULTI, 0, ZormUpdate)
		st.AddBenchmark("Read", 2000*ORM_MULTI, 0, ZormRead)
		st.AddBenchmark("MultiRead limit 1000", 2000*ORM_MULTI, 1000, ZormReadSlice)
		dataSourceConfig := zorm.DataSourceConfig{
			DSN:        ORM_SOURCE,
			DriverName: "mysql",
			DBType:     "mysql",
		}
		zorm.NewBaseDao(&dataSourceConfig)
	}
}

func ZormInsert(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
	})
	b.ResetTimer()
	_, d := zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
		for i := 0; i < b.N; i++ {
			m.Id = 0
			if err := zorm.SaveStruct(ctx, m); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})

	if d != nil {
		fmt.Println(d.Error())
		b.FailNow()
	}

}

func ZormInsertMulti(b *B) {
	panic(fmt.Errorf("Don't support bulk insert"))
}

func ZormUpdate(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		_, d := zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
			return nil, zorm.SaveStruct(ctx, m)
		})
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	})
	b.ResetTimer()
	_, d := zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
		for i := 0; i < b.N; i++ {
			err := zorm.UpdateStruct(ctx, m)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	if d != nil {
		fmt.Println(d.Error())
		b.FailNow()
	}

}

func ZormRead(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		_, d := zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
			return nil, zorm.SaveStruct(ctx, m)
		})
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//查询Struct对象列表
		d := zorm.QueryStruct(context.Background(), zorm.NewSelectFinder(m.TableName()), m)
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}

func ZormReadSlice(b *B) {
	var m *Model
	wrapExecute(b, func() {
		initDB()
		m = NewModel()
		for i := 0; i < b.L; i++ {
			m.Id = 0
			_, d := zorm.Transaction(context.Background(), func(ctx context.Context) (interface{}, error) {
				return nil, zorm.SaveStruct(ctx, m)
			})
			if d != nil {
				fmt.Println(d.Error())
				b.FailNow()
			}
		}
	})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var models []Model
		d := zorm.QueryStructList(context.Background(), zorm.NewSelectFinder(m.TableName()).Append(" order by id asc "), &models, zorm.NewPage())
		if d != nil {
			fmt.Println(d.Error())
			b.FailNow()
		}
	}
}
