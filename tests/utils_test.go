package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mitchellh/mapstructure"
	gotils "github.com/savsgio/gotils/strconv"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v3"
	"graphql-project/domain/model"

	"graphql-project/core"
)

type gqlQuery struct {
	Query string `json:"query"`
}

func Now() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

func getTemplate(name string, params TmplParams) ([]byte, error) {
	if _, err := os.Stat(name); err == nil {
		return parseTemplate(name, params)
	} else if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	} else {
		return nil, err
	}
}

func parseTemplate(fileName string, params TmplParams) ([]byte, error) {
	funcMap := template.FuncMap{
		"NOW": Now,
	}
	tmpl, err := template.New(path.Base(fileName)).Funcs(funcMap).ParseFiles(fileName)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	buffer.Grow(1024)
	if err := tmpl.Execute(&buffer, params()); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func loadTestData(name string, params TmplParams) (m map[string]any, err error) {
	var data []byte
	if data, err = getTemplate(name+".yml", params); err != nil {
		return
	}
	if data == nil {
		if data, err = getTemplate(name+".yaml", params); err != nil {
			return
		}
		if data == nil {
			if data, err = getTemplate(name+".json", params); err != nil {
				return
			}
			if data == nil {
				err = fmt.Errorf("testdata %w `%s.[yml|yaml|json]`", os.ErrNotExist, name)
			} else {
				err = json.Unmarshal(data, &m)
			}
		} else {
			err = yaml.Unmarshal(data, &m)
		}
	} else {
		err = yaml.Unmarshal(data, &m)
	}
	return
}

func loadRequestData(name string, params TmplParams) ([]byte, error) {
	query, err := parseTemplate(name+".gql", params)
	if err != nil {
		return nil, err
	}
	return json.Marshal(gqlQuery{gotils.B2S(query)})
}

func getTables(data map[string]any) []string {
	tables := make([]string, len(data))
	i := 0
	for t := range data {
		tables[i] = t
		i++
	}
	return tables
}

func compareDb(expectedDataFile string, params TmplParams) error {
	expected, err := loadTestData(expectedDataFile, params)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	tables := getTables(expected)
	if len(tables) == 0 {
		return errors.New("no data tables")
	}
	actual, err := getDbData(tables)
	if err != nil {
		return err
	}
	return compare(expected, actual)
}

func createApiRequest(token string) (*fasthttp.Request, error) {
	url := fasthttp.AcquireURI()
	err := url.Parse(nil, gotils.S2B(fmt.Sprintf("http://localhost:%d/graphql", Cfg.Port())))
	if err != nil {
		return nil, err
	}

	request := fasthttp.AcquireRequest()
	request.SetURI(url)

	fasthttp.ReleaseURI(url)

	request.Header.SetMethod(fasthttp.MethodPost)
	request.Header.Add("Authorization", "Bearer "+token)
	return request, err
}

func getResponseError(response *fasthttp.Response) string {
	body := response.Body()
	if len(body) > 0 {
		var m map[string]any
		if json.Unmarshal(body, &m) == nil {
			b, _ := json.MarshalIndent(m, "", "  ")
			return fmt.Sprintf("status code %s(%d)\n%s", fasthttp.StatusMessage(response.StatusCode()), response.StatusCode(), b)
		}
	}
	return fmt.Sprintf("status code %s(%d) %s", fasthttp.StatusMessage(response.StatusCode()), response.StatusCode(), body)
}

type Params = func(map[string]any) any
type TmplParams = func() any

func tmplParams(data map[string]any, params Params) TmplParams {
	if params == nil {
		return func() any {
			return nil
		}
	} else {
		return func() any {
			return params(data)
		}
	}
}

func doTestRequest(token string, testName string, params Params) (map[string]any, error) {
	if data, err := loadTestData(path.Join("testdata", testName, "db-before"), tmplParams(nil, params)); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("load db %v", err)
		}
	} else {
		if err := prepareDbData(data); err != nil {
			return nil, fmt.Errorf("prepare db %v", err)
		}
	}

	reqBody, err := loadRequestData(path.Join("testdata", testName, "request"), tmplParams(nil, params))
	if err != nil {
		return nil, fmt.Errorf("request data %v", err)
	}

	request, err := createApiRequest(token)
	if err != nil {
		return nil, fmt.Errorf("create request %v", err)
	}
	defer fasthttp.ReleaseRequest(request)

	request.SetBody(reqBody)
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	client := &fasthttp.HostClient{
		Addr: fmt.Sprintf("localhost:%d", Cfg.Port()),
	}

	err = client.Do(request, response)

	if err != nil {
		return nil, fmt.Errorf("connection error %v", err)
	}

	var entity map[string]any
	err = json.Unmarshal(response.Body(), &entity)
	if err != nil {
		return nil, fmt.Errorf("unexpected response %s", getResponseError(response))
	}

	expected, err := loadTestData(path.Join("testdata", testName, "response-"+core.IntToStr(response.StatusCode())), tmplParams(entity, params))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("load response %v", err)
		}
	}
	if expected != nil {
		if err := compare(expected, entity); err != nil {
			return nil, fmt.Errorf("unexpected response\n%v", err)
		}
	}

	if err := compareDb(path.Join("testdata", testName, "db-after"), tmplParams(entity, params)); err != nil {
		return nil, fmt.Errorf("db assert failed\n%v", err)
	}

	return entity, nil
}

