package main

import (
	"github.com/heroku/docker-registry-client/registry"
	"gopkg.in/urfave/cli.v2"
	"os"
	"fmt"
	"github.com/docker/distribution/manifest/schema1"
    "github.com/docker/libtrust"
    "docker-image-forwarder/config"
)

func main() {
	srcAuth := config.Auth{}
	dstAuth := config.Auth{}

	privateKey, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		fmt.Println("[ERROR] Private key generating: ", err)
	}

	authFlags := []cli.Flag {
		&cli.StringFlag{
			Name:        "src-host",
			Value:       "https://registry-1.docker.io/",
			Usage:       "Source registry hostname",
			Destination: &srcAuth.RegistryUrl,
			EnvVars: []string{"DIF_SRC_HOST"},
		},
		&cli.StringFlag{
			Name:        "dst-host",
			Value:       "https://registry-1.docker.io/",
			Usage:       "Destination registry hostname",
			Destination: &dstAuth.RegistryUrl,
			EnvVars: []string{"DIF_DST_HOST"},
		},
		&cli.StringFlag{
			Name:        "src-user",
			Value:       "",
			Usage:       "Source registry username",
			Destination: &srcAuth.Username,
			EnvVars: []string{"DIF_SRC_USER"},
		},
		&cli.StringFlag{
			Name:        "dst-user",
			Value:       "",
			Usage:       "Destination registry username",
			Destination: &dstAuth.Username,
			EnvVars: []string{"DIF_DST_USER"},
		},
		&cli.StringFlag{
			Name:        "src-pass",
			Value:       "",
			Usage:       "Source registry password",
			Destination: &srcAuth.Password,
			EnvVars: []string{"DIF_SRC_PASS"},
		},
		&cli.StringFlag{
			Name:        "dst-pass",
			Value:       "",
			Usage:       "Destination registry password",
			Destination: &dstAuth.Password,
			EnvVars: []string{"DIF_DST_PASS"},
		},
	}

	app := &cli.App{
		Commands: []*cli.Command {
			{
				Name:    "forward",
				Aliases: []string{"fwd"},
				Usage:   "Forward docker image from source registry to destination registry",
				Flags:   authFlags,
				ArgsUsage: "[source repository] [destination repository] [tag a] [tag b] ... [tag n] - if no tags given, all tags will be forwarded",
				Action:  func(c *cli.Context) error {
					srcRepo := c.Args().Get(0)
					dstRepo := c.Args().Get(1)

					if len(srcRepo) < 1 || len(dstRepo) < 1 {
						cli.ShowCommandHelpAndExit(c, "forward", 0)
						return nil
					}

					tagsToForward := c.Args().Slice()[2:]

					src, dst := getConnections(srcAuth, dstAuth)
					srcTags, _ := src.Tags(srcRepo)

					for _, tag := range srcTags {
						if len(tagsToForward) > 0 {
							// TODO: refactor
							found := false
							for _, ttf := range tagsToForward {
								if ttf == tag {
									found = true
								}
							}
							if !found {
								continue
							}
						}

						man, _ := src.Manifest(srcRepo, tag)
						for _, layer := range man.FSLayers {
							dstHasLayer, err := dst.HasLayer(dstRepo, layer.BlobSum)
							if nil != err {
								fmt.Println("[ERROR] Has layer: ", err)
							}

							if dstHasLayer {
								continue
							}

							read, err := src.DownloadLayer(srcRepo, layer.BlobSum)
							if nil != err {
								fmt.Println("[ERROR] Download layer: ", err)
							}

							err = dst.UploadLayer(dstRepo, layer.BlobSum, read)
							if nil != err {
								fmt.Println("[ERROR] Upload layer: ", err)
							}
						}

						man.Manifest.Name = dstRepo
						signedMan, err := schema1.Sign(&man.Manifest, privateKey)

						err = dst.PutManifest(dstRepo, tag, signedMan)
						if nil != err {
							fmt.Println("[ERROR] Put manifest: ", err)
						}
					}

					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			cli.ShowAppHelpAndExit(c, 0)
			return nil
		},
		HelpName: "dif",
		Name: "dif",
		Version: "0.1.0",
		Usage: "Docker Image Forwarder",
	}

	app.Run(os.Args)
}

func getConnections(srcAuth config.Auth, dstAuth config.Auth) (*registry.Registry, *registry.Registry) {
	src, err := registry.New(srcAuth.RegistryUrl, srcAuth.Username, srcAuth.Password)
	if nil != err {
		fmt.Println("[ERROR] Source registry connection", err)
	}

	dst, err := registry.New(dstAuth.RegistryUrl, dstAuth.Username, dstAuth.Password)
	if nil != err {
		fmt.Println("[ERROR] Destination registry connection", err)
	}

	return src, dst
}

