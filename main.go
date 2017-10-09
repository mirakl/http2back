package main

import (
	"fmt"
	"github.com/jawher/mow.cli"
	"github.com/nmaupu/http2back/provider"
	"github.com/nmaupu/http2back/server"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	appName             = "http2back"
	defaultProviderDest = "/tmp"
)

func getDefaultProvider() provider.Provider {
	return provider.Filesystem{defaultProviderDest}
}

func main() {
	app := cli.App(appName, "")
	app.Spec = "[--bind=<address>] [--port=<port>]"

	var (
		addr         = app.StringOpt("b bind", "127.0.0.1", "Bind address")
		port         = app.IntOpt("p port", 8080, "Port to listen from")
		providerFunc = getDefaultProvider
	)

	app.Action = func() {
		var (
			vAddr, vProvider, vDest, vUsername, vPassword, vBucket, vRegion, vAwsAccessKey, vAwsSecretKey, vToken string
			vPort                                                                                                 int
		)
		viper.SetConfigName(appName)
		viper.AddConfigPath(fmt.Sprintf("/etc/%s/", appName))
		viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", appName))
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			log.Printf("%s, using defaults\n", err)
		} else {
			vAddr = viper.GetString("bind_address")
			vPort = viper.GetInt("port")
			vProvider = viper.GetString("provider.name")

			if vAddr != "" && *addr == "127.0.0.1" {
				*addr = vAddr
			}
			if vPort > 0 && *port == 8080 {
				*port = vPort
			}

			switch vProvider {
			case "filesystem":
				vDest = viper.GetString("provider.dest")
				providerFunc = func() provider.Provider { return provider.Filesystem{vDest} }
			case "ftp":
				vAddr = viper.GetString("provider.host")
				vUsername = viper.GetString("provider.username")
				vPassword = viper.GetString("provider.password")
				vDest = viper.GetString("provider.dest")
				providerFunc = func() provider.Provider {
					return provider.Ftp{
						Addr:     vAddr,
						Username: vUsername,
						Password: vPassword,
						DestDir:  vDest,
					}
				}
			case "s3":
				vBucket = viper.GetString("provider.bucket")
				vDest = viper.GetString("provider.dest")
				vRegion = viper.GetString("provider.region")
				vAwsAccessKey = viper.GetString("provider.aws-access-key-id")
				vAwsSecretKey = viper.GetString("provider.aws-secret-access-key")
				vToken = viper.GetString("provider.token")
				providerFunc = func() provider.Provider {
					return provider.AwsS3{
						Bucket:             vBucket,
						DestDir:            vDest,
						Region:             vRegion,
						AwsAccessKeyId:     vAwsAccessKey,
						AwsSecretAccessKey: vAwsSecretKey,
						Token:              vToken,
					}
				}
			default:
				vProvider = "filesystem"
				vDest = "/tmp"
			}
		}

		server.Start(port, addr, providerFunc)
	}

	app.Run(os.Args)
}