func getDbData(tables []string) (map[string]any, error) {
	ctx := context.Background()
	connection, err := (*pgxpool.Pool)(DataSource).Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer connection.Release()
	query := "SELECT json_build_object("
	for i, table := range tables {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf(`'%s', (SELECT array_agg(t) FROM (SELECT * FROM "%s" ORDER BY "id") t)`, table, table)
	}
	query += ")"

	rows, err := connection.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var data map[string]any
	if rows.Next() {
		err = rows.Scan(&data)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func prepareDbData(data map[string]any) error {
	ctx := context.Background()
	connection, err := (*pgxpool.Pool)(DataSource).Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	query, err := buildTruncateQuery(data)
	if err != nil {
		return err
	}
	_, err = connection.Exec(ctx, query)
	if err != nil {
		return err
	}

	for tableName, tableData := range data {
		if tableData == nil {
			continue
		}
		rows := tableData.([]any)
		if len(rows) > 0 {
			if _, err = connection.Exec(ctx, `ALTER TABLE "`+tableName+`" DISABLE TRIGGER ALL;`); err != nil {
				return err
			}
			for _, row := range rows {
				query, args := buildInsertQuery(tableName, row.(map[string]any))
				if _, err = connection.Exec(ctx, query, args...); err != nil {
					return err
				}
			}
			if _, err = connection.Exec(ctx, `SELECT setval(pg_get_serial_sequence('`+tableName+`', 'id'), MAX(id)) FROM "`+tableName+`";`); err != nil {
				return err
			}
			if _, err = connection.Exec(ctx, `ALTER TABLE "`+tableName+`" ENABLE TRIGGER ALL;`); err != nil {
				return err
			}
		}
	}
	return nil
}

func buildTruncateQuery(data map[string]any) (string, error) {
	if len(data) == 0 {
		return "", errors.New("no data tables")
	}
	q := "TRUNCATE "
	i := 0
	for t := range data {
		if i > 0 {
			q += ", "
		}
		q += core.Quote(t)
		i++
	}
	q += " RESTART IDENTITY CASCADE;"
	return q, nil
}

func buildInsertQuery(name string, row map[string]any) (string, []any) {
	var fields, values string
	args := make([]any, len(row))
	i := 0
	for n, v := range row {
		if i > 0 {
			fields += ", "
			values += ", "
		}
		fields += core.Quote(n)
		values += "$" + core.IntToStr(i+1)
		args[i] = v
		i++
	}
	return "INSERT INTO " + core.Quote(name) + "(" + fields + ")" + " VALUES(" + values + ");", args
}

func getTokens(m map[string]any, name string) *model.Tokens {
	if d, ok := m["data"].(map[string]any); ok {
		if l, ok := d[name].(map[string]any); ok {
			var tokens model.Tokens
			_ = mapstructure.Decode(l, &tokens)
			return &tokens
		}
	}
	return nil
}
