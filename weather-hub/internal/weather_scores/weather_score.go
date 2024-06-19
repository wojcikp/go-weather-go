package weatherscores

type autoInc struct {
	id int
}

func (a *autoInc) ID() (id int) {
	id = a.id
	a.id++
	return
}

var ai autoInc

type IDbClient interface {
	QueryDb(query string) (any, error)
}

type ScoreValue interface {
	string | float64
}

type IWeatherScore[T ScoreValue] interface {
	GetId() int
	GetName() string
	GetQuery() string
	GetScore(dbClient IDbClient) (T, error)
}

type BaseWeatherScore struct {
	Id int
}

func (wc *BaseWeatherScore) GetId() int {
	return ai.ID()
}
