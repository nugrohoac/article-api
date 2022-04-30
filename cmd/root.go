package cmd

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/nugrohoac/kumparan-assessment/article"
	"github.com/nugrohoac/kumparan-assessment/internal/cache"
	"github.com/nugrohoac/kumparan-assessment/internal/postgresql"
	"log"

	ka "github.com/nugrohoac/kumparan-assessment"

	"github.com/spf13/viper"
)

var (
	dsnPostgres string

	// redis
	redisHost     string
	redisPort     int
	redisPassword string
	redisDB       int

	// PortApp apps
	PortApp = 7070

	configPath *string

	articleRepo ka.ArticleRepository
	// ArticleService .
	ArticleService ka.ArticleService
)

func init() {
	initEnv()
	initApp()
}

func initEnv() {
	configPath = flag.String("config-path", ".", "config path")

	flag.Parse()

	viper.AutomaticEnv()
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(*configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	timezone := "UTC"
	if tz := viper.GetString("timezone"); tz != "" {
		timezone = tz
	}

	dsnPostgres = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s TimeZone=%s sslmode=%s search_path=%s",
		viper.GetString("psql_host"),
		viper.GetString("psql_port"),
		viper.GetString("psql_user"),
		viper.GetString("psql_database"),
		viper.GetString("psql_password"),
		timezone,
		viper.GetString("psql_ssl_mode"),
		viper.GetString("psql_schema"),
	)

	if _portApp := viper.GetInt("port_app"); _portApp != 0 {
		PortApp = _portApp
	}

	redisHost = viper.GetString("redis_host")
	redisPort = viper.GetInt("redis_port")
	redisPassword = viper.GetString("redis_password")
	redisDB = viper.GetInt("redis_db")
}

func initApp() {
	db, err := sql.Open("postgres", dsnPostgres)
	if err != nil {
		log.Fatalln("Error init database connection : ", err)
	}

	addressRedis := fmt.Sprintf("%s:%d", redisHost, redisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     addressRedis,
		Password: redisPassword,
		DB:       redisDB,
	})

	articleRepo = postgresql.NewArticleRepository(db)
	articleRepo = cache.NewArticleRedis(client, articleRepo)

	ArticleService = article.NewArticleService(articleRepo)
}
