package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
)

type ConfirmationPageData struct {
	ProductName string
	IpfsHash    string
}
func createProductHandler(c *gin.Context) {
	productName := c.PostForm("productName")

	fileHeader, err := c.FormFile("productImage")
	if err != nil {
		c.String(http.StatusBadRequest, "Error: Image file is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Error: Could not open uploaded file")
		return
	}
	defer file.Close()
	ipfsHash, err := uploadToIPFS(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error: Failed to upload to IPFS")
		return
	}
	log.Printf("Uploaded image for '%s' to IPFS. CID: %s", productName, ipfsHash)

	data := ConfirmationPageData{
		ProductName: productName,
		IpfsHash:    ipfsHash,
	}

	tmpl, err := template.ParseFiles("io.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error: could not load page template.")
		return
	}

	tmpl.Execute(c.Writer, data)
}
func uploadToIPFS(fileData io.Reader) (string, error) {
	sh := shell.NewShell("localhost:5001")

	cid, err := sh.Add(fileData)
	if err != nil {
		log.Printf("Error adding file to IPFS: %s", err)
		return "", err
	}

	return cid, nil
}