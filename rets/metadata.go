package rets

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"time"

	"github.com/jpfielding/gorets/cmds/common"
	"github.com/jpfielding/gorets/rets"
)

//GetMetadata Get the metadata of a server is, necessary to identify the resources, tables and classes
func GetMetadata(config common.Config) ([]byte, error) {
	metadataOpts := MetadataOptions{
		ID:     "*",
		Format: "COMPACT",
		MType:  "METADATA-SYSTEM",
	}

	// should we throw an err here too?
	session, err := config.Initialize()
	if err != nil {
		return make([]byte, 0), err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	capability, err := rets.Login(ctx, session, rets.LoginRequest{URL: config.URL})
	if err != nil {
		return make([]byte, 0), err
	}
	defer rets.Logout(ctx, session, rets.LogoutRequest{URL: capability.Logout})

	reader, err := rets.MetadataStream(ctx, session, rets.MetadataRequest{
		URL:    capability.GetMetadata,
		Format: metadataOpts.Format,
		MType:  metadataOpts.MType,
		ID:     metadataOpts.ID,
	})

	defer reader.Close()
	if err != nil {
		return make([]byte, 0), err
	}
	return ioutil.ReadAll(reader)
}

// MetadataOptions ...
type MetadataOptions struct {
	MType  string `json:"metadata-type"`
	Format string `json:"format"`
	ID     string `json:"id"`
}

// SetFlags Sets the default values to perform the data collection
func (o *MetadataOptions) SetFlags() {
	flag.StringVar(&o.MType, "mtype", "METADATA-SYSTEM", "The type of metadata requested")
	flag.StringVar(&o.Format, "format", "COMPACT", "Metadata format")
	flag.StringVar(&o.ID, "id", "*", "Metadata identifier")
}

// LoadFrom Load the metadata from a file
func (o *MetadataOptions) LoadFrom(filename string) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	blob, err := ioutil.ReadAll(file)
	err = json.Unmarshal(blob, o)
	if err != nil {
		return err
	}
	return nil
}
