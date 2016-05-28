package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/spf13/viper"
)

// Global Variables
var ClientPrivateKey *rsa.PrivateKey
var ClientPublicKey *rsa.PublicKey
var ClientUser, Server string

// Initialize config
func init() {
	var err error
	// Load RSA private and public key
	ClientPrivateKey, err = getPrivateKey()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	ClientPublicKey = &ClientPrivateKey.PublicKey
	// Initialize config
	viper.SetDefault("ClientUser", "test")
	viper.SetDefault("Server", "127.0.0.1:8080")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("toml")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	ClientUser = viper.GetString("ClientUser")
	Server = "http://" + viper.GetString("Server")
}

// Main function, parses cli args and runs appropriate function
func main() {
	usage := `CS3031 Lab2 Client.

Usage:
  client register <username>
  client upload <filepath> <filename>
  client download <user> <filename> <outputpath>
  client share <filename> <user>...
  client revoke <filename> <user>...
  client -h | --help

Options:
  -h --help     Show this screen.`

	args, _ := docopt.Parse(usage, nil, true, "", false)
	if args["register"].(bool) == true {
		Register(args["<username>"].(string))
	} else if args["upload"].(bool) == true {
		UploadFile(args["<filepath>"].(string), args["<filename>"].(string))
	} else if args["download"].(bool) == true {
		DownloadFile(args["<user>"].(string), args["<filename>"].(string), args["<outputpath>"].(string))
	} else if args["share"].(bool) == true {
		ShareFile(args["<filename>"].(string), args["<user>"].([]string), true)
	} else if args["revoke"].(bool) == true {
		RevokeFile(args["<filename>"].(string), args["<user>"].([]string))
	}
}

// Register user with server, will fail if username is taken
func Register(username string) {
	user := NewUser(username, ClientPublicKey)
	err := user.Register()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully registered")
	os.Exit(0)
}

// Upload a file to server
func UploadFile(filepath string, filename string) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	key, err := generateAESKey()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	encodedData, err := encryptAES(key, data)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	file := NewFile(ClientUser, filename, encodedData)
	err = file.Upload()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	encodedKey, err := encrypt(ClientPublicKey, key)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	filekey := NewFileKey(ClientUser, ClientUser, filename, encodedKey)
	err = filekey.Share()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully uploaded file")
	os.Exit(0)
}

// Download File and decrypt with shared key, output file to given path
// If user doesn't have file access the program will exit with an error message
func DownloadFile(owner string, filename string, outputPath string) {
	file, err := GetFile(owner, filename)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	filekey, err := GetFileKey(owner, filename)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	decodedKey, err := decrypt(ClientPrivateKey, filekey.Key)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	decodedData, err := decryptAES(decodedKey, file.Data)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	err = ioutil.WriteFile(outputPath, decodedData, 0644)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully downloaded file")
	os.Exit(0)
}

// Share file with given users
func ShareFile(filename string, users []string, command bool) {
	// Get shared secret key
	filekey, err := GetFileKey(ClientUser, filename)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	decodedKey, err := decrypt(ClientPrivateKey, filekey.Key)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	// Share file access with given users
	for _, username := range users {
		user, err := GetUser(username)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
		filekey.Id = ""
		filekey.User = username
		encodedKey, err := encrypt(user.PubKey, decodedKey)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
		filekey.Key = encodedKey
		err = filekey.Share()
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
	}
	// If run as terminal command exit with success message
	if command {
		fmt.Println("Successfully shared file")
		os.Exit(0)
	}
}

// Revoke file access for given users
func RevokeFile(filename string, users []string) {
	// Get existing file and key
	file, err := GetFile(ClientUser, filename)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	filekey, err := GetFileKey(ClientUser, filename)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	decodedKey, err := decrypt(ClientPrivateKey, filekey.Key)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	decodedData, err := decryptAES(decodedKey, file.Data)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	// Create new key, re-encrypt and upload file
	newKey, err := generateAESKey()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	encodedData, err := encryptAES(newKey, decodedData)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	file.Data = encodedData
	err = file.Upload()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	encodedKey, err := encrypt(ClientPublicKey, newKey)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	filekey.Key = encodedKey
	err = filekey.Share()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	// Revoke file access for given users
	for _, user := range users {
		filekey = NewFileKey(user, ClientUser, filename, nil)
		err = filekey.Revoke()
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
	}
	// Get remaining file users
	fileUsers, err := GetFileUsers(ClientUser, filename)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	// Reshare file with remaining file users
	ShareFile(filename, fileUsers, false)
	fmt.Println("Successfully revoked file")
	os.Exit(0)
}
