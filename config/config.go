package config

type ServicesConfig struct {
	DynamoDb                     string /*Uri*/
	CustomerKeyService           string /*Uri*/
	SecureReleaseSettingsService string /*Uri*/
	SiteRouter                   string /*Uri*/

	//public static ServicesConfig Get(IConfiguration configuration)
	//{
	//    return configuration.FromPath("Services").ToObject<ServicesConfig>();
	//}
}

type EnvironmentConfig struct {
	EnvironmentName  string
	EnvironmentType  string
	AwsRegion        string
	IsAwsEnvironment bool
}

// TODO
//func (e *EnvironmentConfig) IsAwsEnvironment() bool {
//	return e.EnvironmentType.ToUpper().IndexOf("AWS") > 0
//}

//public static EnvironmentConfig Get(IConfiguration configuration)
//{
//    return configuration.FromPath("Environment").ToObject<EnvironmentConfig>();
//}

type Config struct {
	TestAwsAccessKeyId     string
	TestAwsSecretAccessKey string
}
