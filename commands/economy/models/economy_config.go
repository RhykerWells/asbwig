// Code generated by SQLBoiler 4.18.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// EconomyConfig is an object representing the database table.
type EconomyConfig struct {
	GuildID        string            `boil:"guild_id" json:"guild_id" toml:"guild_id" yaml:"guild_id"`
	Min            int64             `boil:"min" json:"min" toml:"min" yaml:"min"`
	Max            int64             `boil:"max" json:"max" toml:"max" yaml:"max"`
	Maxbet         int64             `boil:"maxbet" json:"maxbet" toml:"maxbet" yaml:"maxbet"`
	Symbol         string            `boil:"symbol" json:"symbol" toml:"symbol" yaml:"symbol"`
	Startbalance   int64             `boil:"startbalance" json:"startbalance" toml:"startbalance" yaml:"startbalance"`
	Workresponses  types.StringArray `boil:"workresponses" json:"workresponses,omitempty" toml:"workresponses" yaml:"workresponses,omitempty"`
	Crimeresponses types.StringArray `boil:"crimeresponses" json:"crimeresponses,omitempty" toml:"crimeresponses" yaml:"crimeresponses,omitempty"`

	R *economyConfigR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L economyConfigL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var EconomyConfigColumns = struct {
	GuildID        string
	Min            string
	Max            string
	Maxbet         string
	Symbol         string
	Startbalance   string
	Workresponses  string
	Crimeresponses string
}{
	GuildID:        "guild_id",
	Min:            "min",
	Max:            "max",
	Maxbet:         "maxbet",
	Symbol:         "symbol",
	Startbalance:   "startbalance",
	Workresponses:  "workresponses",
	Crimeresponses: "crimeresponses",
}

var EconomyConfigTableColumns = struct {
	GuildID        string
	Min            string
	Max            string
	Maxbet         string
	Symbol         string
	Startbalance   string
	Workresponses  string
	Crimeresponses string
}{
	GuildID:        "economy_config.guild_id",
	Min:            "economy_config.min",
	Max:            "economy_config.max",
	Maxbet:         "economy_config.maxbet",
	Symbol:         "economy_config.symbol",
	Startbalance:   "economy_config.startbalance",
	Workresponses:  "economy_config.workresponses",
	Crimeresponses: "economy_config.crimeresponses",
}

// Generated where

type whereHelpertypes_StringArray struct{ field string }

func (w whereHelpertypes_StringArray) EQ(x types.StringArray) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpertypes_StringArray) NEQ(x types.StringArray) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpertypes_StringArray) LT(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_StringArray) LTE(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_StringArray) GT(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_StringArray) GTE(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

func (w whereHelpertypes_StringArray) IsNull() qm.QueryMod { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpertypes_StringArray) IsNotNull() qm.QueryMod {
	return qmhelper.WhereIsNotNull(w.field)
}

var EconomyConfigWhere = struct {
	GuildID        whereHelperstring
	Min            whereHelperint64
	Max            whereHelperint64
	Maxbet         whereHelperint64
	Symbol         whereHelperstring
	Startbalance   whereHelperint64
	Workresponses  whereHelpertypes_StringArray
	Crimeresponses whereHelpertypes_StringArray
}{
	GuildID:        whereHelperstring{field: "\"economy_config\".\"guild_id\""},
	Min:            whereHelperint64{field: "\"economy_config\".\"min\""},
	Max:            whereHelperint64{field: "\"economy_config\".\"max\""},
	Maxbet:         whereHelperint64{field: "\"economy_config\".\"maxbet\""},
	Symbol:         whereHelperstring{field: "\"economy_config\".\"symbol\""},
	Startbalance:   whereHelperint64{field: "\"economy_config\".\"startbalance\""},
	Workresponses:  whereHelpertypes_StringArray{field: "\"economy_config\".\"workresponses\""},
	Crimeresponses: whereHelpertypes_StringArray{field: "\"economy_config\".\"crimeresponses\""},
}

// EconomyConfigRels is where relationship names are stored.
var EconomyConfigRels = struct {
}{}

// economyConfigR is where relationships are stored.
type economyConfigR struct {
}

// NewStruct creates a new relationship struct
func (*economyConfigR) NewStruct() *economyConfigR {
	return &economyConfigR{}
}

// economyConfigL is where Load methods for each relationship are stored.
type economyConfigL struct{}

var (
	economyConfigAllColumns            = []string{"guild_id", "min", "max", "maxbet", "symbol", "startbalance", "workresponses", "crimeresponses"}
	economyConfigColumnsWithoutDefault = []string{"guild_id"}
	economyConfigColumnsWithDefault    = []string{"min", "max", "maxbet", "symbol", "startbalance", "workresponses", "crimeresponses"}
	economyConfigPrimaryKeyColumns     = []string{"guild_id"}
	economyConfigGeneratedColumns      = []string{}
)

type (
	// EconomyConfigSlice is an alias for a slice of pointers to EconomyConfig.
	// This should almost always be used instead of []EconomyConfig.
	EconomyConfigSlice []*EconomyConfig

	economyConfigQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	economyConfigType                 = reflect.TypeOf(&EconomyConfig{})
	economyConfigMapping              = queries.MakeStructMapping(economyConfigType)
	economyConfigPrimaryKeyMapping, _ = queries.BindMapping(economyConfigType, economyConfigMapping, economyConfigPrimaryKeyColumns)
	economyConfigInsertCacheMut       sync.RWMutex
	economyConfigInsertCache          = make(map[string]insertCache)
	economyConfigUpdateCacheMut       sync.RWMutex
	economyConfigUpdateCache          = make(map[string]updateCache)
	economyConfigUpsertCacheMut       sync.RWMutex
	economyConfigUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// OneG returns a single economyConfig record from the query using the global executor.
func (q economyConfigQuery) OneG(ctx context.Context) (*EconomyConfig, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single economyConfig record from the query.
func (q economyConfigQuery) One(ctx context.Context, exec boil.ContextExecutor) (*EconomyConfig, error) {
	o := &EconomyConfig{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for economy_config")
	}

	return o, nil
}

// AllG returns all EconomyConfig records from the query using the global executor.
func (q economyConfigQuery) AllG(ctx context.Context) (EconomyConfigSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all EconomyConfig records from the query.
func (q economyConfigQuery) All(ctx context.Context, exec boil.ContextExecutor) (EconomyConfigSlice, error) {
	var o []*EconomyConfig

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to EconomyConfig slice")
	}

	return o, nil
}

// CountG returns the count of all EconomyConfig records in the query using the global executor
func (q economyConfigQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all EconomyConfig records in the query.
func (q economyConfigQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count economy_config rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table using the global executor.
func (q economyConfigQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q economyConfigQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if economy_config exists")
	}

	return count > 0, nil
}

// EconomyConfigs retrieves all the records using an executor.
func EconomyConfigs(mods ...qm.QueryMod) economyConfigQuery {
	mods = append(mods, qm.From("\"economy_config\""))
	q := NewQuery(mods...)
	if len(queries.GetSelect(q)) == 0 {
		queries.SetSelect(q, []string{"\"economy_config\".*"})
	}

	return economyConfigQuery{q}
}

// FindEconomyConfigG retrieves a single record by ID.
func FindEconomyConfigG(ctx context.Context, guildID string, selectCols ...string) (*EconomyConfig, error) {
	return FindEconomyConfig(ctx, boil.GetContextDB(), guildID, selectCols...)
}

// FindEconomyConfig retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindEconomyConfig(ctx context.Context, exec boil.ContextExecutor, guildID string, selectCols ...string) (*EconomyConfig, error) {
	economyConfigObj := &EconomyConfig{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"economy_config\" where \"guild_id\"=$1", sel,
	)

	q := queries.Raw(query, guildID)

	err := q.Bind(ctx, exec, economyConfigObj)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from economy_config")
	}

	return economyConfigObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *EconomyConfig) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *EconomyConfig) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no economy_config provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(economyConfigColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	economyConfigInsertCacheMut.RLock()
	cache, cached := economyConfigInsertCache[key]
	economyConfigInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			economyConfigAllColumns,
			economyConfigColumnsWithDefault,
			economyConfigColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(economyConfigType, economyConfigMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(economyConfigType, economyConfigMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"economy_config\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"economy_config\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into economy_config")
	}

	if !cached {
		economyConfigInsertCacheMut.Lock()
		economyConfigInsertCache[key] = cache
		economyConfigInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single EconomyConfig record using the global executor.
// See Update for more documentation.
func (o *EconomyConfig) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the EconomyConfig.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *EconomyConfig) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	economyConfigUpdateCacheMut.RLock()
	cache, cached := economyConfigUpdateCache[key]
	economyConfigUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			economyConfigAllColumns,
			economyConfigPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update economy_config, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"economy_config\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, economyConfigPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(economyConfigType, economyConfigMapping, append(wl, economyConfigPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update economy_config row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for economy_config")
	}

	if !cached {
		economyConfigUpdateCacheMut.Lock()
		economyConfigUpdateCache[key] = cache
		economyConfigUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (q economyConfigQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q economyConfigQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for economy_config")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for economy_config")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o EconomyConfigSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o EconomyConfigSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), economyConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"economy_config\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, economyConfigPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in economyConfig slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all economyConfig")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *EconomyConfig) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns, opts...)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *EconomyConfig) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns, opts ...UpsertOptionFunc) error {
	if o == nil {
		return errors.New("models: no economy_config provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(economyConfigColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	economyConfigUpsertCacheMut.RLock()
	cache, cached := economyConfigUpsertCache[key]
	economyConfigUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, _ := insertColumns.InsertColumnSet(
			economyConfigAllColumns,
			economyConfigColumnsWithDefault,
			economyConfigColumnsWithoutDefault,
			nzDefaults,
		)

		update := updateColumns.UpdateColumnSet(
			economyConfigAllColumns,
			economyConfigPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert economy_config, could not build update column list")
		}

		ret := strmangle.SetComplement(economyConfigAllColumns, strmangle.SetIntersect(insert, update))

		conflict := conflictColumns
		if len(conflict) == 0 && updateOnConflict && len(update) != 0 {
			if len(economyConfigPrimaryKeyColumns) == 0 {
				return errors.New("models: unable to upsert economy_config, could not build conflict column list")
			}

			conflict = make([]string, len(economyConfigPrimaryKeyColumns))
			copy(conflict, economyConfigPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"economy_config\"", updateOnConflict, ret, update, conflict, insert, opts...)

		cache.valueMapping, err = queries.BindMapping(economyConfigType, economyConfigMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(economyConfigType, economyConfigMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if errors.Is(err, sql.ErrNoRows) {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert economy_config")
	}

	if !cached {
		economyConfigUpsertCacheMut.Lock()
		economyConfigUpsertCache[key] = cache
		economyConfigUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteG deletes a single EconomyConfig record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *EconomyConfig) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single EconomyConfig record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *EconomyConfig) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no EconomyConfig provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), economyConfigPrimaryKeyMapping)
	sql := "DELETE FROM \"economy_config\" WHERE \"guild_id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from economy_config")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for economy_config")
	}

	return rowsAff, nil
}

