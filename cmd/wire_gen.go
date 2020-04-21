// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package cmd

import (
	"github.com/clevergo/demo/internal/api"
	"github.com/clevergo/demo/internal/api/authz"
	"github.com/clevergo/demo/internal/api/post"
	"github.com/clevergo/demo/internal/api/user"
	"github.com/clevergo/demo/internal/controllers/captcha"
	"github.com/clevergo/demo/internal/core"
	"github.com/clevergo/demo/internal/frontend"
	"github.com/clevergo/demo/internal/frontend/controllers"
	"github.com/clevergo/demo/pkg/access"
	"github.com/google/wire"
)

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Injectors from wire.go:

func initializeServer() (*core.Server, func(), error) {
	logConfig := provideLogConfig()
	logger, cleanup, err := core.NewLogger(logConfig)
	if err != nil {
		return nil, nil, err
	}
	params := provideParams()
	dbConfig := provideDBConfig()
	db, cleanup2, err := core.NewDB(dbConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	sessionConfig := provideSessionConfig()
	redisConfig := provideRedisConfig()
	store := core.NewSessionStore(redisConfig)
	sessionManager := core.NewSessionManager(sessionConfig, store)
	jwtConfig := provideJWTConfig()
	jwtManager := core.NewJWTManager(jwtConfig)
	identityStore := core.NewIdentityStore(db, jwtManager)
	manager := core.NewUserManager(identityStore, sessionManager)
	mailerConfig := provideMailerConfig()
	dialer := core.NewMailer(mailerConfig)
	enforcer, err := core.NewEnforcer(dbConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	accessManager := access.New(enforcer, manager)
	application := frontend.New(logger, params, db, sessionManager, manager, dialer, accessManager)
	captchaConfig := provideCaptchaConfig()
	captchasStore := core.NewCaptchaStore(redisConfig)
	captchasManager := core.NewCaptchaManager(captchaConfig, captchasStore)
	csrfConfig := provideCSRFConfig()
	csrfMiddleware := core.NewCSRFMiddleware(csrfConfig)
	i18NConfig := provideI18NConfig()
	i18nStore := core.NewFileStore(i18NConfig)
	translators, err := core.NewI18N(i18NConfig, i18nStore)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	v := core.NewI18NLanguageParsers(i18NConfig)
	i18NMiddleware := core.NewI18NMiddleware(translators, v)
	gzipMiddleware := core.NewGzipMiddleware()
	sessionMiddleware := core.NewSessionMiddleware(sessionManager)
	m := core.NewMinify()
	minifyMiddleware := core.NewMinifyMiddleware(m)
	loggingMiddleware := core.NewLoggingMiddleware()
	viewConfig := provideViewConfig()
	renderer := core.NewRenderer(viewConfig)
	captchaCaptcha := captcha.New(captchasManager)
	site := controllers.NewSite(application, captchasManager)
	user := controllers.NewUser(application, captchasManager)
	router := provideRouter(application, captchasManager, manager, csrfMiddleware, i18NMiddleware, gzipMiddleware, sessionMiddleware, minifyMiddleware, loggingMiddleware, renderer, captchaCaptcha, site, user)
	server := provideServer(router, logger)
	return server, func() {
		cleanup2()
		cleanup()
	}, nil
}

func initializeAPIServer() (*core.Server, func(), error) {
	logConfig := provideLogConfig()
	logger, cleanup, err := core.NewLogger(logConfig)
	if err != nil {
		return nil, nil, err
	}
	params := provideParams()
	dbConfig := provideDBConfig()
	db, cleanup2, err := core.NewDB(dbConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	sessionConfig := provideSessionConfig()
	redisConfig := provideRedisConfig()
	store := core.NewSessionStore(redisConfig)
	sessionManager := core.NewSessionManager(sessionConfig, store)
	jwtConfig := provideJWTConfig()
	jwtManager := core.NewJWTManager(jwtConfig)
	identityStore := core.NewIdentityStore(db, jwtManager)
	userManager := api.NewUserManager(identityStore)
	mailerConfig := provideMailerConfig()
	dialer := core.NewMailer(mailerConfig)
	enforcer, err := core.NewEnforcer(dbConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	manager := core.NewUserManager(identityStore, sessionManager)
	accessManager := access.New(enforcer, manager)
	application := api.New(logger, params, db, sessionManager, userManager, dialer, accessManager)
	captchaConfig := provideCaptchaConfig()
	captchasStore := core.NewCaptchaStore(redisConfig)
	captchasManager := core.NewCaptchaManager(captchaConfig, captchasStore)
	authenticator := core.NewAuthenticator(identityStore)
	corsConfig := provideCORSConfig()
	corsMiddleware := core.NewCORSMiddleware(corsConfig)
	authzMiddleware := api.NewAuthzMiddleware(enforcer, userManager)
	resource := user.New(application, captchasManager, jwtManager, enforcer)
	postResource := post.New(application)
	authzResource := authz.New(application, enforcer)
	captchaCaptcha := captcha.New(captchasManager)
	server := provideAPIServer(logger, application, captchasManager, jwtManager, userManager, authenticator, corsMiddleware, authzMiddleware, resource, postResource, authzResource, captchaCaptcha)
	return server, func() {
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var superSet = wire.NewSet(
	configSet, core.NewDB, core.NewSessionStore, core.NewSessionManager, core.NewMailer, core.NewLogger, core.NewAuthenticator, core.NewIdentityStore, core.NewUserManager, core.NewCaptchaStore, core.NewCaptchaManager, core.NewI18N, core.NewFileStore, core.NewI18NLanguageParsers, provideRouter, core.NewEnforcer, access.New, core.NewJWTManager, core.MiddlewareSet, core.NewRenderer, captcha.New,
)
