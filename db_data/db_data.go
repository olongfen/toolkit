package db_data

import (
	"context"
	"errors"
	"github.com/olongfen/toolkit/multi/xerror"
	"github.com/olongfen/toolkit/scontext"
	"github.com/olongfen/toolkit/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"strings"
)

// DBData db data
type DBData interface {
	DB(ctx context.Context) *gorm.DB
	ExecTx(ctx context.Context, fc func(context.Context) error) error
	Close() error
}

// ITransaction tx
type ITransaction interface {
	ExecTx(ctx context.Context, fc func(ctx context.Context) error) error
}

// GormSpanKey 包内静态变量
const GormSpanKey = "__gorm_span"
const (
	CallBackBeforeName = "opentracing:before"
	CallBackAfterName  = "opentracing:after"
)

type Data struct {
	db  *gorm.DB
	log *zap.Logger
}

// contextTxKey 事务上下文 key
type contextTxKey struct{}

// NewTransaction new 事务
func NewTransaction(d DBData) ITransaction {
	return d
}

// ExecTx 执行事务
func (d *Data) ExecTx(ctx context.Context, fc func(context.Context) error) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fc(ctx)
	})
}

// DB 获取db
func (d *Data) DB(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(contextTxKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.db.WithContext(ctx)
}

// Close 关闭db连接
func (d *Data) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// NewData new database
func NewData(db *gorm.DB, logger *zap.Logger) (ret DBData, cleanup func()) {
	ret = &Data{
		db:  db,
		log: logger,
	}
	cleanup = func() {
		log.Println("db close")
		if err := ret.Close(); err != nil {
			logger.Error("db close error", zap.Error(err))
		}
	}
	return
}

// OpentracingPlugin 追踪插件
type OpentracingPlugin struct {
}

var _ gorm.Plugin = &OpentracingPlugin{}

func (op *OpentracingPlugin) Name() string {
	return "opentracingPlugin"
}

func (op *OpentracingPlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前 - 并不是都用相同的方法，可以自己自定义
	if err = db.Callback().Create().Before("gorm:before_create").Register(CallBackBeforeName, before); err != nil {
		return
	}
	if err = db.Callback().Query().Before("gorm:query").Register(CallBackBeforeName, before); err != nil {
		return
	}
	if err = db.Callback().Delete().Before("gorm:before_delete").Register(CallBackBeforeName, before); err != nil {
		return
	}
	if err = db.Callback().Update().Before("gorm:setup_reflect_value").Register(CallBackBeforeName, before); err != nil {
		return
	}
	if err = db.Callback().Row().Before("gorm:row").Register(CallBackBeforeName, before); err != nil {
		return
	}
	if err = db.Callback().Raw().Before("gorm:raw").Register(CallBackBeforeName, before); err != nil {
		return
	}

	// 结束后 - 并不是都用相同的方法，可以自己自定义
	if err = db.Callback().Create().After("gorm:after_create").Register(CallBackAfterName, after); err != nil {
		return
	}
	if err = db.Callback().Query().After("gorm:after_query").Register(CallBackAfterName, after); err != nil {
		return
	}
	if err = db.Callback().Delete().After("gorm:after_delete").Register(CallBackAfterName, after); err != nil {
		return
	}
	if err = db.Callback().Update().After("gorm:after_update").Register(CallBackAfterName, after); err != nil {
		return
	}
	if err = db.Callback().Row().After("gorm:row").Register(CallBackAfterName, after); err != nil {
		return
	}
	if err = db.Callback().Raw().After("gorm:raw").Register(CallBackAfterName, after); err != nil {
		return
	}
	return
}

func before(db *gorm.DB) {

	tr := otel.Tracer("gorm-before")

	_, span := tr.Start(db.Statement.Context, "gorm-before")
	// 利用db实例去传递span
	db.InstanceSet(GormSpanKey, span)

}

func after(db *gorm.DB) {
	if db.Error != nil {
		handlerDBError(db)
	}
	_span, exist := db.InstanceGet(GormSpanKey)
	if !exist {
		return
	}
	// 断言类型
	span, ok := _span.(trace.Span)
	if !ok {
		return
	}

	defer span.End()

	if db.Error != nil {
		span.SetAttributes(attribute.Key("gorm_err").String(db.Error.Error()))
	}

	span.SetAttributes(attribute.Key("sql").String(db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)))

}

func handlerDBError(db *gorm.DB) {
	lang := scontext.GetLanguage(db.Statement.Context)
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		db.Error = xerror.NewError(xerror.RecordNotFound, lang)
		return
	}
	if db.Statement.Schema == nil {
		return
	}
	msg := db.Error.Error()
	const (
		code23505 = "23505"
	)
	// 处理数据库错误
	for _, v := range db.Statement.Schema.DBNames {
		if strings.Contains(msg, code23505) && strings.Contains(msg, v) {
			field := db.Statement.Schema.FieldsByDBName[v]
			name := strings.ToLower(field.Name[:1]) + field.Name[1:]
			var errs = xerror.DBErrorResponse{}
			errs[name] = xerror.NewError(xerror.AlreadyExists, lang)
			db.Error = errs
		}
	}
}

// ILike ilike
type ILike clause.Eq

// Build builder sql
func (like ILike) Build(builder clause.Builder) {
	builder.WriteQuoted(like.Column)
	_, err := builder.WriteString(" ILIKE ")
	if err != nil {
		panic(err)
	}
	builder.AddVar(builder, like.Value)
}

// NegationBuild builder sql
func (like ILike) NegationBuild(builder clause.Builder) {
	builder.WriteQuoted(like.Column)
	_, err := builder.WriteString(" NOT LIKE ")
	if err != nil {
		panic(err)
	}
	builder.AddVar(builder, like.Value)
}

// ProcessDBWhere process field symbol
func ProcessDBWhere(column string, value any, symbol string) clause.Expression {
	column = utils.SnakeString(column)
	switch symbol {
	case ">":
		return clause.Gt{Column: column, Value: value}
	case ">=":
		return clause.Gte{Column: column, Value: value}
	case "<":
		return clause.Lt{Column: column, Value: value}
	case "<=":
		return clause.Lte{Column: column, Value: value}
	case "like":
		return clause.Like{Column: column, Value: value}
	case "ilike":
		return ILike{Column: column, Value: value}
	case "in":
		return clause.IN{Column: column, Values: []interface{}{value}}
	case "expr":
		return clause.Expr{
			SQL: column,
		}
	default:
		return clause.Eq{Column: column, Value: value}
	}
}