func (q economyConfigQuery) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all matching rows.
func (q economyConfigQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no economyConfigQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from economy_config")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for economy_config")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o EconomyConfigSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o EconomyConfigSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), economyConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"economy_config\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, economyConfigPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from economyConfig slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for economy_config")
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *EconomyConfig) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no EconomyConfig provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *EconomyConfig) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindEconomyConfig(ctx, exec, o.GuildID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *EconomyConfigSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: empty EconomyConfigSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *EconomyConfigSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := EconomyConfigSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), economyConfigPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"economy_config\".* FROM \"economy_config\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, economyConfigPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in EconomyConfigSlice")
	}

	*o = slice

	return nil
}

// EconomyConfigExistsG checks if the EconomyConfig row exists.
func EconomyConfigExistsG(ctx context.Context, guildID string) (bool, error) {
	return EconomyConfigExists(ctx, boil.GetContextDB(), guildID)
}

// EconomyConfigExists checks if the EconomyConfig row exists.
func EconomyConfigExists(ctx context.Context, exec boil.ContextExecutor, guildID string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"economy_config\" where \"guild_id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, guildID)
	}
	row := exec.QueryRowContext(ctx, sql, guildID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if economy_config exists")
	}

	return exists, nil
}

// Exists checks if the EconomyConfig row exists.
func (o *EconomyConfig) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	return EconomyConfigExists(ctx, exec, o.GuildID)
}
