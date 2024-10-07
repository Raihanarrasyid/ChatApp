package app

import "ChatApp/configs"

type App struct {
	config *configs.Config
}

//	@contact.name	Raihan Arrasyid
//	@contact.email	raihanarrasyid1704@gmail.com

//	@host			configs.Host
//	@BasePath		/api/v1
//	@title			ChatApp API
//	@description	ChatApp API
//	@version		1.0

//	@securityDefinitions.apikey	UserAuthorization
//	@in							header
//	@name						Authorization
//	@description				User	Jwt Token Authorization

//	@securityDefinitions.apikey	AdminAuthorization
//	@in							header
//	@name						Authorization
//	@description				Admin	Jwt Token Authorization

// @externalDocs.description	OpenAPI
//
// @externalDocs.url			https://swagger.io/resources/open-api/
func NewApp(config *configs.Config) *App {
	return &App{config}
}
