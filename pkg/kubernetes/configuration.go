package kubernetes

import (
	"bytes"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/component-base/cli/flag"
	"k8s.io/kubectl/pkg/cmd/config"
)

func ConfigurationView() (string, error) {
	outBuffer := &bytes.Buffer{}
	pathOptions := clientcmd.NewDefaultPathOptions()
	ioStreams := genericiooptions.IOStreams{In: nil, Out: outBuffer, ErrOut: outBuffer}
	o := &config.ViewOptions{
		IOStreams:    ioStreams,
		ConfigAccess: pathOptions,
		PrintFlags:   defaultPrintFlags(),
		Flatten:      true,
		Minify:       true,
		Merge:        flag.True,
	}
	printer, err := o.PrintFlags.ToPrinter()
	if err != nil {
		return "", err
	}
	o.PrintObject = printer.PrintObj
	err = o.Run()
	if err != nil {
		return "", err
	}
	return outBuffer.String(), nil
}
