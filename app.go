 package main

 import (
  "bufio"
  "encoding/base64"
  "fmt"
  "net/http"
  "os"
   "bytes"
   "io/ioutil"
 )

 func main(){
  imgFile, err := os.Open(os.Args[1]) // a QR code image

  if err != nil {
     fmt.Println(err)
     os.Exit(1)
  }

  defer imgFile.Close()

  // create a new buffer base on file size
  fInfo, _ := imgFile.Stat()
  var size int64 = fInfo.Size()
  buf := make([]byte, size)

  // read file content into buffer
  fReader := bufio.NewReader(imgFile)
  fReader.Read(buf)
  
  fmt.Println(len(buf))

  // convert the buffer bytes to base64 string - use buf.Bytes() for new image
  imgBase64Str := base64.StdEncoding.EncodeToString(buf)

//   fmt.Println(imgBase64Str)

   jsonStr := []byte(fmt.Sprintf(`{"device":"countcamera1", "readings":[{"name":"cameraeiamge", "value":"%s"}]}`,imgBase64Str))

//   fmt.Println(len(jsonStr))
   response, err := http.Post("http://localhost:48080/api/v1/event", "application/json", bytes.NewBuffer(jsonStr))
   if err != nil {
       fmt.Printf("The HTTP request failed with error %s\n", err)
   } else {
       data, _ := ioutil.ReadAll(response.Body)
       fmt.Println(string(data))
   }
   fmt.Println("Terminating the application...")
 }
