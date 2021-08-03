package pkg

import "io/fs"

type ProjectGenerator struct {
	fs fs.FS
	logger string
	orm string
	web string
	redis bool
}

type Options func(*ProjectGenerator)

func NewProjectGenerator(fs fs.FS, opts ...Options) *ProjectGenerator {
	pg := &ProjectGenerator{
		fs: fs,
	}
	for _, o := range opts {
		o(pg)
	}
	return pg
}

func WithLogger(logger string) Options {
	return func(pg *ProjectGenerator) {
		pg.logger = logger
	}
}

func WithORM(orm string) Options {
	return func(pg *ProjectGenerator) {
		pg.orm = orm
	}
}

func WithWeb(web string) Options {
	return func(pg *ProjectGenerator) {
		pg.web = web
	}
}

func UseRedis(redis bool) Options {
	return func(pg *ProjectGenerator) {
		pg.redis = redis
	}
}

func (pg *ProjectGenerator) Generate() {
	pg.generateCMD()
	pg.generateConfig()
	if pg.orm != "" || pg.redis {
		pg.generateDB()
	}
	if pg.web != "" {
		pg.generateAPI()
	}
}

func (pg *ProjectGenerator) generateCMD() {

}

func (pg *ProjectGenerator) generateConfig() {

}

func (pg *ProjectGenerator) generateDB() {

}

func (pg *ProjectGenerator) generateAPI() {

}

