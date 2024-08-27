package main

import (
	"fmt"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/users"
)

type dropboxModel struct {
	user users.Client
}

func auth() dropboxModel {
	token := os.Getenv("ACCESS_TOKEN")
	config := dropbox.Config{
		Token:    token,
		LogLevel: dropbox.LogDebug, // if needed, set the desired logging level. Default is off
	}
	dbx := users.New(config)

	return dropboxModel{
		user: dbx,
	}

}

func (d dropboxModel) getAccount() error {
	if resp, err := d.user.GetCurrentAccount(); err != nil {
		return err
	} else {
		fmt.Printf("Name: %v", resp.Name)
	}

	return nil
}
