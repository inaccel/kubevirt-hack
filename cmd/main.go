package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/bombsimon/logrusr/v3"
	"github.com/inaccel/kubevirt-hack/internal"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	corev1 "k8s.io/api/core/v1"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var version string

func main() {
	app := &cli.App{
		Name:    "kubevirt-hack",
		Version: version,
		Usage:   "A self-sufficient runtime for accelerators.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "cert",
				Usage: "SSL certification file",
				Value: "/etc/inaccel/certs/ssl.pem",
			},
			&cli.StringFlag{
				Name:  "key",
				Usage: "SSL key file",
				Value: "/etc/inaccel/private/ssl.key",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "enable debug output",
			},
		},
		Before: func(context *cli.Context) error {
			log.SetOutput(io.Discard)

			logrus.SetFormatter(new(logrus.JSONFormatter))

			if context.Bool("debug") {
				logrus.SetLevel(logrus.DebugLevel)
			}

			return nil
		},
		Action: func(context *cli.Context) error {
			if err := os.MkdirAll(filepath.Join(os.TempDir(), "k8s-webhook-server", "serving-certs"), os.ModePerm); err != nil {
				return err
			}
			if err := os.Symlink(context.String("cert"), filepath.Join(os.TempDir(), "k8s-webhook-server", "serving-certs", "tls.crt")); err != nil {
				return err
			}
			if err := os.Symlink(context.String("key"), filepath.Join(os.TempDir(), "k8s-webhook-server", "serving-certs", "tls.key")); err != nil {
				return err
			}

			config, err := controllerruntime.GetConfig()
			if err != nil {
				return err
			}

			controllerruntime.SetLogger(logrusr.New(logrus.StandardLogger()))
			manager, err := controllerruntime.NewManager(config, controllerruntime.Options{
				WebhookServer: webhook.NewServer(webhook.Options{
					Port: 443,
				}),
			})
			if err != nil {
				return err
			}

			manager.GetWebhookServer().Register("/", admission.WithCustomDefaulter(manager.GetScheme(), new(corev1.Pod), internal.NewPodDefaulter()))

			return manager.GetWebhookServer().Start(context.Context)
		},
		Commands: []*cli.Command{
			initCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
